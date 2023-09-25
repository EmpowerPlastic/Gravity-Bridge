package keeper

import (
	"context"
	"github.com/Gravity-Bridge/Gravity-Bridge/module/x/gravitynft/types"
)

// nolint: exhaustruct
var _ types.QueryServer = Keeper{}

func (k Keeper) Params(c context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	// TODO: implement
	panic("not implemented")
}
