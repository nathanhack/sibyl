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

var enableAllCmd = &cobra.Command{
	Use:   "all",
	Short: "Enables all attributes",
	Long:  `Enables all attributes`,
}

var enableAllAllCmd = &cobra.Command{
	Use:   "all",
	Short: "Enables all things for all stocks",
	Long:  `Enables all things for all stocks`,
	RunE: func(cmd *cobra.Command, args []string) error {
		address, err := cmd.Flags().GetString("serverAddress")
		if err != nil {
			return fmt.Errorf("Could not get server address from passed in arguments: %v\n", err)
		}

		resp, err := resty.R().Put(fmt.Sprintf("%v/stocks/enable/all", address))
		if err != nil {
			return fmt.Errorf("There was an error while enabling all things for all stocks, error: %v\n", err)
		}

		if resp.StatusCode() != http.StatusOK {
			return fmt.Errorf("There was an error while enabling all things for all stocks, statusCode: %v  response: %v\n", resp.StatusCode(), resp)
		}

		var respErrors rest.ErrorState
		err = json.Unmarshal(resp.Body(), &respErrors)
		if err != nil {
			return fmt.Errorf("There was a problem parsing the server response while enabling all things for all stocks: %v  had error:%v\n", string(resp.Body()), err)
		}

		if respErrors.ErrorReturned {
			return fmt.Errorf("There was a problem server side while enabling all things for all stocks: %v\n", respErrors.ErrorReturned)
		}
		fmt.Printf("Successfullly enabled all things for all stocks\n")

		return nil
	},
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
	Short: "Enables history (if downloading is enabled)",
	Long:  `Enables gathering History (if downloading is enabled)`,
}

var enableHistoryAllCmd = &cobra.Command{
	Use:   "all STOCK [STOCK] ...",
	Short: "Enables daily, weekly, monthly, and yearly histories for a particular stock (if downloading is enabled)",
	Long:  `Enables gathering daily, weekly, monthly, and yearly histories for a particular stock (if downloading is enabled)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("must have at least one stock symbol")
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
			resp, err := resty.R().Put(fmt.Sprintf("%v/stocks/enable/history/%v/all", address, s))
			if err != nil {
				return fmt.Errorf("There was an error while enabling all histories for stock %v, error: %v\n", s, err)
			}

			if resp.StatusCode() != http.StatusOK {
				return fmt.Errorf("There was an error while enabling all histories for stock %v, statusCode: %v  response: %v\n", s, resp.StatusCode(), resp)
			}

			var respErrors rest.ErrorState
			err = json.Unmarshal(resp.Body(), &respErrors)
			if err != nil {
				return fmt.Errorf("There was a problem parsing the server response while enabling all histories for stock %v : %v  had error:%v\n", s, string(resp.Body()), err)
			}

			if respErrors.ErrorReturned {
				return fmt.Errorf("There was a problem server side while enabling all histories for stock %v : %v\n", s, respErrors.ErrorReturned)
			}

			fmt.Printf("Successfullly enabled all histories for stock: %v\n", s)
		}
		return nil
	},
}

var enableHistoryDailyCmd = &cobra.Command{
	Use:   "daily STOCK [STOCK] ...",
	Short: "Enables daily history for a particular stock (if downloading is enabled)",
	Long:  `Enables gathering daily history for a particular stock (if downloading is enabled)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("must have at least one stock symbol")
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
			resp, err := resty.R().Put(fmt.Sprintf("%v/stocks/enable/history/%v/daily", address, s))
			if err != nil {
				return fmt.Errorf("There was an error while enabling daily histories for stock %v, error: %v\n", s, err)
			}

			if resp.StatusCode() != http.StatusOK {
				return fmt.Errorf("There was an error while enabling daily histories for stock %v, statusCode: %v  response: %v\n", s, resp.StatusCode(), resp)
			}

			var respErrors rest.ErrorState
			err = json.Unmarshal(resp.Body(), &respErrors)
			if err != nil {
				return fmt.Errorf("There was a problem parsing the server response while enabling daily histories for stock %v : %v  had error:%v\n", s, string(resp.Body()), err)
			}

			if respErrors.ErrorReturned {
				return fmt.Errorf("There was a problem server side while enabling daily histories for stock %v : %v\n", s, respErrors.ErrorReturned)
			}

			fmt.Printf("Successfullly enabled daily histories for stock: %v\n", s)
		}
		return nil
	},
}

