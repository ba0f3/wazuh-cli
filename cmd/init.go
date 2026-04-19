package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/ba0f3/wazuh-cli/internal/config"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize wazuh-cli configuration",
	Long: `Create the wazuh-cli configuration file at ~/.config/wazuh/config.json.

Provide all values via flags for non-interactive use (suitable for scripts and AI agents),
or run without flags for an interactive setup wizard.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		url, _ := cmd.Flags().GetString("url")
		user, _ := cmd.Flags().GetString("user")
		password, _ := cmd.Flags().GetString("password")
		insecure, _ := cmd.Flags().GetBool("insecure")
		output, _ := cmd.Flags().GetString("output")
		configOut, _ := cmd.Flags().GetString("config-path")
		indexerURL, _ := cmd.Flags().GetString("indexer-url")
		indexerUser, _ := cmd.Flags().GetString("indexer-user")
		indexerPassword, _ := cmd.Flags().GetString("indexer-password")
		indexerIndex, _ := cmd.Flags().GetString("indexer-index")

		if configOut == "" {
			configOut = config.DefaultConfigPath()
		}

		// Interactive prompts for any missing values
		reader := bufio.NewReader(os.Stdin)

		if url == "" {
			fmt.Fprint(os.Stderr, "Wazuh API URL (e.g. https://wazuh:55000): ")
			url = strings.TrimSpace(readLine(reader))
		}
		if user == "" {
			fmt.Fprint(os.Stderr, "Username: ")
			user = strings.TrimSpace(readLine(reader))
		}
		if password == "" {
			fmt.Fprint(os.Stderr, "Password: ")
			password = strings.TrimSpace(readLine(reader))
		}

		if url == "" || user == "" || password == "" {
			return fmt.Errorf("url, user, and password are required")
		}

		cfg := &config.Config{
			URL:             url,
			User:            user,
			Password:        password,
			Insecure:        insecure,
			Timeout:         config.DefaultTimeout,
			Output:          output,
			Pretty:          true,
			IndexerURL:      indexerURL,
			IndexerUser:     indexerUser,
			IndexerPassword: indexerPassword,
			IndexerIndex:    indexerIndex,
		}
		if cfg.Output == "" {
			cfg.Output = config.DefaultOutput
		}

		if err := cfg.Save(configOut); err != nil {
			return fmt.Errorf("saving config: %w", err)
		}

		fmt.Fprintf(os.Stderr, "✓ Configuration written to %s\n", configOut)
		fmt.Fprintf(os.Stderr, "  Run 'wazuh-cli manager info' to verify connectivity.\n")
		return nil
	},
}

func readLine(r *bufio.Reader) string {
	line, _ := r.ReadString('\n')
	return strings.TrimRight(line, "\r\n")
}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().String("url", "", "Wazuh API URL")
	initCmd.Flags().String("user", "", "API username")
	initCmd.Flags().String("password", "", "API password")
	initCmd.Flags().Bool("insecure", false, "skip TLS verification")
	initCmd.Flags().String("output", "json", "default output format: json, markdown, raw")
	initCmd.Flags().String("config-path", "", "output config file path (default: ~/.config/wazuh/config.json)")
	initCmd.Flags().String("indexer-url", "", "Wazuh Indexer URL (e.g. https://indexer:9200)")
	initCmd.Flags().String("indexer-user", "", "Indexer username")
	initCmd.Flags().String("indexer-password", "", "Indexer password")
	initCmd.Flags().String("indexer-index", "", "Indexer index pattern (default: wazuh-alerts-4.x-*)")
}
