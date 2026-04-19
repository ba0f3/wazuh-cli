package cmd

import (
	"fmt"

	"github.com/ba0f3/wazuh-cli/internal/client"
	"github.com/spf13/cobra"
)

var groupCmd = &cobra.Command{
	Use:   "group",
	Short: "Manage agent groups",
}

var groupListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all agent groups",
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := globalClient.Get("/groups", nil)
		if err != nil {
			return err
		}
		mustWrite(resp.Items())
		return nil
	},
}

var groupGetCmd = &cobra.Command{
	Use:   "get <NAME>",
	Short: "Get a specific group",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		params := &client.QueryParams{Filters: map[string]string{"groups_list": args[0]}}
		resp, err := globalClient.Get("/groups", params)
		if err != nil {
			return err
		}
		mustWrite(resp.Items())
		return nil
	},
}

var groupCreateCmd = &cobra.Command{
	Use:   "create <NAME>",
	Short: "Create a new agent group",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := globalClient.Post("/groups", map[string]string{"group_id": args[0]})
		if err != nil {
			return err
		}
		mustWrite(resp.Data)
		return nil
	},
}

var groupDeleteCmd = &cobra.Command{
	Use:   "delete <NAME>",
	Short: "Delete an agent group",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		params := &client.QueryParams{Filters: map[string]string{"groups_list": args[0]}}
		resp, err := globalClient.Delete("/groups", params)
		if err != nil {
			return err
		}
		mustWrite(resp.Items())
		return nil
	},
}

var groupAgentsCmd = &cobra.Command{
	Use:   "agents <NAME>",
	Short: "List agents in a group",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := globalClient.Get(fmt.Sprintf("/groups/%s/agents", args[0]), nil)
		if err != nil {
			return err
		}
		mustWrite(resp.Items())
		return nil
	},
}

var groupAddAgentCmd = &cobra.Command{
	Use:   "add-agent <GROUP> <AGENT_ID>",
	Short: "Add an agent to a group",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		params := &client.QueryParams{Filters: map[string]string{"group_id": args[0], "agent_id": args[1]}}
		resp, err := globalClient.Put("/agents/group", params)
		if err != nil {
			return err
		}
		mustWrite(resp.Items())
		return nil
	},
}

var groupRemoveAgentCmd = &cobra.Command{
	Use:   "remove-agent <GROUP> <AGENT_ID>",
	Short: "Remove an agent from a group",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		params := &client.QueryParams{Filters: map[string]string{"group_id": args[0], "agents_list": args[1]}}
		resp, err := globalClient.Delete("/agents/group", params)
		if err != nil {
			return err
		}
		mustWrite(resp.Items())
		return nil
	},
}

var groupConfigCmd = &cobra.Command{
	Use:   "config <NAME>",
	Short: "Get a group's configuration",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := globalClient.Get(fmt.Sprintf("/groups/%s/configuration", args[0]), nil)
		if err != nil {
			return err
		}
		mustWrite(resp.Items())
		return nil
	},
}

func init() {
	rootCmd.AddCommand(groupCmd)
	groupCmd.AddCommand(groupListCmd, groupGetCmd, groupCreateCmd, groupDeleteCmd,
		groupAgentsCmd, groupAddAgentCmd, groupRemoveAgentCmd, groupConfigCmd)
}
