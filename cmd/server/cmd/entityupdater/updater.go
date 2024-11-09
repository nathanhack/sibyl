package entityupdater

import (
	"context"
	"sync"
	"time"

	"github.com/nathanhack/sibyl/agents"
	"github.com/nathanhack/sibyl/cmd/server/cmd/internal"
	"github.com/nathanhack/sibyl/ent"
	"github.com/nathanhack/sibyl/ent/entity"
	"github.com/sirupsen/logrus"
)

func Updater(ctx context.Context, client *ent.Client, agent agents.EntitySourcing, wg *sync.WaitGroup) {
	logrus.Infof("Entity.Updater(%v): Running", agent.Name())
	startupTimer := time.NewTimer(5 * time.Second)
	wg.Add(1)
	defer func() {
		logrus.Infof("Entity.Updater(%v): Stopped", agent.Name())
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
			logrus.Errorf("Entity.Updater(%v): failed to get stocks: %v", agent.Name(), err)
			continue
		}

		updates, err := agent.EntityCreateUpdate(ctx, stocks...)
		if err != nil {
			logrus.Errorf("Entity.Updater(%v): %v", agent.Name(), err)
		}
		for _, update := range updates {
			select {
			case <-ctx.Done():
				return
			default:
				_, err = update.Save(ctx)
				if err != nil {
					logrus.Errorf("Entity.Updater(%v): failed to get stocks: %v", agent.Name(), err)
				}
			}
		}
	}
}
