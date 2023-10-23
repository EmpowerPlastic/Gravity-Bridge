// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: gravitynft/v1/types.proto

package types

import (
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// LastObservedNFTEthereumBlockHeight stores the last observed
// Ethereum block height along with the Cosmos block height that
// it was observed at. These two numbers can be used to project
// outward and always produce batches with timeouts in the future
// even if no Ethereum block height has been relayed for a long time
type LastObservedNFTEthereumBlockHeight struct {
	CosmosBlockHeight   uint64 `protobuf:"varint,1,opt,name=cosmos_block_height,json=cosmosBlockHeight,proto3" json:"cosmos_block_height,omitempty"`
	EthereumBlockHeight uint64 `protobuf:"varint,2,opt,name=ethereum_block_height,json=ethereumBlockHeight,proto3" json:"ethereum_block_height,omitempty"`
}

func (m *LastObservedNFTEthereumBlockHeight) Reset()         { *m = LastObservedNFTEthereumBlockHeight{} }
func (m *LastObservedNFTEthereumBlockHeight) String() string { return proto.CompactTextString(m) }
func (*LastObservedNFTEthereumBlockHeight) ProtoMessage()    {}
func (*LastObservedNFTEthereumBlockHeight) Descriptor() ([]byte, []int) {
	return fileDescriptor_94c6f9c9e4250776, []int{0}
}
func (m *LastObservedNFTEthereumBlockHeight) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *LastObservedNFTEthereumBlockHeight) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_LastObservedNFTEthereumBlockHeight.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *LastObservedNFTEthereumBlockHeight) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LastObservedNFTEthereumBlockHeight.Merge(m, src)
}
func (m *LastObservedNFTEthereumBlockHeight) XXX_Size() int {
	return m.Size()
}
func (m *LastObservedNFTEthereumBlockHeight) XXX_DiscardUnknown() {
	xxx_messageInfo_LastObservedNFTEthereumBlockHeight.DiscardUnknown(m)
}

var xxx_messageInfo_LastObservedNFTEthereumBlockHeight proto.InternalMessageInfo

func (m *LastObservedNFTEthereumBlockHeight) GetCosmosBlockHeight() uint64 {
	if m != nil {
		return m.CosmosBlockHeight
	}
	return 0
}

func (m *LastObservedNFTEthereumBlockHeight) GetEthereumBlockHeight() uint64 {
	if m != nil {
		return m.EthereumBlockHeight
	}
	return 0
}

// This records the relationship between an ERC721 token and the class id
// of the corresponding Cosmos originated asset
type ERC721ToClassId struct {
	Erc721  string `protobuf:"bytes,1,opt,name=erc721,proto3" json:"erc721,omitempty"`
	ClassId string `protobuf:"bytes,2,opt,name=class_id,json=classId,proto3" json:"class_id,omitempty"`
}

func (m *ERC721ToClassId) Reset()         { *m = ERC721ToClassId{} }
func (m *ERC721ToClassId) String() string { return proto.CompactTextString(m) }
func (*ERC721ToClassId) ProtoMessage()    {}
func (*ERC721ToClassId) Descriptor() ([]byte, []int) {
	return fileDescriptor_94c6f9c9e4250776, []int{1}
}
func (m *ERC721ToClassId) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ERC721ToClassId) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ERC721ToClassId.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ERC721ToClassId) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ERC721ToClassId.Merge(m, src)
}
func (m *ERC721ToClassId) XXX_Size() int {
	return m.Size()
}
func (m *ERC721ToClassId) XXX_DiscardUnknown() {
	xxx_messageInfo_ERC721ToClassId.DiscardUnknown(m)
}

var xxx_messageInfo_ERC721ToClassId proto.InternalMessageInfo

func (m *ERC721ToClassId) GetErc721() string {
	if m != nil {
		return m.Erc721
	}
	return ""
}

func (m *ERC721ToClassId) GetClassId() string {
	if m != nil {
		return m.ClassId
	}
	return ""
}

