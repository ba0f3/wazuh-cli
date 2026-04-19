package cmd

import (
	"github.com/spf13/cobra"
)

var activeResponseCmd = &cobra.Command{
	Use:   "active-response",
	Short: "Execute active response commands on agents",
}

var activeResponseRunCmd = &cobra.Command{
	Use:   "run",
	Short: "Run an active response command on an agent",
	RunE: func(cmd *cobra.Command, args []string) error {
		agentID, _ := cmd.Flags().GetString("agent-id")
		command, _ := cmd.Flags().GetString("command")
		arguments, _ := cmd.Flags().GetStringSlice("arguments")
		custom, _ := cmd.Flags().GetBool("custom")

		body := map[string]interface{}{
			"command": command,
			"custom":  custom,
		}
		if len(arguments) > 0 {
			body["arguments"] = arguments
		}
		if agentID != "" {
			body["agents_list"] = []string{agentID}
		}

		resp, err := globalClient.Put("/active-response", body)
		if err != nil {
			return err
		}
		mustWrite(resp.Items())
		return nil
	},
}

func init() {
	rootCmd.AddCommand(activeResponseCmd)
	activeResponseCmd.AddCommand(activeResponseRunCmd)

	activeResponseRunCmd.Flags().String("agent-id", "", "target agent ID (omit for all agents)")
	activeResponseRunCmd.Flags().String("command", "", "active response command name")
	activeResponseRunCmd.Flags().StringSlice("arguments", nil, "command arguments")
	activeResponseRunCmd.Flags().Bool("custom", false, "command is a custom active response")
	_ = activeResponseRunCmd.MarkFlagRequired("command")
}
