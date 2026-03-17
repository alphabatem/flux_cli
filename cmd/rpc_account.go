package cmd

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/alphabatem/flux_cli/dto"
	"github.com/alphabatem/flux_cli/internal/httpclient"
	pb "github.com/alphabatem/flux_cli/internal/yellowstonepb"
	"github.com/alphabatem/flux_cli/output"
	"github.com/spf13/cobra"
)

func init() {
	rpcCmd.AddCommand(rpcAccountCmd)

	rpcAccountShowCmd.Flags().String("encoding", "jsonParsed", "Data encoding: base58, base64, base64+zstd, jsonParsed")
	rpcAccountShowCmd.Flags().String("commitment", "", "Commitment level: processed, confirmed, finalized")
	rpcAccountCmd.AddCommand(rpcAccountShowCmd)

	rpcAccountBalanceCmd.Flags().String("commitment", "", "Commitment level")
	rpcAccountCmd.AddCommand(rpcAccountBalanceCmd)

	rpcAccountMultipleCmd.Flags().String("encoding", "jsonParsed", "Data encoding")
	rpcAccountMultipleCmd.Flags().String("commitment", "", "Commitment level")
	rpcAccountCmd.AddCommand(rpcAccountMultipleCmd)

	rpcAccountProgramCmd.Flags().String("encoding", "jsonParsed", "Data encoding")
	rpcAccountProgramCmd.Flags().String("commitment", "", "Commitment level")
	rpcAccountProgramCmd.Flags().String("filters", "", "JSON array of filters")
	rpcAccountCmd.AddCommand(rpcAccountProgramCmd)

	rpcAccountRentExemptionCmd.Flags().String("commitment", "", "Commitment level")
	rpcAccountCmd.AddCommand(rpcAccountRentExemptionCmd)

	rpcAccountLargestCmd.Flags().String("filter", "", "Filter: circulating or nonCirculating")
	rpcAccountLargestCmd.Flags().String("commitment", "", "Commitment level")
	rpcAccountCmd.AddCommand(rpcAccountLargestCmd)

	rpcAccountWatchCmd.Flags().String("commitment", "processed", "Commitment level: processed, confirmed, finalized")
	rpcAccountWatchCmd.Flags().Duration("timeout", 0*time.Second, "Optional stream timeout (e.g. 30s, 1m); 0 means no timeout")
	rpcAccountCmd.AddCommand(rpcAccountWatchCmd)

	rpcAccountWatchProgramCmd.Flags().String("commitment", "processed", "Commitment level: processed, confirmed, finalized")
	rpcAccountWatchProgramCmd.Flags().Duration("timeout", 0*time.Second, "Optional stream timeout (e.g. 30s, 1m); 0 means no timeout")
	rpcAccountCmd.AddCommand(rpcAccountWatchProgramCmd)

	rpcAccountWatchOwnerCmd.Flags().String("commitment", "processed", "Commitment level: processed, confirmed, finalized")
	rpcAccountWatchOwnerCmd.Flags().Duration("timeout", 0*time.Second, "Optional stream timeout (e.g. 30s, 1m); 0 means no timeout")
	rpcAccountCmd.AddCommand(rpcAccountWatchOwnerCmd)
}

var rpcAccountCmd = &cobra.Command{
	Use:   "account",
	Short: "Account methods - info, balance, program accounts, rent exemption",
}

var rpcAccountShowCmd = &cobra.Command{
	Use:   "show <pubkey>",
	Short: "Get account info (getAccountInfo)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		encoding, _ := cmd.Flags().GetString("encoding")
		commitment, _ := cmd.Flags().GetString("commitment")
		result, err := fluxRPCSvc().GetAccountInfo(args[0], encoding, commitment)
		if err != nil {
			output.PrintError(cmd, "RPC_ERROR", err.Error(), &dto.CLIMeta{Service: "fluxrpc", Endpoint: "getAccountInfo"})
			os.Exit(httpclient.ExitCodeForError(err))
		}
		output.PrintSuccess(cmd, result, &dto.CLIMeta{Service: "fluxrpc", Endpoint: "getAccountInfo"})
	},
}

