package cmd

import (
	"context"
	"encoding/csv"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/go-faster/errors"
	"github.com/nathanhack/sibyl/ent/ogent"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var useCSV, details, all, inactive bool
var startDateStr, endDateStr string

// showCmd represents the get command
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show will retrieve various values from Sibyl server",
	Long:  `Show will retrieve various values from Sibyl server`,
}

var showStocksCmd = &cobra.Command{
	Use:   "stocks",
	Short: "Lists all current stocks(entities) from the Sibyl server",
	Long:  `Lists all current stocks(entities) from the Sibyl server`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := cobra.NoArgs(cmd, args); err != nil {
			return err
		}

		address, err := cmd.Flags().GetString("serverAddress")
		if err != nil {
			return fmt.Errorf("Could not get server address from passed in arguments: %v\n", err)
		}

		client, err := ogent.NewClient(address)
		if err != nil {
			return err
		}

		list, err := getAllEntities(context.Background(), client)
		if err != nil {
			return err
		}

		var data [][]string
		sort.Slice(list, func(i, j int) bool {
			if list[i].Active == list[j].Active {
				return strings.Compare(list[i].Ticker, list[j].Ticker) < 0
			}

			if list[i].Active {
				return true
			}
			return false
		})

		for _, v := range list {
			stockRow := make([]string, 0)
			if v.Active && !all && inactive {
				continue
			}

			if all {
				stockRow = append(stockRow, fmt.Sprint(v.Active))
			}

			if details {
				stockRow = append(stockRow, v.Ticker)
				stockRow = append(stockRow, v.Name)
				stockRow = append(stockRow, v.Description)
				// stockRow = append(stockRow, v.Delisted.String())
				stockRow = append(stockRow, v.ListDate.String())
			} else {
				stockRow = append(stockRow, v.Ticker)
				stockRow = append(stockRow, v.Name)

			}
			data = append(data, stockRow)
		}

		headers := []string{"Ticker", "Name"}

		if details {
			headers = []string{
				"Ticker",
				"Name",
				"Description",
				// "Delisted",
				"List",
			}
			if all {
				headers = append([]string{"Active"}, headers...)
			}
		}
		tableStr, err := makeTable(headers, data, true, useCSV)
		if err != nil {
			return err
		}

		fmt.Println(tableStr)

		return nil
	},
}

func makeTable(headers []string, rows [][]string, showHeaders, useCSV bool) (string, error) {
	tableString := &strings.Builder{}

	if useCSV {
		w := csv.NewWriter(tableString)
		if err := w.Write(headers); err != nil {
			return "", fmt.Errorf("error writing data in csv format: %v", err)
		}
		for _, row := range rows {
			if err := w.Write(row); err != nil {
				return "", fmt.Errorf("error writing data in csv format: %v", err)
			}
		}
		w.Flush()
	} else {
		table := tablewriter.NewWriter(tableString)
		if showHeaders {
			table.SetHeader(headers)
		}

		table.SetAutoWrapText(false)
		table.AppendBulk(rows)
		table.Render()
	}
	return tableString.String(), nil
}

func getTimestamp(str string, formatStr string, defaultValue time.Time) (time.Time, error) {
	//we have two cases either it's
	// a integer or it's a date/datetime string
	// it could also be blank meaning use the defaultValue
	if str == "" {
		return defaultValue, nil
	}

	if num, err := strconv.ParseInt(str, 10, 64); err != nil {
		if datetime, err := time.ParseInLocation(formatStr, str, time.Local); err != nil {
			return defaultValue, fmt.Errorf("problem decoding time information from input: \"%v\", expected an integer or formatting: %v", str, formatStr)
		} else {
			return datetime, nil
		}
	} else {
		return time.Unix(num, 0), nil
	}

}

var showDataSourcesCmd = &cobra.Command{
	Use:   "datasources",
	Short: "Shows the available DataSources",
	Long:  `Shows the available DataSources`,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := ogent.NewClient(serverAddress)
		if err != nil {
			return err
		}
		listDataSource, err := getAllDataSources(context.Background(), client)
		if err != nil {
			return err
		}

		sort.Slice(listDataSource, func(i, j int) bool {
			return listDataSource[i].Name < listDataSource[j].Name
		})

		data := make([][]string, 0)
		for _, ds := range listDataSource {
			data = append(data, []string{strconv.Itoa(ds.ID), ds.Name})
		}
		headers := []string{"ID", "Name"}
		tableStr, err := makeTable(headers, data, true, false)
		if err != nil {
			return err
		}
		fmt.Println(tableStr)

		return nil
	},
}

var showBarsCmd = &cobra.Command{
	Use:   "bars",
	Short: "Shows the bar history for a particular stock",
	Long:  `Shows the bar history for a particular stock`,
}

