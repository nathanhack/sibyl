package barrequester

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/go-faster/errors"
	"github.com/nathanhack/sibyl/cmd/server/cmd/internal"
	"github.com/nathanhack/sibyl/ent"
	"github.com/nathanhack/sibyl/ent/bargroup"
	"github.com/nathanhack/sibyl/ent/barrecord"
	"github.com/nathanhack/sibyl/ent/bartimerange"
	"github.com/nathanhack/sibyl/ent/interval"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
)

func Scrubber(ctx context.Context, client *ent.Client, wg *sync.WaitGroup) {
	logrus.Info("Bars.Scrubber: Running")
	startupTimer := time.NewTimer(5 * time.Second)
	wg.Add(1)
	defer func() {
		logrus.Info("Bars.Scrubber: Stopped")
		wg.Done()
	}()

	// CLEAN
	// look for BarTimeRange with miss-match of BarGroups
	// look for BarGroups with a miss-match of bar records
	// look for BarGroups that have miss-match timerange (shouldn't exist but)

	//CONSOLIDATE
	// look for BarTimeRange the have overlap
	// look for BarGroups that have over lap

	for {
		select {
		case <-ctx.Done():
			return
		case <-startupTimer.C:
		case <-internal.AtMidnight():
		}

		//remove bad BarTimeRanges that are older than 12 hours
		if err := removeBadBarTimeRangesOlderThan(ctx, client, time.Now().Add(-12*time.Hour)); err != nil {
			logrus.Errorf("Bars.Scrubber: %v", err)
		}

		// scrub intervals that are inactive and have unconsolidated BarTimeRanges
		if err := scrub(ctx, client); err != nil {
			logrus.Errorf("Bars.Scrubber: %v", err)
		}
	}
}

// removeBadBarTimeRangesOlderThan deletes BarTimeRanges (and associated BarGroups and BarRecords) that are StatusPending
// and have been around longer than olderThan hours
func removeBadBarTimeRangesOlderThan(ctx context.Context, client *ent.Client, olderThan time.Time) error {
	bars, err := client.BarTimeRange.Query().Where(
		bartimerange.UpdateTimeLTE(olderThan),
		bartimerange.StatusIn(bartimerange.StatusPending),
	).All(ctx)
	if err != nil {
		return errors.Wrap(err, "removeBadBarTimeRangesOlder")
	}

	// auto delete any that are still pending
	for _, bar := range bars {
		err := deleteBarTimeRange(ctx, client, bar)
		if err != nil {
			logrus.Errorf("Bars.Scrubber: BarTimeRange(%v) pending after %v hrs and can not be deleted: %v", bar.ID, olderThan, err)
		}
	}
	return nil
}

// scrub cleans and consolidates intervals that are inactive and have unconsolidated BarTimeRanges
func scrub(ctx context.Context, client *ent.Client) error {
	for {
		// in the scrubber we only care about inactive intervals because active ones are handle in the grabber
		// we query for just one, so that we reduce the possibility of it being changed to active while we're scrubbing it
		stockInterval, err := client.Interval.Query().Where(
			interval.Active(false),
			interval.HasBarsWith(
				bartimerange.StatusNotIn(bartimerange.StatusConsolidated),
			),
		).First(ctx)
		if err != nil {
			if _, ok := err.(*ent.NotFoundError); ok {
				return nil
			}
			return errors.Wrap(err, "scrub")
		}

		if stockInterval == nil {
			return fmt.Errorf("expected to find interval but was nil")
		}

		// We want to get rid of duplication by consolidating BarTimeRange and BarGroups and making sure we
		// have unique BarRecords so we clean then consolidate
		logrus.Debugf("Bars.Scrubber: cleaning")
		err = clean(ctx, client, stockInterval)
		if err != nil {
			logrus.Errorf("Bars.Scrubber: clean error for interval(%v): %v", stockInterval.ID, err)
			continue
		}

		logrus.Debugf("Bars.Scrubber: consolidating")
		err = consolidate(ctx, client, stockInterval)
		if err != nil {
			logrus.Errorf("Bars.Scrubber: consolidate error for interval(%v): %v", stockInterval.ID, err)
			continue
		}
	}
}

