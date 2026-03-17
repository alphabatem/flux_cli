package cmd

import (
	"os"

	"github.com/alphabatem/flux_cli/dto"
	"github.com/alphabatem/flux_cli/internal/httpclient"
	"github.com/alphabatem/flux_cli/output"
	"github.com/spf13/cobra"
)

func init() {
	rpcCmd.AddCommand(rpcTokenCmd)

	rpcTokenBalanceCmd.Flags().String("commitment", "", "Commitment level")
	rpcTokenCmd.AddCommand(rpcTokenBalanceCmd)

	rpcTokenAccountsByOwnerCmd.Flags().String("mint", "", "Filter by mint address")
	rpcTokenAccountsByOwnerCmd.Flags().String("program", "", "Filter by token program ID")
	rpcTokenAccountsByOwnerCmd.Flags().String("encoding", "jsonParsed", "Data encoding")
	rpcTokenAccountsByOwnerCmd.Flags().String("commitment", "", "Commitment level")
	rpcTokenCmd.AddCommand(rpcTokenAccountsByOwnerCmd)

	rpcTokenAccountsByDelegateCmd.Flags().String("mint", "", "Filter by mint address")
	rpcTokenAccountsByDelegateCmd.Flags().String("program", "", "Filter by token program ID")
	rpcTokenAccountsByDelegateCmd.Flags().String("encoding", "jsonParsed", "Data encoding")
	rpcTokenAccountsByDelegateCmd.Flags().String("commitment", "", "Commitment level")
	rpcTokenCmd.AddCommand(rpcTokenAccountsByDelegateCmd)

	rpcTokenLargestCmd.Flags().String("commitment", "", "Commitment level")
	rpcTokenCmd.AddCommand(rpcTokenLargestCmd)

	rpcTokenSupplyCmd.Flags().String("commitment", "", "Commitment level")
	rpcTokenCmd.AddCommand(rpcTokenSupplyCmd)
}

var rpcTokenCmd = &cobra.Command{
	Use:   "token",
	Short: "SPL Token methods - balance, accounts, largest holders, supply",
}

var rpcTokenBalanceCmd = &cobra.Command{
	Use:   "balance <tokenAccount>",
	Short: "Get token account balance (getTokenAccountBalance)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		commitment, _ := cmd.Flags().GetString("commitment")
		result, err := fluxRPCSvc().GetTokenAccountBalance(args[0], commitment)
		if err != nil {
			output.PrintError(cmd, "RPC_ERROR", err.Error(), &dto.CLIMeta{Service: "fluxrpc"})
			os.Exit(httpclient.ExitCodeForError(err))
		}
		output.PrintSuccess(cmd, result, &dto.CLIMeta{Service: "fluxrpc", Endpoint: "getTokenAccountBalance"})
	},
}

var rpcTokenAccountsByOwnerCmd = &cobra.Command{
	Use:   "accounts-by-owner <owner>",
	Short: "Get token accounts owned by address (getTokenAccountsByOwner)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		mint, _ := cmd.Flags().GetString("mint")
		program, _ := cmd.Flags().GetString("program")
		encoding, _ := cmd.Flags().GetString("encoding")
		commitment, _ := cmd.Flags().GetString("commitment")

		filterKey, filterValue := "programId", "TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA"
		if mint != "" {
			filterKey, filterValue = "mint", mint
		} else if program != "" {
			filterKey, filterValue = "programId", program
		}

		result, err := fluxRPCSvc().GetTokenAccountsByOwner(args[0], filterKey, filterValue, encoding, commitment)
		if err != nil {
			output.PrintError(cmd, "RPC_ERROR", err.Error(), &dto.CLIMeta{Service: "fluxrpc"})
			os.Exit(httpclient.ExitCodeForError(err))
		}
		output.PrintSuccess(cmd, result, &dto.CLIMeta{Service: "fluxrpc", Endpoint: "getTokenAccountsByOwner"})
	},
}

var rpcTokenAccountsByDelegateCmd = &cobra.Command{
	Use:   "accounts-by-delegate <delegate>",
	Short: "Get token accounts delegated to address (getTokenAccountsByDelegate)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		mint, _ := cmd.Flags().GetString("mint")
		program, _ := cmd.Flags().GetString("program")
		encoding, _ := cmd.Flags().GetString("encoding")
		commitment, _ := cmd.Flags().GetString("commitment")

		filterKey, filterValue := "programId", "TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA"
		if mint != "" {
			filterKey, filterValue = "mint", mint
		} else if program != "" {
			filterKey, filterValue = "programId", program
		}

		result, err := fluxRPCSvc().GetTokenAccountsByDelegate(args[0], filterKey, filterValue, encoding, commitment)
		if err != nil {
			output.PrintError(cmd, "RPC_ERROR", err.Error(), &dto.CLIMeta{Service: "fluxrpc"})
			os.Exit(httpclient.ExitCodeForError(err))
		}
		output.PrintSuccess(cmd, result, &dto.CLIMeta{Service: "fluxrpc", Endpoint: "getTokenAccountsByDelegate"})
	},
}

var rpcTokenLargestCmd = &cobra.Command{
	Use:   "largest <mint>",
	Short: "Get 20 largest token accounts for a mint (getTokenLargestAccounts)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		commitment, _ := cmd.Flags().GetString("commitment")
		result, err := fluxRPCSvc().GetTokenLargestAccounts(args[0], commitment)
		if err != nil {
			output.PrintError(cmd, "RPC_ERROR", err.Error(), &dto.CLIMeta{Service: "fluxrpc"})
			os.Exit(httpclient.ExitCodeForError(err))
		}
		output.PrintSuccess(cmd, result, &dto.CLIMeta{Service: "fluxrpc", Endpoint: "getTokenLargestAccounts"})
	},
}

var rpcTokenSupplyCmd = &cobra.Command{
	Use:   "supply <mint>",
	Short: "Get total token supply (getTokenSupply)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		commitment, _ := cmd.Flags().GetString("commitment")
		result, err := fluxRPCSvc().GetTokenSupply(args[0], commitment)
		if err != nil {
			output.PrintError(cmd, "RPC_ERROR", err.Error(), &dto.CLIMeta{Service: "fluxrpc"})
			os.Exit(httpclient.ExitCodeForError(err))
		}
		output.PrintSuccess(cmd, result, &dto.CLIMeta{Service: "fluxrpc", Endpoint: "getTokenSupply"})
	},
}