var enableHistoryWeeklyCmd = &cobra.Command{
	Use:   "weekly STOCK [STOCK] ...",
	Short: "Enables weekly history for a particular stock (if downloading is enabled)",
	Long:  `Enables gathering weekly history for a particular stock (if downloading is enabled)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("must have at least one stock symbol")
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
			resp, err := resty.R().Put(fmt.Sprintf("%v/stocks/enable/history/%v/weekly", address, s))
			if err != nil {
				return fmt.Errorf("There was an error while enabling weekly histories for stock %v, error: %v\n", s, err)
			}

			if resp.StatusCode() != http.StatusOK {
				return fmt.Errorf("There was an error while enabling weekly histories for stock %v, statusCode: %v  response: %v\n", s, resp.StatusCode(), resp)
			}

			var respErrors rest.ErrorState
			err = json.Unmarshal(resp.Body(), &respErrors)
			if err != nil {
				return fmt.Errorf("There was a problem parsing the server response while enabling weekly histories for stock %v : %v  had error:%v\n", s, string(resp.Body()), err)
			}

			if respErrors.ErrorReturned {
				return fmt.Errorf("There was a problem server side while enabling weekly histories for stock %v : %v\n", s, respErrors.ErrorReturned)
			}

			fmt.Printf("Successfullly enabled weekly histories for stock: %v\n", s)
		}
		return nil
	},
}
var enableHistoryMonthlyCmd = &cobra.Command{
	Use:   "monthly STOCK [STOCK] ...",
	Short: "Enables monthly history for a particular stock (if downloading is enabled)",
	Long:  `Enables gathering monthly History for a particular stock (if downloading is enabled)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("must have at least one stock symbol")
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
			resp, err := resty.R().Put(fmt.Sprintf("%v/stocks/enable/history/%v/monthly", address, s))
			if err != nil {
				return fmt.Errorf("There was an error while enabling monthly histories for stock %v, error: %v\n", s, err)
			}

			if resp.StatusCode() != http.StatusOK {
				return fmt.Errorf("There was an error while enabling monthly histories for stock %v, statusCode: %v  response: %v\n", s, resp.StatusCode(), resp)
			}

			var respErrors rest.ErrorState
			err = json.Unmarshal(resp.Body(), &respErrors)
			if err != nil {
				return fmt.Errorf("There was a problem parsing the server response while enabling monthly histories for stock %v : %v  had error:%v\n", s, string(resp.Body()), err)
			}

			if respErrors.ErrorReturned {
				return fmt.Errorf("There was a problem server side while enabling monthly histories for stock %v : %v\n", s, respErrors.ErrorReturned)
			}

			fmt.Printf("Successfullly enabled monthly histories for stock: %v\n", s)
		}
		return nil
	},
}
var enableHistoryYearlyCmd = &cobra.Command{
	Use:   "yearly STOCK [STOCK] ...",
	Short: "Enables yearly history for a particular stock (if downloading is enabled)",
	Long:  `Enables gathering yearly history for a particular stock (if downloading is enabled)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("must have at least one stock symbol")
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
			resp, err := resty.R().Put(fmt.Sprintf("%v/stocks/enable/history/%v/yearly", address, s))
			if err != nil {
				return fmt.Errorf("There was an error while enabling yearly histories for stock %v, error: %v\n", s, err)
			}

			if resp.StatusCode() != http.StatusOK {
				return fmt.Errorf("There was an error while enabling yearly histories for stock %v, statusCode: %v  response: %v\n", s, resp.StatusCode(), resp)
			}

			var respErrors rest.ErrorState
			err = json.Unmarshal(resp.Body(), &respErrors)
			if err != nil {
				return fmt.Errorf("There was a problem parsing the server response while enabling yearly histories for stock %v : %v  had error:%v\n", s, string(resp.Body()), err)
			}

			if respErrors.ErrorReturned {
				return fmt.Errorf("There was a problem server side while enabling yearly histories for stock %v : %v\n", s, respErrors.ErrorReturned)
			}

			fmt.Printf("Successfullly enabled yearly histories for stock: %v\n", s)
		}
		return nil
	},
}

var enableAllHistoryCmd = &cobra.Command{
	Use:   "history",
	Short: "Enables history for a all stocks (if downloading is enabled)",
	Long:  `Enables gathering history for all stocks (if downloading is enabled)`,
}

var enableAllHistoryAllCmd = &cobra.Command{
	Use:   "all",
	Short: "Enables daily, weekly, monthly, and yearly histories for a all stocks (if downloading is enabled)",
	Long:  `Enables gathering daily, weekly, monthly, and yearly histories for all stocks (if downloading is enabled)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		address, err := cmd.Flags().GetString("serverAddress")
		if err != nil {
			return fmt.Errorf("Could not get server address from passed in arguments: %v\n", err)
		}

		resp, err := resty.R().Put(fmt.Sprintf("%v/stocks/enable/all/history/all", address))
		if err != nil {
			return fmt.Errorf("There was an error while enabling all histories for all stocks, error: %v\n", err)
		}

		if resp.StatusCode() != http.StatusOK {
			return fmt.Errorf("There was an error while enabling all histories for all stocks, statusCode: %v  response: %v\n", resp.StatusCode(), resp)
		}

		var respErrors rest.ErrorState
		err = json.Unmarshal(resp.Body(), &respErrors)
		if err != nil {
			return fmt.Errorf("There was a problem parsing the server response while enabling all histories for all stocks: %v  had error:%v\n", string(resp.Body()), err)
		}

		if respErrors.ErrorReturned {
			return fmt.Errorf("There was a problem server side while enabling all histories for all stocks: %v\n", respErrors.ErrorReturned)
		}
		fmt.Printf("Successfullly enabled all histories for all stocks\n")

		return nil
	},
}

