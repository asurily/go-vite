// Code generated by protoc-gen-go. DO NOT EDIT.
// source: account_meta.proto

package vitepb

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type AccountMeta struct {
	AccountId            []byte                `protobuf:"bytes,1,opt,name=accountId,proto3" json:"accountId,omitempty"`
	TokenList            []*AccountSimpleToken `protobuf:"bytes,2,rep,name=tokenList,proto3" json:"tokenList,omitempty"`
	XXX_NoUnkeyedLiteral struct{}              `json:"-"`
	XXX_unrecognized     []byte                `json:"-"`
	XXX_sizecache        int32                 `json:"-"`
}

func (m *AccountMeta) Reset()         { *m = AccountMeta{} }
func (m *AccountMeta) String() string { return proto.CompactTextString(m) }
func (*AccountMeta) ProtoMessage()    {}
func (*AccountMeta) Descriptor() ([]byte, []int) {
	return fileDescriptor_account_meta_2cc6bcc8d25eb49e, []int{0}
}
func (m *AccountMeta) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AccountMeta.Unmarshal(m, b)
}
func (m *AccountMeta) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AccountMeta.Marshal(b, m, deterministic)
}
func (dst *AccountMeta) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AccountMeta.Merge(dst, src)
}
func (m *AccountMeta) XXX_Size() int {
	return xxx_messageInfo_AccountMeta.Size(m)
}
func (m *AccountMeta) XXX_DiscardUnknown() {
	xxx_messageInfo_AccountMeta.DiscardUnknown(m)
}

var xxx_messageInfo_AccountMeta proto.InternalMessageInfo

func (m *AccountMeta) GetAccountId() []byte {
	if m != nil {
		return m.AccountId
	}
	return nil
}

func (m *AccountMeta) GetTokenList() []*AccountSimpleToken {
	if m != nil {
		return m.TokenList
	}
	return nil
}

type AccountSimpleToken struct {
	TokenId                []byte   `protobuf:"bytes,1,opt,name=tokenId,proto3" json:"tokenId,omitempty"`
	LastAccountBlockHeight []byte   `protobuf:"bytes,2,opt,name=lastAccountBlockHeight,proto3" json:"lastAccountBlockHeight,omitempty"`
	XXX_NoUnkeyedLiteral   struct{} `json:"-"`
	XXX_unrecognized       []byte   `json:"-"`
	XXX_sizecache          int32    `json:"-"`
}

func (m *AccountSimpleToken) Reset()         { *m = AccountSimpleToken{} }
func (m *AccountSimpleToken) String() string { return proto.CompactTextString(m) }
func (*AccountSimpleToken) ProtoMessage()    {}
func (*AccountSimpleToken) Descriptor() ([]byte, []int) {
	return fileDescriptor_account_meta_2cc6bcc8d25eb49e, []int{1}
}
func (m *AccountSimpleToken) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AccountSimpleToken.Unmarshal(m, b)
}
func (m *AccountSimpleToken) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AccountSimpleToken.Marshal(b, m, deterministic)
}
func (dst *AccountSimpleToken) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AccountSimpleToken.Merge(dst, src)
}
func (m *AccountSimpleToken) XXX_Size() int {
	return xxx_messageInfo_AccountSimpleToken.Size(m)
}
func (m *AccountSimpleToken) XXX_DiscardUnknown() {
	xxx_messageInfo_AccountSimpleToken.DiscardUnknown(m)
}

var xxx_messageInfo_AccountSimpleToken proto.InternalMessageInfo

func (m *AccountSimpleToken) GetTokenId() []byte {
	if m != nil {
		return m.TokenId
	}
	return nil
}

func (m *AccountSimpleToken) GetLastAccountBlockHeight() []byte {
	if m != nil {
		return m.LastAccountBlockHeight
	}
	return nil
}

func init() {
	proto.RegisterType((*AccountMeta)(nil), "vitepb.AccountMeta")
	proto.RegisterType((*AccountSimpleToken)(nil), "vitepb.AccountSimpleToken")
}

func init() { proto.RegisterFile("account_meta.proto", fileDescriptor_account_meta_2cc6bcc8d25eb49e) }

var fileDescriptor_account_meta_2cc6bcc8d25eb49e = []byte{
	// 167 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x4a, 0x4c, 0x4e, 0xce,
	0x2f, 0xcd, 0x2b, 0x89, 0xcf, 0x4d, 0x2d, 0x49, 0xd4, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62,
	0x2b, 0xcb, 0x2c, 0x49, 0x2d, 0x48, 0x52, 0x4a, 0xe5, 0xe2, 0x76, 0x84, 0xc8, 0xfa, 0xa6, 0x96,
	0x24, 0x0a, 0xc9, 0x70, 0x71, 0x42, 0x15, 0x7b, 0xa6, 0x48, 0x30, 0x2a, 0x30, 0x6a, 0xf0, 0x04,
	0x21, 0x04, 0x84, 0x2c, 0xb8, 0x38, 0x4b, 0xf2, 0xb3, 0x53, 0xf3, 0x7c, 0x32, 0x8b, 0x4b, 0x24,
	0x98, 0x14, 0x98, 0x35, 0xb8, 0x8d, 0xa4, 0xf4, 0x20, 0x06, 0xe9, 0x41, 0x4d, 0x09, 0xce, 0xcc,
	0x2d, 0xc8, 0x49, 0x0d, 0x01, 0xa9, 0x0a, 0x42, 0x28, 0x56, 0x4a, 0xe3, 0x12, 0xc2, 0x54, 0x20,
	0x24, 0xc1, 0xc5, 0x0e, 0x56, 0x02, 0xb7, 0x0b, 0xc6, 0x15, 0x32, 0xe3, 0x12, 0xcb, 0x49, 0x2c,
	0x2e, 0x81, 0xea, 0x71, 0xca, 0xc9, 0x4f, 0xce, 0xf6, 0x48, 0xcd, 0x4c, 0xcf, 0x00, 0x59, 0x0b,
	0x52, 0x88, 0x43, 0x36, 0x89, 0x0d, 0xec, 0x3b, 0x63, 0x40, 0x00, 0x00, 0x00, 0xff, 0xff, 0x6f,
	0x77, 0x3f, 0xd5, 0xf3, 0x00, 0x00, 0x00,
}
