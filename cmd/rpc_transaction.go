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
	rpcCmd.AddCommand(rpcTransactionCmd)

	rpcTxShowCmd.Flags().String("encoding", "json", "Encoding: json, jsonParsed, base64, base58")
	rpcTxShowCmd.Flags().String("commitment", "", "Commitment level: confirmed, finalized")
	rpcTransactionCmd.AddCommand(rpcTxShowCmd)

	rpcTxCountCmd.Flags().String("commitment", "", "Commitment level")
	rpcTransactionCmd.AddCommand(rpcTxCountCmd)

	rpcTxSignaturesCmd.Flags().Int("limit", 1000, "Max signatures to return (1-1000)")
	rpcTxSignaturesCmd.Flags().String("before", "", "Start searching backwards from this signature")
	rpcTxSignaturesCmd.Flags().String("until", "", "Search until this signature")
	rpcTxSignaturesCmd.Flags().String("commitment", "", "Commitment level: confirmed, finalized")
	rpcTransactionCmd.AddCommand(rpcTxSignaturesCmd)

	rpcTxStatusCmd.Flags().Bool("search-history", false, "Search full transaction history")
	rpcTransactionCmd.AddCommand(rpcTxStatusCmd)

	rpcTxFeeCmd.Flags().String("commitment", "", "Commitment level")
	rpcTransactionCmd.AddCommand(rpcTxFeeCmd)
}

var rpcTransactionCmd = &cobra.Command{
	Use:   "transaction",
	Short: "Transaction methods - show, count, signatures, status, fee",
}

var rpcTxShowCmd = &cobra.Command{
	Use:   "show <signature>",
	Short: "Get transaction details (getTransaction)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		encoding, _ := cmd.Flags().GetString("encoding")
		commitment, _ := cmd.Flags().GetString("commitment")
		result, err := fluxRPCSvc().GetTransaction(args[0], encoding, commitment)
		if err != nil {
			output.PrintError(cmd, "RPC_ERROR", err.Error(), &dto.CLIMeta{Service: "fluxrpc", Endpoint: "getTransaction"})
			os.Exit(httpclient.ExitCodeForError(err))
		}
		output.PrintSuccess(cmd, result, &dto.CLIMeta{Service: "fluxrpc", Endpoint: "getTransaction"})
	},
}

var rpcTxCountCmd = &cobra.Command{
	Use:   "count",
	Short: "Get total transaction count (getTransactionCount)",
	Run: func(cmd *cobra.Command, args []string) {
		commitment, _ := cmd.Flags().GetString("commitment")
		result, err := fluxRPCSvc().GetTransactionCount(commitment)
		if err != nil {
			output.PrintError(cmd, "RPC_ERROR", err.Error(), &dto.CLIMeta{Service: "fluxrpc"})
			os.Exit(httpclient.ExitCodeForError(err))
		}
		output.PrintSuccess(cmd, result, &dto.CLIMeta{Service: "fluxrpc", Endpoint: "getTransactionCount"})
	},
}

var rpcTxSignaturesCmd = &cobra.Command{
	Use:   "signatures <address>",
	Short: "Get transaction signatures for an address (getSignaturesForAddress)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		limit, _ := cmd.Flags().GetInt("limit")
		before, _ := cmd.Flags().GetString("before")
		until, _ := cmd.Flags().GetString("until")
		commitment, _ := cmd.Flags().GetString("commitment")
		result, err := fluxRPCSvc().GetSignaturesForAddress(args[0], limit, before, until, commitment)
		if err != nil {
			output.PrintError(cmd, "RPC_ERROR", err.Error(), &dto.CLIMeta{Service: "fluxrpc"})
			os.Exit(httpclient.ExitCodeForError(err))
		}
		output.PrintSuccess(cmd, result, &dto.CLIMeta{Service: "fluxrpc", Endpoint: "getSignaturesForAddress"})
	},
}

var rpcTxStatusCmd = &cobra.Command{
	Use:   "status <sig1,sig2,...>",
	Short: "Get signature statuses (getSignatureStatuses)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		sigs := strings.Split(args[0], ",")
		searchHistory, _ := cmd.Flags().GetBool("search-history")
		result, err := fluxRPCSvc().GetSignatureStatuses(sigs, searchHistory)
		if err != nil {
			output.PrintError(cmd, "RPC_ERROR", err.Error(), &dto.CLIMeta{Service: "fluxrpc"})
			os.Exit(httpclient.ExitCodeForError(err))
		}
		output.PrintSuccess(cmd, result, &dto.CLIMeta{Service: "fluxrpc", Endpoint: "getSignatureStatuses"})
	},
}

var rpcTxFeeCmd = &cobra.Command{
	Use:   "fee <base64Message>",
	Short: "Get fee for a transaction message (getFeeForMessage)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		commitment, _ := cmd.Flags().GetString("commitment")
		result, err := fluxRPCSvc().GetFeeForMessage(args[0], commitment)
		if err != nil {
			output.PrintError(cmd, "RPC_ERROR", err.Error(), &dto.CLIMeta{Service: "fluxrpc"})
			os.Exit(httpclient.ExitCodeForError(err))
		}
		output.PrintSuccess(cmd, result, &dto.CLIMeta{Service: "fluxrpc", Endpoint: "getFeeForMessage"})
	},
}
