package types

import (
	sdkerrors "cosmossdk.io/errors"
	gravitytypes "github.com/Gravity-Bridge/Gravity-Bridge/module/x/gravity/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// TODO: Find a proper module name. The original idea of gravitynft won't work because it crashes with the prefix of "grabity"
	// ModuleName is the name of the module
	ModuleName = "erc721gravity"

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

	// ClassIDToERC721Key prefixes the index of Cosmos originated asset denoms to ERC20s
	// [0f4178d32f57c540e8f15e93b90323d7]
	ClassIDToERC721Key = HashString("ClassIDToERC721Key")

	// ERC721ToClassIDKey prefixes the index of Cosmos originated assets ERC721s to ClassIDs
	// [1358ec20e89e72219658c1005cf6f9a0]
	ERC721ToClassIDKey = HashString("ERC721ToClassIDKey")

	// PendingNFTIbcAutoForwards indexes pending SendNFTToCosmos sends via IBC, queued by event nonce
	// [67b20bb5f0a19693de3469e155ef5352]
	PendingNFTIbcAutoForwards = HashString("NFTIbcAutoForwardQueue")

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

func GetClassIDToERC721Key(denom string) []byte {
	return AppendBytes(ClassIDToERC721Key, []byte(denom))
}

func GetERC721ToClassIDKey(erc721 gravitytypes.EthAddress) []byte {
	return AppendBytes(ERC721ToClassIDKey, erc721.GetAddress().Bytes())
}

// GetPendingNFTIbcAutoForwardKey returns the following key format
// prefix		EventNonce
// [0x0][0 0 0 0 0 0 0 1]
func GetPendingNFTIbcAutoForwardKey(eventNonce uint64) []byte {
	return AppendBytes(PendingNFTIbcAutoForwards, UInt64Bytes(eventNonce))
}