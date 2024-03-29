// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: l2ald_bd_render.proto

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

type NetworkInstancesBridgeDomain struct {
	NetworkInstance      []*NetworkInstancesBridgeDomainNetworkInstanceList `protobuf:"bytes,151,rep,name=network_instance,json=networkInstance" json:"network_instance,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                                           `json:"-"`
	XXX_unrecognized     []byte                                             `json:"-"`
	XXX_sizecache        int32                                              `json:"-"`
}

func (m *NetworkInstancesBridgeDomain) Reset()         { *m = NetworkInstancesBridgeDomain{} }
func (m *NetworkInstancesBridgeDomain) String() string { return proto.CompactTextString(m) }
func (*NetworkInstancesBridgeDomain) ProtoMessage()    {}
func (*NetworkInstancesBridgeDomain) Descriptor() ([]byte, []int) {
	return fileDescriptor_90265904e384351a, []int{0}
}
func (m *NetworkInstancesBridgeDomain) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NetworkInstancesBridgeDomain.Unmarshal(m, b)
}
func (m *NetworkInstancesBridgeDomain) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NetworkInstancesBridgeDomain.Marshal(b, m, deterministic)
}
func (m *NetworkInstancesBridgeDomain) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NetworkInstancesBridgeDomain.Merge(m, src)
}
func (m *NetworkInstancesBridgeDomain) XXX_Size() int {
	return xxx_messageInfo_NetworkInstancesBridgeDomain.Size(m)
}
func (m *NetworkInstancesBridgeDomain) XXX_DiscardUnknown() {
	xxx_messageInfo_NetworkInstancesBridgeDomain.DiscardUnknown(m)
}

var xxx_messageInfo_NetworkInstancesBridgeDomain proto.InternalMessageInfo

func (m *NetworkInstancesBridgeDomain) GetNetworkInstance() []*NetworkInstancesBridgeDomainNetworkInstanceList {
	if m != nil {
		return m.NetworkInstance
	}
	return nil
}

type NetworkInstancesBridgeDomainNetworkInstanceList struct {
	Name                 *string                                                    `protobuf:"bytes,51,opt,name=name" json:"name,omitempty"`
	Vlan                 []*NetworkInstancesBridgeDomainNetworkInstanceListVlanList `protobuf:"bytes,151,rep,name=vlan" json:"vlan,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                                                   `json:"-"`
	XXX_unrecognized     []byte                                                     `json:"-"`
	XXX_sizecache        int32                                                      `json:"-"`
}

func (m *NetworkInstancesBridgeDomainNetworkInstanceList) Reset() {
	*m = NetworkInstancesBridgeDomainNetworkInstanceList{}
}
func (m *NetworkInstancesBridgeDomainNetworkInstanceList) String() string {
	return proto.CompactTextString(m)
}
func (*NetworkInstancesBridgeDomainNetworkInstanceList) ProtoMessage() {}
func (*NetworkInstancesBridgeDomainNetworkInstanceList) Descriptor() ([]byte, []int) {
	return fileDescriptor_90265904e384351a, []int{0, 0}
}
func (m *NetworkInstancesBridgeDomainNetworkInstanceList) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NetworkInstancesBridgeDomainNetworkInstanceList.Unmarshal(m, b)
}
func (m *NetworkInstancesBridgeDomainNetworkInstanceList) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NetworkInstancesBridgeDomainNetworkInstanceList.Marshal(b, m, deterministic)
}
func (m *NetworkInstancesBridgeDomainNetworkInstanceList) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NetworkInstancesBridgeDomainNetworkInstanceList.Merge(m, src)
}
func (m *NetworkInstancesBridgeDomainNetworkInstanceList) XXX_Size() int {
	return xxx_messageInfo_NetworkInstancesBridgeDomainNetworkInstanceList.Size(m)
}
func (m *NetworkInstancesBridgeDomainNetworkInstanceList) XXX_DiscardUnknown() {
	xxx_messageInfo_NetworkInstancesBridgeDomainNetworkInstanceList.DiscardUnknown(m)
}

