package keeper

import (
	"context"
	"github.com/Gravity-Bridge/Gravity-Bridge/module/x/gravitynft/types"
)

type msgServer struct {
	Keeper
}

// nolint: exhaustruct
var _ types.MsgServer = Keeper{}

// NewMsgServerImpl returns an implementation of the gov MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

func (k Keeper) SendNFTToCosmosClaim(ctx context.Context, claim *types.MsgSendNFTToCosmosClaim) (*types.MsgSendNFTToCosmosClaimResponse, error) {
	//TODO implement me
	panic("implement me")
}
