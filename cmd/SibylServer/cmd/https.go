// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"context"
	"github.com/nathanhack/sibyl/cmd/SibylServer/cmd/internal"
	"github.com/nathanhack/sibyl/core/database"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"time"
)

// httpsCmd represents the https command
var httpsCmd = &cobra.Command{
	Use:   "https",
	Short: "SibylServer will run with a REST interface via HTTPS",
	Long:  `SibylServer will start up a HTTPS RESTFUL interface for interacting with users`,
	Run: func(cmd *cobra.Command, args []string) {
		//here we start up a few goroutines
		// 1) runs the server to handle api calls
		// 2) takes care of validating stocks
		// 3) takes care of running history queries on the stocks in the database
		// 4) takes care of running intraday queries on the stocks in the database
		// 5) takes care of running stableQuotes queries on the stocks in the database
		// 6) takes care of running quotes queries on the stocks in the database
		// 7) takes care of keeping a cache of stocks and options to reduce database latency issues
		// if starting any of these fails we kill the program

		//first we connect with the database and die if it doesn't work
		db, err := database.ConnectAndEnsureSibylDatabase(context.Background(), cmd.Flag("database").Value.String())
		if err != nil {
			logrus.Errorf("Could not establish connection to database: %v", err)
			os.Exit(-1)
		}

		//
		stockCache := internal.NewStockCache(db)
		if err := stockCache.Run(); err != nil {
			logrus.Errorf("Starting StockCache had an issue: %v", err)
			os.Exit(-1)
		}

		historyGrabber := internal.NewHistoryGrabber(db, stockCache)
		if err := historyGrabber.Run(); err != nil {
			logrus.Errorf("Starting HistoryGrabber had an issue: %v", err)
			os.Exit(-1)
		}

		intradayGrabber := internal.NewIntradayGrabber(db, stockCache)
		if err := intradayGrabber.Run(); err != nil {
			logrus.Errorf("Starting IntradayGrabber had an issue: %v", err)
			os.Exit(-1)
		}

		optionSymbolGrabber := internal.NewOptionSymbolGrabber(db, stockCache)
		if err := optionSymbolGrabber.Run(); err != nil {
			logrus.Errorf("Starting OptionSymbolGrabber had an issue: %v", err)
			os.Exit(-1)
		}

		quoteGrabber := internal.NewQuoteGrabber(db, stockCache)
		if err := quoteGrabber.Run(); err != nil {
			logrus.Errorf("Starting QuoteGrabber had an issue: %v", err)
			os.Exit(-1)
		}

		stableQuoteGrabber := internal.NewStableQuoteGrabber(db, stockCache)
		if err := stableQuoteGrabber.Run(); err != nil {
			logrus.Errorf("Starting StableQuoteGrabber had an issue: %v", err)
			os.Exit(-1)
		}

		stockValidator := internal.NewStockValidator(db, stockCache, optionSymbolGrabber)
		if err := stockValidator.Run(); err != nil {
			logrus.Errorf("Starting StockValidator had an issue: %v", err)
			os.Exit(-1)
		}

		serverDiedCtx, serverDied := context.WithCancel(context.Background())
		httpServer := internal.NewHttpsRestServer(db, cmd.Flag("address").Value.String(), cmd.Flag("publicCert").Value.String(), cmd.Flag("privateKey").Value.String(), stockValidator, serverDied)
		if err := httpServer.Run(); err != nil {
			logrus.Errorf("Starting HttpServer failed: %v", err)
			os.Exit(-1)
		}

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

		//now for each go routine we give upto a 1 min
		// we start with validator and work our way backwards
		stockValidator.Stop(1 * time.Minute)
		stableQuoteGrabber.Stop(1 * time.Minute)
		quoteGrabber.Stop(1 * time.Minute)
		optionSymbolGrabber.Stop(1 * time.Minute)
		intradayGrabber.Stop(1 * time.Minute)
		historyGrabber.Stop(1 * time.Minute)
		stockCache.Stop(1 * time.Minute)
		httpServer.Stop(1 * time.Minute)
		db.Close()

		os.Exit(0)
	},
}

func init() {
	rootCmd.AddCommand(httpsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// httpsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// httpsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	httpsCmd.Flags().String("address", "https://localhost:9090", "The address that SibylServer will be bind to for REST calls")
	httpsCmd.Flags().String("database", "localhost:3306", "The address for the MYSQL Server")
	httpsCmd.Flags().String("publicCert", "./server.crt", "The pathname to the public cert used for SSL")
	httpsCmd.Flags().String("privateKey", "./server.key", "The pathname to the PRIVATE key used for SSL")

	if err := httpsCmd.MarkFlagRequired("publicCert"); err != nil {
		panic("publicCert should be set to required")
	}
	if err := httpsCmd.MarkFlagRequired("privateKey"); err != nil {
		panic("privateKey should be set to required")
	}
}
