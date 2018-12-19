package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/nathanhack/sibyl/rest"
	"github.com/spf13/cobra"
	"gopkg.in/resty.v1"
	"io/ioutil"
	"net/http"
)

var dataBaseFilename string
var databaseCmd = &cobra.Command{
	Use:   "database",
	Short: "Used to backup/restore records",
	Long:  `Used to backup/restore records`,
}

var databaseDownloadCmd = &cobra.Command{
	Use:   "download",
	Short: "To download records as a backup",
	Long:  "To download records as a backup",
}

var databaseDownloadHistory = &cobra.Command{
	Use:   "history",
	Short: "will download the history for all stocks in Sibyl's backend",
	Long:  "will download the history for all stocks in Sibyl's backend",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := cobra.NoArgs(cmd, args); err != nil {
			return err
		}

		address, err := cmd.Flags().GetString("serverAddress")
		if err != nil {
			return fmt.Errorf("Could not get server address from passed in arguments: %v\n", err)
		}

		resp, err := resty.R().Get(fmt.Sprintf("%v/database/download/history", address))
		if err != nil {
			return fmt.Errorf("There was an error sending request to server: %v\n", err)
		} else if resp.StatusCode() != http.StatusOK {
			return fmt.Errorf("There was an error statusCode: %v  response: %v\n", resp.StatusCode(), resp)
		} else {
			var databaseRecords rest.DatabaseRecords
			err := json.Unmarshal(resp.Body(), &databaseRecords)
			if err != nil {
				return fmt.Errorf("There was a problem parsing the server response: %v  had error:%v\n", string(resp.Body()), err)
			} else {
				if databaseRecords.ErrorState.ErrorReturned {
					return fmt.Errorf("There was a problem server side: %v\n", databaseRecords.ErrorState.Error)
				} else {
					if err := ioutil.WriteFile(dataBaseFilename, []byte(databaseRecords.Histories), 0644); err != nil {
						return fmt.Errorf("There was a problem while writing the output file: %v", err)
					}
					fmt.Printf("Successfully downloaded History to file: %v\n", dataBaseFilename)
				}
			}
		}
		return nil
	},
}

var databaseDownloadIntraday = &cobra.Command{
	Use:   "intraday",
	Short: "will download the intraday for all stocks in Sibyl's backend",
	Long:  "will download the intraday for all stocks in Sibyl's backend",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := cobra.NoArgs(cmd, args); err != nil {
			return err
		}

		address, err := cmd.Flags().GetString("serverAddress")
		if err != nil {
			return fmt.Errorf("Could not get server address from passed in arguments: %v\n", err)
		}

		resp, err := resty.R().Get(fmt.Sprintf("%v/database/download/intraday", address))
		if err != nil {
			return fmt.Errorf("There was an error sending request to server: %v\n", err)
		} else if resp.StatusCode() != http.StatusOK {
			return fmt.Errorf("There was an error statusCode: %v  response: %v\n", resp.StatusCode(), resp)
		} else {
			var databaseRecords rest.DatabaseRecords
			err := json.Unmarshal(resp.Body(), &databaseRecords)
			if err != nil {
				return fmt.Errorf("There was a problem parsing the server response: %v  had error:%v\n", string(resp.Body()), err)
			} else {
				if databaseRecords.ErrorState.ErrorReturned {
					return fmt.Errorf("There was a problem server side: %v\n", databaseRecords.ErrorState.Error)
				} else {
					if err := ioutil.WriteFile(dataBaseFilename, []byte(databaseRecords.Intradays), 0644); err != nil {
						return fmt.Errorf("There was a problem while writing the output file: %v", err)
					}
					fmt.Printf("Successfully downloaded Intraday to file: %v\n", dataBaseFilename)
				}
			}
		}
		return nil
	},
}

var databaseDownloadStockCmd = &cobra.Command{
	Use:   "stocks",
	Short: "will download stock related data from Sibyl's backend",
	Long:  "will download stock related data from Sibyl's backend",
}

