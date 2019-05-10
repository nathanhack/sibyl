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

var disableAllCmd = &cobra.Command{
	Use:   "all",
	Short: "Disables all attributes",
	Long:  `Disables all attributes`,
}

var disableAllAllCmd = &cobra.Command{
	Use:   "all",
	Short: "Disables all things for all stocks",
	Long:  `Disables all things for all stocks`,
	RunE: func(cmd *cobra.Command, args []string) error {
		address, err := cmd.Flags().GetString("serverAddress")
		if err != nil {
			return fmt.Errorf("Could not get server address from passed in arguments: %v\n", err)
		}

		resp, err := resty.R().Put(fmt.Sprintf("%v/stocks/disable/all", address))
		if err != nil {
			return fmt.Errorf("There was an error while disabling all things for all stocks, error: %v\n", err)
		}

		if resp.StatusCode() != http.StatusOK {
			return fmt.Errorf("There was an error while disabling all things for all stocks, statusCode: %v response: %v\n", resp.StatusCode(), resp)
		}

		var respErrors rest.ErrorState
		err = json.Unmarshal(resp.Body(), &respErrors)
		if err != nil {
			return fmt.Errorf("There was a problem parsing the server response while disabling all things for all stocks: %v  had error:%v\n", string(resp.Body()), err)
		}

		if respErrors.ErrorReturned {
			return fmt.Errorf("There was a problem server side while disabling all things for all stocks: %v\n", respErrors.ErrorReturned)
		}
		fmt.Printf("Successfullly disabled all things for all stocks\n")

		return nil
	},
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
	Short: "Disables daily, weekly, monthly, and yearly history for a particular stock",
	Long:  `Disables gathering daily, weekly, monthly, and yearly history for a particular stock`,
}

var disableHistoryAllCmd = &cobra.Command{
	Use:   "all STOCK [STOCK] ...",
	Short: "Disables daily, weekly, monthly, and yearly histories for a particular stock",
	Long:  `Disables gathering daily, weekly, monthly, and yearly histories for a particular stock`,
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
			resp, err := resty.R().Put(fmt.Sprintf("%v/stocks/disable/history/%v/all", address, s))
			if err != nil {
				return fmt.Errorf("There was an error while disabling all histories for stock %v, error: %v\n", s, err)
			}

			if resp.StatusCode() != http.StatusOK {
				return fmt.Errorf("There was an error while disabling all histories for stock %v, statusCode: %v response: %v\n", s, resp.StatusCode(), resp)
			}

			var respErrors rest.ErrorState
			err = json.Unmarshal(resp.Body(), &respErrors)
			if err != nil {
				return fmt.Errorf("There was a problem parsing the server response while disabling all histories for stock %v : %v  had error:%v\n", s, string(resp.Body()), err)
			}

			if respErrors.ErrorReturned {
				return fmt.Errorf("There was a problem server side while disabling all histories for stock %v : %v\n", s, respErrors.ErrorReturned)
			}

			fmt.Printf("Successfullly disabled all histories for stock: %v\n", s)
		}
		return nil
	},
}

var disableHistoryDailyCmd = &cobra.Command{
	Use:   "daily STOCK [STOCK] ...",
	Short: "Disables daily history for a particular stock",
	Long:  `Disables gathering daily history for a particular stock`,
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
			resp, err := resty.R().Put(fmt.Sprintf("%v/stocks/disable/history/%v/daily", address, s))
			if err != nil {
				return fmt.Errorf("There was an error while disabling daily histories for stock %v, error: %v\n", s, err)
			}

			if resp.StatusCode() != http.StatusOK {
				return fmt.Errorf("There was an error while disabling daily histories for stock %v, statusCode: %v response: %v\n", s, resp.StatusCode(), resp)
			}

			var respErrors rest.ErrorState
			err = json.Unmarshal(resp.Body(), &respErrors)
			if err != nil {
				return fmt.Errorf("There was a problem parsing the server response while disabling daily histories for stock %v : %v  had error:%v\n", s, string(resp.Body()), err)
			}

			if respErrors.ErrorReturned {
				return fmt.Errorf("There was a problem server side while disabling daily histories for stock %v : %v\n", s, respErrors.ErrorReturned)
			}

			fmt.Printf("Successfullly disabled daily histories for stock: %v\n", s)
		}
		return nil
	},
}