var enableAllHistoryDailyCmd = &cobra.Command{
	Use:   "daily",
	Short: "Enables daily history for a all stocks (if downloading is enabled)",
	Long:  `Enables gathering daily history for all stocks (if downloading is enabled)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		address, err := cmd.Flags().GetString("serverAddress")
		if err != nil {
			return fmt.Errorf("Could not get server address from passed in arguments: %v\n", err)
		}

		resp, err := resty.R().Put(fmt.Sprintf("%v/stocks/enable/all/history/daily", address))
		if err != nil {
			return fmt.Errorf("There was an error while enabling daily histories for all stocks, error: %v\n", err)
		}

		if resp.StatusCode() != http.StatusOK {
			return fmt.Errorf("There was an error while enabling daily histories for all stocks, statusCode: %v  response: %v\n", resp.StatusCode(), resp)
		}

		var respErrors rest.ErrorState
		err = json.Unmarshal(resp.Body(), &respErrors)
		if err != nil {
			return fmt.Errorf("There was a problem parsing the server response while enabling daily histories for all stocks: %v  had error:%v\n", string(resp.Body()), err)
		}

		if respErrors.ErrorReturned {
			return fmt.Errorf("There was a problem server side while enabling daily histories for all stocks: %v\n", respErrors.ErrorReturned)
		}
		fmt.Printf("Successfullly enabled daily histories for all stocks\n")

		return nil
	},
}

var enableAllHistoryWeeklyCmd = &cobra.Command{
	Use:   "weekly",
	Short: "Enables weekly history for a all stocks (if downloading is enabled)",
	Long:  `Enables gathering weekly history for all stocks (if downloading is enabled)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		address, err := cmd.Flags().GetString("serverAddress")
		if err != nil {
			return fmt.Errorf("Could not get server address from passed in arguments: %v\n", err)
		}

		resp, err := resty.R().Put(fmt.Sprintf("%v/stocks/enable/all/history/weekly", address))
		if err != nil {
			return fmt.Errorf("There was an error while enabling weekly histories for all stocks, error: %v\n", err)
		}

		if resp.StatusCode() != http.StatusOK {
			return fmt.Errorf("There was an error while enabling weekly histories for all stocks, statusCode: %v  response: %v\n", resp.StatusCode(), resp)
		}

		var respErrors rest.ErrorState
		err = json.Unmarshal(resp.Body(), &respErrors)
		if err != nil {
			return fmt.Errorf("There was a problem parsing the server response while enabling weekly histories for all stocks: %v had error:%v\n", string(resp.Body()), err)
		}

		if respErrors.ErrorReturned {
			return fmt.Errorf("There was a problem server side while enabling weekly histories for all stocks: %v\n", respErrors.ErrorReturned)
		}
		fmt.Printf("Successfullly enabled weekly histories for all stocks\n")

		return nil
	},
}

