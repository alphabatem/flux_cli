package cmd

import (
	"os"
	"strconv"

	"github.com/alphabatem/flux_cli/dto"
	"github.com/alphabatem/flux_cli/internal/httpclient"
	"github.com/alphabatem/flux_cli/output"
	"github.com/spf13/cobra"
)

func init() {
	rpcCmd.AddCommand(rpcSendCmd)

	rpcSendTxCmd.Flags().Bool("skip-preflight", false, "Skip preflight transaction checks")
	rpcSendTxCmd.Flags().String("encoding", "base64", "Transaction encoding: base58, base64")
	rpcSendTxCmd.Flags().String("commitment", "", "Preflight commitment level")
	rpcSendCmd.AddCommand(rpcSendTxCmd)

	rpcSendSimulateCmd.Flags().Bool("sig-verify", false, "Verify transaction signatures")
	rpcSendSimulateCmd.Flags().String("encoding", "base64", "Transaction encoding: base58, base64")
	rpcSendSimulateCmd.Flags().String("commitment", "", "Commitment level")
	rpcSendCmd.AddCommand(rpcSendSimulateCmd)

	rpcSendAirdropCmd.Flags().String("commitment", "", "Commitment level")
	rpcSendCmd.AddCommand(rpcSendAirdropCmd)
}

var rpcSendCmd = &cobra.Command{
	Use:   "send",
	Short: "Transaction submission - send, simulate, airdrop",
}

var rpcSendTxCmd = &cobra.Command{
	Use:   "transaction <base64Tx>",
	Short: "Submit a signed transaction (sendTransaction)",
	Args:  exactArgsFromUse(),
	Run: func(cmd *cobra.Command, args []string) {
		encoding, _ := cmd.Flags().GetString("encoding")
		skipPreflight, _ := cmd.Flags().GetBool("skip-preflight")
		commitment, _ := cmd.Flags().GetString("commitment")
		result, err := fluxRPCSvc().SendTransaction(args[0], encoding, skipPreflight, commitment)
		if err != nil {
			output.PrintError(cmd, "RPC_ERROR", err.Error(), &dto.CLIMeta{Service: "fluxrpc", Endpoint: "sendTransaction"})
			os.Exit(httpclient.ExitCodeForError(err))
		}
		output.PrintSuccess(cmd, result, &dto.CLIMeta{Service: "fluxrpc", Endpoint: "sendTransaction"})
	},
}

var rpcSendSimulateCmd = &cobra.Command{
	Use:   "simulate <base64Tx>",
	Short: "Simulate a transaction without broadcasting (simulateTransaction)",
	Args:  exactArgsFromUse(),
	Run: func(cmd *cobra.Command, args []string) {
		encoding, _ := cmd.Flags().GetString("encoding")
		sigVerify, _ := cmd.Flags().GetBool("sig-verify")
		commitment, _ := cmd.Flags().GetString("commitment")
		result, err := fluxRPCSvc().SimulateTransaction(args[0], encoding, sigVerify, commitment)
		if err != nil {
			output.PrintError(cmd, "RPC_ERROR", err.Error(), &dto.CLIMeta{Service: "fluxrpc", Endpoint: "simulateTransaction"})
			os.Exit(httpclient.ExitCodeForError(err))
		}
		output.PrintSuccess(cmd, result, &dto.CLIMeta{Service: "fluxrpc", Endpoint: "simulateTransaction"})
	},
}

var rpcSendAirdropCmd = &cobra.Command{
	Use:   "airdrop <pubkey> <lamports>",
	Short: "Request an airdrop (requestAirdrop)",
	Args:  exactArgsFromUse(),
	Run: func(cmd *cobra.Command, args []string) {
		lamports, err := strconv.ParseUint(args[1], 10, 64)
		if err != nil {
			output.PrintError(cmd, "INVALID_ARGS", "lamports must be a number", nil)
			os.Exit(dto.ExitUsageError)
		}
		commitment, _ := cmd.Flags().GetString("commitment")
		result, rpcErr := fluxRPCSvc().RequestAirdrop(args[0], lamports, commitment)
		if rpcErr != nil {
			output.PrintError(cmd, "RPC_ERROR", rpcErr.Error(), &dto.CLIMeta{Service: "fluxrpc", Endpoint: "requestAirdrop"})
			os.Exit(httpclient.ExitCodeForError(rpcErr))
		}
		output.PrintSuccess(cmd, result, &dto.CLIMeta{Service: "fluxrpc", Endpoint: "requestAirdrop"})
	},
}
