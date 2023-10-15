package keeper_test

import (
	gocontext "context"
	"testing"

	"github.com/Gravity-Bridge/Gravity-Bridge/module/app"
	"github.com/Gravity-Bridge/Gravity-Bridge/module/x/gravitynft/keeper"
	"github.com/Gravity-Bridge/Gravity-Bridge/module/x/gravitynft/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

// nolint: exhaustruct
func TestQueryGetNFTAttestations(t *testing.T) {
	input := keeper.CreateTestEnv(t)
	encCfg := app.MakeEncodingConfig()
	k := input.GravityNFTKeeper
	ctx := input.Context

	queryHelper := baseapp.NewQueryServerTestHelper(ctx, encCfg.InterfaceRegistry)
	types.RegisterQueryServer(queryHelper, k)
	queryClient := types.NewQueryClient(queryHelper)

	numAttestations := 10
	createAttestations(t, k, ctx, numAttestations)

	testCases := []struct {
		name      string
		req       *types.QueryNFTAttestationsRequest
		numResult int
		nonces    []uint64
		expectErr bool
	}{
		{
			name:      "no params (all attestations ascending)",
			req:       &types.QueryNFTAttestationsRequest{},
			numResult: numAttestations,
			nonces:    []uint64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			expectErr: false,
		},
		{
			name: "all attestations descending",
			req: &types.QueryNFTAttestationsRequest{
				OrderBy: "desc",
			},
			numResult: numAttestations,
			nonces:    []uint64{10, 9, 8, 7, 6, 5, 4, 3, 2, 1},
			expectErr: false,
		},
		{
			name: "all attestations descending",
			req: &types.QueryNFTAttestationsRequest{
				OrderBy: "desc",
			},
			numResult: numAttestations,
			nonces:    []uint64{10, 9, 8, 7, 6, 5, 4, 3, 2, 1},
			expectErr: false,
		},
		{
			name: "filter by height and limit",
			req: &types.QueryNFTAttestationsRequest{
				Height: 1,
				Limit:  5,
			},
			numResult: 5,
			nonces:    []uint64{1, 2, 3, 4, 5},
			expectErr: false,
		},
		{
			name: "filter by nonce and limit",
			req: &types.QueryNFTAttestationsRequest{
				Nonce: 7,
				Limit: 5,
			},
			numResult: 1,
			nonces:    []uint64{7},
			expectErr: false,
		},
		{
			name: "filter by missing nonce",
			req: &types.QueryNFTAttestationsRequest{
				Nonce: 100000,
				Limit: 5,
			},
			numResult: 0,
			nonces:    []uint64{},
			expectErr: false,
		},
		{
			name: "filter by invalid claim type",
			req: &types.QueryNFTAttestationsRequest{
				ClaimType: "foo",
				Limit:     5,
			},
			numResult: 0,
			nonces:    []uint64{},
			expectErr: false,
		},
	}

	for i, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			result, err := queryClient.GetNFTAttestations(gocontext.Background(), tc.req)
			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Lenf(t, result.NftAttestations, tc.numResult, "unexpected number of results; tc #%d", i)

				nonces := make([]uint64, len(result.NftAttestations))
				for i, att := range result.NftAttestations {
					claim, err := k.UnpackAttestationClaim(&att)
					require.NoError(t, err)
					nonces[i] = claim.GetEventNonce()
				}
				require.Equal(t, tc.nonces, nonces)
			}
		})
	}
}

func createAttestations(t *testing.T, k keeper.Keeper, ctx sdk.Context, length int) {
	t.Helper()

	for i := 0; i < length; i++ {
		nonce := uint64(1 + i)
		msg := types.MsgSendNFTToCosmosClaim{
			EventNonce:     nonce,
			EthBlockHeight: 1,
			TokenContract:  "0x00000000000000000001",
			TokenId:         sdk.NewInt(1 + int64(i)).String(),
			TokenUri:       "",
			EthereumSender: "0x00000000000000000002",
			CosmosReceiver: "0x00000000000000000003",
			Orchestrator:   "0x00000000000000000004",
		}

		any, err := codectypes.NewAnyWithValue(&msg)
		require.NoError(t, err)

		att := &types.NFTAttestation{
			Observed: false,
			Votes:    []string{},
			Height:   uint64(ctx.BlockHeight()),
			Claim:    any,
		}

		hash, err := msg.ClaimHash()
		require.NoError(t, err)

		k.SetAttestation(ctx, nonce, hash, att)
	}
}
