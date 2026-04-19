package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of wazuh-cli",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("wazuh-cli version %s\n", Version)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
