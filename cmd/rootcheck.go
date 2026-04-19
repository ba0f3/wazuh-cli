package cmd

import (
	"fmt"

	"github.com/ba0f3/wazuh-cli/internal/client"
	"github.com/spf13/cobra"
)

var rootcheckCmd = &cobra.Command{
	Use:   "rootcheck",
	Short: "Rootcheck (policy and anomaly detection)",
}

var rootcheckListCmd = &cobra.Command{
	Use:   "list <AGENT_ID>",
	Short: "List rootcheck results for an agent",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := globalClient.Get(fmt.Sprintf("/rootcheck/%s", args[0]), nil)
		if err != nil {
			return err
		}
		mustWrite(resp.Items())
		return nil
	},
}

var rootcheckRunCmd = &cobra.Command{
	Use:   "run <AGENT_ID>",
	Short: "Run rootcheck on an agent",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		params := &client.QueryParams{Filters: map[string]string{"agents_list": args[0]}}
		resp, err := globalClient.Put("/rootcheck", params)
		if err != nil {
			return err
		}
		mustWrite(resp.Items())
		return nil
	},
}

var rootcheckClearCmd = &cobra.Command{
	Use:   "clear <AGENT_ID>",
	Short: "Clear rootcheck database for an agent",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		params := &client.QueryParams{Filters: map[string]string{"agents_list": args[0]}}
		resp, err := globalClient.Delete("/rootcheck", params)
		if err != nil {
			return err
		}
		mustWrite(resp.Items())
		return nil
	},
}

var rootcheckLastScanCmd = &cobra.Command{
	Use:   "last-scan <AGENT_ID>",
	Short: "Get the last rootcheck scan time",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := globalClient.Get(fmt.Sprintf("/rootcheck/%s/last_scan", args[0]), nil)
		if err != nil {
			return err
		}
		mustWrite(resp.Data)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(rootcheckCmd)
	rootcheckCmd.AddCommand(rootcheckListCmd, rootcheckRunCmd, rootcheckClearCmd, rootcheckLastScanCmd)
}