var showBarDailyCmd = &cobra.Command{
	Use:   "daily TICKER [DATA_SOURCE_ID]",
	Args:  cobra.RangeArgs(1, 2),
	Short: "Shows the daily bar history",
	Long:  `Shows the daily bar history.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		startDate, err := getTimestamp(startDateStr, "2006-01-02", time.Time{})
		if err != nil {
			return fmt.Errorf("could not get starting timestamp from passed in arguments: %v", err)
		}

		endDate, err := getTimestamp(endDateStr, "2006-01-02", time.Now().Local())
		if err != nil {
			return fmt.Errorf("could not get ending timestamp from passed in arguments: %v", err)
		}

		datasourceID := 0
		if len(args) >= 2 {
			datasourceID, err = strconv.Atoi(args[1])
			if err != nil {
				return fmt.Errorf("could not parse DATA_SOURCE_ID found '%v' and had error:%v", args[1], err)
			}
		}

		client, err := ogent.NewClient(serverAddress)
		if err != nil {
			return err
		}
		ctx := context.Background()
		ticker := strings.ToUpper(args[0])

		stock, err := getEntity(ctx, ticker, client)
		if err != nil {
			return errors.Wrapf(err, "entity %v not found", ticker)
		}

		//get all the intervals and filter out anything not == daily intervals
		listEntityIntervals, err := getAllIntervals(ctx, client, stock.ID, func(item ogent.EntityIntervalsList) bool {
			return item.Active && item.Interval == ogent.EntityIntervalsListIntervalDaily
		})
		if err != nil {
			return err
		}

		if len(listEntityIntervals) == 1 {
			id := listEntityIntervals[0].DataSourceID
			if datasourceID == 0 {
				datasourceID = id
			} else if datasourceID != id {
				return fmt.Errorf("the passed in Data Source ID(%v) did not match the only Data Source ID(%v) available.", datasourceID, id)
			}

		} else if len(listEntityIntervals) > 1 {
			// if there is more than one then we get data sources to check if the
			// datasource passed in is valid or show them as options if one was not passed in
			_, err := getDataSource(ctx, client, datasourceID)
			if err != nil {
				return err
			}
		}

		intervalID := 0
		for _, ei := range listEntityIntervals {
			if ei.DataSourceID == datasourceID {
				intervalID = ei.ID
			}
		}

		if intervalID == 0 {
			return fmt.Errorf("There doesn't appear to be an interval for 'Daily' data.")
		}

		// with the interval and data source we find the barstimeranges
		return getBarsAndPrint(ctx, client, intervalID, startDate, endDate, false)
	},
}

var showBar1minCmd = &cobra.Command{
	Use:   "1min TICKER [DATA_SOURCE_ID]",
	Args:  cobra.RangeArgs(1, 2),
	Short: "Shows the 1min bar history",
	Long:  `Shows the 1min bar history.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		startDate, err := getTimestamp(startDateStr, "2006-01-02 15:04", time.Time{})
		if err != nil {
			return fmt.Errorf("could not get starting timestamp from passed in arguments: %v", err)
		}

		endDate, err := getTimestamp(endDateStr, "2006-01-02 15:04", time.Now().Local())
		if err != nil {
			return fmt.Errorf("could not get ending timestamp from passed in arguments: %v", err)
		}

		datasourceID := 0
		if len(args) >= 2 {
			datasourceID, err = strconv.Atoi(args[1])
			if err != nil {
				return fmt.Errorf("could not parse DATA_SOURCE_ID found '%v' and had error:%v", args[1], err)
			}
		}

		client, err := ogent.NewClient(serverAddress)
		if err != nil {
			return err
		}
		ctx := context.Background()
		ticker := strings.ToUpper(args[0])

		stock, err := getEntity(ctx, ticker, client)
		if err != nil {
			return errors.Wrapf(err, "entity %v not found", ticker)
		}

		//get all the intervals and filter out anything not == daily intervals
		listEntityIntervals, err := getAllIntervals(ctx, client, stock.ID, func(item ogent.EntityIntervalsList) bool {
			return item.Interval == ogent.EntityIntervalsListInterval1min
		})
		if err != nil {
			return err
		}

		if len(listEntityIntervals) == 1 {
			id := listEntityIntervals[0].DataSourceID
			if datasourceID == 0 {
				datasourceID = id
			} else if datasourceID != id {
				return fmt.Errorf("the passed in Data Source ID(%v) did not match the only Data Source ID(%v) available.", datasourceID, id)
			}

		} else if len(listEntityIntervals) > 1 {
			// if there is more than one then we get data sources to check if the
			// datasource passed in is valid or show them as options if one was not passed in
			_, err := getDataSource(ctx, client, datasourceID)
			if err != nil {
				return err
			}
		}

		intervalID := 0
		for _, ei := range listEntityIntervals {
			if ei.DataSourceID == datasourceID {
				intervalID = ei.ID
			}
		}

		if intervalID == 0 {
			return fmt.Errorf("There doesn't appear to be an interval for '1min' data.")
		}

		// with the interval and data source we find the barstimeranges
		return getBarsAndPrint(ctx, client, intervalID, startDate, endDate, true)
	},
}

