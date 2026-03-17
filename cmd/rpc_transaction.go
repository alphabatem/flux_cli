package cmd

import (
	"os"
	"strings"
	"time"

	"github.com/alphabatem/flux_cli/dto"
	"github.com/alphabatem/flux_cli/internal/httpclient"
	pb "github.com/alphabatem/flux_cli/internal/yellowstonepb"
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

	rpcTxWatchCmd.Flags().String("commitment", "processed", "Commitment level: processed, confirmed, finalized")
	rpcTxWatchCmd.Flags().String("exclude", "", "Comma-separated excluded accounts")
	rpcTxWatchCmd.Flags().String("required", "", "Comma-separated required accounts")
	rpcTxWatchCmd.Flags().Bool("include-votes", false, "Include vote transactions")
	rpcTxWatchCmd.Flags().Bool("include-failed", false, "Include failed transactions")
	rpcTxWatchCmd.Flags().Duration("timeout", 0*time.Second, "Optional stream timeout (e.g. 30s, 1m); 0 means no timeout")
	rpcTransactionCmd.AddCommand(rpcTxWatchCmd)
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

var rpcTxWatchCmd = &cobra.Command{
	Use:   "watch <account1,account2,...>",
	Short: "Stream transaction updates via Yellowstone gRPC",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		commitmentRaw, _ := cmd.Flags().GetString("commitment")
		commitment, err := parseCommitment(commitmentRaw)
		if err != nil {
			failUsage(cmd, err.Error())
		}

		accountInclude := parseCSV(args[0])
		if len(accountInclude) == 0 {
			failUsage(cmd, "at least one include account is required")
		}

		excludeRaw, _ := cmd.Flags().GetString("exclude")
		accountExclude := parseCSV(excludeRaw)
		requiredRaw, _ := cmd.Flags().GetString("required")
		accountRequired := parseCSV(requiredRaw)
		includeVotes, _ := cmd.Flags().GetBool("include-votes")
		includeFailed, _ := cmd.Flags().GetBool("include-failed")
		streamCtx, cancel, err := streamContextFromTimeoutFlag(cmd)
		if err != nil {
			failUsage(cmd, err.Error())
		}
		defer cancel()

		err = yellowstoneSvc().WatchTransactions(
			streamCtx,
			accountInclude,
			accountExclude,
			accountRequired,
			includeVotes,
			includeFailed,
			commitment,
			func(update *pb.SubscribeUpdate) error {
				return printWatchUpdate(cmd, "transaction.watch", update)
			},
		)
		if err != nil {
			output.PrintError(cmd, "WATCH_ERROR", err.Error(), &dto.CLIMeta{Service: "yellowstone", Endpoint: "transaction.watch"})
			os.Exit(dto.ExitGeneralError)
		}
	},
}
