package cmd

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/nathanhack/sibyl/rest"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"gopkg.in/resty.v1"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

var useCSV, details, validOnly, invalidOnly bool
var startT, endT string

// showCmd represents the get command
var showCmd = &cobra.Command{
	Use:       "show",
	Short:     "Show will retrieve various values from SibylServer",
	Long:      `Show will retrieve various values from SibylServer`,
	ValidArgs: []string{"stocks", "history", "intraday"},
}

var showStocksCmd = &cobra.Command{
	Use:   "stocks",
	Short: "Lists all current stocks from the SibylServer",
	Long:  `Lists all current stocks from the SibylServer`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := cobra.NoArgs(cmd, args); err != nil {
			return err
		}

		address, err := cmd.Flags().GetString("serverAddress")
		if err != nil {
			return fmt.Errorf("Could not get server address from passed in arguments: %v\n", err)
		}

		if validOnly && invalidOnly {
			return fmt.Errorf("can not have validOnly only and invalid only, must choose one or the other")
		}

		str := fmt.Sprintf("%v/stocks/get", address)
		resp, err := resty.R().Get(str)
		if err != nil {
			return fmt.Errorf("error while trying to send request: %v", err)
		}

		if resp.StatusCode() != http.StatusOK {
			return fmt.Errorf("the response was not okay make sure URL is correct")
		}

		//next we unmarshal the json
		var stocksResponse rest.StocksResponse
		if err := json.Unmarshal([]byte(resp.Body()), &stocksResponse); err != nil {
			return fmt.Errorf("unmarshal failure: %v", err)
		}

		var data [][]string
		sort.Slice(stocksResponse.Stocks, func(i, j int) bool {
			return strings.Compare(stocksResponse.Stocks[i].Symbol, stocksResponse.Stocks[j].Symbol) < 0
		})

		for _, v := range stocksResponse.Stocks {
			stockRow := make([]string, 0)
			if validOnly && v.Validation != "validOnly" {
				continue
			}
			if invalidOnly && v.Validation != "invalid" {
				continue
			}

			if details {
				hasOptions := "no"
				if v.HasOptions {
					hasOptions = "yes"
				}
				stockRow = append(stockRow, v.Symbol)
				stockRow = append(stockRow, v.Name)
				stockRow = append(stockRow, v.Exchange)
				stockRow = append(stockRow, v.ExchangeDescription)
				stockRow = append(stockRow, hasOptions)
				stockRow = append(stockRow, v.Validation)
				stockRow = append(stockRow, v.DownloadStatus)
				stockRow = append(stockRow, v.QuotesStatus)
				stockRow = append(stockRow, v.StableQuotesStatus)
				stockRow = append(stockRow, v.HistoryStatus)
				stockRow = append(stockRow, v.IntradayHistoryStatus)
			} else {
				stockRow = append(stockRow, v.Symbol)
				stockRow = append(stockRow, v.Name)

			}
			data = append(data, stockRow)
		}

		detailsHeaders := []string{"Symbol", "Name", "Exchange", "Exchange Description", "Has Options", "Validation", "Download Status", "Quote Status", "Stable Quotes Status", "History Status", "Intraday Status"}
		normalHeaders := []string{"Symbol", "Name"}
		if useCSV {
			w := csv.NewWriter(os.Stdout)
			if details {
				if err := w.Write(detailsHeaders); err != nil {
					return fmt.Errorf("error writing data in csv format: %v", err)
				}
			} else {
				if err := w.Write(normalHeaders); err != nil {
					return fmt.Errorf("error writing data in csv format: %v", err)
				}
			}

			for _, row := range data {
				if err := w.Write(row); err != nil {
					return fmt.Errorf("error writing data in csv format: %v", err)
				}
			}
			w.Flush()
		} else {
			table := tablewriter.NewWriter(os.Stdout)
			if details {
				table.SetHeader(detailsHeaders)
			} else {
				table.SetHeader(normalHeaders)
			}
			table.SetAutoWrapText(false)
			table.AppendBulk(data)
			table.Render()
		}

		return nil
	},
}

