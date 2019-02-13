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

// disableCmd represents the disable command
var disableCmd = &cobra.Command{
	Use:   "disable",
	Short: "Disables a set of automated server actions",
	Long:  `Disables a set of actions the SibylServer can perform`,
}

var disableDownloadingCmd = &cobra.Command{
	Use:   "downloading",
	Short: "Disables Downloading for a particular stock",
	Long:  `Disables Downloading for a particular stock`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runDisableBlankCmd(cmd, args, "downloading", "downloading")
	},
}

var disableAllDownloadingCmd = &cobra.Command{
	Use:   "downloading",
	Short: "Disables Downloading for all stocks",
	Long:  `Disables Downloading for all stocks`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runDisableAllBlankCmd(cmd, "downloading", "downloading")
	},
}

var disableHistoryCmd = &cobra.Command{
	Use:   "history",
	Short: "Disables History for a particular stock (if downloading is enabled)",
	Long:  `Disables gathering daily History for a particular stock (if downloading is enabled)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runDisableBlankCmd(cmd, args, "history", "history")
	},
}

var disableAllHistoryCmd = &cobra.Command{
	Use:   "history",
	Short: "Disables History for all stocks (if downloading is enabled)",
	Long:  `Disables gathering daily History for all stocks (if downloading is enabled)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runDisableAllBlankCmd(cmd, "history", "history")
	},
}

var disableIntradayCmd = &cobra.Command{
	Use:   "intraday",
	Short: "Disables Intraday History for a particular stock",
	Long:  `Disables gathering 1 min resolution Intraday History for a particular stock`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runDisableBlankCmd(cmd, args, "intraday", "intraday")
	},
}

var disableAllIntradayCmd = &cobra.Command{
	Use:   "intraday",
	Short: "Disables Intraday History for all stocks",
	Long:  `Disables gathering 1 min resolution Intraday History for all stocks`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runDisableAllBlankCmd(cmd, "intraday", "intraday")
	},
}

var disableQuotesCmd = &cobra.Command{
	Use:   "quotes",
	Short: "Disables Quotes for a particular stock",
	Long:  `Disables gathering 1 min resolution Quotes for a particular stock`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runDisableBlankCmd(cmd, args, "quotes", "quotes")
	},
}

var disableAllQuotesCmd = &cobra.Command{
	Use:   "quotes",
	Short: "Disables Quotes for all stocks",
	Long:  `Disables gathering 1 min resolution Quotes for all stocks`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runDisableAllBlankCmd(cmd, "quotes", "quotes")
	},
}

var disableStableQuotesCmd = &cobra.Command{
	Use:   "stableQuotes",
	Short: "Disables Stable Quotes for a particular stock",
	Long:  `Disables gathering Stable Quotes every day for a particular stock`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runDisableBlankCmd(cmd, args, "stableQuotes", "stableQuotes")
	},
}

var disableAllStableQuotesCmd = &cobra.Command{
	Use:   "stableQuotes",
	Short: "Disables Stable Quotes for all stocks",
	Long:  `Disables gathering Stable Quotes every day for all stocks`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runDisableAllBlankCmd(cmd, "stableQuotes", "stableQuotes")
	},
}

func runDisableBlankCmd(cmd *cobra.Command, args []string, commandName, restEndpoint string) error {
	if len(args) == 0 {
		return fmt.Errorf("'%v' must have at least one stock symbol", commandName)
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
		resp, err := resty.R().Post(fmt.Sprintf("%v/stocks/disable/%v/%v", address, restEndpoint, s))
		if err != nil {
			return fmt.Errorf("There was an error while disabling %v for stock %v, error: %v\n", commandName, s, err)
		} else if resp.StatusCode() != http.StatusOK {
			return fmt.Errorf("There was an error while disabling %v for stock %v, statusCode: %v  response: %v\n", commandName, s, resp.StatusCode(), resp)
		} else {
			var respErrors rest.ErrorState
			err := json.Unmarshal(resp.Body(), &respErrors)
			if err != nil {
				return fmt.Errorf("There was a problem parsing the server response while disabling %v for stock %v : %v  had error:%v\n", commandName, s, string(resp.Body()), err)
			} else {
				if respErrors.ErrorReturned {
					return fmt.Errorf("There was a problem server side while disabling %v for stock %v : %v\n", commandName, s, respErrors.ErrorReturned)
				} else {
					fmt.Printf("Successfullly disabled %v for stock: %v\n", commandName, s)
				}
			}
		}
	}
	return nil
}

