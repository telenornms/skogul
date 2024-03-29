// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: pfe_per_prefix_v6_stats.proto

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

type NetworkInstancesPfePerPrefixV6 struct {
	NetworkInstance      []*NetworkInstancesPfePerPrefixV6NetworkInstanceList `protobuf:"bytes,205,rep,name=network_instance,json=networkInstance" json:"network_instance,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                                             `json:"-"`
	XXX_unrecognized     []byte                                               `json:"-"`
	XXX_sizecache        int32                                                `json:"-"`
}

func (m *NetworkInstancesPfePerPrefixV6) Reset()         { *m = NetworkInstancesPfePerPrefixV6{} }
func (m *NetworkInstancesPfePerPrefixV6) String() string { return proto.CompactTextString(m) }
func (*NetworkInstancesPfePerPrefixV6) ProtoMessage()    {}
func (*NetworkInstancesPfePerPrefixV6) Descriptor() ([]byte, []int) {
	return fileDescriptor_893d61141269cb38, []int{0}
}
func (m *NetworkInstancesPfePerPrefixV6) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NetworkInstancesPfePerPrefixV6.Unmarshal(m, b)
}
func (m *NetworkInstancesPfePerPrefixV6) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NetworkInstancesPfePerPrefixV6.Marshal(b, m, deterministic)
}
func (m *NetworkInstancesPfePerPrefixV6) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NetworkInstancesPfePerPrefixV6.Merge(m, src)
}
func (m *NetworkInstancesPfePerPrefixV6) XXX_Size() int {
	return xxx_messageInfo_NetworkInstancesPfePerPrefixV6.Size(m)
}
func (m *NetworkInstancesPfePerPrefixV6) XXX_DiscardUnknown() {
	xxx_messageInfo_NetworkInstancesPfePerPrefixV6.DiscardUnknown(m)
}

var xxx_messageInfo_NetworkInstancesPfePerPrefixV6 proto.InternalMessageInfo

func (m *NetworkInstancesPfePerPrefixV6) GetNetworkInstance() []*NetworkInstancesPfePerPrefixV6NetworkInstanceList {
	if m != nil {
		return m.NetworkInstance
	}
	return nil
}

type NetworkInstancesPfePerPrefixV6NetworkInstanceList struct {
	Name                 *string                                                    `protobuf:"bytes,206,opt,name=name" json:"name,omitempty"`
	Afts                 *NetworkInstancesPfePerPrefixV6NetworkInstanceListAftsType `protobuf:"bytes,205,opt,name=afts" json:"afts,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                                                   `json:"-"`
	XXX_unrecognized     []byte                                                     `json:"-"`
	XXX_sizecache        int32                                                      `json:"-"`
}

func (m *NetworkInstancesPfePerPrefixV6NetworkInstanceList) Reset() {
	*m = NetworkInstancesPfePerPrefixV6NetworkInstanceList{}
}
func (m *NetworkInstancesPfePerPrefixV6NetworkInstanceList) String() string {
	return proto.CompactTextString(m)
}
func (*NetworkInstancesPfePerPrefixV6NetworkInstanceList) ProtoMessage() {}
func (*NetworkInstancesPfePerPrefixV6NetworkInstanceList) Descriptor() ([]byte, []int) {
	return fileDescriptor_893d61141269cb38, []int{0, 0}
}
func (m *NetworkInstancesPfePerPrefixV6NetworkInstanceList) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NetworkInstancesPfePerPrefixV6NetworkInstanceList.Unmarshal(m, b)
}
func (m *NetworkInstancesPfePerPrefixV6NetworkInstanceList) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NetworkInstancesPfePerPrefixV6NetworkInstanceList.Marshal(b, m, deterministic)
}
func (m *NetworkInstancesPfePerPrefixV6NetworkInstanceList) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NetworkInstancesPfePerPrefixV6NetworkInstanceList.Merge(m, src)
}
func (m *NetworkInstancesPfePerPrefixV6NetworkInstanceList) XXX_Size() int {
	return xxx_messageInfo_NetworkInstancesPfePerPrefixV6NetworkInstanceList.Size(m)
}
func (m *NetworkInstancesPfePerPrefixV6NetworkInstanceList) XXX_DiscardUnknown() {
	xxx_messageInfo_NetworkInstancesPfePerPrefixV6NetworkInstanceList.DiscardUnknown(m)
}

