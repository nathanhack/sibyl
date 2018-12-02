package internal

import (
	"context"
	"fmt"
	"github.com/nathanhack/sibyl/core"
	"github.com/nathanhack/sibyl/core/database"
	"github.com/sirupsen/logrus"
	"time"
)

type HistoryGrabber struct {
	killCtx     context.Context
	kill        context.CancelFunc
	doneCtx     context.Context
	done        context.CancelFunc
	db          *database.SibylDatabase
	symbolCache *SymbolCache
	running     bool
}

func NewHistoryGrabber(db *database.SibylDatabase, symbolCache *SymbolCache) *HistoryGrabber {
	killCtx, kill := context.WithCancel(context.Background())
	doneCtx, done := context.WithCancel(context.Background())
	return &HistoryGrabber{
		killCtx:     killCtx,
		kill:        kill,
		doneCtx:     doneCtx,
		done:        done,
		db:          db,
		symbolCache: symbolCache,
	}
}

func (hg *HistoryGrabber) Run() error {
	if hg.running {
		return fmt.Errorf("HistoryGrabber is already running")
	}

	hg.running = true
	go func(hg *HistoryGrabber) {
		durationToWait := 23 * time.Second
		runHistoryGrabber := make(chan bool, 100)
	mainLoop:
		for {
			select {
			case <-hg.killCtx.Done():
				break mainLoop
			case <-hg.symbolCache.HistorySymbolsChanged:
				runHistoryGrabber <- true
			case <-time.After(durationToWait):
				runHistoryGrabber <- true
			case <-runHistoryGrabber:
				//first we drain the chan
			drainLoop:
				for {
					select {
					case <-runHistoryGrabber:
						continue
					default:
						break drainLoop
					}
				}

				currentTime := time.Now()
				agent, err := hg.db.GetAgent(hg.killCtx)
				if err != nil {
					logrus.Errorf("HistoryGrabber: could not retrieve agent: %v", err)
					continue
				}

				//schedule the next update to be tomorrow morning at 6am
				year, month, day := currentTime.Date()
				durationToWait = time.Date(year, month, day+1, 6, 0, 0, 0, time.Local).Sub(currentTime)

				//we want all the stock history up to yesterdayOrLastWeekday(or the last weekday)
				zeroDate := core.NewDateTypeFromUnix(0)
				today := core.NewDateTypeFromTime(currentTime)
				yesterdayOrLastWeekday := today.AddDate(0, 0, -1)
				for yesterdayOrLastWeekday.IsWeekDay() {
					//since we need at least two days for history quote
					yesterdayOrLastWeekday = yesterdayOrLastWeekday.AddDate(0, 0, -1)
				}
				stockAndDays := make(map[core.StockSymbolType]core.DateType)

				hg.symbolCache.HistorySymbolsMu.RLock()
				historySymbolsToDownload := hg.symbolCache.HistorySymbols
				hg.symbolCache.HistorySymbolsMu.RUnlock()

				for stock := range historySymbolsToDownload {
					//with history we want 10 years
					//TODO make the history time a configuration
					lastDate, err := hg.db.LastHistoryDate(hg.killCtx, stock)
					if err != nil ||
						lastDate == zeroDate {
						if err != nil {
							logrus.Infof("HistoryGrabber: looks like %v or was unable to located it in the database, default history request for 20 years will be used error: %v", stock, err)
						} else {
							logrus.Infof("HistoryGrabber: looks like %v or was unable to located it in the database, default history request for 20 years will be used", stock)
						}
						stockAndDays[stock] = yesterdayOrLastWeekday.AddDate(-20, 0, 0)
					} else {
						if lastDate.Before(yesterdayOrLastWeekday) {
							stockAndDays[stock] = lastDate.AddDate(0, 0, -2)
						}
					}
				}

				for stock, startTime := range stockAndDays {
					select {
					case <-hg.killCtx.Done():
						logrus.Errorf("HistoryGrabber: context canceled")
						break mainLoop
					default:
					}

					if startTime.Before(yesterdayOrLastWeekday) {
						logrus.Infof("HistoryGrabber: getting history info on stock %v", stock)
						historyRecords, err := agent.GetHistory(hg.killCtx, stock, core.DailyTicks, startTime, today)
						if err != nil {
							logrus.Errorf("HistoryGrabber: had a problem getting History Records for %v: %v", stock, err)
							//well we failed to get the data give it another chance soon
							runHistoryGrabber <- true
						}

						if err = hg.db.LoadHistoryRecords(hg.killCtx, historyRecords); err != nil {
							logrus.Errorf("HistoryGrabber: had a problem saving History Records for %v: %v", stock, err)
							//well we failed to load the data give it another chance soon
							runHistoryGrabber <- true
						}
					}
				}
				logrus.Infof("HistoryGrabber: finished a round in %v, next round in %v", time.Since(currentTime), durationToWait)
			}
		}
		hg.done() //signal this is finished
	}(hg)
	return nil
}

func (hg *HistoryGrabber) Stop(waitUpTo time.Duration) {
	//next stop the quoter
	logrus.Infof("Waiting for HistoryGrabber to finish")
	startTime := time.Now()
	hg.kill()
	select {
	case <-hg.doneCtx.Done():
		logrus.Infof("HistoryGrabber finished in %v", time.Since(startTime))
	case <-time.After(waitUpTo):
		logrus.Errorf("HistoryGrabber failed to gracefully finish in %v", time.Since(startTime))
	}
}
