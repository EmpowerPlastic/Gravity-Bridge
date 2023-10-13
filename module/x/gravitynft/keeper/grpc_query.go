package keeper

import (
	"context"
	"github.com/Gravity-Bridge/Gravity-Bridge/module/x/gravitynft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// nolint: exhaustruct
var _ types.QueryServer = Keeper{}

// Params queries the params of the gravity module
func (k Keeper) Params(c context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	params := k.GetParams(sdk.UnwrapSDKContext(c))

	return &types.QueryParamsResponse{Params: params}, nil
}

func (k Keeper) GetLastObservedNFTEthBlock(c context.Context, msg *types.QueryLastObservedNFTEthBlockRequest) (*types.QueryLastObservedNFTEthBlockResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (k Keeper) GetLastObservedNFTEthNonce(c context.Context, msg *types.QueryLastObservedNFTEthNonceRequest) (*types.QueryLastObservedNFTEthNonceResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (k Keeper) GetNFTAttestations(c context.Context, msg *types.QueryNFTAttestationsRequest) (*types.QueryNFTAttestationsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (k Keeper) GetPendingNFTIbcAutoForwards(c context.Context, msg *types.QueryPendingNFTIbcAutoForwards) (*types.QueryPendingNFTIbcAutoForwardsResponse, error) {
	//TODO implement me
	panic("implement me")
}
