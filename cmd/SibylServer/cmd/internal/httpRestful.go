package internal

import (
	"context"
	"fmt"
	"github.com/nathanhack/sibyl/core/database"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type HttpRestServer struct {
	killCtx       context.Context
	kill          context.CancelFunc
	doneCtx       context.Context
	done          context.CancelFunc
	address       string
	running       bool
	serverDied    context.CancelFunc
	serverContext *ServerContext
}

func NewHttpRestServer(db *database.SibylDatabase, stockCache *StockCache, serverAddress string, stockValidator *StockValidator, serverDied context.CancelFunc) *HttpRestServer {
	killCtx, kill := context.WithCancel(context.Background())
	doneCtx, done := context.WithCancel(context.Background())
	sc := ServerContext{
		Ctx:            killCtx,
		db:             db,
		stockValidator: stockValidator,
		stockCache:     stockCache,
	}

	return &HttpRestServer{
		killCtx:       killCtx,
		kill:          kill,
		doneCtx:       doneCtx,
		done:          done,
		address:       serverAddress,
		serverDied:    serverDied,
		serverContext: &sc,
	}
}

func (hrs *HttpRestServer) Run() error {
	if hrs.running {
		return fmt.Errorf("HttpRestServer is already running")
	}

	var err error
	var httpServer *http.Server
	httpServer, err = makeServer(hrs.serverContext, hrs.address)
	if err != nil {
		return fmt.Errorf("HttpRestServer failed to make the server: %v", err)
	}

	serverStoppedCtx, serverStopped := context.WithCancel(context.Background())
	go func(server *http.Server, serverStopped context.CancelFunc) {
		if err := server.ListenAndServe(); err != http.ErrServerClosed && err != nil {
			logrus.Errorf("server had an error: %v", err)
		}
		serverStopped()
	}(httpServer, serverStopped)

	//we sleep to give the go routine a chance to start up
	time.Sleep(100 * time.Millisecond)
	//let's make sure the server actually started
	select {
	case <-serverStoppedCtx.Done():
		return fmt.Errorf("HttpRestServer failed to start the server")
	default:
	}

	//well we made it this far time to start the main loop
	go func(hrs *HttpRestServer, serverStoppedCtx context.Context) {
	mainLoop:
		for {
			select {
			case <-serverStoppedCtx.Done():
				hrs.serverDied()
				break mainLoop
			case <-hrs.killCtx.Done():
				//since the longest task takes less than 1 min that is the timer
				ctx, _ := context.WithTimeout(context.Background(), 1*time.Minute)
				httpServer.Shutdown(ctx)
				break mainLoop
			}
		}
		hrs.done()
	}(hrs, serverStoppedCtx)
	return nil
}

func (hrs *HttpRestServer) Stop(waitUpTo time.Duration) {
	//next stop the quoter
	logrus.Infof("Waiting for HttpRestServer to finish")
	hrs.kill()
	select {
	case <-hrs.doneCtx.Done():
		logrus.Infof("HttpRestServer finished")
	case <-time.After(waitUpTo):
		logrus.Errorf("HttpRestServer failed to gracefully finish")
	}
}
