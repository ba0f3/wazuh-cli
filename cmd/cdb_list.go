package cmd

import (
	"github.com/spf13/cobra"
)

var cdbListCmd = &cobra.Command{
	Use:   "cdb-list",
	Short: "Manage CDB lists",
}

var cdbListListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all CDB list files",
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := globalClient.Get("/lists", nil)
		if err != nil {
			return err
		}
		mustWrite(resp.Items())
		return nil
	},
}

var cdbListGetCmd = &cobra.Command{
	Use:   "get <FILENAME>",
	Short: "Get contents of a CDB list file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := globalClient.Get("/lists/files/"+args[0], nil)
		if err != nil {
			return err
		}
		mustWrite(resp.Items())
		return nil
	},
}

var cdbListDeleteCmd = &cobra.Command{
	Use:   "delete <FILENAME>",
	Short: "Delete a CDB list file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := globalClient.Delete("/lists/files/"+args[0], nil)
		if err != nil {
			return err
		}
		mustWrite(resp.Items())
		return nil
	},
}

func init() {
	rootCmd.AddCommand(cdbListCmd)
	cdbListCmd.AddCommand(cdbListListCmd, cdbListGetCmd, cdbListDeleteCmd)
}
