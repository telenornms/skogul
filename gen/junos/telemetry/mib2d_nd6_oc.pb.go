// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: mib2d_nd6_oc.proto

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

type Nd6InformationMibNd6 struct {
	Ipv6                 *Nd6InformationMibNd6Ipv6Type `protobuf:"bytes,151,opt,name=ipv6" json:"ipv6,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                      `json:"-"`
	XXX_unrecognized     []byte                        `json:"-"`
	XXX_sizecache        int32                         `json:"-"`
}

func (m *Nd6InformationMibNd6) Reset()         { *m = Nd6InformationMibNd6{} }
func (m *Nd6InformationMibNd6) String() string { return proto.CompactTextString(m) }
func (*Nd6InformationMibNd6) ProtoMessage()    {}
func (*Nd6InformationMibNd6) Descriptor() ([]byte, []int) {
	return fileDescriptor_b6dca21213f5e620, []int{0}
}
func (m *Nd6InformationMibNd6) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Nd6InformationMibNd6.Unmarshal(m, b)
}
func (m *Nd6InformationMibNd6) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Nd6InformationMibNd6.Marshal(b, m, deterministic)
}
func (m *Nd6InformationMibNd6) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Nd6InformationMibNd6.Merge(m, src)
}
func (m *Nd6InformationMibNd6) XXX_Size() int {
	return xxx_messageInfo_Nd6InformationMibNd6.Size(m)
}
func (m *Nd6InformationMibNd6) XXX_DiscardUnknown() {
	xxx_messageInfo_Nd6InformationMibNd6.DiscardUnknown(m)
}

var xxx_messageInfo_Nd6InformationMibNd6 proto.InternalMessageInfo

func (m *Nd6InformationMibNd6) GetIpv6() *Nd6InformationMibNd6Ipv6Type {
	if m != nil {
		return m.Ipv6
	}
	return nil
}

type Nd6InformationMibNd6Ipv6Type struct {
	Neighbors            *Nd6InformationMibNd6Ipv6TypeNeighborsType `protobuf:"bytes,151,opt,name=neighbors" json:"neighbors,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                                   `json:"-"`
	XXX_unrecognized     []byte                                     `json:"-"`
	XXX_sizecache        int32                                      `json:"-"`
}

func (m *Nd6InformationMibNd6Ipv6Type) Reset()         { *m = Nd6InformationMibNd6Ipv6Type{} }
func (m *Nd6InformationMibNd6Ipv6Type) String() string { return proto.CompactTextString(m) }
func (*Nd6InformationMibNd6Ipv6Type) ProtoMessage()    {}
func (*Nd6InformationMibNd6Ipv6Type) Descriptor() ([]byte, []int) {
	return fileDescriptor_b6dca21213f5e620, []int{0, 0}
}
func (m *Nd6InformationMibNd6Ipv6Type) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Nd6InformationMibNd6Ipv6Type.Unmarshal(m, b)
}
func (m *Nd6InformationMibNd6Ipv6Type) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Nd6InformationMibNd6Ipv6Type.Marshal(b, m, deterministic)
}
func (m *Nd6InformationMibNd6Ipv6Type) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Nd6InformationMibNd6Ipv6Type.Merge(m, src)
}
func (m *Nd6InformationMibNd6Ipv6Type) XXX_Size() int {
	return xxx_messageInfo_Nd6InformationMibNd6Ipv6Type.Size(m)
}
func (m *Nd6InformationMibNd6Ipv6Type) XXX_DiscardUnknown() {
	xxx_messageInfo_Nd6InformationMibNd6Ipv6Type.DiscardUnknown(m)
}

var xxx_messageInfo_Nd6InformationMibNd6Ipv6Type proto.InternalMessageInfo

func (m *Nd6InformationMibNd6Ipv6Type) GetNeighbors() *Nd6InformationMibNd6Ipv6TypeNeighborsType {
	if m != nil {
		return m.Neighbors
	}
	return nil
}

