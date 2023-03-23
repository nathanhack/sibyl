package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-faster/errors"
	"github.com/nathanhack/sibyl/ent/ogent"

	"github.com/spf13/cobra"
)

var (
	enableAll          bool
	enableDataSourceID int
)

// enableCmd represents the enable command
var enableCmd = &cobra.Command{
	Use:   "enable",
	Short: "Enables a set of automated server actions",
	Long:  `Enables a set of actions the SibylServer will perform`,
}

var enableBarCmd = &cobra.Command{
	Use:   "bar",
	Short: "Enables bar history",
	Long:  `Enables gathering bar history`,
}

var enableBarDailyCmd = &cobra.Command{
	Use:   "daily STOCK [STOCK] ...",
	Short: "Enables daily history for a particular stock (if downloading is enabled)",
	Long:  `Enables gathering daily history for a particular stock (if downloading is enabled)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return enableBars(args, ogent.EntityIntervalsListIntervalDaily)
	},
}

var enableBarMonthlyCmd = &cobra.Command{
	Use:   "monthly STOCK [STOCK] ...",
	Short: "Enables monthly history for a particular stock (if downloading is enabled)",
	Long:  `Enables gathering monthly History for a particular stock (if downloading is enabled)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return enableBars(args, ogent.EntityIntervalsListIntervalMonthly)
	},
}

var enableBarYearlyCmd = &cobra.Command{
	Use:   "yearly STOCK [STOCK] ...",
	Short: "Enables yearly history for a particular stock (if downloading is enabled)",
	Long:  `Enables gathering yearly history for a particular stock (if downloading is enabled)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return enableBars(args, ogent.EntityIntervalsListIntervalYearly)
	},
}

var enableBar1MinCmd = &cobra.Command{
	Use:   "1min STOCK [STOCK] ...",
	Short: "Enables 1 minute intraday for a particular stock(s)",
	Long:  `Enables gathering 1 minute resolution intraday bars for a particular stock(s)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return enableBars(args, ogent.EntityIntervalsListInterval1min)
	},
}

func parserStockList(args []string) []string {
	spaceReplacer := strings.NewReplacer("  ", " ")
	toCommas := strings.NewReplacer(" ,", ",", ", ", ",", " ", ",", ",,", ",")

	noSpace := spaceReplacer.Replace(strings.Join(args, ","))
	commas := toCommas.Replace(noSpace)
	return strings.Split(commas, ",")
}

func enableBars(args []string, stockInterval ogent.EntityIntervalsListInterval) error {
	client, err := ogent.NewClient(serverAddress)
	if err != nil {
		return err
	}
	ctx := context.Background()

	var stocks []string
	if !enableAll {
		if len(args) == 0 {
			return fmt.Errorf("must have at least one stock symbol")
		}
		stocks = parserStockList(args)
	} else {
		entities, err := getAllEntities(ctx, client)
		if err != nil {
			return err
		}
		stocks = make([]string, len(entities))
		for i := range entities {
			stocks[i] = entities[i].Ticker
		}
	}

	if enableDataSourceID == 0 {
		dataSources, err := getAllDataSources(ctx, client)
		if err != nil {
			return fmt.Errorf("error while getting data sources:%v", err)
		}
		if len(dataSources) == 0 {
			return fmt.Errorf("no DataSources available")
		}

		if len(dataSources) > 1 {
			return fmt.Errorf("there are too many DataSources available one must be chosen. (use 'show datasources')")
		}
		enableDataSourceID = dataSources[0].ID
	} else {
		_, err := getDataSource(ctx, client, enableDataSourceID)
		if err != nil {
			return errors.Wrap(err, "error while getting data sources")
		}
	}

	for _, s := range stocks {
		ticker := strings.ToUpper(s)
		err := addOrModifyInterval(ctx, ticker, stockInterval, enableDataSourceID, client, true)
		if err != nil {
			return err
		}
		fmt.Printf("Successfully enabled %v bars for stock: %v\n", stockInterval, s)
	}

	return nil
}

func init() {
	rootCmd.AddCommand(enableCmd)
	enableCmd.AddCommand(enableBarCmd)

	enableBarCmd.AddCommand(enableBar1MinCmd)
	enableBarCmd.AddCommand(enableBarDailyCmd)
	enableBarCmd.AddCommand(enableBarMonthlyCmd)
	enableBarCmd.AddCommand(enableBarYearlyCmd)

	enableBarCmd.PersistentFlags().BoolVarP(&enableAll, "all", "", false, "When used it will apply to all active stocks in the database")
	enableBarCmd.PersistentFlags().IntVarP(&enableDataSourceID, "datasource", "d", 0, "The DataSourceID to use when enabling (defaults to 0 -- use whatever is available)")
}