// clean: cleans up "bad" or miss-matched records associated with a
// particular BarTimeRange specifically looking for:
// 1) BarTimeRange with miss-match of BarGroups
// 2) BarGroups with a miss-match of bar records
func clean(ctx context.Context, client *ent.Client, stockInterval *ent.Interval) error {
	// lets get all the barTimeRanges that haven't been cleaned yet
	barTimeRanges, err := stockInterval.QueryBars().
		Where(bartimerange.StatusIn(bartimerange.StatusCreated)).
		All(ctx)
	if err != nil {
		return errors.Wrap(err, "clean - querying bars")
	}

	for _, bar := range barTimeRanges {
		// anything else was created correctly but we're going to check that
		// the number match the expected
		if err := cleanBarTimeRange(ctx, client, bar); err != nil {
			return errors.Wrap(err, "clean")
		}
	}

	return nil
}

func cleanBarTimeRange(ctx context.Context, client *ent.Client, bar *ent.BarTimeRange) error {
	barGroups, err := bar.QueryGroups().All(ctx)
	if err != nil {
		return errors.Wrapf(err, "cleanBarTimeRange - querying BarGroup of BarTimeRange(%v)", bar.ID)
	}

	if len(barGroups) != bar.Count {
		deleteBarTimeRange(ctx, client, bar)
		return nil
	}

	for _, barGroup := range barGroups {
		records, err := barGroup.QueryRecords().All(ctx)
		if err != nil {
			return errors.Wrap(err, "cleanBarTimeRange - record query")
		}

		if len(records) != barGroup.Count {
			deleteBarTimeRange(ctx, client, bar)
			return nil
		}
	}

	_, err = bar.Update().SetStatus(bartimerange.StatusClean).Save(ctx)
	if err != nil {
		return errors.Wrapf(err, "cleanBarTimeRange - updating status(%v) - BarTimeRange(%v)", bartimerange.StatusClean, bar.ID)
	}
	return nil
}

// consolidate: given an interval it will consolidate all associated BarTimeRanges
func consolidate(ctx context.Context, client *ent.Client, stockInterval *ent.Interval) error {

	// lets get all the bartimeranges
	barTimeRanges, err := stockInterval.QueryBars().
		Where(bartimerange.StatusIn(bartimerange.StatusClean, bartimerange.StatusConsolidated)).
		All(ctx)
	if err != nil {
		return errors.Wrap(err, "consolidate - querying bars")
	}

	logrus.Debugf("consolidate - barTimeRanges %v", barTimeRanges)
	barToTimeIntervals := findUniqueIntervalToBarTimeRanges(barTimeRanges)

	for gi, bars := range barToTimeIntervals {
		if len(bars) == 1 {
			logrus.Debugf("consolidate - skipping consolidation %v-%v", gi, bars)
			// this case will cause useless copy so we
			// update status if needed then continue to the next
			if bars[0].Status == bartimerange.StatusConsolidated {
				continue
			}

			_, err = bars[0].Update().SetStatus(bartimerange.StatusConsolidated).Save(ctx)
			if err != nil {
				return errors.Wrapf(err, "consolidate")
			}
			continue
		}
		logrus.Debugf("consolidate - consolidating %v-%v", gi, bars)
		err := consolidateBarTimeRanges(ctx, client, stockInterval.Interval, gi, bars...)
		if err != nil {
			return errors.Wrapf(err, "consolidate")
		}
	}

	return nil
}

