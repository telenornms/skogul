// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: chassisd-junos-state-chassis-render.proto

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

type StateChassis_248 struct {
	Chassis              *StateChassis_248ChassisType `protobuf:"bytes,149,opt,name=chassis" json:"chassis,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                     `json:"-"`
	XXX_unrecognized     []byte                       `json:"-"`
	XXX_sizecache        int32                        `json:"-"`
}

func (m *StateChassis_248) Reset()         { *m = StateChassis_248{} }
func (m *StateChassis_248) String() string { return proto.CompactTextString(m) }
func (*StateChassis_248) ProtoMessage()    {}
func (*StateChassis_248) Descriptor() ([]byte, []int) {
	return fileDescriptor_dbd3664abf9ea2e5, []int{0}
}
func (m *StateChassis_248) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StateChassis_248.Unmarshal(m, b)
}
func (m *StateChassis_248) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StateChassis_248.Marshal(b, m, deterministic)
}
func (m *StateChassis_248) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StateChassis_248.Merge(m, src)
}
func (m *StateChassis_248) XXX_Size() int {
	return xxx_messageInfo_StateChassis_248.Size(m)
}
func (m *StateChassis_248) XXX_DiscardUnknown() {
	xxx_messageInfo_StateChassis_248.DiscardUnknown(m)
}

var xxx_messageInfo_StateChassis_248 proto.InternalMessageInfo

func (m *StateChassis_248) GetChassis() *StateChassis_248ChassisType {
	if m != nil {
		return m.Chassis
	}
	return nil
}

type StateChassis_248ChassisType struct {
	Modules              *StateChassis_248ChassisTypeModulesType `protobuf:"bytes,150,opt,name=modules" json:"modules,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                                `json:"-"`
	XXX_unrecognized     []byte                                  `json:"-"`
	XXX_sizecache        int32                                   `json:"-"`
}

func (m *StateChassis_248ChassisType) Reset()         { *m = StateChassis_248ChassisType{} }
func (m *StateChassis_248ChassisType) String() string { return proto.CompactTextString(m) }
func (*StateChassis_248ChassisType) ProtoMessage()    {}
func (*StateChassis_248ChassisType) Descriptor() ([]byte, []int) {
	return fileDescriptor_dbd3664abf9ea2e5, []int{0, 0}
}
func (m *StateChassis_248ChassisType) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StateChassis_248ChassisType.Unmarshal(m, b)
}
func (m *StateChassis_248ChassisType) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StateChassis_248ChassisType.Marshal(b, m, deterministic)
}
func (m *StateChassis_248ChassisType) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StateChassis_248ChassisType.Merge(m, src)
}
func (m *StateChassis_248ChassisType) XXX_Size() int {
	return xxx_messageInfo_StateChassis_248ChassisType.Size(m)
}
func (m *StateChassis_248ChassisType) XXX_DiscardUnknown() {
	xxx_messageInfo_StateChassis_248ChassisType.DiscardUnknown(m)
}

var xxx_messageInfo_StateChassis_248ChassisType proto.InternalMessageInfo

func (m *StateChassis_248ChassisType) GetModules() *StateChassis_248ChassisTypeModulesType {
	if m != nil {
		return m.Modules
	}
	return nil
}

type StateChassis_248ChassisTypeModulesType struct {
	Module               []*StateChassis_248ChassisTypeModulesTypeModuleList `protobuf:"bytes,151,rep,name=module" json:"module,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                                            `json:"-"`
	XXX_unrecognized     []byte                                              `json:"-"`
	XXX_sizecache        int32                                               `json:"-"`
}

func (m *StateChassis_248ChassisTypeModulesType) Reset() {
	*m = StateChassis_248ChassisTypeModulesType{}
}
func (m *StateChassis_248ChassisTypeModulesType) String() string { return proto.CompactTextString(m) }
func (*StateChassis_248ChassisTypeModulesType) ProtoMessage()    {}
func (*StateChassis_248ChassisTypeModulesType) Descriptor() ([]byte, []int) {
	return fileDescriptor_dbd3664abf9ea2e5, []int{0, 0, 0}
}
func (m *StateChassis_248ChassisTypeModulesType) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StateChassis_248ChassisTypeModulesType.Unmarshal(m, b)
}
func (m *StateChassis_248ChassisTypeModulesType) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StateChassis_248ChassisTypeModulesType.Marshal(b, m, deterministic)
}
func (m *StateChassis_248ChassisTypeModulesType) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StateChassis_248ChassisTypeModulesType.Merge(m, src)
}
func (m *StateChassis_248ChassisTypeModulesType) XXX_Size() int {
	return xxx_messageInfo_StateChassis_248ChassisTypeModulesType.Size(m)
}
func (m *StateChassis_248ChassisTypeModulesType) XXX_DiscardUnknown() {
	xxx_messageInfo_StateChassis_248ChassisTypeModulesType.DiscardUnknown(m)
}

var xxx_messageInfo_StateChassis_248ChassisTypeModulesType proto.InternalMessageInfo

func (m *StateChassis_248ChassisTypeModulesType) GetModule() []*StateChassis_248ChassisTypeModulesTypeModuleList {
	if m != nil {
		return m.Module
	}
	return nil
}

type StateChassis_248ChassisTypeModulesTypeModuleList struct {
	Name                 *string                                                           `protobuf:"bytes,152,opt,name=name" json:"name,omitempty"`
	MacAddresses         *StateChassis_248ChassisTypeModulesTypeModuleListMacAddressesType `protobuf:"bytes,153,opt,name=mac_addresses,json=macAddresses" json:"mac_addresses,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                                                          `json:"-"`
	XXX_unrecognized     []byte                                                            `json:"-"`
	XXX_sizecache        int32                                                             `json:"-"`
}