type PendingNFTIbcAutoForward struct {
	ForeignReceiver string `protobuf:"bytes,1,opt,name=foreign_receiver,json=foreignReceiver,proto3" json:"foreign_receiver,omitempty"`
	ClassId         string `protobuf:"bytes,2,opt,name=class_id,json=classId,proto3" json:"class_id,omitempty"`
	TokenId         string `protobuf:"bytes,3,opt,name=token_id,json=tokenId,proto3" json:"token_id,omitempty"`
	IbcChannel      string `protobuf:"bytes,4,opt,name=ibc_channel,json=ibcChannel,proto3" json:"ibc_channel,omitempty"`
	EventNonce      uint64 `protobuf:"varint,5,opt,name=event_nonce,json=eventNonce,proto3" json:"event_nonce,omitempty"`
}

func (m *PendingNFTIbcAutoForward) Reset()         { *m = PendingNFTIbcAutoForward{} }
func (m *PendingNFTIbcAutoForward) String() string { return proto.CompactTextString(m) }
func (*PendingNFTIbcAutoForward) ProtoMessage()    {}
func (*PendingNFTIbcAutoForward) Descriptor() ([]byte, []int) {
	return fileDescriptor_94c6f9c9e4250776, []int{2}
}
func (m *PendingNFTIbcAutoForward) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *PendingNFTIbcAutoForward) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_PendingNFTIbcAutoForward.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *PendingNFTIbcAutoForward) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PendingNFTIbcAutoForward.Merge(m, src)
}
func (m *PendingNFTIbcAutoForward) XXX_Size() int {
	return m.Size()
}
func (m *PendingNFTIbcAutoForward) XXX_DiscardUnknown() {
	xxx_messageInfo_PendingNFTIbcAutoForward.DiscardUnknown(m)
}

var xxx_messageInfo_PendingNFTIbcAutoForward proto.InternalMessageInfo

func (m *PendingNFTIbcAutoForward) GetForeignReceiver() string {
	if m != nil {
		return m.ForeignReceiver
	}
	return ""
}

func (m *PendingNFTIbcAutoForward) GetClassId() string {
	if m != nil {
		return m.ClassId
	}
	return ""
}

func (m *PendingNFTIbcAutoForward) GetTokenId() string {
	if m != nil {
		return m.TokenId
	}
	return ""
}

func (m *PendingNFTIbcAutoForward) GetIbcChannel() string {
	if m != nil {
		return m.IbcChannel
	}
	return ""
}

func (m *PendingNFTIbcAutoForward) GetEventNonce() uint64 {
	if m != nil {
		return m.EventNonce
	}
	return 0
}

func init() {
	proto.RegisterType((*LastObservedNFTEthereumBlockHeight)(nil), "gravitynft.v1.LastObservedNFTEthereumBlockHeight")
	proto.RegisterType((*ERC721ToClassId)(nil), "gravitynft.v1.ERC721ToClassId")
	proto.RegisterType((*PendingNFTIbcAutoForward)(nil), "gravitynft.v1.PendingNFTIbcAutoForward")
}

func init() { proto.RegisterFile("gravitynft/v1/types.proto", fileDescriptor_94c6f9c9e4250776) }