var databaseDownloadStockQuoteCmd = &cobra.Command{
	Use:   "quote",
	Short: "will download the stock quotes for all stocks in Sibyl's backend",
	Long:  "will download the stock quotes for all stocks in Sibyl's backend",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := cobra.NoArgs(cmd, args); err != nil {
			return err
		}

		address, err := cmd.Flags().GetString("serverAddress")
		if err != nil {
			return fmt.Errorf("Could not get server address from passed in arguments: %v\n", err)
		}

		resp, err := resty.R().Get(fmt.Sprintf("%v/database/download/stocks/quote", address))
		if err != nil {
			return fmt.Errorf("There was an error sending request to server: %v\n", err)
		} else if resp.StatusCode() != http.StatusOK {
			return fmt.Errorf("There was an error statusCode: %v  response: %v\n", resp.StatusCode(), resp)
		} else {
			var databaseRecords rest.DatabaseRecords
			err := json.Unmarshal(resp.Body(), &databaseRecords)
			if err != nil {
				return fmt.Errorf("There was a problem parsing the server response: %v  had error:%v\n", string(resp.Body()), err)
			} else {
				if databaseRecords.ErrorState.ErrorReturned {
					return fmt.Errorf("There was a problem server side: %v\n", databaseRecords.ErrorState.Error)
				} else {
					if err := ioutil.WriteFile(dataBaseFilename, []byte(databaseRecords.StockQuotes), 0644); err != nil {
						return fmt.Errorf("There was a problem while writing the output file: %v", err)
					}
					fmt.Printf("Successfully downloaded Stock Quotes to file: %v\n", dataBaseFilename)
				}
			}
		}
		return nil
	},
}

var databaseDownloadStockStableCmd = &cobra.Command{
	Use:   "stable",
	Short: "will download the stock stable quotes for all stocks in Sibyl's backend",
	Long:  "will download the stock stable quotes for all stocks in Sibyl's backend",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := cobra.NoArgs(cmd, args); err != nil {
			return err
		}

		address, err := cmd.Flags().GetString("serverAddress")
		if err != nil {
			return fmt.Errorf("Could not get server address from passed in arguments: %v\n", err)
		}

		resp, err := resty.R().Get(fmt.Sprintf("%v/database/download/stocks/stable", address))
		if err != nil {
			return fmt.Errorf("There was an error sending request to server: %v\n", err)
		} else if resp.StatusCode() != http.StatusOK {
			return fmt.Errorf("There was an error statusCode: %v  response: %v\n", resp.StatusCode(), resp)
		} else {
			var databaseRecords rest.DatabaseRecords
			err := json.Unmarshal(resp.Body(), &databaseRecords)
			if err != nil {
				return fmt.Errorf("There was a problem parsing the server response: %v  had error:%v\n", string(resp.Body()), err)
			} else {
				if databaseRecords.ErrorState.ErrorReturned {
					return fmt.Errorf("There was a problem server side: %v\n", databaseRecords.ErrorState.Error)
				} else {
					if err := ioutil.WriteFile(dataBaseFilename, []byte(databaseRecords.StockStableQuotes), 0644); err != nil {
						return fmt.Errorf("There was a problem while writing the output file: %v", err)
					}
					fmt.Printf("Successfully downloaded Stock Stable Quotes to file: %v\n", dataBaseFilename)
				}
			}
		}
		return nil
	},
}

var databaseDownloadOptionsCmd = &cobra.Command{
	Use:   "options",
	Short: "will download stock option related data from Sibyl's backend",
	Long:  "will download stock option related data from Sibyl's backend",
}