func (m *StateChassis_248ChassisTypeModulesTypeModuleList) Reset() {
	*m = StateChassis_248ChassisTypeModulesTypeModuleList{}
}
func (m *StateChassis_248ChassisTypeModulesTypeModuleList) String() string {
	return proto.CompactTextString(m)
}
func (*StateChassis_248ChassisTypeModulesTypeModuleList) ProtoMessage() {}
func (*StateChassis_248ChassisTypeModulesTypeModuleList) Descriptor() ([]byte, []int) {
	return fileDescriptor_dbd3664abf9ea2e5, []int{0, 0, 0, 0}
}
func (m *StateChassis_248ChassisTypeModulesTypeModuleList) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StateChassis_248ChassisTypeModulesTypeModuleList.Unmarshal(m, b)
}
func (m *StateChassis_248ChassisTypeModulesTypeModuleList) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StateChassis_248ChassisTypeModulesTypeModuleList.Marshal(b, m, deterministic)
}
func (m *StateChassis_248ChassisTypeModulesTypeModuleList) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StateChassis_248ChassisTypeModulesTypeModuleList.Merge(m, src)
}
func (m *StateChassis_248ChassisTypeModulesTypeModuleList) XXX_Size() int {
	return xxx_messageInfo_StateChassis_248ChassisTypeModulesTypeModuleList.Size(m)
}
func (m *StateChassis_248ChassisTypeModulesTypeModuleList) XXX_DiscardUnknown() {
	xxx_messageInfo_StateChassis_248ChassisTypeModulesTypeModuleList.DiscardUnknown(m)
}

var xxx_messageInfo_StateChassis_248ChassisTypeModulesTypeModuleList proto.InternalMessageInfo

func (m *StateChassis_248ChassisTypeModulesTypeModuleList) GetName() string {
	if m != nil && m.Name != nil {
		return *m.Name
	}
	return ""
}

func (m *StateChassis_248ChassisTypeModulesTypeModuleList) GetMacAddresses() *StateChassis_248ChassisTypeModulesTypeModuleListMacAddressesType {
	if m != nil {
		return m.MacAddresses
	}
	return nil
}