type Nd6InformationMibNd6Ipv6TypeNeighborsType struct {
	Neighbor             []*Nd6InformationMibNd6Ipv6TypeNeighborsTypeNeighborList `protobuf:"bytes,151,rep,name=neighbor" json:"neighbor,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                                                 `json:"-"`
	XXX_unrecognized     []byte                                                   `json:"-"`
	XXX_sizecache        int32                                                    `json:"-"`
}

func (m *Nd6InformationMibNd6Ipv6TypeNeighborsType) Reset() {
	*m = Nd6InformationMibNd6Ipv6TypeNeighborsType{}
}
func (m *Nd6InformationMibNd6Ipv6TypeNeighborsType) String() string {
	return proto.CompactTextString(m)
}
func (*Nd6InformationMibNd6Ipv6TypeNeighborsType) ProtoMessage() {}
func (*Nd6InformationMibNd6Ipv6TypeNeighborsType) Descriptor() ([]byte, []int) {
	return fileDescriptor_b6dca21213f5e620, []int{0, 0, 0}
}
func (m *Nd6InformationMibNd6Ipv6TypeNeighborsType) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Nd6InformationMibNd6Ipv6TypeNeighborsType.Unmarshal(m, b)
}
func (m *Nd6InformationMibNd6Ipv6TypeNeighborsType) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Nd6InformationMibNd6Ipv6TypeNeighborsType.Marshal(b, m, deterministic)
}
func (m *Nd6InformationMibNd6Ipv6TypeNeighborsType) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Nd6InformationMibNd6Ipv6TypeNeighborsType.Merge(m, src)
}
func (m *Nd6InformationMibNd6Ipv6TypeNeighborsType) XXX_Size() int {
	return xxx_messageInfo_Nd6InformationMibNd6Ipv6TypeNeighborsType.Size(m)
}
func (m *Nd6InformationMibNd6Ipv6TypeNeighborsType) XXX_DiscardUnknown() {
	xxx_messageInfo_Nd6InformationMibNd6Ipv6TypeNeighborsType.DiscardUnknown(m)
}

var xxx_messageInfo_Nd6InformationMibNd6Ipv6TypeNeighborsType proto.InternalMessageInfo

func (m *Nd6InformationMibNd6Ipv6TypeNeighborsType) GetNeighbor() []*Nd6InformationMibNd6Ipv6TypeNeighborsTypeNeighborList {
	if m != nil {
		return m.Neighbor
	}
	return nil
}

type Nd6InformationMibNd6Ipv6TypeNeighborsTypeNeighborList struct {
	Ip                   *string                                                         `protobuf:"bytes,51,opt,name=ip" json:"ip,omitempty"`
	State                *Nd6InformationMibNd6Ipv6TypeNeighborsTypeNeighborListStateType `protobuf:"bytes,151,opt,name=state" json:"state,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                                                        `json:"-"`
	XXX_unrecognized     []byte                                                          `json:"-"`
	XXX_sizecache        int32                                                           `json:"-"`
}

func (m *Nd6InformationMibNd6Ipv6TypeNeighborsTypeNeighborList) Reset() {
	*m = Nd6InformationMibNd6Ipv6TypeNeighborsTypeNeighborList{}
}
func (m *Nd6InformationMibNd6Ipv6TypeNeighborsTypeNeighborList) String() string {
	return proto.CompactTextString(m)
}
func (*Nd6InformationMibNd6Ipv6TypeNeighborsTypeNeighborList) ProtoMessage() {}
func (*Nd6InformationMibNd6Ipv6TypeNeighborsTypeNeighborList) Descriptor() ([]byte, []int) {
	return fileDescriptor_b6dca21213f5e620, []int{0, 0, 0, 0}
}
func (m *Nd6InformationMibNd6Ipv6TypeNeighborsTypeNeighborList) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Nd6InformationMibNd6Ipv6TypeNeighborsTypeNeighborList.Unmarshal(m, b)
}
func (m *Nd6InformationMibNd6Ipv6TypeNeighborsTypeNeighborList) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Nd6InformationMibNd6Ipv6TypeNeighborsTypeNeighborList.Marshal(b, m, deterministic)
}
func (m *Nd6InformationMibNd6Ipv6TypeNeighborsTypeNeighborList) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Nd6InformationMibNd6Ipv6TypeNeighborsTypeNeighborList.Merge(m, src)
}
func (m *Nd6InformationMibNd6Ipv6TypeNeighborsTypeNeighborList) XXX_Size() int {
	return xxx_messageInfo_Nd6InformationMibNd6Ipv6TypeNeighborsTypeNeighborList.Size(m)
}
func (m *Nd6InformationMibNd6Ipv6TypeNeighborsTypeNeighborList) XXX_DiscardUnknown() {
	xxx_messageInfo_Nd6InformationMibNd6Ipv6TypeNeighborsTypeNeighborList.DiscardUnknown(m)
}

