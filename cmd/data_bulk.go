package cmd

import (
	"os"

	"github.com/alphabatem/flux_cli/dto"
	"github.com/alphabatem/flux_cli/internal/httpclient"
	"github.com/alphabatem/flux_cli/output"
	"github.com/spf13/cobra"
)

func init() {
	dataCmd.AddCommand(dataBulkCmd)

	dataBulkCandlesCmd.Flags().String("interval", "5", "Candle interval")
	dataBulkCandlesCmd.Flags().Int64("from", 0, "Start timestamp")
	dataBulkCandlesCmd.Flags().Int64("to", 0, "End timestamp")
	dataBulkCandlesCmd.Flags().Int("count", 100, "Number of candles")
	dataBulkCmd.AddCommand(dataBulkCandlesCmd)
}

var dataBulkCmd = &cobra.Command{
	Use:   "bulk",
	Short: "Bulk data operations",
}

var dataBulkCandlesCmd = &cobra.Command{
	Use:   "candles <mint1,mint2,...>",
	Short: "Get bulk candle data for multiple tokens",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		interval, _ := cmd.Flags().GetString("interval")
		from, _ := cmd.Flags().GetInt64("from")
		to, _ := cmd.Flags().GetInt64("to")
		count, _ := cmd.Flags().GetInt("count")
		result, err := dataStreamSvc().GetBulkCandles(args[0], interval, from, to, count)
		if err != nil {
			output.PrintError(cmd, "API_ERROR", err.Error(), &dto.CLIMeta{Service: "datastream"})
			os.Exit(httpclient.ExitCodeForError(err))
		}
		output.PrintSuccess(cmd, result, &dto.CLIMeta{Service: "datastream"})
	},
}
