package cmd

import (
	stdcontext "context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/alphabatem/flux_cli/dto"
	ctxpkg "github.com/alphabatem/flux_cli/pkg/context"
	"github.com/alphabatem/flux_cli/services"
	"github.com/spf13/cobra"
)

var ctx *ctxpkg.Context

var rootCmd = &cobra.Command{
	Use:   "flux",
	Short: "FluxBeam CLI - unified access to RPC, DataStream, and RugCheck for AI agents",
	Long: `Flux CLI provides AI agents and developers with unified, non-interactive access
to FluxBeam's product suite: RPC (Solana JSON-RPC), DataStream (token & trader data),
and RugCheck (token validation & security reports).

All output defaults to JSON for machine consumption. Use --format table for human-readable output.`,
	SilenceUsage:  true,
	SilenceErrors: true,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		maybeRunFirstTimeSetup(cmd)
		checkForUpdate()
	},
}

func Execute() {
	signalCtx, stop := signal.NotifyContext(stdcontext.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := rootCmd.ExecuteContext(signalCtx); err != nil {
		code := "CLI_ERROR"
		exitCode := dto.ExitGeneralError

		var usageErr *usageError
		if errors.As(err, &usageErr) {
			code = "USAGE_ERROR"
			exitCode = dto.ExitUsageError
		}

		resp := dto.CLIResponse{
			Success: false,
			Error: &dto.CLIError{
				Code:    code,
				Message: err.Error(),
			},
		}
		enc := json.NewEncoder(os.Stderr)
		enc.SetEscapeHTML(false)
		if encodeErr := enc.Encode(resp); encodeErr != nil {
			fmt.Fprintf(os.Stderr, `{"success":false,"error":{"code":"%s","message":"%s"}}%s`, code, err.Error(), "\n")
		}
		os.Exit(exitCode)
	}
}

func init() {
	cobra.OnInitialize(initContext)

	rootCmd.PersistentFlags().String("format", "json", "Output format: json or table")
	rootCmd.PersistentFlags().String("datastream-api-key", "", "DataStream API key (overrides config/env)")
	rootCmd.PersistentFlags().String("fluxrpc-api-key", "", "FluxRPC API key (overrides config/env)")
	rootCmd.PersistentFlags().String("fluxrpc-region", "", "FluxRPC region: eu or us (overrides config/env)")
	rootCmd.PersistentFlags().String("rugcheck-api-key", "", "RugCheck API key (overrides config/env)")

	defaultHelpFunc := rootCmd.HelpFunc()
	rootCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		updateHiddenConfigFlags(rootCmd, loadConfiguredHelpState())
		defaultHelpFunc(cmd, args)
	})
}

func initContext() {
	// Suppress harness log output for clean JSON output
	log.SetOutput(io.Discard)

	var err error
	ctx, err = ctxpkg.NewCtx(
		&services.ConfigService{},
		&services.DataStreamService{},
		&services.FluxRPCService{},
		&services.YellowstoneService{},
		&services.RugCheckService{},
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create context: %v\n", err)
		os.Exit(dto.ExitGeneralError)
	}

	if err := ctx.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to start services: %v\n", err)
		os.Exit(dto.ExitGeneralError)
	}

	// Override API keys from CLI flags (highest priority)
	cfgSvc := ctx.Service(services.CONFIG_SVC).(*services.ConfigService)

	if key, _ := rootCmd.Flags().GetString("datastream-api-key"); key != "" {
		cfgSvc.Config().DataStream.APIKey = key
	}
	if key, _ := rootCmd.Flags().GetString("fluxrpc-api-key"); key != "" {
		cfgSvc.Config().FluxRPC.APIKey = key
	}
	if region, _ := rootCmd.Flags().GetString("fluxrpc-region"); region != "" {
		cfgSvc.Config().FluxRPC.Region = region
	}
	if key, _ := rootCmd.Flags().GetString("rugcheck-api-key"); key != "" {
		cfgSvc.Config().RugCheck.APIKey = key
	}

	updateHiddenConfigFlags(rootCmd, cfgSvc.Config())
}

func updateHiddenConfigFlags(cmd *cobra.Command, cfg *dto.Config) {
	setFlagHidden(cmd, "datastream-api-key", cfg != nil && cfg.DataStream.APIKey != "")
	setFlagHidden(cmd, "fluxrpc-api-key", cfg != nil && cfg.FluxRPC.APIKey != "")
	setFlagHidden(cmd, "rugcheck-api-key", cfg != nil && cfg.RugCheck.APIKey != "")
}

func maybeRunFirstTimeSetup(cmd *cobra.Command) {
	if shouldSkipAutoSetup(cmd) || !isInteractiveTerminal() {
		return
	}

	svc := configSvc()
	if svc.Exists() {
		return
	}

	fmt.Fprintln(os.Stdout, "No Flux CLI config found. Starting first-time setup.")
	if err := runConfigInit(cmd, os.Stdin, os.Stdout, false); err != nil {
		fmt.Fprintf(os.Stderr, "First-time setup failed: %v\n", err)
		os.Exit(dto.ExitConfigError)
	}
	updateHiddenConfigFlags(cmd.Root(), configSvc().Config())
}

func shouldSkipAutoSetup(cmd *cobra.Command) bool {
	path := cmd.CommandPath()
	if path == "flux" {
		return true
	}
	for _, prefix := range []string{"flux config", "flux version", "flux update"} {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}
	return false
}

func isInteractiveTerminal() bool {
	stdinInfo, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	stdoutInfo, err := os.Stdout.Stat()
	if err != nil {
		return false
	}

	return (stdinInfo.Mode()&os.ModeCharDevice) != 0 && (stdoutInfo.Mode()&os.ModeCharDevice) != 0
}

func setFlagHidden(cmd *cobra.Command, name string, hidden bool) {
	flag := cmd.Root().PersistentFlags().Lookup(name)
	if flag == nil {
		return
	}
	flag.Hidden = hidden
}

func loadConfiguredHelpState() *dto.Config {
	cfg := dto.DefaultConfig()

	home, err := os.UserHomeDir()
	if err == nil {
		configPath := filepath.Join(home, ".flux-cli", "config.json")
		if data, readErr := os.ReadFile(configPath); readErr == nil {
			_ = json.Unmarshal(data, &cfg)
		}
	}

	if v := os.Getenv("FLUX_DATASTREAM_API_KEY"); v != "" {
		cfg.DataStream.APIKey = v
	}
	if v := os.Getenv("FLUX_RPC_API_KEY"); v != "" {
		cfg.FluxRPC.APIKey = v
	}
	if v := os.Getenv("FLUX_RUGCHECK_API_KEY"); v != "" {
		cfg.RugCheck.APIKey = v
	}

	return &cfg
}

func configSvc() *services.ConfigService {
	return ctx.Service(services.CONFIG_SVC).(*services.ConfigService)
}

func dataStreamSvc() *services.DataStreamService {
	return ctx.Service(services.DATASTREAM_SVC).(*services.DataStreamService)
}

func fluxRPCSvc() *services.FluxRPCService {
	return ctx.Service(services.FLUXRPC_SVC).(*services.FluxRPCService)
}

func rugCheckSvc() *services.RugCheckService {
	return ctx.Service(services.RUGCHECK_SVC).(*services.RugCheckService)
}

func yellowstoneSvc() *services.YellowstoneService {
	return ctx.Service(services.YELLOWSTONE_SVC).(*services.YellowstoneService)
}