var xxx_messageInfo_NetworkInstancesBridgeDomainNetworkInstanceList proto.InternalMessageInfo

func (m *NetworkInstancesBridgeDomainNetworkInstanceList) GetName() string {
	if m != nil && m.Name != nil {
		return *m.Name
	}
	return ""
}

func (m *NetworkInstancesBridgeDomainNetworkInstanceList) GetVlan() []*NetworkInstancesBridgeDomainNetworkInstanceListVlanList {
	if m != nil {
		return m.Vlan
	}
	return nil
}

type NetworkInstancesBridgeDomainNetworkInstanceListVlanList struct {
	VlanName               *string                                                              `protobuf:"bytes,51,opt,name=vlan_name,json=vlanName" json:"vlan_name,omitempty"`
	VlanId                 *uint32                                                              `protobuf:"varint,52,opt,name=vlan_id,json=vlanId" json:"vlan_id,omitempty"`
	Status                 *string                                                              `protobuf:"bytes,53,opt,name=status" json:"status,omitempty"`
	Vni                    *uint32                                                              `protobuf:"varint,54,opt,name=vni" json:"vni,omitempty"`
	L3Interface            *string                                                              `protobuf:"bytes,55,opt,name=l3_interface,json=l3Interface" json:"l3_interface,omitempty"`
	NumLocalMacEntries     *uint32                                                              `protobuf:"varint,56,opt,name=num_local_mac_entries,json=numLocalMacEntries" json:"num_local_mac_entries,omitempty"`
	NumArReplicatorEntries *uint32                                                              `protobuf:"varint,57,opt,name=num_ar_replicator_entries,json=numArReplicatorEntries" json:"num_ar_replicator_entries,omitempty"`
	EthernetTagId          *uint32                                                              `protobuf:"varint,58,opt,name=ethernet_tag_id,json=ethernetTagId" json:"ethernet_tag_id,omitempty"`
	Member                 []*NetworkInstancesBridgeDomainNetworkInstanceListVlanListMemberList `protobuf:"bytes,161,rep,name=member" json:"member,omitempty"`
	XXX_NoUnkeyedLiteral   struct{}                                                             `json:"-"`
	XXX_unrecognized       []byte                                                               `json:"-"`
	XXX_sizecache          int32                                                                `json:"-"`
}

func (m *NetworkInstancesBridgeDomainNetworkInstanceListVlanList) Reset() {
	*m = NetworkInstancesBridgeDomainNetworkInstanceListVlanList{}
}
func (m *NetworkInstancesBridgeDomainNetworkInstanceListVlanList) String() string {
	return proto.CompactTextString(m)
}
func (*NetworkInstancesBridgeDomainNetworkInstanceListVlanList) ProtoMessage() {}
func (*NetworkInstancesBridgeDomainNetworkInstanceListVlanList) Descriptor() ([]byte, []int) {
	return fileDescriptor_90265904e384351a, []int{0, 0, 0}
}
func (m *NetworkInstancesBridgeDomainNetworkInstanceListVlanList) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NetworkInstancesBridgeDomainNetworkInstanceListVlanList.Unmarshal(m, b)
}
func (m *NetworkInstancesBridgeDomainNetworkInstanceListVlanList) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NetworkInstancesBridgeDomainNetworkInstanceListVlanList.Marshal(b, m, deterministic)
}
func (m *NetworkInstancesBridgeDomainNetworkInstanceListVlanList) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NetworkInstancesBridgeDomainNetworkInstanceListVlanList.Merge(m, src)
}
func (m *NetworkInstancesBridgeDomainNetworkInstanceListVlanList) XXX_Size() int {
	return xxx_messageInfo_NetworkInstancesBridgeDomainNetworkInstanceListVlanList.Size(m)
}
func (m *NetworkInstancesBridgeDomainNetworkInstanceListVlanList) XXX_DiscardUnknown() {
	xxx_messageInfo_NetworkInstancesBridgeDomainNetworkInstanceListVlanList.DiscardUnknown(m)
}

var xxx_messageInfo_NetworkInstancesBridgeDomainNetworkInstanceListVlanList proto.InternalMessageInfo

