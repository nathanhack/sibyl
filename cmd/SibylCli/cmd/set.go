package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/nathanhack/sibyl/rest"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/resty.v1"
)

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Set is used to identify the discount broker and auth information",
	Long:  `Set is used to identify the discount broker and auth information`,
}

var setAllyCredsCmd = &cobra.Command{
	Use:   "ally consumerKey consumerSecret oAuthToken oAuthTokenSecret",
	Short: "Sets the cred values for Ally Invest and sets it as the default broker to use",
	Long:  `Sets the cred values for Ally Invest and sets it as the default broker to use`,
	Args:  cobra.ExactArgs(4),
	RunE: func(cmd *cobra.Command, args []string) error {
		address, err := cmd.Flags().GetString("serverAddress")
		if err != nil {
			return fmt.Errorf("Could not get server address from passed in arguments: %v\n", err)
		}

		if resp, err := resty.R().Post(fmt.Sprintf("%v/agent/ally/%v/%v/%v/%v", address, args[0], args[1], args[2], args[3])); err != nil {
			return fmt.Errorf("There was an error while setting creds, error: %v", err)
		} else {
			var respErrors rest.ErrorState
			err := json.Unmarshal(resp.Body(), &respErrors)
			if err != nil {
				return fmt.Errorf("There was a problem parsing the server response while setting creds for Ally, error:%v and response body: %v\n", err, string(resp.Body()))
			} else {
				if respErrors.ErrorReturned {
					return fmt.Errorf("There was a problem server side while setting creds: %v\n", respErrors.ErrorReturned)
				}
			}
		}

		if resp, err := resty.R().Put(fmt.Sprintf("%v/agent/use/ally", address)); err != nil {
			return fmt.Errorf("There was an error while setting to use Ally, error: %v", err)
		} else {
			var respErrors rest.ErrorState
			err := json.Unmarshal(resp.Body(), &respErrors)
			if err != nil {
				return fmt.Errorf("There was a problem parsing the server response while setting Ally as the default agent, error:%v and response body: %v\n", err, string(resp.Body()))
			} else {
				if respErrors.ErrorReturned {
					return fmt.Errorf("There was a problem server side while setting Ally as the default agent: %v\n", respErrors.ErrorReturned)
				}
			}
		}

		logrus.Infof("Successfully set Ally credentials and Ally to the default agent.")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(setCmd)

	setCmd.AddCommand(setAllyCredsCmd)
}
