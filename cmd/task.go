package cmd

import (
	"github.com/ba0f3/wazuh-cli/internal/client"
	"github.com/spf13/cobra"
)

var taskCmd = &cobra.Command{
	Use:   "task",
	Short: "Query Wazuh tasks (e.g. agent upgrades)",
}

var taskListCmd = &cobra.Command{
	Use:   "list",
	Short: "List tasks",
	RunE: func(cmd *cobra.Command, args []string) error {
		status, _ := cmd.Flags().GetString("status")
		params := &client.QueryParams{Filters: make(map[string]string)}
		if status != "" {
			params.Filters["status"] = status
		}
		resp, err := globalClient.Get("/tasks/status", params)
		if err != nil {
			return err
		}
		mustWrite(resp.Items())
		return nil
	},
}

func init() {
	rootCmd.AddCommand(taskCmd)
	taskCmd.AddCommand(taskListCmd)

	taskListCmd.Flags().String("status", "", "filter by status: In progress, Done, Failed, Cancelled")
}
