package cmd

import (
	"os"
	"strconv"
	"strings"

	"github.com/alphabatem/flux_cli/dto"
	"github.com/alphabatem/flux_cli/internal/httpclient"
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
