package cmd

import (
	"context"
	"fmt"
	"sort"

	"github.com/go-faster/errors"

	"github.com/nathanhack/sibyl/ent/ogent"
	"github.com/ogen-go/ogen/ogenerrors"
)

func getDataSource(ctx context.Context, client *ogent.Client, id int) (*ogent.DataSourceList, error) {
	listDataSource, err := getAllDataSources(ctx, client, func(dsl ogent.DataSourceList) bool {
		return dsl.ID == id
	})
	if err != nil {
		return nil, err
	}

	if len(listDataSource) == 0 {
		return nil, fmt.Errorf("the DataSource ID does not exist (use 'show datasources')")
	}

	return &(listDataSource)[0], nil
}

// getAllDataSources gets all the DataSources, the keepFilter returns true if you want to keep it
func getAllDataSources(ctx context.Context, client *ogent.Client, keepFilter ...func(ogent.DataSourceList) bool) ([]ogent.DataSourceList, error) {
	results := make([]ogent.DataSourceList, 0)
	for i := 1; i < 1<<24; i++ {
		listDataSourceRes, err := client.ListDataSource(ctx, ogent.ListDataSourceParams{
			Page:         ogent.NewOptInt(i),
			ItemsPerPage: ogent.NewOptInt(255),
		})
		if err != nil {
			return nil, err
		}

		listDataSource, ok := listDataSourceRes.(*ogent.ListDataSourceOKApplicationJSON)
		if !ok {
			return nil, fmt.Errorf("expected type %T found: %T", &ogent.ListDataSourceOKApplicationJSON{}, listDataSourceRes)
		}

		if len(*listDataSource) == 0 {
			break
		}

		//run the keep filter if true we KEEP
		for _, filter := range keepFilter {
			for i := 0; i < len(*listDataSource); {
				if !filter((*listDataSource)[i]) {
					lastIndex := len(*listDataSource) - 1
					(*listDataSource)[i] = (*listDataSource)[lastIndex]
					(*listDataSource) = (*listDataSource)[:lastIndex]
					continue
				}
				i++
			}
		}

		results = append(results, *listDataSource...)

	}

	return results, nil
}

func getInterval(ctx context.Context, client *ogent.Client, stockId, dataSourceId int, targetInterval ogent.EntityIntervalsListInterval) (*ogent.EntityIntervalsList, error) {
	intervals, err := getAllIntervals(ctx, client, stockId, func(eil ogent.EntityIntervalsList) bool {
		return eil.DataSourceID == dataSourceId && eil.Interval == targetInterval
	})
	if err != nil {
		return nil, err
	}

	if len(intervals) != 1 {
		return nil, errors.Errorf("expected to find just one interval(%v) for stockId:%v with DataSourceId:%v", targetInterval, stockId, dataSourceId)
	}
	return &intervals[0], nil
}

// getAllIntervals gets all intervals for the ticker, the keepFilter returns true if you want to keep it
func getAllIntervals(ctx context.Context, client *ogent.Client, stockId int, keepFilter ...func(ogent.EntityIntervalsList) bool) ([]ogent.EntityIntervalsList, error) {
	results := make([]ogent.EntityIntervalsList, 0)
	for i := 1; i < 1<<24; i++ {
		listEntityIntervalsRes, err := client.ListEntityIntervals(ctx, ogent.ListEntityIntervalsParams{
			ID:           stockId,
			Page:         ogent.NewOptInt(i),
			ItemsPerPage: ogent.NewOptInt(255),
		})
		if err != nil {
			return nil, unwrapDecodeBodyError(err)
		}

		listEntityIntervals, ok := listEntityIntervalsRes.(*ogent.ListEntityIntervalsOKApplicationJSON)
		if !ok {
			return nil, fmt.Errorf("expected type %T found: %T", &ogent.ListEntityIntervalsOKApplicationJSON{}, listEntityIntervalsRes)
		}

		if len(*listEntityIntervals) == 0 {
			break
		}

		//run the keep filter if true we KEEP
		for _, filter := range keepFilter {
			for i := 0; i < len(*listEntityIntervals); {
				if !filter((*listEntityIntervals)[i]) {
					lastIndex := len(*listEntityIntervals) - 1
					(*listEntityIntervals)[i] = (*listEntityIntervals)[lastIndex]
					(*listEntityIntervals) = (*listEntityIntervals)[:lastIndex]
					continue
				}
				i++
			}
		}

		results = append(results, *listEntityIntervals...)
	}

	return results, nil
}