var disableHistoryWeeklyCmd = &cobra.Command{
	Use:   "weekly STOCK [STOCK] ...",
	Short: "Disables weekly history for a particular stock",
	Long:  `Disables gathering weekly history for a particular stock`,
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
			resp, err := resty.R().Put(fmt.Sprintf("%v/stocks/disable/history/%v/weekly", address, s))
			if err != nil {
				return fmt.Errorf("There was an error while disabling weekly histories for stock %v, error: %v\n", s, err)
			}

			if resp.StatusCode() != http.StatusOK {
				return fmt.Errorf("There was an error while disabling weekly histories for stock %v, statusCode: %v response: %v\n", s, resp.StatusCode(), resp)
			}

			var respErrors rest.ErrorState
			err = json.Unmarshal(resp.Body(), &respErrors)
			if err != nil {
				return fmt.Errorf("There was a problem parsing the server response while disabling weekly histories for stock %v : %v  had error:%v\n", s, string(resp.Body()), err)
			}

			if respErrors.ErrorReturned {
				return fmt.Errorf("There was a problem server side while disabling weekly histories for stock %v : %v\n", s, respErrors.ErrorReturned)
			}

			fmt.Printf("Successfullly disabled weekly histories for stock: %v\n", s)
		}
		return nil
	},
}

var disableHistoryMonthlyCmd = &cobra.Command{
	Use:   "monthly STOCK [STOCK] ...",
	Short: "Disables monthly history for a particular stock",
	Long:  `Disables gathering monthly History for a particular stock`,
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
			resp, err := resty.R().Put(fmt.Sprintf("%v/stocks/disable/history/%v/monthly", address, s))
			if err != nil {
				return fmt.Errorf("There was an error while disabling monthly histories for stock %v, error: %v\n", s, err)
			}

			if resp.StatusCode() != http.StatusOK {
				return fmt.Errorf("There was an error while disabling monthly histories for stock %v, statusCode: %v response: %v\n", s, resp.StatusCode(), resp)
			}

			var respErrors rest.ErrorState
			err = json.Unmarshal(resp.Body(), &respErrors)
			if err != nil {
				return fmt.Errorf("There was a problem parsing the server response while disabling monthly histories for stock %v : %v  had error:%v\n", s, string(resp.Body()), err)
			}

			if respErrors.ErrorReturned {
				return fmt.Errorf("There was a problem server side while disabling monthly histories for stock %v : %v\n", s, respErrors.ErrorReturned)
			}

			fmt.Printf("Successfullly disabled monthly histories for stock: %v\n", s)
		}
		return nil
	},
}

var disableHistoryYearlyCmd = &cobra.Command{
	Use:   "yearly STOCK [STOCK] ...",
	Short: "Disables yearly history for a particular stock",
	Long:  `Disables gathering yearly history for a particular stock`,
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
			resp, err := resty.R().Put(fmt.Sprintf("%v/stocks/disable/history/%v/yearly", address, s))
			if err != nil {
				return fmt.Errorf("There was an error while disabling yearly histories for stock %v, error: %v\n", s, err)
			}

			if resp.StatusCode() != http.StatusOK {
				return fmt.Errorf("There was an error while disabling yearly histories for stock %v, statusCode: %v response: %v\n", s, resp.StatusCode(), resp)
			}

			var respErrors rest.ErrorState
			err = json.Unmarshal(resp.Body(), &respErrors)
			if err != nil {
				return fmt.Errorf("There was a problem parsing the server response while disabling yearly histories for stock %v : %v  had error:%v\n", s, string(resp.Body()), err)
			}

			if respErrors.ErrorReturned {
				return fmt.Errorf("There was a problem server side while disabling yearly histories for stock %v : %v\n", s, respErrors.ErrorReturned)
			}

			fmt.Printf("Successfullly disabled yearly histories for stock: %v\n", s)
		}
		return nil
	},
}

