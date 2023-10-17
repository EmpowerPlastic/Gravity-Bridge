package types

import (
	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	// AttestationVotesPowerThreshold threshold of votes power to succeed
	AttestationVotesPowerThreshold = sdk.NewInt(66)
)

// ValidateBasic validates genesis state by looping through the params and
// calling their validation functions
func (s GenesisState) ValidateBasic() error {
	if err := s.Params.ValidateBasic(); err != nil {
		return sdkerrors.Wrap(err, "params")
	}
	return nil
}

// DefaultGenesisState returns empty genesis state
// nolint: exhaustruct
func DefaultGenesisState() *GenesisState {
	// TODO: ALL OF THIS
	return &GenesisState{
		Params:        DefaultParams(),
		GravityNonces: GravityNonces{},
	}
}

// DefaultParams returns a copy of the default params
func DefaultParams() *Params {
	// TODO: ALL OF THIS
	return &Params{}
}

// ValidateBasic checks that the parameters have valid values.
func (p Params) ValidateBasic() error {
	// TODO: ALL OF THIS
	return nil
}