var xxx_messageInfo_NetworkInstancesPfePerPrefixV6NetworkInstanceList proto.InternalMessageInfo

func (m *NetworkInstancesPfePerPrefixV6NetworkInstanceList) GetName() string {
	if m != nil && m.Name != nil {
		return *m.Name
	}
	return ""
}

func (m *NetworkInstancesPfePerPrefixV6NetworkInstanceList) GetAfts() *NetworkInstancesPfePerPrefixV6NetworkInstanceListAftsType {
	if m != nil {
		return m.Afts
	}
	return nil
}

type NetworkInstancesPfePerPrefixV6NetworkInstanceListAftsType struct {
	Ipv6Unicast          *NetworkInstancesPfePerPrefixV6NetworkInstanceListAftsTypeIpv6UnicastType `protobuf:"bytes,205,opt,name=ipv6_unicast,json=ipv6Unicast" json:"ipv6_unicast,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                                                                  `json:"-"`
	XXX_unrecognized     []byte                                                                    `json:"-"`
	XXX_sizecache        int32                                                                     `json:"-"`
}

func (m *NetworkInstancesPfePerPrefixV6NetworkInstanceListAftsType) Reset() {
	*m = NetworkInstancesPfePerPrefixV6NetworkInstanceListAftsType{}
}
func (m *NetworkInstancesPfePerPrefixV6NetworkInstanceListAftsType) String() string {
	return proto.CompactTextString(m)
}
func (*NetworkInstancesPfePerPrefixV6NetworkInstanceListAftsType) ProtoMessage() {}
func (*NetworkInstancesPfePerPrefixV6NetworkInstanceListAftsType) Descriptor() ([]byte, []int) {
	return fileDescriptor_893d61141269cb38, []int{0, 0, 0}
}
func (m *NetworkInstancesPfePerPrefixV6NetworkInstanceListAftsType) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NetworkInstancesPfePerPrefixV6NetworkInstanceListAftsType.Unmarshal(m, b)
}
func (m *NetworkInstancesPfePerPrefixV6NetworkInstanceListAftsType) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NetworkInstancesPfePerPrefixV6NetworkInstanceListAftsType.Marshal(b, m, deterministic)
}
func (m *NetworkInstancesPfePerPrefixV6NetworkInstanceListAftsType) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NetworkInstancesPfePerPrefixV6NetworkInstanceListAftsType.Merge(m, src)
}
func (m *NetworkInstancesPfePerPrefixV6NetworkInstanceListAftsType) XXX_Size() int {
	return xxx_messageInfo_NetworkInstancesPfePerPrefixV6NetworkInstanceListAftsType.Size(m)
}
func (m *NetworkInstancesPfePerPrefixV6NetworkInstanceListAftsType) XXX_DiscardUnknown() {
	xxx_messageInfo_NetworkInstancesPfePerPrefixV6NetworkInstanceListAftsType.DiscardUnknown(m)
}

var xxx_messageInfo_NetworkInstancesPfePerPrefixV6NetworkInstanceListAftsType proto.InternalMessageInfo

func (m *NetworkInstancesPfePerPrefixV6NetworkInstanceListAftsType) GetIpv6Unicast() *NetworkInstancesPfePerPrefixV6NetworkInstanceListAftsTypeIpv6UnicastType {
	if m != nil {
		return m.Ipv6Unicast
	}
	return nil
}

type NetworkInstancesPfePerPrefixV6NetworkInstanceListAftsTypeIpv6UnicastType struct {
	Ipv6Entry            []*NetworkInstancesPfePerPrefixV6NetworkInstanceListAftsTypeIpv6UnicastTypeIpv6EntryList `protobuf:"bytes,205,rep,name=ipv6_entry,json=ipv6Entry" json:"ipv6_entry,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                                                                                 `json:"-"`
	XXX_unrecognized     []byte                                                                                   `json:"-"`
	XXX_sizecache        int32                                                                                    `json:"-"`
}