var enableAllHistoryMonthlyCmd = &cobra.Command{
	Use:   "monthly",
	Short: "Enables monthly history for a all stocks (if downloading is enabled)",
	Long:  `Enables gathering monthly history for all stocks (if downloading is enabled)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		address, err := cmd.Flags().GetString("serverAddress")
		if err != nil {
			return fmt.Errorf("Could not get server address from passed in arguments: %v\n", err)
		}

		resp, err := resty.R().Put(fmt.Sprintf("%v/stocks/enable/all/history/monthly", address))
		if err != nil {
			return fmt.Errorf("There was an error while enabling monthly histories for all stocks, error: %v\n", err)
		}

		if resp.StatusCode() != http.StatusOK {
			return fmt.Errorf("There was an error while enabling monthly histories for all stocks, statusCode: %v  response: %v\n", resp.StatusCode(), resp)
		}

		var respErrors rest.ErrorState
		err = json.Unmarshal(resp.Body(), &respErrors)
		if err != nil {
			return fmt.Errorf("There was a problem parsing the server response while enabling monthly histories for all stocks: %v had error:%v\n", string(resp.Body()), err)
		}

		if respErrors.ErrorReturned {
			return fmt.Errorf("There was a problem server side while enabling monthly histories for all stocks: %v\n", respErrors.ErrorReturned)
		}
		fmt.Printf("Successfullly enabled monthly histories for all stocks\n")

		return nil
	},
}

var enableAllHistoryYearlyCmd = &cobra.Command{
	Use:   "yearly",
	Short: "Enables yearly history for a all stocks (if downloading is enabled)",
	Long:  `Enables gathering yearly history for all stocks (if downloading is enabled)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		address, err := cmd.Flags().GetString("serverAddress")
		if err != nil {
			return fmt.Errorf("Could not get server address from passed in arguments: %v\n", err)
		}

		resp, err := resty.R().Put(fmt.Sprintf("%v/stocks/enable/all/history/yearly", address))
		if err != nil {
			return fmt.Errorf("There was an error while enabling yearly histories for all stocks, error: %v\n", err)
		}

		if resp.StatusCode() != http.StatusOK {
			return fmt.Errorf("There was an error while enabling yearly histories for all stocks, statusCode: %v  response: %v\n", resp.StatusCode(), resp)
		}

		var respErrors rest.ErrorState
		err = json.Unmarshal(resp.Body(), &respErrors)
		if err != nil {
			return fmt.Errorf("There was a problem parsing the server response while enabling yearly histories for all stocks: %v had error:%v\n", string(resp.Body()), err)
		}

		if respErrors.ErrorReturned {
			return fmt.Errorf("There was a problem server side while enabling yearly histories for all stocks: %v\n", respErrors.ErrorReturned)
		}
		fmt.Printf("Successfullly enabled yearly histories for all stocks\n")

		return nil
	},
}

