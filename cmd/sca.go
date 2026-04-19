package cmd

import (
	"fmt"

	"github.com/ba0f3/wazuh-cli/internal/client"
	"github.com/spf13/cobra"
)

var scaCmd = &cobra.Command{
	Use:   "sca",
	Short: "Security Configuration Assessment (SCA)",
}

var scaListCmd = &cobra.Command{
	Use:   "list <AGENT_ID>",
	Short: "List SCA policies for an agent",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := globalClient.Get(fmt.Sprintf("/sca/%s", args[0]), nil)
		if err != nil {
			return err
		}
		mustWrite(resp.Items())
		return nil
	},
}

var scaGetCmd = &cobra.Command{
	Use:   "get <AGENT_ID> <POLICY_ID>",
	Short: "Get a specific SCA policy result",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		params := &client.QueryParams{Filters: map[string]string{"policy_id": args[1]}}
		resp, err := globalClient.Get(fmt.Sprintf("/sca/%s", args[0]), params)
		if err != nil {
			return err
		}
		mustWrite(resp.Items())
		return nil
	},
}

var scaChecksCmd = &cobra.Command{
	Use:   "checks <AGENT_ID> <POLICY_ID>",
	Short: "List checks for an SCA policy",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := globalClient.Get(fmt.Sprintf("/sca/%s/checks/%s", args[0], args[1]), nil)
		if err != nil {
			return err
		}
		mustWrite(resp.Items())
		return nil
	},
}

func init() {
	rootCmd.AddCommand(scaCmd)
	scaCmd.AddCommand(scaListCmd, scaGetCmd, scaChecksCmd)
}
