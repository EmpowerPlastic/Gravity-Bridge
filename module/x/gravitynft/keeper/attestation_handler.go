package keeper

import (
	sdkerrors "cosmossdk.io/errors"
	"fmt"
	gravitytypes "github.com/Gravity-Bridge/Gravity-Bridge/module/x/gravity/types"
	"github.com/Gravity-Bridge/Gravity-Bridge/module/x/gravitynft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/nft"
	"strconv"
	"strings"
)

// AttestationHandler processes `observed` Attestations
type AttestationHandler struct {
	// NOTE: If you add anything to this struct, add a nil check to ValidateMembers below!
	keeper *Keeper
}

// Check for nil members
func (a AttestationHandler) ValidateMembers() {
	if a.keeper == nil {
		panic("Nil keeper!")
	}
}

// Handle is the entry point for Attestation processing, only attestations with sufficient validator submissions
// should be processed through this function, solidifying their effect in chain state
func (a AttestationHandler) Handle(ctx sdk.Context, att types.NFTAttestation, claim types.EthereumNFTClaim) error {
	switch claim := claim.(type) {

	case *types.MsgSendNFTToCosmosClaim:
		return a.handleSendNFTToCosmos(ctx, *claim)

	// TODO: Add other claim types (like NFT_CLAIM_TYPE_ERC721_DEPLOYED)

	default:
		panic(fmt.Sprintf("Invalid event type for attestations %s", claim.GetType()))
	}
}

