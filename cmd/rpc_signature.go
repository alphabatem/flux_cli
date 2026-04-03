package cmd

import (
	"context"
	"os"
	"time"

	"github.com/alphabatem/flux_cli/dto"
	pb "github.com/alphabatem/flux_cli/internal/yellowstonepb"
	"github.com/alphabatem/flux_cli/output"
	"github.com/spf13/cobra"
)

func init() {
	rpcCmd.AddCommand(rpcSignatureCmd)

	rpcSignatureWatchCmd.Flags().String("commitment", "processed", "Commitment level: processed, confirmed, finalized")
	rpcSignatureWatchCmd.Flags().Duration("timeout", 0*time.Second, "Optional stream timeout (e.g. 30s, 1m); 0 means no timeout")
	rpcSignatureCmd.AddCommand(rpcSignatureWatchCmd)

	rpcSignatureConfirmCmd.Flags().String("commitment", "confirmed", "Commitment level: processed, confirmed, finalized")
	rpcSignatureConfirmCmd.Flags().Duration("timeout", 30*time.Second, "Max time to wait for confirmation (e.g. 15s, 1m)")
	rpcSignatureCmd.AddCommand(rpcSignatureConfirmCmd)
}

var rpcSignatureCmd = &cobra.Command{
	Use:   "signature",
	Short: "Signature operations - watch and confirm",
}

var rpcSignatureWatchCmd = &cobra.Command{
	Use:   "watch <signature>",
	Short: "Stream updates for a specific transaction signature via Yellowstone gRPC",
	Args:  exactArgsFromUse(),
	Run: func(cmd *cobra.Command, args []string) {
		commitmentRaw, _ := cmd.Flags().GetString("commitment")
		commitment, err := parseCommitment(commitmentRaw)
		if err != nil {
			failUsage(cmd, err.Error())
		}

		streamCtx, cancel, err := streamContextFromTimeoutFlag(cmd)
		if err != nil {
			failUsage(cmd, err.Error())
		}
		defer cancel()

		err = yellowstoneSvc().WatchTransactionSignature(streamCtx, args[0], commitment, func(update *pb.SubscribeUpdate) error {
			return printWatchUpdate(cmd, "signature.watch", update)
		})
		if err != nil {
			output.PrintError(cmd, "WATCH_ERROR", err.Error(), &dto.CLIMeta{Service: "yellowstone", Endpoint: "signature.watch"})
			os.Exit(dto.ExitGeneralError)
		}
	},
}

var rpcSignatureConfirmCmd = &cobra.Command{
	Use:   "confirm <signature>",
	Short: "Wait for transaction confirmation via Yellowstone gRPC",
	Args:  exactArgsFromUse(),
	Run: func(cmd *cobra.Command, args []string) {
		commitmentRaw, _ := cmd.Flags().GetString("commitment")
		commitment, err := parseCommitment(commitmentRaw)
		if err != nil {
			failUsage(cmd, err.Error())
		}

		timeout, _ := cmd.Flags().GetDuration("timeout")
		ctx := cmd.Context()
		cancel := func() {}
		if timeout > 0 {
			ctx, cancel = context.WithTimeout(ctx, timeout)
		} else {
			ctx, cancel = context.WithCancel(ctx)
		}
		defer cancel()

		var matched interface{}
		err = yellowstoneSvc().WatchTransactionSignature(ctx, args[0], commitment, func(update *pb.SubscribeUpdate) error {
			if update.GetTransaction() == nil && update.GetTransactionStatus() == nil {
				return nil
			}
			data, convErr := protoUpdateToData(update)
			if convErr != nil {
				return convErr
			}
			matched = data
			cancel()
			return nil
		})
		if err != nil {
			output.PrintError(cmd, "WATCH_ERROR", err.Error(), &dto.CLIMeta{Service: "yellowstone", Endpoint: "signature.confirm"})
			os.Exit(dto.ExitGeneralError)
		}
		if matched == nil {
			if ctx.Err() == context.DeadlineExceeded {
				output.PrintError(cmd, "TIMEOUT", "transaction confirmation timeout reached", &dto.CLIMeta{Service: "yellowstone", Endpoint: "signature.confirm"})
				os.Exit(dto.ExitGeneralError)
			}
			output.PrintError(cmd, "WATCH_ERROR", "no confirmation received", &dto.CLIMeta{Service: "yellowstone", Endpoint: "signature.confirm"})
			os.Exit(dto.ExitGeneralError)
		}

		output.PrintSuccess(cmd, map[string]interface{}{
			"signature": args[0],
			"confirmed": true,
			"update":    matched,
		}, &dto.CLIMeta{Service: "yellowstone", Endpoint: "signature.confirm"})
	},
}
