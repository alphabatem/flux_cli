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
	rpcCmd.AddCommand(rpcBlockCmd)

	rpcBlockShowCmd.Flags().String("encoding", "json", "Encoding: json, jsonParsed, base64, base58")
	rpcBlockShowCmd.Flags().String("tx-detail", "full", "Transaction detail: full, accounts, signatures, none")
	rpcBlockShowCmd.Flags().String("commitment", "", "Commitment level: confirmed, finalized")
	rpcBlockCmd.AddCommand(rpcBlockShowCmd)

	rpcBlockListCmd.Flags().String("commitment", "", "Commitment level")
	rpcBlockCmd.AddCommand(rpcBlockListCmd)

	rpcBlockListLimitCmd.Flags().String("commitment", "", "Commitment level")
	rpcBlockCmd.AddCommand(rpcBlockListLimitCmd)

	rpcBlockHeightCmd.Flags().String("commitment", "", "Commitment level")
	rpcBlockCmd.AddCommand(rpcBlockHeightCmd)

	rpcBlockCmd.AddCommand(rpcBlockTimeCmd)
	rpcBlockCmd.AddCommand(rpcBlockCommitmentCmd)

	rpcBlockProductionCmd.Flags().String("identity", "", "Filter to specific validator identity")
	rpcBlockProductionCmd.Flags().String("commitment", "", "Commitment level")
	rpcBlockCmd.AddCommand(rpcBlockProductionCmd)

	rpcBlockLatestHashCmd.Flags().String("commitment", "", "Commitment level")
	rpcBlockCmd.AddCommand(rpcBlockLatestHashCmd)

	rpcBlockHashValidCmd.Flags().String("commitment", "", "Commitment level")
	rpcBlockCmd.AddCommand(rpcBlockHashValidCmd)
}

var rpcBlockCmd = &cobra.Command{
	Use:   "block",
	Short: "Block methods - info, list, height, time, production, blockhash",
}

var rpcBlockShowCmd = &cobra.Command{
	Use:   "show <slot>",
	Short: "Get block data with transactions (getBlock)",
	Args:  exactArgsFromUse(),
	Run: func(cmd *cobra.Command, args []string) {
		slot, err := strconv.ParseUint(args[0], 10, 64)
		if err != nil {
			output.PrintError(cmd, "INVALID_ARGS", "slot must be a number", nil)
			os.Exit(dto.ExitUsageError)
		}
		encoding, _ := cmd.Flags().GetString("encoding")
		txDetail, _ := cmd.Flags().GetString("tx-detail")
		commitment, _ := cmd.Flags().GetString("commitment")
		result, rpcErr := fluxRPCSvc().GetBlock(slot, encoding, txDetail, commitment)
		if rpcErr != nil {
			output.PrintError(cmd, "RPC_ERROR", rpcErr.Error(), &dto.CLIMeta{Service: "fluxrpc", Endpoint: "getBlock"})
			os.Exit(httpclient.ExitCodeForError(rpcErr))
		}
		output.PrintSuccess(cmd, result, &dto.CLIMeta{Service: "fluxrpc", Endpoint: "getBlock"})
	},
}

var rpcBlockListCmd = &cobra.Command{
	Use:   "list <startSlot> [endSlot]",
	Short: "Get confirmed blocks between slots (getBlocks)",
	Args:  rangeArgsFromUse(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		start, err := strconv.ParseUint(args[0], 10, 64)
		if err != nil {
			output.PrintError(cmd, "INVALID_ARGS", "startSlot must be a number", nil)
			os.Exit(dto.ExitUsageError)
		}
		var end uint64
		if len(args) > 1 {
			end, err = strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				output.PrintError(cmd, "INVALID_ARGS", "endSlot must be a number", nil)
				os.Exit(dto.ExitUsageError)
			}
		}
		commitment, _ := cmd.Flags().GetString("commitment")
		result, rpcErr := fluxRPCSvc().GetBlocks(start, end, commitment)
		if rpcErr != nil {
			output.PrintError(cmd, "RPC_ERROR", rpcErr.Error(), &dto.CLIMeta{Service: "fluxrpc", Endpoint: "getBlocks"})
			os.Exit(httpclient.ExitCodeForError(rpcErr))
		}
		output.PrintSuccess(cmd, result, &dto.CLIMeta{Service: "fluxrpc", Endpoint: "getBlocks"})
	},
}

var rpcBlockListLimitCmd = &cobra.Command{
	Use:   "list-limit <startSlot> <limit>",
	Short: "Get confirmed blocks starting at slot for up to limit blocks (getBlocksWithLimit)",
	Args:  exactArgsFromUse(),
	Run: func(cmd *cobra.Command, args []string) {
		start, err := strconv.ParseUint(args[0], 10, 64)
		if err != nil {
			output.PrintError(cmd, "INVALID_ARGS", "startSlot must be a number", nil)
			os.Exit(dto.ExitUsageError)
		}
		limit, err := strconv.ParseUint(args[1], 10, 64)
		if err != nil {
			output.PrintError(cmd, "INVALID_ARGS", "limit must be a number", nil)
			os.Exit(dto.ExitUsageError)
		}
		commitment, _ := cmd.Flags().GetString("commitment")
		result, rpcErr := fluxRPCSvc().GetBlocksWithLimit(start, limit, commitment)
		if rpcErr != nil {
			output.PrintError(cmd, "RPC_ERROR", rpcErr.Error(), &dto.CLIMeta{Service: "fluxrpc"})
			os.Exit(httpclient.ExitCodeForError(rpcErr))
		}
		output.PrintSuccess(cmd, result, &dto.CLIMeta{Service: "fluxrpc", Endpoint: "getBlocksWithLimit"})
	},
}