var disableAllHistoryCmd = &cobra.Command{
	Use:   "history",
	Short: "Disables History for all stocks",
	Long:  `Disables gathering daily History for all stocks`,
}

var disableAllHistoryAllCmd = &cobra.Command{
	Use:   "all",
	Short: "Disables daily, weekly, monthly, and yearly histories for a all stocks",
	Long:  `Disables gathering daily, weekly, monthly, and yearly histories for all stocks`,
	RunE: func(cmd *cobra.Command, args []string) error {
		address, err := cmd.Flags().GetString("serverAddress")
		if err != nil {
			return fmt.Errorf("Could not get server address from passed in arguments: %v\n", err)
		}

		resp, err := resty.R().Put(fmt.Sprintf("%v/stocks/disable/all/history/all", address))
		if err != nil {
			return fmt.Errorf("There was an error while disabling all histories for all stocks, error: %v\n", err)
		}

		if resp.StatusCode() != http.StatusOK {
			return fmt.Errorf("There was an error while disabling all histories for all stocks, statusCode: %v response: %v\n", resp.StatusCode(), resp)
		}

		var respErrors rest.ErrorState
		err = json.Unmarshal(resp.Body(), &respErrors)
		if err != nil {
			return fmt.Errorf("There was a problem parsing the server response while disabling all histories for all stocks: %v  had error:%v\n", string(resp.Body()), err)
		}

		if respErrors.ErrorReturned {
			return fmt.Errorf("There was a problem server side while disabling all histories for all stocks: %v\n", respErrors.ErrorReturned)
		}
		fmt.Printf("Successfullly disabled all histories for all stocks\n")

		return nil
	},
}

var disableAllHistoryDailyCmd = &cobra.Command{
	Use:   "daily",
	Short: "Disables daily history for a all stocks",
	Long:  `Disables gathering daily history for all stocks`,
	RunE: func(cmd *cobra.Command, args []string) error {
		address, err := cmd.Flags().GetString("serverAddress")
		if err != nil {
			return fmt.Errorf("Could not get server address from passed in arguments: %v\n", err)
		}

		resp, err := resty.R().Put(fmt.Sprintf("%v/stocks/disable/all/history/daily", address))
		if err != nil {
			return fmt.Errorf("There was an error while disabling daily histories for all stocks, error: %v\n", err)
		}

		if resp.StatusCode() != http.StatusOK {
			return fmt.Errorf("There was an error while disabling daily histories for all stocks, statusCode: %v response: %v\n", resp.StatusCode(), resp)
		}

		var respErrors rest.ErrorState
		err = json.Unmarshal(resp.Body(), &respErrors)
		if err != nil {
			return fmt.Errorf("There was a problem parsing the server response while disabling daily histories for all stocks: %v  had error:%v\n", string(resp.Body()), err)
		}

		if respErrors.ErrorReturned {
			return fmt.Errorf("There was a problem server side while disabling daily histories for all stocks: %v\n", respErrors.ErrorReturned)
		}
		fmt.Printf("Successfullly disabled daily histories for all stocks\n")

		return nil
	},
}