var fileDescriptor_94c6f9c9e4250776 = []byte{
	// 388 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x92, 0x41, 0x8f, 0x93, 0x40,
	0x14, 0xc7, 0x8b, 0xd6, 0x6a, 0xc7, 0x98, 0x2a, 0x55, 0x43, 0x3d, 0xa0, 0xe1, 0xa4, 0x07, 0x21,
	0xd4, 0x43, 0xcf, 0x52, 0x5b, 0x6d, 0x62, 0xaa, 0x21, 0x3d, 0x19, 0x13, 0x02, 0x33, 0xaf, 0xc3,
	0xa4, 0x30, 0xd3, 0x0c, 0x03, 0xda, 0x6f, 0xe0, 0xd1, 0x2f, 0xb4, 0xf7, 0x3d, 0xf6, 0xb8, 0xc7,
	0x4d, 0xfb, 0x45, 0x36, 0x0c, 0x24, 0xdb, 0xdd, 0x4d, 0xf6, 0xc6, 0xfb, 0xfd, 0xfe, 0xef, 0xbd,
	0x84, 0x79, 0x68, 0x44, 0x65, 0x5c, 0x31, 0xb5, 0xe3, 0x6b, 0xe5, 0x55, 0xbe, 0xa7, 0x76, 0x5b,
	0x28, 0xdc, 0xad, 0x14, 0x4a, 0x98, 0xcf, 0xae, 0x95, 0x5b, 0xf9, 0x6f, 0x5e, 0x52, 0x41, 0x85,
	0x36, 0x5e, 0xfd, 0xd5, 0x84, 0x9c, 0x7f, 0x06, 0x72, 0xbe, 0xc7, 0x85, 0xfa, 0x91, 0x14, 0x20,
	0x2b, 0x20, 0xcb, 0xf9, 0x6a, 0xa6, 0x52, 0x90, 0x50, 0xe6, 0x41, 0x26, 0xf0, 0xe6, 0x1b, 0x30,
	0x9a, 0x2a, 0xd3, 0x45, 0x43, 0x2c, 0x8a, 0x5c, 0x14, 0x51, 0x52, 0xd3, 0x28, 0xd5, 0xd8, 0x32,
	0xde, 0x19, 0xef, 0xbb, 0xe1, 0x8b, 0x46, 0x9d, 0xe6, 0xc7, 0xe8, 0x15, 0xb4, 0x63, 0x6e, 0x76,
	0x3c, 0xd0, 0x1d, 0x43, 0xb8, 0xbb, 0xc3, 0xf9, 0x82, 0x06, 0xb3, 0x70, 0x3a, 0x19, 0xfb, 0x2b,
	0x31, 0xcd, 0xe2, 0xa2, 0x58, 0x10, 0xf3, 0x35, 0xea, 0x81, 0xc4, 0x93, 0xb1, 0xaf, 0x37, 0xf5,
	0xc3, 0xb6, 0x32, 0x47, 0xe8, 0x09, 0xae, 0x23, 0x11, 0x23, 0x7a, 0x62, 0x3f, 0x7c, 0x8c, 0x9b,
	0x16, 0xe7, 0xcc, 0x40, 0xd6, 0x4f, 0xe0, 0x84, 0x71, 0xba, 0x9c, 0xaf, 0x16, 0x09, 0xfe, 0x5c,
	0x2a, 0x31, 0x17, 0xf2, 0x4f, 0x2c, 0x89, 0xf9, 0x01, 0x3d, 0x5f, 0x0b, 0x09, 0x8c, 0xf2, 0x48,
	0x02, 0x06, 0x56, 0x81, 0x6c, 0x27, 0x0f, 0x5a, 0x1e, 0xb6, 0xf8, 0x9e, 0x15, 0xb5, 0x52, 0x62,
	0x03, 0xbc, 0x56, 0x0f, 0x1b, 0xa5, 0xeb, 0x05, 0x31, 0xdf, 0xa2, 0xa7, 0x2c, 0xc1, 0x11, 0x4e,
	0x63, 0xce, 0x21, 0xb3, 0xba, 0xda, 0x22, 0x96, 0xe0, 0x69, 0x43, 0xea, 0x00, 0x54, 0xc0, 0x55,
	0xc4, 0x05, 0xc7, 0x60, 0x3d, 0xd2, 0xbf, 0x03, 0x69, 0xb4, 0xac, 0x49, 0xf0, 0xfb, 0xfc, 0x60,
	0x1b, 0xfb, 0x83, 0x6d, 0x5c, 0x1e, 0x6c, 0xe3, 0xff, 0xd1, 0xee, 0xec, 0x8f, 0x76, 0xe7, 0xe2,
	0x68, 0x77, 0x7e, 0x05, 0x94, 0xa9, 0xb4, 0x4c, 0x5c, 0x2c, 0x72, 0xef, 0x6b, 0xf3, 0xb4, 0x1f,
	0x03, 0xc9, 0x08, 0x85, 0xdb, 0x65, 0x2e, 0x48, 0x99, 0x81, 0xf7, 0xd7, 0x3b, 0x39, 0x0e, 0x7d,
	0x19, 0x49, 0x4f, 0xbf, 0xfa, 0xa7, 0xab, 0x00, 0x00, 0x00, 0xff, 0xff, 0x06, 0x28, 0x01, 0x02,
	0x37, 0x02, 0x00, 0x00,
}