func getBarsAndPrint(ctx context.Context, client *ogent.Client, intervalID int, startDate, endDate time.Time, addTime bool) error {
	listBarTimeRanges, err := getAllBarTimeRanges(ctx, client, intervalID, func(ibl ogent.IntervalBarsList) bool {
		if endDate.Before(ibl.Start) ||
			ibl.End.Before(startDate) ||
			ibl.Status != ogent.IntervalBarsListStatusConsolidated {
			return false
		}
		return true
	})
	if err != nil {
		return err
	}

	var data [][]string
	records := make([]ogent.BarGroupRecordsList, 0)
	for _, bar := range listBarTimeRanges {
		// then with those we get the BarGroups
		listBarTimeRangeGroups, err := getAllBarGroups(ctx, client, bar.ID, func(btrgl ogent.BarTimeRangeGroupsList) bool {
			if endDate.Before(btrgl.First) || btrgl.Last.Before(startDate) {
				return false
			}
			return true
		})
		if err != nil {
			return err
		}

		// then get those records for the remaining BarGroups
		for _, group := range listBarTimeRangeGroups {
			listBarGroupRecords, err := getAllBarRecords(ctx, client, group.ID, func(bgrl ogent.BarGroupRecordsList) bool {
				if endDate.Before(bgrl.Timestamp) || bgrl.Timestamp.Before(startDate) {
					return false
				}
				return true
			})
			if err != nil {
				return err
			}
			records = append(records, listBarGroupRecords...)
		}
	}

	sort.Slice(records, func(i, j int) bool {
		return records[i].Timestamp.Before(records[j].Timestamp)
	})

	formatStr := "2006-01-02"
	if addTime {
		formatStr = "2006-01-02 15:04"
	}
	for _, v := range records {
		data = append(data, []string{
			strconv.Itoa(len(data)),
			v.Timestamp.Local().Format(formatStr),
			strconv.FormatFloat(v.Close, 'f', 2, 64),
			strconv.FormatFloat(v.High, 'f', 2, 64),
			strconv.FormatFloat(v.Low, 'f', 2, 64),
			strconv.FormatFloat(v.Open, 'f', 2, 64),
			strconv.Itoa(int(v.Volume)),
			strconv.FormatInt(int64(v.Transactions), 10),
		})
	}

	// now we process the results in to a table
	headers := []string{"Num", "Date Time", "Close", "High", "Low", "Open", "Volume", "Transactions"}

	tableStr, err := makeTable(headers, data, true, useCSV)
	if err != nil {
		return err
	}
	fmt.Println(tableStr)
	return nil
}

var showSplitsCmd = &cobra.Command{
	Use:   "splits",
	Short: "Shows the split history for a particular stock",
	Long:  `Shows the split history for a particular stock`,
}

var showDividendsCmd = &cobra.Command{
	Use:   "dividends",
	Short: "Shows the dividend history for a particular stock",
	Long:  `Shows the dividend history for a particular stock`,
}

func init() {
	rootCmd.AddCommand(showCmd)

	showCmd.AddCommand(showStocksCmd, showDataSourcesCmd)
	showCmd.PersistentFlags().BoolVarP(&useCSV, "csv", "", false, "makes output CSV formatted")

	showStocksCmd.Flags().BoolVarP(&details, "details", "d", false, "show all details about the stocks")
	showStocksCmd.Flags().BoolVar(&all, "all", false, "show all stocks active and inactive")
	showStocksCmd.Flags().BoolVarP(&inactive, "inactive", "i", false, "show ONLY inactive stocks")

	showCmd.AddCommand(showBarsCmd)
	showBarsCmd.AddCommand(showBarDailyCmd, showBar1minCmd)

	showBarDailyCmd.Flags().StringVarP(&startDateStr, "start", "s", "", "starting Date: an integer timestamp or date formatted like '2006-01-02'. If omitted assumes 0 (beginning of time)")
	showBarDailyCmd.Flags().StringVarP(&endDateStr, "end", "e", "", "ending Date: an integer timestamp or date formatted like '2006-01-02'. If omitted assumes today")

	showBar1minCmd.Flags().StringVarP(&startDateStr, "start", "s", "", "starting Date and Time: an integer timestamp or date formatted like '2006-01-02 15:23'. If omitted assumes 0 (beginning of time)")
	showBar1minCmd.Flags().StringVarP(&endDateStr, "end", "e", "", "ending Date and Time: an integer timestamp or date formatted like '2006-01-02 15:23'. If omitted assumes now")

}