var disableAllHistoryWeeklyCmd = &cobra.Command{
	Use:   "weekly",
	Short: "Disables weekly history for a all stocks",
	Long:  `Disables gathering weekly history for all stocks`,
	RunE: func(cmd *cobra.Command, args []string) error {
		address, err := cmd.Flags().GetString("serverAddress")
		if err != nil {
			return fmt.Errorf("Could not get server address from passed in arguments: %v\n", err)
		}

		resp, err := resty.R().Put(fmt.Sprintf("%v/stocks/disable/all/history/weekly", address))
		if err != nil {
			return fmt.Errorf("There was an error while disabling weekly histories for all stocks, error: %v\n", err)
		}

		if resp.StatusCode() != http.StatusOK {
			return fmt.Errorf("There was an error while disabling weekly histories for all stocks, statusCode: %v response: %v\n", resp.StatusCode(), resp)
		}

		var respErrors rest.ErrorState
		err = json.Unmarshal(resp.Body(), &respErrors)
		if err != nil {
			return fmt.Errorf("There was a problem parsing the server response while disabling weekly histories for all stocks: %v had error:%v\n", string(resp.Body()), err)
		}

		if respErrors.ErrorReturned {
			return fmt.Errorf("There was a problem server side while disabling weekly histories for all stocks: %v\n", respErrors.ErrorReturned)
		}
		fmt.Printf("Successfullly disabled weekly histories for all stocks\n")

		return nil
	},
}

var disableAllHistoryMonthlyCmd = &cobra.Command{
	Use:   "monthly",
	Short: "Disables monthly history for a all stocks",
	Long:  `Disables gathering monthly history for all stocks`,
	RunE: func(cmd *cobra.Command, args []string) error {
		address, err := cmd.Flags().GetString("serverAddress")
		if err != nil {
			return fmt.Errorf("Could not get server address from passed in arguments: %v\n", err)
		}

		resp, err := resty.R().Put(fmt.Sprintf("%v/stocks/disable/all/history/monthly", address))
		if err != nil {
			return fmt.Errorf("There was an error while disabling monthly histories for all stocks, error: %v\n", err)
		}

		if resp.StatusCode() != http.StatusOK {
			return fmt.Errorf("There was an error while disabling monthly histories for all stocks, statusCode: %v response: %v\n", resp.StatusCode(), resp)
		}

		var respErrors rest.ErrorState
		err = json.Unmarshal(resp.Body(), &respErrors)
		if err != nil {
			return fmt.Errorf("There was a problem parsing the server response while disabling monthly histories for all stocks: %v had error:%v\n", string(resp.Body()), err)
		}

		if respErrors.ErrorReturned {
			return fmt.Errorf("There was a problem server side while disabling monthly histories for all stocks: %v\n", respErrors.ErrorReturned)
		}
		fmt.Printf("Successfullly disabled monthly histories for all stocks\n")

		return nil
	},
}

var disableAllHistoryYearlyCmd = &cobra.Command{
	Use:   "yearly",
	Short: "Disables yearly history for a all stocks",
	Long:  `Disables gathering yearly history for all stocks`,
	RunE: func(cmd *cobra.Command, args []string) error {
		address, err := cmd.Flags().GetString("serverAddress")
		if err != nil {
			return fmt.Errorf("Could not get server address from passed in arguments: %v\n", err)
		}

		resp, err := resty.R().Put(fmt.Sprintf("%v/stocks/disable/all/history/yearly", address))
		if err != nil {
			return fmt.Errorf("There was an error while disabling yearly histories for all stocks, error: %v\n", err)
		}

		if resp.StatusCode() != http.StatusOK {
			return fmt.Errorf("There was an error while disabling yearly histories for all stocks, statusCode: %v response: %v\n", resp.StatusCode(), resp)
		}

		var respErrors rest.ErrorState
		err = json.Unmarshal(resp.Body(), &respErrors)
		if err != nil {
			return fmt.Errorf("There was a problem parsing the server response while disabling yearly histories for all stocks: %v had error:%v\n", string(resp.Body()), err)
		}

		if respErrors.ErrorReturned {
			return fmt.Errorf("There was a problem server side while disabling yearly histories for all stocks: %v\n", respErrors.ErrorReturned)
		}
		fmt.Printf("Successfullly disabled yearly histories for all stocks\n")

		return nil
	},
}

var disableIntradayCmd = &cobra.Command{
	Use:   "intraday",
	Short: "Intraday history gathering",
	Long:  `Intraday history gathering`,
}

