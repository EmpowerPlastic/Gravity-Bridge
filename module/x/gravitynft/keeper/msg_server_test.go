package keeper

import (
	gravitytypes "github.com/Gravity-Bridge/Gravity-Bridge/module/x/gravity/types"
	"github.com/Gravity-Bridge/Gravity-Bridge/module/x/gravitynft/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSendNFTToCosmosClaim(t *testing.T) {
	input, ctx := SetupFiveValChain(t)

	privKey, err := crypto.GenerateKey()
	require.NoError(t, err)

	ethAddress, err := gravitytypes.NewEthAddress(crypto.PubkeyToAddress(privKey.PublicKey).String())
	require.NoError(t, err)

	input.GravityKeeper.SetEthAddressForValidator(ctx, ValAddrs[0], *ethAddress)
	input.GravityKeeper.SetOrchestratorValidator(ctx, ValAddrs[0], AccAddrs[0])

	sv := msgServer{input.GravityNFTKeeper}
	_, err = sv.SendNFTToCosmosClaim(ctx, &types.MsgSendNFTToCosmosClaim{
		EventNonce:     1,
		EthBlockHeight: 0,
		TokenContract:  types.NonemptyEthAddress(),
		TokenId:        "1",
		TokenUri:       "",
		EthereumSender: types.NonemptyEthAddress(),
		CosmosReceiver: types.NonemptySdkAccAddress().String(),
		Orchestrator:   AccAddrs[0].String(),
	})
	assert.NoError(t, err)

	attestations := input.GravityNFTKeeper.GetMostRecentAttestations(ctx, 10)
	assert.Len(t, attestations, 1)
	assert.Equal(t, attestations[0].Votes, []string{ValAddrs[0].String()})
	assert.Equal(t, attestations[0].Observed, false)
	assert.Equal(t, attestations[0].Height, uint64(ctx.BlockHeight()))

	lastEventNonce := input.GravityNFTKeeper.GetLastEventNonceByValidator(ctx, ValAddrs[0])
	assert.Equal(t, lastEventNonce, uint64(1))

	// TODO: Set up a test somewhere else that tests the processing of the claim after being "seen"
	// TODO: ^ Maybe in some kind of end block test?
}

func TestSendNFTToCosmosClaimWithNonOrchestrator(t *testing.T) {
	input, ctx := SetupFiveValChain(t)

	privKey, err := crypto.GenerateKey()
	require.NoError(t, err)

	ethAddress, err := gravitytypes.NewEthAddress(crypto.PubkeyToAddress(privKey.PublicKey).String())
	require.NoError(t, err)

	input.GravityKeeper.SetEthAddressForValidator(ctx, ValAddrs[0], *ethAddress)
	input.GravityKeeper.SetOrchestratorValidator(ctx, ValAddrs[0], AccAddrs[0])

	sv := msgServer{input.GravityNFTKeeper}
	_, err = sv.SendNFTToCosmosClaim(ctx, &types.MsgSendNFTToCosmosClaim{
		EventNonce:     1,
		EthBlockHeight: 0,
		TokenContract:  types.NonemptyEthAddress(),
		TokenId:        "1",
		TokenUri:       "",
		EthereumSender: types.NonemptyEthAddress(),
		CosmosReceiver: types.NonemptySdkAccAddress().String(),
		Orchestrator:   AccAddrs[1].String(),
	})
	assert.ErrorContains(t, err, "Orchstrator validator not in set")
}
