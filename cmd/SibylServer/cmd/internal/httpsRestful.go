package internal

import (
	"context"
	"fmt"
	"github.com/nathanhack/sibyl/core/database"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type HttpsRestServer struct {
	killCtx            context.Context
	kill               context.CancelFunc
	doneCtx            context.Context
	done               context.CancelFunc
	address            string
	running            bool
	serverDied         context.CancelFunc
	serverContext      *ServerContext
	publicCertPathname string
	privateKeyPathname string
}

func NewHttpsRestServer(db *database.SibylDatabase, serverAddress, publicCertPathname, privateKeyPathname string, stockValidator *StockValidator, serverDied context.CancelFunc) *HttpsRestServer {
	killCtx, kill := context.WithCancel(context.Background())
	doneCtx, done := context.WithCancel(context.Background())
	sc := ServerContext{
		Ctx:            killCtx,
		db:             db,
		stockValidator: stockValidator,
	}

	return &HttpsRestServer{
		killCtx:            killCtx,
		kill:               kill,
		doneCtx:            doneCtx,
		done:               done,
		address:            serverAddress,
		serverDied:         serverDied,
		serverContext:      &sc,
		publicCertPathname: publicCertPathname,
		privateKeyPathname: privateKeyPathname,
	}
}

func (hrs *HttpsRestServer) Run() error {
	if hrs.running {
		return fmt.Errorf("HttpsRestServer is already running")
	}

	var err error
	var httpServer *http.Server
	httpServer, err = makeServer(hrs.serverContext, hrs.address)
	if err != nil {
		return fmt.Errorf("HttpsRestServer failed to make the server: %v", err)
	}

	serverStoppedCtx, serverStopped := context.WithCancel(context.Background())
	go func(hrs *HttpsRestServer, server *http.Server, serverStopped context.CancelFunc) {
		if err := server.ListenAndServeTLS(hrs.publicCertPathname, hrs.privateKeyPathname); err != http.ErrServerClosed && err != nil {
			logrus.Errorf("server had an error: %v", err)
		}
		serverStopped()
	}(hrs, httpServer, serverStopped)

	//we sleep to give the go routine a chance to start up
	time.Sleep(100 * time.Millisecond)
	//let's make sure the server actually started
	select {
	case <-serverStoppedCtx.Done():
		return fmt.Errorf("HttpsRestServer: failed to start the server")
	default:
	}

	//well we made it this far time to start the main loop
	go func(hrs *HttpsRestServer, serverStoppedCtx context.Context) {
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

func (hrs *HttpsRestServer) Stop(waitUpTo time.Duration) {
	//next stop the quoter
	logrus.Infof("Waiting for HttpsRestServer to finish")
	hrs.kill()
	select {
	case <-hrs.doneCtx.Done():
		logrus.Infof("HttpsRestServer finished")
	case <-time.After(waitUpTo):
		logrus.Errorf("HttpsRestServer failed to gracefully finish")
	}
}
