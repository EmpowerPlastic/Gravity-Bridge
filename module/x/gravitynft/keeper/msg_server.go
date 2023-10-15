package keeper

import (
	"context"
	"cosmossdk.io/errors"
	"github.com/Gravity-Bridge/Gravity-Bridge/module/x/gravitynft/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

type msgServer struct {
	Keeper
}

// nolint: exhaustruct
var _ types.MsgServer = msgServer{}

// NewMsgServerImpl returns an implementation of the gov MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

func (k msgServer) UpdateParams(c context.Context, msg *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	if k.authority != msg.Authority {
		return nil, errors.Wrapf(sdkerrors.ErrUnauthorized, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
	}

	ctx := sdk.UnwrapSDKContext(c)

	if err := k.setParams(ctx, msg.Params); err != nil {
		return nil, err
	}

	return &types.MsgUpdateParamsResponse{}, nil
}

func (k msgServer) SendNFTToCosmosClaim(c context.Context, msg *types.MsgSendNFTToCosmosClaim) (*types.MsgSendNFTToCosmosClaimResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	if err := k.gravityKeeper.CheckOrchestratorValidatorInSet(ctx, msg.Orchestrator); err != nil {
		return nil, sdkerrors.Wrap(err, "Orchstrator validator not in set")
	}
	any, err := codectypes.NewAnyWithValue(msg)
	if err != nil {
		return nil, errors.Wrap(err, "Could not check Any value")
	}

	err = k.claimHandlerCommon(ctx, any, msg)
	if err != nil {
		return nil, err
	}

	return &types.MsgSendNFTToCosmosClaimResponse{}, nil
}

// claimHandlerCommon is an internal function that provides common code for processing claims once they are
// translated from the message to the Ethereum claim interface
func (k msgServer) claimHandlerCommon(ctx sdk.Context, msgAny *codectypes.Any, msg types.EthereumNFTClaim) error {
	// Add the claim to the store
	_, err := k.Attest(ctx, msg, msgAny)
	if err != nil {
		return errors.Wrap(err, "create attestation")
	}
	hash, err := msg.ClaimHash()
	if err != nil {
		return errors.Wrap(err, "unable to compute claim hash")
	}

	// Emit the handle message event
	return ctx.EventManager().EmitTypedEvent(
		&types.EventClaim{
			Message:       string(msg.GetType()),
			ClaimHash:     string(hash),
			AttestationId: string(types.GetAttestationKey(msg.GetEventNonce(), hash)),
		},
	)
}

func (k msgServer) SendNFTToEth(c context.Context, msg *types.MsgSendNFTToEth) (*types.MsgSendNFTToEthResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (k msgServer) ExecuteIbcNFTAutoForwards(c context.Context, msg *types.MsgExecuteIbcNFTAutoForwards) (*types.MsgExecuteIbcNFTAutoForwardsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (k msgServer) SendNFTToEthClaim(c context.Context, msg *types.MsgSendNFTToEthClaim) (*types.MsgSendNFTToEthClaimResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (k msgServer) ERC721DeployedClaim(c context.Context, msg *types.MsgERC721DeployedClaim) (*types.MsgERC721DeployedClaimResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (k msgServer) CancelSendNFTToEth(c context.Context, msg *types.MsgCancelSendNFTToEth) (*types.MsgCancelSendNFTToEthResponse, error) {
	//TODO implement me
	panic("implement me")
}