func (m *NetworkInstancesPfePerPrefixV6NetworkInstanceListAftsTypeIpv6UnicastType) Reset() {
	*m = NetworkInstancesPfePerPrefixV6NetworkInstanceListAftsTypeIpv6UnicastType{}
}
func (m *NetworkInstancesPfePerPrefixV6NetworkInstanceListAftsTypeIpv6UnicastType) String() string {
	return proto.CompactTextString(m)
}
func (*NetworkInstancesPfePerPrefixV6NetworkInstanceListAftsTypeIpv6UnicastType) ProtoMessage() {}
func (*NetworkInstancesPfePerPrefixV6NetworkInstanceListAftsTypeIpv6UnicastType) Descriptor() ([]byte, []int) {
	return fileDescriptor_893d61141269cb38, []int{0, 0, 0, 0}
}
func (m *NetworkInstancesPfePerPrefixV6NetworkInstanceListAftsTypeIpv6UnicastType) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NetworkInstancesPfePerPrefixV6NetworkInstanceListAftsTypeIpv6UnicastType.Unmarshal(m, b)
}
func (m *NetworkInstancesPfePerPrefixV6NetworkInstanceListAftsTypeIpv6UnicastType) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NetworkInstancesPfePerPrefixV6NetworkInstanceListAftsTypeIpv6UnicastType.Marshal(b, m, deterministic)
}
func (m *NetworkInstancesPfePerPrefixV6NetworkInstanceListAftsTypeIpv6UnicastType) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NetworkInstancesPfePerPrefixV6NetworkInstanceListAftsTypeIpv6UnicastType.Merge(m, src)
}
func (m *NetworkInstancesPfePerPrefixV6NetworkInstanceListAftsTypeIpv6UnicastType) XXX_Size() int {
	return xxx_messageInfo_NetworkInstancesPfePerPrefixV6NetworkInstanceListAftsTypeIpv6UnicastType.Size(m)
}
func (m *NetworkInstancesPfePerPrefixV6NetworkInstanceListAftsTypeIpv6UnicastType) XXX_DiscardUnknown() {
	xxx_messageInfo_NetworkInstancesPfePerPrefixV6NetworkInstanceListAftsTypeIpv6UnicastType.DiscardUnknown(m)
}

var xxx_messageInfo_NetworkInstancesPfePerPrefixV6NetworkInstanceListAftsTypeIpv6UnicastType proto.InternalMessageInfo

func (m *NetworkInstancesPfePerPrefixV6NetworkInstanceListAftsTypeIpv6UnicastType) GetIpv6Entry() []*NetworkInstancesPfePerPrefixV6NetworkInstanceListAftsTypeIpv6UnicastTypeIpv6EntryList {
	if m != nil {
		return m.Ipv6Entry
	}
	return nil
}