var disableIntradayAllCmd = &cobra.Command{
	Use:   "all STOCK [STOCK] ...",
	Short: "Disables all intraday histories for a particular stock(s)",
	Long:  `Disables gathering Ticks, 1 min, 5 min resolution intraday histories for a particular stock(s)`,
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
			resp, err := resty.R().Put(fmt.Sprintf("%v/stocks/disable/intraday/%v/all", address, s))
			if err != nil {
				return fmt.Errorf("There was an error while disabling all intraday histories for stock %v, error: %v\n", s, err)
			}

			if resp.StatusCode() != http.StatusOK {
				return fmt.Errorf("There was an error while disabling all intraday histories for stock %v, statusCode: %v response: %v\n", s, resp.StatusCode(), resp)
			}

			var respErrors rest.ErrorState
			err = json.Unmarshal(resp.Body(), &respErrors)
			if err != nil {
				return fmt.Errorf("There was a problem parsing the server response while disabling all intraday histories for stock %v : %v had error:%v\n", s, string(resp.Body()), err)
			}

			if respErrors.ErrorReturned {
				return fmt.Errorf("There was a problem server side while disabling all intraday histories for stock %v : %v\n", s, respErrors.ErrorReturned)
			}

			fmt.Printf("Successfullly disabled all intraday histories for stock: %v\n", s)
		}
		return nil
	},
}

var disableIntradayTickCmd = &cobra.Command{
	Use:   "tick STOCK [STOCK] ...",
	Short: "Disables tick intraday for a particular stock(s)",
	Long:  `Disables gathering tick resolution intraday for a particular stock(s)`,
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
			resp, err := resty.R().Put(fmt.Sprintf("%v/stocks/disable/intraday/%v/tick", address, s))
			if err != nil {
				return fmt.Errorf("There was an error while disabling tick intraday histories for stock %v, error: %v\n", s, err)
			}

			if resp.StatusCode() != http.StatusOK {
				return fmt.Errorf("There was an error while disabling tick intraday histories for stock %v, statusCode: %v response: %v\n", s, resp.StatusCode(), resp)
			}

			var respErrors rest.ErrorState
			err = json.Unmarshal(resp.Body(), &respErrors)
			if err != nil {
				return fmt.Errorf("There was a problem parsing the server response while disabling tick intraday histories for stock %v : %v had error:%v\n", s, string(resp.Body()), err)
			}

			if respErrors.ErrorReturned {
				return fmt.Errorf("There was a problem server side while disabling tick intraday histories for stock %v : %v\n", s, respErrors.ErrorReturned)
			}

			fmt.Printf("Successfullly disabled tick intraday histories for stock: %v\n", s)
		}
		return nil
	},
}

var disableIntraday1MinCmd = &cobra.Command{
	Use:   "1min STOCK [STOCK] ...",
	Short: "Disables 1 minute intraday for a particular stock(s)",
	Long:  `Disables gathering 1 minute resolution intraday for a particular stock(s)`,
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
			resp, err := resty.R().Put(fmt.Sprintf("%v/stocks/disable/intraday/%v/1min", address, s))
			if err != nil {
				return fmt.Errorf("There was an error while disabling 1 min intraday histories for stock %v, error: %v\n", s, err)
			}

			if resp.StatusCode() != http.StatusOK {
				return fmt.Errorf("There was an error while disabling 1 min intraday histories for stock %v, statusCode: %v response: %v\n", s, resp.StatusCode(), resp)
			}

			var respErrors rest.ErrorState
			err = json.Unmarshal(resp.Body(), &respErrors)
			if err != nil {
				return fmt.Errorf("There was a problem parsing the server response while disabling 1 min intraday histories for stock %v : %v had error:%v\n", s, string(resp.Body()), err)
			}

			if respErrors.ErrorReturned {
				return fmt.Errorf("There was a problem server side while disabling 1 min intraday histories for stock %v : %v\n", s, respErrors.ErrorReturned)
			}

			fmt.Printf("Successfullly disabled 1 min intraday histories for stock: %v\n", s)
		}
		return nil
	},
}

