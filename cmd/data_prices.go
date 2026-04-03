package cmd

import (
	"os"

	"github.com/alphabatem/flux_cli/dto"
	"github.com/alphabatem/flux_cli/internal/httpclient"
	"github.com/alphabatem/flux_cli/output"
	"github.com/spf13/cobra"
)

func init() {
	dataCmd.AddCommand(dataPricesCmd)
}

var dataPricesCmd = &cobra.Command{
	Use:   "prices <mint1,mint2,...>",
	Short: "Get bulk token prices (max 255 tokens)",
	Args:  exactArgsFromUse(),
	Run: func(cmd *cobra.Command, args []string) {
		result, err := dataStreamSvc().GetPrices(args[0])
		if err != nil {
			output.PrintError(cmd, "API_ERROR", err.Error(), &dto.CLIMeta{Service: "datastream", Endpoint: "/prices"})
			os.Exit(httpclient.ExitCodeForError(err))
		}
		output.PrintSuccess(cmd, result, &dto.CLIMeta{Service: "datastream", Endpoint: "/prices"})
	},
}
