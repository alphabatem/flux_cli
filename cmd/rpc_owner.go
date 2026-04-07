package cmd

import (
	"os"
	"time"

	"github.com/alphabatem/flux_cli/dto"
	pb "github.com/alphabatem/flux_cli/internal/yellowstonepb"
	"github.com/alphabatem/flux_cli/output"
	"github.com/spf13/cobra"
)

func init() {
	rpcCmd.AddCommand(rpcOwnerCmd)

	rpcOwnerWatchCmd.Flags().String("commitment", "processed", "Commitment level: processed, confirmed, finalized")
	rpcOwnerWatchCmd.Flags().Duration("timeout", 0*time.Second, "Optional stream timeout (e.g. 30s, 1m); 0 means no timeout")
	rpcOwnerCmd.AddCommand(rpcOwnerWatchCmd)
}

var rpcOwnerCmd = &cobra.Command{
	Use:   "owner",
	Short: "Owner methods - watch accounts by owner",
}

var rpcOwnerWatchCmd = &cobra.Command{
	Use:   "watch <owner1,owner2,...>",
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
			return printWatchUpdate(cmd, "owner.watch", update)
		})
		if err != nil {
			output.PrintError(cmd, "WATCH_ERROR", err.Error(), &dto.CLIMeta{Service: "yellowstone", Endpoint: "owner.watch"})
			os.Exit(dto.ExitGeneralError)
		}
	},
}