type NetworkInstancesPfePerPrefixV6NetworkInstanceListAftsTypeIpv6UnicastTypeIpv6EntryList struct {
	Prefix               *string                                                                                         `protobuf:"bytes,207,opt,name=prefix" json:"prefix,omitempty"`
	State                *NetworkInstancesPfePerPrefixV6NetworkInstanceListAftsTypeIpv6UnicastTypeIpv6EntryListStateType `protobuf:"bytes,205,opt,name=state" json:"state,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                                                                                        `json:"-"`
	XXX_unrecognized     []byte                                                                                          `json:"-"`
	XXX_sizecache        int32                                                                                           `json:"-"`
}

func (m *NetworkInstancesPfePerPrefixV6NetworkInstanceListAftsTypeIpv6UnicastTypeIpv6EntryList) Reset() {
	*m = NetworkInstancesPfePerPrefixV6NetworkInstanceListAftsTypeIpv6UnicastTypeIpv6EntryList{}
}
func (m *NetworkInstancesPfePerPrefixV6NetworkInstanceListAftsTypeIpv6UnicastTypeIpv6EntryList) String() string {
	return proto.CompactTextString(m)
}
func (*NetworkInstancesPfePerPrefixV6NetworkInstanceListAftsTypeIpv6UnicastTypeIpv6EntryList) ProtoMessage() {
}
func (*NetworkInstancesPfePerPrefixV6NetworkInstanceListAftsTypeIpv6UnicastTypeIpv6EntryList) Descriptor() ([]byte, []int) {
	return fileDescriptor_893d61141269cb38, []int{0, 0, 0, 0, 0}
}
func (m *NetworkInstancesPfePerPrefixV6NetworkInstanceListAftsTypeIpv6UnicastTypeIpv6EntryList) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NetworkInstancesPfePerPrefixV6NetworkInstanceListAftsTypeIpv6UnicastTypeIpv6EntryList.Unmarshal(m, b)
}
func (m *NetworkInstancesPfePerPrefixV6NetworkInstanceListAftsTypeIpv6UnicastTypeIpv6EntryList) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NetworkInstancesPfePerPrefixV6NetworkInstanceListAftsTypeIpv6UnicastTypeIpv6EntryList.Marshal(b, m, deterministic)
}
func (m *NetworkInstancesPfePerPrefixV6NetworkInstanceListAftsTypeIpv6UnicastTypeIpv6EntryList) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NetworkInstancesPfePerPrefixV6NetworkInstanceListAftsTypeIpv6UnicastTypeIpv6EntryList.Merge(m, src)
}
func (m *NetworkInstancesPfePerPrefixV6NetworkInstanceListAftsTypeIpv6UnicastTypeIpv6EntryList) XXX_Size() int {
	return xxx_messageInfo_NetworkInstancesPfePerPrefixV6NetworkInstanceListAftsTypeIpv6UnicastTypeIpv6EntryList.Size(m)
}
func (m *NetworkInstancesPfePerPrefixV6NetworkInstanceListAftsTypeIpv6UnicastTypeIpv6EntryList) XXX_DiscardUnknown() {
	xxx_messageInfo_NetworkInstancesPfePerPrefixV6NetworkInstanceListAftsTypeIpv6UnicastTypeIpv6EntryList.DiscardUnknown(m)
}

var xxx_messageInfo_NetworkInstancesPfePerPrefixV6NetworkInstanceListAftsTypeIpv6UnicastTypeIpv6EntryList proto.InternalMessageInfo

func (m *NetworkInstancesPfePerPrefixV6NetworkInstanceListAftsTypeIpv6UnicastTypeIpv6EntryList) GetPrefix() string {
	if m != nil && m.Prefix != nil {
		return *m.Prefix
	}
	return ""
}

func (m *NetworkInstancesPfePerPrefixV6NetworkInstanceListAftsTypeIpv6UnicastTypeIpv6EntryList) GetState() *NetworkInstancesPfePerPrefixV6NetworkInstanceListAftsTypeIpv6UnicastTypeIpv6EntryListStateType {
	if m != nil {
		return m.State
	}
	return nil
}

type NetworkInstancesPfePerPrefixV6NetworkInstanceListAftsTypeIpv6UnicastTypeIpv6EntryListStateType struct {
	PacketsForwarded     *uint64  `protobuf:"varint,208,opt,name=packets_forwarded,json=packetsForwarded" json:"packets_forwarded,omitempty"`
	OctetsForwarded      *uint64  `protobuf:"varint,209,opt,name=octets_forwarded,json=octetsForwarded" json:"octets_forwarded,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *NetworkInstancesPfePerPrefixV6NetworkInstanceListAftsTypeIpv6UnicastTypeIpv6EntryListStateType) Reset() {
	*m = NetworkInstancesPfePerPrefixV6NetworkInstanceListAftsTypeIpv6UnicastTypeIpv6EntryListStateType{}
}
func (m *NetworkInstancesPfePerPrefixV6NetworkInstanceListAftsTypeIpv6UnicastTypeIpv6EntryListStateType) String() string {
	return proto.CompactTextString(m)
}
func (*NetworkInstancesPfePerPrefixV6NetworkInstanceListAftsTypeIpv6UnicastTypeIpv6EntryListStateType) ProtoMessage() {
}
func (*NetworkInstancesPfePerPrefixV6NetworkInstanceListAftsTypeIpv6UnicastTypeIpv6EntryListStateType) Descriptor() ([]byte, []int) {
	return fileDescriptor_893d61141269cb38, []int{0, 0, 0, 0, 0, 0}
}
func (m *NetworkInstancesPfePerPrefixV6NetworkInstanceListAftsTypeIpv6UnicastTypeIpv6EntryListStateType) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NetworkInstancesPfePerPrefixV6NetworkInstanceListAftsTypeIpv6UnicastTypeIpv6EntryListStateType.Unmarshal(m, b)
}
func (m *NetworkInstancesPfePerPrefixV6NetworkInstanceListAftsTypeIpv6UnicastTypeIpv6EntryListStateType) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NetworkInstancesPfePerPrefixV6NetworkInstanceListAftsTypeIpv6UnicastTypeIpv6EntryListStateType.Marshal(b, m, deterministic)
}
func (m *NetworkInstancesPfePerPrefixV6NetworkInstanceListAftsTypeIpv6UnicastTypeIpv6EntryListStateType) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NetworkInstancesPfePerPrefixV6NetworkInstanceListAftsTypeIpv6UnicastTypeIpv6EntryListStateType.Merge(m, src)
}
func (m *NetworkInstancesPfePerPrefixV6NetworkInstanceListAftsTypeIpv6UnicastTypeIpv6EntryListStateType) XXX_Size() int {
	return xxx_messageInfo_NetworkInstancesPfePerPrefixV6NetworkInstanceListAftsTypeIpv6UnicastTypeIpv6EntryListStateType.Size(m)
}
func (m *NetworkInstancesPfePerPrefixV6NetworkInstanceListAftsTypeIpv6UnicastTypeIpv6EntryListStateType) XXX_DiscardUnknown() {
	xxx_messageInfo_NetworkInstancesPfePerPrefixV6NetworkInstanceListAftsTypeIpv6UnicastTypeIpv6EntryListStateType.DiscardUnknown(m)
}

