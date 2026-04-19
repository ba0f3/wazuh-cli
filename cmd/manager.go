package cmd

import (
	"github.com/ba0f3/wazuh-cli/internal/client"
	"github.com/spf13/cobra"
)

var managerCmd = &cobra.Command{
	Use:   "manager",
	Short: "Wazuh manager (server) information and control",
}

var managerInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Get manager information",
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := globalClient.Get("/manager/info", nil)
		if err != nil {
			return err
		}
		mustWrite(resp.Data)
		return nil
	},
}

var managerStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Get manager daemon status",
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := globalClient.Get("/manager/status", nil)
		if err != nil {
			return err
		}
		mustWrite(resp.Data)
		return nil
	},
}

var managerConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Get manager active configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		component, _ := cmd.Flags().GetString("component")
		configuration, _ := cmd.Flags().GetString("configuration")
		var path string
		if component != "" && configuration != "" {
			path = "/manager/configuration/" + component + "/" + configuration
		} else {
			path = "/manager/configuration"
		}
		resp, err := globalClient.Get(path, nil)
		if err != nil {
			return err
		}
		mustWrite(resp.Data)
		return nil
	},
}

var managerStatsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Get manager statistics",
	RunE: func(cmd *cobra.Command, args []string) error {
		date, _ := cmd.Flags().GetString("date")
		params := &client.QueryParams{Filters: make(map[string]string)}
		if date != "" {
			params.Filters["date"] = date
		}
		resp, err := globalClient.Get("/manager/stats", params)
		if err != nil {
			return err
		}
		mustWrite(resp.Items())
		return nil
	},
}

var managerLogsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Get manager logs",
	RunE: func(cmd *cobra.Command, args []string) error {
		level, _ := cmd.Flags().GetString("level")
		tag, _ := cmd.Flags().GetString("tag")
		params := &client.QueryParams{Filters: make(map[string]string)}
		if level != "" {
			params.Filters["level"] = level
		}
		if tag != "" {
			params.Filters["tag"] = tag
		}
		resp, err := globalClient.Get("/manager/logs", params)
		if err != nil {
			return err
		}
		mustWrite(resp.Items())
		return nil
	},
}

var managerLogSummaryCmd = &cobra.Command{
	Use:   "log-summary",
	Short: "Get a summary of manager logs",
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := globalClient.Get("/manager/logs/summary", nil)
		if err != nil {
			return err
		}
		mustWrite(resp.Data)
		return nil
	},
}

var managerRestartCmd = &cobra.Command{
	Use:   "restart",
	Short: "Restart the manager",
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := globalClient.Put("/manager/restart", nil)
		if err != nil {
			return err
		}
		mustWrite(resp.Data)
		return nil
	},
}

var managerValidationCmd = &cobra.Command{
	Use:   "validation",
	Short: "Validate manager configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := globalClient.Get("/manager/configuration/validation", nil)
		if err != nil {
			return err
		}
		mustWrite(resp.Data)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(managerCmd)
	managerCmd.AddCommand(managerInfoCmd, managerStatusCmd, managerConfigCmd,
		managerStatsCmd, managerLogsCmd, managerLogSummaryCmd,
		managerRestartCmd, managerValidationCmd)

	managerConfigCmd.Flags().String("component", "", "configuration component")
	managerConfigCmd.Flags().String("configuration", "", "configuration name")

	managerStatsCmd.Flags().String("date", "", "stats date (YYYYMMDD)")

	managerLogsCmd.Flags().String("level", "", "log level: error, warning, info")
	managerLogsCmd.Flags().String("tag", "", "log tag/module")
}
