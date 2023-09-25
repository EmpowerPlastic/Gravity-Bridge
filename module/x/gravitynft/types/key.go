package types

import (
	sdkerrors "cosmossdk.io/errors"
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

	// LastEventNonceByValidatorKey indexes lateset event nonce by validator
	// [0xeefcb999cc3d7b80b052b55106a6ba5e]
	LastEventNonceByValidatorKey = HashString("LastEventNonceByValidatorKey")

	// LastObservedEventNonceKey indexes the latest event nonce
	// [0xa34e56ab6fab9ee91e82ba216bfeb759]
	LastObservedEventNonceKey = HashString("LastObservedEventNonceKey")
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