var enableIntradayCmd = &cobra.Command{
	Use:   "intraday",
	Short: "Intraday history gathering (if downloading is enabled)",
	Long:  `Intraday history gathering (if downloading is enabled)`,
}

var enableIntradayAllCmd = &cobra.Command{
	Use:   "all STOCK [STOCK] ...",
	Short: "Enables All Intraday Histories for a particular stock(s) (if downloading is enabled)",
	Long:  `Enables gathering Ticks, 1 min, 5 min resolution Intraday Histories for a particular stock(s) (if downloading is enabled)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("must have at least one stock symbol")
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
			resp, err := resty.R().Put(fmt.Sprintf("%v/stocks/enable/intraday/%v/all", address, s))
			if err != nil {
				return fmt.Errorf("There was an error while enabling all intraday histories for stock %v, error: %v\n", s, err)
			}

			if resp.StatusCode() != http.StatusOK {
				return fmt.Errorf("There was an error while enabling all intraday histories for stock %v, statusCode: %v  response: %v\n", s, resp.StatusCode(), resp)
			}

			var respErrors rest.ErrorState
			err = json.Unmarshal(resp.Body(), &respErrors)
			if err != nil {
				return fmt.Errorf("There was a problem parsing the server response while enabling all intraday histories for stock %v : %v had error:%v\n", s, string(resp.Body()), err)
			}

			if respErrors.ErrorReturned {
				return fmt.Errorf("There was a problem server side while enabling all intraday histories for stock %v : %v\n", s, respErrors.ErrorReturned)
			}

			fmt.Printf("Successfullly enabled all intraday histories for stock: %v\n", s)
		}
		return nil
	},
}

var enableIntradayTickCmd = &cobra.Command{
	Use:   "tick STOCK [STOCK] ...",
	Short: "Enables tick intraday for a particular stock(s) (if downloading is enabled)",
	Long:  `Enables gathering tick resolution intraday for a particular stock(s) (if downloading is enabled)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("must have at least one stock symbol")
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
			resp, err := resty.R().Put(fmt.Sprintf("%v/stocks/enable/intraday/%v/tick", address, s))
			if err != nil {
				return fmt.Errorf("There was an error while enabling tick intraday histories for stock %v, error: %v\n", s, err)
			}

			if resp.StatusCode() != http.StatusOK {
				return fmt.Errorf("There was an error while enabling tick intraday histories for stock %v, statusCode: %v  response: %v\n", s, resp.StatusCode(), resp)
			}

			var respErrors rest.ErrorState
			err = json.Unmarshal(resp.Body(), &respErrors)
			if err != nil {
				return fmt.Errorf("There was a problem parsing the server response while enabling tick intraday histories for stock %v : %v had error:%v\n", s, string(resp.Body()), err)
			}

			if respErrors.ErrorReturned {
				return fmt.Errorf("There was a problem server side while enabling tick intraday histories for stock %v : %v\n", s, respErrors.ErrorReturned)
			}

			fmt.Printf("Successfullly enabled tick intraday histories for stock: %v\n", s)
		}
		return nil
	},
}

