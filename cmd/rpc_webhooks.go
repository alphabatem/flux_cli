package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rpcCmd.AddCommand(rpcWebhooksCmd)

	// TODO: Implement webhook CRUD when FluxRPC webhook API is finalized
	// rpcWebhooksCmd.AddCommand(rpcWebhooksListCmd)
	// rpcWebhooksCmd.AddCommand(rpcWebhooksGetCmd)
	// rpcWebhooksCmd.AddCommand(rpcWebhooksCreateCmd)
	// rpcWebhooksCmd.AddCommand(rpcWebhooksUpdateCmd)
	// rpcWebhooksCmd.AddCommand(rpcWebhooksDeleteCmd)
}

var rpcWebhooksCmd = &cobra.Command{
	Use:   "webhooks",
	Short: "FluxRPC webhook management (coming soon)",
}
