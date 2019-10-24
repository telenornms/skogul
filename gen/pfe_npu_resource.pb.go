// Code generated by protoc-gen-go. DO NOT EDIT.
// source: pfe_npu_resource.proto

package gen

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
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

type JunosPfeNpu struct {
	NpuMemory            []*JunosPfeNpuNpuMemoryList `protobuf:"bytes,151,rep,name=npu_memory,json=npuMemory" json:"npu_memory,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                    `json:"-"`
	XXX_unrecognized     []byte                      `json:"-"`
	XXX_sizecache        int32                       `json:"-"`
}

func (m *JunosPfeNpu) Reset()         { *m = JunosPfeNpu{} }
func (m *JunosPfeNpu) String() string { return proto.CompactTextString(m) }
func (*JunosPfeNpu) ProtoMessage()    {}
func (*JunosPfeNpu) Descriptor() ([]byte, []int) {
	return fileDescriptor_cdd8cb97f7b6deb1, []int{0}
}

func (m *JunosPfeNpu) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_JunosPfeNpu.Unmarshal(m, b)
}
func (m *JunosPfeNpu) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_JunosPfeNpu.Marshal(b, m, deterministic)
}
func (m *JunosPfeNpu) XXX_Merge(src proto.Message) {
	xxx_messageInfo_JunosPfeNpu.Merge(m, src)
}
func (m *JunosPfeNpu) XXX_Size() int {
	return xxx_messageInfo_JunosPfeNpu.Size(m)
}
func (m *JunosPfeNpu) XXX_DiscardUnknown() {
	xxx_messageInfo_JunosPfeNpu.DiscardUnknown(m)
}

var xxx_messageInfo_JunosPfeNpu proto.InternalMessageInfo

func (m *JunosPfeNpu) GetNpuMemory() []*JunosPfeNpuNpuMemoryList {
	if m != nil {
		return m.NpuMemory
	}
	return nil
}

type JunosPfeNpuNpuMemoryList struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *JunosPfeNpuNpuMemoryList) Reset()         { *m = JunosPfeNpuNpuMemoryList{} }
func (m *JunosPfeNpuNpuMemoryList) String() string { return proto.CompactTextString(m) }
func (*JunosPfeNpuNpuMemoryList) ProtoMessage()    {}
func (*JunosPfeNpuNpuMemoryList) Descriptor() ([]byte, []int) {
	return fileDescriptor_cdd8cb97f7b6deb1, []int{0, 0}
}

func (m *JunosPfeNpuNpuMemoryList) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_JunosPfeNpuNpuMemoryList.Unmarshal(m, b)
}
func (m *JunosPfeNpuNpuMemoryList) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_JunosPfeNpuNpuMemoryList.Marshal(b, m, deterministic)
}
func (m *JunosPfeNpuNpuMemoryList) XXX_Merge(src proto.Message) {
	xxx_messageInfo_JunosPfeNpuNpuMemoryList.Merge(m, src)
}
func (m *JunosPfeNpuNpuMemoryList) XXX_Size() int {
	return xxx_messageInfo_JunosPfeNpuNpuMemoryList.Size(m)
}
func (m *JunosPfeNpuNpuMemoryList) XXX_DiscardUnknown() {
	xxx_messageInfo_JunosPfeNpuNpuMemoryList.DiscardUnknown(m)
}

var xxx_messageInfo_JunosPfeNpuNpuMemoryList proto.InternalMessageInfo

var E_JnprJunosPfeNpuExt = &proto.ExtensionDesc{
	ExtendedType:  (*JuniperNetworksSensors)(nil),
	ExtensionType: (*JunosPfeNpu)(nil),
	Field:         59,
	Name:          "jnpr_junos_pfe_npu_ext",
	Tag:           "bytes,59,opt,name=jnpr_junos_pfe_npu_ext",
	Filename:      "pfe_npu_resource.proto",
}

func init() {
	proto.RegisterType((*JunosPfeNpu)(nil), "junos_pfe_npu")
	proto.RegisterType((*JunosPfeNpuNpuMemoryList)(nil), "junos_pfe_npu.npu_memory_list")
	proto.RegisterExtension(E_JnprJunosPfeNpuExt)
}

func init() { proto.RegisterFile("pfe_npu_resource.proto", fileDescriptor_cdd8cb97f7b6deb1) }

var fileDescriptor_cdd8cb97f7b6deb1 = []byte{
	// 189 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x2b, 0x48, 0x4b, 0x8d,
	0xcf, 0x2b, 0x28, 0x8d, 0x2f, 0x4a, 0x2d, 0xce, 0x2f, 0x2d, 0x4a, 0x4e, 0xd5, 0x2b, 0x28, 0xca,
	0x2f, 0xc9, 0x97, 0x12, 0x2e, 0x49, 0xcd, 0x49, 0xcd, 0x4d, 0x2d, 0x29, 0xaa, 0x8c, 0x2f, 0xc9,
	0x2f, 0x80, 0x08, 0x2a, 0x25, 0x71, 0xf1, 0x66, 0x95, 0xe6, 0xe5, 0x17, 0xc7, 0x43, 0x35, 0x09,
	0xd9, 0x71, 0x71, 0x81, 0xf4, 0xe6, 0xa6, 0xe6, 0xe6, 0x17, 0x55, 0x4a, 0x4c, 0x67, 0x54, 0x60,
	0xd6, 0xe0, 0x36, 0x92, 0xd3, 0x43, 0x51, 0xa4, 0x87, 0x50, 0x11, 0x9f, 0x93, 0x59, 0x5c, 0x12,
	0xc4, 0x99, 0x57, 0x50, 0xea, 0x0b, 0xe6, 0x4b, 0x09, 0x72, 0xf1, 0xa3, 0xc9, 0x5a, 0x45, 0x73,
	0x89, 0x65, 0xe5, 0x15, 0x14, 0xc5, 0xa3, 0x98, 0x11, 0x9f, 0x5a, 0x51, 0x22, 0x24, 0xae, 0xe7,
	0x55, 0x9a, 0x97, 0x59, 0x90, 0x5a, 0xe4, 0x97, 0x5a, 0x52, 0x9e, 0x5f, 0x94, 0x5d, 0x1c, 0x9c,
	0x9a, 0x57, 0x9c, 0x5f, 0x54, 0x2c, 0x61, 0xad, 0xc0, 0xa8, 0xc1, 0x6d, 0xc4, 0x87, 0x6a, 0x6d,
	0x90, 0x10, 0xc8, 0x18, 0x2f, 0x90, 0x50, 0x40, 0x5a, 0xaa, 0x5f, 0x41, 0xa9, 0x6b, 0x45, 0x09,
	0x20, 0x00, 0x00, 0xff, 0xff, 0x53, 0x09, 0xb2, 0xdb, 0xee, 0x00, 0x00, 0x00,
}