var enableIntraday1MinCmd = &cobra.Command{
	Use:   "1min STOCK [STOCK] ...",
	Short: "Enables 1 minute intraday for a particular stock(s) (if downloading is enabled)",
	Long:  `Enables gathering 1 minute resolution intraday for a particular stock(s) (if downloading is enabled)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("must have at least one stock symbol")
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
			resp, err := resty.R().Put(fmt.Sprintf("%v/stocks/enable/intraday/%v/1min", address, s))
			if err != nil {
				return fmt.Errorf("There was an error while enabling 1 min intraday histories for stock %v, error: %v\n", s, err)
			}

			if resp.StatusCode() != http.StatusOK {
				return fmt.Errorf("There was an error while enabling 1 min intraday histories for stock %v, statusCode: %v  response: %v\n", s, resp.StatusCode(), resp)
			}

			var respErrors rest.ErrorState
			err = json.Unmarshal(resp.Body(), &respErrors)
			if err != nil {
				return fmt.Errorf("There was a problem parsing the server response while enabling 1 min intraday histories for stock %v : %v had error:%v\n", s, string(resp.Body()), err)
			}

			if respErrors.ErrorReturned {
				return fmt.Errorf("There was a problem server side while enabling 1 min intraday histories for stock %v : %v\n", s, respErrors.ErrorReturned)
			}

			fmt.Printf("Successfullly enabled 1 min intraday histories for stock: %v\n", s)
		}
		return nil
	},
}

var enableIntraday5MinCmd = &cobra.Command{
	Use:   "5min STOCK [STOCK] ...",
	Short: "Enables 5 minute intraday for a particular stock(s) (if downloading is enabled)",
	Long:  `Enables gathering 5 minute resolution intraday for a particular stock(s) (if downloading is enabled)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("must have at least one stock symbol")
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
			resp, err := resty.R().Put(fmt.Sprintf("%v/stocks/enable/intraday/%v/tick", address, s))
			if err != nil {
				return fmt.Errorf("There was an error while enabling 5 min intraday histories for stock %v, error: %v\n", s, err)
			}

			if resp.StatusCode() != http.StatusOK {
				return fmt.Errorf("There was an error while enabling 5 min intraday histories for stock %v, statusCode: %v  response: %v\n", s, resp.StatusCode(), resp)
			}

			var respErrors rest.ErrorState
			err = json.Unmarshal(resp.Body(), &respErrors)
			if err != nil {
				return fmt.Errorf("There was a problem parsing the server response while enabling 5 min intraday histories for stock %v : %v had error:%v\n", s, string(resp.Body()), err)
			}

			if respErrors.ErrorReturned {
				return fmt.Errorf("There was a problem server side while enabling 5 min intraday histories for stock %v : %v\n", s, respErrors.ErrorReturned)
			}

			fmt.Printf("Successfullly enabled 5 min intraday histories for stock: %v\n", s)
		}
		return nil
	},
}

var enableAllIntradayCmd = &cobra.Command{
	Use:   "intraday",
	Short: "Intraday history gathering (if downloading is enabled)",
	Long:  `Intraday history gathering (if downloading is enabled)`,
}

var enableAllIntradayAllCmd = &cobra.Command{
	Use:   "all",
	Short: "Enables all intraday gathering for all stocks (if downloading is enabled)",
	Long:  `Enables gathering ticks, 1 min and 5 min resolution intraday for all stocks (if downloading is enabled)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		address, err := cmd.Flags().GetString("serverAddress")
		if err != nil {
			return fmt.Errorf("Could not get server address from passed in arguments: %v\n", err)
		}

		resp, err := resty.R().Put(fmt.Sprintf("%v/stocks/enable/all/intraday/all", address))
		if err != nil {
			return fmt.Errorf("There was an error while enabling all intraday histories for all stocks, error: %v\n", err)
		}

		if resp.StatusCode() != http.StatusOK {
			return fmt.Errorf("There was an error while enabling all intraday histories for all stocks, statusCode: %v  response: %v\n", resp.StatusCode(), resp)
		}

		var respErrors rest.ErrorState
		err = json.Unmarshal(resp.Body(), &respErrors)
		if err != nil {
			return fmt.Errorf("There was a problem parsing the server response while enabling all intraday histories for all stocks: %v had error:%v\n", string(resp.Body()), err)
		}

		if respErrors.ErrorReturned {
			return fmt.Errorf("There was a problem server side while enabling all intraday histories for all stocks: %v\n", respErrors.ErrorReturned)
		}
		fmt.Printf("Successfullly enabled all intraday histories for all stocks\n")

		return nil
	},
}

var enableAllIntradayTickCmd = &cobra.Command{
	Use:   "tick",
	Short: "Enables tick intraday gathering for all stocks (if downloading is enabled)",
	Long:  `Enables gathering ticks resolution intraday for all stocks (if downloading is enabled)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		address, err := cmd.Flags().GetString("serverAddress")
		if err != nil {
			return fmt.Errorf("Could not get server address from passed in arguments: %v\n", err)
		}

		resp, err := resty.R().Put(fmt.Sprintf("%v/stocks/enable/all/intraday/tick", address))
		if err != nil {
			return fmt.Errorf("There was an error while enabling tick intraday histories for all stocks, error: %v\n", err)
		}

		if resp.StatusCode() != http.StatusOK {
			return fmt.Errorf("There was an error while enabling tick intraday histories for all stocks, statusCode: %v  response: %v\n", resp.StatusCode(), resp)
		}

		var respErrors rest.ErrorState
		err = json.Unmarshal(resp.Body(), &respErrors)
		if err != nil {
			return fmt.Errorf("There was a problem parsing the server response while enabling tick intraday histories for all stocks: %v had error:%v\n", string(resp.Body()), err)
		}

		if respErrors.ErrorReturned {
			return fmt.Errorf("There was a problem server side while enabling tick intraday histories for all stocks: %v\n", respErrors.ErrorReturned)
		}
		fmt.Printf("Successfullly enabled tick intraday histories for all stocks\n")

		return nil
	},
}

