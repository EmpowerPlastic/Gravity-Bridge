package keeper

import (
	gravitykeeper "github.com/Gravity-Bridge/Gravity-Bridge/module/x/gravity/keeper"
	"github.com/Gravity-Bridge/Gravity-Bridge/module/x/gravitynft/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdkstore "github.com/cosmos/cosmos-sdk/store/types"
)

// Keeper maintains the link to storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	// NOTE: If you add anything to this struct, add a nil check to ValidateMembers below!
	storeKey   sdkstore.StoreKey // Unexposed key to access store from sdk.Context

	// NOTE: If you add anything to this struct, add a nil check to ValidateMembers below!
	cdc               codec.BinaryCodec // The wire codec for binary encoding/decoding.
	gravityKeeper  *gravitykeeper.Keeper
	// TODO: Add stuffs
}

// Check for nil members
func (k Keeper) validateMembers() {
	if k.gravityKeeper == nil {
		panic("Nil bankKeeper!")
	}
}

// TODO: Create NewKeeper method and add it to app.go

// NewKeeper returns a new instance of the gravity keeper
func NewKeeper(storeKey sdkstore.StoreKey, cdc codec.BinaryCodec, gravityKeeper *gravitykeeper.Keeper) Keeper {
	k := Keeper{
		storeKey:   storeKey,
		cdc:               cdc,
		gravityKeeper:  gravityKeeper,
	}

	k.validateMembers()

	return k
}

func (k Keeper) UnpackAttestationClaim(att *types.NFTAttestation) (types.EthereumNFTClaim, error) {
	var msg types.EthereumNFTClaim
	err := k.cdc.UnpackAny(att.Claim, &msg)
	if err != nil {
		return nil, err
	} else {
		return msg, nil
	}
}

// prefixRange turns a prefix into a (start, end) range. The start is the given prefix value and
// the end is calculated by adding 1 bit to the start value. Nil is not allowed as prefix.
// Example: []byte{1, 3, 4} becomes []byte{1, 3, 5}
// []byte{15, 42, 255, 255} becomes []byte{15, 43, 0, 0}
//
// In case of an overflow the end is set to nil.
// Example: []byte{255, 255, 255, 255} becomes nil
// MARK finish-batches: this is where some crazy shit happens
func prefixRange(prefix []byte) ([]byte, []byte) {
	if prefix == nil {
		panic("nil key not allowed")
	}
	// special case: no prefix is whole range
	if len(prefix) == 0 {
		return nil, nil
	}

	// copy the prefix and update last byte
	end := make([]byte, len(prefix))
	copy(end, prefix)
	l := len(end) - 1
	end[l]++

	// wait, what if that overflowed?....
	for end[l] == 0 && l > 0 {
		l--
		end[l]++
	}

	// okay, funny guy, you gave us FFF, no end to this range...
	if l == 0 && end[0] == 0 {
		end = nil
	}
	return prefix, end
}