// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: spu_cpu_util.proto

package telemetry

import (
	fmt "fmt"
	proto "github.com/gogo/protobuf/proto"
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
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type JunosPfeSpuCpu struct {
	Security             *JunosPfeSpuCpuSecurityType `protobuf:"bytes,151,opt,name=security" json:"security,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                    `json:"-"`
	XXX_unrecognized     []byte                      `json:"-"`
	XXX_sizecache        int32                       `json:"-"`
}

func (m *JunosPfeSpuCpu) Reset()         { *m = JunosPfeSpuCpu{} }
func (m *JunosPfeSpuCpu) String() string { return proto.CompactTextString(m) }
func (*JunosPfeSpuCpu) ProtoMessage()    {}
func (*JunosPfeSpuCpu) Descriptor() ([]byte, []int) {
	return fileDescriptor_d0270bcf540ee61c, []int{0}
}
func (m *JunosPfeSpuCpu) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_JunosPfeSpuCpu.Unmarshal(m, b)
}
func (m *JunosPfeSpuCpu) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_JunosPfeSpuCpu.Marshal(b, m, deterministic)
}
func (m *JunosPfeSpuCpu) XXX_Merge(src proto.Message) {
	xxx_messageInfo_JunosPfeSpuCpu.Merge(m, src)
}
func (m *JunosPfeSpuCpu) XXX_Size() int {
	return xxx_messageInfo_JunosPfeSpuCpu.Size(m)
}
func (m *JunosPfeSpuCpu) XXX_DiscardUnknown() {
	xxx_messageInfo_JunosPfeSpuCpu.DiscardUnknown(m)
}

var xxx_messageInfo_JunosPfeSpuCpu proto.InternalMessageInfo

func (m *JunosPfeSpuCpu) GetSecurity() *JunosPfeSpuCpuSecurityType {
	if m != nil {
		return m.Security
	}
	return nil
}

type JunosPfeSpuCpuSecurityType struct {
	Spu                  *JunosPfeSpuCpuSecurityTypeSpuType `protobuf:"bytes,151,opt,name=spu" json:"spu,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                           `json:"-"`
	XXX_unrecognized     []byte                             `json:"-"`
	XXX_sizecache        int32                              `json:"-"`
}

func (m *JunosPfeSpuCpuSecurityType) Reset()         { *m = JunosPfeSpuCpuSecurityType{} }
func (m *JunosPfeSpuCpuSecurityType) String() string { return proto.CompactTextString(m) }
func (*JunosPfeSpuCpuSecurityType) ProtoMessage()    {}
func (*JunosPfeSpuCpuSecurityType) Descriptor() ([]byte, []int) {
	return fileDescriptor_d0270bcf540ee61c, []int{0, 0}
}
func (m *JunosPfeSpuCpuSecurityType) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_JunosPfeSpuCpuSecurityType.Unmarshal(m, b)
}
func (m *JunosPfeSpuCpuSecurityType) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_JunosPfeSpuCpuSecurityType.Marshal(b, m, deterministic)
}
func (m *JunosPfeSpuCpuSecurityType) XXX_Merge(src proto.Message) {
	xxx_messageInfo_JunosPfeSpuCpuSecurityType.Merge(m, src)
}
func (m *JunosPfeSpuCpuSecurityType) XXX_Size() int {
	return xxx_messageInfo_JunosPfeSpuCpuSecurityType.Size(m)
}
func (m *JunosPfeSpuCpuSecurityType) XXX_DiscardUnknown() {
	xxx_messageInfo_JunosPfeSpuCpuSecurityType.DiscardUnknown(m)
}

var xxx_messageInfo_JunosPfeSpuCpuSecurityType proto.InternalMessageInfo

func (m *JunosPfeSpuCpuSecurityType) GetSpu() *JunosPfeSpuCpuSecurityTypeSpuType {
	if m != nil {
		return m.Spu
	}
	return nil
}

