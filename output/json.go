package output

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/alphabatem/flux_cli/dto"
)

func printJSON(response *dto.CLIResponse) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)
	if err := enc.Encode(response); err != nil {
		fmt.Fprintf(os.Stderr, "error encoding JSON: %v\n", err)
		os.Exit(dto.ExitGeneralError)
	}
}
