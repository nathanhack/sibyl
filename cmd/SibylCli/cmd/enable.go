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

// enableCmd represents the enable command
var enableCmd = &cobra.Command{
	Use:   "enable",
	Short: "Enables a set of automated server actions",
	Long:  `Enables a set of actions the SibylServer will perform`,
}

var enableDownloadingCmd = &cobra.Command{
	Use:   "downloading",
	Short: "Enables downloading for a particular stock",
	Long:  `Enables downloading for a particular stock`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runEnableBlankCmd(cmd, args, "downloading", "downloading")
	},
}

var enableAllDownloadingCmd = &cobra.Command{
	Use:   "downloading",
	Short: "Enables downloading for all stocks",
	Long:  `Enables downloading for all stocks`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runEnableAllBlankCmd(cmd, "downloading", "downloading")
	},
}

var enableHistoryCmd = &cobra.Command{
	Use:   "history",
	Short: "Enables history for a particular stock (if downloading is enabled)",
	Long:  `Enables gathering daily History for a particular stock (if downloading is enabled)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runEnableBlankCmd(cmd, args, "history", "history")
	},
}

var enableAllHistoryCmd = &cobra.Command{
	Use:   "history",
	Short: "Enables history for a all stocks (if downloading is enabled)",
	Long:  `Enables gathering daily History for all stocks (if downloading is enabled)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runEnableAllBlankCmd(cmd, "history", "history")
	},
}

var enableIntradayCmd = &cobra.Command{
	Use:   "intraday",
	Short: "Enables Intraday History for a particular stock (if downloading is enabled)",
	Long:  `Enables gathering 1 min resolution Intraday History for a particular stock (if downloading is enabled)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runEnableBlankCmd(cmd, args, "intraday", "intraday")
	},
}

var enableAllIntradayCmd = &cobra.Command{
	Use:   "intraday",
	Short: "Enables Intraday History for all stocks (if downloading is enabled)",
	Long:  `Enables gathering 1 min resolution Intraday History for all stocks (if downloading is enabled)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runEnableAllBlankCmd(cmd, "intraday", "intraday")
	},
}

var enableQuotesCmd = &cobra.Command{
	Use:   "quotes",
	Short: "Enables Quotes for a particular stock (if downloading is enabled)",
	Long:  `Enables gathering 1 min resolution Quotes for a particular stock (if downloading is enabled)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runEnableBlankCmd(cmd, args, "quotes", "quotes")
	},
}

var enableAllQuotesCmd = &cobra.Command{
	Use:   "quotes",
	Short: "Enables Quotes for all stocks (if downloading is enabled)",
	Long:  `Enables gathering 1 min resolution Quotes for all stocks (if downloading is enabled)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runEnableAllBlankCmd(cmd, "quotes", "quotes")
	},
}

var enableStableQuotesCmd = &cobra.Command{
	Use:   "stableQuotes",
	Short: "Enables Stable Quotes for a particular stock (if downloading is enabled)",
	Long:  `Enables gathering Stable Quotes every day for a particular stock (if downloading is enabled)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runEnableBlankCmd(cmd, args, "stableQuotes", "stableQuotes")
	},
}

var enableAllStableQuotesCmd = &cobra.Command{
	Use:   "stableQuotes",
	Short: "Enables Stable Quotes for all stocks (if downloading is enabled)",
	Long:  `Enables gathering Stable Quotes every day for all stocks (if downloading is enabled)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runEnableBlankCmd(cmd, args, "stableQuotes", "stableQuotes")
	},
}