// Upon acceptance of sufficient validator SendNFTToCosmosClaim claims: transfer NFT to the appropriate cosmos account
// The cosmos receiver can be a native account (e.g. gravity1abc...) or a foreign account (e.g. cosmos1abc...)
// In the event of a native receiver, x/nft module handles the transfer, otherwise an IBC transfer is initiated
func (a AttestationHandler) handleSendNFTToCosmos(ctx sdk.Context, claim types.MsgSendNFTToCosmosClaim) error {
	invalidAddress := false
	// Validate the receiver as a valid bech32 address
	receiverAddress, addressErr := gravitytypes.IBCAddressFromBech32(claim.CosmosReceiver)

	if addressErr != nil {
		invalidAddress = true
		hash, er := claim.ClaimHash()
		if er != nil {
			return sdkerrors.Wrapf(er, "Unable to log error %v, could not compute ClaimHash for claim %v: %v", addressErr, claim, er)
		}

		a.keeper.logger(ctx).Error("Invalid SendNFTToCosmos receiver",
			"address", receiverAddress,
			"cause", addressErr.Error(),
			"claim type", claim.GetType(),
			"id", types.GetAttestationKey(claim.GetEventNonce(), hash),
			"nonce", fmt.Sprint(claim.GetEventNonce()),
		)
	}

	tokenAddress, errTokenAddress := gravitytypes.NewEthAddress(claim.TokenContract)
	ethereumSender, errEthereumSender := gravitytypes.NewEthAddress(claim.EthereumSender)
	// nil address is not possible unless the validators get together and submit
	// a bogus event, this would create lost tokens stuck in the bridge
	// and not accessible to anyone
	if errTokenAddress != nil {
		hash, er := claim.ClaimHash()
		if er != nil {
			return sdkerrors.Wrapf(er, "Unable to log error %v, could not compute ClaimHash for claim %v: %v", errTokenAddress, claim, er)
		}
		a.keeper.logger(ctx).Error("Invalid token contract",
			"cause", errTokenAddress.Error(),
			"claim type", claim.GetType(),
			"id", types.GetAttestationKey(claim.GetEventNonce(), hash),
			"nonce", fmt.Sprint(claim.GetEventNonce()),
		)
		return sdkerrors.Wrap(errTokenAddress, "invalid token contract on claim")
	}
	// likewise nil sender would have to be caused by a bogus event
	if errEthereumSender != nil {
		hash, er := claim.ClaimHash()
		if er != nil {
			return sdkerrors.Wrapf(er, "Unable to log error %v, could not compute ClaimHash for claim %v: %v", errEthereumSender, claim, er)
		}
		a.keeper.logger(ctx).Error("Invalid ethereum sender",
			"cause", errEthereumSender.Error(),
			"claim type", claim.GetType(),
			"id", types.GetAttestationKey(claim.GetEventNonce(), hash),
			"nonce", fmt.Sprint(claim.GetEventNonce()),
		)
		return sdkerrors.Wrap(errTokenAddress, "invalid ethereum sender on claim")
	}

	// TODO: IS THIS OK? UISING THE GRAVITY MODULE'S BLACKLIST?
	// Block blacklisted asset transfers using the gravity module's black list
	// (these funds are unrecoverable for the blacklisted sender, they will instead be sent to community pool)
	if a.keeper.gravityKeeper.IsOnBlacklist(ctx, *ethereumSender) {
		hash, er := claim.ClaimHash()
		if er != nil {
			return sdkerrors.Wrapf(er, "Unable to log blacklisted error, could not compute ClaimHash for claim %v: %v", claim, er)
		}
		a.keeper.logger(ctx).Error("Invalid SendToCosmos: receiver is blacklisted",
			"address", receiverAddress,
			"claim type", claim.GetType(),
			"id", types.GetAttestationKey(claim.GetEventNonce(), hash),
			"nonce", fmt.Sprint(claim.GetEventNonce()),
		)
		invalidAddress = true
	}

	// Check if coin is Cosmos-originated asset and get denom
	isCosmosOriginated, classID := a.keeper.ERC721ToClassIDLookup(ctx, *tokenAddress)
	nftToken := nft.NFT{
		ClassId: classID,
		Id:      claim.TokenId,
		Uri:     claim.TokenUri,
	}
	moduleAddr := a.keeper.accountKeeper.GetModuleAddress(types.ModuleName)
	if !isCosmosOriginated { // We need to mint eth-originated NFT (aka vouchers)
		if err := a.mintEthereumOriginatedNFTVouchers(ctx, moduleAddr, claim, nftToken); err != nil {
			// TODO: Evaluate closely, if we can't mint an ethereum voucher, what should we do?
			return err
		}
	}

	if !invalidAddress { // address appears valid, attempt to send minted/locked coins to receiver
		// Failure to send will result in NFT transfer to community pool
		ibcForwardQueued, err := a.sendNFTToCosmosAccount(ctx, claim, receiverAddress, nftToken)
		_ = ibcForwardQueued

		// TODO: Add this stuff type of  later for production
		// TODO: Perform module balance assertions
		/*if err != nil || ibcForwardQueued { // ibc forward enqueue and errors should not send tokens to anyone
			a.assertNothingSent(ctx, moduleAddr, preSendBalance, denom)
		} else { // No error, local send -> assert send had right amount
			a.assertSentAmount(ctx, moduleAddr, preSendBalance, denom, claim.Amount)
		}*/

		if err != nil { // trigger send to community pool
			invalidAddress = true
		}
	}

	// for whatever reason above, blacklisted, invalid string, etc this deposit is not valid
	// we can't send the tokens back on the Ethereum side, and if we don't put them somewhere on
	// the cosmos side they will be lost an inaccessible even though they are locked in the bridge.
	// so we deposit the tokens into the community pool for later use via governance vote
	if invalidAddress {
		if err := a.keeper.SendNFTToCommunityPool(ctx, nftToken); err != nil {
			hash, _ := claim.ClaimHash()
			a.keeper.logger(ctx).Error("Failed community pool send",
				"cause", err.Error(),
				"claim type", claim.GetType(),
				"id", types.GetAttestationKey(claim.GetEventNonce(), hash),
				"nonce", fmt.Sprint(claim.GetEventNonce()),
			)
			return sdkerrors.Wrap(err, "failed to send to Community pool")
		}

		if err := ctx.EventManager().EmitTypedEvent(
			&types.EventInvalidSendNFTToCosmosReceiver{
				Contract: tokenAddress.GetAddress().Hex(),
				ClassId:  nftToken.ClassId,
				TokenId:  nftToken.Id,
				Nonce:    strconv.Itoa(int(claim.GetEventNonce())),
				Sender:   claim.EthereumSender,
			},
		); err != nil {
			return err
		}

	} else {
		if err := ctx.EventManager().EmitTypedEvent(
			&types.EventSendNFTToCosmos{
				Contract: claim.TokenContract,
				ClassId:  nftToken.ClassId,
				TokenId:  nftToken.Id,
				Nonce:    strconv.Itoa(int(claim.GetEventNonce())),
			},
		); err != nil {
			return err
		}
	}

	return nil
}

