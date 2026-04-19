package cmd

import (
	"github.com/ba0f3/wazuh-cli/internal/client"
	"github.com/spf13/cobra"
)

var clusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "Wazuh cluster management",
}

var clusterStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Get cluster status",
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := globalClient.Get("/cluster/status", nil)
		if err != nil {
			return err
		}
		mustWrite(resp.Data)
		return nil
	},
}

var clusterNodesCmd = &cobra.Command{
	Use:   "nodes",
	Short: "List cluster nodes",
	RunE: func(cmd *cobra.Command, args []string) error {
		nodeType, _ := cmd.Flags().GetString("type")
		params := &client.QueryParams{Filters: make(map[string]string)}
		if nodeType != "" {
			params.Filters["type"] = nodeType
		}
		resp, err := globalClient.Get("/cluster/nodes", params)
		if err != nil {
			return err
		}
		mustWrite(resp.Items())
		return nil
	},
}

var clusterNodeCmd = &cobra.Command{
	Use:   "node <NAME>",
	Short: "Get information about a specific cluster node",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := globalClient.Get("/cluster/"+args[0]+"/info", nil)
		if err != nil {
			return err
		}
		mustWrite(resp.Data)
		return nil
	},
}

var clusterHealthCmd = &cobra.Command{
	Use:   "health",
	Short: "Get cluster health check",
	RunE: func(cmd *cobra.Command, args []string) error {
		node, _ := cmd.Flags().GetString("node")
		params := &client.QueryParams{Filters: make(map[string]string)}
		if node != "" {
			params.Filters["nodes_list"] = node
		}
		resp, err := globalClient.Get("/cluster/healthcheck", params)
		if err != nil {
			return err
		}
		mustWrite(resp.Data)
		return nil
	},
}

var clusterConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Get cluster configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := globalClient.Get("/cluster/local/config", nil)
		if err != nil {
			return err
		}
		mustWrite(resp.Data)
		return nil
	},
}

var clusterRestartCmd = &cobra.Command{
	Use:   "restart",
	Short: "Restart cluster nodes",
	RunE: func(cmd *cobra.Command, args []string) error {
		node, _ := cmd.Flags().GetString("node")
		params := &client.QueryParams{Filters: make(map[string]string)}
		if node != "" {
			params.Filters["nodes_list"] = node
		}
		resp, err := globalClient.Put("/cluster/restart", params)
		if err != nil {
			return err
		}
		mustWrite(resp.Items())
		return nil
	},
}

func init() {
	rootCmd.AddCommand(clusterCmd)
	clusterCmd.AddCommand(clusterStatusCmd, clusterNodesCmd, clusterNodeCmd,
		clusterHealthCmd, clusterConfigCmd, clusterRestartCmd)

	clusterNodesCmd.Flags().String("type", "", "filter by node type: master, worker")
	clusterHealthCmd.Flags().String("node", "", "specific node name")
	clusterRestartCmd.Flags().String("node", "", "specific node name (omit for all)")
}
