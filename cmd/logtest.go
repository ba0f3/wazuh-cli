package cmd

import (
	"github.com/spf13/cobra"
)

var logtestCmd = &cobra.Command{
	Use:   "logtest",
	Short: "Test log decoding and rule matching",
}

var logtestRunCmd = &cobra.Command{
	Use:   "run",
	Short: "Send an event through the log test engine",
	RunE: func(cmd *cobra.Command, args []string) error {
		event, _ := cmd.Flags().GetString("event")
		sessionID, _ := cmd.Flags().GetString("session-id")

		body := map[string]interface{}{
			"event":      event,
			"log_format": "syslog",
			"location":   "logtest",
		}
		if sessionID != "" {
			body["token"] = sessionID
		}

		resp, err := globalClient.Put("/logtest", body)
		if err != nil {
			return err
		}
		mustWrite(resp.Data)
		return nil
	},
}

var logtestEndCmd = &cobra.Command{
	Use:   "end",
	Short: "End a logtest session",
	RunE: func(cmd *cobra.Command, args []string) error {
		sessionID, _ := cmd.Flags().GetString("session-id")
		resp, err := globalClient.Delete("/logtest/sessions/"+sessionID, nil)
		if err != nil {
			return err
		}
		mustWrite(resp.Data)
		return nil
	},
}

var ciscatCmd = &cobra.Command{
	Use:   "ciscat",
	Short: "CIS-CAT assessment results",
}

var ciscatResultsCmd = &cobra.Command{
	Use:   "results <AGENT_ID>",
	Short: "Get CIS-CAT results for an agent",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := globalClient.Get("/ciscat/"+args[0]+"/results", nil)
		if err != nil {
			return err
		}
		mustWrite(resp.Items())
		return nil
	},
}

func init() {
	rootCmd.AddCommand(logtestCmd)
	logtestCmd.AddCommand(logtestRunCmd, logtestEndCmd)

	logtestRunCmd.Flags().String("event", "", "log event string to test")
	logtestRunCmd.Flags().String("session-id", "", "existing logtest session token")
	_ = logtestRunCmd.MarkFlagRequired("event")

	logtestEndCmd.Flags().String("session-id", "", "logtest session token to close")
	_ = logtestEndCmd.MarkFlagRequired("session-id")

	rootCmd.AddCommand(ciscatCmd)
	ciscatCmd.AddCommand(ciscatResultsCmd)
}
