package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"github.com/flytam/filenamify"
	"github.com/nathanhack/sibyl/ent/ogent"
	"github.com/nathanhack/threadpool"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var dumpCmd = &cobra.Command{
	Use:   "dump",
	Short: "dump all data to a directory",
	Long:  `dump all data to a directory`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if _, err := os.Stat(dumpOutputDir); os.IsExist(err) {
			empty, err := directoryIsEmpty(dumpOutputDir)
			if err != nil {
				return errors.Wrapf(err, "checking directory")
			}
			if !empty {
				return errors.New("the output directory must be empty")
			}
		}

		err := os.MkdirAll(dumpOutputDir, os.ModePerm)
		if err != nil {
			return err
		}

		client, err := ogent.NewClient(serverAddress)
		if err != nil {
			return err
		}
		ctx := context.Background()

		allDataSources, err := getAllDataSources(ctx, client, func(dsl ogent.DataSourceList) bool { return true })
		if err != nil {
			return err
		}

		err = toJsonFile(allDataSources, filepath.Join(dumpOutputDir, "dataSources.json"))
		if err != nil {
			return err
		}

		_, err = getDataSource(ctx, client, dumpDataSourceId)
		if err != nil {
			return err
		}

		logrus.Info("Getting all entities")
		entities, err := getAllEntities(ctx, client)
		if err != nil {
			return err
		}

		sort.Slice(entities, func(i, j int) bool {
			return entities[i].Ticker < entities[j].Ticker
		})

		err = toJsonFile(entities, filepath.Join(dumpOutputDir, "entities.json"))
		if err != nil {
			return err
		}

		logrus.Info("Getting all market hours")
		marketDays, err := getAllMarketHours(ctx, client)
		if err != nil {
			return err
		}

		err = toJsonFile(marketDays, filepath.Join(dumpOutputDir, "markethours.json"))
		if err != nil {
			return err
		}

		pool := threadpool.New(ctx, 3, len(entities))
		errs := make([]error, 0)
		mux := sync.Mutex{}
		for _, entity := range entities {
			en := entity
			pool.Add(func() {
				err := dumpProcessStock(ctx, client, en, allDataSources, marketDays)
				mux.Lock()
				if err != nil {
					errs = append(errs, errors.Wrapf(err, "%v", en.Ticker))
					logrus.Errorf("%v:%v", en.Ticker, err)
				}

				logrus.Infof("Done: %v", en.Ticker)
				mux.Unlock()
			})

		}
		pool.Wait()

		defer func() {
			for _, err := range errs {
				logrus.Error(err)
			}
		}()

		return nil
	},
}

func toJsonFile(object interface{}, pathname string) error {
	bs, err := json.Marshal(object)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(pathname, bs, 0644)
	if err != nil {
		return err
	}
	return nil
}

func dumpGetBars(ctx context.Context, client *ogent.Client, entityId, datasourceId int, intervalValue ogent.EntityIntervalsListInterval) ([]ogent.BarGroupRecordsList, error) {
	interval1Min, err := getInterval(ctx, client, entityId, datasourceId, intervalValue)
	if err != nil {
		return nil, errors.Wrapf(err, "getInterval 1min")
	}

	barTimeRange1Min, err := getLargestBarTimeRange(ctx, client, interval1Min.ID)
	if err != nil {
		return nil, errors.Wrapf(err, "getLargestBarTimeRange 1min")
	}

	records1Min, err := getAllBarTimeRangeBarRecords(ctx, client, barTimeRange1Min.ID)
	if err != nil {
		return nil, errors.Wrapf(err, "getAllBarTimeRangeBarRecords 1min")
	}

	return records1Min, nil
}

func directoryIsEmpty(dir string) (bool, error) {
	f, err := os.Open(dir)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1) // Or f.Readdir(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err // Either not empty or error, suits both cases
}