// consolidateBarTimeRanges take a all the overlapping bars and moves all BarGroups to the largest one (or older) and update it with the new TimeInterval.
// Once moved and updated, it will then consolidate the BarGroups
func consolidateBarTimeRanges(ctx context.Context, client *ent.Client, intervalSize interval.Interval, barInterval internal.TimeInterval, bars ...*ent.BarTimeRange) (err error) {
	//with the biggest range of time
	sort.Slice(bars, func(i, j int) bool {
		cmp := int(bars[i].End.Sub(bars[i].Start) - bars[j].End.Sub(bars[j].Start))
		if cmp == 0 {
			if bars[i].Start.Equal(bars[j].Start) {
				return bars[i].ID < bars[j].ID
			}

			return bars[i].Start.Before(bars[j].Start)
		}
		return cmp > 0
	})

	//we pick the first one to move everything to
	barIds := make([]int, 0, len(bars)-1)
	for _, bar := range bars[1:] {
		barIds = append(barIds, bar.ID)
	}

	// move BarGroup to th
	num, err := client.BarGroup.Update().
		Where(bargroup.HasTimeRangeWith(bartimerange.IDIn(barIds...))).
		SetTimeRangeID(bars[0].ID).
		Save(ctx)
	if err != nil {
		return errors.Wrap(err, "consolidateBarTimeRanges Update1")
	}

	// update the BarTimeRange id of the groups
	_, err = bars[0].Update().
		AddCount(num).
		SetStart(barInterval.Start).
		SetEnd(barInterval.End).
		Save(ctx)
	if err != nil {
		return errors.Wrap(err, "consolidateBarTimeRanges Update2")
	}

	groups, err := client.BarGroup.Query().
		Where(bargroup.HasTimeRangeWith(bartimerange.ID(bars[0].ID))).
		All(ctx)
	if err != nil {
		return errors.Wrap(err, "consolidateBarTimeRanges Query")
	}

	logrus.Debugf("consolidateBarTimeRanges findUniqueIntervalToBarGroups(%v) for bar(%v)", len(groups), bars[0].ID)
	groupIntervals := findUniqueIntervalToBarGroups(intervalSize, groups)

	for i, gs := range groupIntervals {
		if len(gs) == 1 {
			continue
		}
		logrus.Debugf("consolidateBarTimeRanges consolidating %v groups", len(gs))
		err = consolidateBarTimeRangesGroups(ctx, client, bars[0].ID, i, gs)
		if err != nil {
			return errors.Wrap(err, "consolidateBarTimeRanges")
		}
	}

	num, err = client.BarTimeRange.Delete().
		Where(bartimerange.IDIn(barIds...)).
		Exec(ctx)
	if err != nil {
		return errors.Wrap(err, "consolidateBarTimeRanges Delete")
	}

	if num != len(barIds) {
		return fmt.Errorf("consolidateBarTimeRanges deleted bars unexpected count %v expected %v", num, len(barIds))
	}

	return nil
}

