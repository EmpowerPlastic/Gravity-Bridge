package cli

import (
	sdkerrors "cosmossdk.io/errors"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
	"strconv"

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
		CmdSendNFTToCosmosClaim(),
		CmdExecutePendingNFTIbcAutoForwards(),
	}...)

	return gravitynftTxCmd
}

// TODO: REMOVE THIS. This is here only to test this locally without having a full orchestrator setup going
func CmdSendNFTToCosmosClaim() *cobra.Command {
	// nolint: exhaustruct
	cmd := &cobra.Command{
		Use:   "send-nft-to-cosmos-claim [event_nonce] [eth_block_height] [token_contract] [token_id] [token_uri] [ethereum_sender] [cosmos_receiver] [orchestrator]",
		Short: "Make a MsgSendNFTToCosmosClaim",
		Args:  cobra.ExactArgs(8),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			sender := cliCtx.GetFromAddress()
			if sender.String() == "" {
				return fmt.Errorf("from address must be specified")
			}

			eventNonce, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}
			ethBlockHeight, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}
			tokenContract := args[2]
			tokenId := args[3]
			tokenUri := args[4]
			ethereumSender := args[5]
			cosmosReceiver := args[6]
			orchestrator := args[7]

			msg := types.MsgSendNFTToCosmosClaim{
				EventNonce: eventNonce,
				EthBlockHeight: ethBlockHeight,
				TokenContract: tokenContract,
				TokenId: tokenId,
				TokenUri: tokenUri,
				EthereumSender: ethereumSender,
				CosmosReceiver: cosmosReceiver,
				Orchestrator: orchestrator,
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			// Send it
			return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), &msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// CmdExecutePendingNFTIbcAutoForwards Executes a number of queued IBC Auto Forwards. When users perform a Send NFT to Cosmos
// with a registered foreign address prefix (e.g. canto1... cre1...), their funds will be locked in the gravitynft module
// until their pending forward is executed. This will send the NFT to the equivalent gravity-prefixed account and then
// immediately create an NFT IBC transfer to the destination chain to the original foreign account. If there is an IBC
// failure, the funds will be deposited on the gravitynft-prefixed account.
func CmdExecutePendingNFTIbcAutoForwards() *cobra.Command {
	// nolint: exhaustruct
	cmd := &cobra.Command{
		Use:   "execute-pending-nft-ibc-auto-forwards [forwards-to-execute]",
		Short: "Executes a given number of IBC Auto-Forwards",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			sender := cliCtx.GetFromAddress()
			if sender.String() == "" {
				return fmt.Errorf("from address must be specified")
			}
			forwardsToClear, err := strconv.ParseUint(args[0], 10, 0)
			if err != nil {
				return sdkerrors.Wrap(err, "Unable to parse forwards-to-execute as an non-negative integer")
			}
			msg := types.MsgExecuteIbcNFTAutoForwards{
				ForwardsToClear: forwardsToClear,
				Executor:        cliCtx.GetFromAddress().String(),
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			// Send it
			return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), &msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}