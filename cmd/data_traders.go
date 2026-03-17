package cmd

import (
	"os"

	"github.com/alphabatem/flux_cli/dto"
	"github.com/alphabatem/flux_cli/internal/httpclient"
	"github.com/alphabatem/flux_cli/output"
	"github.com/spf13/cobra"
)

func init() {
	dataCmd.AddCommand(dataTradersCmd)

	dataTradersTopCmd.Flags().Int("limit", 100, "Number of results (max 100)")
	dataTradersTopCmd.Flags().Int64("from", 0, "Start timestamp")
	dataTradersTopCmd.Flags().Int64("to", 0, "End timestamp")
	dataTradersCmd.AddCommand(dataTradersTopCmd)

	dataTradersCmd.AddCommand(dataTradersDetailCmd)

	dataTradersTokensCmd.Flags().Int("limit", 100, "Number of results (max 100)")
	dataTradersCmd.AddCommand(dataTradersTokensCmd)

	dataTradersPnLCmd.Flags().Int("limit", 100, "Number of results")
	dataTradersCmd.AddCommand(dataTradersPnLCmd)

	dataTradersTradesCmd.Flags().Int("limit", 100, "Number of results (max 100)")
	dataTradersTradesCmd.Flags().Int64("from", 0, "Start timestamp")
	dataTradersTradesCmd.Flags().Int64("to", 0, "End timestamp")
	dataTradersCmd.AddCommand(dataTradersTradesCmd)
}

var dataTradersCmd = &cobra.Command{
	Use:   "traders",
	Short: "Trader data - top traders, details, holdings, PnL, trades",
}

var dataTradersTopCmd = &cobra.Command{
	Use:   "top",
	Short: "Get top traders",
	Run: func(cmd *cobra.Command, args []string) {
		limit, _ := cmd.Flags().GetInt("limit")
		from, _ := cmd.Flags().GetInt64("from")
		to, _ := cmd.Flags().GetInt64("to")
		result, err := dataStreamSvc().GetTradersTop(limit, from, to)
		if err != nil {
			output.PrintError(cmd, "API_ERROR", err.Error(), &dto.CLIMeta{Service: "datastream"})
			os.Exit(httpclient.ExitCodeForError(err))
		}
		output.PrintSuccess(cmd, result, &dto.CLIMeta{Service: "datastream"})
	},
}

var dataTradersDetailCmd = &cobra.Command{
	Use:   "detail <wallet>",
	Short: "Get trader details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		result, err := dataStreamSvc().GetTraderDetail(args[0])
		if err != nil {
			output.PrintError(cmd, "API_ERROR", err.Error(), &dto.CLIMeta{Service: "datastream"})
			os.Exit(httpclient.ExitCodeForError(err))
		}
		output.PrintSuccess(cmd, result, &dto.CLIMeta{Service: "datastream"})
	},
}

var dataTradersTokensCmd = &cobra.Command{
	Use:   "tokens <wallet>",
	Short: "Get token holdings for a trader",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		limit, _ := cmd.Flags().GetInt("limit")
		result, err := dataStreamSvc().GetTraderTokens(args[0], limit)
		if err != nil {
			output.PrintError(cmd, "API_ERROR", err.Error(), &dto.CLIMeta{Service: "datastream"})
			os.Exit(httpclient.ExitCodeForError(err))
		}
		output.PrintSuccess(cmd, result, &dto.CLIMeta{Service: "datastream"})
	},
}

var dataTradersPnLCmd = &cobra.Command{
	Use:   "pnl <wallet>",
	Short: "Get PnL summary for a trader",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		limit, _ := cmd.Flags().GetInt("limit")
		result, err := dataStreamSvc().GetTraderPnL(args[0], limit)
		if err != nil {
			output.PrintError(cmd, "API_ERROR", err.Error(), &dto.CLIMeta{Service: "datastream"})
			os.Exit(httpclient.ExitCodeForError(err))
		}
		output.PrintSuccess(cmd, result, &dto.CLIMeta{Service: "datastream"})
	},
}

var dataTradersTradesCmd = &cobra.Command{
	Use:   "trades <wallet>",
	Short: "Get trade history for a trader",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		limit, _ := cmd.Flags().GetInt("limit")
		from, _ := cmd.Flags().GetInt64("from")
		to, _ := cmd.Flags().GetInt64("to")
		result, err := dataStreamSvc().GetTraderTrades(args[0], limit, from, to)
		if err != nil {
			output.PrintError(cmd, "API_ERROR", err.Error(), &dto.CLIMeta{Service: "datastream"})
			os.Exit(httpclient.ExitCodeForError(err))
		}
		output.PrintSuccess(cmd, result, &dto.CLIMeta{Service: "datastream"})
	},
}