var xxx_messageInfo_Nd6InformationMibNd6Ipv6TypeNeighborsTypeNeighborList proto.InternalMessageInfo

func (m *Nd6InformationMibNd6Ipv6TypeNeighborsTypeNeighborList) GetIp() string {
	if m != nil && m.Ip != nil {
		return *m.Ip
	}
	return ""
}

func (m *Nd6InformationMibNd6Ipv6TypeNeighborsTypeNeighborList) GetState() *Nd6InformationMibNd6Ipv6TypeNeighborsTypeNeighborListStateType {
	if m != nil {
		return m.State
	}
	return nil
}

type Nd6InformationMibNd6Ipv6TypeNeighborsTypeNeighborListStateType struct {
	Ip                   *string  `protobuf:"bytes,51,opt,name=ip" json:"ip,omitempty"`
	LinkLayerAddress     *string  `protobuf:"bytes,52,opt,name=link_layer_address,json=linkLayerAddress" json:"link_layer_address,omitempty"`
	Origin               *string  `protobuf:"bytes,53,opt,name=origin" json:"origin,omitempty"`
	IsRouter             *bool    `protobuf:"varint,54,opt,name=is_router,json=isRouter" json:"is_router,omitempty"`
	NeighborState        *string  `protobuf:"bytes,55,opt,name=neighbor_state,json=neighborState" json:"neighbor_state,omitempty"`
	TableId              *uint32  `protobuf:"varint,61,opt,name=table_id,json=tableId" json:"table_id,omitempty"`
	IsSecure             *bool    `protobuf:"varint,62,opt,name=is_secure,json=isSecure" json:"is_secure,omitempty"`
	Expiry               *uint32  `protobuf:"varint,64,opt,name=expiry" json:"expiry,omitempty"`
	IsPublish            *bool    `protobuf:"varint,63,opt,name=is_publish,json=isPublish" json:"is_publish,omitempty"`
	InterfaceName        *string  `protobuf:"bytes,65,opt,name=interface_name,json=interfaceName" json:"interface_name,omitempty"`
	LogicalRouterId      *uint32  `protobuf:"varint,66,opt,name=logical_router_id,json=logicalRouterId" json:"logical_router_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Nd6InformationMibNd6Ipv6TypeNeighborsTypeNeighborListStateType) Reset() {
	*m = Nd6InformationMibNd6Ipv6TypeNeighborsTypeNeighborListStateType{}
}
func (m *Nd6InformationMibNd6Ipv6TypeNeighborsTypeNeighborListStateType) String() string {
	return proto.CompactTextString(m)
}
func (*Nd6InformationMibNd6Ipv6TypeNeighborsTypeNeighborListStateType) ProtoMessage() {}
func (*Nd6InformationMibNd6Ipv6TypeNeighborsTypeNeighborListStateType) Descriptor() ([]byte, []int) {
	return fileDescriptor_b6dca21213f5e620, []int{0, 0, 0, 0, 0}
}
func (m *Nd6InformationMibNd6Ipv6TypeNeighborsTypeNeighborListStateType) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Nd6InformationMibNd6Ipv6TypeNeighborsTypeNeighborListStateType.Unmarshal(m, b)
}
func (m *Nd6InformationMibNd6Ipv6TypeNeighborsTypeNeighborListStateType) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Nd6InformationMibNd6Ipv6TypeNeighborsTypeNeighborListStateType.Marshal(b, m, deterministic)
}
func (m *Nd6InformationMibNd6Ipv6TypeNeighborsTypeNeighborListStateType) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Nd6InformationMibNd6Ipv6TypeNeighborsTypeNeighborListStateType.Merge(m, src)
}
func (m *Nd6InformationMibNd6Ipv6TypeNeighborsTypeNeighborListStateType) XXX_Size() int {
	return xxx_messageInfo_Nd6InformationMibNd6Ipv6TypeNeighborsTypeNeighborListStateType.Size(m)
}
func (m *Nd6InformationMibNd6Ipv6TypeNeighborsTypeNeighborListStateType) XXX_DiscardUnknown() {
	xxx_messageInfo_Nd6InformationMibNd6Ipv6TypeNeighborsTypeNeighborListStateType.DiscardUnknown(m)
}

var xxx_messageInfo_Nd6InformationMibNd6Ipv6TypeNeighborsTypeNeighborListStateType proto.InternalMessageInfo

func (m *Nd6InformationMibNd6Ipv6TypeNeighborsTypeNeighborListStateType) GetIp() string {
	if m != nil && m.Ip != nil {
		return *m.Ip
	}
	return ""
}

func (m *Nd6InformationMibNd6Ipv6TypeNeighborsTypeNeighborListStateType) GetLinkLayerAddress() string {
	if m != nil && m.LinkLayerAddress != nil {
		return *m.LinkLayerAddress
	}
	return ""
}

func (m *Nd6InformationMibNd6Ipv6TypeNeighborsTypeNeighborListStateType) GetOrigin() string {
	if m != nil && m.Origin != nil {
		return *m.Origin
	}
	return ""
}

func (m *Nd6InformationMibNd6Ipv6TypeNeighborsTypeNeighborListStateType) GetIsRouter() bool {
	if m != nil && m.IsRouter != nil {
		return *m.IsRouter
	}
	return false
}

func (m *Nd6InformationMibNd6Ipv6TypeNeighborsTypeNeighborListStateType) GetNeighborState() string {
	if m != nil && m.NeighborState != nil {
		return *m.NeighborState
	}
	return ""
}

func (m *Nd6InformationMibNd6Ipv6TypeNeighborsTypeNeighborListStateType) GetTableId() uint32 {
	if m != nil && m.TableId != nil {
		return *m.TableId
	}
	return 0
}

func (m *Nd6InformationMibNd6Ipv6TypeNeighborsTypeNeighborListStateType) GetIsSecure() bool {
	if m != nil && m.IsSecure != nil {
		return *m.IsSecure
	}
	return false
}

func (m *Nd6InformationMibNd6Ipv6TypeNeighborsTypeNeighborListStateType) GetExpiry() uint32 {
	if m != nil && m.Expiry != nil {
		return *m.Expiry
	}
	return 0
}

func (m *Nd6InformationMibNd6Ipv6TypeNeighborsTypeNeighborListStateType) GetIsPublish() bool {
	if m != nil && m.IsPublish != nil {
		return *m.IsPublish
	}
	return false
}

func (m *Nd6InformationMibNd6Ipv6TypeNeighborsTypeNeighborListStateType) GetInterfaceName() string {
	if m != nil && m.InterfaceName != nil {
		return *m.InterfaceName
	}
	return ""
}

func (m *Nd6InformationMibNd6Ipv6TypeNeighborsTypeNeighborListStateType) GetLogicalRouterId() uint32 {
	if m != nil && m.LogicalRouterId != nil {
		return *m.LogicalRouterId
	}
	return 0
}

var E_JnprNd6InformationMibNd6Ext = &proto.ExtensionDesc{
	ExtendedType:  (*JuniperNetworksSensors)(nil),
	ExtensionType: (*Nd6InformationMibNd6)(nil),
	Field:         54,
	Name:          "jnpr_nd6_information_mib_nd6_ext",
	Tag:           "bytes,54,opt,name=jnpr_nd6_information_mib_nd6_ext",
	Filename:      "mib2d_nd6_oc.proto",
}

func init() {
	proto.RegisterType((*Nd6InformationMibNd6)(nil), "nd6_information_mib_nd6")
	proto.RegisterType((*Nd6InformationMibNd6Ipv6Type)(nil), "nd6_information_mib_nd6.ipv6_type")
	proto.RegisterType((*Nd6InformationMibNd6Ipv6TypeNeighborsType)(nil), "nd6_information_mib_nd6.ipv6_type.neighbors_type")
	proto.RegisterType((*Nd6InformationMibNd6Ipv6TypeNeighborsTypeNeighborList)(nil), "nd6_information_mib_nd6.ipv6_type.neighbors_type.neighbor_list")
	proto.RegisterType((*Nd6InformationMibNd6Ipv6TypeNeighborsTypeNeighborListStateType)(nil), "nd6_information_mib_nd6.ipv6_type.neighbors_type.neighbor_list.state_type")
	proto.RegisterExtension(E_JnprNd6InformationMibNd6Ext)
}

func init() { proto.RegisterFile("mib2d_nd6_oc.proto", fileDescriptor_b6dca21213f5e620) }

var fileDescriptor_b6dca21213f5e620 = []byte{
	// 495 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x52, 0xdf, 0x6a, 0x13, 0x4f,
	0x18, 0x65, 0xd3, 0xa6, 0xbf, 0xe4, 0x2b, 0xed, 0x4f, 0x47, 0xb4, 0x63, 0x8a, 0x10, 0x0a, 0x42,
	0x10, 0x59, 0x30, 0xea, 0x16, 0x04, 0xed, 0x1f, 0xf0, 0x22, 0x45, 0x43, 0xd9, 0xdc, 0x0a, 0xc3,
	0x6e, 0xf6, 0x6b, 0xfa, 0xd9, 0xdd, 0x99, 0x65, 0x66, 0xa2, 0xc9, 0xad, 0x2f, 0x21, 0xf8, 0x3e,
	0x3e, 0x8a, 0xe0, 0x03, 0xf8, 0x00, 0x32, 0xb3, 0xc9, 0x2e, 0xbd, 0x08, 0x22, 0x5e, 0x9e, 0x73,
	0xbe, 0x73, 0xe6, 0x7c, 0x33, 0x03, 0xac, 0xa0, 0x74, 0x98, 0x09, 0x99, 0x45, 0x42, 0x4d, 0xc3,
	0x52, 0x2b, 0xab, 0x7a, 0xf7, 0x2c, 0xe6, 0x58, 0xa0, 0xd5, 0x4b, 0x61, 0x55, 0x59, 0x91, 0x47,
	0xdf, 0x76, 0xe0, 0xc0, 0x4d, 0x91, 0xbc, 0x52, 0xba, 0x48, 0x2c, 0x29, 0x29, 0x0a, 0x4a, 0x9d,
	0x93, 0x1d, 0xc3, 0x36, 0x95, 0x9f, 0x22, 0xfe, 0x35, 0xe8, 0x07, 0x83, 0xdd, 0xe1, 0x51, 0xb8,
	0x61, 0x30, 0x74, 0x53, 0xc2, 0x2e, 0x4b, 0x8c, 0xbd, 0xa1, 0xf7, 0xbd, 0x0d, 0xdd, 0x9a, 0x63,
	0x97, 0xd0, 0x95, 0x48, 0xb3, 0xeb, 0x54, 0x69, 0xb3, 0xce, 0x7a, 0xf6, 0xe7, 0xac, 0xb0, 0x36,
	0x55, 0xd1, 0x4d, 0x48, 0xef, 0xc7, 0x36, 0xec, 0xdf, 0x56, 0xd9, 0x07, 0xe8, 0xac, 0x19, 0x77,
	0xc6, 0xd6, 0x60, 0x77, 0x78, 0xf2, 0xd7, 0x67, 0xd4, 0x50, 0xe4, 0x64, 0x6c, 0x5c, 0x27, 0xf6,
	0x7e, 0x6d, 0xc1, 0xde, 0x2d, 0x8d, 0xdd, 0x87, 0x16, 0x95, 0xfc, 0x79, 0x3f, 0x18, 0x74, 0xcf,
	0xdb, 0x5f, 0x4e, 0x5b, 0x9d, 0x20, 0x6e, 0x51, 0xc9, 0x12, 0x68, 0x1b, 0x9b, 0x58, 0x5c, 0xef,
	0x79, 0xf1, 0x8f, 0x1d, 0x42, 0x9f, 0x56, 0x5d, 0x40, 0x95, 0xdc, 0xfb, 0xd9, 0x02, 0x68, 0x58,
	0xb6, 0xdf, 0x14, 0xf1, 0x0d, 0x9e, 0x02, 0xcb, 0x49, 0xde, 0x88, 0x3c, 0x59, 0xa2, 0x16, 0x49,
	0x96, 0x69, 0x34, 0x86, 0xbf, 0xf0, 0xfa, 0x1d, 0xa7, 0xbc, 0x73, 0xc2, 0x59, 0xc5, 0xb3, 0x07,
	0xb0, 0xa3, 0x34, 0xcd, 0x48, 0xf2, 0x97, 0x7e, 0x62, 0x85, 0xd8, 0x21, 0x74, 0xc9, 0x08, 0xad,
	0xe6, 0x16, 0x35, 0x8f, 0xfa, 0xc1, 0xa0, 0x13, 0x77, 0xc8, 0xc4, 0x1e, 0xb3, 0xc7, 0xcd, 0xed,
	0x8b, 0x6a, 0xdb, 0x63, 0x6f, 0xae, 0xaf, 0x68, 0xe2, 0x48, 0xf6, 0x10, 0x3a, 0x36, 0x49, 0x73,
	0x14, 0x94, 0xf1, 0xd7, 0xfd, 0x60, 0xb0, 0x17, 0xff, 0xe7, 0xf1, 0x28, 0x5b, 0xc5, 0x1b, 0x9c,
	0xce, 0x35, 0xf2, 0x37, 0xeb, 0xf8, 0x89, 0xc7, 0xae, 0x13, 0x2e, 0x4a, 0xd2, 0x4b, 0x7e, 0xea,
	0x5d, 0x2b, 0xc4, 0x1e, 0x01, 0x90, 0x11, 0xe5, 0x3c, 0xcd, 0xc9, 0x5c, 0xf3, 0x13, 0xef, 0xea,
	0x92, 0xb9, 0xac, 0x08, 0xd7, 0x8a, 0xa4, 0x45, 0x7d, 0x95, 0x4c, 0x51, 0xc8, 0xa4, 0x40, 0x7e,
	0x56, 0xb5, 0xaa, 0xd9, 0x71, 0x52, 0x20, 0x7b, 0x02, 0x77, 0x73, 0x35, 0xa3, 0x69, 0x92, 0xaf,
	0xd6, 0x73, 0xf5, 0xce, 0xfd, 0x41, 0xff, 0xaf, 0x84, 0x6a, 0xcd, 0x51, 0xf6, 0x6a, 0x01, 0xfd,
	0x8f, 0xb2, 0xd4, 0x62, 0xc3, 0x1b, 0x0a, 0x5c, 0x58, 0x76, 0x10, 0x5e, 0xcc, 0x25, 0x95, 0xa8,
	0xc7, 0x68, 0x3f, 0x2b, 0x7d, 0x63, 0x26, 0x28, 0x8d, 0xfb, 0xe9, 0x91, 0xff, 0x00, 0x7c, 0xd3,
	0x07, 0x88, 0x0f, 0x5d, 0xf4, 0x38, 0x8b, 0x46, 0x8d, 0xf6, 0x9e, 0xd2, 0x71, 0x16, 0xbd, 0x5d,
	0xd8, 0xdf, 0x01, 0x00, 0x00, 0xff, 0xff, 0xaa, 0xf7, 0x42, 0x2e, 0xc0, 0x03, 0x00, 0x00,
}
