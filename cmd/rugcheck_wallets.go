package cmd

import (
	"os"

	"github.com/alphabatem/flux_cli/dto"
	"github.com/alphabatem/flux_cli/internal/httpclient"
	"github.com/alphabatem/flux_cli/output"
	"github.com/spf13/cobra"
)

func init() {
	rugcheckCmd.AddCommand(rugcheckWalletCmd)
	rugcheckCmd.AddCommand(rugcheckChainsCmd)
}

var rugcheckWalletCmd = &cobra.Command{
	Use:   "wallet <chain> <address>",
	Short: "Assess wallet risk profile",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		result, err := rugCheckSvc().GetWalletRisk(args[0], args[1])
		if err != nil {
			output.PrintError(cmd, "API_ERROR", err.Error(), &dto.CLIMeta{Service: "rugcheck"})
			os.Exit(httpclient.ExitCodeForError(err))
		}
		output.PrintSuccess(cmd, result, &dto.CLIMeta{Service: "rugcheck"})
	},
}

var rugcheckChainsCmd = &cobra.Command{
	Use:   "chains",
	Short: "List all supported blockchain networks",
	Run: func(cmd *cobra.Command, args []string) {
		result, err := rugCheckSvc().GetChains()
		if err != nil {
			output.PrintError(cmd, "API_ERROR", err.Error(), &dto.CLIMeta{Service: "rugcheck"})
			os.Exit(httpclient.ExitCodeForError(err))
		}
		output.PrintSuccess(cmd, result, &dto.CLIMeta{Service: "rugcheck"})
	},
}
