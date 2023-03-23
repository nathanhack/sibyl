package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/nathanhack/sibyl/ent/ogent"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add STOCK [STOCK] ...",
	Short: "Adds one or more stock by ticker to server",
	Long:  `The add command adds a one or more stocks`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("'add' requires at least one stock symbol to add")
		}

		spaceReplacer := strings.NewReplacer("  ", " ")
		toCommas := strings.NewReplacer(" ,", ",", ", ", ",", " ", ",", ",,", ",")
		stocks := make([]string, 0)

		noSpace := spaceReplacer.Replace(strings.Join(args, ","))
		commas := toCommas.Replace(noSpace)
		stocks = append(stocks, strings.Split(commas, ",")...)

		address, err := cmd.Flags().GetString("serverAddress")
		if err != nil {
			return fmt.Errorf("Could not get server address from passed in arguments: %v\n", err)
		}

		client, err := ogent.NewClient(address)
		if err != nil {
			return err
		}
		ctx := context.Background()

		for _, s := range stocks {
			err := client.AddTicker(ctx, ogent.AddTickerParams{
				Ticker: s,
			})
			if err != nil {
				fmt.Printf("Failed on stock: %v\n", s)
				return err
			}

			fmt.Printf("Successfully requested addition of stock: %v\n", s)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
