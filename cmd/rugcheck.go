package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(rugcheckCmd)
}

var rugcheckCmd = &cobra.Command{
	Use:   "rugcheck",
	Short: "RugCheck - token validation, security reports, and wallet risk assessment",
}
