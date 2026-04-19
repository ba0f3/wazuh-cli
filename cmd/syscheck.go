package cmd

import (
	"fmt"

	"github.com/ba0f3/wazuh-cli/internal/client"
	"github.com/spf13/cobra"
)

var syscheckCmd = &cobra.Command{
	Use:   "syscheck",
	Short: "File Integrity Monitoring (FIM/Syscheck)",
}

var syscheckListCmd = &cobra.Command{
	Use:   "list <AGENT_ID>",
	Short: "List FIM database entries for an agent",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		file, _ := cmd.Flags().GetString("file")
		params := &client.QueryParams{Filters: make(map[string]string)}
		if file != "" {
			params.Filters["filename"] = file
		}
		resp, err := globalClient.Get(fmt.Sprintf("/syscheck/%s", args[0]), params)
		if err != nil {
			return err
		}
		mustWrite(resp.Items())
		return nil
	},
}

var syscheckRunCmd = &cobra.Command{
	Use:   "run <AGENT_ID>",
	Short: "Run FIM scan on an agent",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		params := &client.QueryParams{Filters: map[string]string{"agents_list": args[0]}}
		resp, err := globalClient.Put("/syscheck", params)
		if err != nil {
			return err
		}
		mustWrite(resp.Items())
		return nil
	},
}

var syscheckClearCmd = &cobra.Command{
	Use:   "clear <AGENT_ID>",
	Short: "Clear FIM database for an agent",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		params := &client.QueryParams{Filters: map[string]string{"agents_list": args[0]}}
		resp, err := globalClient.Delete("/syscheck", params)
		if err != nil {
			return err
		}
		mustWrite(resp.Items())
		return nil
	},
}

var syscheckLastScanCmd = &cobra.Command{
	Use:   "last-scan <AGENT_ID>",
	Short: "Get the time of the last FIM scan",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := globalClient.Get(fmt.Sprintf("/syscheck/%s/last_scan", args[0]), nil)
		if err != nil {
			return err
		}
		mustWrite(resp.Data)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(syscheckCmd)
	syscheckCmd.AddCommand(syscheckListCmd, syscheckRunCmd, syscheckClearCmd, syscheckLastScanCmd)

	syscheckListCmd.Flags().String("file", "", "filter by file path")
}
