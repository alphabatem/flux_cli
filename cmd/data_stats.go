package cmd

import (
	"os"

	"github.com/alphabatem/flux_cli/dto"
	"github.com/alphabatem/flux_cli/internal/httpclient"
	"github.com/alphabatem/flux_cli/output"
	"github.com/spf13/cobra"
)

func init() {
	dataCmd.AddCommand(dataStatsCmd)
	dataStatsCmd.AddCommand(dataStatsNewCmd)

	dataStatsTopCmd.Flags().Int("limit", 100, "Number of results (max 200)")
	dataStatsTopCmd.Flags().String("interval", "30", "Time interval in minutes")
	dataStatsCmd.AddCommand(dataStatsTopCmd)

	dataStatsTrendingCmd.Flags().Int("limit", 100, "Number of results (max 100)")
	dataStatsTrendingCmd.Flags().String("interval", "30", "Time interval in minutes")
	dataStatsCmd.AddCommand(dataStatsTrendingCmd)

	dataStatsVolumeCmd.Flags().Int("limit", 100, "Number of results (max 200)")
	dataStatsCmd.AddCommand(dataStatsVolumeCmd)
}

var dataStatsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Market statistics - new tokens, top, trending, volume",
}

var dataStatsNewCmd = &cobra.Command{
	Use:   "new",
	Short: "Get newly detected tokens",
	Run: func(cmd *cobra.Command, args []string) {
		result, err := dataStreamSvc().GetStatsNew()
		if err != nil {
			output.PrintError(cmd, "API_ERROR", err.Error(), &dto.CLIMeta{Service: "datastream", Endpoint: "/stats/tokens/new"})
			os.Exit(httpclient.ExitCodeForError(err))
		}
		output.PrintSuccess(cmd, result, &dto.CLIMeta{Service: "datastream", Endpoint: "/stats/tokens/new"})
	},
}

var dataStatsTopCmd = &cobra.Command{
	Use:   "top",
	Short: "Get top tokens by volume",
	Run: func(cmd *cobra.Command, args []string) {
		limit, _ := cmd.Flags().GetInt("limit")
		interval, _ := cmd.Flags().GetString("interval")
		result, err := dataStreamSvc().GetStatsTop(limit, interval)
		if err != nil {
			output.PrintError(cmd, "API_ERROR", err.Error(), &dto.CLIMeta{Service: "datastream"})
			os.Exit(httpclient.ExitCodeForError(err))
		}
		output.PrintSuccess(cmd, result, &dto.CLIMeta{Service: "datastream"})
	},
}

var dataStatsTrendingCmd = &cobra.Command{
	Use:   "trending",
	Short: "Get trending tokens",
	Run: func(cmd *cobra.Command, args []string) {
		limit, _ := cmd.Flags().GetInt("limit")
		interval, _ := cmd.Flags().GetString("interval")
		result, err := dataStreamSvc().GetStatsTrending(limit, interval)
		if err != nil {
			output.PrintError(cmd, "API_ERROR", err.Error(), &dto.CLIMeta{Service: "datastream"})
			os.Exit(httpclient.ExitCodeForError(err))
		}
		output.PrintSuccess(cmd, result, &dto.CLIMeta{Service: "datastream"})
	},
}

var dataStatsVolumeCmd = &cobra.Command{
	Use:   "volume",
	Short: "Get tokens by volume",
	Run: func(cmd *cobra.Command, args []string) {
		limit, _ := cmd.Flags().GetInt("limit")
		result, err := dataStreamSvc().GetStatsVolume(limit)
		if err != nil {
			output.PrintError(cmd, "API_ERROR", err.Error(), &dto.CLIMeta{Service: "datastream"})
			os.Exit(httpclient.ExitCodeForError(err))
		}
		output.PrintSuccess(cmd, result, &dto.CLIMeta{Service: "datastream"})
	},
}
