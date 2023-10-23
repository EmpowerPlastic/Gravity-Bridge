package keeper

import (
	sdkerrors "cosmossdk.io/errors"
	"fmt"
	gravitytypes "github.com/Gravity-Bridge/Gravity-Bridge/module/x/gravity/types"
	"github.com/Gravity-Bridge/Gravity-Bridge/module/x/gravitynft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) GetCosmosOriginatedClassID(ctx sdk.Context, tokenContract gravitytypes.EthAddress) (string, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetERC721ToClassIDKey(tokenContract))

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

// ClassIDToERC721Lookup returns (bool isCosmosOriginated, EthAddress ERC721, err)
// Using this information, you can see if an asset is native to Cosmos or Ethereum,
// and get its corresponding ERC721 address.
func (k Keeper) ClassIDToERC721Lookup(ctx sdk.Context, classID string) (bool, *gravitytypes.EthAddress, error) {
	// First try parsing the ERC20 out of the denom
	tc1, err := types.GravityDenomToERC721(classID)

	if err != nil {
		// Look up ERC721 contract in index and error if it's not in there.
		tc2, exists := k.GetCosmosOriginatedERC721(ctx, classID)
		if !exists {
			return false, nil,
				sdkerrors.Wrap(types.ErrInvalid, fmt.Sprintf("classID not a gravity voucher nft: %s, and also not in cosmos-originated ERC721 index", err))
		}
		// This is a cosmos-originated asset
		return true, tc2, nil
	}

	// This is an ethereum-originated asset
	return false, tc1, nil
}

func (k Keeper) GetCosmosOriginatedERC721(ctx sdk.Context, classID string) (*gravitytypes.EthAddress, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetClassIDToERC721Key(classID))
	if bz != nil {
		ethAddr, err := gravitytypes.NewEthAddressFromBytes(bz)
		if err != nil {
			panic(fmt.Errorf("discovered invalid ERC721 address under key %v", string(bz)))
		}

		return ethAddr, true
	}
	return nil, false
}