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
	rpcCmd.AddCommand(rpcEpochCmd)

	rpcEpochInfoCmd.Flags().String("commitment", "", "Commitment level")
	rpcEpochCmd.AddCommand(rpcEpochInfoCmd)

	rpcEpochCmd.AddCommand(rpcEpochScheduleCmd)

	rpcEpochLeaderScheduleCmd.Flags().String("identity", "", "Filter to specific validator")
	rpcEpochLeaderScheduleCmd.Flags().String("commitment", "", "Commitment level")
	rpcEpochCmd.AddCommand(rpcEpochLeaderScheduleCmd)
}

var rpcEpochCmd = &cobra.Command{
	Use:   "epoch",
	Short: "Epoch methods - info, schedule, leader schedule",
}

var rpcEpochInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Get current epoch info (getEpochInfo)",
	Run: func(cmd *cobra.Command, args []string) {
		commitment, _ := cmd.Flags().GetString("commitment")
		result, err := fluxRPCSvc().GetEpochInfo(commitment)
		if err != nil {
			output.PrintError(cmd, "RPC_ERROR", err.Error(), &dto.CLIMeta{Service: "fluxrpc"})
			os.Exit(httpclient.ExitCodeForError(err))
		}
		output.PrintSuccess(cmd, result, &dto.CLIMeta{Service: "fluxrpc", Endpoint: "getEpochInfo"})
	},
}

var rpcEpochScheduleCmd = &cobra.Command{
	Use:   "schedule",
	Short: "Get epoch schedule configuration (getEpochSchedule)",
	Run: func(cmd *cobra.Command, args []string) {
		result, err := fluxRPCSvc().GetEpochSchedule()
		if err != nil {
			output.PrintError(cmd, "RPC_ERROR", err.Error(), &dto.CLIMeta{Service: "fluxrpc"})
			os.Exit(httpclient.ExitCodeForError(err))
		}
		output.PrintSuccess(cmd, result, &dto.CLIMeta{Service: "fluxrpc", Endpoint: "getEpochSchedule"})
	},
}

var rpcEpochLeaderScheduleCmd = &cobra.Command{
	Use:   "leader-schedule [slot]",
	Short: "Get leader schedule for an epoch (getLeaderSchedule)",
	Args:  maximumArgsFromUse(1),
	Run: func(cmd *cobra.Command, args []string) {
		var slot uint64
		if len(args) > 0 {
			var err error
			slot, err = strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				output.PrintError(cmd, "INVALID_ARGS", "slot must be a number", nil)
				os.Exit(dto.ExitUsageError)
			}
		}
		identity, _ := cmd.Flags().GetString("identity")
		commitment, _ := cmd.Flags().GetString("commitment")
		result, rpcErr := fluxRPCSvc().GetLeaderSchedule(slot, identity, commitment)
		if rpcErr != nil {
			output.PrintError(cmd, "RPC_ERROR", rpcErr.Error(), &dto.CLIMeta{Service: "fluxrpc"})
			os.Exit(httpclient.ExitCodeForError(rpcErr))
		}
		output.PrintSuccess(cmd, result, &dto.CLIMeta{Service: "fluxrpc", Endpoint: "getLeaderSchedule"})
	},
}
