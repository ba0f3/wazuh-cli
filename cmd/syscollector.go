package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var syscollectorCmd = &cobra.Command{
	Use:   "syscollector",
	Short: "System inventory (Syscollector)",
}

func syscollectorGet(endpoint string) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		resp, err := globalClient.Get(fmt.Sprintf("/syscollector/%s/%s", args[0], endpoint), nil)
		if err != nil {
			return err
		}
		mustWrite(resp.Items())
		return nil
	}
}

var syscollectorHardwareCmd = &cobra.Command{Use: "hardware <AGENT_ID>", Short: "Hardware information", Args: cobra.ExactArgs(1), RunE: syscollectorGet("hardware")}
var syscollectorOsCmd = &cobra.Command{Use: "os <AGENT_ID>", Short: "Operating system information", Args: cobra.ExactArgs(1), RunE: syscollectorGet("os")}
var syscollectorPackagesCmd = &cobra.Command{Use: "packages <AGENT_ID>", Short: "Installed packages", Args: cobra.ExactArgs(1), RunE: syscollectorGet("packages")}
var syscollectorProcessesCmd = &cobra.Command{Use: "processes <AGENT_ID>", Short: "Running processes", Args: cobra.ExactArgs(1), RunE: syscollectorGet("processes")}
var syscollectorPortsCmd = &cobra.Command{Use: "ports <AGENT_ID>", Short: "Open ports", Args: cobra.ExactArgs(1), RunE: syscollectorGet("ports")}
var syscollectorNetaddrCmd = &cobra.Command{Use: "netaddr <AGENT_ID>", Short: "Network addresses", Args: cobra.ExactArgs(1), RunE: syscollectorGet("netaddr")}
var syscollectorNetifaceCmd = &cobra.Command{Use: "netiface <AGENT_ID>", Short: "Network interfaces", Args: cobra.ExactArgs(1), RunE: syscollectorGet("netiface")}
var syscollectorHotfixesCmd = &cobra.Command{Use: "hotfixes <AGENT_ID>", Short: "Windows hotfixes", Args: cobra.ExactArgs(1), RunE: syscollectorGet("hotfixes")}

func init() {
	rootCmd.AddCommand(syscollectorCmd)
	syscollectorCmd.AddCommand(
		syscollectorHardwareCmd, syscollectorOsCmd, syscollectorPackagesCmd,
		syscollectorProcessesCmd, syscollectorPortsCmd, syscollectorNetaddrCmd,
		syscollectorNetifaceCmd, syscollectorHotfixesCmd,
	)
}
