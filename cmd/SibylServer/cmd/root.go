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
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"time"

	"github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/spf13/cobra"
)

var logDirectory string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "SibylServer",
	Short: "The Sibyl suite backend tool for stock data acquisition, Display Only views and manual trade execution",
	Long:  `SibylServer is the Sibyl suite's backend server pursuant to being Display Only is used to acquire data from a discount broker, display data and perform trades.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {

	responseErrorLog := logDirectory + "/SibylServer.log"
	if _, err := os.Stat(logDirectory); os.IsNotExist(err) {
		if err := os.MkdirAll(logDirectory, 0775); err != nil {
			fmt.Printf("problems with creating the logging directory: %v", err)
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

	err := rootCmd.Execute()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&logDirectory, "logs", "./", "specifies the directory to store logs defaults to current directory")
}