var rpcBlockHeightCmd = &cobra.Command{
	Use:   "height",
	Short: "Get current block height (getBlockHeight)",
	Run: func(cmd *cobra.Command, args []string) {
		commitment, _ := cmd.Flags().GetString("commitment")
		result, err := fluxRPCSvc().GetBlockHeight(commitment)
		if err != nil {
			output.PrintError(cmd, "RPC_ERROR", err.Error(), &dto.CLIMeta{Service: "fluxrpc"})
			os.Exit(httpclient.ExitCodeForError(err))
		}
		output.PrintSuccess(cmd, result, &dto.CLIMeta{Service: "fluxrpc", Endpoint: "getBlockHeight"})
	},
}

var rpcBlockTimeCmd = &cobra.Command{
	Use:   "time <slot>",
	Short: "Get estimated production time for a block (getBlockTime)",
	Args:  exactArgsFromUse(),
	Run: func(cmd *cobra.Command, args []string) {
		slot, err := strconv.ParseUint(args[0], 10, 64)
		if err != nil {
			output.PrintError(cmd, "INVALID_ARGS", "slot must be a number", nil)
			os.Exit(dto.ExitUsageError)
		}
		result, rpcErr := fluxRPCSvc().GetBlockTime(slot)
		if rpcErr != nil {
			output.PrintError(cmd, "RPC_ERROR", rpcErr.Error(), &dto.CLIMeta{Service: "fluxrpc"})
			os.Exit(httpclient.ExitCodeForError(rpcErr))
		}
		output.PrintSuccess(cmd, result, &dto.CLIMeta{Service: "fluxrpc", Endpoint: "getBlockTime"})
	},
}

var rpcBlockCommitmentCmd = &cobra.Command{
	Use:   "commitment <slot>",
	Short: "Get commitment status for a block (getBlockCommitment)",
	Args:  exactArgsFromUse(),
	Run: func(cmd *cobra.Command, args []string) {
		slot, err := strconv.ParseUint(args[0], 10, 64)
		if err != nil {
			output.PrintError(cmd, "INVALID_ARGS", "slot must be a number", nil)
			os.Exit(dto.ExitUsageError)
		}
		result, rpcErr := fluxRPCSvc().GetBlockCommitment(slot)
		if rpcErr != nil {
			output.PrintError(cmd, "RPC_ERROR", rpcErr.Error(), &dto.CLIMeta{Service: "fluxrpc"})
			os.Exit(httpclient.ExitCodeForError(rpcErr))
		}
		output.PrintSuccess(cmd, result, &dto.CLIMeta{Service: "fluxrpc", Endpoint: "getBlockCommitment"})
	},
}

var rpcBlockProductionCmd = &cobra.Command{
	Use:   "production",
	Short: "Get block production metrics (getBlockProduction)",
	Run: func(cmd *cobra.Command, args []string) {
		identity, _ := cmd.Flags().GetString("identity")
		commitment, _ := cmd.Flags().GetString("commitment")
		result, err := fluxRPCSvc().GetBlockProduction(identity, commitment)
		if err != nil {
			output.PrintError(cmd, "RPC_ERROR", err.Error(), &dto.CLIMeta{Service: "fluxrpc"})
			os.Exit(httpclient.ExitCodeForError(err))
		}
		output.PrintSuccess(cmd, result, &dto.CLIMeta{Service: "fluxrpc", Endpoint: "getBlockProduction"})
	},
}

var rpcBlockLatestHashCmd = &cobra.Command{
	Use:   "latest-blockhash",
	Short: "Get latest blockhash and last valid block height (getLatestBlockhash)",
	Run: func(cmd *cobra.Command, args []string) {
		commitment, _ := cmd.Flags().GetString("commitment")
		result, err := fluxRPCSvc().GetLatestBlockhash(commitment)
		if err != nil {
			output.PrintError(cmd, "RPC_ERROR", err.Error(), &dto.CLIMeta{Service: "fluxrpc"})
			os.Exit(httpclient.ExitCodeForError(err))
		}
		output.PrintSuccess(cmd, result, &dto.CLIMeta{Service: "fluxrpc", Endpoint: "getLatestBlockhash"})
	},
}

var rpcBlockHashValidCmd = &cobra.Command{
	Use:   "blockhash-valid <blockhash>",
	Short: "Check if a blockhash is still valid (isBlockhashValid)",
	Args:  exactArgsFromUse(),
	Run: func(cmd *cobra.Command, args []string) {
		commitment, _ := cmd.Flags().GetString("commitment")
		result, err := fluxRPCSvc().IsBlockhashValid(args[0], commitment)
		if err != nil {
			output.PrintError(cmd, "RPC_ERROR", err.Error(), &dto.CLIMeta{Service: "fluxrpc"})
			os.Exit(httpclient.ExitCodeForError(err))
		}
		output.PrintSuccess(cmd, result, &dto.CLIMeta{Service: "fluxrpc", Endpoint: "isBlockhashValid"})
	},
}
