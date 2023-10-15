package types

import (
	sdkerrors "cosmossdk.io/errors"
	"fmt"
	gravitytypes "github.com/Gravity-Bridge/Gravity-Bridge/module/x/gravity/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/tendermint/tendermint/crypto/tmhash"
)

// nolint: exhaustruct
var (
	_ sdk.Msg = &MsgUpdateParams{}
	_ sdk.Msg = &MsgSendNFTToCosmosClaim{}
	_ sdk.Msg = &MsgERC721DeployedClaim{}
	_ sdk.Msg = &MsgExecuteIbcNFTAutoForwards{}
	_ sdk.Msg = &MsgSendNFTToEth{}
	_ sdk.Msg = &MsgSendNFTToEthClaim{}
	_ sdk.Msg = &MsgCancelSendNFTToEth{}
)

func (msg *MsgUpdateParams) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return sdkerrors.Wrap(errors.ErrInvalidAddress, "authority")
	}

	if err := msg.Params.Validate(); err != nil {
		return sdkerrors.Wrap(err, "params")
	}
	return nil
}

func (msg *MsgUpdateParams) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{acc}
}

func (msg *MsgSendNFTToCosmosClaim) ValidateBasic() error {
	if err := gravitytypes.ValidateEthAddress(msg.EthereumSender); err != nil {
		return sdkerrors.Wrap(err, "eth sender")
	}
	if err := gravitytypes.ValidateEthAddress(msg.TokenContract); err != nil {
		return sdkerrors.Wrap(err, "nft contract address")
	}
	if _, err := sdk.AccAddressFromBech32(msg.Orchestrator); err != nil {
		return sdkerrors.Wrap(errors.ErrInvalidAddress, "orchestrator")
	}
	// note the destination address is intentionally not validated here, since
	// MsgSendNFTToCosmosClaim has it's destination as a string many invalid inputs are possible
	// the orchestrator will convert these invalid deposits to simply the string invalid'
	// this is done because the oracle requires an event be processed on Cosmos for each event
	// nonce on the Ethereum side, otherwise (A) the oracle will never proceed and (B) the funds
	// sent with the invalid deposit will forever be lost, with no representation minted anywhere
	// on cosmos. The attestation handler deals with this by managing invalid deposits and placing
	// them into the community pool
	if msg.EventNonce == 0 {
		return fmt.Errorf("nonce == 0")
	}
	return nil
}

func (msg *MsgSendNFTToCosmosClaim) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(msg.Orchestrator)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{acc}
}

func (m *MsgERC721DeployedClaim) ValidateBasic() error {
	//TODO implement me
	panic("implement me")
}

func (m *MsgERC721DeployedClaim) GetSigners() []sdk.AccAddress {
	//TODO implement me
	panic("implement me")
}

func (m *MsgExecuteIbcNFTAutoForwards) ValidateBasic() error {
	//TODO implement me
	panic("implement me")
}

func (m *MsgExecuteIbcNFTAutoForwards) GetSigners() []sdk.AccAddress {
	//TODO implement me
	panic("implement me")
}

func (m *MsgSendNFTToEth) ValidateBasic() error {
	//TODO implement me
	panic("implement me")
}

func (m *MsgSendNFTToEth) GetSigners() []sdk.AccAddress {
	//TODO implement me
	panic("implement me")
}

func (m *MsgSendNFTToEthClaim) ValidateBasic() error {
	//TODO implement me
	panic("implement me")
}

func (m *MsgSendNFTToEthClaim) GetSigners() []sdk.AccAddress {
	//TODO implement me
	panic("implement me")
}

func (m *MsgCancelSendNFTToEth) ValidateBasic() error {
	//TODO implement me
	panic("implement me")
}

func (m *MsgCancelSendNFTToEth) GetSigners() []sdk.AccAddress {
	//TODO implement me
	panic("implement me")
}

// EthereumNFTClaim represents a claim on ethereum state
type EthereumNFTClaim interface {
	// All Ethereum claims that we relay from the Gravity contract and into the module
	// have a nonce that is strictly increasing and unique, since this nonce is
	// issued by the Ethereum contract it is immutable and must be agreed on by all validators
	// any disagreement on what claim goes to what nonce means someone is lying.
	GetEventNonce() uint64
	// The block height that the claimed event occurred on. This EventNonce provides sufficient
	// ordering for the execution of all claims. The block height is used only for batchTimeouts + logicTimeouts
	// when we go to create a new batch we set the timeout some number of batches out from the last
	// known height plus projected block progress since then.
	GetEthBlockHeight() uint64
	// the delegate address of the claimer, for MsgDepositClaim and MsgWithdrawClaim
	// this is sent in as the sdk.AccAddress of the delegated key. it is up to the user
	// to disambiguate this into a sdk.ValAddress
	GetClaimer() sdk.AccAddress
	// Which type of claim this is
	GetType() NFTClaimType
	ValidateBasic() error
	// The claim hash of this claim. This is used to store these claims and also used to check if two different
	// validators claims agree. Therefore it's extremely important that this include all elements of the claim
	// with the exception of the orchestrator who sent it in, which will be used as a different part of the index
	ClaimHash() ([]byte, error)
	// Sets the orchestrator value on the claim
	SetOrchestrator(sdk.AccAddress)
}

// nolint: exhaustruct
var (
	_ EthereumNFTClaim = &MsgSendNFTToCosmosClaim{}
	_ EthereumNFTClaim = &MsgERC721DeployedClaim{}
	_ EthereumNFTClaim = &MsgSendNFTToEthClaim{}
)

func (msg *MsgSendNFTToCosmosClaim) GetClaimer() sdk.AccAddress {
	err := msg.ValidateBasic()
	if err != nil {
		panic("MsgSendNFTToCosmosClaim failed ValidateBasic! Should have been handled earlier")
	}

	val, err := sdk.AccAddressFromBech32(msg.Orchestrator)
	if err != nil {
		panic(err)
	}

	return val
}

func (msg *MsgSendNFTToCosmosClaim) GetType() NFTClaimType {
	return NFT_CLAIM_TYPE_SEND_NFT_TO_COSMOS
}

func (msg *MsgSendNFTToCosmosClaim) ClaimHash() ([]byte, error) {
	path := fmt.Sprintf("%d/%d/%s/%s/%s/%s/%s", msg.EventNonce, msg.EthBlockHeight, msg.TokenContract, msg.TokenId, msg.TokenUri, msg.EthereumSender, msg.CosmosReceiver)
	return tmhash.Sum([]byte(path)), nil
}

func (msg *MsgSendNFTToCosmosClaim) SetOrchestrator(orchestrator sdk.AccAddress) {
	msg.Orchestrator = orchestrator.String()
}

func (m *MsgERC721DeployedClaim) GetClaimer() sdk.AccAddress {
	//TODO implement me
	panic("implement me")
}

func (m *MsgERC721DeployedClaim) GetType() NFTClaimType {
	//TODO implement me
	panic("implement me")
}

func (m *MsgERC721DeployedClaim) ClaimHash() ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MsgERC721DeployedClaim) SetOrchestrator(address sdk.AccAddress) {
	//TODO implement me
	panic("implement me")
}

func (m *MsgSendNFTToEthClaim) GetClaimer() sdk.AccAddress {
	//TODO implement me
	panic("implement me")
}

func (m *MsgSendNFTToEthClaim) GetType() NFTClaimType {
	//TODO implement me
	panic("implement me")
}

func (m *MsgSendNFTToEthClaim) ClaimHash() ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MsgSendNFTToEthClaim) SetOrchestrator(address sdk.AccAddress) {
	//TODO implement me
	panic("implement me")
}