var databaseDownloadOptionsQuoteCmd = &cobra.Command{
	Use:   "quote",
	Short: "will download the options quotes for all stocks in Sibyl's backend",
	Long:  "will download the options quotes for all stocks in Sibyl's backend",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := cobra.NoArgs(cmd, args); err != nil {
			return err
		}

		address, err := cmd.Flags().GetString("serverAddress")
		if err != nil {
			return fmt.Errorf("Could not get server address from passed in arguments: %v\n", err)
		}

		resp, err := resty.R().Get(fmt.Sprintf("%v/database/download/options/quote", address))
		if err != nil {
			return fmt.Errorf("There was an error sending request to server: %v\n", err)
		} else if resp.StatusCode() != http.StatusOK {
			return fmt.Errorf("There was an error statusCode: %v  response: %v\n", resp.StatusCode(), resp)
		} else {
			var databaseRecords rest.DatabaseRecords
			err := json.Unmarshal(resp.Body(), &databaseRecords)
			if err != nil {
				return fmt.Errorf("There was a problem parsing the server response: %v  had error:%v\n", string(resp.Body()), err)
			} else {
				if databaseRecords.ErrorState.ErrorReturned {
					return fmt.Errorf("There was a problem server side: %v\n", databaseRecords.ErrorState.Error)
				} else {
					if err := ioutil.WriteFile(dataBaseFilename, []byte(databaseRecords.OptionQuotes), 0644); err != nil {
						return fmt.Errorf("There was a problem while writing the output file: %v", err)
					}
					fmt.Printf("Successfully downloaded Option Quotes to file: %v\n", dataBaseFilename)
				}
			}
		}
		return nil
	},
}

var databaseDownloadOptionsStableCmd = &cobra.Command{
	Use:   "stable",
	Short: "will download the options stable quotes for all stocks in Sibyl's backend",
	Long:  "will download the options stable quotes for all stocks in Sibyl's backend",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := cobra.NoArgs(cmd, args); err != nil {
			return err
		}

		address, err := cmd.Flags().GetString("serverAddress")
		if err != nil {
			return fmt.Errorf("Could not get server address from passed in arguments: %v\n", err)
		}

		resp, err := resty.R().Get(fmt.Sprintf("%v/database/download/options/stable", address))
		if err != nil {
			return fmt.Errorf("There was an error sending request to server: %v\n", err)
		} else if resp.StatusCode() != http.StatusOK {
			return fmt.Errorf("There was an error statusCode: %v  response: %v\n", resp.StatusCode(), resp)
		} else {
			var databaseRecords rest.DatabaseRecords
			err := json.Unmarshal(resp.Body(), &databaseRecords)
			if err != nil {
				return fmt.Errorf("There was a problem parsing the server response: %v  had error:%v\n", string(resp.Body()), err)
			} else {
				if databaseRecords.ErrorState.ErrorReturned {
					return fmt.Errorf("There was a problem server side: %v\n", databaseRecords.ErrorState.Error)
				} else {
					if err := ioutil.WriteFile(dataBaseFilename, []byte(databaseRecords.OptionStableQuotes), 0644); err != nil {
						return fmt.Errorf("There was a problem while writing the output file: %v", err)
					}
					fmt.Printf("Successfully downloaded Option Stable Quotes to file: %v\n", dataBaseFilename)
				}
			}
		}
		return nil
	},
}

var databaseUploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "To upload records to restore from backup",
	Long:  "To upload records to restore from backup",
}

var databaseUploadHistory = &cobra.Command{
	Use:   "history",
	Short: "will download the history for all stocks in Sibyl's backend",
	Long:  "will download the history for all stocks in Sibyl's backend",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := cobra.NoArgs(cmd, args); err != nil {
			return err
		}

		address, err := cmd.Flags().GetString("serverAddress")
		if err != nil {
			return fmt.Errorf("Could not get server address from passed in arguments: %v\n", err)
		}

		bytes, err := ioutil.ReadFile(dataBaseFilename)
		if err != nil {
			return fmt.Errorf("Could not read file to upload: %v", err)
		}

		jsonBytes, err := json.Marshal(rest.DatabaseRecords{Histories: string(bytes)})
		if err != nil {
			return fmt.Errorf("An error occured while encoding for upload: %v", err)
		}

		resp, err := resty.R().SetBody(jsonBytes).Post(fmt.Sprintf("%v/database/upload/history", address))
		if err != nil {
			return fmt.Errorf("There was an error sending request to server: %v\n", err)
		} else if resp.StatusCode() != http.StatusOK {
			return fmt.Errorf("There was an error statusCode: %v  response: %v\n", resp.StatusCode(), resp)
		} else {
			fmt.Println("Successfully uploaded History")
		}
		return nil
	},
}