func (a AttestationHandler) mintEthereumOriginatedNFTVouchers(
	ctx sdk.Context, moduleAddr sdk.AccAddress, claim types.MsgSendNFTToCosmosClaim, nftToken nft.NFT,
) error {
	hasClass := a.keeper.nftKeeper.HasClass(ctx, nftToken.ClassId)
	if !hasClass {
		class := nft.Class{
			Id: nftToken.ClassId,
		}
		if err := a.keeper.nftKeeper.SaveClass(ctx, class); err != nil {
			// in this case we have lost the NFT! They are in the bridge, but not
			// in the community pool or out in some users balance, every instance of this
			// error needs to be detected and resolved
			hash, _ := claim.ClaimHash()
			a.keeper.logger(ctx).Error("Failed creating NFT class",
				"cause", err.Error(),
				"claim type", claim.GetType(),
				"id", types.GetAttestationKey(claim.GetEventNonce(), hash),
				"nonce", fmt.Sprint(claim.GetEventNonce()),
			)
			return sdkerrors.Wrapf(err, "mint vouchers NFT: %s s", nftToken.ClassId, nftToken.Id)
		}
	}

	if err := a.keeper.nftKeeper.Mint(ctx, nftToken, moduleAddr); err != nil {
		// in this case we have lost the NFT! They are in the bridge, but not
		// in the community pool or out in some users balance, every instance of this
		// error needs to be detected and resolved
		hash, _ := claim.ClaimHash()
		a.keeper.logger(ctx).Error("Failed minting NFT",
			"cause", err.Error(),
			"claim type", claim.GetType(),
			"id", types.GetAttestationKey(claim.GetEventNonce(), hash),
			"nonce", fmt.Sprint(claim.GetEventNonce()),
		)
		return sdkerrors.Wrapf(err, "mint vouchers NFT: %s s", nftToken.ClassId, nftToken.Id)
	}

	return nil
}

func (a AttestationHandler) sendNFTToCosmosAccount(ctx sdk.Context, claim types.MsgSendNFTToCosmosClaim, receiver sdk.AccAddress, nftToken nft.NFT) (ibcForwardQueued bool, err error) {
	accountPrefix, err := gravitytypes.GetPrefixFromBech32(claim.CosmosReceiver)
	if err != nil {
		hash, _ := claim.ClaimHash()
		a.keeper.logger(ctx).Error("Invalid bech32 CosmosReceiver",
			"cause", err.Error(), "address", receiver,
			"claimType", claim.GetType(),
			"id", types.GetAttestationKey(claim.GetEventNonce(), hash),
			"nonce", fmt.Sprint(claim.GetEventNonce()),
		)
		return false, err
	}
	nativePrefix, err := a.keeper.bech32IbcKeeper.GetNativeHrp(ctx)
	if err != nil {
		// In a real environment bech32ibc panics on InitGenesis and on Send with their bech32ics20 module, which
		// prevents all MsgSend + MsgMultiSend transfers, in a testing environment it is possible to hit this condition,
		// so we should panic as well. This will cause a chain halt, and prevent attestation handling until prefix is set
		panic("SendToCosmos failure: bech32ibc NativeHrp has not been set!")
	}

	if accountPrefix == nativePrefix { // Send to a native gravity account
		return false, a.sendNFTToLocalAddress(ctx, claim, receiver, nftToken)
	} else { // Try to send tokens to IBC chain, fall back to native send on errors
		// TODO: Implement this properly as soon as we have fixed the bech32-ibc implementation
		// panic("not implemented just yet! need to figure out if I can skip that pesky ibc-forwarding queue crap")
		// The hrp hack here is ugly as shit, but the hrp thing doesn't seem to allow more than one ibc channel per account prefix...
		// stars... -> IBC channel ICS20
		// nftstars -> IBC channel ICS721
		hrpPrefix := types.AccountPrefixForERC721Hrp(accountPrefix)
		hrpIbcRecord, err := a.keeper.bech32IbcKeeper.GetHrpIbcRecord(ctx, hrpPrefix)
		if err != nil {
			hash, _ := claim.ClaimHash()
			a.keeper.logger(ctx).Error("Unregistered foreign prefix",
				"cause", err.Error(),
				"address", receiver,
				"hrp prefix", hrpPrefix,
				"claim type", claim.GetType(),
				"id", types.GetAttestationKey(claim.GetEventNonce(), hash),
				"nonce", fmt.Sprint(claim.GetEventNonce()),
			)

			// Fall back to sending tokens to native account
			return false, sdkerrors.Wrap(
				a.sendNFTToLocalAddress(ctx, claim, receiver, nftToken),
				"Unregistered foreign prefix, send via x/nft",
			)
		}

		// Add the SendERC721ToCosmos to the Pending IBC Auto-Forward Queue, which when processed will send the funds to a
		// local address before sending via IBC
		err = a.addERC721ToIbcAutoForwardQueue(ctx, receiver, accountPrefix, nftToken, hrpIbcRecord.SourceChannel, claim)

		if err != nil {
			a.keeper.logger(ctx).Error(
				"SendERC721ToCosmos IBC auto forwarding failed, sending to local gravity account instead",
				"cosmos-receiver", claim.CosmosReceiver, "class-id", nftToken.ClassId, "token-id", nftToken.Id,
				"ethereum-contract", claim.TokenContract, "sender", claim.EthereumSender, "event-nonce", claim.EventNonce,
			)
			// Fall back to sending tokens to native account
			return false, sdkerrors.Wrap(
				a.sendNFTToLocalAddress(ctx, claim, receiver, nftToken),
				"Unregistered foreign prefix, send via x/nft",
			)
		}
		return true, nil
	}
}