func (m *NetworkInstancesBridgeDomainNetworkInstanceListVlanList) GetVlanName() string {
	if m != nil && m.VlanName != nil {
		return *m.VlanName
	}
	return ""
}

func (m *NetworkInstancesBridgeDomainNetworkInstanceListVlanList) GetVlanId() uint32 {
	if m != nil && m.VlanId != nil {
		return *m.VlanId
	}
	return 0
}

func (m *NetworkInstancesBridgeDomainNetworkInstanceListVlanList) GetStatus() string {
	if m != nil && m.Status != nil {
		return *m.Status
	}
	return ""
}

func (m *NetworkInstancesBridgeDomainNetworkInstanceListVlanList) GetVni() uint32 {
	if m != nil && m.Vni != nil {
		return *m.Vni
	}
	return 0
}

func (m *NetworkInstancesBridgeDomainNetworkInstanceListVlanList) GetL3Interface() string {
	if m != nil && m.L3Interface != nil {
		return *m.L3Interface
	}
	return ""
}

func (m *NetworkInstancesBridgeDomainNetworkInstanceListVlanList) GetNumLocalMacEntries() uint32 {
	if m != nil && m.NumLocalMacEntries != nil {
		return *m.NumLocalMacEntries
	}
	return 0
}

func (m *NetworkInstancesBridgeDomainNetworkInstanceListVlanList) GetNumArReplicatorEntries() uint32 {
	if m != nil && m.NumArReplicatorEntries != nil {
		return *m.NumArReplicatorEntries
	}
	return 0
}

func (m *NetworkInstancesBridgeDomainNetworkInstanceListVlanList) GetEthernetTagId() uint32 {
	if m != nil && m.EthernetTagId != nil {
		return *m.EthernetTagId
	}
	return 0
}

func (m *NetworkInstancesBridgeDomainNetworkInstanceListVlanList) GetMember() []*NetworkInstancesBridgeDomainNetworkInstanceListVlanListMemberList {
	if m != nil {
		return m.Member
	}
	return nil
}

