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
	rpcCmd.AddCommand(rpcStakingCmd)

	rpcStakingVoteAccountsCmd.Flags().String("vote-pubkey", "", "Filter to specific vote account")
	rpcStakingVoteAccountsCmd.Flags().String("commitment", "", "Commitment level")
	rpcStakingCmd.AddCommand(rpcStakingVoteAccountsCmd)

	rpcStakingMinDelegationCmd.Flags().String("commitment", "", "Commitment level")
	rpcStakingCmd.AddCommand(rpcStakingMinDelegationCmd)

	rpcStakingInflationGovCmd.Flags().String("commitment", "", "Commitment level")
	rpcStakingCmd.AddCommand(rpcStakingInflationGovCmd)

	rpcStakingCmd.AddCommand(rpcStakingInflationRateCmd)

	rpcStakingInflationRewardCmd.Flags().Uint64("epoch", 0, "Epoch for which to return rewards")
	rpcStakingInflationRewardCmd.Flags().String("commitment", "", "Commitment level")
	rpcStakingCmd.AddCommand(rpcStakingInflationRewardCmd)
}

var rpcStakingCmd = &cobra.Command{
	Use:   "staking",
	Short: "Staking/inflation methods - vote accounts, delegation, inflation rates, rewards",
}

var rpcStakingVoteAccountsCmd = &cobra.Command{
	Use:   "vote-accounts",
	Short: "Get vote accounts with activation info (getVoteAccounts)",
	Run: func(cmd *cobra.Command, args []string) {
		votePubkey, _ := cmd.Flags().GetString("vote-pubkey")
		commitment, _ := cmd.Flags().GetString("commitment")
		result, err := fluxRPCSvc().GetVoteAccounts(votePubkey, commitment)
		if err != nil {
			output.PrintError(cmd, "RPC_ERROR", err.Error(), &dto.CLIMeta{Service: "fluxrpc"})
			os.Exit(httpclient.ExitCodeForError(err))
		}
		output.PrintSuccess(cmd, result, &dto.CLIMeta{Service: "fluxrpc", Endpoint: "getVoteAccounts"})
	},
}

var rpcStakingMinDelegationCmd = &cobra.Command{
	Use:   "min-delegation",
	Short: "Get minimum stake delegation (getStakeMinimumDelegation)",
	Run: func(cmd *cobra.Command, args []string) {
		commitment, _ := cmd.Flags().GetString("commitment")
		result, err := fluxRPCSvc().GetStakeMinimumDelegation(commitment)
		if err != nil {
			output.PrintError(cmd, "RPC_ERROR", err.Error(), &dto.CLIMeta{Service: "fluxrpc"})
			os.Exit(httpclient.ExitCodeForError(err))
		}
		output.PrintSuccess(cmd, result, &dto.CLIMeta{Service: "fluxrpc", Endpoint: "getStakeMinimumDelegation"})
	},
}

var rpcStakingInflationGovCmd = &cobra.Command{
	Use:   "inflation-governor",
	Short: "Get inflation governor parameters (getInflationGovernor)",
	Run: func(cmd *cobra.Command, args []string) {
		commitment, _ := cmd.Flags().GetString("commitment")
		result, err := fluxRPCSvc().GetInflationGovernor(commitment)
		if err != nil {
			output.PrintError(cmd, "RPC_ERROR", err.Error(), &dto.CLIMeta{Service: "fluxrpc"})
			os.Exit(httpclient.ExitCodeForError(err))
		}
		output.PrintSuccess(cmd, result, &dto.CLIMeta{Service: "fluxrpc", Endpoint: "getInflationGovernor"})
	},
}

var rpcStakingInflationRateCmd = &cobra.Command{
	Use:   "inflation-rate",
	Short: "Get current inflation rate (getInflationRate)",
	Run: func(cmd *cobra.Command, args []string) {
		result, err := fluxRPCSvc().GetInflationRate()
		if err != nil {
			output.PrintError(cmd, "RPC_ERROR", err.Error(), &dto.CLIMeta{Service: "fluxrpc"})
			os.Exit(httpclient.ExitCodeForError(err))
		}
		output.PrintSuccess(cmd, result, &dto.CLIMeta{Service: "fluxrpc", Endpoint: "getInflationRate"})
	},
}

var rpcStakingInflationRewardCmd = &cobra.Command{
	Use:   "inflation-reward <addr1,addr2,...>",
	Short: "Get inflation rewards for addresses (getInflationReward)",
	Args:  exactArgsFromUse(),
	Run: func(cmd *cobra.Command, args []string) {
		addresses := strings.Split(args[0], ",")
		epoch, _ := cmd.Flags().GetUint64("epoch")
		commitment, _ := cmd.Flags().GetString("commitment")
		result, err := fluxRPCSvc().GetInflationReward(addresses, epoch, commitment)
		if err != nil {
			output.PrintError(cmd, "RPC_ERROR", err.Error(), &dto.CLIMeta{Service: "fluxrpc"})
			os.Exit(httpclient.ExitCodeForError(err))
		}
		output.PrintSuccess(cmd, result, &dto.CLIMeta{Service: "fluxrpc", Endpoint: "getInflationReward"})
	},
}