func consolidateBarTimeRangesGroups(ctx context.Context, client *ent.Client, timeRangeID int, groupInterval internal.TimeInterval, groups []*ent.BarGroup) error {
	sort.Slice(groups, func(i, j int) bool {
		cmp := int(groups[i].Last.Sub(groups[i].First) - groups[j].Last.Sub(groups[j].First))
		if cmp < 0 {
			return false
		}
		if cmp > 0 {
			return true
		}

		if groups[i].First.Equal(groups[j].First) {
			return groups[i].ID < groups[j].ID
		}

		return groups[i].First.Before(groups[j].First)
	})

	groupIds := make([]int, 0, len(groups)-1)
	for _, group := range groups[1:] {
		groupIds = append(groupIds, group.ID)
	}

	//get all the BarRecords
	num, err := client.BarRecord.Update().
		Where(barrecord.HasGroupWith(bargroup.IDIn(groupIds...))).
		SetGroupID(groups[0].ID).
		Save(ctx)
	if err != nil {
		return errors.Wrap(err, "consolidateBarTimeRangesGroups Update1")
	}

	//now we do the move, move records to the one group
	_, err = groups[0].Update().
		SetFirst(groupInterval.Start).
		SetLast(groupInterval.End).
		AddCount(num).
		Save(ctx)
	if err != nil {
		return errors.Wrap(err, "consolidateBarTimeRangesGroups Update2")
	}
	records, err := client.BarRecord.Query().
		Where(barrecord.HasGroupWith(bargroup.ID(groups[0].ID))).
		All(ctx)
	if err != nil {
		return errors.Wrap(err, "consolidateBarTimeRangesGroups Query")
	}

	logrus.Debugf("consolidateBarTimeRangesGroups consolidating %v records", len(records))
	err = consolidateBarTimeRangesGroupsRecords(ctx, client, groups[0].ID, records...)
	if err != nil {
		return errors.Wrap(err, "consolidateBarTimeRangesGroups")
	}
	num, err = client.BarGroup.Delete().Where(
		bargroup.IDIn(groupIds...),
	).Exec(ctx)
	if err != nil {
		return errors.Wrap(err, "consolidateBarTimeRangesGroups Delete")
	}

	if num != len(groupIds) {
		return fmt.Errorf("consolidateBarTimeRangesGroups deleted groups unexpected count %v expected %v", num, len(groupIds))
	}

	return nil
}

func consolidateBarTimeRangesGroupsRecords(ctx context.Context, client *ent.Client, groupId int, records ...*ent.BarRecord) error {
	// we sort by time and id with there's a record with the same Timestamp we favor the older record
	sort.Slice(records, func(i, j int) bool {
		if records[i].Timestamp.Equal(records[j].Timestamp) {
			return records[i].ID < records[j].ID
		}
		return records[i].Timestamp.Before(records[j].Timestamp)
	})

	//then we just do a simple loop to find the dups
	delRecordIds := make([]int, 0)
	for i := 1; i < len(records); i++ {
		if records[i-1].Timestamp.Equal(records[i].Timestamp) {
			delRecordIds = append(delRecordIds, records[i].ID)
			slices.Delete(records, i, i+1)
			i--
		}
	}

	err := client.BarGroup.UpdateOneID(groupId).
		AddCount(-len(delRecordIds)).
		RemoveRecordIDs(delRecordIds...).Exec(ctx)
	if err != nil {
		return errors.Wrap(err, "consolidateBarTimeRangesGroupsRecords UpdateOneID")
	}

	num, err := client.BarRecord.Delete().Where(
		barrecord.IDIn(delRecordIds...),
	).Exec(ctx)
	if err != nil {
		return errors.Wrap(err, "consolidateBarTimeRangesGroupsRecords Delete")
	}

	if num != len(delRecordIds) {
		return fmt.Errorf("consolidateBarTimeRangesGroupsRecords deleted records unexpected count %v expected %v", num, len(delRecordIds))
	}
	return nil
}

// findUniqueIntervalToBarTimeRanges will determine all the intervals and return a map with intersecting BarTimeRanges
// that overlap in value of the map
func findUniqueIntervalToBarTimeRanges(ranges []*ent.BarTimeRange) map[internal.TimeInterval][]*ent.BarTimeRange {
	if len(ranges) == 0 {
		return map[internal.TimeInterval][]*ent.BarTimeRange{}
	}

	type item struct {
		I internal.TimeInterval
		B []*ent.BarTimeRange
	}

	m := make([]*item, len(ranges))
	for i, r := range ranges {
		m[i] = &item{
			I: internal.TimeInterval{r.Start, r.End},
			B: []*ent.BarTimeRange{ranges[i]},
		}
	}

	// sort by Start time.Time, and when Start's are equal sort by longest, when equal in length sort by id
	sort.Slice(m, func(i, j int) bool {
		if m[i].B[0].Start.Equal(m[j].B[0].Start) {
			if m[i].B[0].End.Equal(m[j].B[0].End) {
				return m[i].B[0].ID < m[j].B[0].ID
			}
			return m[i].B[0].End.After(m[j].B[0].End)
		}
		return m[i].B[0].Start.Before(m[j].B[0].Start)
	})

	result := make(map[internal.TimeInterval][]*ent.BarTimeRange)
	j := 0
	for i := range m {
		if i == j {
			if i == len(m)-1 {
				result[m[i].I] = m[i].B
			}
			continue
		}

		b1 := m[j]
		b2 := m[i]

		if internal.IntervalOverlap(b1.I, b2.I) {
			b1.I = internal.IntervalUnion(b1.I, b2.I)[0]

			b1.B = append(b1.B, b2.B...)
			b2.B = nil

			if i != len(m)-1 {
				continue
			}
		}

		//else jth item is one to add to the results, and we reduce the interval
		result[b1.I] = b1.B

		//if this is the last time we'll see i then add it as well
		if i == len(m)-1 && len(b2.B) > 0 {
			result[b2.I] = b2.B
		}

		//lastly assign j=i and go to the next i
		j = i
	}
	return result
}

