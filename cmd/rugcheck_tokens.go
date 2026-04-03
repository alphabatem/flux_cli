package cmd

import (
	"os"

	"github.com/alphabatem/flux_cli/dto"
	"github.com/alphabatem/flux_cli/internal/httpclient"
	"github.com/alphabatem/flux_cli/output"
	"github.com/spf13/cobra"
)

func init() {
	rugcheckCmd.AddCommand(rugcheckReportCmd)
	rugcheckCmd.AddCommand(rugcheckSummaryCmd)

	rugcheckScanCmd.Flags().Bool("include-dexscreener", false, "Include DexScreener data")
	rugcheckCmd.AddCommand(rugcheckScanCmd)

	rugcheckSearchCmd.Flags().String("chain", "", "Filter by blockchain")
	rugcheckSearchCmd.Flags().Int("page-size", 0, "Results per page")
	rugcheckCmd.AddCommand(rugcheckSearchCmd)
}

var rugcheckReportCmd = &cobra.Command{
	Use:   "report <mint>",
	Short: "Get detailed security report for a token",
	Args:  exactArgsFromUse(),
	Run: func(cmd *cobra.Command, args []string) {
		result, err := rugCheckSvc().GetReport(args[0])
		if err != nil {
			output.PrintError(cmd, "API_ERROR", err.Error(), &dto.CLIMeta{Service: "rugcheck"})
			os.Exit(httpclient.ExitCodeForError(err))
		}
		output.PrintSuccess(cmd, result, &dto.CLIMeta{Service: "rugcheck"})
	},
}

var rugcheckSummaryCmd = &cobra.Command{
	Use:   "summary <mint>",
	Short: "Get summary security report for a token",
	Args:  exactArgsFromUse(),
	Run: func(cmd *cobra.Command, args []string) {
		result, err := rugCheckSvc().GetSummary(args[0])
		if err != nil {
			output.PrintError(cmd, "API_ERROR", err.Error(), &dto.CLIMeta{Service: "rugcheck"})
			os.Exit(httpclient.ExitCodeForError(err))
		}
		output.PrintSuccess(cmd, result, &dto.CLIMeta{Service: "rugcheck"})
	},
}

var rugcheckScanCmd = &cobra.Command{
	Use:   "scan <chain> <address>",
	Short: "Real-time token security scan",
	Args:  exactArgsFromUse(),
	Run: func(cmd *cobra.Command, args []string) {
		includeDex, _ := cmd.Flags().GetBool("include-dexscreener")
		result, err := rugCheckSvc().Scan(args[0], args[1], includeDex)
		if err != nil {
			output.PrintError(cmd, "API_ERROR", err.Error(), &dto.CLIMeta{Service: "rugcheck"})
			os.Exit(httpclient.ExitCodeForError(err))
		}
		output.PrintSuccess(cmd, result, &dto.CLIMeta{Service: "rugcheck"})
	},
}

var rugcheckSearchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search tokens by name or symbol",
	Args:  exactArgsFromUse(),
	Run: func(cmd *cobra.Command, args []string) {
		chain, _ := cmd.Flags().GetString("chain")
		pageSize, _ := cmd.Flags().GetInt("page-size")
		result, err := rugCheckSvc().Search(args[0], chain, pageSize)
		if err != nil {
			output.PrintError(cmd, "API_ERROR", err.Error(), &dto.CLIMeta{Service: "rugcheck"})
			os.Exit(httpclient.ExitCodeForError(err))
		}
		output.PrintSuccess(cmd, result, &dto.CLIMeta{Service: "rugcheck"})
	},
}