var enableAllIntraday1MinCmd = &cobra.Command{
	Use:   "1min",
	Short: "Enables 1 minute intraday gathering for all stocks (if downloading is enabled)",
	Long:  `Enables gathering 1 minute resolution Intraday for all stocks (if downloading is enabled)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		address, err := cmd.Flags().GetString("serverAddress")
		if err != nil {
			return fmt.Errorf("Could not get server address from passed in arguments: %v\n", err)
		}

		resp, err := resty.R().Put(fmt.Sprintf("%v/stocks/enable/all/intraday/1min", address))
		if err != nil {
			return fmt.Errorf("There was an error while enabling 1 min intraday histories for all stocks, error: %v\n", err)
		}

		if resp.StatusCode() != http.StatusOK {
			return fmt.Errorf("There was an error while enabling 1 min intraday histories for all stocks, statusCode: %v  response: %v\n", resp.StatusCode(), resp)
		}

		var respErrors rest.ErrorState
		err = json.Unmarshal(resp.Body(), &respErrors)
		if err != nil {
			return fmt.Errorf("There was a problem parsing the server response while enabling 1 min intraday histories for all stocks: %v had error:%v\n", string(resp.Body()), err)
		}

		if respErrors.ErrorReturned {
			return fmt.Errorf("There was a problem server side while enabling 1 min intraday histories for all stocks: %v\n", respErrors.ErrorReturned)
		}
		fmt.Printf("Successfullly enabled 1 min intraday histories for all stocks\n")

		return nil
	},
}

var enableAllIntraday5MinCmd = &cobra.Command{
	Use:   "5min",
	Short: "Enables 5 minute intraday gathering for all stocks (if downloading is enabled)",
	Long:  `Enables gathering 5 minute resolution intraday for all stocks (if downloading is enabled)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		address, err := cmd.Flags().GetString("serverAddress")
		if err != nil {
			return fmt.Errorf("Could not get server address from passed in arguments: %v\n", err)
		}

		resp, err := resty.R().Put(fmt.Sprintf("%v/stocks/enable/all/intraday/5min", address))
		if err != nil {
			return fmt.Errorf("There was an error while enabling 5 min intraday histories for all stocks, error: %v\n", err)
		}

		if resp.StatusCode() != http.StatusOK {
			return fmt.Errorf("There was an error while enabling 5 min intraday histories for all stocks, statusCode: %v  response: %v\n", resp.StatusCode(), resp)
		}

		var respErrors rest.ErrorState
		err = json.Unmarshal(resp.Body(), &respErrors)
		if err != nil {
			return fmt.Errorf("There was a problem parsing the server response while enabling 5 min intraday histories for all stocks: %v had error:%v\n", string(resp.Body()), err)
		}

		if respErrors.ErrorReturned {
			return fmt.Errorf("There was a problem server side while enabling 5 min intraday histories for all stocks: %v\n", respErrors.ErrorReturned)
		}
		fmt.Printf("Successfullly enabled 5 min intraday histories for all stocks\n")

		return nil
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
		}

		if resp.StatusCode() != http.StatusOK {
			return fmt.Errorf("There was an error while enabling %v for stock %v, statusCode: %v  response: %v\n", commandName, s, resp.StatusCode(), resp)
		}

		var respErrors rest.ErrorState
		err = json.Unmarshal(resp.Body(), &respErrors)
		if err != nil {
			return fmt.Errorf("There was a problem parsing the server response while enabling %v for stock %v : %v had error:%v\n", commandName, s, string(resp.Body()), err)
		}

		if respErrors.ErrorReturned {
			return fmt.Errorf("There was a problem server side while enabling %v for stock %v : %v\n", commandName, s, respErrors.ErrorReturned)
		}

		fmt.Printf("Successfullly enabled %v for stock: %v\n", commandName, s)
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
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("There was an error while enabling %v for all stocks, statusCode: %v  response: %v\n", commandName, resp.StatusCode(), resp)
	}

	var respErrors rest.ErrorState
	err = json.Unmarshal(resp.Body(), &respErrors)
	if err != nil {
		return fmt.Errorf("There was a problem parsing the server response while enabling %v for all stocks: %v had error:%v\n", commandName, string(resp.Body()), err)
	}

	if respErrors.ErrorReturned {
		return fmt.Errorf("There was a problem server side while enabling %v for all stocks: %v\n", commandName, respErrors.ErrorReturned)
	}
	fmt.Printf("Successfullly enabled %v for all stocks\n", commandName)

	return nil
}

