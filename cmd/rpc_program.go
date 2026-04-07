package cmd

import (
	"os"
	"time"

	"github.com/alphabatem/flux_cli/dto"
	"github.com/alphabatem/flux_cli/internal/httpclient"
	pb "github.com/alphabatem/flux_cli/internal/yellowstonepb"
	"github.com/alphabatem/flux_cli/output"
	"github.com/spf13/cobra"
)

func init() {
	rpcCmd.AddCommand(rpcProgramCmd)

	rpcProgramAccountsCmd.Flags().String("encoding", "jsonParsed", "Data encoding")
	rpcProgramAccountsCmd.Flags().String("commitment", "", "Commitment level")
	rpcProgramAccountsCmd.Flags().String("filters", "", "JSON array of filters")
	rpcProgramCmd.AddCommand(rpcProgramAccountsCmd)

	rpcProgramWatchCmd.Flags().String("commitment", "processed", "Commitment level: processed, confirmed, finalized")
	rpcProgramWatchCmd.Flags().Duration("timeout", 0*time.Second, "Optional stream timeout (e.g. 30s, 1m); 0 means no timeout")
	rpcProgramCmd.AddCommand(rpcProgramWatchCmd)
}

var rpcProgramCmd = &cobra.Command{
	Use:   "program",
	Short: "Program methods - accounts, watch program updates",
}

var rpcProgramAccountsCmd = &cobra.Command{
	Use:   "accounts <programId>",
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

var rpcProgramWatchCmd = &cobra.Command{
	Use:   "watch <programId1,programId2,...>",
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
			return printWatchUpdate(cmd, "program.watch", update)
		})
		if err != nil {
			output.PrintError(cmd, "WATCH_ERROR", err.Error(), &dto.CLIMeta{Service: "yellowstone", Endpoint: "program.watch"})
			os.Exit(dto.ExitGeneralError)
		}
	},
}