var xxx_messageInfo_NetworkInstancesPfePerPrefixV6NetworkInstanceListAftsTypeIpv6UnicastTypeIpv6EntryListStateType proto.InternalMessageInfo

func (m *NetworkInstancesPfePerPrefixV6NetworkInstanceListAftsTypeIpv6UnicastTypeIpv6EntryListStateType) GetPacketsForwarded() uint64 {
	if m != nil && m.PacketsForwarded != nil {
		return *m.PacketsForwarded
	}
	return 0
}

func (m *NetworkInstancesPfePerPrefixV6NetworkInstanceListAftsTypeIpv6UnicastTypeIpv6EntryListStateType) GetOctetsForwarded() uint64 {
	if m != nil && m.OctetsForwarded != nil {
		return *m.OctetsForwarded
	}
	return 0
}

var E_JnprNetworkInstancesPfePerPrefixV6Ext = &proto.ExtensionDesc{
	ExtendedType:  (*JuniperNetworksSensors)(nil),
	ExtensionType: (*NetworkInstancesPfePerPrefixV6)(nil),
	Field:         144,
	Name:          "jnpr_network_instances_pfe_per_prefix_v6_ext",
	Tag:           "bytes,144,opt,name=jnpr_network_instances_pfe_per_prefix_v6_ext",
	Filename:      "pfe_per_prefix_v6_stats.proto",
}

