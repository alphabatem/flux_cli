package cmd

import (
	"os"

	"github.com/alphabatem/flux_cli/dto"
	"github.com/alphabatem/flux_cli/internal/httpclient"
	"github.com/alphabatem/flux_cli/output"
	"github.com/spf13/cobra"
)

func init() {
	dataCmd.AddCommand(dataTokensCmd)

	dataTokensListCmd.Flags().Int("limit", 100, "Number of tokens to return (max 100)")
	dataTokensListCmd.Flags().String("cursor", "", "Pagination cursor")
	dataTokensCmd.AddCommand(dataTokensListCmd)

	dataTokensCmd.AddCommand(dataTokensGetCmd)
	dataTokensCmd.AddCommand(dataTokensDetailsCmd)

	dataTokensCandlesCmd.Flags().String("interval", "1", "Candle interval")
	dataTokensCandlesCmd.Flags().Int64("from", 0, "Start timestamp")
	dataTokensCandlesCmd.Flags().Int64("to", 0, "End timestamp")
	dataTokensCandlesCmd.Flags().Int("count", 10, "Number of candles")
	dataTokensCmd.AddCommand(dataTokensCandlesCmd)

	dataTokensTradersCmd.Flags().Int("limit", 100, "Number of traders (max 100)")
	dataTokensTradersCmd.Flags().Int("page", 0, "Page number")
	dataTokensCmd.AddCommand(dataTokensTradersCmd)

	dataTokensTradesCmd.Flags().Int("limit", 100, "Number of trades (max 100)")
	dataTokensTradesCmd.Flags().Int64("from", 0, "Start timestamp")
	dataTokensTradesCmd.Flags().Int64("to", 0, "End timestamp")
	dataTokensCmd.AddCommand(dataTokensTradesCmd)

	dataTokensCmd.AddCommand(dataTokensPriceCmd)
	dataTokensCmd.AddCommand(dataTokensMcapCmd)
	dataTokensCmd.AddCommand(dataTokensStatsCmd)
}

var dataTokensCmd = &cobra.Command{
	Use:   "tokens",
	Short: "Token data - details, candles, traders, trades, price, stats",
}

var dataTokensListCmd = &cobra.Command{
	Use:   "list",
	Short: "List tokens",
	Run: func(cmd *cobra.Command, args []string) {
		limit, _ := cmd.Flags().GetInt("limit")
		cursor, _ := cmd.Flags().GetString("cursor")
		result, err := dataStreamSvc().GetTokens(cursor, limit)
		if err != nil {
			output.PrintError(cmd, "API_ERROR", err.Error(), &dto.CLIMeta{Service: "datastream", Endpoint: "/tokens"})
			os.Exit(httpclient.ExitCodeForError(err))
		}
		output.PrintSuccess(cmd, result, &dto.CLIMeta{Service: "datastream", Endpoint: "/tokens"})
	},
}

var dataTokensGetCmd = &cobra.Command{
	Use:   "get <mint>",
	Short: "Get token by mint address",
	Args:  exactArgsFromUse(),
	Run: func(cmd *cobra.Command, args []string) {
		result, err := dataStreamSvc().GetToken(args[0])
		if err != nil {
			output.PrintError(cmd, "API_ERROR", err.Error(), &dto.CLIMeta{Service: "datastream"})
			os.Exit(httpclient.ExitCodeForError(err))
		}
		output.PrintSuccess(cmd, result, &dto.CLIMeta{Service: "datastream"})
	},
}

var dataTokensDetailsCmd = &cobra.Command{
	Use:   "details <mint>",
	Short: "Get extended token details with price and market cap",
	Args:  exactArgsFromUse(),
	Run: func(cmd *cobra.Command, args []string) {
		result, err := dataStreamSvc().GetTokenDetails(args[0])
		if err != nil {
			output.PrintError(cmd, "API_ERROR", err.Error(), &dto.CLIMeta{Service: "datastream"})
			os.Exit(httpclient.ExitCodeForError(err))
		}
		output.PrintSuccess(cmd, result, &dto.CLIMeta{Service: "datastream"})
	},
}

