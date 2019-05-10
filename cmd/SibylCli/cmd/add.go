package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/nathanhack/sibyl/rest"
	"github.com/spf13/cobra"
	"gopkg.in/resty.v1"
	"net/http"
	"strings"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add STOCK [STOCK] ...",
	Short: "Adds a stock to analyze and trade with",
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

		for _, s := range stocks {
			resp, err := resty.R().Post(fmt.Sprintf("%v/stocks/add/%v", address, s))
			if err != nil {
				return fmt.Errorf("There was an error while adding stock %v, error: %v\n", s, err)
			}
			if resp.StatusCode() != http.StatusOK {
				return fmt.Errorf("There was an error while adding stock %v, statusCode: %v  response: %v\n", s, resp.StatusCode(), resp)
			}

			var respErrors rest.ErrorState
			err = json.Unmarshal(resp.Body(), &respErrors)
			if err != nil {
				return fmt.Errorf("There was a problem parsing the server response while adding stock %v : %v  had error:%v\n", s, string(resp.Body()), err)
			}

			if respErrors.ErrorReturned {
				return fmt.Errorf("There was a problem server side while adding stock %s: %v\n", s, respErrors.ErrorReturned)
			}

			fmt.Printf("Successfullly added stock: %v\n", s)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