var disableIntraday5MinCmd = &cobra.Command{
	Use:   "5min STOCK [STOCK] ...",
	Short: "Disables 5 minute intraday for a particular stock(s)",
	Long:  `Disables gathering 5 minute resolution intraday for a particular stock(s)`,
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
			resp, err := resty.R().Put(fmt.Sprintf("%v/stocks/disable/intraday/%v/tick", address, s))
			if err != nil {
				return fmt.Errorf("There was an error while disabling 5 min intraday histories for stock %v, error: %v\n", s, err)
			}

			if resp.StatusCode() != http.StatusOK {
				return fmt.Errorf("There was an error while disabling 5 min intraday histories for stock %v, statusCode: %v response: %v\n", s, resp.StatusCode(), resp)
			}

			var respErrors rest.ErrorState
			err = json.Unmarshal(resp.Body(), &respErrors)
			if err != nil {
				return fmt.Errorf("There was a problem parsing the server response while disabling 5 min intraday histories for stock %v : %v had error:%v\n", s, string(resp.Body()), err)
			}

			if respErrors.ErrorReturned {
				return fmt.Errorf("There was a problem server side while disabling 5 min intraday histories for stock %v : %v\n", s, respErrors.ErrorReturned)
			}

			fmt.Printf("Successfullly disabled 5 min intraday histories for stock: %v\n", s)
		}
		return nil
	},
}

var disableAllIntradayCmd = &cobra.Command{
	Use:   "intraday",
	Short: "Intraday history gathering",
	Long:  `Intraday history gathering`,
}

var disableAllIntradayAllCmd = &cobra.Command{
	Use:   "all",
	Short: "Disables all intraday gathering for all stocks",
	Long:  `Disables gathering ticks, 1 min and 5 min resolution intraday for all stocks`,
	RunE: func(cmd *cobra.Command, args []string) error {
		address, err := cmd.Flags().GetString("serverAddress")
		if err != nil {
			return fmt.Errorf("Could not get server address from passed in arguments: %v\n", err)
		}

		resp, err := resty.R().Put(fmt.Sprintf("%v/stocks/disable/all/intraday/all", address))
		if err != nil {
			return fmt.Errorf("There was an error while disabling all intraday histories for all stocks, error: %v\n", err)
		}

		if resp.StatusCode() != http.StatusOK {
			return fmt.Errorf("There was an error while disabling all intraday histories for all stocks, statusCode: %v response: %v\n", resp.StatusCode(), resp)
		}

		var respErrors rest.ErrorState
		err = json.Unmarshal(resp.Body(), &respErrors)
		if err != nil {
			return fmt.Errorf("There was a problem parsing the server response while disabling all intraday histories for all stocks: %v had error:%v\n", string(resp.Body()), err)
		}

		if respErrors.ErrorReturned {
			return fmt.Errorf("There was a problem server side while disabling all intraday histories for all stocks: %v\n", respErrors.ErrorReturned)
		}
		fmt.Printf("Successfullly disabled all intraday histories for all stocks\n")

		return nil
	},
}

var disableAllIntradayTickCmd = &cobra.Command{
	Use:   "tick",
	Short: "Disables tick intraday gathering for all stocks",
	Long:  `Disables gathering ticks resolution intraday for all stocks`,
	RunE: func(cmd *cobra.Command, args []string) error {
		address, err := cmd.Flags().GetString("serverAddress")
		if err != nil {
			return fmt.Errorf("Could not get server address from passed in arguments: %v\n", err)
		}

		resp, err := resty.R().Put(fmt.Sprintf("%v/stocks/disable/all/intraday/tick", address))
		if err != nil {
			return fmt.Errorf("There was an error while disabling tick intraday histories for all stocks, error: %v\n", err)
		}

		if resp.StatusCode() != http.StatusOK {
			return fmt.Errorf("There was an error while disabling tick intraday histories for all stocks, statusCode: %v response: %v\n", resp.StatusCode(), resp)
		}

		var respErrors rest.ErrorState
		err = json.Unmarshal(resp.Body(), &respErrors)
		if err != nil {
			return fmt.Errorf("There was a problem parsing the server response while disabling tick intraday histories for all stocks: %v had error:%v\n", string(resp.Body()), err)
		}

		if respErrors.ErrorReturned {
			return fmt.Errorf("There was a problem server side while disabling tick intraday histories for all stocks: %v\n", respErrors.ErrorReturned)
		}
		fmt.Printf("Successfullly disabled tick intraday histories for all stocks\n")

		return nil
	},
}

