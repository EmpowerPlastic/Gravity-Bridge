package keeper_test

import (
	"bytes"
	sdkmath "cosmossdk.io/math"
	gravityparams "github.com/Gravity-Bridge/Gravity-Bridge/module/app/params"
	gravitykeeper "github.com/Gravity-Bridge/Gravity-Bridge/module/x/gravity/keeper"
	gravitytypes "github.com/Gravity-Bridge/Gravity-Bridge/module/x/gravity/types"
	gravitynftkeeper "github.com/Gravity-Bridge/Gravity-Bridge/module/x/gravitynft/keeper"
	"github.com/Gravity-Bridge/Gravity-Bridge/module/x/gravitynft/types"
	"github.com/Gravity-Bridge/Gravity-Bridge/module/x/internft"
	bech32ibckeeper "github.com/althea-net/bech32-ibc/x/bech32ibc/keeper"
	bech32ibctypes "github.com/althea-net/bech32-ibc/x/bech32ibc/types"
	ibcnfttransferkeeper "github.com/bianjieai/nft-transfer/keeper"
	ibcnfttransfertypes "github.com/bianjieai/nft-transfer/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	ccodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	ccrypto "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/std"
	"github.com/cosmos/cosmos-sdk/store"
	sdkstore "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting"
	"github.com/cosmos/cosmos-sdk/x/bank"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/capability"
	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	distrclient "github.com/cosmos/cosmos-sdk/x/distribution/client"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/evidence"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	govtypesv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/cosmos/cosmos-sdk/x/mint"
	nfttypes "github.com/cosmos/cosmos-sdk/x/nft"
	nftkeeper "github.com/cosmos/cosmos-sdk/x/nft/keeper"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramsclient "github.com/cosmos/cosmos-sdk/x/params/client"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	paramsproposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	upgradeclient "github.com/cosmos/cosmos-sdk/x/upgrade/client"
	upgradekeeper "github.com/cosmos/cosmos-sdk/x/upgrade/keeper"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	ibctransferkeeper "github.com/cosmos/ibc-go/v5/modules/apps/transfer/keeper"
	ibctransfertypes "github.com/cosmos/ibc-go/v5/modules/apps/transfer/types"
	ibchost "github.com/cosmos/ibc-go/v5/modules/core/24-host"
	ibckeeper "github.com/cosmos/ibc-go/v5/modules/core/keeper"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmversion "github.com/tendermint/tendermint/proto/tendermint/version"
	dbm "github.com/tendermint/tm-db"
	"testing"
	"time"
)

var (
	// ModuleBasics is a mock module basic manager for testing
	ModuleBasics = module.NewBasicManager(
		auth.AppModuleBasic{},
		genutil.AppModuleBasic{},
		bank.AppModuleBasic{},
		capability.AppModuleBasic{},
		staking.AppModuleBasic{},
		mint.AppModuleBasic{},
		distribution.AppModuleBasic{},
		gov.NewAppModuleBasic(
			[]govclient.ProposalHandler{
				paramsclient.ProposalHandler, distrclient.ProposalHandler, upgradeclient.LegacyProposalHandler, upgradeclient.LegacyCancelProposalHandler,
			},
		),
		params.AppModuleBasic{},
		crisis.AppModuleBasic{},
		slashing.AppModuleBasic{},
		upgrade.AppModuleBasic{},
		evidence.AppModuleBasic{},
		vesting.AppModuleBasic{},
	)
)