var databaseUploadIntraday = &cobra.Command{
	Use:   "intraday",
	Short: "will download the history for all stocks in Sibyl's backend",
	Long:  "will download the history for all stocks in Sibyl's backend",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := cobra.NoArgs(cmd, args); err != nil {
			return err
		}

		address, err := cmd.Flags().GetString("serverAddress")
		if err != nil {
			return fmt.Errorf("Could not get server address from passed in arguments: %v\n", err)
		}

		bytes, err := ioutil.ReadFile(dataBaseFilename)
		if err != nil {
			return fmt.Errorf("Could not read file to upload: %v", err)
		}

		jsonBytes, err := json.Marshal(rest.DatabaseRecords{Intradays: string(bytes)})
		if err != nil {
			return fmt.Errorf("An error occured while encoding for upload: %v", err)
		}

		resp, err := resty.R().SetBody(jsonBytes).Post(fmt.Sprintf("%v/database/upload/intraday", address))
		if err != nil {
			return fmt.Errorf("There was an error sending request to server: %v\n", err)
		} else if resp.StatusCode() != http.StatusOK {
			return fmt.Errorf("There was an error statusCode: %v  response: %v\n", resp.StatusCode(), resp)
		} else {
			fmt.Println("Successfully uploaded Intraday")
		}
		return nil
	},
}

var databaseUploadStockCmd = &cobra.Command{
	Use:   "stocks",
	Short: "",
	Long:  "",
}

var databaseUploadStockQuoteCmd = &cobra.Command{
	Use:   "quote",
	Short: "",
	Long:  "",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := cobra.NoArgs(cmd, args); err != nil {
			return err
		}

		address, err := cmd.Flags().GetString("serverAddress")
		if err != nil {
			return fmt.Errorf("Could not get server address from passed in arguments: %v\n", err)
		}

		bytes, err := ioutil.ReadFile(dataBaseFilename)
		if err != nil {
			return fmt.Errorf("Could not read file to upload: %v", err)
		}

		jsonBytes, err := json.Marshal(rest.DatabaseRecords{StockQuotes: string(bytes)})
		if err != nil {
			return fmt.Errorf("An error occured while encoding for upload: %v", err)
		}

		resp, err := resty.R().SetBody(jsonBytes).Post(fmt.Sprintf("%v/database/upload/stocks/quote", address))
		if err != nil {
			return fmt.Errorf("There was an error sending request to server: %v\n", err)
		} else if resp.StatusCode() != http.StatusOK {
			return fmt.Errorf("There was an error statusCode: %v  response: %v\n", resp.StatusCode(), resp)
		} else {
			fmt.Println("Successfully uploaded Stock Quotes")
		}
		return nil
	},
}

var databaseUploadStockStableCmd = &cobra.Command{
	Use:   "stable",
	Short: "",
	Long:  "",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := cobra.NoArgs(cmd, args); err != nil {
			return err
		}

		address, err := cmd.Flags().GetString("serverAddress")
		if err != nil {
			return fmt.Errorf("Could not get server address from passed in arguments: %v\n", err)
		}

		bytes, err := ioutil.ReadFile(dataBaseFilename)
		if err != nil {
			return fmt.Errorf("Could not read file to upload: %v", err)
		}

		jsonBytes, err := json.Marshal(rest.DatabaseRecords{StockStableQuotes: string(bytes)})
		if err != nil {
			return fmt.Errorf("An error occured while encoding for upload: %v", err)
		}

		resp, err := resty.R().SetBody(jsonBytes).Post(fmt.Sprintf("%v/database/upload/stocks/stable", address))
		if err != nil {
			return fmt.Errorf("There was an error sending request to server: %v\n", err)
		} else if resp.StatusCode() != http.StatusOK {
			return fmt.Errorf("There was an error statusCode: %v  response: %v\n", resp.StatusCode(), resp)
		} else {
			fmt.Println("Successfully uploaded Stock Stable Quotes")
		}
		return nil
	},
}