func (m *LastObservedNFTEthereumBlockHeight) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *LastObservedNFTEthereumBlockHeight) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *LastObservedNFTEthereumBlockHeight) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.EthereumBlockHeight != 0 {
		i = encodeVarintTypes(dAtA, i, uint64(m.EthereumBlockHeight))
		i--
		dAtA[i] = 0x10
	}
	if m.CosmosBlockHeight != 0 {
		i = encodeVarintTypes(dAtA, i, uint64(m.CosmosBlockHeight))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *ERC721ToClassId) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ERC721ToClassId) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ERC721ToClassId) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.ClassId) > 0 {
		i -= len(m.ClassId)
		copy(dAtA[i:], m.ClassId)
		i = encodeVarintTypes(dAtA, i, uint64(len(m.ClassId)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Erc721) > 0 {
		i -= len(m.Erc721)
		copy(dAtA[i:], m.Erc721)
		i = encodeVarintTypes(dAtA, i, uint64(len(m.Erc721)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *PendingNFTIbcAutoForward) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *PendingNFTIbcAutoForward) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *PendingNFTIbcAutoForward) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.EventNonce != 0 {
		i = encodeVarintTypes(dAtA, i, uint64(m.EventNonce))
		i--
		dAtA[i] = 0x28
	}
	if len(m.IbcChannel) > 0 {
		i -= len(m.IbcChannel)
		copy(dAtA[i:], m.IbcChannel)
		i = encodeVarintTypes(dAtA, i, uint64(len(m.IbcChannel)))
		i--
		dAtA[i] = 0x22
	}
	if len(m.TokenId) > 0 {
		i -= len(m.TokenId)
		copy(dAtA[i:], m.TokenId)
		i = encodeVarintTypes(dAtA, i, uint64(len(m.TokenId)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.ClassId) > 0 {
		i -= len(m.ClassId)
		copy(dAtA[i:], m.ClassId)
		i = encodeVarintTypes(dAtA, i, uint64(len(m.ClassId)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.ForeignReceiver) > 0 {
		i -= len(m.ForeignReceiver)
		copy(dAtA[i:], m.ForeignReceiver)
		i = encodeVarintTypes(dAtA, i, uint64(len(m.ForeignReceiver)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintTypes(dAtA []byte, offset int, v uint64) int {
	offset -= sovTypes(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *LastObservedNFTEthereumBlockHeight) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.CosmosBlockHeight != 0 {
		n += 1 + sovTypes(uint64(m.CosmosBlockHeight))
	}
	if m.EthereumBlockHeight != 0 {
		n += 1 + sovTypes(uint64(m.EthereumBlockHeight))
	}
	return n
}

func (m *ERC721ToClassId) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Erc721)
	if l > 0 {
		n += 1 + l + sovTypes(uint64(l))
	}
	l = len(m.ClassId)
	if l > 0 {
		n += 1 + l + sovTypes(uint64(l))
	}
	return n
}

func (m *PendingNFTIbcAutoForward) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.ForeignReceiver)
	if l > 0 {
		n += 1 + l + sovTypes(uint64(l))
	}
	l = len(m.ClassId)
	if l > 0 {
		n += 1 + l + sovTypes(uint64(l))
	}
	l = len(m.TokenId)
	if l > 0 {
		n += 1 + l + sovTypes(uint64(l))
	}
	l = len(m.IbcChannel)
	if l > 0 {
		n += 1 + l + sovTypes(uint64(l))
	}
	if m.EventNonce != 0 {
		n += 1 + sovTypes(uint64(m.EventNonce))
	}
	return n
}

func sovTypes(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozTypes(x uint64) (n int) {
	return sovTypes(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *LastObservedNFTEthereumBlockHeight) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTypes
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: LastObservedNFTEthereumBlockHeight: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: LastObservedNFTEthereumBlockHeight: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field CosmosBlockHeight", wireType)
			}
			m.CosmosBlockHeight = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.CosmosBlockHeight |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field EthereumBlockHeight", wireType)
			}
			m.EthereumBlockHeight = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.EthereumBlockHeight |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipTypes(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTypes
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *ERC721ToClassId) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTypes
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: ERC721ToClassId: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ERC721ToClassId: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Erc721", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Erc721 = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ClassId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ClassId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTypes(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTypes
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *PendingNFTIbcAutoForward) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTypes
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: PendingNFTIbcAutoForward: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: PendingNFTIbcAutoForward: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ForeignReceiver", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ForeignReceiver = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ClassId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ClassId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TokenId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.TokenId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field IbcChannel", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.IbcChannel = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field EventNonce", wireType)
			}
			m.EventNonce = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.EventNonce |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipTypes(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTypes
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipTypes(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowTypes
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthTypes
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupTypes
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthTypes
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthTypes        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowTypes          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupTypes = fmt.Errorf("proto: unexpected end of group")
)
