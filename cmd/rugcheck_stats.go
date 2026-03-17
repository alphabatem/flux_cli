package cmd

import (
	"os"

	"github.com/alphabatem/flux_cli/dto"
	"github.com/alphabatem/flux_cli/internal/httpclient"
	"github.com/alphabatem/flux_cli/output"
	"github.com/spf13/cobra"
)

func init() {
	rugcheckCmd.AddCommand(rugcheckStatsCmd)
	rugcheckStatsCmd.AddCommand(rugcheckStatsNewCmd)
	rugcheckStatsCmd.AddCommand(rugcheckStatsRecentCmd)
	rugcheckStatsCmd.AddCommand(rugcheckStatsTrendingCmd)
	rugcheckStatsCmd.AddCommand(rugcheckStatsVerifiedCmd)
}

var rugcheckStatsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Token statistics - new, recent, trending, verified",
}

var rugcheckStatsNewCmd = &cobra.Command{
	Use:   "new",
	Short: "Get recently detected tokens",
	Run: func(cmd *cobra.Command, args []string) {
		result, err := rugCheckSvc().GetStatsNew()
		if err != nil {
			output.PrintError(cmd, "API_ERROR", err.Error(), &dto.CLIMeta{Service: "rugcheck"})
			os.Exit(httpclient.ExitCodeForError(err))
		}
		output.PrintSuccess(cmd, result, &dto.CLIMeta{Service: "rugcheck"})
	},
}

var rugcheckStatsRecentCmd = &cobra.Command{
	Use:   "recent",
	Short: "Get most viewed tokens in the past 24 hours",
	Run: func(cmd *cobra.Command, args []string) {
		result, err := rugCheckSvc().GetStatsRecent()
		if err != nil {
			output.PrintError(cmd, "API_ERROR", err.Error(), &dto.CLIMeta{Service: "rugcheck"})
			os.Exit(httpclient.ExitCodeForError(err))
		}
		output.PrintSuccess(cmd, result, &dto.CLIMeta{Service: "rugcheck"})
	},
}

var rugcheckStatsTrendingCmd = &cobra.Command{
	Use:   "trending",
	Short: "Get trending tokens in the past 24 hours",
	Run: func(cmd *cobra.Command, args []string) {
		result, err := rugCheckSvc().GetStatsTrending()
		if err != nil {
			output.PrintError(cmd, "API_ERROR", err.Error(), &dto.CLIMeta{Service: "rugcheck"})
			os.Exit(httpclient.ExitCodeForError(err))
		}
		output.PrintSuccess(cmd, result, &dto.CLIMeta{Service: "rugcheck"})
	},
}

var rugcheckStatsVerifiedCmd = &cobra.Command{
	Use:   "verified",
	Short: "Get recently verified tokens",
	Run: func(cmd *cobra.Command, args []string) {
		result, err := rugCheckSvc().GetStatsVerified()
		if err != nil {
			output.PrintError(cmd, "API_ERROR", err.Error(), &dto.CLIMeta{Service: "rugcheck"})
			os.Exit(httpclient.ExitCodeForError(err))
		}
		output.PrintSuccess(cmd, result, &dto.CLIMeta{Service: "rugcheck"})
	},
}