// findUniqueIntervalToBarTimeRanges will determine all the intervals and return a map with intersecting BarTimeRanges
// that overlap in value of the map
func findUniqueIntervalToBarGroups(intervalSize interval.Interval, groups []*ent.BarGroup) map[internal.TimeInterval][]*ent.BarGroup {
	if len(groups) == 0 {
		return map[internal.TimeInterval][]*ent.BarGroup{}
	}

	type item struct {
		I internal.TimeInterval
		G []*ent.BarGroup
	}

	m := make([]*item, len(groups))
	for i, g := range groups {
		m[i] = &item{
			I: expandTimeInterval(internal.TimeInterval{g.First, g.Last}, intervalSize),
			G: []*ent.BarGroup{groups[i]},
		}
	}

	// sort by First time.Time, and when First's are equal sort by longest, when equal in length sort by id
	sort.Slice(m, func(i, j int) bool {
		if m[i].G[0].First.Equal(m[j].G[0].First) {
			if m[i].G[0].Last.Equal(m[j].G[0].Last) {
				return m[i].G[0].ID < m[j].G[0].ID
			}
			return m[i].G[0].Last.After(m[j].G[0].Last)
		}
		return m[i].G[0].First.Before(m[j].G[0].First)
	})

	result := make(map[internal.TimeInterval][]*ent.BarGroup)
	//now loop through and those next to one another should overlap if they need to be merged
	j := 0
	for i := range m {
		if i == j {
			if i == len(m)-1 {
				result[m[i].I] = m[i].G
			}
			continue
		}

		g1 := m[j]
		g2 := m[i]

		//if there's overlap we group it and expand the interval for the next run
		if internal.IntervalOverlap(g1.I, g2.I) {
			g1.I = internal.IntervalUnion(g1.I, g2.I)[0]

			g1.G = append(g1.G, g2.G...)
			g2.G = nil

			if i != len(m)-1 {
				continue
			}
		}

		//else jth item is one to add to the results, and we reduce the interval
		result[reduceTimeInterval(g1.I, intervalSize)] = g1.G

		//if this is the last time we'll see i then add it as well
		if i == len(m)-1 && len(g2.G) > 0 {
			result[reduceTimeInterval(g2.I, intervalSize)] = g2.G
		}

		//lastly assign j=i and go to the next i
		j = i
	}
	return result
}

func expandTimeInterval(g internal.TimeInterval, currentIntervalSize interval.Interval) internal.TimeInterval {
	switch currentIntervalSize {
	case interval.Interval1min:
		return internal.TimeInterval{g.Start, g.End.Add(time.Minute)}
	case interval.IntervalDaily:
		return internal.TimeInterval{g.Start, g.End.AddDate(0, 0, 1)}
	case interval.IntervalMonthly:
		panic("this needs to be validated it works")
		return internal.TimeInterval{g.Start, g.End.AddDate(0, 1, 0)}
	case interval.IntervalYearly:
		panic("this needs to be validated it works")
		return internal.TimeInterval{g.Start, g.End.AddDate(1, 0, 0)}
	}
	return g
}