func (a AttestationHandler) sendNFTToLocalAddress(
	ctx sdk.Context, claim types.MsgSendNFTToCosmosClaim, receiver sdk.AccAddress, nftToken nft.NFT,
) (err error) {
	err = a.keeper.nftKeeper.Transfer(ctx, nftToken.ClassId, nftToken.Id, receiver)
	if err != nil {
		// Well, that can't be good
		hash, _ := claim.ClaimHash()
		a.keeper.logger(ctx).Error("ERC721 transfer failed",
			"cause", err.Error(),
			"claim type", claim.GetType(),
			"id", types.GetAttestationKey(claim.GetEventNonce(), hash),
			"nonce", fmt.Sprint(claim.GetEventNonce()),
		)
	} else {
		a.keeper.logger(ctx).Info("SendERC721ToCosmos to local gravity receiver", "ethSender", claim.EthereumSender,
			"receiver", receiver, "contract", claim.TokenContract, "classId", nftToken.ClassId, "nftId", nftToken.Id,
			"nonce", claim.EventNonce, "ethContract", claim.TokenContract, "ethBlockHeight", claim.EthBlockHeight,
			"cosmosBlockHeight", ctx.BlockHeight(),
		)
		if err := ctx.EventManager().EmitTypedEvent(&types.EventSendNFTToCosmosLocal{
			Nonce:    fmt.Sprint(claim.EventNonce),
			Receiver: receiver.String(),
			ClassId:  nftToken.ClassId,
			TokenId:  nftToken.Id,
		}); err != nil {
			return err
		}
	}

	return err
}

func (a AttestationHandler) addERC721ToIbcAutoForwardQueue(
	ctx sdk.Context,
	receiver sdk.AccAddress,
	accountPrefix string,
	nftToken nft.NFT,
	channel string,
	claim types.MsgSendNFTToCosmosClaim,
) error {
	if strings.TrimSpace(accountPrefix) == "" {
		panic("invalid call to addToIbcAutoForwardQueue: provided accountPrefix is empty!")
	}
	acctPrefix, err := gravitytypes.GetPrefixFromBech32(claim.CosmosReceiver)
	if err != nil || acctPrefix != accountPrefix {
		panic(fmt.Sprintf("invalid call to addToIbcAutoForwardQueue: invalid or inaccurate accountPrefix %s for receiver %s!", accountPrefix, claim.CosmosReceiver))
	}

	forward := types.PendingNFTIbcAutoForward{
		ForeignReceiver: claim.CosmosReceiver,
		ClassId:         nftToken.ClassId,
		TokenId:         nftToken.Id,
		IbcChannel:      channel,
		EventNonce:      claim.EventNonce,
	}

	return a.keeper.AddPendingNFTPendingIbcAutoForward(ctx, forward)
}