type NetworkInstancesBridgeDomainNetworkInstanceListVlanListMemberList struct {
	Interface            *string  `protobuf:"bytes,61,opt,name=interface" json:"interface,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *NetworkInstancesBridgeDomainNetworkInstanceListVlanListMemberList) Reset() {
	*m = NetworkInstancesBridgeDomainNetworkInstanceListVlanListMemberList{}
}
func (m *NetworkInstancesBridgeDomainNetworkInstanceListVlanListMemberList) String() string {
	return proto.CompactTextString(m)
}
func (*NetworkInstancesBridgeDomainNetworkInstanceListVlanListMemberList) ProtoMessage() {}
func (*NetworkInstancesBridgeDomainNetworkInstanceListVlanListMemberList) Descriptor() ([]byte, []int) {
	return fileDescriptor_90265904e384351a, []int{0, 0, 0, 0}
}
func (m *NetworkInstancesBridgeDomainNetworkInstanceListVlanListMemberList) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NetworkInstancesBridgeDomainNetworkInstanceListVlanListMemberList.Unmarshal(m, b)
}
func (m *NetworkInstancesBridgeDomainNetworkInstanceListVlanListMemberList) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NetworkInstancesBridgeDomainNetworkInstanceListVlanListMemberList.Marshal(b, m, deterministic)
}
func (m *NetworkInstancesBridgeDomainNetworkInstanceListVlanListMemberList) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NetworkInstancesBridgeDomainNetworkInstanceListVlanListMemberList.Merge(m, src)
}
func (m *NetworkInstancesBridgeDomainNetworkInstanceListVlanListMemberList) XXX_Size() int {
	return xxx_messageInfo_NetworkInstancesBridgeDomainNetworkInstanceListVlanListMemberList.Size(m)
}
func (m *NetworkInstancesBridgeDomainNetworkInstanceListVlanListMemberList) XXX_DiscardUnknown() {
	xxx_messageInfo_NetworkInstancesBridgeDomainNetworkInstanceListVlanListMemberList.DiscardUnknown(m)
}

var xxx_messageInfo_NetworkInstancesBridgeDomainNetworkInstanceListVlanListMemberList proto.InternalMessageInfo

func (m *NetworkInstancesBridgeDomainNetworkInstanceListVlanListMemberList) GetInterface() string {
	if m != nil && m.Interface != nil {
		return *m.Interface
	}
	return ""
}

var E_JnprNetworkInstancesBridgeDomainExt = &proto.ExtensionDesc{
	ExtendedType:  (*JuniperNetworksSensors)(nil),
	ExtensionType: (*NetworkInstancesBridgeDomain)(nil),
	Field:         115,
	Name:          "jnpr_network_instances_bridge_domain_ext",
	Tag:           "bytes,115,opt,name=jnpr_network_instances_bridge_domain_ext",
	Filename:      "l2ald_bd_render.proto",
}

func init() {
	proto.RegisterType((*NetworkInstancesBridgeDomain)(nil), "network_instances_bridge_domain")
	proto.RegisterType((*NetworkInstancesBridgeDomainNetworkInstanceList)(nil), "network_instances_bridge_domain.network_instance_list")
	proto.RegisterType((*NetworkInstancesBridgeDomainNetworkInstanceListVlanList)(nil), "network_instances_bridge_domain.network_instance_list.vlan_list")
	proto.RegisterType((*NetworkInstancesBridgeDomainNetworkInstanceListVlanListMemberList)(nil), "network_instances_bridge_domain.network_instance_list.vlan_list.member_list")
	proto.RegisterExtension(E_JnprNetworkInstancesBridgeDomainExt)
}

func init() { proto.RegisterFile("l2ald_bd_render.proto", fileDescriptor_90265904e384351a) }

var fileDescriptor_90265904e384351a = []byte{
	// 463 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x92, 0xc1, 0x6e, 0x13, 0x31,
	0x10, 0x86, 0xb5, 0x24, 0x0d, 0xcd, 0x84, 0xaa, 0x95, 0x51, 0x5a, 0x37, 0x42, 0x22, 0x14, 0x09,
	0xad, 0x84, 0xb4, 0x12, 0x09, 0x14, 0x5a, 0x09, 0xa9, 0x54, 0xf4, 0x10, 0x54, 0x72, 0x58, 0xe0,
	0x6c, 0x39, 0xbb, 0x43, 0x30, 0xd8, 0xde, 0x95, 0xed, 0x2d, 0xe5, 0xca, 0x13, 0x70, 0xe3, 0xcc,
	0x93, 0xf0, 0x00, 0x1c, 0x79, 0x21, 0x64, 0xef, 0x26, 0x05, 0x8a, 0x14, 0x89, 0xde, 0x76, 0xbe,
	0xf9, 0xff, 0xd9, 0xdf, 0x1e, 0x43, 0x5f, 0x8e, 0xb8, 0xcc, 0xd9, 0x2c, 0x67, 0x06, 0x75, 0x8e,
	0x26, 0x29, 0x4d, 0xe1, 0x8a, 0xc1, 0x4d, 0x87, 0x12, 0x15, 0x3a, 0xf3, 0x89, 0xb9, 0xa2, 0xac,
	0xe1, 0xde, 0xcf, 0x35, 0xb8, 0xad, 0xd1, 0x7d, 0x2c, 0xcc, 0x07, 0x26, 0xb4, 0x75, 0x5c, 0x67,
	0x68, 0xd9, 0xcc, 0x88, 0x7c, 0x8e, 0x2c, 0x2f, 0x14, 0x17, 0x9a, 0xcc, 0x60, 0xeb, 0x6f, 0x09,
	0xfd, 0x1a, 0x0d, 0x5b, 0x71, 0x6f, 0xb4, 0x9f, 0xac, 0x30, 0x5f, 0xea, 0x33, 0x29, 0xac, 0x4b,
	0x37, 0x1b, 0x3c, 0x69, 0xe8, 0xe0, 0x47, 0x1b, 0xfa, 0xff, 0x94, 0x92, 0x5d, 0x68, 0x6b, 0xae,
	0x90, 0x8e, 0x87, 0x51, 0xdc, 0x3d, 0x5e, 0xfb, 0x7c, 0x74, 0x6d, 0x3d, 0x4a, 0x03, 0x22, 0x6f,
	0xa0, 0x7d, 0x26, 0xb9, 0x5e, 0x84, 0x39, 0xfa, 0xbf, 0x30, 0x89, 0x9f, 0x51, 0xc7, 0x0a, 0xe3,
	0x06, 0xdf, 0x5b, 0xd0, 0x5d, 0x32, 0xb2, 0xd7, 0x14, 0x97, 0x43, 0xac, 0x7b, 0x3e, 0xf5, 0x41,
	0x76, 0xe0, 0x7a, 0xd0, 0x88, 0x9c, 0x3e, 0x1c, 0x46, 0xf1, 0x46, 0xda, 0xf1, 0xe5, 0x24, 0x27,
	0xdb, 0xd0, 0xb1, 0x8e, 0xbb, 0xca, 0xd2, 0x47, 0xde, 0x99, 0x36, 0x15, 0xd9, 0x82, 0xd6, 0x99,
	0x16, 0x74, 0x3f, 0x88, 0xfd, 0x27, 0xb9, 0x03, 0x37, 0xe4, 0x98, 0x09, 0xed, 0xd0, 0xbc, 0xe5,
	0x19, 0xd2, 0xc7, 0x41, 0xdf, 0x93, 0xe3, 0xc9, 0x02, 0x91, 0x07, 0xd0, 0xd7, 0x95, 0x62, 0xb2,
	0xc8, 0xb8, 0x64, 0x8a, 0x67, 0x0c, 0xb5, 0x33, 0x02, 0x2d, 0x7d, 0x12, 0xc6, 0x10, 0x5d, 0xa9,
	0x53, 0xdf, 0x7b, 0xc9, 0xb3, 0x93, 0xba, 0x43, 0x0e, 0x60, 0xd7, 0x5b, 0xb8, 0x61, 0x06, 0x4b,
	0x29, 0x32, 0xee, 0x0a, 0xb3, 0xb4, 0x1d, 0x04, 0xdb, 0xb6, 0xae, 0xd4, 0x33, 0x93, 0x2e, 0xdb,
	0x0b, 0xeb, 0x3d, 0xd8, 0x44, 0xf7, 0x0e, 0x8d, 0x46, 0xc7, 0x1c, 0x9f, 0xfb, 0xb3, 0x1d, 0x06,
	0xc3, 0xc6, 0x02, 0xbf, 0xe6, 0xf3, 0x49, 0x4e, 0x10, 0x3a, 0x0a, 0xd5, 0x0c, 0x0d, 0xfd, 0x56,
	0xaf, 0xe1, 0xf4, 0xaa, 0x6b, 0x48, 0xea, 0x79, 0xf5, 0x4a, 0x9a, 0xe1, 0x83, 0xfb, 0xd0, 0xfb,
	0x0d, 0x93, 0x5b, 0xd0, 0xbd, 0xb8, 0xab, 0xa7, 0xe1, 0xae, 0x2e, 0xc0, 0xe1, 0x97, 0x08, 0xe2,
	0xf7, 0xba, 0x34, 0x6c, 0x45, 0x12, 0x86, 0xe7, 0x8e, 0xec, 0x24, 0x2f, 0x2a, 0x2d, 0x4a, 0x34,
	0xd3, 0x5a, 0x6c, 0x5f, 0xa1, 0xb6, 0x85, 0xb1, 0xd4, 0x0e, 0xa3, 0xb8, 0x37, 0x1a, 0xae, 0x3a,
	0x4e, 0x7a, 0xd7, 0xff, 0x6a, 0xfa, 0xe7, 0x83, 0xb6, 0xc7, 0x41, 0xf2, 0x3c, 0x28, 0x4e, 0xce,
	0xdd, 0xaf, 0x00, 0x00, 0x00, 0xff, 0xff, 0x5b, 0x88, 0x9f, 0x64, 0x95, 0x03, 0x00, 0x00,
}