func init() {
	proto.RegisterType((*NetworkInstancesPfePerPrefixV6)(nil), "network_instances_pfe_per_prefix_v6")
	proto.RegisterType((*NetworkInstancesPfePerPrefixV6NetworkInstanceList)(nil), "network_instances_pfe_per_prefix_v6.network_instance_list")
	proto.RegisterType((*NetworkInstancesPfePerPrefixV6NetworkInstanceListAftsType)(nil), "network_instances_pfe_per_prefix_v6.network_instance_list.afts_type")
	proto.RegisterType((*NetworkInstancesPfePerPrefixV6NetworkInstanceListAftsTypeIpv6UnicastType)(nil), "network_instances_pfe_per_prefix_v6.network_instance_list.afts_type.ipv6_unicast_type")
	proto.RegisterType((*NetworkInstancesPfePerPrefixV6NetworkInstanceListAftsTypeIpv6UnicastTypeIpv6EntryList)(nil), "network_instances_pfe_per_prefix_v6.network_instance_list.afts_type.ipv6_unicast_type.ipv6_entry_list")
	proto.RegisterType((*NetworkInstancesPfePerPrefixV6NetworkInstanceListAftsTypeIpv6UnicastTypeIpv6EntryListStateType)(nil), "network_instances_pfe_per_prefix_v6.network_instance_list.afts_type.ipv6_unicast_type.ipv6_entry_list.state_type")
	proto.RegisterExtension(E_JnprNetworkInstancesPfePerPrefixV6Ext)
}

func init() { proto.RegisterFile("pfe_per_prefix_v6_stats.proto", fileDescriptor_893d61141269cb38) }