var databaseUploadOptionsCmd = &cobra.Command{
	Use:   "options",
	Short: "",
	Long:  "",
}

var databaseUploadOptionsQuoteCmd = &cobra.Command{
	Use:   "quote",
	Short: "",
	Long:  "",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := cobra.NoArgs(cmd, args); err != nil {
			return err
		}

		address, err := cmd.Flags().GetString("serverAddress")
		if err != nil {
			return fmt.Errorf("Could not get server address from passed in arguments: %v\n", err)
		}

		bytes, err := ioutil.ReadFile(dataBaseFilename)
		if err != nil {
			return fmt.Errorf("Could not read file to upload: %v", err)
		}

		jsonBytes, err := json.Marshal(rest.DatabaseRecords{OptionQuotes: string(bytes)})
		if err != nil {
			return fmt.Errorf("An error occured while encoding for upload: %v", err)
		}

		resp, err := resty.R().SetBody(jsonBytes).Post(fmt.Sprintf("%v/database/upload/options/quote", address))
		if err != nil {
			return fmt.Errorf("There was an error sending request to server: %v\n", err)
		} else if resp.StatusCode() != http.StatusOK {
			return fmt.Errorf("There was an error statusCode: %v response: %v\n", resp.StatusCode(), resp)
		} else {
			fmt.Println("Successfully uploaded Option Quotes")
		}
		return nil
	},
}

var databaseUploadOptionsStableCmd = &cobra.Command{
	Use:   "stable",
	Short: "",
	Long:  "",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := cobra.NoArgs(cmd, args); err != nil {
			return err
		}

		address, err := cmd.Flags().GetString("serverAddress")
		if err != nil {
			return fmt.Errorf("Could not get server address from passed in arguments: %v\n", err)
		}

		bytes, err := ioutil.ReadFile(dataBaseFilename)
		if err != nil {
			return fmt.Errorf("Could not read file to upload: %v", err)
		}

		jsonBytes, err := json.Marshal(rest.DatabaseRecords{OptionStableQuotes: string(bytes)})
		if err != nil {
			return fmt.Errorf("An error occured while encoding for upload: %v", err)
		}

		resp, err := resty.R().SetBody(jsonBytes).Post(fmt.Sprintf("%v/database/upload/options/stable", address))
		if err != nil {
			return fmt.Errorf("There was an error sending request to server: %v\n", err)
		} else if resp.StatusCode() != http.StatusOK {
			return fmt.Errorf("There was an error statusCode: %v response: %v\n", resp.StatusCode(), resp)
		} else {
			fmt.Println("Successfully uploaded Option Stable Quotes")
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(databaseCmd)
	databaseCmd.AddCommand(databaseDownloadCmd, databaseUploadCmd)

	databaseCmd.PersistentFlags().StringVarP(&dataBaseFilename, "file", "f", "./file.backup", "the pathname to input or output file")

	databaseDownloadCmd.AddCommand(databaseDownloadHistory, databaseDownloadIntraday, databaseDownloadStockCmd, databaseDownloadOptionsCmd)
	databaseDownloadStockCmd.AddCommand(databaseDownloadStockQuoteCmd, databaseDownloadStockStableCmd)
	databaseDownloadOptionsCmd.AddCommand(databaseDownloadOptionsQuoteCmd, databaseDownloadOptionsStableCmd)

	databaseUploadCmd.AddCommand(databaseUploadHistory, databaseUploadIntraday, databaseUploadStockCmd, databaseUploadOptionsCmd)
	databaseUploadStockCmd.AddCommand(databaseUploadStockQuoteCmd, databaseUploadStockStableCmd)
	databaseUploadOptionsCmd.AddCommand(databaseUploadOptionsQuoteCmd, databaseUploadOptionsStableCmd)

}