type JunosPfeSpuCpuSecurityTypeSpuType struct {
	SpuName              *string                                     `protobuf:"bytes,51,opt,name=spu_name,json=spuName" json:"spu_name,omitempty"`
	Cpu                  []*JunosPfeSpuCpuSecurityTypeSpuTypeCpuList `protobuf:"bytes,151,rep,name=cpu" json:"cpu,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                                    `json:"-"`
	XXX_unrecognized     []byte                                      `json:"-"`
	XXX_sizecache        int32                                       `json:"-"`
}

func (m *JunosPfeSpuCpuSecurityTypeSpuType) Reset()         { *m = JunosPfeSpuCpuSecurityTypeSpuType{} }
func (m *JunosPfeSpuCpuSecurityTypeSpuType) String() string { return proto.CompactTextString(m) }
func (*JunosPfeSpuCpuSecurityTypeSpuType) ProtoMessage()    {}
func (*JunosPfeSpuCpuSecurityTypeSpuType) Descriptor() ([]byte, []int) {
	return fileDescriptor_d0270bcf540ee61c, []int{0, 0, 0}
}
func (m *JunosPfeSpuCpuSecurityTypeSpuType) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_JunosPfeSpuCpuSecurityTypeSpuType.Unmarshal(m, b)
}
func (m *JunosPfeSpuCpuSecurityTypeSpuType) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_JunosPfeSpuCpuSecurityTypeSpuType.Marshal(b, m, deterministic)
}
func (m *JunosPfeSpuCpuSecurityTypeSpuType) XXX_Merge(src proto.Message) {
	xxx_messageInfo_JunosPfeSpuCpuSecurityTypeSpuType.Merge(m, src)
}
func (m *JunosPfeSpuCpuSecurityTypeSpuType) XXX_Size() int {
	return xxx_messageInfo_JunosPfeSpuCpuSecurityTypeSpuType.Size(m)
}
func (m *JunosPfeSpuCpuSecurityTypeSpuType) XXX_DiscardUnknown() {
	xxx_messageInfo_JunosPfeSpuCpuSecurityTypeSpuType.DiscardUnknown(m)
}

var xxx_messageInfo_JunosPfeSpuCpuSecurityTypeSpuType proto.InternalMessageInfo

func (m *JunosPfeSpuCpuSecurityTypeSpuType) GetSpuName() string {
	if m != nil && m.SpuName != nil {
		return *m.SpuName
	}
	return ""
}

func (m *JunosPfeSpuCpuSecurityTypeSpuType) GetCpu() []*JunosPfeSpuCpuSecurityTypeSpuTypeCpuList {
	if m != nil {
		return m.Cpu
	}
	return nil
}

type JunosPfeSpuCpuSecurityTypeSpuTypeCpuList struct {
	CpuName              *string  `protobuf:"bytes,52,opt,name=cpu_name,json=cpuName" json:"cpu_name,omitempty"`
	Utilization          *uint32  `protobuf:"varint,53,opt,name=utilization" json:"utilization,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *JunosPfeSpuCpuSecurityTypeSpuTypeCpuList) Reset() {
	*m = JunosPfeSpuCpuSecurityTypeSpuTypeCpuList{}
}
func (m *JunosPfeSpuCpuSecurityTypeSpuTypeCpuList) String() string { return proto.CompactTextString(m) }
func (*JunosPfeSpuCpuSecurityTypeSpuTypeCpuList) ProtoMessage()    {}
func (*JunosPfeSpuCpuSecurityTypeSpuTypeCpuList) Descriptor() ([]byte, []int) {
	return fileDescriptor_d0270bcf540ee61c, []int{0, 0, 0, 0}
}
func (m *JunosPfeSpuCpuSecurityTypeSpuTypeCpuList) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_JunosPfeSpuCpuSecurityTypeSpuTypeCpuList.Unmarshal(m, b)
}
func (m *JunosPfeSpuCpuSecurityTypeSpuTypeCpuList) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_JunosPfeSpuCpuSecurityTypeSpuTypeCpuList.Marshal(b, m, deterministic)
}
func (m *JunosPfeSpuCpuSecurityTypeSpuTypeCpuList) XXX_Merge(src proto.Message) {
	xxx_messageInfo_JunosPfeSpuCpuSecurityTypeSpuTypeCpuList.Merge(m, src)
}
func (m *JunosPfeSpuCpuSecurityTypeSpuTypeCpuList) XXX_Size() int {
	return xxx_messageInfo_JunosPfeSpuCpuSecurityTypeSpuTypeCpuList.Size(m)
}
func (m *JunosPfeSpuCpuSecurityTypeSpuTypeCpuList) XXX_DiscardUnknown() {
	xxx_messageInfo_JunosPfeSpuCpuSecurityTypeSpuTypeCpuList.DiscardUnknown(m)
}

var xxx_messageInfo_JunosPfeSpuCpuSecurityTypeSpuTypeCpuList proto.InternalMessageInfo

func (m *JunosPfeSpuCpuSecurityTypeSpuTypeCpuList) GetCpuName() string {
	if m != nil && m.CpuName != nil {
		return *m.CpuName
	}
	return ""
}

func (m *JunosPfeSpuCpuSecurityTypeSpuTypeCpuList) GetUtilization() uint32 {
	if m != nil && m.Utilization != nil {
		return *m.Utilization
	}
	return 0
}

var E_JnprJunosPfeSpuCpuExt = &proto.ExtensionDesc{
	ExtendedType:  (*JuniperNetworksSensors)(nil),
	ExtensionType: (*JunosPfeSpuCpu)(nil),
	Field:         130,
	Name:          "jnpr_junos_pfe_spu_cpu_ext",
	Tag:           "bytes,130,opt,name=jnpr_junos_pfe_spu_cpu_ext",
	Filename:      "spu_cpu_util.proto",
}

