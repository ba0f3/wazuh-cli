package cmd

import (
	"github.com/spf13/cobra"
)

var securityCmd = &cobra.Command{
	Use:   "security",
	Short: "Wazuh RBAC — users, roles, policies, rules",
}

// ── Users ────────────────────────────────────────────────────────────────────

var securityUserCmd = &cobra.Command{Use: "user", Short: "Manage security users"}

var securityUserListCmd = &cobra.Command{
	Use: "list", Short: "List all security users",
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := globalClient.Get("/security/users", nil)
		if err != nil {
			return err
		}
		mustWrite(resp.Items())
		return nil
	},
}

var securityUserGetCmd = &cobra.Command{
	Use: "get <ID>", Short: "Get a security user by ID",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := globalClient.Get("/security/users/"+args[0], nil)
		if err != nil {
			return err
		}
		mustWrite(resp.Items())
		return nil
	},
}

var securityUserCreateCmd = &cobra.Command{
	Use: "create", Short: "Create a security user",
	RunE: func(cmd *cobra.Command, args []string) error {
		username, _ := cmd.Flags().GetString("username")
		password, _ := cmd.Flags().GetString("password")
		resp, err := globalClient.Post("/security/users", map[string]string{
			"username": username, "password": password,
		})
		if err != nil {
			return err
		}
		mustWrite(resp.Items())
		return nil
	},
}

var securityUserUpdateCmd = &cobra.Command{
	Use: "update <ID>", Short: "Update a security user",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		password, _ := cmd.Flags().GetString("password")
		resp, err := globalClient.Put("/security/users/"+args[0], map[string]string{"password": password})
		if err != nil {
			return err
		}
		mustWrite(resp.Items())
		return nil
	},
}

var securityUserDeleteCmd = &cobra.Command{
	Use: "delete <ID>", Short: "Delete a security user",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := globalClient.Delete("/security/users?user_ids="+args[0], nil)
		if err != nil {
			return err
		}
		mustWrite(resp.Items())
		return nil
	},
}

// ── Roles ─────────────────────────────────────────────────────────────────────

var securityRoleCmd = &cobra.Command{Use: "role", Short: "Manage security roles"}

var securityRoleListCmd = &cobra.Command{
	Use: "list", Short: "List all security roles",
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := globalClient.Get("/security/roles", nil)
		if err != nil {
			return err
		}
		mustWrite(resp.Items())
		return nil
	},
}

var securityRoleGetCmd = &cobra.Command{
	Use: "get <ID>", Short: "Get a security role by ID",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := globalClient.Get("/security/roles?role_ids="+args[0], nil)
		if err != nil {
			return err
		}
		mustWrite(resp.Items())
		return nil
	},
}

var securityRoleCreateCmd = &cobra.Command{
	Use: "create", Short: "Create a security role",
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		resp, err := globalClient.Post("/security/roles", map[string]string{"name": name})
		if err != nil {
			return err
		}
		mustWrite(resp.Items())
		return nil
	},
}

var securityRoleDeleteCmd = &cobra.Command{
	Use: "delete <ID>", Short: "Delete a security role",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := globalClient.Delete("/security/roles?role_ids="+args[0], nil)
		if err != nil {
			return err
		}
		mustWrite(resp.Items())
		return nil
	},
}

// ── Policies ──────────────────────────────────────────────────────────────────

var securityPolicyCmd = &cobra.Command{Use: "policy", Short: "Manage security policies"}

var securityPolicyListCmd = &cobra.Command{
	Use: "list", Short: "List all security policies",
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := globalClient.Get("/security/policies", nil)
		if err != nil {
			return err
		}
		mustWrite(resp.Items())
		return nil
	},
}

var securityPolicyGetCmd = &cobra.Command{
	Use: "get <ID>", Short: "Get a security policy",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := globalClient.Get("/security/policies?policy_ids="+args[0], nil)
		if err != nil {
			return err
		}
		mustWrite(resp.Items())
		return nil
	},
}

var securityPolicyDeleteCmd = &cobra.Command{
	Use: "delete <ID>", Short: "Delete a security policy",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := globalClient.Delete("/security/policies?policy_ids="+args[0], nil)
		if err != nil {
			return err
		}
		mustWrite(resp.Items())
		return nil
	},
}

// ── Security Rules ─────────────────────────────────────────────────────────────

var securityRuleCmd2 = &cobra.Command{Use: "rule", Short: "Manage security rules"}

var securityRuleListCmd = &cobra.Command{
	Use: "list", Short: "List all security rules",
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := globalClient.Get("/security/rules", nil)
		if err != nil {
			return err
		}
		mustWrite(resp.Items())
		return nil
	},
}

var securityRuleGetCmd = &cobra.Command{
	Use: "get <ID>", Short: "Get a security rule",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := globalClient.Get("/security/rules?rule_ids="+args[0], nil)
		if err != nil {
			return err
		}
		mustWrite(resp.Items())
		return nil
	},
}

func init() {
	rootCmd.AddCommand(securityCmd)

	// Users
	securityCmd.AddCommand(securityUserCmd)
	securityUserCmd.AddCommand(securityUserListCmd, securityUserGetCmd,
		securityUserCreateCmd, securityUserUpdateCmd, securityUserDeleteCmd)
	securityUserCreateCmd.Flags().String("username", "", "username")
	securityUserCreateCmd.Flags().String("password", "", "password")
	_ = securityUserCreateCmd.MarkFlagRequired("username")
	_ = securityUserCreateCmd.MarkFlagRequired("password")
	securityUserUpdateCmd.Flags().String("password", "", "new password")
	_ = securityUserUpdateCmd.MarkFlagRequired("password")

	// Roles
	securityCmd.AddCommand(securityRoleCmd)
	securityRoleCmd.AddCommand(securityRoleListCmd, securityRoleGetCmd,
		securityRoleCreateCmd, securityRoleDeleteCmd)
	securityRoleCreateCmd.Flags().String("name", "", "role name")
	_ = securityRoleCreateCmd.MarkFlagRequired("name")

	// Policies
	securityCmd.AddCommand(securityPolicyCmd)
	securityPolicyCmd.AddCommand(securityPolicyListCmd, securityPolicyGetCmd, securityPolicyDeleteCmd)

	// Security rules
	securityCmd.AddCommand(securityRuleCmd2)
	securityRuleCmd2.AddCommand(securityRuleListCmd, securityRuleGetCmd)
}
