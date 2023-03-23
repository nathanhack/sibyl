package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var serverAddress string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "SibylCli",
	Short: "SibylCli is a commandline tool to interface with SibylServer",
	Long:  `SibylCli is the Sibyl suite's commandline tool used to manage the Sibyl backend server SibylServer.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&serverAddress,"serverAddress", "a", "http://localhost:9090", "sets the address for the SibylServer")
}
