package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/slices"
	"testing"
)

func TestMsgSendNFTToCosmosClaimValidateBasic(t *testing.T) {
	var (
		exampleSender               = NonemptyEthAddress()
		exampleTokenContract        = NonemptyEthAddress()
		exampleReceiver             = NonemptySdkAccAddress().String()
		exampleOrchestrator         = NonemptySdkAccAddress().String()
		exampleNonce         uint64 = NonzeroUint64()
	)

	tests := map[string]struct {
		ethereumSender   string
		tokenContract    string
		receiver         string
		orchestrator     string
		expErrorContains string
		nonce            uint64
	}{
		"happy path": {
			ethereumSender:   exampleSender,
			tokenContract:    exampleTokenContract,
			receiver:         exampleReceiver,
			orchestrator:     exampleOrchestrator,
			nonce:            exampleNonce,
			expErrorContains: "",
		},
		"invalid receiver is fine": {
			ethereumSender:   exampleSender,
			tokenContract:    exampleTokenContract,
			receiver:         "invalid",
			orchestrator:     exampleOrchestrator,
			nonce:            exampleNonce,
			expErrorContains: "",
		},
		"invalid ethereum sender": {
			ethereumSender:   "invalid",
			tokenContract:    exampleTokenContract,
			receiver:         exampleReceiver,
			orchestrator:     exampleOrchestrator,
			nonce:            exampleNonce,
			expErrorContains: "eth sender: invalid hex",
		},
		"invalid nft contract": {
			ethereumSender:   exampleSender,
			tokenContract:    "invalid",
			receiver:         exampleReceiver,
			orchestrator:     exampleOrchestrator,
			nonce:            exampleNonce,
			expErrorContains: "nft contract address: invalid hex",
		},
		"invalid orchestrator": {
			ethereumSender:   exampleSender,
			tokenContract:    exampleTokenContract,
			receiver:         exampleReceiver,
			orchestrator:     "invalid",
			nonce:            exampleNonce,
			expErrorContains: "orchestrator: invalid address",
		},
		"invalid nonce": {
			ethereumSender:   exampleSender,
			tokenContract:    exampleTokenContract,
			receiver:         exampleReceiver,
			orchestrator:     exampleOrchestrator,
			nonce:            0,
			expErrorContains: "nonce == 0",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			msg := MsgSendNFTToCosmosClaim{
				EventNonce:     tc.nonce,
				EthBlockHeight: 0,
				TokenContract:  tc.tokenContract,
				TokenId:        "1",
				TokenUri:       "1",
				EthereumSender: tc.ethereumSender,
				CosmosReceiver: tc.receiver,
				Orchestrator:   tc.orchestrator,
			}
			err := msg.ValidateBasic()
			if tc.expErrorContains != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expErrorContains)
				return
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// Ensures that ClaimHash changes when members of MsgSendNFTToCosmosClaim change
// The only field which MUST NOT affect ClaimHash is Orchestrator
func TestMsgSendNFTToCosmosClaimHash(t *testing.T) {
	base := MsgSendNFTToCosmosClaim{
		EventNonce:     0,
		EthBlockHeight: 0,
		TokenContract:  "",
		TokenId:        "",
		TokenUri:       "",
		EthereumSender: "",
		CosmosReceiver: "",
		Orchestrator:   "",
	}

	// Copy and populate base with values, saving orchestrator for a special check
	orchestrator := NonemptySdkAccAddress()
	mNonce := base
	mNonce.EventNonce = NonzeroUint64()
	mBlock := base
	mBlock.EthBlockHeight = NonzeroUint64()
	mCtr := base
	mCtr.TokenContract = NonemptyEthAddress()
	mId := base
	mId.TokenId = "5"
	mUri := base
	mUri.TokenUri = "https://www.gravitybridge.net"
	mSend := base
	mSend.EthereumSender = NonemptyEthAddress()
	mRecv := base
	mRecv.CosmosReceiver = NonemptySdkAccAddress().String()

	hashes := getClaimHashStrings(t, &base, &mNonce, &mBlock, &mCtr, &mId, &mUri, &mSend, &mRecv)
	baseH := hashes[0]
	rest := hashes[1:]
	// Assert that the base claim hash differs from all the rest
	require.False(t, slices.Contains(rest, baseH))

	newClaims := setOrchestratorOnClaims(orchestrator, &base, &mNonce, &mBlock, &mCtr, &mId, &mUri, &mSend, &mRecv)
	newHashes := getClaimHashStrings(t, newClaims...)
	// Assert that the claims with orchestrator set do not change the hashes
	require.Equal(t, hashes, newHashes)
}

// Gets the ClaimHash() output from every claims member and casts it to a string, panicing on any errors
func getClaimHashStrings(t *testing.T, claims ...EthereumNFTClaim) (hashes []string) {
	for _, claim := range claims {
		hash, e := claim.ClaimHash()
		require.NoError(t, e)
		hashes = append(hashes, string(hash))
	}
	return
}

// Calls SetOrchestrator on every claims member, passing orch as the value
func setOrchestratorOnClaims(orch sdk.AccAddress, claims ...EthereumNFTClaim) (ret []EthereumNFTClaim) {
	for _, claim := range claims {
		clam := claim
		clam.SetOrchestrator(orch)
		ret = append(ret, clam)
	}
	return
}
