package keeper

import (
	"fmt"
	"github.com/Gravity-Bridge/Gravity-Bridge/module/x/gravitynft/types"
	bech32ibctypes "github.com/althea-net/bech32-ibc/x/bech32ibc/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

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