func reduceTimeInterval(g internal.TimeInterval, currentIntervalSize interval.Interval) internal.TimeInterval {
	switch currentIntervalSize {
	case interval.Interval1min:
		return internal.TimeInterval{g.Start, g.End.Add(-time.Minute)}
	case interval.IntervalDaily:
		return internal.TimeInterval{g.Start, g.End.AddDate(0, 0, -1)}
	case interval.IntervalMonthly:
		panic("this needs to be validated it works")
		return internal.TimeInterval{g.Start, g.End.AddDate(0, -1, 0)}
	case interval.IntervalYearly:
		panic("this needs to be validated it works")
		return internal.TimeInterval{g.Start, g.End.AddDate(-1, 0, 0)}
	}
	return g
}

// deleteBarTimeRange: delete the BarTimeRange and all associated BarGroups
func deleteBarTimeRange(ctx context.Context, client *ent.Client, barTimeRange *ent.BarTimeRange) error {

	//we'll set the status to pending just in case we get stopped in the middle
	barTimeRangeUpdated, err := barTimeRange.Update().SetStatus(bartimerange.StatusPending).Save(ctx)
	if err != nil {
		return errors.Wrapf(err, "delete - update BarTimeRange(%v)", barTimeRange.ID)
	}

	//we find all the BarGroups, then BarRecords
	// then working backward delete from the DB
	groupsFromBar, err := barTimeRangeUpdated.QueryGroups().All(ctx)
	if err != nil {
		return errors.Wrapf(err, "delete - group query for BarTimeRange(%v)", barTimeRangeUpdated.ID)
	}

	groupsFromQuery, err := client.BarGroup.Query().Where(
		bargroup.HasTimeRangeWith(
			bartimerange.ID(barTimeRangeUpdated.ID),
		),
	).All(ctx)
	if err != nil {
		return errors.Wrapf(err, "delete - group verify query for BarTimeRange(%v)", barTimeRangeUpdated.ID)
	}

	if len(groupsFromBar) != len(groupsFromQuery) {
		return fmt.Errorf("consistency error len(groupsFromBar(%v))!=len(groupsFromQuery(%v))", len(groupsFromBar), len(groupsFromQuery))
	}

	for _, group := range groupsFromBar {
		records, err := group.QueryRecords().All(ctx)
		if err != nil {
			return errors.Wrapf(err, "delete - record query for BarGroup(%v)", group.ID)
		}

		num, err := client.BarRecord.Delete().Where(
			barrecord.HasGroupWith(
				bargroup.ID(group.ID),
			),
		).Exec(ctx)
		if err != nil {
			return errors.Wrapf(err, "delete - delete record for BarGroup(%v)", group.ID)
		}

		if len(records) != num {
			return fmt.Errorf("delete - records deleted expected %v found %v for BarGroup(%v)", len(records), num, group.ID)
		}
	}

	num, err := client.BarGroup.Delete().Where(
		bargroup.HasTimeRangeWith(
			bartimerange.ID(barTimeRangeUpdated.ID),
		),
	).Exec(ctx)
	if err != nil {
		return errors.Wrapf(err, "delete - delete BarGroup for BarTimeRange(%v)", barTimeRangeUpdated.ID)
	}
	if len(groupsFromBar) != num {
		return fmt.Errorf("delete - groups expected %v found %v for BarTimeRange(%v)", len(groupsFromBar), num, barTimeRangeUpdated.ID)
	}

	err = client.BarTimeRange.DeleteOne(barTimeRangeUpdated).Exec(ctx)
	if err != nil {
		return errors.Wrapf(err, "delete - BarTimeRange(%v)", barTimeRangeUpdated.ID)
	}

	return nil
}