func dumpProcessStock(ctx context.Context, client *ogent.Client, entity ogent.EntityList, datasources []ogent.DataSourceList, marketDays []ogent.MarketHoursList) error {
	extraPath, err := filenamify.Filenamify(fmt.Sprintf("%v_%v", entity.Ticker, entity.Name), filenamify.Options{})
	if err != nil {
		return err
	}

	stockDir := path.Join(dumpOutputDir, extraPath)
	err = os.MkdirAll(stockDir, os.ModePerm)
	if err != nil {
		return err
	}

	err = toJsonFile(&entity, path.Join(stockDir, "entity.json"))
	if err != nil {
		return err
	}

	logrus.Infof("Getting all dividends for %v", entity.Ticker)
	dividends, err := getAllDividends(ctx, client, entity.ID)
	if err != nil {
		return errors.Wrapf(err, "getAllDividends")
	}
	err = toJsonFile(dividends, filepath.Join(stockDir, "dividends.json"))
	if err != nil {
		return err
	}

	logrus.Infof("Getting all splits for %v", entity.Ticker)
	splits, err := getAllSplits(ctx, client, entity.ID)
	if err != nil {
		return err
	}
	err = toJsonFile(splits, filepath.Join(stockDir, "splits.json"))
	if err != nil {
		return err
	}

	for _, ds := range datasources {
		logrus.Infof("Getting all daily BarRecords for %v", entity.Ticker)
		recordsDaily, err := dumpGetBars(ctx, client, entity.ID, ds.ID, ogent.EntityIntervalsListIntervalDaily)
		if err != nil {
			return errors.Wrapf(err, "getAllBarTimeRangeBarRecords daily")
		}

		logrus.Infof("Getting all 1min BarRecords for %v", entity.Ticker)
		records1Min, err := dumpGetBars(ctx, client, entity.ID, ds.ID, ogent.EntityIntervalsListInterval1min)
		if err != nil {
			return errors.Wrapf(err, "getAllBarTimeRangeBarRecords 1min")
		}

		logrus.Infof("Processing %v", entity.Ticker)
		//sort all the data sources
		sort.Slice(marketDays, func(i, j int) bool {
			return marketDays[i].Date.Before(marketDays[i].Date)
		})

		sort.Slice(records1Min, func(i, j int) bool {
			return records1Min[i].Timestamp.Before(records1Min[j].Timestamp)
		})

		sort.Slice(recordsDaily, func(i, j int) bool {
			return recordsDaily[i].Timestamp.Before(recordsDaily[j].Timestamp)
		})

		fullDir := path.Join(stockDir, fmt.Sprint(ds.ID))
		os.MkdirAll(fullDir, os.ModePerm)
		if err != nil {
			return err
		}

		err = toJsonFile(recordsDaily, filepath.Join(fullDir, "daily.json"))
		if err != nil {
			return err
		}

		//starting with the entities List date and incrementing until the last daily record
		// we create dataset output
		tomorrow := time.Now().Local().Truncate(24*time.Hour).AddDate(0, 0, 1)
		for currentDate := time.Date(entity.ListDate.Year(), entity.ListDate.Month(), entity.ListDate.Day(), 0, 0, 0, 0, time.Local); currentDate.Before(tomorrow); currentDate = currentDate.AddDate(0, 0, 1) {

			if currentDate.Weekday() == time.Saturday || currentDate.Weekday() == time.Sunday {
				continue
			}

			marketHoursDayOfIndex := sort.Search(len(marketDays), func(i int) bool {
				return marketDays[i].Date.Equal(currentDate) || marketDays[i].Date.After(currentDate)
			})

			if marketHoursDayOfIndex == len(marketDays) ||
				!marketDays[marketHoursDayOfIndex].Date.Equal(currentDate) {
				continue
			}

			records, good := dumpSliceRecordsDayOf(currentDate, records1Min)
			if !good {
				continue
			}

			err = toJsonFile(records, filepath.Join(fullDir, fmt.Sprintf("1min_%v.json", currentDate.Format("2006_01_02"))))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// dumpSliceRecordsDayOf returns the a slice of records with only from the date
func dumpSliceRecordsDayOf(date time.Time, records []ogent.BarGroupRecordsList) ([]ogent.BarGroupRecordsList, bool) {
	// so they should have already been sorted

	recordsIndex := sort.Search(len(records), func(i int) bool {
		bgrl := records[i]
		t := time.Date(bgrl.Timestamp.Year(), bgrl.Timestamp.Month(), bgrl.Timestamp.Day(), 0, 0, 0, 0, time.Local)
		return t.Equal(date) || date.Before(t)
	})

	if recordsIndex == len(records) {
		return nil, false
	}

	// for we need to make sure we have values for each day the market was open
	nextDay := date.AddDate(0, 0, 1)
	nextRecordsIndex := sort.Search(len(records), func(i int) bool {
		bgrl := records[i]
		t := time.Date(bgrl.Timestamp.Year(), bgrl.Timestamp.Month(), bgrl.Timestamp.Day(), 0, 0, 0, 0, time.Local)
		return t.Equal(nextDay) || nextDay.Before(t)
	})

	return records[recordsIndex:nextRecordsIndex], true
}

var dumpOutputDir string
var dumpDataSourceId int

func init() {
	rootCmd.AddCommand(dumpCmd)

	dumpCmd.Flags().StringVarP(&dumpOutputDir, "output", "o", "./output", "directory to store output")
	dumpCmd.MarkPersistentFlagDirname("output")

	dumpCmd.PersistentFlags().IntVarP(&dumpDataSourceId, "datasource", "d", 0, "The DataSourceID to use (defaults to 0 -- all data sources)")
	dumpCmd.MarkPersistentFlagRequired("datasource")
}
