package keeper

import (
	"fmt"
	gravitykeeper "github.com/Gravity-Bridge/Gravity-Bridge/module/x/gravity/keeper"
	"github.com/Gravity-Bridge/Gravity-Bridge/module/x/gravitynft/types"
	bech32ibckeeper "github.com/althea-net/bech32-ibc/x/bech32ibc/keeper"
	ibcnfttransferkeeper "github.com/bianjieai/nft-transfer/keeper"
	"github.com/cosmos/cosmos-sdk/codec"
	sdkstore "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/nft"
	nftkeeper "github.com/cosmos/cosmos-sdk/x/nft/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	"github.com/tendermint/tendermint/libs/log"
)

// Keeper maintains the link to storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	// NOTE: If you add anything to this struct, add a nil check to ValidateMembers below!
	storeKey sdkstore.StoreKey // Unexposed key to access store from sdk.Context

	// NOTE: If you add anything to this struct, add a nil check to ValidateMembers below!
	cdc                  codec.BinaryCodec // The wire codec for binary encoding/decoding.
	gravityKeeper        *gravitykeeper.Keeper
	StakingKeeper        *stakingkeeper.Keeper
	accountKeeper        *authkeeper.AccountKeeper
	bech32IbcKeeper      *bech32ibckeeper.Keeper
	nftKeeper            *nftkeeper.Keeper
	ibcNftTransferKeeper *ibcnfttransferkeeper.Keeper
	// TODO: Add stuffs

	AttestationHandler interface {
		Handle(sdk.Context, types.NFTAttestation, types.EthereumNFTClaim) error
	}
}

// Check for nil members
func (k Keeper) validateMembers() {
	if k.gravityKeeper == nil {
		panic("Nil bankKeeper!")
	}
	if k.StakingKeeper == nil {
		panic("Nil StakingKeeper!")
	}
	if k.accountKeeper == nil {
		panic("Nil accountKeeper!")
	}
	if k.bech32IbcKeeper == nil {
		panic("Nil bech32IbcKeeper!")
	}
	if k.nftKeeper == nil {
		panic("Nil nftKeeper!")
	}
	if k.ibcNftTransferKeeper == nil {
		panic("Nil ibcNftTransferKeeper!")
	}
}

// TODO: Create NewKeeper method and add it to app.go

// NewKeeper returns a new instance of the gravity keeper
func NewKeeper(
	storeKey sdkstore.StoreKey,
	cdc codec.BinaryCodec,
	gravityKeeper *gravitykeeper.Keeper,
	stakingKeeper *stakingkeeper.Keeper,
	accKeeper *authkeeper.AccountKeeper,
	bech32IbcKeeper *bech32ibckeeper.Keeper,
	nftKeeper *nftkeeper.Keeper,
	ibcNftTransferKeeper *ibcnfttransferkeeper.Keeper,
) Keeper {
	k := Keeper{
		storeKey:             storeKey,
		cdc:                  cdc,
		gravityKeeper:        gravityKeeper,
		StakingKeeper:        stakingKeeper,
		accountKeeper:        accKeeper,
		bech32IbcKeeper:      bech32IbcKeeper,
		nftKeeper:            nftKeeper,
		ibcNftTransferKeeper: ibcNftTransferKeeper,
	}
	attestationHandler := AttestationHandler{keeper: &k}
	attestationHandler.ValidateMembers()
	k.AttestationHandler = attestationHandler
	k.validateMembers()

	return k
}

////////////////////////
/////// HELPERS ////////
////////////////////////

// SendNFTToCommunityPool handles incorrect SendToCosmos calls to the community pool, since the calls
// have already been made on Ethereum there's nothing we can do to reverse them, and we should at least
// make use of the tokens which would otherwise be lost
func (k Keeper) SendNFTToCommunityPool(ctx sdk.Context, nftToken nft.NFT) error {
	communityPoolAddress := k.accountKeeper.GetModuleAddress(distrtypes.ModuleName)
	if communityPoolAddress == nil {
		return sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "transfer NFT to community pool failed: module account %s does not exist", distrtypes.ModuleName)
	}
	if err := k.nftKeeper.Transfer(ctx, nftToken.ClassId, nftToken.Id, communityPoolAddress); err != nil {
		return sdkerrors.Wrapf(err, "transfer NFT %s %s to community pool failed", nftToken.ClassId, nftToken.Id)
	}

	return nil
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

// logger returns a module-specific logger.
func (k Keeper) logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
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
