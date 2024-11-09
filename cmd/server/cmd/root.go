package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"time"

	entsql "entgo.io/ent/dialect/sql"

	"github.com/go-faster/jx"
	"github.com/nathanhack/sibyl/agents"
	"github.com/nathanhack/sibyl/agents/alpaca"
	"github.com/nathanhack/sibyl/agents/polygonio"
	"github.com/nathanhack/sibyl/cmd/server/cmd/add"
	"github.com/nathanhack/sibyl/cmd/server/cmd/barrequester"
	"github.com/nathanhack/sibyl/cmd/server/cmd/config"
	"github.com/nathanhack/sibyl/cmd/server/cmd/dividendrequester"
	"github.com/nathanhack/sibyl/cmd/server/cmd/entityupdater"
	"github.com/nathanhack/sibyl/cmd/server/cmd/markethoursrequester"
	"github.com/nathanhack/sibyl/cmd/server/cmd/splitrequester"
	"github.com/nathanhack/sibyl/ent"
	"github.com/nathanhack/sibyl/ent/migrate"
	"github.com/nathanhack/sibyl/ent/ogent"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/spf13/cobra"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
)

const (
	defaultConfigPathname = ".sibyl.json"
)

var cfgFilepath string

type handler struct {
	*ogent.OgentHandler
	db    *sql.DB
	agent agents.EntitySearcher
}

type searchCacheResult struct {
	result    *ogent.SearchTickerOK
	timestamp time.Time
}

var searchCache = map[string]searchCacheResult{}
var searchCacheMux = sync.RWMutex{}

func (h handler) SearchTicker(ctx context.Context, params ogent.SearchTickerParams) (*ogent.SearchTickerOK, error) {
	searchCacheMux.RLock()
	x, has := searchCache[params.Ticker]
	searchCacheMux.RUnlock()
	if !has || time.Since(x.timestamp) > 24*time.Hour {
		go func() {
			searchCacheMux.Lock()
			searchCache[params.Ticker] = searchCacheResult{timestamp: time.Now()}
			searchCacheMux.Unlock()

			results, err := h.agent.EntitySearch(context.Background(), strings.ToUpper(params.Ticker), 15)

			tmp := searchCacheResult{
				result:    &ogent.SearchTickerOK{},
				timestamp: time.Now(),
			}

			if err != nil {
				var e jx.Encoder
				e.Str(fmt.Sprint(err))

				tmp.result.Errors = e.Bytes()
			}

			tmp.result.Results = make([]ogent.SearchTickerOKResultsItem, len(results))

			for i := range tmp.result.Results {
				tmp.result.Results[i] = ogent.SearchTickerOKResultsItem{
					Ticker: results[i].Ticker,
					Name:   results[i].Name,
				}
			}

			logrus.Debugf("Results(%v) for /rest/search/%v", len(tmp.result.Results), params.Ticker)
			searchCacheMux.Lock()
			searchCache[params.Ticker] = tmp
			searchCacheMux.Unlock()
		}()

		var e jx.Encoder
		e.Str("processing")
		return &ogent.SearchTickerOK{
			Status: e.Bytes(),
		}, nil
	}

	if x.result == nil {
		var e jx.Encoder
		e.Str("processing")
		return &ogent.SearchTickerOK{
			Status: e.Bytes(),
		}, nil
	}

	return &ogent.SearchTickerOK{
		Results: x.result.Results,
		Errors:  x.result.Errors,
	}, nil
}

var tickerAdd = make(chan string, 1000)