type StateChassis_248ChassisTypeModulesTypeModuleListMacAddressesType struct {
	PublicBaseAddress    *string  `protobuf:"bytes,154,opt,name=public_base_address,json=publicBaseAddress" json:"public_base_address,omitempty"`
	PublicAddressCount   *uint32  `protobuf:"varint,155,opt,name=public_address_count,json=publicAddressCount" json:"public_address_count,omitempty"`
	PrivateBaseAddress   *string  `protobuf:"bytes,156,opt,name=private_base_address,json=privateBaseAddress" json:"private_base_address,omitempty"`
	PrivateAddressCount  *uint32  `protobuf:"varint,157,opt,name=private_address_count,json=privateAddressCount" json:"private_address_count,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StateChassis_248ChassisTypeModulesTypeModuleListMacAddressesType) Reset() {
	*m = StateChassis_248ChassisTypeModulesTypeModuleListMacAddressesType{}
}
func (m *StateChassis_248ChassisTypeModulesTypeModuleListMacAddressesType) String() string {
	return proto.CompactTextString(m)
}
func (*StateChassis_248ChassisTypeModulesTypeModuleListMacAddressesType) ProtoMessage() {}
func (*StateChassis_248ChassisTypeModulesTypeModuleListMacAddressesType) Descriptor() ([]byte, []int) {
	return fileDescriptor_dbd3664abf9ea2e5, []int{0, 0, 0, 0, 0}
}
func (m *StateChassis_248ChassisTypeModulesTypeModuleListMacAddressesType) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StateChassis_248ChassisTypeModulesTypeModuleListMacAddressesType.Unmarshal(m, b)
}
func (m *StateChassis_248ChassisTypeModulesTypeModuleListMacAddressesType) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StateChassis_248ChassisTypeModulesTypeModuleListMacAddressesType.Marshal(b, m, deterministic)
}
func (m *StateChassis_248ChassisTypeModulesTypeModuleListMacAddressesType) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StateChassis_248ChassisTypeModulesTypeModuleListMacAddressesType.Merge(m, src)
}
func (m *StateChassis_248ChassisTypeModulesTypeModuleListMacAddressesType) XXX_Size() int {
	return xxx_messageInfo_StateChassis_248ChassisTypeModulesTypeModuleListMacAddressesType.Size(m)
}
func (m *StateChassis_248ChassisTypeModulesTypeModuleListMacAddressesType) XXX_DiscardUnknown() {
	xxx_messageInfo_StateChassis_248ChassisTypeModulesTypeModuleListMacAddressesType.DiscardUnknown(m)
}

var xxx_messageInfo_StateChassis_248ChassisTypeModulesTypeModuleListMacAddressesType proto.InternalMessageInfo

func (m *StateChassis_248ChassisTypeModulesTypeModuleListMacAddressesType) GetPublicBaseAddress() string {
	if m != nil && m.PublicBaseAddress != nil {
		return *m.PublicBaseAddress
	}
	return ""
}

func (m *StateChassis_248ChassisTypeModulesTypeModuleListMacAddressesType) GetPublicAddressCount() uint32 {
	if m != nil && m.PublicAddressCount != nil {
		return *m.PublicAddressCount
	}
	return 0
}

func (m *StateChassis_248ChassisTypeModulesTypeModuleListMacAddressesType) GetPrivateBaseAddress() string {
	if m != nil && m.PrivateBaseAddress != nil {
		return *m.PrivateBaseAddress
	}
	return ""
}

func (m *StateChassis_248ChassisTypeModulesTypeModuleListMacAddressesType) GetPrivateAddressCount() uint32 {
	if m != nil && m.PrivateAddressCount != nil {
		return *m.PrivateAddressCount
	}
	return 0
}

var E_JnprStateChassis_248Ext = &proto.ExtensionDesc{
	ExtendedType:  (*JuniperNetworksSensors)(nil),
	ExtensionType: (*StateChassis_248)(nil),
	Field:         248,
	Name:          "jnpr_state_chassis_248_ext",
	Tag:           "bytes,248,opt,name=jnpr_state_chassis_248_ext",
	Filename:      "chassisd-junos-state-chassis-render.proto",
}

func init() {
	proto.RegisterType((*StateChassis_248)(nil), "state_chassis_248")
	proto.RegisterType((*StateChassis_248ChassisType)(nil), "state_chassis_248.chassis_type")
	proto.RegisterType((*StateChassis_248ChassisTypeModulesType)(nil), "state_chassis_248.chassis_type.modules_type")
	proto.RegisterType((*StateChassis_248ChassisTypeModulesTypeModuleList)(nil), "state_chassis_248.chassis_type.modules_type.module_list")
	proto.RegisterType((*StateChassis_248ChassisTypeModulesTypeModuleListMacAddressesType)(nil), "state_chassis_248.chassis_type.modules_type.module_list.mac_addresses_type")
	proto.RegisterExtension(E_JnprStateChassis_248Ext)
}

func init() {
	proto.RegisterFile("chassisd-junos-state-chassis-render.proto", fileDescriptor_dbd3664abf9ea2e5)
}

var fileDescriptor_dbd3664abf9ea2e5 = []byte{
	// 403 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x50, 0x4b, 0x8e, 0xd3, 0x30,
	0x18, 0x56, 0xfa, 0xa0, 0xe0, 0xb6, 0x8b, 0xba, 0x3c, 0xac, 0x6c, 0xa8, 0x58, 0x15, 0x89, 0x04,
	0x11, 0xba, 0xa8, 0xba, 0x82, 0x56, 0x08, 0x09, 0x24, 0x24, 0xd2, 0x03, 0x58, 0x6e, 0xf2, 0x4b,
	0x4d, 0x49, 0xe2, 0x60, 0x3b, 0xd0, 0x6e, 0xb9, 0x00, 0x2b, 0xde, 0xcc, 0x51, 0xe6, 0x2e, 0x23,
	0xcd, 0x25, 0x66, 0x37, 0xa3, 0xc4, 0x8e, 0x94, 0x2a, 0x8b, 0xd1, 0xcc, 0xce, 0xff, 0xf7, 0xf4,
	0xff, 0xa3, 0xa7, 0xc1, 0x96, 0x49, 0x19, 0xc9, 0xd0, 0xd9, 0xe5, 0x29, 0x97, 0x8e, 0x54, 0x4c,
	0x81, 0x63, 0x40, 0x47, 0x40, 0x1a, 0x82, 0x70, 0x33, 0xc1, 0x15, 0xb7, 0xc7, 0x0a, 0x62, 0x48,
	0x40, 0x89, 0x03, 0x55, 0x3c, 0xd3, 0xe0, 0x93, 0xef, 0x5d, 0x34, 0x2a, 0x3d, 0xd4, 0x78, 0xa8,
	0x37, 0x9b, 0xe3, 0x05, 0xea, 0x99, 0x91, 0xfc, 0xb0, 0x26, 0xd6, 0xb4, 0xef, 0x3d, 0x76, 0x1b,
	0x2a, 0xb7, 0x7a, 0xab, 0x43, 0x06, 0x7e, 0x65, 0xb0, 0x4f, 0x3b, 0x68, 0x50, 0x67, 0xf0, 0x5b,
	0xd4, 0x4b, 0x78, 0x98, 0xc7, 0x20, 0xc9, 0x4f, 0x1d, 0xf6, 0xec, 0x9a, 0x30, 0xd7, 0xe8, 0x4d,
	0xb2, 0x99, 0xec, 0xf3, 0x36, 0x1a, 0xd4, 0x19, 0xfc, 0x11, 0xdd, 0xd1, 0x33, 0xf9, 0x65, 0x4d,
	0xda, 0xd3, 0xbe, 0x37, 0xbf, 0x49, 0xb0, 0x19, 0x68, 0x1c, 0x49, 0xe5, 0x9b, 0x20, 0xfb, 0xb2,
	0x85, 0xfa, 0x35, 0x1c, 0xdb, 0xa8, 0x93, 0xb2, 0x04, 0xc8, 0xef, 0xe2, 0xe7, 0xf7, 0x96, 0xdd,
	0x6f, 0xaf, 0x5a, 0x77, 0x2d, 0xbf, 0xc4, 0xf0, 0x67, 0x34, 0x4c, 0x58, 0x40, 0x59, 0x18, 0x0a,
	0x90, 0x12, 0x24, 0xf9, 0xa3, 0xd7, 0x7b, 0x7f, 0xdb, 0x5f, 0xb8, 0x47, 0x71, 0x7a, 0xfb, 0x41,
	0xc2, 0x82, 0xd7, 0x15, 0x64, 0x9f, 0x59, 0x08, 0x37, 0x45, 0xf8, 0x39, 0x1a, 0x67, 0xf9, 0x26,
	0x8e, 0x02, 0xba, 0x61, 0x12, 0x2a, 0x96, 0xfc, 0x2d, 0x3f, 0xed, 0x8f, 0x34, 0xb7, 0x64, 0x12,
	0x4c, 0x12, 0x7e, 0x81, 0xee, 0x1b, 0x83, 0xd1, 0xd2, 0x80, 0xe7, 0xa9, 0x22, 0xff, 0x0a, 0xc7,
	0xd0, 0xc7, 0x9a, 0x34, 0xea, 0x55, 0x41, 0x95, 0x16, 0x11, 0x7d, 0x29, 0x16, 0x3b, 0x2a, 0xf9,
	0xaf, 0x4b, 0xb0, 0x21, 0xeb, 0x2d, 0x2f, 0xd1, 0x83, 0xca, 0x72, 0x5c, 0x73, 0xa2, 0x6b, 0xc6,
	0x86, 0xad, 0xf7, 0x2c, 0xb6, 0xc8, 0xde, 0xa5, 0x99, 0xa0, 0x8d, 0x1b, 0x52, 0xd8, 0x2b, 0xfc,
	0xc8, 0x7d, 0x97, 0xa7, 0x51, 0x06, 0xe2, 0x03, 0xa8, 0xaf, 0x5c, 0x7c, 0x92, 0x6b, 0x48, 0x25,
	0x17, 0x92, 0x5c, 0xe8, 0xdb, 0xe3, 0xe6, 0xed, 0xfd, 0x87, 0x45, 0xde, 0xba, 0x80, 0x57, 0x1a,
	0xf5, 0x66, 0xf3, 0x37, 0x7b, 0x75, 0x15, 0x00, 0x00, 0xff, 0xff, 0xd7, 0xac, 0x72, 0x2e, 0x3c,
	0x03, 0x00, 0x00,
}