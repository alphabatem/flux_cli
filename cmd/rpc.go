package cmd

import (
	"encoding/json"
	"os"

	"github.com/alphabatem/flux_cli/dto"
	"github.com/alphabatem/flux_cli/internal/httpclient"
	"github.com/alphabatem/flux_cli/output"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(rpcCmd)
	rpcCmd.AddCommand(rpcCallCmd)
}

var rpcCmd = &cobra.Command{
	Use:   "rpc",
	Short: "FluxRPC - Solana JSON-RPC access (accounts, blocks, transactions, tokens, slots, epochs, network, staking)",
}

var rpcCallCmd = &cobra.Command{
	Use:   "call <method> [params-json]",
	Short: "Execute an arbitrary JSON-RPC method",
	Long: `Execute any Solana JSON-RPC method directly.

Examples:
  flux rpc call getHealth
  flux rpc call getBalance '["pubkey123"]'
  flux rpc call getBlock '[12345, {"encoding": "json"}]'`,
	Args: cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		method := args[0]
		var params []interface{}
		if len(args) > 1 {
			if err := json.Unmarshal([]byte(args[1]), &params); err != nil {
				output.PrintError(cmd, "INVALID_PARAMS", "params must be a valid JSON array: "+err.Error(), nil)
				os.Exit(dto.ExitUsageError)
			}
		}
		result, err := fluxRPCSvc().Call(method, params)
		if err != nil {
			output.PrintError(cmd, "RPC_ERROR", err.Error(), &dto.CLIMeta{Service: "fluxrpc", Endpoint: method})
			os.Exit(httpclient.ExitCodeForError(err))
		}
		output.PrintSuccess(cmd, result, &dto.CLIMeta{Service: "fluxrpc", Endpoint: method})
	},
}