func getTimestamp(str string, formatStr string, defaultValue int64) (int64, error) {
	//we have two cases either it's
	// a integer or it's a date/datetime string
	// it could also be blank meaning use the defaultValue
	if str == "" {
		return defaultValue, nil
	}

	if num, err := strconv.ParseInt(str, 10, 64); err != nil {
		if datetime, err := time.Parse(formatStr, str); err != nil {
			return defaultValue, fmt.Errorf("problem decoding time information from input: \"%v\", expected an integer or formating: %v", str, formatStr)
		} else {
			return datetime.Unix(), nil
		}
	} else {
		return num, nil
	}

}

var showHistoryCmd = &cobra.Command{
	Use:   "history [stockSymbol]",
	Args:  cobra.ExactArgs(1),
	Short: "Shows the history for a particular stock from the SibylServer",
	Long:  `Shows the history for a particular stock from the SibylServer`,
	RunE: func(cmd *cobra.Command, args []string) error {

		address, err := cmd.Flags().GetString("serverAddress")
		if err != nil {
			return fmt.Errorf("could not get server address from passed in arguments: %v", err)
		}

		startTimestamp, err := getTimestamp(startT, "2006-01-02", 0)
		if err != nil {
			return fmt.Errorf("could not get starting timestamp from passed in arguments: %v", err)
		}

		endTimestamp, err := getTimestamp(endT, "2006-01-02", time.Now().Local().Unix())
		if err != nil {
			return fmt.Errorf("could not get ending timestamp from passed in arguments: %v", err)
		}

		resp, err := resty.R().Get(fmt.Sprintf("%v/history/%v/%v/%v", address, args[0], startTimestamp, endTimestamp))
		if err != nil {
			return fmt.Errorf("error occured during request from server:%v", err)
		}

		var histories rest.Histories
		resp.Body()
		if err := json.Unmarshal(resp.Body(), &histories); err != nil {
			return fmt.Errorf("error occured during decoding the json from the server: %v", err)
		}

		if histories.ErrorState.ErrorReturned {
			return fmt.Errorf("error occured server side: %v", histories.ErrorState.Error)
		}

		var data [][]string

		for _, v := range histories.Histories {

			data = append(data, []string{
				time.Unix(v.Timestamp, 0).Format("20060102"), string(v.ClosePrice), v.HighPrice, string(v.LowPrice), string(v.OpenPrice), string(v.Volume)})
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Date Time", "Close", "High", "Low", "Open", "Volume"})

		table.AppendBulk(data)
		table.Render()
		return nil
	},
}

var showIntradayCmd = &cobra.Command{
	Use:   "intraday [stockSymbol]",
	Args:  cobra.ExactArgs(1),
	Short: "Shows the intraday history for a particular stock from the SibylServer",
	Long:  `Shows the intraday history for a particular stock from the SibylServer`,
	RunE: func(cmd *cobra.Command, args []string) error {
		address, err := cmd.Flags().GetString("serverAddress")
		if err != nil {
			return fmt.Errorf("could not get server address from passed in arguments: %v\n", err)
		}

		startTimestamp, err := getTimestamp(startT, "2006-01-02T15:04:05", 0)
		if err != nil {
			return fmt.Errorf("could not get starting timestamp from passed in arguments: %v", err)
		}

		endTimestamp, err := getTimestamp(endT, "2006-01-02T15:04:05", time.Now().Local().Unix())
		if err != nil {
			return fmt.Errorf("could not get ending timestamp from passed in arguments: %v", err)
		}

		resp, err := resty.R().Get(fmt.Sprintf("%v/intraday/%v/%v/%v", address, args[0], startTimestamp, endTimestamp))
		if err != nil {
			return fmt.Errorf("error occured during request from server:%v", err)
		}

		var intradays rest.Intradays
		resp.Body()
		if err := json.Unmarshal(resp.Body(), &intradays); err != nil {
			return fmt.Errorf("error occured during decoding the json from the server: %v", err)
		}

		if intradays.ErrorState.ErrorReturned {
			return fmt.Errorf("error occured server side: %v", intradays.ErrorState.Error)
		}

		var data [][]string

		for _, v := range intradays.Intradays {

			data = append(data, []string{
				time.Unix(v.Timestamp, 0).Format("20060102150405"), v.HighPrice, string(v.LastPrice), string(v.LowPrice), string(v.OpenPrice), string(v.Volume)})
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Date", "High", "Last", "Low", "Open", "Volume"})

		table.AppendBulk(data)
		table.Render()
		return nil
	},
}

var showCredsCmd = &cobra.Command{
	Use:   "creds",
	Short: "Shows the current cred information from the SibylServer",
	Long:  `Shows the current cred information from the SibylServer`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := cobra.NoArgs(cmd, args); err != nil {
			return err
		}

		address, err := cmd.Flags().GetString("serverAddress")
		if err != nil {
			return fmt.Errorf("could not get server address from passed in arguments: %v\n", err)
		}

		resp, err := resty.R().Get(fmt.Sprintf("%v/agent/creds", address))
		if err != nil {
			return fmt.Errorf("error occured during request from server:%v", err)
		}

		var creds rest.Creds
		resp.Body()
		if err := json.Unmarshal(resp.Body(), &creds); err != nil {
			return fmt.Errorf("error occured during decoding the json from the server: %v", err)
		}

		if creds.ErrorState.ErrorReturned {
			return fmt.Errorf("error occured server side: %v", creds.ErrorState.Error)
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Agent Selection",
			"Consumer Key",
			"Consumer Secret",
			"Token",
			"Token Secret",
			"Url Redirect",
			"Access Token",
			"Refresh Token",
			"Expire Timestamp"})

		table.Append([]string{creds.AgentSelection,
			creds.ConsumerKey,
			creds.ConsumerSecret,
			creds.Token,
			creds.TokenSecret,
			creds.UrlRedirect,
			creds.AccessToken,
			creds.RefreshToken,
			strconv.FormatInt(creds.ExpireTimestamp, 10),
			strconv.FormatInt(creds.RefreshExpireTimestamp, 10),
		})
		table.Render()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(showCmd)

	showCmd.AddCommand(showStocksCmd)
	showCmd.PersistentFlags().BoolVarP(&useCSV, "csv", "", false, "makes output CSV formatted")

	showStocksCmd.Flags().BoolVarP(&details, "details", "d", false, "show all details about the stocks")
	showStocksCmd.Flags().BoolVarP(&validOnly, "valid", "v", false, "show ONLY validated stocks")
	showStocksCmd.Flags().BoolVarP(&invalidOnly, "invalid", "i", false, "show ONLY invalidated stocks")

	showCmd.AddCommand(showHistoryCmd)

	showHistoryCmd.Flags().StringVarP(&startT, "startTimestamp", "s", "", "sets the start Date by integer timestamp or via a string using the following values & format: 2006-01-02 if omitted assumes 0 (beginning of time)")
	showHistoryCmd.Flags().StringVarP(&endT, "endTimestamp", "e", "", "sets the end Date by integer timestamp or via a string using the following values & format: 2006-01-02  if omitted assumes today")

	showCmd.AddCommand(showIntradayCmd)
	showIntradayCmd.Flags().StringVarP(&startT, "startTimestamp", "s", "", "sets the start Date & Time by integer timestamp or via a string using the following values & format: 2006-01-02T15:04:05 if omitted assumes 0 (beginning of time)")
	showIntradayCmd.Flags().StringVarP(&endT, "endTimestamp", "e", "", "sets the end Date & Time by integer timestamp or via a string using the following values & format: 2006-01-02T15:04:05 if omitted assumes today")

	showCmd.AddCommand(showCredsCmd)
}
