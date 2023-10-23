package types

import (
	"fmt"
	gravitytypes "github.com/Gravity-Bridge/Gravity-Bridge/module/x/gravity/types"
	"strings"
)

const (
	// GravityDenomPrefix indicates the prefix for all assets minted by this module
	GravityNFTClassIDPrefix = ModuleName

	// GravityDenomSeparator is the separator for gravity denoms
	GravityNFTClassIDSeparator = ""

	// ETHContractAddressLen is the length of contract address strings
	ETHContractAddressLen = 42

	// GravityNFTClassIDLen is the length of the denoms generated by the gravity module
	GravityNFTClassIDLen = len(GravityNFTClassIDPrefix) + len(GravityNFTClassIDSeparator) + ETHContractAddressLen
)

/////////////////////////
// ERC721              //
/////////////////////////

// GravityERC721ClassId converts an EthAddress to a gravity cosmos class id for ERC721 tokens
func GravityERC721ClassId(tokenContract gravitytypes.EthAddress) string {
	return fmt.Sprintf("%s%s%s", GravityNFTClassIDPrefix, GravityNFTClassIDSeparator, strings.ToLower(tokenContract.GetAddress().Hex()))
}

// TODO: Need to be checked and tested!!!
// GravityDenomToERC721 converts a gravity cosmos denom to an EthAddress
func GravityDenomToERC721(classID string) (*gravitytypes.EthAddress, error) {
	fullPrefix := GravityNFTClassIDPrefix + GravityNFTClassIDSeparator
	if !strings.HasPrefix(classID, fullPrefix) {
		return nil, fmt.Errorf("classID prefix(%s) not equal to expected(%s)", classID, fullPrefix)
	}
	contract := strings.TrimPrefix(classID, fullPrefix)
	ethAddr, err := gravitytypes.NewEthAddress(contract)
	switch {
	case err != nil:
		return nil, fmt.Errorf("error(%s) validating ethereum contract address", err)
	case len(classID) != GravityNFTClassIDLen:
		return nil, fmt.Errorf("len(classID)(%d) not equal to GravityDenomLen(%d)", len(classID), GravityNFTClassIDLen)
	default:
		return ethAddr, nil
	}
}