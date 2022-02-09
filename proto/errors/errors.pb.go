// Code generated by protoc-gen-go. DO NOT EDIT.
// source: errors/errors.proto

package errors

import (
	fmt "fmt"
	proto "google.golang.org/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type Error struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Code                 int32    `protobuf:"varint,2,opt,name=code,proto3" json:"code,omitempty"`
	Detail               string   `protobuf:"bytes,3,opt,name=detail,proto3" json:"detail,omitempty"`
	Status               string   `protobuf:"bytes,4,opt,name=status,proto3" json:"status,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Error) Reset()         { *m = Error{} }
func (m *Error) String() string { return proto.CompactTextString(m) }
func (*Error) ProtoMessage()    {}
func (*Error) Descriptor() ([]byte, []int) {
	return fileDescriptor_85c4eef3398a32b2, []int{0}
}

func (m *Error) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Error.Unmarshal(m, b)
}
func (m *Error) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Error.Marshal(b, m, deterministic)
}
func (m *Error) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Error.Merge(m, src)
}
func (m *Error) XXX_Size() int {
	return xxx_messageInfo_Error.Size(m)
}
func (m *Error) XXX_DiscardUnknown() {
	xxx_messageInfo_Error.DiscardUnknown(m)
}

var xxx_messageInfo_Error proto.InternalMessageInfo

func (m *Error) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Error) GetCode() int32 {
	if m != nil {
		return m.Code
	}
	return 0
}

func (m *Error) GetDetail() string {
	if m != nil {
		return m.Detail
	}
	return ""
}

func (m *Error) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

func init() {
	proto.RegisterType((*Error)(nil), "errors.Error")
}

func init() { proto.RegisterFile("errors/errors.proto", fileDescriptor_85c4eef3398a32b2) }

var fileDescriptor_85c4eef3398a32b2 = []byte{
	// 150 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x4e, 0x2d, 0x2a, 0xca,
	0x2f, 0x2a, 0xd6, 0x87, 0x50, 0x7a, 0x05, 0x45, 0xf9, 0x25, 0xf9, 0x42, 0x6c, 0x10, 0x9e, 0x52,
	0x34, 0x17, 0xab, 0x2b, 0x88, 0x25, 0xc4, 0xc7, 0xc5, 0x94, 0x99, 0x22, 0xc1, 0xa8, 0xc0, 0xa8,
	0xc1, 0x19, 0xc4, 0x94, 0x99, 0x22, 0x24, 0xc4, 0xc5, 0x92, 0x9c, 0x9f, 0x92, 0x2a, 0xc1, 0xa4,
	0xc0, 0xa8, 0xc1, 0x1a, 0x04, 0x66, 0x0b, 0x89, 0x71, 0xb1, 0xa5, 0xa4, 0x96, 0x24, 0x66, 0xe6,
	0x48, 0x30, 0x83, 0xd5, 0x41, 0x79, 0x20, 0xf1, 0xe2, 0x92, 0xc4, 0x92, 0xd2, 0x62, 0x09, 0x16,
	0x88, 0x38, 0x84, 0xe7, 0xa4, 0x1f, 0xa5, 0x9b, 0x9e, 0x59, 0x92, 0x51, 0x9a, 0xa4, 0x97, 0x9c,
	0x9f, 0xab, 0x9f, 0x9b, 0x99, 0x5c, 0x94, 0x0f, 0x25, 0xcb, 0x8c, 0xf5, 0xc1, 0xee, 0x80, 0x3a,
	0xca, 0x1a, 0x42, 0x25, 0xb1, 0x81, 0x05, 0x8d, 0x01, 0x01, 0x00, 0x00, 0xff, 0xff, 0xc4, 0x06,
	0x4a, 0xd5, 0xb3, 0x00, 0x00, 0x00,
}
