package output

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/alphabatem/flux_cli/dto"
	"github.com/spf13/cobra"
)

func printJSON(cmd *cobra.Command, response *dto.CLIResponse) {
	enc := json.NewEncoder(os.Stdout)
	indent, _ := cmd.Flags().GetBool("json-indent")
	if indent {
		enc.SetIndent("", "  ")
	}
	enc.SetEscapeHTML(false)
	if err := enc.Encode(response); err != nil {
		fmt.Fprintf(os.Stderr, "error encoding JSON: %v\n", err)
		os.Exit(dto.ExitGeneralError)
	}
}