var (
	// TestingStakeParams is a set of staking params for testing
	TestingStakeParams = stakingtypes.Params{
		UnbondingTime:     100,
		MaxValidators:     10,
		MaxEntries:        10,
		HistoricalEntries: 10000,
		BondDenom:         "stake",
		MinCommissionRate: sdk.ZeroDec(),
	}

	// TestingGravityParams is a set of gravity params for testing
	TestingGravityParams = gravitytypes.Params{
		GravityId:                    "testgravityid",
		ContractSourceHash:           "62328f7bc12efb28f86111d08c29b39285680a906ea0e524e0209d6f6657b713",
		BridgeEthereumAddress:        "0x8858eeb3dfffa017d4bce9801d340d36cf895ccf",
		BridgeChainId:                11,
		SignedValsetsWindow:          10,
		SignedBatchesWindow:          10,
		SignedLogicCallsWindow:       10,
		TargetBatchTimeout:           60001,
		AverageBlockTime:             5000,
		AverageEthereumBlockTime:     15000,
		SlashFractionValset:          sdk.NewDecWithPrec(1, 2),
		SlashFractionBatch:           sdk.NewDecWithPrec(1, 2),
		SlashFractionLogicCall:       sdk.Dec{},
		UnbondSlashingValsetsWindow:  15,
		SlashFractionBadEthSignature: sdk.NewDecWithPrec(1, 2),
		ValsetReward:                 sdk.Coin{Denom: "", Amount: sdk.ZeroInt()},
		BridgeActive:                 true,
		EthereumBlacklist:            []string{},
		MinChainFeeBasisPoints:       0,
	}

	// ConsPrivKeys generate ed25519 ConsPrivKeys to be used for validator operator keys
	ConsPrivKeys = []ccrypto.PrivKey{
		ed25519.GenPrivKey(),
		ed25519.GenPrivKey(),
		ed25519.GenPrivKey(),
		ed25519.GenPrivKey(),
		ed25519.GenPrivKey(),
	}

	// ConsPubKeys holds the consensus public keys to be used for validator operator keys
	ConsPubKeys = []ccrypto.PubKey{
		ConsPrivKeys[0].PubKey(),
		ConsPrivKeys[1].PubKey(),
		ConsPrivKeys[2].PubKey(),
		ConsPrivKeys[3].PubKey(),
		ConsPrivKeys[4].PubKey(),
	}

	// AccPrivKeys generate secp256k1 pubkeys to be used for account pub keys
	AccPrivKeys = []ccrypto.PrivKey{
		secp256k1.GenPrivKey(),
		secp256k1.GenPrivKey(),
		secp256k1.GenPrivKey(),
		secp256k1.GenPrivKey(),
		secp256k1.GenPrivKey(),
	}

	// AccPubKeys holds the pub keys for the account keys
	AccPubKeys = []ccrypto.PubKey{
		AccPrivKeys[0].PubKey(),
		AccPrivKeys[1].PubKey(),
		AccPrivKeys[2].PubKey(),
		AccPrivKeys[3].PubKey(),
		AccPrivKeys[4].PubKey(),
	}

	// AccAddrs holds the sdk.AccAddresses
	AccAddrs = []sdk.AccAddress{
		sdk.AccAddress(AccPubKeys[0].Address()),
		sdk.AccAddress(AccPubKeys[1].Address()),
		sdk.AccAddress(AccPubKeys[2].Address()),
		sdk.AccAddress(AccPubKeys[3].Address()),
		sdk.AccAddress(AccPubKeys[4].Address()),
	}

	// ValAddrs holds the sdk.ValAddresses
	ValAddrs = []sdk.ValAddress{
		sdk.ValAddress(AccPubKeys[0].Address()),
		sdk.ValAddress(AccPubKeys[1].Address()),
		sdk.ValAddress(AccPubKeys[2].Address()),
		sdk.ValAddress(AccPubKeys[3].Address()),
		sdk.ValAddress(AccPubKeys[4].Address()),
	}

	// AccPubKeys holds the pub keys for the account keys
	OrchPubKeys = []ccrypto.PubKey{
		OrchPrivKeys[0].PubKey(),
		OrchPrivKeys[1].PubKey(),
		OrchPrivKeys[2].PubKey(),
		OrchPrivKeys[3].PubKey(),
		OrchPrivKeys[4].PubKey(),
	}

	// Orchestrator private keys
	OrchPrivKeys = []ccrypto.PrivKey{
		secp256k1.GenPrivKey(),
		secp256k1.GenPrivKey(),
		secp256k1.GenPrivKey(),
		secp256k1.GenPrivKey(),
		secp256k1.GenPrivKey(),
	}

	// AccAddrs holds the sdk.AccAddresses
	OrchAddrs = []sdk.AccAddress{
		sdk.AccAddress(OrchPubKeys[0].Address()),
		sdk.AccAddress(OrchPubKeys[1].Address()),
		sdk.AccAddress(OrchPubKeys[2].Address()),
		sdk.AccAddress(OrchPubKeys[3].Address()),
		sdk.AccAddress(OrchPubKeys[4].Address()),
	}

	// TODO: generate the eth priv keys here and
	// derive the address from them.

	// EthAddrs holds etheruem addresses
	EthAddrs = []gethcommon.Address{
		gethcommon.BytesToAddress(bytes.Repeat([]byte{byte(1)}, 20)),
		gethcommon.BytesToAddress(bytes.Repeat([]byte{byte(2)}, 20)),
		gethcommon.BytesToAddress(bytes.Repeat([]byte{byte(3)}, 20)),
		gethcommon.BytesToAddress(bytes.Repeat([]byte{byte(4)}, 20)),
		gethcommon.BytesToAddress(bytes.Repeat([]byte{byte(5)}, 20)),
	}

	// InitTokens holds the number of tokens to initialize an account with
	InitTokens = sdk.TokensFromConsensusPower(110, sdk.DefaultPowerReduction)

	// InitCoins holds the number of coins to initialize an account with
	InitCoins = sdk.NewCoins(sdk.NewCoin(TestingStakeParams.BondDenom, InitTokens))

	// StakingAmount holds the staking power to start a validator with
	StakingAmount = sdk.TokensFromConsensusPower(10, sdk.DefaultPowerReduction)
)

