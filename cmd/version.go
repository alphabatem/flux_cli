package cmd

import (
	"github.com/alphabatem/flux_cli/output"
	"github.com/spf13/cobra"
)

var (
	Version   = "dev"
	Commit    = "none"
	BuildDate = "unknown"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the CLI version",
	Run: func(cmd *cobra.Command, args []string) {
		output.PrintSuccess(cmd, map[string]string{
			"version":    Version,
			"commit":     Commit,
			"build_date": BuildDate,
		}, nil)
	},
}