var fileDescriptor_893d61141269cb38 = []byte{
	// 426 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xbc, 0x93, 0xcd, 0x8a, 0x13, 0x41,
	0x14, 0x85, 0xa9, 0x99, 0xce, 0x60, 0x6e, 0x84, 0x64, 0x4a, 0xc4, 0xa6, 0x41, 0x19, 0xfc, 0x81,
	0x20, 0x43, 0x2f, 0x66, 0x91, 0xc5, 0xac, 0x44, 0x1c, 0x41, 0x17, 0x43, 0x68, 0x71, 0xc0, 0x55,
	0xd1, 0xf4, 0xdc, 0x82, 0x36, 0x49, 0x75, 0x59, 0x75, 0xf3, 0xe7, 0xd2, 0x85, 0xb8, 0x74, 0xef,
	0xd6, 0x37, 0x71, 0x1d, 0x7f, 0x5e, 0xc0, 0x47, 0xf0, 0x19, 0xa4, 0xba, 0xba, 0x13, 0x93, 0xb8,
	0x08, 0xa2, 0x2e, 0xf3, 0xdd, 0x93, 0x73, 0x4e, 0xdf, 0xdb, 0x0d, 0x37, 0xb5, 0x44, 0xa1, 0xd1,
	0x08, 0x6d, 0x50, 0xe6, 0x33, 0x31, 0xe9, 0x09, 0x4b, 0x29, 0xd9, 0x58, 0x9b, 0x82, 0x8a, 0xe8,
	0x1a, 0xe1, 0x10, 0x47, 0x48, 0x66, 0x2e, 0xa8, 0xd0, 0x1e, 0xde, 0xfe, 0x71, 0x00, 0x77, 0x14,
	0xd2, 0xb4, 0x30, 0x03, 0x91, 0x2b, 0x4b, 0xa9, 0xca, 0xd0, 0x8a, 0x2d, 0x23, 0x2e, 0xa1, 0xb3,
	0x29, 0x0b, 0x17, 0xec, 0x68, 0xbf, 0xdb, 0x3a, 0x39, 0x8d, 0x77, 0x30, 0xd8, 0xd2, 0x88, 0x61,
	0x6e, 0x29, 0x69, 0x57, 0xf8, 0x49, 0x45, 0xa3, 0xef, 0x0d, 0xb8, 0xfe, 0x5b, 0x29, 0x8f, 0x20,
	0x50, 0xe9, 0x08, 0xc3, 0xcf, 0xec, 0x88, 0x75, 0x9b, 0x0f, 0x1b, 0x6f, 0x1e, 0xec, 0x5d, 0x61,
	0x49, 0xc9, 0xf8, 0x0b, 0x08, 0x52, 0x49, 0xd6, 0x35, 0x62, 0xdd, 0xd6, 0xc9, 0xa3, 0x3f, 0x6f,
	0x14, 0x3b, 0x1f, 0x41, 0x73, 0x8d, 0x49, 0x69, 0x19, 0x2d, 0x02, 0x68, 0x2e, 0x19, 0x7f, 0x0d,
	0x57, 0x73, 0x3d, 0xe9, 0x89, 0xb1, 0xca, 0xb3, 0xd4, 0x52, 0x1d, 0x78, 0xf1, 0x37, 0x02, 0xe3,
	0x5f, 0x9d, 0x7d, 0x85, 0x96, 0x43, 0xcf, 0x3d, 0x89, 0x3e, 0xed, 0xc3, 0xe1, 0x96, 0x84, 0xbf,
	0x65, 0x00, 0x25, 0x45, 0x45, 0x66, 0x5e, 0xdf, 0x44, 0xfe, 0x9b, 0x42, 0xf1, 0x2a, 0xc9, 0xdf,
	0xaf, 0xe9, 0xc0, 0x99, 0xfb, 0x1d, 0x7d, 0xdc, 0x83, 0xf6, 0xc6, 0x98, 0xdf, 0x82, 0x03, 0x9f,
	0x16, 0x7e, 0x59, 0xbb, 0x5a, 0x45, 0xf9, 0x3b, 0x06, 0x0d, 0xf7, 0x8a, 0x62, 0xbd, 0xc8, 0x57,
	0xff, 0xa7, 0x77, 0x5c, 0x86, 0xfa, 0x1d, 0xfb, 0x02, 0x91, 0x04, 0x58, 0x41, 0x7e, 0x0c, 0x87,
	0x3a, 0xcd, 0x06, 0x48, 0x56, 0xc8, 0xc2, 0x4c, 0x53, 0x73, 0x89, 0x97, 0xe1, 0x57, 0xd7, 0x31,
	0x48, 0x3a, 0xd5, 0xe4, 0x71, 0x3d, 0xe0, 0xf7, 0xa1, 0x53, 0x64, 0xb4, 0x2e, 0xfe, 0xe6, 0xc5,
	0x6d, 0x3f, 0x58, 0x6a, 0x4f, 0x3f, 0x30, 0x38, 0x7e, 0xa9, 0xb4, 0x11, 0x3b, 0x3c, 0xa8, 0xc0,
	0x19, 0xf1, 0x1b, 0xf1, 0xd3, 0xb1, 0xca, 0x35, 0x9a, 0x73, 0xff, 0x07, 0xfb, 0x0c, 0x95, 0x2d,
	0x8c, 0x0d, 0xdf, 0xfb, 0x95, 0xdd, 0xdd, 0x65, 0x65, 0xc9, 0x3d, 0x97, 0x79, 0xbe, 0xfe, 0xb1,
	0xd9, 0xbe, 0xc4, 0x3e, 0x9a, 0x7e, 0x29, 0xba, 0xe8, 0x9d, 0xcd, 0xe8, 0x67, 0x00, 0x00, 0x00,
	0xff, 0xff, 0x2e, 0x3a, 0xf4, 0x80, 0x43, 0x04, 0x00, 0x00,
}