var rpcAccountBalanceCmd = &cobra.Command{
	Use:   "balance <pubkey>",
	Short: "Get account balance in lamports (getBalance)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		commitment, _ := cmd.Flags().GetString("commitment")
		result, err := fluxRPCSvc().GetBalance(args[0], commitment)
		if err != nil {
			output.PrintError(cmd, "RPC_ERROR", err.Error(), &dto.CLIMeta{Service: "fluxrpc", Endpoint: "getBalance"})
			os.Exit(httpclient.ExitCodeForError(err))
		}
		output.PrintSuccess(cmd, result, &dto.CLIMeta{Service: "fluxrpc", Endpoint: "getBalance"})
	},
}

var rpcAccountMultipleCmd = &cobra.Command{
	Use:   "multiple <pubkey1,pubkey2,...>",
	Short: "Get multiple accounts in one request (getMultipleAccounts)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pubkeys := strings.Split(args[0], ",")
		encoding, _ := cmd.Flags().GetString("encoding")
		commitment, _ := cmd.Flags().GetString("commitment")
		result, err := fluxRPCSvc().GetMultipleAccounts(pubkeys, encoding, commitment)
		if err != nil {
			output.PrintError(cmd, "RPC_ERROR", err.Error(), &dto.CLIMeta{Service: "fluxrpc", Endpoint: "getMultipleAccounts"})
			os.Exit(httpclient.ExitCodeForError(err))
		}
		output.PrintSuccess(cmd, result, &dto.CLIMeta{Service: "fluxrpc", Endpoint: "getMultipleAccounts"})
	},
}

var rpcAccountProgramCmd = &cobra.Command{
	Use:   "program <programId>",
	Short: "Get all accounts owned by a program (getProgramAccounts)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		encoding, _ := cmd.Flags().GetString("encoding")
		filters, _ := cmd.Flags().GetString("filters")
		commitment, _ := cmd.Flags().GetString("commitment")
		result, err := fluxRPCSvc().GetProgramAccounts(args[0], encoding, filters, commitment)
		if err != nil {
			output.PrintError(cmd, "RPC_ERROR", err.Error(), &dto.CLIMeta{Service: "fluxrpc", Endpoint: "getProgramAccounts"})
			os.Exit(httpclient.ExitCodeForError(err))
		}
		output.PrintSuccess(cmd, result, &dto.CLIMeta{Service: "fluxrpc", Endpoint: "getProgramAccounts"})
	},
}

var rpcAccountRentExemptionCmd = &cobra.Command{
	Use:   "rent-exemption <dataLength>",
	Short: "Get minimum lamports for rent exemption (getMinimumBalanceForRentExemption)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		length, err := strconv.ParseUint(args[0], 10, 64)
		if err != nil {
			output.PrintError(cmd, "INVALID_ARGS", "dataLength must be a number", nil)
			os.Exit(dto.ExitUsageError)
		}
		commitment, _ := cmd.Flags().GetString("commitment")
		result, rpcErr := fluxRPCSvc().GetMinimumBalanceForRentExemption(length, commitment)
		if rpcErr != nil {
			output.PrintError(cmd, "RPC_ERROR", rpcErr.Error(), &dto.CLIMeta{Service: "fluxrpc"})
			os.Exit(httpclient.ExitCodeForError(rpcErr))
		}
		output.PrintSuccess(cmd, result, &dto.CLIMeta{Service: "fluxrpc", Endpoint: "getMinimumBalanceForRentExemption"})
	},
}

