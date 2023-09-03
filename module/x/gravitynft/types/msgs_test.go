package types

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMsgSendNFTToCosmosClaimValidateBasic(t *testing.T) {
	var (
		exampleSender        = "0xb462864E395d88d6bc7C5dd5F3F5eb4cc2599255"
		exampleTokenContract = "0x71C7656EC7ab88b098defB751B7401B5f6d8976F"
		exampleReceiver      = "gravity1dh3pl2e9swm478h085yt2tqxnvd8gxd2r3mp47"
		exampleOrchestrator  = "gravity1agdhen5mlf683ruat6egarzljpsceutjllnf7n"
	)

	tests := map[string]struct {
		ethereumSender   string
		tokenContract    string
		receiver         string
		orchestrator     string
		expErrorContains string
	}{
		"happy path": {
			ethereumSender:   exampleSender,
			tokenContract:    exampleTokenContract,
			receiver:         exampleReceiver,
			orchestrator:     exampleOrchestrator,
			expErrorContains: "",
		},
		"invalid ethereum sender": {
			ethereumSender:   "invalid",
			tokenContract:    exampleTokenContract,
			receiver:         exampleReceiver,
			orchestrator:     exampleOrchestrator,
			expErrorContains: "eth sender: invalid hex",
		},
		"invalid nft contract": {
			ethereumSender:   exampleSender,
			tokenContract:    "invalid",
			receiver:         exampleReceiver,
			orchestrator:     exampleOrchestrator,
			expErrorContains: "nft contract address: invalid hex",
		},
		// TODO: Not sure if destination should be validated or not
		"invalid orchestrator": {
			ethereumSender:   exampleSender,
			tokenContract:    exampleTokenContract,
			receiver:         exampleReceiver,
			orchestrator:     "invalid",
			expErrorContains: "orchestrator: invalid address",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			msg := MsgSendNFTToCosmosClaim{
				EventNonce:     0,
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
			}
		})
	}
}
