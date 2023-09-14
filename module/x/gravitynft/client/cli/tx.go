package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"

	"github.com/Gravity-Bridge/Gravity-Bridge/module/x/gravitynft/keeper"
	"github.com/Gravity-Bridge/Gravity-Bridge/module/x/gravitynft/types"
)

// GetTxCmd bundles all the subcmds together so they appear under `gravitynft tx`
func GetTxCmd() *cobra.Command {
	// needed for governance proposal txs in cli case
	// internal check prevents double registration in node case
	keeper.RegisterProposalTypes()

	// nolint: exhaustruct
	gravitynftTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Gravitynft transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	gravitynftTxCmd.AddCommand([]*cobra.Command{
		// TODO: Add tx Cmds
	}...)

	return gravitynftTxCmd
}