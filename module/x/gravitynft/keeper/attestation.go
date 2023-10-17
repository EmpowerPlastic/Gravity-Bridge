package keeper

import (
	sdkerrors "cosmossdk.io/errors"
	"fmt"
	"github.com/Gravity-Bridge/Gravity-Bridge/module/x/gravitynft/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"sort"
)

func (k Keeper) Attest(
	ctx sdk.Context,
	claim types.EthereumNFTClaim,
	anyClaim *codectypes.Any,
) (*types.NFTAttestation, error) {
	val, found := k.gravityKeeper.GetOrchestratorValidator(ctx, claim.GetClaimer())
	if !found {
		panic("Could not find ValAddr for delegate key, should be checked by now")
	}
	valAddr := val.GetOperator()
	if err := sdk.VerifyAddressFormat(valAddr); err != nil {
		return nil, sdkerrors.Wrap(err, "invalid orchestrator validator address")
	}
	// Check that the nonce of this event is exactly one higher than the last nonce stored by this validator.
	// We check the event nonce in processAttestation as well,
	// but checking it here gives individual eth signers a chance to retry,
	// and prevents validators from submitting two claims with the same nonce.
	// This prevents there being two attestations with the same nonce that get 2/3s of the votes
	// in the endBlocker.
	lastEventNonce := k.GetLastEventNonceByValidator(ctx, valAddr)
	if claim.GetEventNonce() != lastEventNonce+1 {
		return nil, fmt.Errorf(types.ErrNonContiguousEventNonce.Error(), lastEventNonce+1, claim.GetEventNonce())
	}

	// Tries to get an attestation with the same eventNonce and claim as the claim that was submitted.
	hash, err := claim.ClaimHash()
	if err != nil {
		return nil, sdkerrors.Wrap(err, "unable to compute claim hash")
	}
	att := k.GetAttestation(ctx, claim.GetEventNonce(), hash)

	// If it does not exist, create a new one.
	if att == nil {
		att = &types.NFTAttestation{
			Observed: false,
			Votes:    []string{},
			Height:   uint64(ctx.BlockHeight()),
			Claim:    anyClaim,
		}
	}

	ethClaim, err := k.UnpackAttestationClaim(att)
	if err != nil {
		panic(fmt.Sprintf("could not unpack stored attestation claim, %v", err))
	}

	if ethClaim.GetEthBlockHeight() == claim.GetEthBlockHeight() {

		// Add the validator's vote to this attestation
		att.Votes = append(att.Votes, valAddr.String())

		k.SetAttestation(ctx, claim.GetEventNonce(), hash, att)
		k.SetLastEventNonceByValidator(ctx, valAddr, claim.GetEventNonce())

		return att, nil
	} else {
		return nil, fmt.Errorf("invalid height - this claim's height is %v while the stored height is %v", claim.GetEthBlockHeight(), ethClaim.GetEthBlockHeight())
	}
}

