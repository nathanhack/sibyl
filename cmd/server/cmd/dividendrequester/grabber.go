package dividendrequester

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

func Grabber(ctx context.Context, client *ent.Client, agent agents.DividendRequester, wg *sync.WaitGroup) {
	logrus.Infof("Dividend.Grabber(%v): Running", agent.Name())
	startupTimer := time.NewTimer(5 * time.Second)
	wg.Add(1)
	defer func() {
		logrus.Infof("Dividend.Grabber(%v): Stopped", agent.Name())
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
			logrus.Errorf("Dividend.Grabber(%v): failed to get stocks: %v", agent.Name(), err)
			continue
		}

		for _, stock := range stocks {
			select {
			case <-ctx.Done():
				return
			default:
			}

			// we take the payDates and use that to determine if we already have the dividend recorded
			dividends, err := stock.QueryDividends().All(ctx)
			if err != nil {
				logrus.Errorf("Dividend.Grabber(%v): QueryDividends: %v", agent.Name(), err)
				continue
			}

			sort.Slice(dividends, func(i, j int) bool {
				return dividends[i].PayDate.Before(dividends[j].PayDate)
			})

			startDate := stock.ListDate
			if len(dividends) > 0 {
				startDate = dividends[len(dividends)-1].PayDate
			}

			dividendCreates, payDates, err := agent.DividendRequest(ctx, stock.Ticker, startDate, time.Now())
			if err != nil {
				logrus.Errorf("Dividend.Grabber(%v): DividendRequest: %v", agent.Name(), err)
				continue
			}

			err = processDividendData(ctx, stock.ID, dividendCreates, payDates, dividends)
			if err != nil {
				logrus.Errorf("Dividend.Grabber(%v): processDividendData: %v", agent.Name(), err)
			}
		}
	}
}

func processDividendData(ctx context.Context, stockId int, dividendCreates []*ent.DividendCreate, payDates []time.Time, dividends []*ent.Dividend) error {
	for index := range payDates {
		i := slices.IndexFunc(dividends, func(d *ent.Dividend) bool {
			return d.PayDate.Equal(payDates[index])
		})
		if i < 0 {
			_, err := dividendCreates[index].AddStockIDs(stockId).Save(ctx)
			if err != nil {
				return err
			}
		}

	}
	return nil
}