func runEnableBlankCmd(cmd *cobra.Command, args []string, commandName, restEndpoint string) error {
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
		resp, err := resty.R().Put(fmt.Sprintf("%v/stocks/enable/%v/%v", address, restEndpoint, s))
		if err != nil {
			return fmt.Errorf("There was an error while enabling %v for stock %v, error: %v\n", commandName, s, err)
		} else if resp.StatusCode() != http.StatusOK {
			return fmt.Errorf("There was an error while enabling %v for stock %v, statusCode: %v  response: %v\n", commandName, s, resp.StatusCode(), resp)
		} else {
			var respErrors rest.ErrorState
			err := json.Unmarshal(resp.Body(), &respErrors)
			if err != nil {
				return fmt.Errorf("There was a problem parsing the server response while enabling %v for stock %v : %v  had error:%v\n", commandName, s, string(resp.Body()), err)
			} else {
				if respErrors.ErrorReturned {
					return fmt.Errorf("There was a problem server side while enabling %v for stock %v : %v\n", commandName, s, respErrors.ErrorReturned)
				} else {
					fmt.Printf("Successfullly enabled %v for stock: %v\n", commandName, s)
				}
			}
		}
	}
	return nil
}

func runEnableAllBlankCmd(cmd *cobra.Command, commandName, restEndpoint string) error {
	address, err := cmd.Flags().GetString("serverAddress")
	if err != nil {
		return fmt.Errorf("Could not get server address from passed in arguments: %v\n", err)
	}

	resp, err := resty.R().Put(fmt.Sprintf("%v/stocks/enable/all/%v", address, restEndpoint))
	if err != nil {
		return fmt.Errorf("There was an error while enabling %v for all stocks, error: %v\n", commandName, err)
	} else if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("There was an error while enabling %v for all stocks, statusCode: %v  response: %v\n", commandName, resp.StatusCode(), resp)
	} else {
		var respErrors rest.ErrorState
		err := json.Unmarshal(resp.Body(), &respErrors)
		if err != nil {
			return fmt.Errorf("There was a problem parsing the server response while enabling %v for all stocks: %v  had error:%v\n", commandName, string(resp.Body()), err)
		} else {
			if respErrors.ErrorReturned {
				return fmt.Errorf("There was a problem server side while enabling %v for all stocks: %v\n", commandName, respErrors.ErrorReturned)
			} else {
				fmt.Printf("Successfullly enabled %v for all stocks\n", commandName)
			}
		}
	}

	return nil
}

var enableAllCmd = &cobra.Command{
	Use:   "all",
	Short: "Enables all attributes",
	Long:  `Enables all attributes`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := cobra.NoArgs(cmd, args); err != nil {
			return err
		}

		address, err := cmd.Flags().GetString("serverAddress")
		if err != nil {
			return fmt.Errorf("Could not get server address from passed in arguments: %v\n", err)
		}

		resp, err := resty.R().Put(fmt.Sprintf("%v/stocks/enable/all", address))
		if err != nil {
			return fmt.Errorf("There was an error while enabling all stocks error: %v\n", err)
		} else if resp.StatusCode() != http.StatusOK {
			return fmt.Errorf("There was an error while enabling all stocks, statusCode: %v  response: %v\n", resp.StatusCode(), resp)
		} else {
			var respErrors rest.ErrorState
			err := json.Unmarshal(resp.Body(), &respErrors)
			if err != nil {
				return fmt.Errorf("There was a problem parsing the server response while enabling all stocks: %v  had error:%v\n", string(resp.Body()), err)
			} else {
				if respErrors.ErrorReturned {
					return fmt.Errorf("There was a problem server side while enabling all stocks: %v\n", respErrors.Error)
				} else {
					fmt.Printf("Successfullly enabled all stocks\n")
				}
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(enableCmd)
	enableCmd.AddCommand(enableDownloadingCmd)
	enableCmd.AddCommand(enableHistoryCmd)
	enableCmd.AddCommand(enableIntradayCmd)
	enableCmd.AddCommand(enableQuotesCmd)
	enableCmd.AddCommand(enableStableQuotesCmd)
	enableCmd.AddCommand(enableAllCmd)

	enableAllCmd.AddCommand(enableAllDownloadingCmd)
	enableAllCmd.AddCommand(enableAllHistoryCmd)
	enableAllCmd.AddCommand(enableAllIntradayCmd)
	enableAllCmd.AddCommand(enableAllQuotesCmd)
	enableAllCmd.AddCommand(enableAllStableQuotesCmd)
}