// TestInput stores the various keepers required to test gravity
type TestInput struct {
	GravityNFTKeeper  gravitynftkeeper.Keeper
	GravityKeeper     gravitykeeper.Keeper
	AccountKeeper     authkeeper.AccountKeeper
	StakingKeeper     stakingkeeper.Keeper
	SlashingKeeper    slashingkeeper.Keeper
	DistKeeper        distrkeeper.Keeper
	BankKeeper        bankkeeper.BaseKeeper
	GovKeeper         govkeeper.Keeper
	IbcTransferKeeper ibctransferkeeper.Keeper
	NftKeeper 		  nftkeeper.Keeper
	Bech32IBCKeeper   bech32ibckeeper.Keeper
	Context           sdk.Context
	Marshaler         codec.Codec
	LegacyAmino       *codec.LegacyAmino
	GravityStoreKey   *sdkstore.KVStoreKey
}

// CreateTestEnv creates the keeper testing environment for gravity
func CreateTestEnv(t *testing.T) TestInput {
	t.Helper()

	// Initialize store keys
	gravityNFTKey := sdk.NewKVStoreKey(types.StoreKey)
	gravityKey := sdk.NewKVStoreKey(gravitytypes.StoreKey)
	keyAcc := sdk.NewKVStoreKey(authtypes.StoreKey)
	keyStaking := sdk.NewKVStoreKey(stakingtypes.StoreKey)
	keyBank := sdk.NewKVStoreKey(banktypes.StoreKey)
	keyDistro := sdk.NewKVStoreKey(distrtypes.StoreKey)
	keyParams := sdk.NewKVStoreKey(paramstypes.StoreKey)
	tkeyParams := sdk.NewTransientStoreKey(paramstypes.TStoreKey)
	keyGov := sdk.NewKVStoreKey(govtypes.StoreKey)
	keySlashing := sdk.NewKVStoreKey(slashingtypes.StoreKey)
	keyCapability := sdk.NewKVStoreKey(capabilitytypes.StoreKey)
	keyUpgrade := sdk.NewKVStoreKey(upgradetypes.StoreKey)
	keyIbc := sdk.NewKVStoreKey(ibchost.StoreKey)
	keyIbcTransfer := sdk.NewKVStoreKey(ibctransfertypes.StoreKey)
	keyBech32Ibc := sdk.NewKVStoreKey(bech32ibctypes.StoreKey)
	keyNft := sdk.NewKVStoreKey(nfttypes.StoreKey)
	keyIbcNftTransfer := sdk.NewKVStoreKey(ibcnfttransfertypes.StoreKey)

	// Initialize memory database and mount stores on it
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(gravityNFTKey, sdkstore.StoreTypeIAVL, db)
	ms.MountStoreWithDB(gravityKey, sdkstore.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyAcc, sdkstore.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyParams, sdkstore.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyStaking, sdkstore.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyBank, sdkstore.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyDistro, sdkstore.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeyParams, sdkstore.StoreTypeTransient, db)
	ms.MountStoreWithDB(keyGov, sdkstore.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keySlashing, sdkstore.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyCapability, sdkstore.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyUpgrade, sdkstore.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyIbc, sdkstore.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyIbcTransfer, sdkstore.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyBech32Ibc, sdkstore.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyNft, sdkstore.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyIbcNftTransfer, sdkstore.StoreTypeIAVL, db)
	err := ms.LoadLatestVersion()
	require.Nil(t, err)

	// Create sdk.Context
	ctx := sdk.NewContext(ms, tmproto.Header{
		Version: tmversion.Consensus{
			Block: 0,
			App:   0,
		},
		ChainID: "",
		Height:  1234567,
		Time:    time.Date(2020, time.April, 22, 12, 0, 0, 0, time.UTC),
		LastBlockId: tmproto.BlockID{
			Hash: []byte{},
			PartSetHeader: tmproto.PartSetHeader{
				Total: 0,
				Hash:  []byte{},
			},
		},
		LastCommitHash:     []byte{},
		DataHash:           []byte{},
		ValidatorsHash:     []byte{},
		NextValidatorsHash: []byte{},
		ConsensusHash:      []byte{},
		AppHash:            []byte{},
		LastResultsHash:    []byte{},
		EvidenceHash:       []byte{},
		ProposerAddress:    []byte{},
	}, false, log.TestingLogger())

	cdc := MakeTestCodec()
	marshaler := MakeTestMarshaler()

	paramsKeeper := paramskeeper.NewKeeper(marshaler, cdc, keyParams, tkeyParams)
	paramsKeeper.Subspace(authtypes.ModuleName)
	paramsKeeper.Subspace(banktypes.ModuleName)
	paramsKeeper.Subspace(stakingtypes.ModuleName)
	paramsKeeper.Subspace(distrtypes.ModuleName)
	paramsKeeper.Subspace(govtypes.ModuleName)
	paramsKeeper.Subspace(gravitytypes.DefaultParamspace)
	paramsKeeper.Subspace(slashingtypes.ModuleName)
	paramsKeeper.Subspace(ibchost.ModuleName)
	paramsKeeper.Subspace(ibctransfertypes.ModuleName)
	paramsKeeper.Subspace(bech32ibctypes.ModuleName)

	// this is also used to initialize module accounts for all the map keys
	maccPerms := map[string][]string{
		authtypes.FeeCollectorName:     nil,
		distrtypes.ModuleName:          nil,
		stakingtypes.BondedPoolName:    {authtypes.Burner, authtypes.Staking},
		stakingtypes.NotBondedPoolName: {authtypes.Burner, authtypes.Staking},
		govtypes.ModuleName:            {authtypes.Burner},
		gravitytypes.ModuleName:        {authtypes.Minter, authtypes.Burner},
		ibctransfertypes.ModuleName:    {authtypes.Minter, authtypes.Burner},
		nfttypes.ModuleName:            nil,
		types.ModuleName: 				nil,
	}

	accountKeeper := authkeeper.NewAccountKeeper(
		marshaler,
		keyAcc, // target store
		getSubspace(paramsKeeper, authtypes.ModuleName),
		authtypes.ProtoBaseAccount, // prototype
		maccPerms,
		sdk.GetConfig().GetBech32AccountAddrPrefix(),
	)

	blockedAddr := make(map[string]bool, len(maccPerms))
	for acc := range maccPerms {
		blockedAddr[authtypes.NewModuleAddress(acc).String()] = true
	}
	bankKeeper := bankkeeper.NewBaseKeeper(
		marshaler,
		keyBank,
		accountKeeper,
		getSubspace(paramsKeeper, banktypes.ModuleName),
		blockedAddr,
	)
	bankKeeper.SetParams(ctx, banktypes.Params{
		SendEnabled:        []*banktypes.SendEnabled{},
		DefaultSendEnabled: true,
	})

	stakingKeeper := stakingkeeper.NewKeeper(marshaler, keyStaking, accountKeeper, bankKeeper, getSubspace(paramsKeeper, stakingtypes.ModuleName))
	stakingKeeper.SetParams(ctx, TestingStakeParams)

	distKeeper := distrkeeper.NewKeeper(marshaler, keyDistro, getSubspace(paramsKeeper, distrtypes.ModuleName), accountKeeper, bankKeeper, stakingKeeper, authtypes.FeeCollectorName)
	distKeeper.SetParams(ctx, distrtypes.DefaultParams())

	// set genesis items required for distribution
	distKeeper.SetFeePool(ctx, distrtypes.InitialFeePool())

	// set up initial accounts
	for name, perms := range maccPerms {
		mod := authtypes.NewEmptyModuleAccount(name, perms...)
		if name == distrtypes.ModuleName {
			// some big pot to pay out
			amt := sdk.NewCoins(sdk.NewInt64Coin("stake", 500000))
			err = bankKeeper.MintCoins(ctx, gravitytypes.ModuleName, amt)
			require.NoError(t, err)
			err = bankKeeper.SendCoinsFromModuleToModule(ctx, gravitytypes.ModuleName, mod.Name, amt)

			// distribution module balance must be outstanding rewards + community pool in order to pass
			// invariants checks, therefore we must add any amount we add to the module balance to the fee pool
			feePool := distKeeper.GetFeePool(ctx)
			newCoins := feePool.CommunityPool.Add(sdk.NewDecCoinsFromCoins(amt...)...)
			feePool.CommunityPool = newCoins
			distKeeper.SetFeePool(ctx, feePool)

			require.NoError(t, err)
		}
		accountKeeper.SetModuleAccount(ctx, mod)
	}

	stakeAddr := authtypes.NewModuleAddress(stakingtypes.BondedPoolName)
	moduleAcct := accountKeeper.GetAccount(ctx, stakeAddr)
	require.NotNil(t, moduleAcct)

	msgServiceRouter := baseapp.NewMsgServiceRouter()

	govRouter := govtypesv1beta1.NewRouter().
		AddRoute(paramsproposal.RouterKey, params.NewParamChangeProposalHandler(paramsKeeper)).
		AddRoute(govtypes.RouterKey, govtypesv1beta1.ProposalHandler)

	govConfig := govtypes.DefaultConfig()
	govKeeper := govkeeper.NewKeeper(
		marshaler, keyGov, getSubspace(paramsKeeper, govtypes.ModuleName).WithKeyTable(govtypesv1.ParamKeyTable()), accountKeeper, bankKeeper, stakingKeeper, govRouter,
		msgServiceRouter, govConfig,
	)

	govKeeper.SetProposalID(ctx, govtypesv1beta1.DefaultStartingProposalID)
	govKeeper.SetDepositParams(ctx, govtypesv1.DefaultDepositParams())
	govKeeper.SetVotingParams(ctx, govtypesv1.DefaultVotingParams())
	govKeeper.SetTallyParams(ctx, govtypesv1.DefaultTallyParams())

	slashingKeeper := slashingkeeper.NewKeeper(
		marshaler,
		keySlashing,
		&stakingKeeper,
		getSubspace(paramsKeeper, slashingtypes.ModuleName).WithKeyTable(slashingtypes.ParamKeyTable()),
	)

	bApp := *baseapp.NewBaseApp("test", log.TestingLogger(), db, MakeTestEncodingConfig().TxConfig.TxDecoder())
	upgradeKeeper := upgradekeeper.NewKeeper(
		make(map[int64]bool),
		keyUpgrade,
		marshaler,
		"",
		&bApp,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	memKeys := sdk.NewMemoryStoreKeys(capabilitytypes.MemStoreKey)
	capabilityKeeper := *capabilitykeeper.NewKeeper(
		marshaler,
		keyCapability,
		memKeys[capabilitytypes.MemStoreKey],
	)

	scopedIbcKeeper := capabilityKeeper.ScopeToModule(ibchost.ModuleName)
	ibcKeeper := *ibckeeper.NewKeeper(
		marshaler,
		keyIbc,
		getSubspace(paramsKeeper, ibchost.ModuleName),
		stakingKeeper,
		upgradeKeeper,
		scopedIbcKeeper,
	)

	scopedTransferKeeper := capabilityKeeper.ScopeToModule(ibctransfertypes.ModuleName)
	ibcTransferKeeper := ibctransferkeeper.NewKeeper(
		marshaler, keyIbcTransfer, getSubspace(paramsKeeper, ibctransfertypes.ModuleName),
		ibcKeeper.ChannelKeeper, ibcKeeper.ChannelKeeper, &ibcKeeper.PortKeeper,
		accountKeeper, bankKeeper, scopedTransferKeeper,
	)

	nftKeeper := nftkeeper.NewKeeper(
		keyNft, marshaler, accountKeeper, bankKeeper,
	)

	scopedNFTTransferKeeper := capabilityKeeper.ScopeToModule(ibcnfttransfertypes.ModuleName)
	nftTransferKeeper := ibcnfttransferkeeper.NewKeeper(
		marshaler, keyIbcNftTransfer, authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		ibcKeeper.ChannelKeeper, ibcKeeper.ChannelKeeper, &ibcKeeper.PortKeeper, accountKeeper,
		internft.NewInterNftKeeperWrapper(&nftKeeper), scopedNFTTransferKeeper)

	bech32IbcKeeper := *bech32ibckeeper.NewKeeper(
		ibcKeeper.ChannelKeeper, marshaler, keyBech32Ibc,
		ibcTransferKeeper, nftTransferKeeper,
	)
	// Set the native prefix to the "gravity" value we like in module/config/config.go
	err = bech32IbcKeeper.SetNativeHrp(ctx, sdk.GetConfig().GetBech32AccountAddrPrefix())
	if err != nil {
		panic("Test Env Creation failure, could not set native hrp")
	}

	gravityKeeper := gravitykeeper.NewKeeper(gravityKey, getSubspace(paramsKeeper, gravitytypes.DefaultParamspace), marshaler, &bankKeeper,
		&stakingKeeper, &slashingKeeper, &distKeeper, &accountKeeper, &ibcTransferKeeper, &bech32IbcKeeper)

	stakingKeeper = *stakingKeeper.SetHooks(
		stakingtypes.NewMultiStakingHooks(
			distKeeper.Hooks(),
			slashingKeeper.Hooks(),
			gravityKeeper.Hooks(),
		),
	)

	// set gravityIDs for batches and tx items, simulating genesis setup
	gravityKeeper.SetLatestValsetNonce(ctx, 0)
	gravityKeeper.SetLastObservedEventNonce(ctx, 0)
	gravityKeeper.SetLastSlashedValsetNonce(ctx, 0)
	gravityKeeper.SetLastSlashedBatchBlock(ctx, 0)
	gravityKeeper.SetLastSlashedLogicCallBlock(ctx, 0)
	gravityKeeper.SetID(ctx, 0, gravitytypes.KeyLastTXPoolID)
	gravityKeeper.SetID(ctx, 0, gravitytypes.KeyLastOutgoingBatchID)

	gravityKeeper.SetParams(ctx, TestingGravityParams)

	gravityNFTKeeper := gravitynftkeeper.NewKeeper(gravityNFTKey, marshaler, &gravityKeeper, &stakingKeeper, &accountKeeper,
		&bech32IbcKeeper, &nftKeeper, &nftTransferKeeper, authtypes.NewModuleAddress(govtypes.ModuleName).String())
	// TODO: Set params and any other genesis stuffs

	testInput := TestInput{
		GravityNFTKeeper: gravityNFTKeeper,
		GravityKeeper:     gravityKeeper,
		AccountKeeper:     accountKeeper,
		StakingKeeper:     stakingKeeper,
		SlashingKeeper:    slashingKeeper,
		DistKeeper:        distKeeper,
		BankKeeper:        bankKeeper,
		GovKeeper:         govKeeper,
		IbcTransferKeeper: ibcTransferKeeper,
		NftKeeper: 		   nftKeeper,
		Bech32IBCKeeper:   bech32IbcKeeper,
		Context:           ctx,
		Marshaler:         marshaler,
		LegacyAmino:       cdc,
		GravityStoreKey:   gravityKey,
	}
	// check invariants before starting
	testInput.Context.Logger().Info("Asserting invariants on new test env")
	testInput.AssertInvariants()
	return testInput
}

// MakeTestCodec creates a legacy amino codec for testing
func MakeTestCodec() *codec.LegacyAmino {
	var cdc = codec.NewLegacyAmino()
	auth.AppModuleBasic{}.RegisterLegacyAminoCodec(cdc)
	bank.AppModuleBasic{}.RegisterLegacyAminoCodec(cdc)
	staking.AppModuleBasic{}.RegisterLegacyAminoCodec(cdc)
	distribution.AppModuleBasic{}.RegisterLegacyAminoCodec(cdc)
	sdk.RegisterLegacyAminoCodec(cdc)
	ccodec.RegisterCrypto(cdc)
	params.AppModuleBasic{}.RegisterLegacyAminoCodec(cdc)
	types.RegisterCodec(cdc)
	return cdc
}

// MakeTestMarshaler creates a proto codec for use in testing
func MakeTestMarshaler() codec.Codec {
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	std.RegisterInterfaces(interfaceRegistry)
	ModuleBasics.RegisterInterfaces(interfaceRegistry)
	types.RegisterInterfaces(interfaceRegistry)
	return codec.NewProtoCodec(interfaceRegistry)
}

// getSubspace returns a param subspace for a given module name.
func getSubspace(k paramskeeper.Keeper, moduleName string) paramstypes.Subspace {
	subspace, _ := k.GetSubspace(moduleName)
	return subspace
}

func MakeTestEncodingConfig() gravityparams.EncodingConfig {
	encodingConfig := gravityparams.MakeEncodingConfig()
	std.RegisterLegacyAminoCodec(encodingConfig.Amino)
	std.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	ModuleBasics.RegisterLegacyAminoCodec(encodingConfig.Amino)
	ModuleBasics.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	return encodingConfig
}

// AssertInvariants tests each modules invariants individually, this is easier than
// dealing with all the init required to get the crisis keeper working properly by
// running appModuleBasic for every module and allowing them to register their invariants
func (t TestInput) AssertInvariants() {
	distrInvariantFunc := distrkeeper.AllInvariants(t.DistKeeper)
	bankInvariantFunc := bankkeeper.AllInvariants(t.BankKeeper)
	govInvariantFunc := govkeeper.AllInvariants(t.GovKeeper, t.BankKeeper)
	stakeInvariantFunc := stakingkeeper.AllInvariants(t.StakingKeeper)
	gravInvariantFunc := gravitykeeper.AllInvariants(t.GravityKeeper)
	// TODO: ADD IF WE ADD INVARIANTS
	// I would not reccommed spending too much time on invariants as it is getting heavily reworked in the SDK

	invariantStr, invariantViolated := distrInvariantFunc(t.Context)
	if invariantViolated {
		panic(invariantStr)
	}
	invariantStr, invariantViolated = bankInvariantFunc(t.Context)
	if invariantViolated {
		panic(invariantStr)
	}
	invariantStr, invariantViolated = govInvariantFunc(t.Context)
	if invariantViolated {
		panic(invariantStr)
	}
	invariantStr, invariantViolated = stakeInvariantFunc(t.Context)
	if invariantViolated {
		panic(invariantStr)
	}
	invariantStr, invariantViolated = gravInvariantFunc(t.Context)
	if invariantViolated {
		panic(invariantStr)
	}

	t.Context.Logger().Info("All invariants successful")
}

// SetupFiveValChain does all the initialization for a 5 Validator chain using the keys here
func SetupFiveValChain(t *testing.T) (TestInput, sdk.Context) {
	t.Helper()
	input := CreateTestEnv(t)

	// Set the params for our modules
	input.StakingKeeper.SetParams(input.Context, TestingStakeParams)

	// Initialize each of the validators
	stakingMsgSvr := stakingkeeper.NewMsgServerImpl(input.StakingKeeper)
	for i := range []int{0, 1, 2, 3, 4} {

		// Initialize the account for the key
		acc := input.AccountKeeper.NewAccount(
			input.Context,
			authtypes.NewBaseAccount(AccAddrs[i], AccPubKeys[i], uint64(i), 0),
		)

		// Set the balance for the account
		require.NoError(t, input.BankKeeper.MintCoins(input.Context, gravitytypes.ModuleName, InitCoins))
		require.NoError(t, input.BankKeeper.SendCoinsFromModuleToAccount(input.Context, gravitytypes.ModuleName, acc.GetAddress(), InitCoins))

		// Set the account in state
		input.AccountKeeper.SetAccount(input.Context, acc)

		// Create a validator for that account using some of the tokens in the account
		// and the staking handler
		_, err := stakingMsgSvr.CreateValidator(
			input.Context,
			NewTestMsgCreateValidator(ValAddrs[i], ConsPubKeys[i], StakingAmount),
		)

		// Return error if one exists
		require.NoError(t, err)
	}

	// Run the staking endblocker to ensure valset is correct in state
	staking.EndBlocker(input.Context, input.StakingKeeper)

	// Register eth addresses and orchestrator address for each validator
	for i, addr := range ValAddrs {
		ethAddr, err := gravitytypes.NewEthAddress(EthAddrs[i].String())
		if err != nil {
			panic("found invalid address in EthAddrs")
		}
		input.GravityKeeper.SetEthAddressForValidator(input.Context, addr, *ethAddr)

		input.GravityKeeper.SetOrchestratorValidator(input.Context, addr, OrchAddrs[i])
	}

	// Return the test input
	return input, input.Context
}

func NewTestMsgCreateValidator(address sdk.ValAddress, pubKey ccrypto.PubKey, amt sdkmath.Int) *stakingtypes.MsgCreateValidator {
	commission := stakingtypes.NewCommissionRates(sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec())
	out, err := stakingtypes.NewMsgCreateValidator(
		address, pubKey, sdk.NewCoin("stake", amt),
		stakingtypes.Description{
			Moniker:         "",
			Identity:        "",
			Website:         "",
			SecurityContact: "",
			Details:         "",
		}, commission, sdk.OneInt(),
	)
	if err != nil {
		panic(err)
	}
	return out
}