// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: firewall.proto

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

//
// Top-level message
//
type Firewall struct {
	FirewallStats        []*FirewallStats `protobuf:"bytes,1,rep,name=firewall_stats,json=firewallStats" json:"firewall_stats,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *Firewall) Reset()         { *m = Firewall{} }
func (m *Firewall) String() string { return proto.CompactTextString(m) }
func (*Firewall) ProtoMessage()    {}
func (*Firewall) Descriptor() ([]byte, []int) {
	return fileDescriptor_00e54131a1710129, []int{0}
}
func (m *Firewall) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Firewall.Unmarshal(m, b)
}
func (m *Firewall) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Firewall.Marshal(b, m, deterministic)
}
func (m *Firewall) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Firewall.Merge(m, src)
}
func (m *Firewall) XXX_Size() int {
	return xxx_messageInfo_Firewall.Size(m)
}
func (m *Firewall) XXX_DiscardUnknown() {
	xxx_messageInfo_Firewall.DiscardUnknown(m)
}

var xxx_messageInfo_Firewall proto.InternalMessageInfo

func (m *Firewall) GetFirewallStats() []*FirewallStats {
	if m != nil {
		return m.FirewallStats
	}
	return nil
}

//
// Firewall filter statistics
//
type FirewallStats struct {
	FilterName *string `protobuf:"bytes,1,req,name=filter_name,json=filterName" json:"filter_name,omitempty"`
	// The Unix timestamp (seconds since 00:00:00 UTC 1970-01-01) of
	// last filter state change event such as filter add, filter change,
	// or counter clear.
	Timestamp                *uint64                     `protobuf:"varint,2,opt,name=timestamp" json:"timestamp,omitempty"`
	MemoryUsage              []*MemoryUsage              `protobuf:"bytes,3,rep,name=memory_usage,json=memoryUsage" json:"memory_usage,omitempty"`
	CounterStats             []*CounterStats             `protobuf:"bytes,4,rep,name=counter_stats,json=counterStats" json:"counter_stats,omitempty"`
	PolicerStats             []*PolicerStats             `protobuf:"bytes,5,rep,name=policer_stats,json=policerStats" json:"policer_stats,omitempty"`
	HierarchicalPolicerStats []*HierarchicalPolicerStats `protobuf:"bytes,6,rep,name=hierarchical_policer_stats,json=hierarchicalPolicerStats" json:"hierarchical_policer_stats,omitempty"`
	XXX_NoUnkeyedLiteral     struct{}                    `json:"-"`
	XXX_unrecognized         []byte                      `json:"-"`
	XXX_sizecache            int32                       `json:"-"`
}

func (m *FirewallStats) Reset()         { *m = FirewallStats{} }
func (m *FirewallStats) String() string { return proto.CompactTextString(m) }
func (*FirewallStats) ProtoMessage()    {}
func (*FirewallStats) Descriptor() ([]byte, []int) {
	return fileDescriptor_00e54131a1710129, []int{1}
}
func (m *FirewallStats) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FirewallStats.Unmarshal(m, b)
}
func (m *FirewallStats) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FirewallStats.Marshal(b, m, deterministic)
}
func (m *FirewallStats) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FirewallStats.Merge(m, src)
}
func (m *FirewallStats) XXX_Size() int {
	return xxx_messageInfo_FirewallStats.Size(m)
}
func (m *FirewallStats) XXX_DiscardUnknown() {
	xxx_messageInfo_FirewallStats.DiscardUnknown(m)
}

var xxx_messageInfo_FirewallStats proto.InternalMessageInfo

func (m *FirewallStats) GetFilterName() string {
	if m != nil && m.FilterName != nil {
		return *m.FilterName
	}
	return ""
}

func (m *FirewallStats) GetTimestamp() uint64 {
	if m != nil && m.Timestamp != nil {
		return *m.Timestamp
	}
	return 0
}

func (m *FirewallStats) GetMemoryUsage() []*MemoryUsage {
	if m != nil {
		return m.MemoryUsage
	}
	return nil
}

func (m *FirewallStats) GetCounterStats() []*CounterStats {
	if m != nil {
		return m.CounterStats
	}
	return nil
}

func (m *FirewallStats) GetPolicerStats() []*PolicerStats {
	if m != nil {
		return m.PolicerStats
	}
	return nil
}

func (m *FirewallStats) GetHierarchicalPolicerStats() []*HierarchicalPolicerStats {
	if m != nil {
		return m.HierarchicalPolicerStats
	}
	return nil
}

//
// Memory usage
//
type MemoryUsage struct {
	// The router has typically several types of memories (e.g. CPU's memory,
	// ASIC's forwarding memories) in which the firewall object is written.
	// This field indicates the name of the memory subsystem whose utilization
	// is being reported.
	Name *string `protobuf:"bytes,1,req,name=name" json:"name,omitempty"`
	// The amount of the memory allocated in bytes to the filter
	Allocated            *uint64  `protobuf:"varint,2,opt,name=allocated" json:"allocated,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MemoryUsage) Reset()         { *m = MemoryUsage{} }
func (m *MemoryUsage) String() string { return proto.CompactTextString(m) }
func (*MemoryUsage) ProtoMessage()    {}
func (*MemoryUsage) Descriptor() ([]byte, []int) {
	return fileDescriptor_00e54131a1710129, []int{2}
}
func (m *MemoryUsage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MemoryUsage.Unmarshal(m, b)
}
func (m *MemoryUsage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MemoryUsage.Marshal(b, m, deterministic)
}
func (m *MemoryUsage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MemoryUsage.Merge(m, src)
}
func (m *MemoryUsage) XXX_Size() int {
	return xxx_messageInfo_MemoryUsage.Size(m)
}
func (m *MemoryUsage) XXX_DiscardUnknown() {
	xxx_messageInfo_MemoryUsage.DiscardUnknown(m)
}

var xxx_messageInfo_MemoryUsage proto.InternalMessageInfo

func (m *MemoryUsage) GetName() string {
	if m != nil && m.Name != nil {
		return *m.Name
	}
	return ""
}

func (m *MemoryUsage) GetAllocated() uint64 {
	if m != nil && m.Allocated != nil {
		return *m.Allocated
	}
	return 0
}

//
// Counter statistics
//
type CounterStats struct {
	// Counter name
	Name *string `protobuf:"bytes,1,req,name=name" json:"name,omitempty"`
	// The total number of packets seen by the counter
	Packets *uint64 `protobuf:"varint,2,opt,name=packets" json:"packets,omitempty"`
	// The total number of bytes seen by the counter
	Bytes                *uint64  `protobuf:"varint,3,opt,name=bytes" json:"bytes,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CounterStats) Reset()         { *m = CounterStats{} }
func (m *CounterStats) String() string { return proto.CompactTextString(m) }
func (*CounterStats) ProtoMessage()    {}
func (*CounterStats) Descriptor() ([]byte, []int) {
	return fileDescriptor_00e54131a1710129, []int{3}
}
func (m *CounterStats) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CounterStats.Unmarshal(m, b)
}
func (m *CounterStats) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CounterStats.Marshal(b, m, deterministic)
}
func (m *CounterStats) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CounterStats.Merge(m, src)
}
func (m *CounterStats) XXX_Size() int {
	return xxx_messageInfo_CounterStats.Size(m)
}
func (m *CounterStats) XXX_DiscardUnknown() {
	xxx_messageInfo_CounterStats.DiscardUnknown(m)
}

var xxx_messageInfo_CounterStats proto.InternalMessageInfo

func (m *CounterStats) GetName() string {
	if m != nil && m.Name != nil {
		return *m.Name
	}
	return ""
}

func (m *CounterStats) GetPackets() uint64 {
	if m != nil && m.Packets != nil {
		return *m.Packets
	}
	return 0
}

func (m *CounterStats) GetBytes() uint64 {
	if m != nil && m.Bytes != nil {
		return *m.Bytes
	}
	return 0
}

//
// Policer statistics
//
type PolicerStats struct {
	// Policer instance name
	Name *string `protobuf:"bytes,1,req,name=name" json:"name,omitempty"`
	// The total number of packets marked out-of-specification by the policer
	OutOfSpecPackets *uint64 `protobuf:"varint,2,opt,name=out_of_spec_packets,json=outOfSpecPackets" json:"out_of_spec_packets,omitempty"`
	// The total number of bytes marked out-of-specification by the policer
	OutOfSpecBytes *uint64 `protobuf:"varint,3,opt,name=out_of_spec_bytes,json=outOfSpecBytes" json:"out_of_spec_bytes,omitempty"`
	// Additional statistics when enhanced policer statistics are available
	ExtendedPolicerStats *ExtendedPolicerStats `protobuf:"bytes,4,opt,name=extended_policer_stats,json=extendedPolicerStats" json:"extended_policer_stats,omitempty"`
	XXX_NoUnkeyedLiteral struct{}              `json:"-"`
	XXX_unrecognized     []byte                `json:"-"`
	XXX_sizecache        int32                 `json:"-"`
}

func (m *PolicerStats) Reset()         { *m = PolicerStats{} }
func (m *PolicerStats) String() string { return proto.CompactTextString(m) }
func (*PolicerStats) ProtoMessage()    {}
func (*PolicerStats) Descriptor() ([]byte, []int) {
	return fileDescriptor_00e54131a1710129, []int{4}
}
func (m *PolicerStats) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PolicerStats.Unmarshal(m, b)
}
func (m *PolicerStats) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PolicerStats.Marshal(b, m, deterministic)
}
func (m *PolicerStats) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PolicerStats.Merge(m, src)
}
func (m *PolicerStats) XXX_Size() int {
	return xxx_messageInfo_PolicerStats.Size(m)
}
func (m *PolicerStats) XXX_DiscardUnknown() {
	xxx_messageInfo_PolicerStats.DiscardUnknown(m)
}

var xxx_messageInfo_PolicerStats proto.InternalMessageInfo

func (m *PolicerStats) GetName() string {
	if m != nil && m.Name != nil {
		return *m.Name
	}
	return ""
}

func (m *PolicerStats) GetOutOfSpecPackets() uint64 {
	if m != nil && m.OutOfSpecPackets != nil {
		return *m.OutOfSpecPackets
	}
	return 0
}

func (m *PolicerStats) GetOutOfSpecBytes() uint64 {
	if m != nil && m.OutOfSpecBytes != nil {
		return *m.OutOfSpecBytes
	}
	return 0
}

func (m *PolicerStats) GetExtendedPolicerStats() *ExtendedPolicerStats {
	if m != nil {
		return m.ExtendedPolicerStats
	}
	return nil
}

//
// Extended policer statistics when enhanced policer statistics are available
//
type ExtendedPolicerStats struct {
	// The total number of packets subjected to policing
	OfferedPackets *uint64 `protobuf:"varint,1,opt,name=offered_packets,json=offeredPackets" json:"offered_packets,omitempty"`
	// The total number of bytes subjected to policing
	OfferedBytes *uint64 `protobuf:"varint,2,opt,name=offered_bytes,json=offeredBytes" json:"offered_bytes,omitempty"`
	// The total number of packets not discarded by the policer
	TransmittedPackets *uint64 `protobuf:"varint,3,opt,name=transmitted_packets,json=transmittedPackets" json:"transmitted_packets,omitempty"`
	// The total number of bytes not discarded by the policer
	TransmittedBytes     *uint64  `protobuf:"varint,4,opt,name=transmitted_bytes,json=transmittedBytes" json:"transmitted_bytes,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ExtendedPolicerStats) Reset()         { *m = ExtendedPolicerStats{} }
func (m *ExtendedPolicerStats) String() string { return proto.CompactTextString(m) }
func (*ExtendedPolicerStats) ProtoMessage()    {}
func (*ExtendedPolicerStats) Descriptor() ([]byte, []int) {
	return fileDescriptor_00e54131a1710129, []int{5}
}
func (m *ExtendedPolicerStats) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ExtendedPolicerStats.Unmarshal(m, b)
}
func (m *ExtendedPolicerStats) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ExtendedPolicerStats.Marshal(b, m, deterministic)
}
func (m *ExtendedPolicerStats) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ExtendedPolicerStats.Merge(m, src)
}
func (m *ExtendedPolicerStats) XXX_Size() int {
	return xxx_messageInfo_ExtendedPolicerStats.Size(m)
}
func (m *ExtendedPolicerStats) XXX_DiscardUnknown() {
	xxx_messageInfo_ExtendedPolicerStats.DiscardUnknown(m)
}