var disableAllIntraday1MinCmd = &cobra.Command{
	Use:   "1min",
	Short: "Disables 1 minute intraday gathering for all stocks",
	Long:  `Disables gathering 1 minute resolution Intraday for all stocks`,
	RunE: func(cmd *cobra.Command, args []string) error {
		address, err := cmd.Flags().GetString("serverAddress")
		if err != nil {
			return fmt.Errorf("Could not get server address from passed in arguments: %v\n", err)
		}

		resp, err := resty.R().Put(fmt.Sprintf("%v/stocks/disable/all/intraday/1min", address))
		if err != nil {
			return fmt.Errorf("There was an error while disabling 1 min intraday histories for all stocks, error: %v\n", err)
		}

		if resp.StatusCode() != http.StatusOK {
			return fmt.Errorf("There was an error while disabling 1 min intraday histories for all stocks, statusCode: %v response: %v\n", resp.StatusCode(), resp)
		}

		var respErrors rest.ErrorState
		err = json.Unmarshal(resp.Body(), &respErrors)
		if err != nil {
			return fmt.Errorf("There was a problem parsing the server response while disabling 1 min intraday histories for all stocks: %v had error:%v\n", string(resp.Body()), err)
		}

		if respErrors.ErrorReturned {
			return fmt.Errorf("There was a problem server side while disabling 1 min intraday histories for all stocks: %v\n", respErrors.ErrorReturned)
		}
		fmt.Printf("Successfullly disabled 1 min intraday histories for all stocks\n")

		return nil
	},
}

var disableAllIntraday5MinCmd = &cobra.Command{
	Use:   "5min",
	Short: "Disables 5 minute intraday gathering for all stocks",
	Long:  `Disables gathering 5 minute resolution intraday for all stocks`,
	RunE: func(cmd *cobra.Command, args []string) error {
		address, err := cmd.Flags().GetString("serverAddress")
		if err != nil {
			return fmt.Errorf("Could not get server address from passed in arguments: %v\n", err)
		}

		resp, err := resty.R().Put(fmt.Sprintf("%v/stocks/disable/all/intraday/5min", address))
		if err != nil {
			return fmt.Errorf("There was an error while disabling 5 min intraday histories for all stocks, error: %v\n", err)
		}

		if resp.StatusCode() != http.StatusOK {
			return fmt.Errorf("There was an error while disabling 5 min intraday histories for all stocks, statusCode: %v response: %v\n", resp.StatusCode(), resp)
		}

		var respErrors rest.ErrorState
		err = json.Unmarshal(resp.Body(), &respErrors)
		if err != nil {
			return fmt.Errorf("There was a problem parsing the server response while disabling 5 min intraday histories for all stocks: %v had error:%v\n", string(resp.Body()), err)
		}

		if respErrors.ErrorReturned {
			return fmt.Errorf("There was a problem server side while disabling 5 min intraday histories for all stocks: %v\n", respErrors.ErrorReturned)
		}
		fmt.Printf("Successfullly disabled 5 min intraday histories for all stocks\n")

		return nil
	},
}

var disableQuotesCmd = &cobra.Command{
	Use:   "quotes",
	Short: "Disables quotes for a particular stock",
	Long:  `Disables gathering 1 min resolution quotes for a particular stock`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runDisableBlankCmd(cmd, args, "quotes", "quotes")
	},
}

