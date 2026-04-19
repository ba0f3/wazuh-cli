package cmd

import (
	"github.com/ba0f3/wazuh-cli/internal/client"
	"github.com/spf13/cobra"
)

var ruleCmd = &cobra.Command{
	Use:   "rule",
	Short: "Manage Wazuh rules",
}

var ruleListCmd = &cobra.Command{
	Use:   "list",
	Short: "List rules",
	RunE: func(cmd *cobra.Command, args []string) error {
		status, _ := cmd.Flags().GetString("status")
		group, _ := cmd.Flags().GetString("group")
		level, _ := cmd.Flags().GetString("level")
		file, _ := cmd.Flags().GetString("file")
		limit, _ := cmd.Flags().GetInt("limit")
		offset, _ := cmd.Flags().GetInt("offset")

		params := &client.QueryParams{Limit: limit, Offset: offset, Filters: make(map[string]string)}
		if status != "" {
			params.Filters["status"] = status
		}
		if group != "" {
			params.Filters["group"] = group
		}
		if level != "" {
			params.Filters["level"] = level
		}
		if file != "" {
			params.Filters["filename"] = file
		}

		resp, err := globalClient.Get("/rules", params)
		if err != nil {
			return err
		}
		mustWrite(resp.Items())
		return nil
	},
}

var ruleGetCmd = &cobra.Command{
	Use:   "get <RULE_ID>",
	Short: "Get a specific rule by ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		params := &client.QueryParams{Filters: map[string]string{"rule_ids": args[0]}}
		resp, err := globalClient.Get("/rules", params)
		if err != nil {
			return err
		}
		mustWrite(resp.Items())
		return nil
	},
}

var ruleFilesCmd = &cobra.Command{
	Use:   "files",
	Short: "List rule files",
	RunE: func(cmd *cobra.Command, args []string) error {
		status, _ := cmd.Flags().GetString("status")
		params := &client.QueryParams{}
		if status != "" {
			params.Filters = map[string]string{"status": status}
		}
		resp, err := globalClient.Get("/rules/files", params)
		if err != nil {
			return err
		}
		mustWrite(resp.Items())
		return nil
	},
}

var ruleGroupsCmd = &cobra.Command{
	Use:   "groups",
	Short: "List all rule groups",
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := globalClient.Get("/rules/groups", nil)
		if err != nil {
			return err
		}
		mustWrite(resp.Items())
		return nil
	},
}

func init() {
	rootCmd.AddCommand(ruleCmd)
	ruleCmd.AddCommand(ruleListCmd, ruleGetCmd, ruleFilesCmd, ruleGroupsCmd)

	ruleListCmd.Flags().String("status", "", "filter by status: enabled, disabled")
	ruleListCmd.Flags().String("group", "", "filter by rule group")
	ruleListCmd.Flags().String("level", "", "filter by level range (e.g. 5-10)")
	ruleListCmd.Flags().String("file", "", "filter by file name")
	ruleListCmd.Flags().Int("limit", 500, "maximum results")
	ruleListCmd.Flags().Int("offset", 0, "pagination offset")

	ruleFilesCmd.Flags().String("status", "", "filter by status: enabled, disabled")
}
