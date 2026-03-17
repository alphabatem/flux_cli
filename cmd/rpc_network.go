package cmd

import (
	"os"
	"strings"

	"github.com/alphabatem/flux_cli/dto"
	"github.com/alphabatem/flux_cli/internal/httpclient"
	"github.com/alphabatem/flux_cli/output"
	"github.com/spf13/cobra"
)

func init() {
	rpcCmd.AddCommand(rpcNetworkCmd)
	rpcNetworkCmd.AddCommand(rpcNetworkHealthCmd)
	rpcNetworkCmd.AddCommand(rpcNetworkVersionCmd)
	rpcNetworkCmd.AddCommand(rpcNetworkIdentityCmd)
	rpcNetworkCmd.AddCommand(rpcNetworkNodesCmd)
	rpcNetworkCmd.AddCommand(rpcNetworkGenesisCmd)
	rpcNetworkCmd.AddCommand(rpcNetworkFirstBlockCmd)

	rpcNetworkPerformanceCmd.Flags().Int("limit", 720, "Number of samples (max 720)")
	rpcNetworkCmd.AddCommand(rpcNetworkPerformanceCmd)

	rpcNetworkCmd.AddCommand(rpcNetworkPriorityFeesCmd)

	rpcNetworkSupplyCmd.Flags().String("commitment", "", "Commitment level")
	rpcNetworkCmd.AddCommand(rpcNetworkSupplyCmd)
}

var rpcNetworkCmd = &cobra.Command{
	Use:   "network",
	Short: "Network/system methods - health, version, nodes, genesis, performance, supply",
}

var rpcNetworkHealthCmd = &cobra.Command{
	Use:   "health",
	Short: "Check RPC node health (getHealth)",
	Run: func(cmd *cobra.Command, args []string) {
		result, err := fluxRPCSvc().GetHealth()
		if err != nil {
			output.PrintError(cmd, "RPC_ERROR", err.Error(), &dto.CLIMeta{Service: "fluxrpc"})
			os.Exit(httpclient.ExitCodeForError(err))
		}
		output.PrintSuccess(cmd, result, &dto.CLIMeta{Service: "fluxrpc", Endpoint: "getHealth"})
	},
}

var rpcNetworkVersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Get node software version (getVersion)",
	Run: func(cmd *cobra.Command, args []string) {
		result, err := fluxRPCSvc().GetVersion()
		if err != nil {
			output.PrintError(cmd, "RPC_ERROR", err.Error(), &dto.CLIMeta{Service: "fluxrpc"})
			os.Exit(httpclient.ExitCodeForError(err))
		}
		output.PrintSuccess(cmd, result, &dto.CLIMeta{Service: "fluxrpc", Endpoint: "getVersion"})
	},
}

var rpcNetworkIdentityCmd = &cobra.Command{
	Use:   "identity",
	Short: "Get RPC node identity (getIdentity)",
	Run: func(cmd *cobra.Command, args []string) {
		result, err := fluxRPCSvc().GetIdentity()
		if err != nil {
			output.PrintError(cmd, "RPC_ERROR", err.Error(), &dto.CLIMeta{Service: "fluxrpc"})
			os.Exit(httpclient.ExitCodeForError(err))
		}
		output.PrintSuccess(cmd, result, &dto.CLIMeta{Service: "fluxrpc", Endpoint: "getIdentity"})
	},
}

var rpcNetworkNodesCmd = &cobra.Command{
	Use:   "nodes",
	Short: "Get cluster nodes (getClusterNodes)",
	Run: func(cmd *cobra.Command, args []string) {
		result, err := fluxRPCSvc().GetClusterNodes()
		if err != nil {
			output.PrintError(cmd, "RPC_ERROR", err.Error(), &dto.CLIMeta{Service: "fluxrpc"})
			os.Exit(httpclient.ExitCodeForError(err))
		}
		output.PrintSuccess(cmd, result, &dto.CLIMeta{Service: "fluxrpc", Endpoint: "getClusterNodes"})
	},
}

var rpcNetworkGenesisCmd = &cobra.Command{
	Use:   "genesis-hash",
	Short: "Get genesis hash (getGenesisHash)",
	Run: func(cmd *cobra.Command, args []string) {
		result, err := fluxRPCSvc().GetGenesisHash()
		if err != nil {
			output.PrintError(cmd, "RPC_ERROR", err.Error(), &dto.CLIMeta{Service: "fluxrpc"})
			os.Exit(httpclient.ExitCodeForError(err))
		}
		output.PrintSuccess(cmd, result, &dto.CLIMeta{Service: "fluxrpc", Endpoint: "getGenesisHash"})
	},
}

var rpcNetworkFirstBlockCmd = &cobra.Command{
	Use:   "first-available-block",
	Short: "Get first available block (getFirstAvailableBlock)",
	Run: func(cmd *cobra.Command, args []string) {
		result, err := fluxRPCSvc().GetFirstAvailableBlock()
		if err != nil {
			output.PrintError(cmd, "RPC_ERROR", err.Error(), &dto.CLIMeta{Service: "fluxrpc"})
			os.Exit(httpclient.ExitCodeForError(err))
		}
		output.PrintSuccess(cmd, result, &dto.CLIMeta{Service: "fluxrpc", Endpoint: "getFirstAvailableBlock"})
	},
}

var rpcNetworkPerformanceCmd = &cobra.Command{
	Use:   "performance",
	Short: "Get recent performance samples (getRecentPerformanceSamples)",
	Run: func(cmd *cobra.Command, args []string) {
		limit, _ := cmd.Flags().GetInt("limit")
		result, err := fluxRPCSvc().GetRecentPerformanceSamples(limit)
		if err != nil {
			output.PrintError(cmd, "RPC_ERROR", err.Error(), &dto.CLIMeta{Service: "fluxrpc"})
			os.Exit(httpclient.ExitCodeForError(err))
		}
		output.PrintSuccess(cmd, result, &dto.CLIMeta{Service: "fluxrpc", Endpoint: "getRecentPerformanceSamples"})
	},
}

var rpcNetworkPriorityFeesCmd = &cobra.Command{
	Use:   "priority-fees [addr1,addr2,...]",
	Short: "Get recent prioritization fees (getRecentPrioritizationFees)",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var addresses []string
		if len(args) > 0 {
			addresses = strings.Split(args[0], ",")
		}
		result, err := fluxRPCSvc().GetRecentPrioritizationFees(addresses)
		if err != nil {
			output.PrintError(cmd, "RPC_ERROR", err.Error(), &dto.CLIMeta{Service: "fluxrpc"})
			os.Exit(httpclient.ExitCodeForError(err))
		}
		output.PrintSuccess(cmd, result, &dto.CLIMeta{Service: "fluxrpc", Endpoint: "getRecentPrioritizationFees"})
	},
}

var rpcNetworkSupplyCmd = &cobra.Command{
	Use:   "supply",
	Short: "Get total, circulating, and non-circulating supply (getSupply)",
	Run: func(cmd *cobra.Command, args []string) {
		commitment, _ := cmd.Flags().GetString("commitment")
		result, err := fluxRPCSvc().GetSupply(commitment)
		if err != nil {
			output.PrintError(cmd, "RPC_ERROR", err.Error(), &dto.CLIMeta{Service: "fluxrpc"})
			os.Exit(httpclient.ExitCodeForError(err))
		}
		output.PrintSuccess(cmd, result, &dto.CLIMeta{Service: "fluxrpc", Endpoint: "getSupply"})
	},
}