var disableAllQuotesCmd = &cobra.Command{
	Use:   "quotes",
	Short: "Disables quotes for all stocks",
	Long:  `Disables gathering 1 min resolution quotes for all stocks`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runDisableAllBlankCmd(cmd, "quotes", "quotes")
	},
}

var disableStableQuotesCmd = &cobra.Command{
	Use:   "stableQuotes",
	Short: "Disables stable quotes for a particular stock",
	Long:  `Disables gathering stable quotes every day for a particular stock`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runDisableBlankCmd(cmd, args, "stableQuotes", "stableQuotes")
	},
}

var disableAllStableQuotesCmd = &cobra.Command{
	Use:   "stableQuotes",
	Short: "Disables stable quotes for all stocks",
	Long:  `Disables gathering stable quotes every day for all stocks`,
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
		resp, err := resty.R().Put(fmt.Sprintf("%v/stocks/disable/%v/%v", address, restEndpoint, s))
		if err != nil {
			return fmt.Errorf("There was an error while disabling %v for stock %v, error: %v\n", commandName, s, err)
		} else if resp.StatusCode() != http.StatusOK {
			return fmt.Errorf("There was an error while disabling %v for stock %v, statusCode: %v response: %v\n", commandName, s, resp.StatusCode(), resp)
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

	resp, err := resty.R().Put(fmt.Sprintf("%v/stocks/disable/all/%v", address, restEndpoint))
	if err != nil {
		return fmt.Errorf("There was an error while disabling %v for all stocks, error: %v\n", commandName, err)
	} else if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("There was an error while disabling %v for all stocks, statusCode: %v response: %v\n", commandName, resp.StatusCode(), resp)
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

func init() {
	rootCmd.AddCommand(disableCmd)
	disableCmd.AddCommand(disableDownloadingCmd)
	disableCmd.AddCommand(disableHistoryCmd)
	disableCmd.AddCommand(disableIntradayCmd)
	disableCmd.AddCommand(disableQuotesCmd)
	disableCmd.AddCommand(disableStableQuotesCmd)
	disableCmd.AddCommand(disableAllCmd)

	disableIntradayCmd.AddCommand(disableIntradayAllCmd)
	disableIntradayCmd.AddCommand(disableIntradayTickCmd)
	disableIntradayCmd.AddCommand(disableIntraday1MinCmd)
	disableIntradayCmd.AddCommand(disableIntraday5MinCmd)

	disableHistoryCmd.AddCommand(disableHistoryAllCmd)
	disableHistoryCmd.AddCommand(disableHistoryDailyCmd)
	disableHistoryCmd.AddCommand(disableHistoryWeeklyCmd)
	disableHistoryCmd.AddCommand(disableHistoryMonthlyCmd)
	disableHistoryCmd.AddCommand(disableHistoryYearlyCmd)

	disableAllCmd.AddCommand(disableAllAllCmd)
	disableAllCmd.AddCommand(disableAllDownloadingCmd)
	disableAllCmd.AddCommand(disableAllHistoryCmd)
	disableAllCmd.AddCommand(disableAllIntradayCmd)
	disableAllCmd.AddCommand(disableAllQuotesCmd)
	disableAllCmd.AddCommand(disableAllStableQuotesCmd)

	disableAllIntradayCmd.AddCommand(disableAllIntradayAllCmd)
	disableAllIntradayCmd.AddCommand(disableAllIntradayTickCmd)
	disableAllIntradayCmd.AddCommand(disableAllIntraday1MinCmd)
	disableAllIntradayCmd.AddCommand(disableAllIntraday5MinCmd)

	disableAllHistoryCmd.AddCommand(disableAllHistoryAllCmd)
	disableAllHistoryCmd.AddCommand(disableAllHistoryDailyCmd)
	disableAllHistoryCmd.AddCommand(disableAllHistoryWeeklyCmd)
	disableAllHistoryCmd.AddCommand(disableAllHistoryMonthlyCmd)
	disableAllHistoryCmd.AddCommand(disableAllHistoryYearlyCmd)

}
