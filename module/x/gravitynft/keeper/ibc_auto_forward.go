package keeper

import (
	"fmt"
	gravitytypes "github.com/Gravity-Bridge/Gravity-Bridge/module/x/gravity/types"
	"github.com/Gravity-Bridge/Gravity-Bridge/module/x/gravitynft/types"
	bech32ibctypes "github.com/althea-net/bech32-ibc/x/bech32ibc/types"
	ibcnfttransfertypes "github.com/bianjieai/nft-transfer/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/nft"
	ibctransfertypes "github.com/cosmos/ibc-go/v5/modules/apps/transfer/types"
	ibcclienttypes "github.com/cosmos/ibc-go/v5/modules/core/02-client/types"
	"time"
)

// PendingIbcAutoForwards returns an ordered slice of the queued IBC Auto-Forward sends to IBC-enabled chains
func (k Keeper) PendingNFTIbcAutoForwards(ctx sdk.Context, limit uint64) []*types.PendingNFTIbcAutoForward {
	forwards := make([]*types.PendingNFTIbcAutoForward, 0)

	k.IteratePendingNFTIbcAutoForwards(ctx, func(key []byte, forward *types.PendingNFTIbcAutoForward) (stop bool) {
		forwards = append(forwards, forward)
		if limit != 0 && uint64(len(forwards)) >= limit {
			return true
		}
		return false
	})

	return forwards
}

// IteratePendingNFTIbcAutoForwards executes the given callback on each PendingIbcAutoForward in the store
// cb should return true to stop iteration, false to continue
func (k Keeper) IteratePendingNFTIbcAutoForwards(ctx sdk.Context, cb func(key []byte, forward *types.PendingNFTIbcAutoForward) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	prefixStore := prefix.NewStore(store, types.PendingNFTIbcAutoForwards)
	iter := prefixStore.Iterator(nil, nil)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		forward := new(types.PendingNFTIbcAutoForward)
		k.cdc.MustUnmarshal(iter.Value(), forward)

		if cb(iter.Key(), forward) {
			break
		}
	}
}

func (k Keeper) AddPendingNFTPendingIbcAutoForward(ctx sdk.Context, forward types.PendingNFTIbcAutoForward) error {
	if err := k.ValidatePendingERC721IbcAutoForward(ctx, forward); err != nil {
		return err
	}

	store := ctx.KVStore(k.storeKey)
	key := types.GetPendingNFTIbcAutoForwardKey(forward.EventNonce)
	if store.Has(key) {
		return sdkerrors.Wrapf(types.ErrDuplicate,
			"Pending IBC NFT Auto-Forward Queue already has an entry with nonce %v", forward.EventNonce,
		)
	}
	store.Set(key, k.cdc.MustMarshal(&forward))

	k.logger(ctx).Info("SendToCosmos Pending IBC Auto-Forward", "ibcReceiver", forward.ForeignReceiver,
		"class-id", forward.ClassId, "token-id", forward.TokenId,
		"ibc-port", k.ibcNftTransferKeeper.GetPort(ctx), "ibcChannel", forward.IbcChannel, "claimNonce", forward.EventNonce,
		"cosmosBlockTime", ctx.BlockTime(), "cosmosBlockHeight", ctx.BlockHeight(),
	)

	return ctx.EventManager().EmitTypedEvent(&types.EventSendNFTToCosmosPendingIbcAutoForward{
		Nonce:    fmt.Sprint(forward.EventNonce),
		Receiver: forward.ForeignReceiver,
		ClassId:  forward.ClassId,
		TokenId:  forward.TokenId,
		Channel:  forward.IbcChannel,
	})
}

// ValidatePendingERC721IbcAutoForward performs basic validation, asserts the nonce is not ahead of what gravity is aware of,
// requires ForeignReceiver's bech32 prefix to be registered and match with IbcChannel, and gravity module must have the
// funds to meet this forward amount
func (k Keeper) ValidatePendingERC721IbcAutoForward(ctx sdk.Context, forward types.PendingNFTIbcAutoForward) error {
	if err := forward.ValidateBasic(); err != nil {
		return err
	}

	latestEventNonce := k.GetLastObservedEventNonce(ctx)
	if forward.EventNonce > latestEventNonce {
		return sdkerrors.Wrap(types.ErrInvalid, "EventNonce must be <= latest observed event nonce")
	}
	prefix, _, err := bech32.DecodeAndConvert(forward.ForeignReceiver)
	if err != nil { // Covered by ValidateBasic, but check anyway to avoid linter issues
		return sdkerrors.Wrapf(err, "ForeignReceiver %s is not a valid bech32 address", forward.ForeignReceiver)
	}
	hrpRecord, err := k.bech32IbcKeeper.GetHrpIbcRecord(ctx, prefix)
	if err != nil {
		return sdkerrors.Wrapf(bech32ibctypes.ErrInvalidHRP, "ForeignReciever %s has an invalid or unregistered prefix: %s", forward.ForeignReceiver, prefix)
	}
	if forward.IbcChannel != hrpRecord.NftSourceChannel {
		return sdkerrors.Wrapf(types.ErrMismatched, "IbcChannel %s does not match the registered prefix's IBC channel %v",
			forward.IbcChannel, hrpRecord.String(),
		)
	}
	modAcc := k.accountKeeper.GetModuleAccount(ctx, types.ModuleName).GetAddress()
	owner := k.nftKeeper.GetOwner(ctx, forward.ClassId, forward.TokenId)
	if !owner.Equals(modAcc) {
		return sdkerrors.Wrapf(
			sdkerrors.ErrInsufficientFunds, "gravitynft module account does not have nft token %s from class id %s",
			forward.TokenId, forward.ClassId,
		)
	}

	return nil
}

