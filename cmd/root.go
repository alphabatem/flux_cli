package cmd

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/alphabatem/flux_cli/dto"
	"github.com/alphabatem/flux_cli/pkg/conte
	"github.com/alphabatem/flux_cli/services"
	"github.com/spf13/cobra"
)

var ctx *context.Context

var rootCmd = &cobra.Command{
	Use:   "flux",
	Short: "FluxBeam CLI - unified access to RPC, DataStream, and RugCheck for AI agents",
	Long: `Flux CLI provides AI agents and developers with unified, non-interactive access
to FluxBeam's product suite: RPC (Solana JSON-RPC), DataStream (token & trader data),
and RugCheck (token validation & security reports).

All output defaults to JSON for machine consumption. Use --format table for human-readable output.`,
	SilenceUsage:  true,
	SilenceErrors: true,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, `{"success":false,"error":{"code":"CLI_ERROR","message":"%s"}}%s`, err.Error(), "\n")
		os.Exit(dto.ExitGeneralError)
	}
}

func init() {
	cobra.OnInitialize(initContext)

	rootCmd.PersistentFlags().String("format", "json", "Output format: json or table")
	rootCmd.PersistentFlags().String("datastream-api-key", "", "DataStream API key (overrides config/env)")
	rootCmd.PersistentFlags().String("fluxrpc-api-key", "", "FluxRPC API key (overrides config/env)")
	rootCmd.PersistentFlags().String("fluxrpc-region", "", "FluxRPC region: eu or us (overrides config/env)")
	rootCmd.PersistentFlags().String("rugcheck-api-key", "", "RugCheck API key (overrides config/env)")
}

func initContext() {
	// Suppress harness log output for clean JSON output
	log.SetOutput(io.Discard)

	var err error
	ctx, err = context.NewCtx(
		&services.ConfigService{},
		&services.DataStreamService{},
		&services.FluxRPCService{},
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