func runDisableAllBlankCmd(cmd *cobra.Command, commandName, restEndpoint string) error {
	address, err := cmd.Flags().GetString("serverAddress")
	if err != nil {
		return fmt.Errorf("Could not get server address from passed in arguments: %v\n", err)
	}

	resp, err := resty.R().Post(fmt.Sprintf("%v/stocks/disable/all/%v", address, restEndpoint))
	if err != nil {
		return fmt.Errorf("There was an error while disabling %v for all stocks, error: %v\n", commandName, err)
	} else if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("There was an error while disabling %v for all stocks, statusCode: %v  response: %v\n", commandName, resp.StatusCode(), resp)
	} else {
		var respErrors rest.ErrorState
		err := json.Unmarshal(resp.Body(), &respErrors)
		if err != nil {
			return fmt.Errorf("There was a problem parsing the server response while disabling %v for all stocks: %v  had error:%v\n", commandName, string(resp.Body()), err)
		} else {
			if respErrors.ErrorReturned {
				return fmt.Errorf("There was a problem server side while disabling %v for all stocks: %v\n", commandName, respErrors.ErrorReturned)
			} else {
				fmt.Printf("Successfullly disabled %v for all stocks\n", commandName)
			}
		}
	}
	return nil
}

var disableAllCmd = &cobra.Command{
	Use:   "all",
	Short: "Enables all attributes",
	Long:  `Enables all attributes`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			return fmt.Errorf("No args allowed for 'all'")
		}

		address, err := cmd.Flags().GetString("serverAddress")
		if err != nil {
			return fmt.Errorf("Could not get server address from passed in arguments: %v\n", err)
		}

		resp, err := resty.R().Put(fmt.Sprintf("%v/stocks/disable/all", address))
		if err != nil {
			return fmt.Errorf("There was an error while disabling all stocks error: %v\n", err)
		} else if resp.StatusCode() != http.StatusOK {
			return fmt.Errorf("There was an error while disabling all stocks, statusCode: %v  response: %v\n", resp.StatusCode(), resp)
		} else {
			var respErrors rest.ErrorState
			err := json.Unmarshal(resp.Body(), &respErrors)
			if err != nil {
				return fmt.Errorf("There was a problem parsing the server response while disabling all stocks: %v  had error:%v\n", string(resp.Body()), err)
			} else {
				if respErrors.ErrorReturned {
					return fmt.Errorf("There was a problem server side while disabling all stocks: %v\n", respErrors.Error)
				} else {
					fmt.Printf("Successfullly disabled all stocks\n")
				}
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(disableCmd)
	disableCmd.AddCommand(disableDownloadingCmd)
	disableCmd.AddCommand(disableHistoryCmd)
	disableCmd.AddCommand(disableIntradayCmd)
	disableCmd.AddCommand(disableQuotesCmd)
	disableCmd.AddCommand(disableStableQuotesCmd)

	disableCmd.AddCommand(disableAllCmd)
	disableAllCmd.AddCommand(disableAllDownloadingCmd)
	disableAllCmd.AddCommand(disableAllHistoryCmd)
	disableAllCmd.AddCommand(disableAllIntradayCmd)
	disableAllCmd.AddCommand(disableAllQuotesCmd)
	disableAllCmd.AddCommand(disableAllStableQuotesCmd)

}
