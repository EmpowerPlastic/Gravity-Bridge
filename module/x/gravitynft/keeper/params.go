package keeper

import (
	"github.com/Gravity-Bridge/Gravity-Bridge/module/x/gravitynft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ParamsKey)
	if bz == nil {
		return types.Params{}
	}

	var p types.Params
	k.cdc.MustUnmarshal(bz, &p)
	return p
}

func (k Keeper) setParams(ctx sdk.Context, p types.Params) error {
	if err := p.Validate(); err != nil {
		return err
	}

	store := ctx.KVStore(k.storeKey)
	bz, err := k.cdc.Marshal(&p)
	if err != nil {
		return err
	}
	store.Set(types.ParamsKey, bz)

	return nil
}