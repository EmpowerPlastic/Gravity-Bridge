package cli

import (
	"github.com/Gravity-Bridge/Gravity-Bridge/module/x/gravitynft/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/spf13/cobra"
	"strconv"
)

const (
	FlagOrder     = "order"
	FlagClaimType = "claim-type"
	FlagNonce     = "nonce"
	FlagEthHeight = "eth-height"
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
		CmdGetNFTAttestations(),
		GetCmdPendingNFTIbcAutoForwards(),
		CmdGetLastObservedNFTEthBlock(),
		CmdGetLastObservedNFTEthNonce(),
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

// CmdGetNFTAttestations fetches the most recently created Attestations in the store (only the most recent 1000 are available)
// up to an optional limit
func CmdGetNFTAttestations() *cobra.Command {
	short := "Query gravity current and historical attestations (only the most recent 1000 are stored)"
	long := short + "\n\n" + "Optionally provide a limit to reduce the number of attestations returned"
	// nolint: exhaustruct
	cmd := &cobra.Command{
		Use:   "attestations [optional limit]",
		Args:  cobra.MaximumNArgs(1),
		Short: short,
		Long:  long,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			var limit uint64
			// Limit is 0 or whatever the user put in
			if len(args) == 0 || args[0] == "" {
				limit = 0
			} else {
				limit, err = strconv.ParseUint(args[0], 10, 64)
				if err != nil {
					return err
				}
			}
			orderBy, err := cmd.Flags().GetString(FlagOrder)
			if err != nil {
				return err
			}
			claimType, err := cmd.Flags().GetString(FlagClaimType)
			if err != nil {
				return err
			}
			nonce, err := cmd.Flags().GetUint64(FlagNonce)
			if err != nil {
				return err
			}
			height, err := cmd.Flags().GetUint64(FlagEthHeight)
			if err != nil {
				return err
			}

			req := &types.QueryNFTAttestationsRequest{
				Limit:     limit,
				OrderBy:   orderBy,
				ClaimType: claimType,
				Nonce:     nonce,
				Height:    height,
			}
			res, err := queryClient.GetNFTAttestations(cmd.Context(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	// Global flags
	flags.AddQueryFlagsToCmd(cmd)
	// Local flags
	cmd.Flags().String(FlagOrder, "asc", "order attestations by eth block height: set to 'desc' for reverse ordering")
	cmd.Flags().String(FlagClaimType, "", "which types of claims to filter, empty for all or one of: CLAIM_TYPE_SEND_TO_COSMOS, CLAIM_TYPE_BATCH_SEND_TO_ETH, CLAIM_TYPE_ERC20_DEPLOYED, CLAIM_TYPE_LOGIC_CALL_EXECUTED, CLAIM_TYPE_VALSET_UPDATED")
	cmd.Flags().Uint64(FlagNonce, 0, "the exact nonce to find, 0 for any")
	cmd.Flags().Uint64(FlagEthHeight, 0, "the exact ethereum block height an event happened at, 0 for any")

	return cmd
}

// GetCmdPendingNFTIbcAutoForwards fetches the next NFT IBC auto forwards to be executed, up to an optional limit
func GetCmdPendingNFTIbcAutoForwards() *cobra.Command {
	// nolint: exhaustruct
	cmd := &cobra.Command{
		Use:   "pending-nft-ibc-auto-forwards [optional limit]",
		Short: "Query SendNFTToCosmos transactions waiting to be forwarded over IBC",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			var limit uint64 = 0
			if args[0] != "" {
				var err error
				limit, err = strconv.ParseUint(args[0], 10, 0)
				if err != nil {
					return sdkerrors.Wrapf(err, "Unable to parse limit from %v", args[0])
				}
			}

			req := &types.QueryPendingNFTIbcAutoForwardsRequest{Limit: limit}
			res, err := queryClient.GetPendingNFTIbcAutoForwards(cmd.Context(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// CmdGetLastObservedNFTEthBlock fetches the Ethereum block height for the most recent "observed" Attestation, indicating
// the state of Cosmos consensus on the submitted Ethereum events
// nolint: dupl
func CmdGetLastObservedNFTEthBlock() *cobra.Command {
	short := "Query the last observed NFT Ethereum block height"
	long := short + "\n\n" +
		"This value is expected to lag the actual Ethereum block height significantly due to 1. Ethereum Finality and 2. Consensus mirroring the state on Ethereum"
	// nolint: exhaustruct
	cmd := &cobra.Command{
		Use:   "last-observed-nft-eth-block",
		Short: short,
		Long:  long,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			req := &types.QueryLastObservedNFTEthBlockRequest{}
			res, err := queryClient.GetLastObservedNFTEthBlock(cmd.Context(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// CmdGetLastObservedNFTEthNonce fetches the Ethereum event nonce for the most recent "observed" Attestation, indicating
// // the state of Cosmos consensus on the submitted Ethereum events
// nolint: dupl
func CmdGetLastObservedNFTEthNonce() *cobra.Command {
	short := "Query the last observed NFT Ethereum event nonce"
	long := short + "\n\n" +
		"This this is likely to lag the last executed event a little due to 1. Ethereum Finality and 2. Consensus mirroring the Ethereum state"
	// nolint: exhaustruct
	cmd := &cobra.Command{
		Use:   "last-observed-nft-eth-nonce",
		Short: short,
		Long:  long,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			req := &types.QueryLastObservedNFTEthNonceRequest{}
			res, err := queryClient.GetLastObservedNFTEthNonce(cmd.Context(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}