func (h handler) AddTicker(ctx context.Context, params ogent.AddTickerParams) error {
	add.AddEntity <- params.Ticker
	return nil
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "server",
	Short: "The Sibyl suite backend tool for stock data acquisition, Display Only views and manual trade execution",
	Long:  `server is the Sibyl suite's backend server pursuant to being Display Only is used to acquire data from a discount broker, display data and perform trades.`,
	Run: func(cmd *cobra.Command, args []string) {
		// ctx, cancel := context.WithCancel(context.Background())

		setupLogger()

		//here we start up a few goroutines
		// 1) runs the server to handle api calls
		// 2) takes care of validating stocks
		// 3) takes care of running history queries on the stocks in the database
		// 4) takes care of running intraday queries on the stocks in the database
		// 5) takes care of running stableQuotes queries on the options in the database
		// 6) takes care of running stableQuotes queries on the stocks in the database
		// 7) takes care of running quotes queries on the options in the database
		// 8) takes care of running quotes queries on the stocks in the database
		// 9) takes care of keeping a cache of stocks and options to reduce database latency issues
		// if starting any of these fails we kill the program

		//first we connect with the database and die if it doesn't work
		logrus.Info("Establishing database context")
		drv, err := entsql.Open(config.State.Database.Dialect, config.State.Database.DSN)
		if err != nil {
			logrus.Fatalf("Could not establish connection to database: %v", err)
		}
		drv.DB().SetConnMaxLifetime(10 * time.Minute)
		logrus.Info("Database driver")
		client := ent.NewClient(ent.Driver(drv))

		ctx, cancel := context.WithCancel(context.Background())

		// On a fresh install the database will not exist so we create it if it's not there
		_, err = drv.DB().ExecContext(ctx, "CREATE DATABASE IF NOT EXISTS `sibyl`;")
		if err != nil {
			logrus.Fatalf("Create Database Error: %v", err)
		}

		// Run the migrations.
		err = client.Schema.Create(context.Background(),
			migrate.WithDropColumn(true),
			migrate.WithDropIndex(true),
		)
		if err != nil {
			logrus.Fatalf("Migration had an error: %v", err)
		}

		polyAgent, err := polygonio.New(ctx, client,
			config.State.Agents.PolygonIO.ApiKey,
			polygonio.Plan(config.State.Agents.PolygonIO.Plan))
		if err != nil {
			logrus.Fatalf("Unable to create Polygon.IO Agent error: %v", err)
		}

		var alpacaAgent *alpaca.Alpaca
		if strings.ToLower(config.State.Agents.Alpaca.Endpoint) == "live" {
			alpacaAgent, err = alpaca.New(ctx, client,
				config.State.Agents.Alpaca.Live.ApiKeyID,
				config.State.Agents.Alpaca.Live.SecretKey,
				config.State.Agents.Alpaca.Live.Url,
				alpaca.Plan(config.State.Agents.Alpaca.Plan))
		} else {
			alpacaAgent, err = alpaca.New(ctx, client,
				config.State.Agents.Alpaca.Paper.ApiKeyID,
				config.State.Agents.Alpaca.Paper.SecretKey,
				config.State.Agents.Alpaca.Paper.Url,
				alpaca.Plan(config.State.Agents.Alpaca.Plan))
		}
		if err != nil {
			logrus.Fatalf("Unable to create Alpaca Agent error: %v", err)
		}

		go add.Entity(ctx, client, alpacaAgent)
		wg := sync.WaitGroup{}
		go entityupdater.Updater(ctx, client, alpacaAgent, &wg)
		go entityupdater.Updater(ctx, client, polyAgent, &wg)
		go barrequester.Grabber(ctx, client, polyAgent, &wg)
		go barrequester.Grabber(ctx, client, alpacaAgent, &wg)
		go barrequester.Scrubber(ctx, client, &wg)
		go dividendrequester.Grabber(ctx, client, alpacaAgent, &wg)
		go splitrequester.Grabber(ctx, client, alpacaAgent, &wg)
		go markethoursrequester.Grabber(ctx, client, alpacaAgent, &wg)

		// Start listening.
		h := handler{
			OgentHandler: ogent.NewOgentHandler(client),
			db:           drv.DB(),
			agent:        polyAgent,
		}
		// handler, err := ogent.NewServer(ogent.NewOgentHandler(client))
		handler, err := ogent.NewServer(h)
		if err != nil {
			logrus.Fatal(err)
		}

		address, err := url.Parse(config.State.Address)
		logrus.Infof("Server on: %v", address)
		server := &http.Server{
			Addr:    address.Host,
			Handler: handler,
		}
		serverDiedCtx, serverDied := context.WithCancel(context.Background())
		go func() {
			err := server.ListenAndServe()
			if err != nil && err != http.ErrServerClosed {
				logrus.Fatal(err)
				serverDied()
			}
		}()

		signalInterruptChan := make(chan os.Signal, 1)
		signal.Notify(signalInterruptChan, os.Interrupt)
		signalKillChan := make(chan os.Signal, 1)
		signal.Notify(signalKillChan, os.Kill)

	mainLoop:
		for {
			select {
			case <-signalKillChan:
				logrus.Infof("Received a Kill signal stopping server.")
				break mainLoop
			case <-signalInterruptChan:
				logrus.Infof("Received a Interrupt signal stopping server.")
				break mainLoop
			case <-serverDiedCtx.Done():
				logrus.Errorf("Server Died unexpectedly.")
				break mainLoop
			}
		}

		shutdownCtx, _ := context.WithDeadline(ctx, time.Now().Add(15*time.Second))
		logrus.Info("Shutting down REST Server")
		err = server.Shutdown(shutdownCtx)
		if err != nil {
			logrus.Error(err)
		}
		cancel()
		logrus.Info("Shutting down BarGrabber")
		wg.Wait()
		err = client.Close()
		if err != nil {
			logrus.Error(err)
		}


	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFilepath, "config", "", "config file (default is $CWD/.sibyl.json)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetDefault("database.dialect", "sqlite3")
	viper.SetDefault("database.dsn", "./sqlite.db?_fk=1")
	viper.SetDefault("address", "http://localhost:9090")
	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.directory", "./")

	if cfgFilepath != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFilepath)
	} else {
		// Search config in home directory with name ".sibyl.json" (without extension).
		viper.AddConfigPath("./")
		viper.SetConfigFile(defaultConfigPathname)
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		logrus.Info("Using config file:", viper.ConfigFileUsed())

	} else {
		if err := viper.WriteConfigAs(defaultConfigPathname); err != nil {
			logrus.Error(err)
			os.Exit(-1)
		}
	}
	if err := viper.WriteConfig(); err != nil {
		logrus.Error(err)
		os.Exit(-1)
	}

	if err := viper.Unmarshal(&config.State); err != nil {
		logrus.Error(err)
		os.Exit(-1)
	}
}

