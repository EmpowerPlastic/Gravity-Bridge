package types

import (
	sdkerrors "cosmossdk.io/errors"
	gravitytypes "github.com/Gravity-Bridge/Gravity-Bridge/module/x/gravity/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName is the name of the module
	ModuleName = "gravitynft"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName

	// RouterKey is the module name router key
	RouterKey = ModuleName

	// QuerierRoute to be used for querierer msgs
	QuerierRoute = ModuleName
)

var (
	// OracleAttestationKey attestation details by nonce and validator address
	// i.e. gravityvaloper1ahx7f8wyertuus9r20284ej0asrs085ceqtfnm
	// An attestation can be thought of as the 'event to be executed' while
	// the Claims are an individual validator saying that they saw an event
	// occur the Attestation is 'the event' that multiple claims vote on and
	// eventually executes
	// [0x0bfa165ff4ef558b3d0b62ea4d4a46c5]
	OracleAttestationKey = HashString("OracleAttestationKey")

	// LastObservedEthereumBlockHeightKey indexes the latest Ethereum block height
	// [0x83a283a6c3390f1526250df45e9ef8c6]
	LastObservedEthereumBlockHeightKey = HashString("LastObservedEthereumBlockHeightKey")

	// LastEventNonceByValidatorKey indexes lateset event nonce by validator
	// [0xeefcb999cc3d7b80b052b55106a6ba5e]
	LastEventNonceByValidatorKey = HashString("LastEventNonceByValidatorKey")

	// LastObservedEventNonceKey indexes the latest event nonce
	// [0xa34e56ab6fab9ee91e82ba216bfeb759]
	LastObservedEventNonceKey = HashString("LastObservedEventNonceKey")

	// ERC721ToDenomKey prefixes the index of Cosmos originated assets ERC721s to denoms
	// [caf3c584d2930aa14feb256fca05ee1f]
	ERC721ToDenomKey = HashString("ERC721ToDenomKey")

	// PendingNFTIbcAutoForwards indexes pending SendNFTToCosmos sends via IBC, queued by event nonce
	// [bbe244784132579f2e9b3b8577b57ad3]
	PendingNFTIbcAutoForwards = HashString("IbcNFTAutoForwardQueue")

	// ParamsKey indexes the parameters of the module
	// [4e2242f221531924f77b7250660af487]
	ParamsKey = HashString("Params")
)

// GetAttestationKey returns the following key format
// prefix     nonce                             claim-details-hash
// [0x0][0 0 0 0 0 0 0 1][fd1af8cec6c67fcf156f1b61fdf91ebc04d05484d007436e75342fc05bbff35a]
// An attestation is an event multiple people are voting on, this function needs the claim
// details because each Attestation is aggregating all claims of a specific event, lets say
// validator X and validator y were making different claims about the same event nonce
// Note that the claim hash does NOT include the claimer address and only identifies an event
func GetAttestationKey(eventNonce uint64, claimHash []byte) []byte {
	return AppendBytes(OracleAttestationKey, UInt64Bytes(eventNonce), claimHash)
}

// GetLastEventNonceByValidatorKey indexes latest event nonce by validator
// GetLastEventNonceByValidatorKey returns the following key format
// prefix              cosmos-validator
// [0x0][gravity1ahx7f8wyertuus9r20284ej0asrs085ceqtfnm]
func GetLastEventNonceByValidatorKey(validator sdk.ValAddress) []byte {
	if err := sdk.VerifyAddressFormat(validator); err != nil {
		panic(sdkerrors.Wrap(err, "invalid validator address"))
	}
	return AppendBytes(LastEventNonceByValidatorKey, validator.Bytes())
}

func GetERC721ToDenomKey(erc721 gravitytypes.EthAddress) []byte {
	return AppendBytes(ERC721ToDenomKey, erc721.GetAddress().Bytes())
}

// GetPendingNFTIbcAutoForwardKey returns the following key format
// prefix		EventNonce
// [0x0][0 0 0 0 0 0 0 1]
func GetPendingNFTIbcAutoForwardKey(eventNonce uint64) []byte {
	return AppendBytes(PendingNFTIbcAutoForwards, UInt64Bytes(eventNonce))
}