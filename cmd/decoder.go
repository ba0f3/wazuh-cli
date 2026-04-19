package cmd

import (
	"github.com/ba0f3/wazuh-cli/internal/client"
	"github.com/spf13/cobra"
)

var decoderCmd = &cobra.Command{
	Use:   "decoder",
	Short: "Manage Wazuh decoders",
}

var decoderListCmd = &cobra.Command{
	Use:   "list",
	Short: "List decoders",
	RunE: func(cmd *cobra.Command, args []string) error {
		status, _ := cmd.Flags().GetString("status")
		file, _ := cmd.Flags().GetString("file")
		limit, _ := cmd.Flags().GetInt("limit")
		offset, _ := cmd.Flags().GetInt("offset")

		params := &client.QueryParams{Limit: limit, Offset: offset, Filters: make(map[string]string)}
		if status != "" {
			params.Filters["status"] = status
		}
		if file != "" {
			params.Filters["filename"] = file
		}

		resp, err := globalClient.Get("/decoders", params)
		if err != nil {
			return err
		}
		mustWrite(resp.Items())
		return nil
	},
}

var decoderGetCmd = &cobra.Command{
	Use:   "get <NAME>",
	Short: "Get a specific decoder by name",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		params := &client.QueryParams{Filters: map[string]string{"decoder_names": args[0]}}
		resp, err := globalClient.Get("/decoders", params)
		if err != nil {
			return err
		}
		mustWrite(resp.Items())
		return nil
	},
}

var decoderFilesCmd = &cobra.Command{
	Use:   "files",
	Short: "List decoder files",
	RunE: func(cmd *cobra.Command, args []string) error {
		status, _ := cmd.Flags().GetString("status")
		params := &client.QueryParams{}
		if status != "" {
			params.Filters = map[string]string{"status": status}
		}
		resp, err := globalClient.Get("/decoders/files", params)
		if err != nil {
			return err
		}
		mustWrite(resp.Items())
		return nil
	},
}

func init() {
	rootCmd.AddCommand(decoderCmd)
	decoderCmd.AddCommand(decoderListCmd, decoderGetCmd, decoderFilesCmd)

	decoderListCmd.Flags().String("status", "", "filter by status: enabled, disabled")
	decoderListCmd.Flags().String("file", "", "filter by file name")
	decoderListCmd.Flags().Int("limit", 500, "maximum results")
	decoderListCmd.Flags().Int("offset", 0, "pagination offset")

	decoderFilesCmd.Flags().String("status", "", "filter by status: enabled, disabled")
}
