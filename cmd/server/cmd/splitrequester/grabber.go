package splitrequester

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/nathanhack/sibyl/agents"
	"github.com/nathanhack/sibyl/cmd/server/cmd/internal"
	"github.com/nathanhack/sibyl/ent"
	"github.com/nathanhack/sibyl/ent/entity"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
)

func Grabber(ctx context.Context, client *ent.Client, agent agents.SplitRequester, wg *sync.WaitGroup) {
	logrus.Infof("Split.Grabber(%v): Running", agent.Name())
	startupTimer := time.NewTimer(5 * time.Second)
	wg.Add(1)
	defer func() {
		logrus.Infof("Split.Grabber(%v): Stopped", agent.Name())
		wg.Done()
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case <-startupTimer.C:
		case <-internal.AtMidnight():
		}

		//we first get all the stocks
		stocks, err := client.Entity.Query().Where(
			entity.Active(true),
		).All(ctx)

		if err != nil {
			logrus.Errorf("Split.Grabber(%v): failed to get stocks: %v", agent.Name(), err)
			continue
		}

		for _, stock := range stocks {
			select {
			case <-ctx.Done():
				return
			default:
			}

			// we take the executionDates and use that to determine if we already have the split recorded
			splits, err := stock.QuerySplits().All(ctx)
			if err != nil {
				logrus.Errorf("Split.Grabber(%v): QueryDividends: %v", agent.Name(), err)
				continue
			}

			sort.Slice(splits, func(i, j int) bool {
				return splits[i].ExecutionDate.Before(splits[j].ExecutionDate)
			})

			startDate := stock.ListDate
			if len(splits) > 0 {
				startDate = splits[len(splits)-1].ExecutionDate
			}

			splitCreates, executionDates, err := agent.SplitRequest(ctx, stock.Ticker, startDate, time.Now())
			if err != nil {
				logrus.Errorf("Split.Grabber(%v): DividendRequest: %v", agent.Name(), err)
				continue
			}

			err = processSplitData(ctx, stock.ID, splitCreates, executionDates, splits)
			if err != nil {
				logrus.Errorf("Split.Grabber(%v): processDividendData: %v", agent.Name(), err)
			}
		}
	}
}

func processSplitData(ctx context.Context, stockId int, splitCreates []*ent.SplitCreate, executionDates []time.Time, splits []*ent.Split) error {
	for index := range executionDates {
		i := slices.IndexFunc(splits, func(d *ent.Split) bool {
			return d.ExecutionDate.Equal(executionDates[index])
		})
		if i < 0 {
			_, err := splitCreates[index].SetStockID(stockId).Save(ctx)
			if err != nil {
				return err
			}
		}

	}
	return nil
}