// ProcessPendingIbcNFTAutoForwards processes and dequeues many pending IBC NFT Auto-Forwards, either sending the NFTs to their
// respective destination chains or on error sending the funds to the local gravitynft-prefixed account
// See ProcessNextPendingIbcNFTAutoForward for more details
func (k Keeper) ProcessPendingIbcNFTAutoForwards(ctx sdk.Context, forwardsToClear uint64) error {
	for i := uint64(0); i < forwardsToClear; i++ {
		stop, err := k.ProcessNextPendingIbcNFTAutoForward(ctx)

		if err != nil {
			return sdkerrors.Wrapf(err, "unable to process Pending IBC Auto-Forward number %v", i)
		}
		if stop {
			break
		}
	}

	return nil
}

func (k Keeper) ProcessNextPendingIbcNFTAutoForward(ctx sdk.Context) (stop bool, err error) {
	forward := k.GetNextPendingIbcNFTAutoForward(ctx)
	if forward == nil {
		return true, nil // No forwards to process, exit early
	}
	if err := forward.ValidateBasic(); err != nil { // double-check the forward before sending it
		// Fail this tx
		panic(fmt.Sprintf("Invalid forward found in Pending IBC Auto-Forward queue: %s", err.Error()))
	}
	// Point of no return: the funds will be sent somewhere, either the IBC address, local address or the community pool
	err = k.deletePendingIbcNFTAutoForward(ctx, forward.EventNonce)
	if err != nil {
		// Fail this tx
		panic(fmt.Sprintf("Discovered nonexistent Pending IBC Auto-Forward in the queue %s", forward.String()))
	}

	portId := k.ibcNftTransferKeeper.GetPort(ctx)

	// This local gravity user receives the nft if the ibc transaction fails
	var fallback sdk.AccAddress
	fallback, err = gravitytypes.IBCAddressFromBech32(forward.ForeignReceiver)
	if err != nil {
		panic(fmt.Sprintf("Invalid ForeignReceiver found in Pending IBC Auto-Forward queue: %s [[%+v]]", err.Error(), forward))
	}

	if err := k.nftKeeper.Transfer(ctx, forward.ClassId, forward.TokenId, fallback); err != nil {
		nftToken := nft.NFT{
			Id:      forward.TokenId,
			ClassId: forward.ClassId,
		}
		return false, k.SendNFTToCommunityPool(ctx, nftToken)
	}

	timeoutTime := thirtyDaysInFuture(ctx) // Set the ibc transfer to expire ~one month from now
	msgTransfer := createERC721IbcMsgTransfer(portId, *forward, fallback.String(), uint64(timeoutTime.UnixNano()))

	// Make the ibc-transfer attempt
	wCtx := sdk.WrapSDKContext(ctx)
	_, recoverableErr := k.ibcNftTransferKeeper.Transfer(wCtx, &msgTransfer)
	ctx = sdk.UnwrapSDKContext(wCtx)

	// Log + emit event
	if recoverableErr == nil {
		k.logEmitNFTIbcForwardSuccessEvent(ctx, *forward, msgTransfer)
	} else {
		// NFT have already been sent to the fallback user, emit a failure log
		k.logEmitNFTIbcForwardFailureEvent(ctx, *forward, recoverableErr)
	}
	return false, nil // Error case has been handled, funds are in receiver's control locally or on IBC chain
}

// GetNextPendingERC721IbcAutoForward returns the first pending IBC Auto-Forward in the queue
func (k Keeper) GetNextPendingIbcNFTAutoForward(ctx sdk.Context) *types.PendingNFTIbcAutoForward {
	store := ctx.KVStore(k.storeKey)
	iter := store.Iterator(prefixRange(types.PendingNFTIbcAutoForwards))
	defer iter.Close()
	if iter.Valid() {
		var forward types.PendingNFTIbcAutoForward
		k.cdc.MustUnmarshal(iter.Value(), &forward)

		return &forward
	}
	return nil
}

