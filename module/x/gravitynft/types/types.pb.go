// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: gravitynft/v1/types.proto

package types

import (
	fmt "fmt"
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
	return fileDescriptor_94c6f9c9e4250776, []int{1}
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
	proto.RegisterType((*PendingNFTIbcAutoForward)(nil), "gravitynft.v1.PendingNFTIbcAutoForward")
}

func init() { proto.RegisterFile("gravitynft/v1/types.proto", fileDescriptor_94c6f9c9e4250776) }

var fileDescriptor_94c6f9c9e4250776 = []byte{
	// 343 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x64, 0x91, 0xc1, 0x4a, 0xeb, 0x40,
	0x14, 0x86, 0x9b, 0x7b, 0x7b, 0xaf, 0x76, 0x44, 0xd4, 0x14, 0x21, 0xdd, 0x44, 0xe9, 0x4a, 0x17,
	0x26, 0x54, 0x9f, 0xc0, 0x88, 0xd5, 0x82, 0x54, 0x09, 0x5d, 0x89, 0x10, 0x32, 0x33, 0xa7, 0xc9,
	0xd0, 0x64, 0xa6, 0xcc, 0x4c, 0xa2, 0x7d, 0x03, 0x97, 0xbe, 0x90, 0x7b, 0x97, 0x5d, 0xba, 0x94,
	0xf6, 0x45, 0x24, 0x93, 0x82, 0x55, 0x97, 0xe7, 0xfb, 0xce, 0xc9, 0x4f, 0xe6, 0x47, 0x9d, 0x44,
	0xc6, 0x25, 0xd3, 0x33, 0x3e, 0xd6, 0x7e, 0xd9, 0xf3, 0xf5, 0x6c, 0x0a, 0xca, 0x9b, 0x4a, 0xa1,
	0x85, 0xbd, 0xfd, 0xa5, 0xbc, 0xb2, 0xd7, 0x7d, 0xb6, 0x50, 0xf7, 0x26, 0x56, 0xfa, 0x16, 0x2b,
	0x90, 0x25, 0xd0, 0x61, 0x7f, 0x74, 0xa9, 0x53, 0x90, 0x50, 0xe4, 0x41, 0x26, 0xc8, 0xe4, 0x1a,
	0x58, 0x92, 0x6a, 0xdb, 0x43, 0x6d, 0x22, 0x54, 0x2e, 0x54, 0x84, 0x2b, 0x1a, 0xa5, 0x06, 0x3b,
	0xd6, 0xa1, 0x75, 0xd4, 0x0c, 0xf7, 0x6a, 0xb5, 0xbe, 0x7f, 0x8a, 0xf6, 0x61, 0xf5, 0x99, 0xef,
	0x17, 0x7f, 0xcc, 0x45, 0x1b, 0x7e, 0x67, 0x74, 0x5f, 0x2d, 0xe4, 0xdc, 0x01, 0xa7, 0x8c, 0x27,
	0xc3, 0xfe, 0x68, 0x80, 0xc9, 0x79, 0xa1, 0x45, 0x5f, 0xc8, 0xc7, 0x58, 0x52, 0xfb, 0x18, 0xed,
	0x8e, 0x85, 0x04, 0x96, 0xf0, 0x48, 0x02, 0x01, 0x56, 0x82, 0x34, 0xe9, 0xad, 0x70, 0x67, 0xc5,
	0xc3, 0x15, 0xb6, 0x3b, 0x68, 0x93, 0x64, 0xb1, 0x52, 0x11, 0xa3, 0x26, 0xae, 0x15, 0x6e, 0x98,
	0x79, 0x40, 0x2b, 0xa5, 0xc5, 0x04, 0x78, 0xa5, 0xfe, 0xd6, 0xca, 0xcc, 0x03, 0x6a, 0x1f, 0xa0,
	0x2d, 0x86, 0x49, 0x44, 0xd2, 0x98, 0x73, 0xc8, 0x9c, 0xa6, 0xb1, 0x88, 0x61, 0x72, 0x51, 0x93,
	0x6a, 0x01, 0x4a, 0xe0, 0x3a, 0xe2, 0x82, 0x13, 0x70, 0xfe, 0x99, 0x1f, 0x41, 0x06, 0x0d, 0x2b,
	0x12, 0x3c, 0xbc, 0x2d, 0x5c, 0x6b, 0xbe, 0x70, 0xad, 0x8f, 0x85, 0x6b, 0xbd, 0x2c, 0xdd, 0xc6,
	0x7c, 0xe9, 0x36, 0xde, 0x97, 0x6e, 0xe3, 0x3e, 0x48, 0x98, 0x4e, 0x0b, 0xec, 0x11, 0x91, 0xfb,
	0x57, 0xf5, 0xf3, 0x9f, 0x04, 0x92, 0xd1, 0x04, 0x7e, 0x8e, 0xb9, 0xa0, 0x45, 0x06, 0xfe, 0x93,
	0xbf, 0x56, 0xa0, 0x69, 0x0f, 0xff, 0x37, 0xf5, 0x9d, 0x7d, 0x06, 0x00, 0x00, 0xff, 0xff, 0x8d,
	0xea, 0x94, 0xb9, 0xdb, 0x01, 0x00, 0x00,
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