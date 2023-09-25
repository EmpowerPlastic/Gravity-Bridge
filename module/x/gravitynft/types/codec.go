package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

// ModuleCdc is the codec for the module
var ModuleCdc = codec.NewLegacyAmino()

func init() {
	RegisterCodec(ModuleCdc)
}

// RegisterInterfaces registers the interfaces for the proto stuff
// nolint: exhaustruct
func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgSendNFTToCosmosClaim{},
	)

	registry.RegisterInterface(
		"gravitynft.v1beta1.EthereumNFTClaim",
		(*EthereumNFTClaim)(nil),
		&MsgSendNFTToCosmosClaim{},
	)

	// TODO: If we need any gov stuff?
	// registry.RegisterImplementations((*govtypesv1beta1.Content)(nil), &UnhaltBridgeProposal{}, &AirdropProposal{}, &IBCMetadataProposal{})

	// TODO: If we need anything like this
	// registry.RegisterInterface("gravity.v1beta1.EthereumSigned", (*EthereumSigned)(nil), &Valset{}, &OutgoingTxBatch{}, &OutgoingLogicCall{})

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

// RegisterCodec registers concrete types on the Amino codec
// nolint: exhaustruct
func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterInterface((*EthereumNFTClaim)(nil), nil)
	cdc.RegisterConcrete(&MsgSendNFTToCosmosClaim{}, "gravitynft/MsgSendNFTToCosmosClaim", nil)
}