func setupLogger() {
	responseErrorLog := filepath.Join(config.State.Logging.Directory + "server.log")
	logrus.Infof("Logging to %v", responseErrorLog)
	if _, err := os.Stat(config.State.Logging.Directory); os.IsNotExist(err) {
		if err := os.MkdirAll(config.State.Logging.Directory, 0775); err != nil {
			fmt.Printf("problems with creating the logging directory(%v): %v\n", config.State.Logging.Directory, err)
			return
		}
	}

	writer, _ := rotatelogs.New(
		responseErrorLog+".%Y%m%d%H%M",
		rotatelogs.WithLinkName(responseErrorLog),
		rotatelogs.WithMaxAge(time.Hour*24),         // one day logs
		rotatelogs.WithRotationTime(time.Hour*24*7), // for seven days
	)

	hook := lfshook.NewHook(
		lfshook.WriterMap{
			logrus.InfoLevel:  writer,
			logrus.DebugLevel: writer,
			logrus.ErrorLevel: writer,
			logrus.WarnLevel:  writer,
			logrus.FatalLevel: writer,
		},
		&logrus.JSONFormatter{},
	)

	logrus.AddHook(hook)

	//add line numbers
	//logrus.SetReportCaller(true)
	//format the line output
	//logrus.SetFormatter(&zt_formatter.ZtFormatter{
	//	//CallerPrettyfier: func(f *runtime.Frame) (string, string) { return strconv.Itoa(f.Line), f.File },
	//	CallerPrettyfier: func(f *runtime.Frame) (string, string) { return "", "" },
	//})

	//timestamp format
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	logrus.SetFormatter(customFormatter)
	customFormatter.FullTimestamp = true

	if strings.ToLower(config.State.Logging.Level) == "debug" {
		logrus.Info("logging level set to: DEBUG")
		logrus.SetLevel(logrus.DebugLevel)
	}
}
