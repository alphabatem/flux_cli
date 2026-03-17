package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(dataCmd)
}

var dataCmd = &cobra.Command{
	Use:   "data",
	Short: "FluxBeam DataStream - token prices, stats, traders, and market data",
}
