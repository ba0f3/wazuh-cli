package cmd

import (
	"fmt"

	"github.com/ba0f3/wazuh-cli/internal/client"
	"github.com/spf13/cobra"
)

var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: "Manage Wazuh agents",
}

var agentListCmd = &cobra.Command{
	Use:   "list",
	Short: "List agents",
	RunE: func(cmd *cobra.Command, args []string) error {
		status, _ := cmd.Flags().GetString("status")
		group, _ := cmd.Flags().GetString("group")
		query, _ := cmd.Flags().GetString("query")
		limit, _ := cmd.Flags().GetInt("limit")
		offset, _ := cmd.Flags().GetInt("offset")

		params := &client.QueryParams{Limit: limit, Offset: offset, Query: query}
		if status != "" {
			params.Filters = map[string]string{"status": status}
		}
		if group != "" {
			if params.Filters == nil {
				params.Filters = make(map[string]string)
			}
			params.Filters["group"] = group
		}

		resp, err := globalClient.Get("/agents", params)
		if err != nil {
			return err
		}
		mustWrite(resp.Items())
		return nil
	},
}

var agentGetCmd = &cobra.Command{
	Use:   "get <AGENT_ID>",
	Short: "Get a specific agent by ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := globalClient.Get(fmt.Sprintf("/agents/%s", args[0]), nil)
		if err != nil {
			return err
		}
		mustWrite(resp.Items())
		return nil
	},
}

var agentDeleteCmd = &cobra.Command{
	Use:   "delete <AGENT_ID>",
	Short: "Delete an agent",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		params := &client.QueryParams{Filters: map[string]string{"agents_list": args[0]}}
		resp, err := globalClient.Delete("/agents", params)
		if err != nil {
			return err
		}
		mustWrite(resp.Items())
		return nil
	},
}

var agentRestartCmd = &cobra.Command{
	Use:   "restart <AGENT_ID>",
	Short: "Restart an agent",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := globalClient.Put(fmt.Sprintf("/agents/%s/restart", args[0]), nil)
		if err != nil {
			return err
		}
		mustWrite(resp.Items())
		return nil
	},
}

var agentSummaryCmd = &cobra.Command{
	Use:   "summary",
	Short: "Get agent status summary",
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := globalClient.Get("/agents/summary/status", nil)
		if err != nil {
			return err
		}
		mustWrite(resp.Data)
		return nil
	},
}

var agentKeyCmd = &cobra.Command{
	Use:   "key <AGENT_ID>",
	Short: "Get agent key",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := globalClient.Get(fmt.Sprintf("/agents/%s/key", args[0]), nil)
		if err != nil {
			return err
		}
		mustWrite(resp.Items())
		return nil
	},
}

var agentUpgradeCmd = &cobra.Command{
	Use:   "upgrade <AGENT_ID>",
	Short: "Upgrade an agent",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		version, _ := cmd.Flags().GetString("version")
		body := map[string]interface{}{"agents_list": []string{args[0]}}
		if version != "" {
			body["version"] = version
		}
		resp, err := globalClient.Put("/agents/upgrade", body)
		if err != nil {
			return err
		}
		mustWrite(resp.Items())
		return nil
	},
}

var agentConfigCmd = &cobra.Command{
	Use:   "config <AGENT_ID>",
	Short: "Get agent active configuration",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		component, _ := cmd.Flags().GetString("component")
		configuration, _ := cmd.Flags().GetString("configuration")
		path := fmt.Sprintf("/agents/%s/config/%s/%s", args[0], component, configuration)
		resp, err := globalClient.Get(path, nil)
		if err != nil {
			return err
		}
		mustWrite(resp.Data)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(agentCmd)
	agentCmd.AddCommand(agentListCmd, agentGetCmd, agentDeleteCmd, agentRestartCmd,
		agentSummaryCmd, agentKeyCmd, agentUpgradeCmd, agentConfigCmd)

	agentListCmd.Flags().String("status", "", "filter by status: active, disconnected, pending, never_connected")
	agentListCmd.Flags().String("group", "", "filter by group name")
	agentListCmd.Flags().String("query", "", "WQL query filter")
	agentListCmd.Flags().Int("limit", 500, "maximum number of results")
	agentListCmd.Flags().Int("offset", 0, "pagination offset")

	agentUpgradeCmd.Flags().String("version", "", "target version")

	agentConfigCmd.Flags().String("component", "", "component name (e.g. agent, logcollector)")
	agentConfigCmd.Flags().String("configuration", "", "configuration name (e.g. client, localfile)")
	_ = agentConfigCmd.MarkFlagRequired("component")
	_ = agentConfigCmd.MarkFlagRequired("configuration")
}