// TryAttestation checks if an attestation has enough votes to be applied to the consensus state
// and has not already been marked Observed, then calls processAttestation to actually apply it to the state,
// and then marks it Observed and emits an event.
func (k Keeper) TryAttestation(ctx sdk.Context, att *types.NFTAttestation) {
	claim, err := k.UnpackAttestationClaim(att)
	if err != nil {
		panic("could not cast to claim")
	}
	hash, err := claim.ClaimHash()
	if err != nil {
		panic("unable to compute claim hash")
	}
	// If the attestation has not yet been Observed, sum up the votes and see if it is ready to apply to the state.
	// This conditional stops the attestation from accidentally being applied twice.
	if !att.Observed {
		// Sum the current powers of all validators who have voted and see if it passes the current threshold
		// TODO: The different integer types and math here needs a careful review
		totalPower := k.StakingKeeper.GetLastTotalPower(ctx)
		requiredPower := types.AttestationVotesPowerThreshold.Mul(totalPower).Quo(sdk.NewInt(100))
		attestationPower := sdk.NewInt(0)
		for _, validator := range att.Votes {
			val, err := sdk.ValAddressFromBech32(validator)
			if err != nil {
				panic(err)
			}
			validatorPower := k.StakingKeeper.GetLastValidatorPower(ctx, val)
			// Add it to the attestation power's sum
			attestationPower = attestationPower.Add(sdk.NewInt(validatorPower))
			// If the power of all the validators that have voted on the attestation is higher or equal to the threshold,
			// process the attestation, set Observed to true, and break
			if attestationPower.GT(requiredPower) {
				lastEventNonce := k.GetLastObservedEventNonce(ctx)
				// this check is performed at the next level up so this should never panic
				// outside of programmer error.
				if claim.GetEventNonce() != lastEventNonce+1 {
					panic("attempting to apply events to state out of order")
				}
				k.SetLastObservedEventNonce(ctx, claim.GetEventNonce())
				k.SetLastObservedEthereumBlockHeight(ctx, claim.GetEthBlockHeight())

				att.Observed = true
				k.SetAttestation(ctx, claim.GetEventNonce(), hash, att)

				k.processAttestation(ctx, att, claim)
				k.emitObservedEvent(ctx, att, claim)

				break
			}
		}
	} else {
		// We panic here because this should never happen
		panic("attempting to process observed attestation")
	}
}

// processAttestation actually applies the attestation to the consensus state
func (k Keeper) processAttestation(ctx sdk.Context, att *types.NFTAttestation, claim types.EthereumNFTClaim) {
	hash, err := claim.ClaimHash()
	if err != nil {
		panic("unable to compute claim hash")
	}
	// then execute in a new Tx so that we can store state on failure
	xCtx, commit := ctx.CacheContext()
	if err := k.AttestationHandler.Handle(xCtx, *att, claim); err != nil { // execute with a transient storage
		// If the attestation fails, something has gone wrong and we can't recover it. Log and move on
		// The attestation will still be marked "Observed", allowing the oracle to progress properly
		k.logger(ctx).Error("attestation failed",
			"cause", err.Error(),
			"claim type", claim.GetType(),
			"id", types.GetAttestationKey(claim.GetEventNonce(), hash),
			"nonce", fmt.Sprint(claim.GetEventNonce()),
		)
	} else {
		commit() // persist transient storage
	}
}

// emitObservedEvent emits an event with information about an attestation that has been applied to
// consensus state.
func (k Keeper) emitObservedEvent(ctx sdk.Context, att *types.NFTAttestation, claim types.EthereumNFTClaim) {
	hash, err := claim.ClaimHash()
	if err != nil {
		panic(sdkerrors.Wrap(err, "unable to compute claim hash"))
	}

	err = ctx.EventManager().EmitTypedEvent(
		&types.NFTEventObservation{
			AttestationType: string(claim.GetType()),
			// TODO: Add Bridge Contract Address
			// BridgeContract:  k.GetBridgeContractAddress(ctx).GetAddress().Hex(),
			// TODO: Add Bridge Chain ID
			// BridgeChainId:   strconv.Itoa(int(k.GetBridgeChainID(ctx))),
			AttestationId:   string(types.GetAttestationKey(claim.GetEventNonce(), hash)),
			Nonce:           fmt.Sprint(claim.GetEventNonce()),
		},
	)
	if err != nil {
		panic(err)
	}
}

// SetAttestation sets the attestation in the store
func (k Keeper) SetAttestation(ctx sdk.Context, eventNonce uint64, claimHash []byte, att *types.NFTAttestation) {
	store := ctx.KVStore(k.storeKey)
	aKey := types.GetAttestationKey(eventNonce, claimHash)
	store.Set(aKey, k.cdc.MustMarshal(att))
}

// GetAttestation return an attestation given a nonce
func (k Keeper) GetAttestation(ctx sdk.Context, eventNonce uint64, claimHash []byte) *types.NFTAttestation {
	store := ctx.KVStore(k.storeKey)
	aKey := types.GetAttestationKey(eventNonce, claimHash)
	bz := store.Get(aKey)
	if len(bz) == 0 {
		return nil
	}
	var att types.NFTAttestation
	k.cdc.MustUnmarshal(bz, &att)
	return &att
}

