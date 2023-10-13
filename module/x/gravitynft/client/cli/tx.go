package cli

import (
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