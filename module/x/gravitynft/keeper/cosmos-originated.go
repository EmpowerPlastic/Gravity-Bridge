package keeper

import (
	gravitytypes "github.com/Gravity-Bridge/Gravity-Bridge/module/x/gravity/types"
	"github.com/Gravity-Bridge/Gravity-Bridge/module/x/gravitynft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) GetCosmosOriginatedClassID(ctx sdk.Context, tokenContract gravitytypes.EthAddress) (string, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetERC721ToDenomKey(tokenContract))

	if bz != nil {
		return string(bz), true
	}
	return "", false
}

// ERC721ToClassIDLookup returns (bool isCosmosOriginated, string classID, err)
// Using this information, you can see if an ERC721 address representing an asset is native to Cosmos or Ethereum,
// and get its corresponding classID
func (k Keeper) ERC721ToClassIDLookup(ctx sdk.Context, tokenContract gravitytypes.EthAddress) (bool, string) {
	// First try looking up tokenContract in index
	dn1, exists := k.GetCosmosOriginatedClassID(ctx, tokenContract)
	if exists {
		// It is a cosmos originated asset
		return true, dn1
	}

	// If it is not in there, it is not a cosmos originated token, turn the ERC20 into a gravity denom
	return false, types.GravityERC721ClassId(tokenContract)
}