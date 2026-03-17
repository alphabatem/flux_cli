package cmd

import (
	"os"

	"github.com/alphabatem/flux_cli/dto"
	"github.com/alphabatem/flux_cli/output"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configListCmd)
	configCmd.AddCommand(configPathCmd)
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage CLI configuration (API keys, base URLs, output format)",
}

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a configuration value",
	Long: `Set a configuration value using dot-notation keys.

Available keys:
  datastream.api_key    DataStream API key
  datastream.base_url   DataStream base URL
  fluxrpc.api_key       FluxRPC API key
  fluxrpc.region        FluxRPC region: eu or us (default: us)
  fluxrpc.base_url      FluxRPC base URL (overridden by region)
  rugcheck.api_key      RugCheck API key
  rugcheck.base_url     RugCheck base URL
  output.format         Default output format (json or table)`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if err := configSvc().Set(args[0], args[1]); err != nil {
			output.PrintError(cmd, "CONFIG_ERROR", err.Error(), nil)
			os.Exit(dto.ExitConfigError)
		}
		output.PrintSuccess(cmd, map[string]string{
			"key":   args[0],
			"value": args[1],
		}, nil)
	},
}

var configGetCmd = &cobra.Command{
	Use:   "get <key>",
	Short: "Get a configuration value",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		value, err := configSvc().Get(args[0])
		if err != nil {
			output.PrintError(cmd, "CONFIG_ERROR", err.Error(), nil)
			os.Exit(dto.ExitConfigError)
		}
		output.PrintSuccess(cmd, map[string]string{
			"key":   args[0],
			"value": value,
		}, nil)
	},
}

var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all configuration values (API keys redacted)",
	Run: func(cmd *cobra.Command, args []string) {
		output.PrintSuccess(cmd, configSvc().List(), nil)
	},
}

var configPathCmd = &cobra.Command{
	Use:   "path",
	Short: "Print the configuration file path",
	Run: func(cmd *cobra.Command, args []string) {
		output.PrintSuccess(cmd, map[string]string{
			"path": configSvc().ConfigPath(),
		}, nil)
	},
}
