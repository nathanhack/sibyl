package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/nathanhack/sibyl/ent/ogent"
	"github.com/spf13/cobra"
)

var (
	disableAll          bool
	disableDataSourceID int
)

// enableCmd represents the enable command
var disableCmd = &cobra.Command{
	Use:   "disable",
	Short: "Disables a set of automated server actions",
	Long:  `Disables a set of actions the SibylServer will perform`,
}

var disableBarCmd = &cobra.Command{
	Use:   "bar",
	Short: "Disables bar history",
	Long:  `Disables gathering bar history`,
}

var disableBarDailyCmd = &cobra.Command{
	Use:   "daily STOCK [STOCK] ...",
	Short: "Disables daily history for a particular stock (if downloading is enabled)",
	Long:  `Disables gathering daily history for a particular stock (if downloading is enabled)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return disableBars(args, ogent.EntityIntervalsListIntervalDaily)
	},
}

var disableBarMonthlyCmd = &cobra.Command{
	Use:   "monthly STOCK [STOCK] ...",
	Short: "Disables monthly history for a particular stock (if downloading is enabled)",
	Long:  `Disables gathering monthly History for a particular stock (if downloading is enabled)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("must have at least one stock symbol")
		}
		stocks := parserStockList(args)
		return disableBars(stocks, ogent.EntityIntervalsListIntervalMonthly)
	},
}

var disableBarYearlyCmd = &cobra.Command{
	Use:   "yearly STOCK [STOCK] ...",
	Short: "Disables yearly history for a particular stock (if downloading is enabled)",
	Long:  `Disables gathering yearly history for a particular stock (if downloading is enabled)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("must have at least one stock symbol")
		}
		stocks := parserStockList(args)
		return disableBars(stocks, ogent.EntityIntervalsListIntervalYearly)
	},
}

var disableBar1MinCmd = &cobra.Command{
	Use:   "1min STOCK [STOCK] ...",
	Short: "Disables 1 minute intraday for a particular stock(s)",
	Long:  `Disables gathering 1 minute resolution intraday bars for a particular stock(s)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("must have at least one stock symbol")
		}
		stocks := parserStockList(args)
		return disableBars(stocks, ogent.EntityIntervalsListInterval1min)
	},
}

func disableBars(args []string, stockInterval ogent.EntityIntervalsListInterval) error {
	client, err := ogent.NewClient(serverAddress)
	if err != nil {
		return err
	}
	ctx := context.Background()

	var stocks []string
	if !disableAll {
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

	dataSourceToUse := make([]int, 0)
	if disableDataSourceID == 0 {
		dataSources, err := getAllDataSources(ctx, client)
		if err != nil {
			return fmt.Errorf("error while getting data sources:%v", err)
		}
		if len(dataSources) == 0 {
			return fmt.Errorf("no DataSources available")
		}

		for _, d := range dataSources {
			dataSourceToUse = append(dataSourceToUse, d.ID)
		}
	} else {
		_, err := getDataSource(ctx, client, disableDataSourceID)
		if err != nil {
			return fmt.Errorf("error while getting data sources:%v", err)
		}
		dataSourceToUse = append(dataSourceToUse, disableDataSourceID)
	}

	for _, s := range stocks {
		ticker := strings.ToUpper(s)
		for _, dataSource := range dataSourceToUse {
			err := addOrModifyInterval(ctx, ticker, stockInterval, dataSource, client, false)
			if err != nil {
				return err
			}
		}

		fmt.Printf("Successfully disabled %v bars for stock: %v\n", stockInterval, s)
	}
	return nil
}

func init() {
	rootCmd.AddCommand(disableCmd)
	disableCmd.AddCommand(disableBarCmd)

	disableBarCmd.AddCommand(disableBar1MinCmd)
	disableBarCmd.AddCommand(disableBarDailyCmd)
	disableBarCmd.AddCommand(disableBarMonthlyCmd)
	disableBarCmd.AddCommand(disableBarYearlyCmd)

	disableBarCmd.PersistentFlags().BoolVarP(&disableAll, "all", "", false, "When used it will apply to all active stocks in the database")
	disableBarCmd.PersistentFlags().IntVarP(&disableDataSourceID, "datasource", "d", 0, "The DataSourceID to use when disabling (defaults to 0 all datasources)")
}
