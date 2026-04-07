package cmd

import (
	"os"

	"github.com/alphabatem/flux_cli/dto"
	"github.com/alphabatem/flux_cli/internal/httpclient"
	"github.com/alphabatem/flux_cli/output"
	"github.com/spf13/cobra"
)

func init() {
	dataCmd.AddCommand(dataHoldersCmd)

	dataHoldersListCmd.Flags().Int("limit", 100, "Number of holders (max 100)")
	dataHoldersListCmd.Flags().Int("page", 0, "Page number")
	dataHoldersCmd.AddCommand(dataHoldersListCmd)

	dataHoldersTopCmd.Flags().Int("limit", 100, "Number of top holders (max 1000)")
	dataHoldersCmd.AddCommand(dataHoldersTopCmd)

	dataHoldersCmd.AddCommand(dataHoldersDistCmd)
	dataHoldersCmd.AddCommand(dataHoldersCountCmd)
}

var dataHoldersCmd = &cobra.Command{
	Use:   "holders",
	Short: "Holder data - list, top, distribution, count",
}

var dataHoldersListCmd = &cobra.Command{
	Use:   "list <mint>",
	Short: "Get holders for a token",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		limit, _ := cmd.Flags().GetInt("limit")
		page, _ := cmd.Flags().GetInt("page")
		result, err := dataStreamSvc().GetTokenHolders(args[0], limit, page)
		if err != nil {
			output.PrintError(cmd, "API_ERROR", err.Error(), &dto.CLIMeta{Service: "datastream"})
			os.Exit(httpclient.ExitCodeForError(err))
		}
		output.PrintSuccess(cmd, result, &dto.CLIMeta{Service: "datastream"})
	},
}

var dataHoldersTopCmd = &cobra.Command{
	Use:   "top <mint>",
	Short: "Get top holders for a token sorted by balance",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		limit, _ := cmd.Flags().GetInt("limit")
		result, err := dataStreamSvc().GetTokenHoldersTop(args[0], limit)
		if err != nil {
			output.PrintError(cmd, "API_ERROR", err.Error(), &dto.CLIMeta{Service: "datastream"})
			os.Exit(httpclient.ExitCodeForError(err))
		}
		output.PrintSuccess(cmd, result, &dto.CLIMeta{Service: "datastream"})
	},
}

var dataHoldersDistCmd = &cobra.Command{
	Use:   "distribution <mint>",
	Short: "Get holder distribution for a token",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		result, err := dataStreamSvc().GetTokenHoldersDistribution(args[0])
		if err != nil {
			output.PrintError(cmd, "API_ERROR", err.Error(), &dto.CLIMeta{Service: "datastream"})
			os.Exit(httpclient.ExitCodeForError(err))
		}
		output.PrintSuccess(cmd, result, &dto.CLIMeta{Service: "datastream"})
	},
}

var dataHoldersCountCmd = &cobra.Command{
	Use:   "count <mint>",
	Short: "Get holder count for a token",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		result, err := dataStreamSvc().GetTokenHoldersCount(args[0])
		if err != nil {
			output.PrintError(cmd, "API_ERROR", err.Error(), &dto.CLIMeta{Service: "datastream"})
			os.Exit(httpclient.ExitCodeForError(err))
		}
		output.PrintSuccess(cmd, result, &dto.CLIMeta{Service: "datastream"})
	},
}