func getEntity(ctx context.Context, ticker string, client *ogent.Client) (*ogent.EntityList, error) {
	stocks, err := getAllEntities(ctx, client)
	if err != nil {
		return nil, err
	}

	for _, stock := range stocks {
		if stock.Ticker == ticker {
			tmp := stock
			return &tmp, nil
		}
	}
	return nil, fmt.Errorf("ticker: %s not found", ticker)
}

func getAllEntities(ctx context.Context, client *ogent.Client) ([]ogent.EntityList, error) {
	results := make([]ogent.EntityList, 0)
	for i := 1; i < 1<<24; i++ {
		res, err := client.ListEntity(ctx, ogent.ListEntityParams{
			Page:         ogent.NewOptInt(i),
			ItemsPerPage: ogent.NewOptInt(255),
		})
		if err != nil {
			return nil, unwrapDecodeBodyError(err)
		}

		reslist, ok := res.(*ogent.ListEntityOKApplicationJSON)
		if !ok {
			return nil, fmt.Errorf("expected type %T found: %T", &ogent.ListEntityOKApplicationJSON{}, res)
		}

		if len(*reslist) == 0 {
			break
		}

		results = append(results, *reslist...)
	}
	return results, nil
}

func unwrapDecodeBodyError(err error) error {
	tmp := err
	for i := 0; i < 3; i++ {
		if tmp != nil {
			if v, ok := err.(*ogenerrors.DecodeBodyError); ok {
				return fmt.Errorf("%s", v.Body)
			}
			tmp = errors.Unwrap(tmp)
		}
	}
	return err
}

func getAllBarTimeRanges(ctx context.Context, client *ogent.Client, intervalId int, keepFilter ...func(ogent.IntervalBarsList) bool) ([]ogent.IntervalBarsList, error) {
	results := make([]ogent.IntervalBarsList, 0)
	for i := 1; i < 1<<24; i++ {
		listIntervalBarsRes, err := client.ListIntervalBars(ctx, ogent.ListIntervalBarsParams{
			ID:           intervalId,
			Page:         ogent.NewOptInt(i),
			ItemsPerPage: ogent.NewOptInt(255),
		})
		if err != nil {
			return nil, err
		}

		listBarGroups, ok := listIntervalBarsRes.(*ogent.ListIntervalBarsOKApplicationJSON)
		if !ok {
			return nil, fmt.Errorf("expected type %T found: %T", &ogent.ListIntervalBarsOKApplicationJSON{}, listIntervalBarsRes)
		}

		if len(*listBarGroups) == 0 {
			break
		}

		//run the keep filter if true we KEEP
		for _, filter := range keepFilter {
			for i := 0; i < len(*listBarGroups); {
				if !filter((*listBarGroups)[i]) {
					lastIndex := len(*listBarGroups) - 1
					(*listBarGroups)[i] = (*listBarGroups)[lastIndex]
					(*listBarGroups) = (*listBarGroups)[:lastIndex]
					continue
				}
				i++
			}
		}

		results = append(results, *listBarGroups...)
	}

	return results, nil
}

func getAllBarGroups(ctx context.Context, client *ogent.Client, barId int, keepFilter ...func(ogent.BarTimeRangeGroupsList) bool) ([]ogent.BarTimeRangeGroupsList, error) {
	results := make([]ogent.BarTimeRangeGroupsList, 0)
	for i := 1; i < 1<<24; i++ {
		listBarTimeRangeGroupsRes, err := client.ListBarTimeRangeGroups(ctx, ogent.ListBarTimeRangeGroupsParams{
			ID:           barId,
			Page:         ogent.NewOptInt(i),
			ItemsPerPage: ogent.NewOptInt(255),
		})
		if err != nil {
			return nil, err
		}

		listBarGroups, ok := listBarTimeRangeGroupsRes.(*ogent.ListBarTimeRangeGroupsOKApplicationJSON)
		if !ok {
			return nil, fmt.Errorf("expected type %T found: %T", &ogent.ListBarTimeRangeGroupsOKApplicationJSON{}, listBarTimeRangeGroupsRes)
		}

		if len(*listBarGroups) == 0 {
			break
		}

		//run the keep filter if true we KEEP
		for _, filter := range keepFilter {
			for i := 0; i < len(*listBarGroups); {
				if !filter((*listBarGroups)[i]) {
					lastIndex := len(*listBarGroups) - 1
					(*listBarGroups)[i] = (*listBarGroups)[lastIndex]
					(*listBarGroups) = (*listBarGroups)[:lastIndex]
					continue
				}
				i++
			}
		}

		results = append(results, *listBarGroups...)

	}

	return results, nil
}

