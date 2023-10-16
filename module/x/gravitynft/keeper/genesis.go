package keeper

import (
	"github.com/Gravity-Bridge/Gravity-Bridge/module/x/gravitynft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis starts a chain from a genesis state
func InitGenesis(ctx sdk.Context, k Keeper, data types.GenesisState) {
	if err := k.setParams(ctx, *data.Params); err != nil {
		panic(err)
	}

	k.SetLastObservedEventNonce(ctx, data.GravityNonces.LastObservedNftNonce)
	// TODO: restore slash send NFT to ETH

	// TODO: Implement me!
}

// ExportGenesis exports all the state needed to restart the chain
// from the current state of the chain
func ExportGenesis(ctx sdk.Context, k Keeper) types.GenesisState {
	// TODO: Implement me!
	panic("ExportGenesis not implemented yet")
}