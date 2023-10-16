package keeper_test

import (
	"bytes"
	"fmt"
	gravitynftkeeper "github.com/Gravity-Bridge/Gravity-Bridge/module/x/gravitynft/keeper"
	"github.com/Gravity-Bridge/Gravity-Bridge/module/x/gravitynft/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetAndDeleteAttestation(t *testing.T) {
	input := CreateTestEnv(t)
	k := input.GravityNFTKeeper
	ctx := input.Context

	length := 10
	_, _, hashes := createAttestations(t, length, k, ctx)

	// Get created attestations
	for i := 0; i < length; i++ {
		nonce := uint64(1 + i)
		att := k.GetAttestation(ctx, nonce, hashes[i])
		require.NotNil(t, att)
	}

	recentAttestations := k.GetMostRecentAttestations(ctx, uint64(length))
	require.True(t, len(recentAttestations) == length)

	// Delete last 3 attestations
	var nilAtt *types.NFTAttestation
	for i := 7; i < length; i++ {
		nonce := uint64(1 + i)
		att := k.GetAttestation(ctx, nonce, hashes[i])
		k.DeleteAttestation(ctx, *att)

		att = k.GetAttestation(ctx, nonce, hashes[i])
		require.Equal(t, nilAtt, att)
	}
	recentAttestations = k.GetMostRecentAttestations(ctx, uint64(10))
	require.True(t, len(recentAttestations) == 7)

	// Check all attestations again
	for i := 0; i < 7; i++ {
		nonce := uint64(1 + i)
		att := k.GetAttestation(ctx, nonce, hashes[i])
		require.NotNil(t, att)
	}
	for i := 7; i < length; i++ {
		nonce := uint64(1 + i)
		att := k.GetAttestation(ctx, nonce, hashes[i])
		require.Equal(t, nilAtt, att)
	}
}

// Sets up 10 attestations and checks that they are returned in the correct order
func TestGetMostRecentAttestations(t *testing.T) {
	input := CreateTestEnv(t)
	defer func() { input.Context.Logger().Info("Asserting invariants at test end"); input.AssertInvariants() }()

	k := input.GravityNFTKeeper
	ctx := input.Context

	length := 10
	msgs, anys, _ := createAttestations(t, length, k, ctx)

	recentAttestations := k.GetMostRecentAttestations(ctx, uint64(length))
	require.True(t, len(recentAttestations) == length,
		"recentAttestations should have len %v but instead has %v", length, len(recentAttestations))
	for n, attest := range recentAttestations {
		require.Equal(t, attest.Claim.GetCachedValue(), anys[n].GetCachedValue(),
			"The %vth claim does not match our message: claim %v\n message %v", n, attest.Claim, msgs[n])
	}
}

// TODO: Test the nonces and last observered etc

func createAttestations(t *testing.T, length int, k gravitynftkeeper.Keeper, ctx sdktypes.Context) ([]types.MsgSendNFTToCosmosClaim, []codectypes.Any, [][]byte) {
	msgs := make([]types.MsgSendNFTToCosmosClaim, 0, length)
	anys := make([]codectypes.Any, 0, length)
	hashes := make([][]byte, 0, length)
	for i := 0; i < length; i++ {
		nonce := uint64(1 + i)

		contract := common.BytesToAddress(bytes.Repeat([]byte{0x1}, 20)).String()
		sender := common.BytesToAddress(bytes.Repeat([]byte{0x2}, 20)).String()
		orch := sdktypes.AccAddress(bytes.Repeat([]byte{0x3}, 20)).String()
		receiver := sdktypes.AccAddress(bytes.Repeat([]byte{0x4}, 20)).String()
		msg := types.MsgSendNFTToCosmosClaim{
			EventNonce:     nonce,
			EthBlockHeight: 1,
			TokenContract:  contract,
			TokenId:        "1",
			TokenUri:       "testtokenuri",
			EthereumSender: sender,
			CosmosReceiver: receiver,
			Orchestrator:   orch,
		}
		msgs = append(msgs, msg)

		any, err := codectypes.NewAnyWithValue(&msg)
		require.NoError(t, err)
		anys = append(anys, *any)
		att := &types.NFTAttestation{
			Observed: false,
			Votes:    []string{},
			Height:   uint64(ctx.BlockHeight()),
			Claim:    any,
		}
		unpackedClaim, err := k.UnpackAttestationClaim(att)
		if err != nil {
			panic(fmt.Sprintf("Bad new attestation: %s", err.Error()))
		}
		err = unpackedClaim.ValidateBasic()
		if err != nil {
			panic(fmt.Sprintf("Bad claim discovered: %s", err.Error()))
		}
		hash, err := msg.ClaimHash()
		hashes = append(hashes, hash)
		require.NoError(t, err)
		k.SetAttestation(ctx, nonce, hash, att)
	}

	return msgs, anys, hashes
}