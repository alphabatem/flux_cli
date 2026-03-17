package output

import (
	"github.com/alphabatem/flux_cli/dto"
	"github.com/spf13/cobra"
)

// Print outputs the CLIResponse in the format specified by the --format flag.
func Print(cmd *cobra.Command, response *dto.CLIResponse) {
	format, _ := cmd.Flags().GetString("format")
	if format == "" {
		format = "json"
	}

	switch format {
	case "table":
		printTable(response)
	default:
		printJSON(response)
	}
}

// PrintSuccess is a convenience function for successful responses.
func PrintSuccess(cmd *cobra.Command, data interface{}, meta *dto.CLIMeta) {
	Print(cmd, &dto.CLIResponse{
		Success: true,
		Data:    data,
		Meta:    meta,
	})
}

// PrintError is a convenience function for error responses.
func PrintError(cmd *cobra.Command, code string, message string, meta *dto.CLIMeta) {
	Print(cmd, &dto.CLIResponse{
		Success: false,
		Error: &dto.CLIError{
			Code:    code,
			Message: message,
		},
		Meta: meta,
	})
}
