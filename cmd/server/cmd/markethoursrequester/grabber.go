package markethoursrequester

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/nathanhack/sibyl/agents"
	"github.com/nathanhack/sibyl/cmd/server/cmd/internal"
	"github.com/nathanhack/sibyl/ent"
	"github.com/nathanhack/sibyl/ent/markethours"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
)

func Grabber(ctx context.Context, client *ent.Client, agent agents.MarketHoursRequester, wg *sync.WaitGroup) {
	logrus.Infof("MarketHours.Grabber(%v): Running", agent.Name())
	startupTimer := time.NewTimer(5 * time.Second)
	wg.Add(1)
	defer func() {
		logrus.Infof("MarketHours.Grabber(%v): Stopped", agent.Name())
		wg.Done()
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case <-startupTimer.C:
		case <-time.After(24 * time.Hour):
		}

		startDate := time.Date(1900, 0, 0, 0, 0, 0, 0, time.Local)
		endDate := time.Now().AddDate(0, 0, 1)

		currentIntervals := []internal.TimeInterval{{Start: startDate, End: endDate}}

		marketInfo, err := client.MarketInfo.Query().Only(ctx)
		if err != nil {
			//if it was just not found then we continue on otherwise we bail
			if _, ok := err.(*ent.NotFoundError); !ok {
				logrus.Errorf("MarketHours.Grabber(%v): Query: %v", agent.Name(), err)
				continue
			}
		}

		if marketInfo != nil {
			currentIntervals = internal.IntervalDifferenceSlice(currentIntervals,
				internal.TimeInterval{Start: marketInfo.HoursStart, End: marketInfo.HoursEnd})
		}

		for _, i := range currentIntervals {
			creates, dates, err := agent.MarketHoursRequest(ctx, i.Start, i.End)
			if err != nil {
				logrus.Errorf("MarketHours.Grabber(%v): MarketHoursRequest: %v", agent.Name(), err)
				continue
			}

			err = processData(ctx, client, creates, dates, i, marketInfo)
			if err != nil {
				logrus.Errorf("MarketHours.Grabber(%v): processData: %v", agent.Name(), err)
			}
		}
		logrus.Infof("MarketHours.Grabber(%v) completed", agent.Name())
	}
}

func processData(ctx context.Context, client *ent.Client, creates []*ent.MarketHoursCreate, dates []time.Time, tInterval internal.TimeInterval, marketInfo *ent.MarketInfo) error {
	if len(creates) == 0 {
		return nil
	}
	var err error
	hours := []*ent.MarketHours{}

	if marketInfo != nil {
		hours, err = marketInfo.QueryHours().Where(
			markethours.DateGTE(tInterval.Start.AddDate(0, 0, -5)),
			markethours.DateLTE(tInterval.End),
		).Order(ent.Asc()).All(ctx)
		if err != nil {
			return err
		}
	}

	//remove any that we already have
	newHours := make([]*ent.MarketHoursCreate, 0)
	for index, d := range dates {
		//we make sure it's not one already there
		i := slices.IndexFunc(hours, func(mh *ent.MarketHours) bool {
			return mh.Date.Equal(d)
		})
		if i >= 0 {
			continue
		}

		newHours = append(newHours, creates[index])
	}

	// if there's none to add
	if len(newHours) == 0 {
		return nil
	}

	// union the intervals
	u := []internal.TimeInterval{tInterval}
	if marketInfo != nil {
		u = internal.IntervalUnion(
			internal.TimeInterval{Start: marketInfo.HoursStart, End: marketInfo.HoursEnd},
			tInterval,
		)
	}

	// it should only be 1 if not we got a problem
	if len(u) != 1 {
		return fmt.Errorf("expected 1 interval found: %v", u)
	}

	// add the hours
	hs, err := client.MarketHours.CreateBulk(newHours...).Save(ctx)
	if err != nil {
		return err
	}

	if marketInfo == nil {
		return client.MarketInfo.Create().
			SetHoursStart(u[0].Start).
			SetHoursEnd(u[0].End).
			AddHours(hs...).
			Exec(ctx)
	}

	// lastly update the info with the new interval
	return marketInfo.Update().
		SetHoursStart(u[0].Start).
		SetHoursEnd(u[0].End).
		AddHours(hs...).
		Exec(ctx)
}
