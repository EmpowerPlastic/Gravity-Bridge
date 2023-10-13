package types

import (
	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/cosmos/cosmos-sdk/types/errors"
)

// ValidateBasic checks the ForeignReceiver is valid and foreign, the Amount is non-zero, the IbcChannel is
// non-empty, and the EventNonce is non-zero
func (p PendingNFTIbcAutoForward) ValidateBasic() error {
	prefix, _, err := bech32.DecodeAndConvert(p.ForeignReceiver)
	if err != nil {
		return sdkerrors.Wrapf(err, "ForeignReceiver is not a valid bech32 address")
	}
	nativePrefix := sdk.GetConfig().GetBech32AccountAddrPrefix()
	if prefix == nativePrefix {
		return sdkerrors.Wrapf(errors.ErrInvalidAddress, "ForeignReceiver cannot have the native chain prefix")
	}

	if p.ClassId == "" {
		return sdkerrors.Wrap(ErrInvalid, "ClassId must not be an empty string")
	}

	if p.TokenId == "" {
		return sdkerrors.Wrap(ErrInvalid, "TokenId must not be an empty string")
	}

	if p.IbcChannel == "" {
		return sdkerrors.Wrap(ErrInvalid, "IbcChannel must not be an empty string")
	}

	if p.EventNonce == 0 {
		return sdkerrors.Wrap(ErrInvalid, "EventNonce must be non-zero")
	}

	return nil
}