func getAllBarRecords(ctx context.Context, client *ogent.Client, groupId int, keepFilter ...func(ogent.BarGroupRecordsList) bool) ([]ogent.BarGroupRecordsList, error) {
	results := make([]ogent.BarGroupRecordsList, 0)
	for i := 1; i < 1<<24; i++ {
		listBarGroupRecordsRes, err := client.ListBarGroupRecords(ctx, ogent.ListBarGroupRecordsParams{
			ID:           groupId,
			Page:         ogent.NewOptInt(i),
			ItemsPerPage: ogent.NewOptInt(255),
		})
		if err != nil {
			return nil, err
		}

		listBarRecords, ok := listBarGroupRecordsRes.(*ogent.ListBarGroupRecordsOKApplicationJSON)
		if !ok {
			return nil, fmt.Errorf("expected type %T found: %T", &ogent.ListBarGroupRecordsOKApplicationJSON{}, listBarGroupRecordsRes)
		}

		if len(*listBarRecords) == 0 {
			break
		}

		//run the keep filter if true we KEEP
		for _, filter := range keepFilter {
			for i := 0; i < len(*listBarRecords); {
				if !filter((*listBarRecords)[i]) {
					lastIndex := len(*listBarRecords) - 1
					(*listBarRecords)[i] = (*listBarRecords)[lastIndex]
					(*listBarRecords) = (*listBarRecords)[:lastIndex]
					continue
				}
				i++
			}
		}

		results = append(results, *listBarRecords...)
	}

	return results, nil
}

func getLargestBarTimeRange(ctx context.Context, client *ogent.Client, intervalId int) (*ogent.IntervalBarsList, error) {

	ranges, err := getAllBarTimeRanges(ctx, client, intervalId, func(ibl ogent.IntervalBarsList) bool {
		return ibl.Status == ogent.IntervalBarsListStatusConsolidated
	})
	if err != nil {
		return nil, err
	}

	if len(ranges) == 0 {
		return nil, fmt.Errorf("expected at least one BarTimeRange")
	}

	// sort for the longest, oldest start, oldest id
	sort.Slice(ranges, func(i, j int) bool {
		idiff := ranges[i].End.Sub(ranges[i].Start)
		jdiff := ranges[j].End.Sub(ranges[j].Start)

		if idiff == jdiff {
			if ranges[i].Start.Equal(ranges[j].Start) {
				return ranges[i].ID < ranges[j].ID
			}
			return ranges[i].Start.Before(ranges[j].Start)
		}

		return idiff > jdiff
	})

	return &ranges[0], nil
}

func getAllBarTimeRangeBarRecords(ctx context.Context, client *ogent.Client, barTimeRangeId int) ([]ogent.BarGroupRecordsList, error) {
	results := make([]ogent.BarGroupRecordsList, 0)
	groups, err := getAllBarGroups(ctx, client, barTimeRangeId, func(btrgl ogent.BarTimeRangeGroupsList) bool { return true })
	if err != nil {
		return nil, err
	}

	for _, group := range groups {
		records, err := getAllBarRecords(ctx, client, group.ID, func(bgrl ogent.BarGroupRecordsList) bool { return true })
		if err != nil {
			return nil, err
		}

		results = append(results, records...)
	}

	return results, nil
}

