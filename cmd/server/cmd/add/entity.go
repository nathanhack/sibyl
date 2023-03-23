package add

import (
	"context"
	"strings"

	"github.com/nathanhack/sibyl/agents"
	"github.com/nathanhack/sibyl/ent"
	"github.com/nathanhack/sibyl/ent/entity"
	"github.com/sirupsen/logrus"
)

var AddEntity = make(chan string, 10000)
var filtered = make(chan string, 1000)

func Entity(ctx context.Context, entClient *ent.Client, agent agents.EntitySearcher) {
	logrus.Info("AddEntity Running")

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case x := <-AddEntity:
				x = strings.ToUpper(x)

				stock, _ := entClient.Entity.Query().Where(
					entity.Ticker(x),
				).Only(ctx)

				if stock != nil {
					// then we already have it so don't try and get it
					logrus.Infof("AddEntity: duplicate request received to add %v", x)
					continue
				}
				logrus.Infof("AddEntity: received request to add %v", x)
				filtered <- x
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			logrus.Info("AddEntity Stopping")
			return
		case x := <-filtered:
			entity, err := agent.Entity(ctx, x)
			if err != nil {
				logrus.Errorf("AddEntity(%v): %v", x, err)
				continue
			}
			_, err = entity.Save(ctx)
			if err != nil {
				logrus.Errorf("AddEntity(%v): %v", x, err)
				continue
			}
			logrus.Infof("AddEntity: %v added", x)
		}
	}
}
