package keeper

import (
	"context"
	sdkerrors "cosmossdk.io/errors"
	"github.com/Gravity-Bridge/Gravity-Bridge/module/x/gravitynft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"strings"
)

// nolint: exhaustruct
var _ types.QueryServer = Keeper{}

const queryAttentionsLimit uint64 = 1000

// Params queries the params of the gravity module
func (k Keeper) Params(c context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	params := k.GetParams(sdk.UnwrapSDKContext(c))

	return &types.QueryParamsResponse{Params: params}, nil
}

func (k Keeper) GetLastObservedNFTEthBlock(c context.Context, req *types.QueryLastObservedNFTEthBlockRequest) (*types.QueryLastObservedNFTEthBlockResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	ethHeight := k.GetLastObservedEthereumBlockHeight(ctx)
	return &types.QueryLastObservedNFTEthBlockResponse{Block: ethHeight.EthereumBlockHeight}, nil
}

func (k Keeper) GetLastObservedNFTEthNonce(c context.Context, msg *types.QueryLastObservedNFTEthNonceRequest) (*types.QueryLastObservedNFTEthNonceResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	nonce := k.GetLastObservedEventNonce(ctx)
	return &types.QueryLastObservedNFTEthNonceResponse{Nonce: nonce}, nil
}

func (k Keeper) GetNFTAttestations(c context.Context, req *types.QueryNFTAttestationsRequest) (*types.QueryNFTAttestationsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	iterator := k.IterateAttestations

	limit := req.Limit
	if limit == 0 || limit > queryAttentionsLimit {
		limit = queryAttentionsLimit
	}

	var (
		attestations []types.NFTAttestation
		count        uint64
		iterErr      error
	)

	reverse := strings.EqualFold(req.OrderBy, "desc")
	filter := req.Height > 0 || req.Nonce > 0 || req.ClaimType != ""

	iterator(ctx, reverse, func(_ []byte, att types.NFTAttestation) (abort bool) {
		claim, err := k.UnpackAttestationClaim(&att)
		if err != nil {
			iterErr = sdkerrors.Wrap(errors.ErrUnpackAny, "failed to unmarshal Ethereum claim")
			return true
		}

		var match bool
		switch {
		case filter && claim.GetEthBlockHeight() == req.Height:
			attestations = append(attestations, att)
			match = true

		case filter && claim.GetEventNonce() == req.Nonce:
			attestations = append(attestations, att)
			match = true

		case filter && claim.GetType().String() == req.ClaimType:
			attestations = append(attestations, att)
			match = true

		case !filter:
			// No filter provided, so we include the attestation. This is equivalent
			// to providing no query params or just limit and/or order_by.
			attestations = append(attestations, att)
			match = true
		}

		if match {
			count++
			if count >= limit {
				return true
			}
		}

		return false
	})
	if iterErr != nil {
		return nil, iterErr
	}

	return &types.QueryNFTAttestationsResponse{Attestations: attestations}, nil
}

func (k Keeper) GetPendingNFTIbcAutoForwards(c context.Context, req *types.QueryPendingNFTIbcAutoForwards) (*types.QueryPendingNFTIbcAutoForwardsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	pendingForwards := k.PendingNFTIbcAutoForwards(ctx, req.Limit)
	return &types.QueryPendingNFTIbcAutoForwardsResponse{PendingIbcAutoForwards: pendingForwards}, nil
}
