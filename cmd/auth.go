package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/ba0f3/wazuh-cli/internal/client"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authentication and token management",
}

var authLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate and cache a new JWT token",
	Long: `Prompt for credentials and exchange them for a JWT token.
This is safer than setting the password in the config file.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		url, _ := cmd.Flags().GetString("url")
		user, _ := cmd.Flags().GetString("user")

		if url == "" && globalCfg != nil {
			url = globalCfg.URL
		}
		if user == "" && globalCfg != nil {
			user = globalCfg.User
		}

		if url == "" {
			fmt.Print("Wazuh API URL: ")
			fmt.Scanln(&url)
		}
		if user == "" {
			fmt.Print("Username: ")
			fmt.Scanln(&user)
		}

		fmt.Print("Password: ")
		bytePassword, err := term.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			return err
		}
		fmt.Println() // New line after password entry
		password := strings.TrimSpace(string(bytePassword))

		// Set overrides and re-init client
		// We copy current flagOverrides to ensure --insecure etc are preserved
		loginCfg := flagOverrides
		loginCfg.URL = url
		loginCfg.User = user
		loginCfg.Password = password

		// Re-initialize the client with the new credentials
		c, err := client.NewClient(&loginCfg)
		if err != nil {
			return err
		}

		// Manually trigger authentication
		_, err = c.Auth().Token()
		if err != nil {
			return err
		}

		fmt.Fprintf(os.Stderr, "✓ Authenticated successfully as %s\n", user)
		fmt.Fprintf(os.Stderr, "  Token cached at %s\n", c.Auth().TokenPath())

		// Optionally save URL/User to config (but NOT password)
		save, _ := cmd.Flags().GetBool("save")
		if save {
			flagOverrides.Password = "" // Don't save password
			if err := flagOverrides.Save(""); err != nil {
				return err
			}
			fmt.Fprintf(os.Stderr, "✓ Saved URL and User to config file\n")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(authCmd)
	authCmd.AddCommand(authLoginCmd)

	authLoginCmd.Flags().String("url", "", "Wazuh API URL")
	authLoginCmd.Flags().String("user", "", "API username")
	authLoginCmd.Flags().Bool("save", false, "save URL and User to config file (won't save password)")
}
