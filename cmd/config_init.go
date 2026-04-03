package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/alphabatem/flux_cli/dto"
	"github.com/alphabatem/flux_cli/output"
	"github.com/spf13/cobra"
)

func init() {
	configCmd.AddCommand(configInitCmd)
	configInitCmd.Flags().Bool("force", false, "Run setup even if config already exists")
}

var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Run interactive first-time configuration",
	Run: func(cmd *cobra.Command, args []string) {
		force, _ := cmd.Flags().GetBool("force")
		if err := runConfigInit(cmd, os.Stdin, os.Stdout, force); err != nil {
			output.PrintError(cmd, "CONFIG_ERROR", err.Error(), nil)
			os.Exit(dto.ExitConfigError)
		}
	},
}

func runConfigInit(cmd *cobra.Command, in io.Reader, out io.Writer, force bool) error {
	svc := configSvc()
	if svc.Exists() && !force {
		fmt.Fprintln(out, "Config already exists. Re-run with --force to update it interactively.")
		return nil
	}

	cfg := *svc.Config()
	reader := bufio.NewReader(in)

	fmt.Fprintln(out, "Flux CLI setup")
	fmt.Fprintln(out, "Press Enter to skip any API key.")

	cfg.DataStream.APIKey = promptLine(reader, out, "DataStream API key", cfg.DataStream.APIKey, true)
	cfg.FluxRPC.APIKey = promptLine(reader, out, "FluxRPC API key", cfg.FluxRPC.APIKey, true)
	cfg.FluxRPC.Region = promptChoice(reader, out, "FluxRPC region", cfg.FluxRPC.Region, []string{"us", "eu"})
	cfg.RugCheck.APIKey = promptLine(reader, out, "RugCheck API key", cfg.RugCheck.APIKey, true)
	cfg.Output.Format = promptChoice(reader, out, "Default output format", cfg.Output.Format, []string{"json", "table"})

	if err := svc.SaveConfig(cfg); err != nil {
		return err
	}

	fmt.Fprintf(out, "Saved config to %s\n", svc.ConfigPath())
	return nil
}

func promptLine(reader *bufio.Reader, out io.Writer, label, current string, secret bool) string {
	suffix := ""
	if current != "" {
		if secret {
			suffix = " [current: set]"
		} else {
			suffix = fmt.Sprintf(" [current: %s]", current)
		}
	}
	fmt.Fprintf(out, "%s%s: ", label, suffix)
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)
	if line == "" {
		return current
	}
	return line
}

func promptChoice(reader *bufio.Reader, out io.Writer, label, current string, choices []string) string {
	if current == "" && len(choices) > 0 {
		current = choices[0]
	}

	for {
		fmt.Fprintf(out, "%s [%s] (default: %s): ", label, strings.Join(choices, "/"), current)
		line, _ := reader.ReadString('\n')
		line = strings.TrimSpace(strings.ToLower(line))
		if line == "" {
			return current
		}
		for _, choice := range choices {
			if line == choice {
				return line
			}
		}
		fmt.Fprintf(out, "Invalid value %q. Choose one of: %s\n", line, strings.Join(choices, ", "))
	}
}
