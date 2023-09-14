package cli

import (
	"github.com/Gravity-Bridge/Gravity-Bridge/module/x/gravitynft/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"
)

// GetQueryCmd bundles all the query subcmds together so they appear under `gravitynf query` or `gravitynft q`
func GetQueryCmd() *cobra.Command {
	gravitynftQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the gravitynft module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	gravitynftQueryCmd.AddCommand([]*cobra.Command{
		// TODO: Add query Cmds
	}...)

	return gravitynftQueryCmd
}