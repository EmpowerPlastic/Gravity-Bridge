package cli

import (
	"github.com/Gravity-Bridge/Gravity-Bridge/module/x/gravitynft/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
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
		GetCmdQueryParams(),
	}...)

	return gravitynftQueryCmd
}

// GetCmdQueryParams fetches the current Gravity module params
func GetCmdQueryParams() *cobra.Command {
	// nolint: exhaustruct
	cmd := &cobra.Command{
		Use:   "params",
		Args:  cobra.NoArgs,
		Short: "Query gravitynft params",
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.Params(cmd.Context(), &types.QueryParamsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(&res.Params)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
