package keeper_test

import (
	gravitytypes "github.com/Gravity-Bridge/Gravity-Bridge/module/x/gravity/types"
	"github.com/Gravity-Bridge/Gravity-Bridge/module/x/gravitynft"
	gravitynftkeeper "github.com/Gravity-Bridge/Gravity-Bridge/module/x/gravitynft/keeper"
	"github.com/Gravity-Bridge/Gravity-Bridge/module/x/gravitynft/types"
	bech32ibctypes "github.com/althea-net/bech32-ibc/x/bech32ibc/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestUpdateParams(t *testing.T) {
	input, ctx := SetupFiveValChain(t)

	sv := gravitynftkeeper.MsgServer{Keeper: input.GravityNFTKeeper}
	msg := &types.MsgUpdateParams{
		Authority: AccAddrs[0].String(),
		Params: types.Params{
			// TODO: ADD REAL PARAMS
		},
	}
	_, err := sv.UpdateParams(ctx, msg)
	assert.ErrorContains(t, err, "invalid authority")

	msg.Authority = authtypes.NewModuleAddress(govtypes.ModuleName).String()
	_, err = sv.UpdateParams(ctx, msg)
	assert.NoError(t, err)

	params := input.GravityNFTKeeper.GetParams(ctx)
	_ = params
	// TODO: CHECK THE PARAMS ARE WHAT WE EXPECT
}

func TestSendNFTToCosmosClaimLocal(t *testing.T) {
	var (
		myBlockTime      = time.Date(2020, 9, 14, 15, 20, 10, 0, time.UTC)
		contract = types.NonemptyEthAddress()
		ethSender = types.NonemptyEthAddress()
		cosmosReceiver = types.NonemptySdkAccAddress().String()
		tokenId = "1"
		tokenUri = "testtokenuri"
	)

	input, ctx := SetupFiveValChain(t)

	sv := gravitynftkeeper.MsgServer{Keeper: input.GravityNFTKeeper}

	for _, v := range OrchAddrs {
		// each msg goes into its own block
		ctx = ctx.WithBlockTime(myBlockTime)

		claim := types.MsgSendNFTToCosmosClaim{
			EventNonce:     1,
			EthBlockHeight: 0,
			TokenContract:  contract,
			TokenId:        tokenId,
			TokenUri:       tokenUri,
			EthereumSender: ethSender,
			CosmosReceiver: cosmosReceiver,
			Orchestrator:   v.String(),
		}
		_, err := sv.SendNFTToCosmosClaim(ctx, &claim)
		assert.NoError(t, err)

		gravitynft.EndBlocker(ctx, input.GravityNFTKeeper)
		require.NoError(t, err)

		// and attestation persisted
		hash, err := claim.ClaimHash()
		require.NoError(t, err)
		a := input.GravityNFTKeeper.GetAttestation(ctx, uint64(1), hash)
		require.NotNil(t, a)

		// Test to reject duplicate deposit
		// when
		ctx = ctx.WithBlockTime(myBlockTime)
		_, err = sv.SendNFTToCosmosClaim(ctx, &claim)
		gravitynft.EndBlocker(ctx, input.GravityNFTKeeper)
		// then
		require.Error(t, err)
	}


	attestations := input.GravityNFTKeeper.GetMostRecentAttestations(ctx, 10)
	assert.Len(t, attestations, 1)

	classes := input.NftKeeper.GetClasses(ctx)
	assert.Len(t, classes, 1)
	assert.Equal(t, classes[0].Id, "erc721gravity" + contract)
	assert.Equal(t, classes[0].Name, "")
	assert.Equal(t, classes[0].Symbol, "")
	assert.Equal(t, classes[0].Description, "")
	assert.Equal(t, classes[0].Uri, "")
	assert.Equal(t, classes[0].UriHash, "")
	assert.Nil(t, classes[0].Data)

	nfts := input.NftKeeper.GetNFTsOfClass(ctx, "erc721gravity" + contract)
	assert.Len(t, nfts, 1)

	nft, found := input.NftKeeper.GetNFT(ctx, "erc721gravity" + contract, tokenId)
	assert.True(t, found)
	assert.Equal(t, nft.Id, tokenId)
	assert.Equal(t, nft.ClassId, "erc721gravity" + contract)
	assert.Equal(t, nft.Uri, tokenUri)
	assert.Nil(t, nft.Data)

	owner := input.NftKeeper.GetOwner(ctx, "erc721gravity" + contract, tokenId)
	assert.Equal(t, owner.String(), cosmosReceiver)
}