var rpcAccountLargestCmd = &cobra.Command{
	Use:   "largest",
	Short: "Get 20 largest accounts by balance (getLargestAccounts)",
	Run: func(cmd *cobra.Command, args []string) {
		filter, _ := cmd.Flags().GetString("filter")
		commitment, _ := cmd.Flags().GetString("commitment")
		result, err := fluxRPCSvc().GetLargestAccounts(filter, commitment)
		if err != nil {
			output.PrintError(cmd, "RPC_ERROR", err.Error(), &dto.CLIMeta{Service: "fluxrpc"})
			os.Exit(httpclient.ExitCodeForError(err))
		}
		output.PrintSuccess(cmd, result, &dto.CLIMeta{Service: "fluxrpc", Endpoint: "getLargestAccounts"})
	},
}

var rpcAccountWatchCmd = &cobra.Command{
	Use:   "watch <pubkey1,pubkey2,...>",
	Short: "Stream account updates via Yellowstone gRPC",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		commitmentRaw, _ := cmd.Flags().GetString("commitment")
		commitment, err := parseCommitment(commitmentRaw)
		if err != nil {
			failUsage(cmd, err.Error())
		}

		accounts := parseCSV(args[0])
		if len(accounts) == 0 {
			failUsage(cmd, "at least one account pubkey is required")
		}

		streamCtx, cancel, err := streamContextFromTimeoutFlag(cmd)
		if err != nil {
			failUsage(cmd, err.Error())
		}
		defer cancel()

		err = yellowstoneSvc().WatchAccounts(streamCtx, accounts, commitment, func(update *pb.SubscribeUpdate) error {
			return printWatchUpdate(cmd, "account.watch", update)
		})
		if err != nil {
			output.PrintError(cmd, "WATCH_ERROR", err.Error(), &dto.CLIMeta{Service: "yellowstone", Endpoint: "account.watch"})
			os.Exit(dto.ExitGeneralError)
		}
	},
}

var rpcAccountWatchProgramCmd = &cobra.Command{
	Use:   "watch-program <programId1,programId2,...>",
	Short: "Stream updates for all accounts owned by program(s) via Yellowstone gRPC",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		commitmentRaw, _ := cmd.Flags().GetString("commitment")
		commitment, err := parseCommitment(commitmentRaw)
		if err != nil {
			failUsage(cmd, err.Error())
		}

		owners := parseCSV(args[0])
		if len(owners) == 0 {
			failUsage(cmd, "at least one program id is required")
		}

		streamCtx, cancel, err := streamContextFromTimeoutFlag(cmd)
		if err != nil {
			failUsage(cmd, err.Error())
		}
		defer cancel()

		err = yellowstoneSvc().WatchProgramOwners(streamCtx, owners, commitment, func(update *pb.SubscribeUpdate) error {
			return printWatchUpdate(cmd, "account.watch-program", update)
		})
		if err != nil {
			output.PrintError(cmd, "WATCH_ERROR", err.Error(), &dto.CLIMeta{Service: "yellowstone", Endpoint: "account.watch-program"})
			os.Exit(dto.ExitGeneralError)
		}
	},
}

var rpcAccountWatchOwnerCmd = &cobra.Command{
	Use:   "watch-owner <owner1,owner2,...>",
	Short: "Stream updates for accounts owned by owner/program via Yellowstone gRPC",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		commitmentRaw, _ := cmd.Flags().GetString("commitment")
		commitment, err := parseCommitment(commitmentRaw)
		if err != nil {
			failUsage(cmd, err.Error())
		}

		owners := parseCSV(args[0])
		if len(owners) == 0 {
			failUsage(cmd, "at least one owner is required")
		}

		streamCtx, cancel, err := streamContextFromTimeoutFlag(cmd)
		if err != nil {
			failUsage(cmd, err.Error())
		}
		defer cancel()

		err = yellowstoneSvc().WatchProgramOwners(streamCtx, owners, commitment, func(update *pb.SubscribeUpdate) error {
			return printWatchUpdate(cmd, "account.watch-owner", update)
		})
		if err != nil {
			output.PrintError(cmd, "WATCH_ERROR", err.Error(), &dto.CLIMeta{Service: "yellowstone", Endpoint: "account.watch-owner"})
			os.Exit(dto.ExitGeneralError)
		}
	},
}