func addOrModifyInterval(ctx context.Context, ticker string, interval ogent.EntityIntervalsListInterval, dataSourceID int, client *ogent.Client, activeState bool) error {
	// first we check the datasource we need to see if it exists
	_, err := getDataSource(ctx, client, dataSourceID)
	if err != nil {
		return errors.Wrap(err, "datasource not found")
	}

	stock, err := getEntity(ctx, ticker, client)
	if err != nil {
		return errors.Wrapf(err, "entity %v not found", ticker)
	}

	// now we determine if the interval exists if not we create it
	intervals, err := getAllIntervals(ctx, client, stock.ID, func(eil ogent.EntityIntervalsList) bool {
		return eil.DataSourceID == dataSourceID && eil.Interval == interval
	})
	if err != nil {
		return errors.Wrap(err, "finding intervals")
	}

	if len(intervals) == 0 {
		//doesn't exist so we make it
		_, err := client.CreateInterval(ctx, &ogent.CreateIntervalReq{
			Active:       activeState,
			Interval:     ogent.CreateIntervalReqInterval(interval),
			StockID:      stock.ID,
			DataSourceID: dataSourceID,
			DataSource:   dataSourceID,
			Stock:        stock.ID,
		})
		if err != nil {
			return errors.Wrap(err, "creating interval")
		}
		return nil
	}

	_, err = client.UpdateInterval(ctx, &ogent.UpdateIntervalReq{
		Active: ogent.NewOptBool(activeState),
	}, ogent.UpdateIntervalParams{
		ID: intervals[0].ID,
	})

	if err != nil {
		return errors.Wrap(err, "updating interval")
	}

	return nil
}

func getAllMarketHours(ctx context.Context, client *ogent.Client) ([]ogent.MarketHoursList, error) {
	results := make([]ogent.MarketHoursList, 0)
	for i := 1; i < 1<<24; i++ {
		recordsRes, err := client.ListMarketHours(ctx, ogent.ListMarketHoursParams{
			Page:         ogent.NewOptInt(i),
			ItemsPerPage: ogent.NewOptInt(255),
		})
		if err != nil {
			return nil, err
		}

		records, ok := recordsRes.(*ogent.ListMarketHoursOKApplicationJSON)
		if !ok {
			return nil, fmt.Errorf("expected type %T found: %T", &ogent.ListMarketHoursOKApplicationJSON{}, recordsRes)
		}

		if len(*records) == 0 {
			break
		}

		results = append(results, *records...)
	}
	return results, nil
}

func getAllDividends(ctx context.Context, client *ogent.Client, stockId int) ([]ogent.EntityDividendsList, error) {
	results := make([]ogent.EntityDividendsList, 0)
	for i := 1; i < 1<<24; i++ {
		dividendsRes, err := client.ListEntityDividends(ctx, ogent.ListEntityDividendsParams{
			ID:           stockId,
			Page:         ogent.NewOptInt(i),
			ItemsPerPage: ogent.NewOptInt(255),
		})
		if err != nil {
			return nil, unwrapDecodeBodyError(err)
		}

		dividends, ok := dividendsRes.(*ogent.ListEntityDividendsOKApplicationJSON)
		if !ok {
			return nil, fmt.Errorf("expected type %T found: %T", &ogent.ListEntityDividendsOKApplicationJSON{}, dividendsRes)
		}

		if len(*dividends) == 0 {
			break
		}

		results = append(results, *dividends...)
	}

	return results, nil
}

func getAllSplits(ctx context.Context, client *ogent.Client, stockId int) ([]ogent.EntitySplitsList, error) {
	results := make([]ogent.EntitySplitsList, 0)
	for i := 1; i < 1<<24; i++ {
		splitsRes, err := client.ListEntitySplits(ctx, ogent.ListEntitySplitsParams{
			ID:           stockId,
			Page:         ogent.NewOptInt(i),
			ItemsPerPage: ogent.NewOptInt(255),
		})
		if err != nil {
			return nil, unwrapDecodeBodyError(err)
		}

		splits, ok := splitsRes.(*ogent.ListEntitySplitsOKApplicationJSON)
		if !ok {
			return nil, fmt.Errorf("expected type %T found: %T", &ogent.ListEntitySplitsOKApplicationJSON{}, splitsRes)
		}

		if len(*splits) == 0 {
			break
		}

		results = append(results, *splits...)
	}

	return results, nil
}