// deletePendingIbcNFTAutoForward removes a single pending IBC Auto-Forward send to an IBC-enabled chain from the store
// WARNING: this should only be called while clearing the queue in ClearNextPendingIbcAutoForward
func (k Keeper) deletePendingIbcNFTAutoForward(ctx sdk.Context, eventNonce uint64) error {
	store := ctx.KVStore(k.storeKey)
	key := types.GetPendingNFTIbcAutoForwardKey(eventNonce)
	if !store.Has(key) {
		return sdkerrors.Wrapf(types.ErrInvalid, "No PendingIbcAutoForward with nonce %v in the store", eventNonce)
	}
	store.Delete(key)
	return nil
}

func createERC721IbcMsgTransfer(portId string, forward types.PendingNFTIbcAutoForward, sender string, timeoutTimestampNs uint64) ibcnfttransfertypes.MsgTransfer {
	zeroHeight := ibcclienttypes.Height{}
	return *ibcnfttransfertypes.NewMsgTransfer(
		portId,
		forward.IbcChannel,
		forward.ClassId,
		[]string{forward.TokenId},
		sender,
		forward.ForeignReceiver,
		zeroHeight, // Do not use block height based timeout
		timeoutTimestampNs,
		"NFT IBC Auto-Forwarded by Gravity Bridge",
	)
}

// thirtyDaysInFuture creates a time.Time exactly 30 days from the last BlockTime for use in createIbcMsgTransfer
func thirtyDaysInFuture(ctx sdk.Context) time.Time {
	approxNow := ctx.BlockTime()
	// Get the offset from zero of 30 days in the future
	return approxNow.Add(time.Hour * 24 * 30)
}

// logEmitNFTIbcForwardSuccessEvent logs for successful IBC Auto-Forwarding and emits a
// EventSendERC721ToCosmosExecutedIbcAutoForward type event
func (k Keeper) logEmitNFTIbcForwardSuccessEvent(
	ctx sdk.Context,
	forward types.PendingNFTIbcAutoForward,
	msgTransfer ibcnfttransfertypes.MsgTransfer,
) {
	k.logger(ctx).Info("SendERC721ToCosmos IBC Auto-Forward", "ibcReceiver", forward.ForeignReceiver,
		"class-id", forward.ClassId, "token-id", forward.TokenId, "ibc-port", msgTransfer.SourcePort, "ibcChannel", forward.IbcChannel,
		"timeoutHeight", msgTransfer.TimeoutHeight.String(), "timeoutTimestamp", msgTransfer.TimeoutTimestamp,
		"claimNonce", forward.EventNonce, "cosmosBlockHeight", ctx.BlockHeight(),
	)

	ctx.EventManager().EmitTypedEvent(&types.EventSendNFTToCosmosExecutedIbcAutoForward{
		Nonce:         fmt.Sprint(forward.EventNonce),
		Receiver:      forward.ForeignReceiver,
		ClassId:       forward.ClassId,
		TokenId:       forward.TokenId,
		Channel:       forward.IbcChannel,
		TimeoutTime:   fmt.Sprint(msgTransfer.TimeoutTimestamp),
		TimeoutHeight: msgTransfer.TimeoutHeight.String(),
	})
}

// logEmitNFTIbcForwardFailureEvent logs failed IBC Auto-Forwarding and emits a EventSendToCosmosLocal type event
func (k Keeper) logEmitNFTIbcForwardFailureEvent(ctx sdk.Context, forward types.PendingNFTIbcAutoForward, err error) {
	var localReceiver sdk.AccAddress
	localReceiver, _ = gravitytypes.IBCAddressFromBech32(forward.ForeignReceiver) // checked valid bech32 receiver earlier
	k.logger(ctx).Error("SendToCosmos IBC Auto-Forward Failure: funds sent to local address",
		"localReceiver", localReceiver, "class-id", forward.ClassId, "token-id", forward.TokenId,
		"failedIbcPort", ibctransfertypes.PortID, "failedIbcChannel", forward.IbcChannel,
		"claimNonce", forward.EventNonce, "cosmosBlockHeight", ctx.BlockHeight(), "err", err,
	)

	ctx.EventManager().EmitTypedEvent(&types.EventSendNFTToCosmosLocal{
		Nonce:    fmt.Sprint(forward.EventNonce),
		Receiver: forward.ForeignReceiver,
		ClassId:  forward.ClassId,
		TokenId:  forward.TokenId,
	})
}
