package cmd

import (
	"github.com/ba0f3/wazuh-cli/internal/client"
	"github.com/spf13/cobra"
)

var mitreCmd = &cobra.Command{
	Use:   "mitre",
	Short: "Query the MITRE ATT&CK framework data",
}

var mitreListCmd = &cobra.Command{
	Use:   "list",
	Short: "List MITRE ATT&CK techniques",
	RunE: func(cmd *cobra.Command, args []string) error {
		phase, _ := cmd.Flags().GetString("phase")
		params := &client.QueryParams{Filters: make(map[string]string)}
		if phase != "" {
			params.Filters["phase_name"] = phase
		}
		resp, err := globalClient.Get("/mitre/techniques", params)
		if err != nil {
			return err
		}
		mustWrite(resp.Items())
		return nil
	},
}

var mitreGetCmd = &cobra.Command{
	Use:   "get <ID>",
	Short: "Get a specific MITRE technique by ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		params := &client.QueryParams{Filters: map[string]string{"ids": args[0]}}
		resp, err := globalClient.Get("/mitre/techniques", params)
		if err != nil {
			return err
		}
		mustWrite(resp.Items())
		return nil
	},
}

func init() {
	rootCmd.AddCommand(mitreCmd)
	mitreCmd.AddCommand(mitreListCmd, mitreGetCmd)

	mitreListCmd.Flags().String("phase", "", "filter by kill chain phase name")
}