func TestSendNFTToCosmosClaimIBC(t *testing.T) {
	var (
		myBlockTime      = time.Date(2020, 9, 14, 15, 20, 10, 0, time.UTC)
		contract = types.NonemptyEthAddress()
		ethSender = types.NonemptyEthAddress()
		cosmosReceiver = "stars12zdmrp96wr0sccskearc0j6c6u4lpx4rmvnm2n"
		tokenId = "1"
		tokenUri = "testtokenuri"
	)

	input, ctx := SetupFiveValChain(t)
	input.Bech32IBCKeeper.SetHrpIbcRecords(ctx, []bech32ibctypes.HrpIbcRecord{
		{
			Hrp:                   "stars",
			FungibleSourceChannel: "",
			NftSourceChannel:      "channel-1337",
		},
	})

	sv := gravitynftkeeper.MsgServer{Keeper: input.GravityNFTKeeper}

	for _, v := range OrchAddrs {
		// each msg goes into its own block
		ctx = ctx.WithBlockTime(myBlockTime)

		claim := types.MsgSendNFTToCosmosClaim{
			EventNonce:     1,
			EthBlockHeight: 0,
			TokenContract:  contract,
			TokenId:        tokenId,
			TokenUri:       tokenUri,
			EthereumSender: ethSender,
			CosmosReceiver: cosmosReceiver,
			Orchestrator:   v.String(),
		}
		_, err := sv.SendNFTToCosmosClaim(ctx, &claim)
		assert.NoError(t, err)

		gravitynft.EndBlocker(ctx, input.GravityNFTKeeper)
		require.NoError(t, err)

		// and attestation persisted
		hash, err := claim.ClaimHash()
		require.NoError(t, err)
		a := input.GravityNFTKeeper.GetAttestation(ctx, uint64(1), hash)
		require.NotNil(t, a)

		// Test to reject duplicate deposit
		// when
		ctx = ctx.WithBlockTime(myBlockTime)
		_, err = sv.SendNFTToCosmosClaim(ctx, &claim)
		gravitynft.EndBlocker(ctx, input.GravityNFTKeeper)
		// then
		require.Error(t, err)
	}


	attestations := input.GravityNFTKeeper.GetMostRecentAttestations(ctx, 10)
	assert.Len(t, attestations, 1)

	forwards := input.GravityNFTKeeper.PendingNFTIbcAutoForwards(ctx, 2)
	assert.Len(t, forwards, 1)
	assert.Equal(t, forwards[0].ForeignReceiver, cosmosReceiver)
	assert.Equal(t, forwards[0].ClassId, "erc721gravity" + contract)
	assert.Equal(t, forwards[0].TokenId, tokenId)
	assert.Equal(t, forwards[0].IbcChannel, "channel-1337")
	assert.Equal(t, forwards[0].EventNonce, uint64(1))
}


func TestSendNFTToCosmosClaimWithNonOrchestrator(t *testing.T) {
	input, ctx := SetupFiveValChain(t)

	privKey, err := crypto.GenerateKey()
	require.NoError(t, err)

	ethAddress, err := gravitytypes.NewEthAddress(crypto.PubkeyToAddress(privKey.PublicKey).String())
	require.NoError(t, err)

	input.GravityKeeper.SetEthAddressForValidator(ctx, ValAddrs[0], *ethAddress)
	input.GravityKeeper.SetOrchestratorValidator(ctx, ValAddrs[0], AccAddrs[0])

	sv := gravitynftkeeper.MsgServer{Keeper: input.GravityNFTKeeper}
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