// DeleteAttestation deletes the given attestation
func (k Keeper) DeleteAttestation(ctx sdk.Context, att types.NFTAttestation) {
	claim, err := k.UnpackAttestationClaim(&att)
	if err != nil {
		panic("Bad Attestation in DeleteAttestation")
	}
	hash, err := claim.ClaimHash()
	if err != nil {
		panic(sdkerrors.Wrap(err, "unable to compute claim hash"))
	}
	store := ctx.KVStore(k.storeKey)

	store.Delete(types.GetAttestationKey(claim.GetEventNonce(), hash))
}

// SetLastObservedEthereumBlockHeight sets the block height in the store.
func (k Keeper) SetLastObservedEthereumBlockHeight(ctx sdk.Context, ethereumHeight uint64) {
	store := ctx.KVStore(k.storeKey)
	previous := k.GetLastObservedEthereumBlockHeight(ctx)
	if previous.EthereumBlockHeight > ethereumHeight {
		panic("Attempt to roll back Ethereum block height!")
	}
	height := types.LastObservedNFTEthereumBlockHeight{
		EthereumBlockHeight: ethereumHeight,
		CosmosBlockHeight:   uint64(ctx.BlockHeight()),
	}
	store.Set(types.LastObservedEthereumBlockHeightKey, k.cdc.MustMarshal(&height))
}

// GetLastObservedEthereumBlockHeight height gets the block height to of the last observed attestation from
// the store
func (k Keeper) GetLastObservedEthereumBlockHeight(ctx sdk.Context) types.LastObservedNFTEthereumBlockHeight {
	store := ctx.KVStore(k.storeKey)
	bytes := store.Get(types.LastObservedEthereumBlockHeightKey)

	if len(bytes) == 0 {
		return types.LastObservedNFTEthereumBlockHeight{
			CosmosBlockHeight:   0,
			EthereumBlockHeight: 0,
		}
	}
	height := types.LastObservedNFTEthereumBlockHeight{
		CosmosBlockHeight:   0,
		EthereumBlockHeight: 0,
	}
	k.cdc.MustUnmarshal(bytes, &height)
	return height
}

// SetLastObservedEventNonce sets the latest observed event nonce
func (k Keeper) SetLastObservedEventNonce(ctx sdk.Context, nonce uint64) {
	store := ctx.KVStore(k.storeKey)
	last := k.GetLastObservedEventNonce(ctx)
	// event nonce must increase, unless it's zero at which point allow zero to be set
	// as many times as needed (genesis test setup etc)
	zeroCase := last == 0 && nonce == 0
	if last >= nonce && !zeroCase {
		panic("Event nonce going backwards or replay!")
	}
	store.Set(types.LastObservedEventNonceKey, types.UInt64Bytes(nonce))
}

// GetLastEventNonceByValidator returns the latest event nonce for a given validator
func (k Keeper) GetLastEventNonceByValidator(ctx sdk.Context, validator sdk.ValAddress) uint64 {
	if err := sdk.VerifyAddressFormat(validator); err != nil {
		panic(sdkerrors.Wrap(err, "invalid validator address"))
	}
	store := ctx.KVStore(k.storeKey)
	bytes := store.Get(types.GetLastEventNonceByValidatorKey(validator))

	if len(bytes) == 0 {
		// in the case that we have no existing value this is the first
		// time a validator is submitting a claim. Since we don't want to force
		// them to replay the entire history of all events ever we can't start
		// at zero
		lastEventNonce := k.GetLastObservedEventNonce(ctx)
		if lastEventNonce >= 1 {
			return lastEventNonce - 1
		} else {
			return 0
		}
	}
	return types.UInt64FromBytesUnsafe(bytes)
}

// GetLastObservedEventNonce returns the latest observed event nonce
func (k Keeper) GetLastObservedEventNonce(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bytes := store.Get(types.LastObservedEventNonceKey)

	if len(bytes) == 0 {
		return 0
	}
	if len(bytes) > 8 {
		panic("Last observed event nonce is not a uint64!")
	}
	return types.UInt64FromBytesUnsafe(bytes)
}