var dataTokensCandlesCmd = &cobra.Command{
	Use:   "candles <mint>",
	Short: "Get OHLCV candle data for a token",
	Args:  exactArgsFromUse(),
	Run: func(cmd *cobra.Command, args []string) {
		interval, _ := cmd.Flags().GetString("interval")
		from, _ := cmd.Flags().GetInt64("from")
		to, _ := cmd.Flags().GetInt64("to")
		count, _ := cmd.Flags().GetInt("count")
		result, err := dataStreamSvc().GetTokenCandles(args[0], interval, from, to, count)
		if err != nil {
			output.PrintError(cmd, "API_ERROR", err.Error(), &dto.CLIMeta{Service: "datastream"})
			os.Exit(httpclient.ExitCodeForError(err))
		}
		output.PrintSuccess(cmd, result, &dto.CLIMeta{Service: "datastream"})
	},
}

var dataTokensTradersCmd = &cobra.Command{
	Use:   "traders <mint>",
	Short: "Get trader stats for a token",
	Args:  exactArgsFromUse(),
	Run: func(cmd *cobra.Command, args []string) {
		limit, _ := cmd.Flags().GetInt("limit")
		page, _ := cmd.Flags().GetInt("page")
		result, err := dataStreamSvc().GetTokenTraders(args[0], limit, page, false)
		if err != nil {
			output.PrintError(cmd, "API_ERROR", err.Error(), &dto.CLIMeta{Service: "datastream"})
			os.Exit(httpclient.ExitCodeForError(err))
		}
		output.PrintSuccess(cmd, result, &dto.CLIMeta{Service: "datastream"})
	},
}

var dataTokensTradesCmd = &cobra.Command{
	Use:   "trades <mint>",
	Short: "Get trade history for a token",
	Args:  exactArgsFromUse(),
	Run: func(cmd *cobra.Command, args []string) {
		limit, _ := cmd.Flags().GetInt("limit")
		from, _ := cmd.Flags().GetInt64("from")
		to, _ := cmd.Flags().GetInt64("to")
		result, err := dataStreamSvc().GetTokenTrades(args[0], limit, from, to)
		if err != nil {
			output.PrintError(cmd, "API_ERROR", err.Error(), &dto.CLIMeta{Service: "datastream"})
			os.Exit(httpclient.ExitCodeForError(err))
		}
		output.PrintSuccess(cmd, result, &dto.CLIMeta{Service: "datastream"})
	},
}

var dataTokensPriceCmd = &cobra.Command{
	Use:   "price <mint>",
	Short: "Get current price for a token",
	Args:  exactArgsFromUse(),
	Run: func(cmd *cobra.Command, args []string) {
		result, err := dataStreamSvc().GetTokenPrice(args[0])
		if err != nil {
			output.PrintError(cmd, "API_ERROR", err.Error(), &dto.CLIMeta{Service: "datastream"})
			os.Exit(httpclient.ExitCodeForError(err))
		}
		output.PrintSuccess(cmd, result, &dto.CLIMeta{Service: "datastream"})
	},
}

var dataTokensMcapCmd = &cobra.Command{
	Use:   "mcap <mint>",
	Short: "Get market cap for a token",
	Args:  exactArgsFromUse(),
	Run: func(cmd *cobra.Command, args []string) {
		result, err := dataStreamSvc().GetTokenMcap(args[0])
		if err != nil {
			output.PrintError(cmd, "API_ERROR", err.Error(), &dto.CLIMeta{Service: "datastream"})
			os.Exit(httpclient.ExitCodeForError(err))
		}
		output.PrintSuccess(cmd, result, &dto.CLIMeta{Service: "datastream"})
	},
}

var dataTokensStatsCmd = &cobra.Command{
	Use:   "stats <mint>",
	Short: "Get volume change stats for a token",
	Args:  exactArgsFromUse(),
	Run: func(cmd *cobra.Command, args []string) {
		result, err := dataStreamSvc().GetTokenStats(args[0])
		if err != nil {
			output.PrintError(cmd, "API_ERROR", err.Error(), &dto.CLIMeta{Service: "datastream"})
			os.Exit(httpclient.ExitCodeForError(err))
		}
		output.PrintSuccess(cmd, result, &dto.CLIMeta{Service: "datastream"})
	},
}
