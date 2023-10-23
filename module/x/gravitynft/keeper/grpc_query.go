package keeper

import (
	"context"
	"strings"

	sdkerrors "cosmossdk.io/errors"
	gravitytypes "github.com/Gravity-Bridge/Gravity-Bridge/module/x/gravity/types"
	"github.com/Gravity-Bridge/Gravity-Bridge/module/x/gravitynft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
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
	return &types.QueryLastObservedNFTEthBlockResponse{LastObservedNftEthBlock: ethHeight.EthereumBlockHeight}, nil
}

func (k Keeper) GetLastObservedNFTEthNonce(c context.Context, msg *types.QueryLastObservedNFTEthNonceRequest) (*types.QueryLastObservedNFTEthNonceResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	nonce := k.GetLastObservedEventNonce(ctx)
	return &types.QueryLastObservedNFTEthNonceResponse{LastObservedNftEthNonce: nonce}, nil
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

	return &types.QueryNFTAttestationsResponse{NftAttestations: attestations}, nil
}

func (k Keeper) GetPendingNFTIbcAutoForwards(c context.Context, req *types.QueryPendingNFTIbcAutoForwardsRequest) (*types.QueryPendingNFTIbcAutoForwardsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	pendingForwards := k.PendingNFTIbcAutoForwards(ctx, req.Limit)
	return &types.QueryPendingNFTIbcAutoForwardsResponse{PendingNftIbcAutoForwards: pendingForwards}, nil
}

func (k Keeper) LastNFTEventNonceByAddr(c context.Context, req *types.QueryLastNFTEventNonceByAddrRequest) (*types.QueryLastNFTEventNonceByAddrResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	var ret types.QueryLastNFTEventNonceByAddrResponse
	addr, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return nil, sdkerrors.Wrap(errors.ErrInvalidAddress, req.Address)
	}
	validator, found := k.gravityKeeper.GetOrchestratorValidator(ctx, addr)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrUnknown, "address")
	}
	if err := sdk.VerifyAddressFormat(validator.GetOperator()); err != nil {
		return nil, sdkerrors.Wrap(err, "invalid validator address")
	}
	lastEventNonce := k.GetLastEventNonceByValidator(ctx, validator.GetOperator())
	ret.LastNftEventNonce = lastEventNonce
	return &ret, nil
}

func (k Keeper) OutgoingSendNFTToEths(c context.Context, req *types.QueryOutgoingSendNFTToEthsRequest) (*types.QueryOutgoingSendNFTToEthsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (k Keeper) ERC721ToClassId(c context.Context, req *types.QueryERC721ToClassIdRequest) (*types.QueryERC721ToClassIdResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	ethAddr, err := gravitytypes.NewEthAddress(req.Erc721)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "invalid ERC721 in request: %s", req.Erc721)
	}
	cosmosOriginated, classID := k.ERC721ToClassIDLookup(ctx, *ethAddr)
	var ret types.QueryERC721ToClassIdResponse
	ret.ClassId = classID
	ret.CosmosOriginated = cosmosOriginated

	return &ret, nil
}

func (k Keeper) ClassIdToERC721(c context.Context, req *types.QueryClassIdToERC721Request) (*types.QueryClassIdToERC721Response, error) {
	ctx := sdk.UnwrapSDKContext(c)
	cosmosOriginated, erc721, err := k.ClassIDToERC721Lookup(ctx, req.ClassId)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "invalid classID (%v) queried", req.ClassId)
	}
	var ret types.QueryClassIdToERC721Response
	ret.Erc721 = erc721.GetAddress().Hex()
	ret.CosmosOriginated = cosmosOriginated

	return &ret, err
}
