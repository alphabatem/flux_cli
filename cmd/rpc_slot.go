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
	rpcCmd.AddCommand(rpcSlotCmd)

	rpcSlotShowCmd.Flags().String("commitment", "", "Commitment level")
	rpcSlotCmd.AddCommand(rpcSlotShowCmd)

	rpcSlotLeaderCmd.Flags().String("commitment", "", "Commitment level")
	rpcSlotCmd.AddCommand(rpcSlotLeaderCmd)

	rpcSlotCmd.AddCommand(rpcSlotLeadersCmd)
	rpcSlotCmd.AddCommand(rpcSlotMinLedgerCmd)
	rpcSlotCmd.AddCommand(rpcSlotMaxRetransmitCmd)
	rpcSlotCmd.AddCommand(rpcSlotMaxShredCmd)
	rpcSlotCmd.AddCommand(rpcSlotHighestSnapshotCmd)
}

var rpcSlotCmd = &cobra.Command{
	Use:   "slot",
	Short: "Slot methods - current, leader, leaders, ledger, snapshot",
}

var rpcSlotShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Get current slot (getSlot)",
	Run: func(cmd *cobra.Command, args []string) {
		commitment, _ := cmd.Flags().GetString("commitment")
		result, err := fluxRPCSvc().GetSlot(commitment)
		if err != nil {
			output.PrintError(cmd, "RPC_ERROR", err.Error(), &dto.CLIMeta{Service: "fluxrpc"})
			os.Exit(httpclient.ExitCodeForError(err))
		}
		output.PrintSuccess(cmd, result, &dto.CLIMeta{Service: "fluxrpc", Endpoint: "getSlot"})
	},
}

var rpcSlotLeaderCmd = &cobra.Command{
	Use:   "leader",
	Short: "Get current slot leader (getSlotLeader)",
	Run: func(cmd *cobra.Command, args []string) {
		commitment, _ := cmd.Flags().GetString("commitment")
		result, err := fluxRPCSvc().GetSlotLeader(commitment)
		if err != nil {
			output.PrintError(cmd, "RPC_ERROR", err.Error(), &dto.CLIMeta{Service: "fluxrpc"})
			os.Exit(httpclient.ExitCodeForError(err))
		}
		output.PrintSuccess(cmd, result, &dto.CLIMeta{Service: "fluxrpc", Endpoint: "getSlotLeader"})
	},
}

var rpcSlotLeadersCmd = &cobra.Command{
	Use:   "leaders <startSlot> <limit>",
	Short: "Get slot leaders for a range (getSlotLeaders)",
	Args:  cobra.ExactArgs(2),
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
		result, rpcErr := fluxRPCSvc().GetSlotLeaders(start, limit)
		if rpcErr != nil {
			output.PrintError(cmd, "RPC_ERROR", rpcErr.Error(), &dto.CLIMeta{Service: "fluxrpc"})
			os.Exit(httpclient.ExitCodeForError(rpcErr))
		}
		output.PrintSuccess(cmd, result, &dto.CLIMeta{Service: "fluxrpc", Endpoint: "getSlotLeaders"})
	},
}

var rpcSlotMinLedgerCmd = &cobra.Command{
	Use:   "minimum-ledger",
	Short: "Get lowest available ledger slot (minimumLedgerSlot)",
	Run: func(cmd *cobra.Command, args []string) {
		result, err := fluxRPCSvc().MinimumLedgerSlot()
		if err != nil {
			output.PrintError(cmd, "RPC_ERROR", err.Error(), &dto.CLIMeta{Service: "fluxrpc"})
			os.Exit(httpclient.ExitCodeForError(err))
		}
		output.PrintSuccess(cmd, result, &dto.CLIMeta{Service: "fluxrpc", Endpoint: "minimumLedgerSlot"})
	},
}

var rpcSlotMaxRetransmitCmd = &cobra.Command{
	Use:   "max-retransmit",
	Short: "Get max retransmit slot (getMaxRetransmitSlot)",
	Run: func(cmd *cobra.Command, args []string) {
		result, err := fluxRPCSvc().GetMaxRetransmitSlot()
		if err != nil {
			output.PrintError(cmd, "RPC_ERROR", err.Error(), &dto.CLIMeta{Service: "fluxrpc"})
			os.Exit(httpclient.ExitCodeForError(err))
		}
		output.PrintSuccess(cmd, result, &dto.CLIMeta{Service: "fluxrpc", Endpoint: "getMaxRetransmitSlot"})
	},
}

var rpcSlotMaxShredCmd = &cobra.Command{
	Use:   "max-shred-insert",
	Short: "Get max shred insert slot (getMaxShredInsertSlot)",
	Run: func(cmd *cobra.Command, args []string) {
		result, err := fluxRPCSvc().GetMaxShredInsertSlot()
		if err != nil {
			output.PrintError(cmd, "RPC_ERROR", err.Error(), &dto.CLIMeta{Service: "fluxrpc"})
			os.Exit(httpclient.ExitCodeForError(err))
		}
		output.PrintSuccess(cmd, result, &dto.CLIMeta{Service: "fluxrpc", Endpoint: "getMaxShredInsertSlot"})
	},
}

var rpcSlotHighestSnapshotCmd = &cobra.Command{
	Use:   "highest-snapshot",
	Short: "Get highest snapshot slot (getHighestSnapshotSlot)",
	Run: func(cmd *cobra.Command, args []string) {
		result, err := fluxRPCSvc().GetHighestSnapshotSlot()
		if err != nil {
			output.PrintError(cmd, "RPC_ERROR", err.Error(), &dto.CLIMeta{Service: "fluxrpc"})
			os.Exit(httpclient.ExitCodeForError(err))
		}
		output.PrintSuccess(cmd, result, &dto.CLIMeta{Service: "fluxrpc", Endpoint: "getHighestSnapshotSlot"})
	},
}