// SetLastEventNonceByValidator sets the latest event nonce for a give validator
func (k Keeper) SetLastEventNonceByValidator(ctx sdk.Context, validator sdk.ValAddress, nonce uint64) {
	if err := sdk.VerifyAddressFormat(validator); err != nil {
		panic(sdkerrors.Wrap(err, "invalid validator address"))
	}
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetLastEventNonceByValidatorKey(validator), types.UInt64Bytes(nonce))
}

// GetMostRecentAttestations returns sorted (by nonce) attestations up to a provided limit number of attestations
// Note: calls GetAttestationMapping in the hopes that there are potentially many attestations
// which are distributed between few nonces to minimize sorting time
func (k Keeper) GetMostRecentAttestations(ctx sdk.Context, limit uint64) []types.NFTAttestation {
	attestationMapping, keys := k.GetAttestationMapping(ctx)
	attestations := make([]types.NFTAttestation, 0, limit)

	// Iterate the nonces and collect the attestations
	count := 0
	for _, nonce := range keys {
		if count >= int(limit) {
			break
		}
		for _, att := range attestationMapping[nonce] {
			if count >= int(limit) {
				break
			}
			attestations = append(attestations, att)
			count++
		}
	}

	return attestations
}

// GetAttestationMapping returns a mapping of eventnonce -> attestations at that nonce
// it also returns a pre-sorted array of the keys, this assists callers of this function
// by providing a deterministic iteration order. You should always iterate over ordered keys
// if you are iterating this map at all.
func (k Keeper) GetAttestationMapping(ctx sdk.Context) (attestationMapping map[uint64][]types.NFTAttestation, orderedKeys []uint64) {
	attestationMapping = make(map[uint64][]types.NFTAttestation)
	k.IterateAttestations(ctx, false, func(_ []byte, att types.NFTAttestation) bool {
		claim, err := k.UnpackAttestationClaim(&att)
		if err != nil {
			panic("couldn't cast to claim")
		}

		if val, ok := attestationMapping[claim.GetEventNonce()]; !ok {
			attestationMapping[claim.GetEventNonce()] = []types.NFTAttestation{att}
		} else {
			attestationMapping[claim.GetEventNonce()] = append(val, att)
		}
		return false
	})
	orderedKeys = make([]uint64, 0, len(attestationMapping))
	for k := range attestationMapping {
		orderedKeys = append(orderedKeys, k)
	}
	sort.Slice(orderedKeys, func(i, j int) bool { return orderedKeys[i] < orderedKeys[j] })

	return
}

// IterateAttestations iterates through all attestations executing a given callback on each discovered attestation
// If reverse is true, attestations will be returned in descending order by key (aka by event nonce and then claim hash)
// cb should return true to stop iteration, false to continue
func (k Keeper) IterateAttestations(ctx sdk.Context, reverse bool, cb func(key []byte, att types.NFTAttestation) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	keyPrefix := types.OracleAttestationKey

	var iter storetypes.Iterator
	if reverse {
		iter = store.ReverseIterator(prefixRange(keyPrefix))
	} else {
		iter = store.Iterator(prefixRange(keyPrefix))
	}
	defer func(iter storetypes.Iterator) {
		err := iter.Close()
		if err != nil {
			panic("Unable to close attestation iterator!")
		}
	}(iter)

	for ; iter.Valid(); iter.Next() {
		att := types.NFTAttestation{
			Observed: false,
			Votes:    []string{},
			Height:   0,
			Claim: &codectypes.Any{
				TypeUrl:              "",
				Value:                []byte{},
				XXX_NoUnkeyedLiteral: struct{}{},
				XXX_unrecognized:     []byte{},
				XXX_sizecache:        0,
			},
		}
		k.cdc.MustUnmarshal(iter.Value(), &att)
		// cb returns true to stop early
		if cb(iter.Key(), att) {
			return
		}
	}
}