func init() {
	rootCmd.AddCommand(enableCmd)
	enableCmd.AddCommand(enableDownloadingCmd)
	enableCmd.AddCommand(enableHistoryCmd)
	enableCmd.AddCommand(enableIntradayCmd)
	enableCmd.AddCommand(enableQuotesCmd)
	enableCmd.AddCommand(enableStableQuotesCmd)
	enableCmd.AddCommand(enableAllCmd)

	enableIntradayCmd.AddCommand(enableIntradayAllCmd)
	enableIntradayCmd.AddCommand(enableIntradayTickCmd)
	enableIntradayCmd.AddCommand(enableIntraday1MinCmd)
	enableIntradayCmd.AddCommand(enableIntraday5MinCmd)

	enableHistoryCmd.AddCommand(enableHistoryAllCmd)
	enableHistoryCmd.AddCommand(enableHistoryDailyCmd)
	enableHistoryCmd.AddCommand(enableHistoryWeeklyCmd)
	enableHistoryCmd.AddCommand(enableHistoryMonthlyCmd)
	enableHistoryCmd.AddCommand(enableHistoryYearlyCmd)

	enableAllCmd.AddCommand(enableAllAllCmd)
	enableAllCmd.AddCommand(enableAllDownloadingCmd)
	enableAllCmd.AddCommand(enableAllHistoryCmd)
	enableAllCmd.AddCommand(enableAllIntradayCmd)
	enableAllCmd.AddCommand(enableAllQuotesCmd)
	enableAllCmd.AddCommand(enableAllStableQuotesCmd)

	enableAllIntradayCmd.AddCommand(enableAllIntradayAllCmd)
	enableAllIntradayCmd.AddCommand(enableAllIntradayTickCmd)
	enableAllIntradayCmd.AddCommand(enableAllIntraday1MinCmd)
	enableAllIntradayCmd.AddCommand(enableAllIntraday5MinCmd)

	enableAllHistoryCmd.AddCommand(enableAllHistoryAllCmd)
	enableAllHistoryCmd.AddCommand(enableAllHistoryDailyCmd)
	enableAllHistoryCmd.AddCommand(enableAllHistoryWeeklyCmd)
	enableAllHistoryCmd.AddCommand(enableAllHistoryMonthlyCmd)
	enableAllHistoryCmd.AddCommand(enableAllHistoryYearlyCmd)

}
