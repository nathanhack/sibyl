package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/nathanhack/sibyl/rest"
	"gopkg.in/resty.v1"
	"net/http"
	"strings"

	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Deletes one or many stocks",
	Long:  `Deletes one or many stocks`,

	RunE: func(cmd *cobra.Command, args []string) error {
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
			resp, err := resty.R().Delete(fmt.Sprintf("%v/stocks/%v", address, s))
			if err != nil {
				return fmt.Errorf("There was an error while deleting stock %v, error: %v\n", s, err)
			} else if resp.StatusCode() != http.StatusOK {
				return fmt.Errorf("There was an error while deleting stock %v, statusCode: %v  response: %v\n", s, resp.StatusCode(), resp)
			} else {
				var respErrors rest.ErrorState
				err := json.Unmarshal(resp.Body(), &respErrors)
				if err != nil {
					return fmt.Errorf("There was a problme parsing the server response: %v  had error:%v\n", string(resp.Body()), err)
				} else {
					if respErrors.ErrorReturned {
						return fmt.Errorf("There was a problem server side: %v\n", respErrors.ErrorReturned)
					} else {
						fmt.Printf("Successfullly deleted stock: %v\n", s)
					}
				}
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