func init() {
	proto.RegisterType((*JunosPfeSpuCpu)(nil), "junos_pfe_spu_cpu")
	proto.RegisterType((*JunosPfeSpuCpuSecurityType)(nil), "junos_pfe_spu_cpu.security_type")
	proto.RegisterType((*JunosPfeSpuCpuSecurityTypeSpuType)(nil), "junos_pfe_spu_cpu.security_type.spu_type")
	proto.RegisterType((*JunosPfeSpuCpuSecurityTypeSpuTypeCpuList)(nil), "junos_pfe_spu_cpu.security_type.spu_type.cpu_list")
	proto.RegisterExtension(E_JnprJunosPfeSpuCpuExt)
}

func init() { proto.RegisterFile("spu_cpu_util.proto", fileDescriptor_d0270bcf540ee61c) }

var fileDescriptor_d0270bcf540ee61c = []byte{
	// 302 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x90, 0xc1, 0x4a, 0x33, 0x31,
	0x14, 0x85, 0x99, 0x96, 0x9f, 0xce, 0x9f, 0xd2, 0x85, 0x11, 0x71, 0x9c, 0xd5, 0xe0, 0x6a, 0xdc,
	0x64, 0x51, 0x75, 0x23, 0x0a, 0xa2, 0x74, 0xd3, 0x45, 0x91, 0xf4, 0x01, 0x42, 0x19, 0x6e, 0x25,
	0xb5, 0x4d, 0x2e, 0xc9, 0x0d, 0xb6, 0x2e, 0xbb, 0x76, 0xef, 0xcb, 0xf8, 0x32, 0xbe, 0x89, 0xa4,
	0xd3, 0x11, 0x65, 0x16, 0xba, 0xbb, 0x7c, 0xb9, 0xe7, 0x9c, 0xdc, 0xc3, 0xb8, 0xc7, 0xa0, 0x2a,
	0x0c, 0x2a, 0x90, 0x5e, 0x0a, 0x74, 0x96, 0x6c, 0x7e, 0x48, 0xb0, 0x84, 0x15, 0x90, 0xdb, 0x28,
	0xb2, 0x58, 0xc3, 0xd3, 0x8f, 0x0e, 0x3b, 0x58, 0x04, 0x63, 0xbd, 0xc2, 0x39, 0xa8, 0xbd, 0x8a,
	0xdf, 0xb0, 0xd4, 0x43, 0x15, 0x9c, 0xa6, 0x4d, 0xf6, 0x96, 0x14, 0x49, 0xd9, 0x1f, 0x16, 0xa2,
	0xb5, 0x26, 0x9a, 0x1d, 0x45, 0x1b, 0x04, 0xf9, 0x25, 0xc9, 0x5f, 0x3b, 0x6c, 0xf0, 0xe3, 0x8d,
	0x5f, 0xb3, 0xae, 0xc7, 0xd0, 0x78, 0x9d, 0xfd, 0xe6, 0x25, 0x22, 0xdd, 0x99, 0x46, 0x59, 0xfe,
	0x9e, 0xb0, 0xb4, 0x21, 0xfc, 0xa4, 0x9e, 0xcd, 0x6c, 0x05, 0xd9, 0x79, 0x91, 0x94, 0xff, 0x65,
	0xcf, 0x63, 0x98, 0xcc, 0x56, 0xc0, 0x47, 0xac, 0x5b, 0xd5, 0x29, 0xdd, 0xb2, 0x3f, 0x1c, 0xfe,
	0x39, 0x45, 0xc4, 0xa6, 0x96, 0xda, 0x93, 0x8c, 0xfa, 0x7c, 0xc2, 0xd2, 0x06, 0xf0, 0xa2, 0x9e,
	0x77, 0x69, 0x17, 0x31, 0xed, 0xee, 0xdf, 0xf6, 0xb6, 0x93, 0x26, 0xb2, 0x57, 0xed, 0x43, 0x0b,
	0xd6, 0x8f, 0x25, 0xeb, 0x97, 0x19, 0x69, 0x6b, 0xb2, 0xcb, 0x22, 0x29, 0x07, 0xf2, 0x3b, 0xba,
	0x7a, 0x64, 0xf9, 0xc2, 0xa0, 0x53, 0xad, 0xef, 0x28, 0x58, 0x13, 0x3f, 0x16, 0xe3, 0x60, 0x34,
	0x82, 0x9b, 0x00, 0x3d, 0x5b, 0xf7, 0xe4, 0xa7, 0x60, 0xbc, 0x75, 0x3e, 0xdb, 0xd6, 0x65, 0xf1,
	0xf6, 0x19, 0xf2, 0x28, 0xfa, 0x8d, 0x23, 0x7e, 0x98, 0xc3, 0x14, 0xc3, 0x3d, 0x86, 0xd1, 0x9a,
	0x3e, 0x03, 0x00, 0x00, 0xff, 0xff, 0xa8, 0x3f, 0x0d, 0xd0, 0xf6, 0x01, 0x00, 0x00,
}