var xxx_messageInfo_ExtendedPolicerStats proto.InternalMessageInfo

func (m *ExtendedPolicerStats) GetOfferedPackets() uint64 {
	if m != nil && m.OfferedPackets != nil {
		return *m.OfferedPackets
	}
	return 0
}

func (m *ExtendedPolicerStats) GetOfferedBytes() uint64 {
	if m != nil && m.OfferedBytes != nil {
		return *m.OfferedBytes
	}
	return 0
}

func (m *ExtendedPolicerStats) GetTransmittedPackets() uint64 {
	if m != nil && m.TransmittedPackets != nil {
		return *m.TransmittedPackets
	}
	return 0
}

func (m *ExtendedPolicerStats) GetTransmittedBytes() uint64 {
	if m != nil && m.TransmittedBytes != nil {
		return *m.TransmittedBytes
	}
	return 0
}

//
// Hierarchical policer statistics
//
type HierarchicalPolicerStats struct {
	// Hierarchical policer instance name
	Name *string `protobuf:"bytes,1,req,name=name" json:"name,omitempty"`
	// The total number of packets marked out-of-specification by
	// the premium policer
	PremiumPackets *uint64 `protobuf:"varint,2,opt,name=premium_packets,json=premiumPackets" json:"premium_packets,omitempty"`
	// The total number of bytes marked out-of-specification by
	// the premium policer
	PremiumBytes *uint64 `protobuf:"varint,3,opt,name=premium_bytes,json=premiumBytes" json:"premium_bytes,omitempty"`
	// The total number of packets marked out-of-specification by
	// the aggregate policer
	AggregatePackets *uint64 `protobuf:"varint,4,opt,name=aggregate_packets,json=aggregatePackets" json:"aggregate_packets,omitempty"`
	// The total number of bytes marked out-of-specification by
	// the aggregate policer
	AggregateBytes       *uint64  `protobuf:"varint,5,opt,name=aggregate_bytes,json=aggregateBytes" json:"aggregate_bytes,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *HierarchicalPolicerStats) Reset()         { *m = HierarchicalPolicerStats{} }
func (m *HierarchicalPolicerStats) String() string { return proto.CompactTextString(m) }
func (*HierarchicalPolicerStats) ProtoMessage()    {}
func (*HierarchicalPolicerStats) Descriptor() ([]byte, []int) {
	return fileDescriptor_00e54131a1710129, []int{6}
}
func (m *HierarchicalPolicerStats) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_HierarchicalPolicerStats.Unmarshal(m, b)
}
func (m *HierarchicalPolicerStats) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_HierarchicalPolicerStats.Marshal(b, m, deterministic)
}
func (m *HierarchicalPolicerStats) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HierarchicalPolicerStats.Merge(m, src)
}
func (m *HierarchicalPolicerStats) XXX_Size() int {
	return xxx_messageInfo_HierarchicalPolicerStats.Size(m)
}
func (m *HierarchicalPolicerStats) XXX_DiscardUnknown() {
	xxx_messageInfo_HierarchicalPolicerStats.DiscardUnknown(m)
}

var xxx_messageInfo_HierarchicalPolicerStats proto.InternalMessageInfo

func (m *HierarchicalPolicerStats) GetName() string {
	if m != nil && m.Name != nil {
		return *m.Name
	}
	return ""
}

func (m *HierarchicalPolicerStats) GetPremiumPackets() uint64 {
	if m != nil && m.PremiumPackets != nil {
		return *m.PremiumPackets
	}
	return 0
}

func (m *HierarchicalPolicerStats) GetPremiumBytes() uint64 {
	if m != nil && m.PremiumBytes != nil {
		return *m.PremiumBytes
	}
	return 0
}

func (m *HierarchicalPolicerStats) GetAggregatePackets() uint64 {
	if m != nil && m.AggregatePackets != nil {
		return *m.AggregatePackets
	}
	return 0
}

func (m *HierarchicalPolicerStats) GetAggregateBytes() uint64 {
	if m != nil && m.AggregateBytes != nil {
		return *m.AggregateBytes
	}
	return 0
}

var E_JnprFirewallExt = &proto.ExtensionDesc{
	ExtendedType:  (*JuniperNetworksSensors)(nil),
	ExtensionType: (*Firewall)(nil),
	Field:         6,
	Name:          "jnpr_firewall_ext",
	Tag:           "bytes,6,opt,name=jnpr_firewall_ext",
	Filename:      "firewall.proto",
}

func init() {
	proto.RegisterType((*Firewall)(nil), "Firewall")
	proto.RegisterType((*FirewallStats)(nil), "FirewallStats")
	proto.RegisterType((*MemoryUsage)(nil), "MemoryUsage")
	proto.RegisterType((*CounterStats)(nil), "CounterStats")
	proto.RegisterType((*PolicerStats)(nil), "PolicerStats")
	proto.RegisterType((*ExtendedPolicerStats)(nil), "ExtendedPolicerStats")
	proto.RegisterType((*HierarchicalPolicerStats)(nil), "HierarchicalPolicerStats")
	proto.RegisterExtension(E_JnprFirewallExt)
}

func init() { proto.RegisterFile("firewall.proto", fileDescriptor_00e54131a1710129) }

var fileDescriptor_00e54131a1710129 = []byte{
	// 582 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x54, 0x4d, 0x6f, 0xd3, 0x40,
	0x10, 0x95, 0xf3, 0x01, 0xcd, 0xc4, 0x49, 0x9a, 0x4d, 0x01, 0xb7, 0x1c, 0x88, 0x8c, 0x84, 0x22,
	0x0e, 0x06, 0x45, 0xc0, 0x81, 0x13, 0x14, 0x15, 0x21, 0x50, 0x43, 0xe5, 0x08, 0x71, 0xb4, 0x16,
	0x67, 0x9c, 0x98, 0xda, 0x5e, 0x6b, 0xbd, 0x56, 0x93, 0x2b, 0x3f, 0x90, 0x9f, 0xc0, 0xbf, 0xe0,
	0x8e, 0xec, 0xb5, 0x9d, 0x35, 0x75, 0x73, 0x9c, 0xf7, 0xe6, 0xcd, 0xdb, 0xb7, 0x63, 0x2f, 0x0c,
	0x3d, 0x9f, 0xe3, 0x0d, 0x0d, 0x02, 0x2b, 0xe6, 0x4c, 0xb0, 0xb3, 0x89, 0xc0, 0x00, 0x43, 0x14,
	0x7c, 0xe7, 0x08, 0x16, 0x4b, 0xd0, 0x7c, 0x0f, 0x47, 0x1f, 0x8b, 0x36, 0xf2, 0x7a, 0x2f, 0x71,
	0x12, 0x41, 0x45, 0x62, 0x68, 0xd3, 0xf6, 0xac, 0x3f, 0x1f, 0x5a, 0x65, 0xcb, 0x32, 0x43, 0xed,
	0x81, 0xa7, 0x96, 0xe6, 0xef, 0x16, 0x0c, 0x6a, 0x0d, 0xe4, 0x19, 0xf4, 0x3d, 0x3f, 0x10, 0xc8,
	0x9d, 0x88, 0x86, 0x68, 0x68, 0xd3, 0xd6, 0xac, 0x77, 0xde, 0xfd, 0xf5, 0xae, 0x75, 0xa4, 0xd9,
	0x20, 0x99, 0x05, 0x0d, 0x91, 0x3c, 0x85, 0x9e, 0xf0, 0x43, 0x4c, 0x04, 0x0d, 0x63, 0xa3, 0x35,
	0xd5, 0x66, 0x9d, 0xbc, 0xeb, 0x58, 0xb3, 0xf7, 0x38, 0x79, 0x01, 0x7a, 0x88, 0x21, 0xe3, 0x3b,
	0x27, 0x4d, 0xe8, 0x1a, 0x8d, 0x76, 0x7e, 0x26, 0xdd, 0xba, 0xcc, 0xc1, 0x6f, 0x19, 0x66, 0xf7,
	0xc3, 0x7d, 0x41, 0xe6, 0x30, 0x70, 0x59, 0x1a, 0x65, 0xf6, 0x32, 0x45, 0x27, 0x57, 0x0c, 0xac,
	0x0f, 0x12, 0x95, 0x21, 0x74, 0x57, 0xa9, 0x32, 0x4d, 0xcc, 0x02, 0xdf, 0xad, 0x34, 0xdd, 0x42,
	0x73, 0x25, 0xd1, 0x42, 0x13, 0x2b, 0x15, 0xf9, 0x0e, 0x67, 0x1b, 0x1f, 0x39, 0xe5, 0xee, 0xc6,
	0x77, 0x69, 0xe0, 0xd4, 0x07, 0xdc, 0xcb, 0x07, 0x9c, 0x5a, 0x9f, 0x94, 0x96, 0xda, 0x30, 0x63,
	0x73, 0x07, 0x63, 0x5e, 0x42, 0x5f, 0x09, 0x47, 0x4e, 0xa1, 0x73, 0xfb, 0x1a, 0x73, 0x28, 0xbb,
	0x40, 0x1a, 0x04, 0xcc, 0xa5, 0x02, 0x57, 0xca, 0x05, 0x4e, 0x35, 0x7b, 0x8f, 0x9b, 0x6b, 0xd0,
	0xd5, 0xe4, 0x87, 0xe6, 0x3d, 0x81, 0xfb, 0x31, 0x75, 0xaf, 0x51, 0x24, 0xca, 0x34, 0x43, 0xb3,
	0x4b, 0x94, 0x3c, 0x86, 0xee, 0x8f, 0x9d, 0xc0, 0xc4, 0x68, 0xab, 0xb4, 0xc4, 0xcc, 0x3f, 0x1a,
	0xe8, 0x6a, 0x90, 0x43, 0x4e, 0xaf, 0x60, 0xc2, 0x52, 0xe1, 0x30, 0xcf, 0x49, 0x62, 0x74, 0x9d,
	0x46, 0xd7, 0x63, 0x96, 0x8a, 0xaf, 0xde, 0x32, 0x46, 0xf7, 0xaa, 0xb0, 0x7f, 0x09, 0x63, 0x55,
	0xd5, 0x70, 0x94, 0x61, 0xa5, 0x39, 0xcf, 0x48, 0xf2, 0x05, 0x1e, 0xe2, 0x56, 0x60, 0xb4, 0xc2,
	0xd5, 0x7f, 0x0b, 0xea, 0x4c, 0xb5, 0x59, 0x7f, 0xfe, 0xc0, 0xba, 0x28, 0xe8, 0xda, 0x72, 0x4e,
	0xb0, 0x01, 0xcd, 0x02, 0x9e, 0x34, 0xb5, 0x13, 0x0b, 0x46, 0xcc, 0xf3, 0x90, 0x67, 0x26, 0x45,
	0x12, 0xad, 0x7e, 0x2a, 0xc9, 0x96, 0x39, 0x9e, 0xc3, 0xa0, 0xec, 0x97, 0x19, 0x6a, 0xb9, 0xf5,
	0x82, 0x93, 0x09, 0xde, 0xc0, 0x44, 0x70, 0x1a, 0x25, 0xa1, 0x2f, 0x84, 0x32, 0xbf, 0x96, 0x9a,
	0x28, 0x1d, 0xa5, 0xc7, 0x1c, 0xc6, 0xaa, 0x4e, 0xfa, 0x74, 0x6a, 0xf7, 0xab, 0xf0, 0xb9, 0x97,
	0xf9, 0x57, 0x03, 0xe3, 0xae, 0x0f, 0xf6, 0xd0, 0x36, 0x2d, 0x18, 0xc5, 0x1c, 0x43, 0x3f, 0x0d,
	0x9b, 0x37, 0x39, 0x2c, 0x58, 0x25, 0x7f, 0xd9, 0xdf, 0xb0, 0x43, 0xbd, 0xe0, 0x64, 0xfe, 0x39,
	0x8c, 0xe9, 0x7a, 0xcd, 0x71, 0x4d, 0x05, 0x56, 0xd3, 0xeb, 0x39, 0x2a, 0xbe, 0x9c, 0x6f, 0xc1,
	0x68, 0xaf, 0x91, 0x0e, 0xdd, 0xda, 0x79, 0x2a, 0x36, 0xf7, 0x78, 0xbb, 0x80, 0xf1, 0xcf, 0x28,
	0xe6, 0x4e, 0xf5, 0xfc, 0xe1, 0x56, 0x90, 0x47, 0xd6, 0xe7, 0x34, 0xf2, 0x63, 0xe4, 0x0b, 0x14,
	0x37, 0x8c, 0x5f, 0x27, 0x4b, 0x8c, 0x12, 0xc6, 0xb3, 0x5f, 0x3b, 0xfb, 0x72, 0x7a, 0xd5, 0xab,
	0x68, 0x8f, 0x32, 0x71, 0x59, 0x5d, 0x6c, 0xc5, 0xbf, 0x00, 0x00, 0x00, 0xff, 0xff, 0x6a, 0x36,
	0x5c, 0xd8, 0x7b, 0x05, 0x00, 0x00,
}
