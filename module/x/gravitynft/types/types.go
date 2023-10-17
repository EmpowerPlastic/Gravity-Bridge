package types

import (
	"crypto/md5"
	"encoding/binary"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (p Params) Validate() error {
	// TODO: When we have some real params to work with add validation
	return nil
}

// UInt64FromBytesUnsafe create uint from binary big endian representation
// Note: This is unsafe because the function will panic if provided over 8 bytes
func UInt64FromBytesUnsafe(s []byte) uint64 {
	if len(s) > 8 {
		panic("Invalid uint64 bytes passed to UInt64FromBytes!")
	}
	return binary.BigEndian.Uint64(s)
}

// UInt64Bytes uses the SDK byte marshaling to encode a uint64
func UInt64Bytes(n uint64) []byte {
	return sdk.Uint64ToBigEndian(n)
}

// Hashing string using cryptographic MD5 function
// returns 128bit(16byte) value
func HashString(input string) []byte {
	md5 := md5.New()
	md5.Write([]byte(input))
	return md5.Sum(nil)
}

func AppendBytes(args ...[]byte) []byte {
	length := 0
	for _, v := range args {
		length += len(v)
	}

	res := make([]byte, length)

	length = 0
	for _, v := range args {
		copy(res[length:length+len(v)], v)
		length += len(v)
	}

	return res
}
