// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: logical_port.proto

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
type LogicalPort struct {
	InterfaceInfo        []*LogicalInterfaceInfo `protobuf:"bytes,1,rep,name=interface_info,json=interfaceInfo" json:"interface_info,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                `json:"-"`
	XXX_unrecognized     []byte                  `json:"-"`
	XXX_sizecache        int32                   `json:"-"`
}

func (m *LogicalPort) Reset()         { *m = LogicalPort{} }
func (m *LogicalPort) String() string { return proto.CompactTextString(m) }
func (*LogicalPort) ProtoMessage()    {}
func (*LogicalPort) Descriptor() ([]byte, []int) {
	return fileDescriptor_ed53654dcd9b9a05, []int{0}
}
func (m *LogicalPort) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LogicalPort.Unmarshal(m, b)
}
func (m *LogicalPort) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LogicalPort.Marshal(b, m, deterministic)
}
func (m *LogicalPort) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LogicalPort.Merge(m, src)
}
func (m *LogicalPort) XXX_Size() int {
	return xxx_messageInfo_LogicalPort.Size(m)
}
func (m *LogicalPort) XXX_DiscardUnknown() {
	xxx_messageInfo_LogicalPort.DiscardUnknown(m)
}

var xxx_messageInfo_LogicalPort proto.InternalMessageInfo

func (m *LogicalPort) GetInterfaceInfo() []*LogicalInterfaceInfo {
	if m != nil {
		return m.InterfaceInfo
	}
	return nil
}

//
// Logical Interaface information
//
type LogicalInterfaceInfo struct {
	// Logical interface name (e.g. xe-0/0/0.0)
	IfName *string `protobuf:"bytes,1,req,name=if_name,json=ifName" json:"if_name,omitempty"`
	// Time reset
	InitTime *uint64 `protobuf:"varint,2,req,name=init_time,json=initTime" json:"init_time,omitempty"`
	// Global Index
	SnmpIfIndex *uint32 `protobuf:"varint,3,opt,name=snmp_if_index,json=snmpIfIndex" json:"snmp_if_index,omitempty"`
	// Name of the aggregate bundle
	ParentAeName *string `protobuf:"bytes,4,opt,name=parent_ae_name,json=parentAeName" json:"parent_ae_name,omitempty"`
	// Inbound traffic statistics
	IngressStats *IngressInterfaceStats `protobuf:"bytes,5,opt,name=ingress_stats,json=ingressStats" json:"ingress_stats,omitempty"`
	// Outbound traffic statistics
	EgressStats *EgressInterfaceStats `protobuf:"bytes,6,opt,name=egress_stats,json=egressStats" json:"egress_stats,omitempty"`
	// Link state UP\DOWN etc.
	OpState *OperationalState `protobuf:"bytes,7,opt,name=op_state,json=opState" json:"op_state,omitempty"`
	// administrative status, i.e.. enabled/disabled
	AdministractiveStatus *string `protobuf:"bytes,8,opt,name=administractive_status,json=administractiveStatus" json:"administractive_status,omitempty"`
	// Description of the interface
	Description *string `protobuf:"bytes,9,opt,name=description" json:"description,omitempty"`
	// This corresponds to the ifLastChange object in the standard interface MIB
	LastChange *uint32 `protobuf:"varint,10,opt,name=last_change,json=lastChange" json:"last_change,omitempty"`
	// This corresponds to the ifHighSpeed object in the standard interface MIB
	HighSpeed *uint32 `protobuf:"varint,11,opt,name=high_speed,json=highSpeed" json:"high_speed,omitempty"`
	// Ingress queue information
	IngressQueueInfo []*LogicalInterfaceQueueStats `protobuf:"bytes,12,rep,name=ingress_queue_info,json=ingressQueueInfo" json:"ingress_queue_info,omitempty"`
	// Egress queue information
	EgressQueueInfo          []*LogicalInterfaceQueueStats `protobuf:"bytes,13,rep,name=egress_queue_info,json=egressQueueInfo" json:"egress_queue_info,omitempty"`
	AggregatedInstanceMember *string                       `protobuf:"bytes,14,opt,name=aggregated_instance_member,json=aggregatedInstanceMember" json:"aggregated_instance_member,omitempty"`
	XXX_NoUnkeyedLiteral     struct{}                      `json:"-"`
	XXX_unrecognized         []byte                        `json:"-"`
	XXX_sizecache            int32                         `json:"-"`
}

func (m *LogicalInterfaceInfo) Reset()         { *m = LogicalInterfaceInfo{} }
func (m *LogicalInterfaceInfo) String() string { return proto.CompactTextString(m) }
func (*LogicalInterfaceInfo) ProtoMessage()    {}
func (*LogicalInterfaceInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_ed53654dcd9b9a05, []int{1}
}
func (m *LogicalInterfaceInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LogicalInterfaceInfo.Unmarshal(m, b)
}
func (m *LogicalInterfaceInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LogicalInterfaceInfo.Marshal(b, m, deterministic)
}
func (m *LogicalInterfaceInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LogicalInterfaceInfo.Merge(m, src)
}
func (m *LogicalInterfaceInfo) XXX_Size() int {
	return xxx_messageInfo_LogicalInterfaceInfo.Size(m)
}
func (m *LogicalInterfaceInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_LogicalInterfaceInfo.DiscardUnknown(m)
}

var xxx_messageInfo_LogicalInterfaceInfo proto.InternalMessageInfo

func (m *LogicalInterfaceInfo) GetIfName() string {
	if m != nil && m.IfName != nil {
		return *m.IfName
	}
	return ""
}

func (m *LogicalInterfaceInfo) GetInitTime() uint64 {
	if m != nil && m.InitTime != nil {
		return *m.InitTime
	}
	return 0
}

func (m *LogicalInterfaceInfo) GetSnmpIfIndex() uint32 {
	if m != nil && m.SnmpIfIndex != nil {
		return *m.SnmpIfIndex
	}
	return 0
}

func (m *LogicalInterfaceInfo) GetParentAeName() string {
	if m != nil && m.ParentAeName != nil {
		return *m.ParentAeName
	}
	return ""
}

func (m *LogicalInterfaceInfo) GetIngressStats() *IngressInterfaceStats {
	if m != nil {
		return m.IngressStats
	}
	return nil
}

func (m *LogicalInterfaceInfo) GetEgressStats() *EgressInterfaceStats {
	if m != nil {
		return m.EgressStats
	}
	return nil
}

func (m *LogicalInterfaceInfo) GetOpState() *OperationalState {
	if m != nil {
		return m.OpState
	}
	return nil
}

func (m *LogicalInterfaceInfo) GetAdministractiveStatus() string {
	if m != nil && m.AdministractiveStatus != nil {
		return *m.AdministractiveStatus
	}
	return ""
}

func (m *LogicalInterfaceInfo) GetDescription() string {
	if m != nil && m.Description != nil {
		return *m.Description
	}
	return ""
}

func (m *LogicalInterfaceInfo) GetLastChange() uint32 {
	if m != nil && m.LastChange != nil {
		return *m.LastChange
	}
	return 0
}

func (m *LogicalInterfaceInfo) GetHighSpeed() uint32 {
	if m != nil && m.HighSpeed != nil {
		return *m.HighSpeed
	}
	return 0
}

func (m *LogicalInterfaceInfo) GetIngressQueueInfo() []*LogicalInterfaceQueueStats {
	if m != nil {
		return m.IngressQueueInfo
	}
	return nil
}

func (m *LogicalInterfaceInfo) GetEgressQueueInfo() []*LogicalInterfaceQueueStats {
	if m != nil {
		return m.EgressQueueInfo
	}
	return nil
}

func (m *LogicalInterfaceInfo) GetAggregatedInstanceMember() string {
	if m != nil && m.AggregatedInstanceMember != nil {
		return *m.AggregatedInstanceMember
	}
	return ""
}

//
//  Interface inbound/Ingress traffic statistics
//
type IngressInterfaceStats struct {
	// Count of packets
	IfPackets *uint64 `protobuf:"varint,1,req,name=if_packets,json=ifPackets" json:"if_packets,omitempty"`
	// Count of bytes
	IfOctets *uint64 `protobuf:"varint,2,req,name=if_octets,json=ifOctets" json:"if_octets,omitempty"`
	// Count of unicast packets
	IfUcastPackets *uint64 `protobuf:"varint,3,opt,name=if_ucast_packets,json=ifUcastPackets" json:"if_ucast_packets,omitempty"`
	// Count of multicast packets
	IfMcastPackets       *uint64                      `protobuf:"varint,4,opt,name=if_mcast_packets,json=ifMcastPackets" json:"if_mcast_packets,omitempty"`
	IfFcStats            []*ForwardingClassAccounting `protobuf:"bytes,5,rep,name=if_fc_stats,json=ifFcStats" json:"if_fc_stats,omitempty"`
	IfFaStats            []*FamilyAccounting          `protobuf:"bytes,6,rep,name=if_fa_stats,json=ifFaStats" json:"if_fa_stats,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                     `json:"-"`
	XXX_unrecognized     []byte                       `json:"-"`
	XXX_sizecache        int32                        `json:"-"`
}

func (m *IngressInterfaceStats) Reset()         { *m = IngressInterfaceStats{} }
func (m *IngressInterfaceStats) String() string { return proto.CompactTextString(m) }
func (*IngressInterfaceStats) ProtoMessage()    {}
func (*IngressInterfaceStats) Descriptor() ([]byte, []int) {
	return fileDescriptor_ed53654dcd9b9a05, []int{2}
}
func (m *IngressInterfaceStats) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_IngressInterfaceStats.Unmarshal(m, b)
}
func (m *IngressInterfaceStats) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_IngressInterfaceStats.Marshal(b, m, deterministic)
}
func (m *IngressInterfaceStats) XXX_Merge(src proto.Message) {
	xxx_messageInfo_IngressInterfaceStats.Merge(m, src)
}
func (m *IngressInterfaceStats) XXX_Size() int {
	return xxx_messageInfo_IngressInterfaceStats.Size(m)
}
func (m *IngressInterfaceStats) XXX_DiscardUnknown() {
	xxx_messageInfo_IngressInterfaceStats.DiscardUnknown(m)
}

var xxx_messageInfo_IngressInterfaceStats proto.InternalMessageInfo

func (m *IngressInterfaceStats) GetIfPackets() uint64 {
	if m != nil && m.IfPackets != nil {
		return *m.IfPackets
	}
	return 0
}

func (m *IngressInterfaceStats) GetIfOctets() uint64 {
	if m != nil && m.IfOctets != nil {
		return *m.IfOctets
	}
	return 0
}

func (m *IngressInterfaceStats) GetIfUcastPackets() uint64 {
	if m != nil && m.IfUcastPackets != nil {
		return *m.IfUcastPackets
	}
	return 0
}

func (m *IngressInterfaceStats) GetIfMcastPackets() uint64 {
	if m != nil && m.IfMcastPackets != nil {
		return *m.IfMcastPackets
	}
	return 0
}

func (m *IngressInterfaceStats) GetIfFcStats() []*ForwardingClassAccounting {
	if m != nil {
		return m.IfFcStats
	}
	return nil
}

func (m *IngressInterfaceStats) GetIfFaStats() []*FamilyAccounting {
	if m != nil {
		return m.IfFaStats
	}
	return nil
}

//
//  Interface outbound/Egress traffic statistics
//
type EgressInterfaceStats struct {
	// Count of packets
	IfPackets *uint64 `protobuf:"varint,1,req,name=if_packets,json=ifPackets" json:"if_packets,omitempty"`
	// Count of bytes
	IfOctets *uint64 `protobuf:"varint,2,req,name=if_octets,json=ifOctets" json:"if_octets,omitempty"`
	// Count of unicast packets
	IfUcastPackets *uint64 `protobuf:"varint,3,opt,name=if_ucast_packets,json=ifUcastPackets" json:"if_ucast_packets,omitempty"`
	// Count of multicast packets
	IfMcastPackets       *uint64             `protobuf:"varint,4,opt,name=if_mcast_packets,json=ifMcastPackets" json:"if_mcast_packets,omitempty"`
	IfFaStats            []*FamilyAccounting `protobuf:"bytes,5,rep,name=if_fa_stats,json=ifFaStats" json:"if_fa_stats,omitempty"`
	XXX_NoUnkeyedLiteral struct{}            `json:"-"`
	XXX_unrecognized     []byte              `json:"-"`
	XXX_sizecache        int32               `json:"-"`
}

func (m *EgressInterfaceStats) Reset()         { *m = EgressInterfaceStats{} }
func (m *EgressInterfaceStats) String() string { return proto.CompactTextString(m) }
func (*EgressInterfaceStats) ProtoMessage()    {}
func (*EgressInterfaceStats) Descriptor() ([]byte, []int) {
	return fileDescriptor_ed53654dcd9b9a05, []int{3}
}
func (m *EgressInterfaceStats) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EgressInterfaceStats.Unmarshal(m, b)
}
func (m *EgressInterfaceStats) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EgressInterfaceStats.Marshal(b, m, deterministic)
}
func (m *EgressInterfaceStats) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EgressInterfaceStats.Merge(m, src)
}
func (m *EgressInterfaceStats) XXX_Size() int {
	return xxx_messageInfo_EgressInterfaceStats.Size(m)
}
func (m *EgressInterfaceStats) XXX_DiscardUnknown() {
	xxx_messageInfo_EgressInterfaceStats.DiscardUnknown(m)
}

var xxx_messageInfo_EgressInterfaceStats proto.InternalMessageInfo

func (m *EgressInterfaceStats) GetIfPackets() uint64 {
	if m != nil && m.IfPackets != nil {
		return *m.IfPackets
	}
	return 0
}

func (m *EgressInterfaceStats) GetIfOctets() uint64 {
	if m != nil && m.IfOctets != nil {
		return *m.IfOctets
	}
	return 0
}

func (m *EgressInterfaceStats) GetIfUcastPackets() uint64 {
	if m != nil && m.IfUcastPackets != nil {
		return *m.IfUcastPackets
	}
	return 0
}

func (m *EgressInterfaceStats) GetIfMcastPackets() uint64 {
	if m != nil && m.IfMcastPackets != nil {
		return *m.IfMcastPackets
	}
	return 0
}

func (m *EgressInterfaceStats) GetIfFaStats() []*FamilyAccounting {
	if m != nil {
		return m.IfFaStats
	}
	return nil
}

//
//  Interface operational State details
//
type OperationalState struct {
	// If the link is up/down
	OperationalStatus    *string  `protobuf:"bytes,1,opt,name=operational_status,json=operationalStatus" json:"operational_status,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *OperationalState) Reset()         { *m = OperationalState{} }
func (m *OperationalState) String() string { return proto.CompactTextString(m) }
func (*OperationalState) ProtoMessage()    {}
func (*OperationalState) Descriptor() ([]byte, []int) {
	return fileDescriptor_ed53654dcd9b9a05, []int{4}
}
func (m *OperationalState) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_OperationalState.Unmarshal(m, b)
}
func (m *OperationalState) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_OperationalState.Marshal(b, m, deterministic)
}
func (m *OperationalState) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OperationalState.Merge(m, src)
}
func (m *OperationalState) XXX_Size() int {
	return xxx_messageInfo_OperationalState.Size(m)
}
func (m *OperationalState) XXX_DiscardUnknown() {
	xxx_messageInfo_OperationalState.DiscardUnknown(m)
}

var xxx_messageInfo_OperationalState proto.InternalMessageInfo

func (m *OperationalState) GetOperationalStatus() string {
	if m != nil && m.OperationalStatus != nil {
		return *m.OperationalStatus
	}
	return ""
}

//
//  Interface forwarding class accounting
//
type ForwardingClassAccounting struct {
	// Interface protocol
	IfFamily *string `protobuf:"bytes,1,opt,name=if_family,json=ifFamily" json:"if_family,omitempty"`
	// Forwarding class number
	FcNumber *uint32 `protobuf:"varint,2,opt,name=fc_number,json=fcNumber" json:"fc_number,omitempty"`
	// Count of packets
	IfPackets *uint64 `protobuf:"varint,3,opt,name=if_packets,json=ifPackets" json:"if_packets,omitempty"`
	// Count of bytes
	IfOctets             *uint64  `protobuf:"varint,4,opt,name=if_octets,json=ifOctets" json:"if_octets,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ForwardingClassAccounting) Reset()         { *m = ForwardingClassAccounting{} }
func (m *ForwardingClassAccounting) String() string { return proto.CompactTextString(m) }
func (*ForwardingClassAccounting) ProtoMessage()    {}
func (*ForwardingClassAccounting) Descriptor() ([]byte, []int) {
	return fileDescriptor_ed53654dcd9b9a05, []int{5}
}
func (m *ForwardingClassAccounting) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ForwardingClassAccounting.Unmarshal(m, b)
}
func (m *ForwardingClassAccounting) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ForwardingClassAccounting.Marshal(b, m, deterministic)
}
func (m *ForwardingClassAccounting) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ForwardingClassAccounting.Merge(m, src)
}
func (m *ForwardingClassAccounting) XXX_Size() int {
	return xxx_messageInfo_ForwardingClassAccounting.Size(m)
}
func (m *ForwardingClassAccounting) XXX_DiscardUnknown() {
	xxx_messageInfo_ForwardingClassAccounting.DiscardUnknown(m)
}

var xxx_messageInfo_ForwardingClassAccounting proto.InternalMessageInfo

func (m *ForwardingClassAccounting) GetIfFamily() string {
	if m != nil && m.IfFamily != nil {
		return *m.IfFamily
	}
	return ""
}

func (m *ForwardingClassAccounting) GetFcNumber() uint32 {
	if m != nil && m.FcNumber != nil {
		return *m.FcNumber
	}
	return 0
}

func (m *ForwardingClassAccounting) GetIfPackets() uint64 {
	if m != nil && m.IfPackets != nil {
		return *m.IfPackets
	}
	return 0
}

func (m *ForwardingClassAccounting) GetIfOctets() uint64 {
	if m != nil && m.IfOctets != nil {
		return *m.IfOctets
	}
	return 0
}

//
//  logical Interface family stats accounting
//
type FamilyAccounting struct {
	// Interface protocol
	IfFamily *string `protobuf:"bytes,1,opt,name=if_family,json=ifFamily" json:"if_family,omitempty"`
	// Count of packets
	IfPackets *uint64 `protobuf:"varint,2,opt,name=if_packets,json=ifPackets" json:"if_packets,omitempty"`
	// Count of v4 bytes
	IfOctets *uint64 `protobuf:"varint,3,opt,name=if_octets,json=ifOctets" json:"if_octets,omitempty"`
	// Count of v6 packets
	IfV6Packets *uint64 `protobuf:"varint,4,opt,name=if_v6_packets,json=ifV6Packets" json:"if_v6_packets,omitempty"`
	// Count of v6 bytes
	IfV6Octets *uint64 `protobuf:"varint,5,opt,name=if_v6_octets,json=ifV6Octets" json:"if_v6_octets,omitempty"`
	// Count of multicast packets
	IfMcastPackets *uint64 `protobuf:"varint,6,opt,name=if_mcast_packets,json=ifMcastPackets" json:"if_mcast_packets,omitempty"`
	// Count of multicast bytes
	IfMcastOctets        *uint64  `protobuf:"varint,7,opt,name=if_mcast_octets,json=ifMcastOctets" json:"if_mcast_octets,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *FamilyAccounting) Reset()         { *m = FamilyAccounting{} }
func (m *FamilyAccounting) String() string { return proto.CompactTextString(m) }
func (*FamilyAccounting) ProtoMessage()    {}
func (*FamilyAccounting) Descriptor() ([]byte, []int) {
	return fileDescriptor_ed53654dcd9b9a05, []int{6}
}
func (m *FamilyAccounting) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FamilyAccounting.Unmarshal(m, b)
}
func (m *FamilyAccounting) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FamilyAccounting.Marshal(b, m, deterministic)
}
func (m *FamilyAccounting) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FamilyAccounting.Merge(m, src)
}
func (m *FamilyAccounting) XXX_Size() int {
	return xxx_messageInfo_FamilyAccounting.Size(m)
}
func (m *FamilyAccounting) XXX_DiscardUnknown() {
	xxx_messageInfo_FamilyAccounting.DiscardUnknown(m)
}

var xxx_messageInfo_FamilyAccounting proto.InternalMessageInfo

func (m *FamilyAccounting) GetIfFamily() string {
	if m != nil && m.IfFamily != nil {
		return *m.IfFamily
	}
	return ""
}

func (m *FamilyAccounting) GetIfPackets() uint64 {
	if m != nil && m.IfPackets != nil {
		return *m.IfPackets
	}
	return 0
}

func (m *FamilyAccounting) GetIfOctets() uint64 {
	if m != nil && m.IfOctets != nil {
		return *m.IfOctets
	}
	return 0
}

func (m *FamilyAccounting) GetIfV6Packets() uint64 {
	if m != nil && m.IfV6Packets != nil {
		return *m.IfV6Packets
	}
	return 0
}

func (m *FamilyAccounting) GetIfV6Octets() uint64 {
	if m != nil && m.IfV6Octets != nil {
		return *m.IfV6Octets
	}
	return 0
}

func (m *FamilyAccounting) GetIfMcastPackets() uint64 {
	if m != nil && m.IfMcastPackets != nil {
		return *m.IfMcastPackets
	}
	return 0
}

func (m *FamilyAccounting) GetIfMcastOctets() uint64 {
	if m != nil && m.IfMcastOctets != nil {
		return *m.IfMcastOctets
	}
	return 0
}

//
// Interface queue statistics
//
type LogicalInterfaceQueueStats struct {
	// Queue number
	QueueNumber *uint32 `protobuf:"varint,1,opt,name=queue_number,json=queueNumber" json:"queue_number,omitempty"`
	// The total number of packets that have been added to this queue
	Packets *uint64 `protobuf:"varint,2,opt,name=packets" json:"packets,omitempty"`
	// The total number of bytes that have been added to this queue
	Bytes *uint64 `protobuf:"varint,3,opt,name=bytes" json:"bytes,omitempty"`
	// The total number of tail dropped packets
	TailDropPackets *uint64 `protobuf:"varint,4,opt,name=tail_drop_packets,json=tailDropPackets" json:"tail_drop_packets,omitempty"`
	// The total number of rate-limited packets
	RateLimitDropPackets *uint64 `protobuf:"varint,5,opt,name=rate_limit_drop_packets,json=rateLimitDropPackets" json:"rate_limit_drop_packets,omitempty"`
	// The total number of rate-limited bytes
	RateLimitDropBytes *uint64 `protobuf:"varint,6,opt,name=rate_limit_drop_bytes,json=rateLimitDropBytes" json:"rate_limit_drop_bytes,omitempty"`
	// The total number of red-dropped packets
	RedDropPackets *uint64 `protobuf:"varint,7,opt,name=red_drop_packets,json=redDropPackets" json:"red_drop_packets,omitempty"`
	// The total number of red-dropped bytes
	RedDropBytes *uint64 `protobuf:"varint,8,opt,name=red_drop_bytes,json=redDropBytes" json:"red_drop_bytes,omitempty"`
	// Average queue depth, in packets
	AverageBufferOccupancy *uint64 `protobuf:"varint,9,opt,name=average_buffer_occupancy,json=averageBufferOccupancy" json:"average_buffer_occupancy,omitempty"`
	// Current queue depth, in packets
	CurrentBufferOccupancy *uint64 `protobuf:"varint,10,opt,name=current_buffer_occupancy,json=currentBufferOccupancy" json:"current_buffer_occupancy,omitempty"`
	// The max measured queue depth, in packets, across all measurements since boot
	PeakBufferOccupancy *uint64 `protobuf:"varint,11,opt,name=peak_buffer_occupancy,json=peakBufferOccupancy" json:"peak_buffer_occupancy,omitempty"`
	// Allocated buffer size
	AllocatedBufferSize  *uint64  `protobuf:"varint,12,opt,name=allocated_buffer_size,json=allocatedBufferSize" json:"allocated_buffer_size,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *LogicalInterfaceQueueStats) Reset()         { *m = LogicalInterfaceQueueStats{} }
func (m *LogicalInterfaceQueueStats) String() string { return proto.CompactTextString(m) }
func (*LogicalInterfaceQueueStats) ProtoMessage()    {}
func (*LogicalInterfaceQueueStats) Descriptor() ([]byte, []int) {
	return fileDescriptor_ed53654dcd9b9a05, []int{7}
}
func (m *LogicalInterfaceQueueStats) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LogicalInterfaceQueueStats.Unmarshal(m, b)
}
func (m *LogicalInterfaceQueueStats) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LogicalInterfaceQueueStats.Marshal(b, m, deterministic)
}
func (m *LogicalInterfaceQueueStats) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LogicalInterfaceQueueStats.Merge(m, src)
}
func (m *LogicalInterfaceQueueStats) XXX_Size() int {
	return xxx_messageInfo_LogicalInterfaceQueueStats.Size(m)
}
func (m *LogicalInterfaceQueueStats) XXX_DiscardUnknown() {
	xxx_messageInfo_LogicalInterfaceQueueStats.DiscardUnknown(m)
}

var xxx_messageInfo_LogicalInterfaceQueueStats proto.InternalMessageInfo

func (m *LogicalInterfaceQueueStats) GetQueueNumber() uint32 {
	if m != nil && m.QueueNumber != nil {
		return *m.QueueNumber
	}
	return 0
}

func (m *LogicalInterfaceQueueStats) GetPackets() uint64 {
	if m != nil && m.Packets != nil {
		return *m.Packets
	}
	return 0
}

func (m *LogicalInterfaceQueueStats) GetBytes() uint64 {
	if m != nil && m.Bytes != nil {
		return *m.Bytes
	}
	return 0
}

func (m *LogicalInterfaceQueueStats) GetTailDropPackets() uint64 {
	if m != nil && m.TailDropPackets != nil {
		return *m.TailDropPackets
	}
	return 0
}

func (m *LogicalInterfaceQueueStats) GetRateLimitDropPackets() uint64 {
	if m != nil && m.RateLimitDropPackets != nil {
		return *m.RateLimitDropPackets
	}
	return 0
}

func (m *LogicalInterfaceQueueStats) GetRateLimitDropBytes() uint64 {
	if m != nil && m.RateLimitDropBytes != nil {
		return *m.RateLimitDropBytes
	}
	return 0
}

func (m *LogicalInterfaceQueueStats) GetRedDropPackets() uint64 {
	if m != nil && m.RedDropPackets != nil {
		return *m.RedDropPackets
	}
	return 0
}

func (m *LogicalInterfaceQueueStats) GetRedDropBytes() uint64 {
	if m != nil && m.RedDropBytes != nil {
		return *m.RedDropBytes
	}
	return 0
}

func (m *LogicalInterfaceQueueStats) GetAverageBufferOccupancy() uint64 {
	if m != nil && m.AverageBufferOccupancy != nil {
		return *m.AverageBufferOccupancy
	}
	return 0
}

func (m *LogicalInterfaceQueueStats) GetCurrentBufferOccupancy() uint64 {
	if m != nil && m.CurrentBufferOccupancy != nil {
		return *m.CurrentBufferOccupancy
	}
	return 0
}

func (m *LogicalInterfaceQueueStats) GetPeakBufferOccupancy() uint64 {
	if m != nil && m.PeakBufferOccupancy != nil {
		return *m.PeakBufferOccupancy
	}
	return 0
}

func (m *LogicalInterfaceQueueStats) GetAllocatedBufferSize() uint64 {
	if m != nil && m.AllocatedBufferSize != nil {
		return *m.AllocatedBufferSize
	}
	return 0
}

var E_JnprLogicalInterfaceExt = &proto.ExtensionDesc{
	ExtendedType:  (*JuniperNetworksSensors)(nil),
	ExtensionType: (*LogicalPort)(nil),
	Field:         7,
	Name:          "jnprLogicalInterfaceExt",
	Tag:           "bytes,7,opt,name=jnprLogicalInterfaceExt",
	Filename:      "logical_port.proto",
}

func init() {
	proto.RegisterType((*LogicalPort)(nil), "LogicalPort")
	proto.RegisterType((*LogicalInterfaceInfo)(nil), "LogicalInterfaceInfo")
	proto.RegisterType((*IngressInterfaceStats)(nil), "IngressInterfaceStats")
	proto.RegisterType((*EgressInterfaceStats)(nil), "EgressInterfaceStats")
	proto.RegisterType((*OperationalState)(nil), "OperationalState")
	proto.RegisterType((*ForwardingClassAccounting)(nil), "ForwardingClassAccounting")
	proto.RegisterType((*FamilyAccounting)(nil), "FamilyAccounting")
	proto.RegisterType((*LogicalInterfaceQueueStats)(nil), "logicalInterfaceQueueStats")
	proto.RegisterExtension(E_JnprLogicalInterfaceExt)
}

func init() { proto.RegisterFile("logical_port.proto", fileDescriptor_ed53654dcd9b9a05) }

var fileDescriptor_ed53654dcd9b9a05 = []byte{
	// 1021 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xd4, 0x56, 0x5f, 0x6f, 0x1b, 0xc5,
	0x17, 0xd5, 0x3a, 0x75, 0xe3, 0x5c, 0xff, 0x89, 0x33, 0x8d, 0x93, 0xfd, 0xa5, 0xfa, 0x51, 0xcb,
	0x42, 0xc2, 0xa8, 0xd4, 0x55, 0x2b, 0x51, 0x95, 0x52, 0x09, 0x92, 0x90, 0x20, 0x43, 0x9b, 0x14,
	0x07, 0x78, 0x1d, 0x4d, 0xd6, 0x77, 0x9c, 0x21, 0xbb, 0x33, 0xcb, 0xec, 0x38, 0x6d, 0xfa, 0x88,
	0xc4, 0xc7, 0xe1, 0x05, 0xf1, 0x25, 0xf8, 0x40, 0xbc, 0xa3, 0x99, 0xdd, 0x75, 0x76, 0x37, 0x0e,
	0xe4, 0x95, 0x37, 0xfb, 0xde, 0x73, 0xce, 0xcc, 0x9c, 0x7b, 0x76, 0x77, 0x80, 0x84, 0x6a, 0x26,
	0x02, 0x16, 0xd2, 0x58, 0x69, 0x33, 0x8a, 0xb5, 0x32, 0x6a, 0xe7, 0x9e, 0xc1, 0x10, 0x23, 0x34,
	0xfa, 0x92, 0x1a, 0x15, 0xa7, 0xc5, 0xc1, 0xb7, 0xd0, 0x7c, 0x95, 0x42, 0xdf, 0x28, 0x6d, 0xc8,
	0x4b, 0xe8, 0x08, 0x69, 0x50, 0x73, 0x16, 0x20, 0x15, 0x92, 0x2b, 0xdf, 0xeb, 0xaf, 0x0c, 0x9b,
	0x4f, 0x7b, 0xa3, 0x0c, 0x35, 0xce, 0xbb, 0x63, 0xc9, 0xd5, 0xa4, 0x2d, 0x8a, 0x7f, 0x07, 0x7f,
	0xd6, 0x61, 0x73, 0x19, 0x8e, 0x7c, 0x00, 0xab, 0x82, 0x53, 0xc9, 0x22, 0xf4, 0xbd, 0x7e, 0x6d,
	0xb8, 0xb6, 0x57, 0xff, 0xe5, 0xcb, 0x5a, 0xc3, 0x9b, 0xdc, 0x15, 0xfc, 0x88, 0x45, 0x48, 0x06,
	0xb0, 0x26, 0xa4, 0x30, 0xd4, 0x88, 0x08, 0xfd, 0x5a, 0xbf, 0x36, 0xbc, 0xe3, 0x10, 0x5d, 0x6f,
	0xd2, 0xb0, 0xf5, 0xef, 0x45, 0x84, 0xe4, 0x63, 0x68, 0x27, 0x32, 0x8a, 0xa9, 0xe0, 0x54, 0xc8,
	0x29, 0xbe, 0xf3, 0x57, 0xfa, 0xde, 0xb0, 0x9d, 0x2b, 0x35, 0x6d, 0x6f, 0xcc, 0xc7, 0xb6, 0x43,
	0x1e, 0x42, 0x27, 0x66, 0x1a, 0xa5, 0xa1, 0x0c, 0xd3, 0x55, 0xef, 0xf4, 0xbd, 0xab, 0x55, 0x5b,
	0x69, 0x73, 0x17, 0xdd, 0xda, 0x9f, 0x43, 0x5b, 0xc8, 0x99, 0xc6, 0x24, 0xa1, 0x89, 0x61, 0x26,
	0xf1, 0xeb, 0x7d, 0x6f, 0xd8, 0x7c, 0xba, 0x35, 0x1a, 0xa7, 0xd5, 0xc5, 0x49, 0x4e, 0x6c, 0x77,
	0xd2, 0xca, 0xc0, 0xee, 0x1f, 0x79, 0x0e, 0x2d, 0x2c, 0x72, 0xef, 0x3a, 0x6e, 0x6f, 0x74, 0xb0,
	0x8c, 0xda, 0xc4, 0x02, 0xf3, 0x13, 0x68, 0xa8, 0xd8, 0xb1, 0xd0, 0x5f, 0x75, 0xac, 0x8d, 0xd1,
	0x71, 0x8c, 0x9a, 0x19, 0xa1, 0x24, 0x0b, 0x2d, 0x08, 0x27, 0xab, 0x2a, 0x76, 0x3f, 0xc8, 0xa7,
	0xb0, 0xc5, 0xa6, 0x91, 0x90, 0x22, 0x31, 0x9a, 0x05, 0x46, 0x5c, 0xa0, 0xa3, 0xce, 0x13, 0xbf,
	0x61, 0x4f, 0x36, 0xe9, 0x55, 0xba, 0x27, 0xae, 0x49, 0xfa, 0xd0, 0x9c, 0x62, 0x12, 0x68, 0x11,
	0x5b, 0x55, 0x7f, 0xcd, 0x61, 0x8b, 0x25, 0xf2, 0x00, 0x9a, 0x21, 0x4b, 0x0c, 0x0d, 0xce, 0x98,
	0x9c, 0xa1, 0x0f, 0xd6, 0xd3, 0x09, 0xd8, 0xd2, 0xbe, 0xab, 0x90, 0xff, 0x03, 0x9c, 0x89, 0xd9,
	0x19, 0x4d, 0x62, 0xc4, 0xa9, 0xdf, 0x74, 0xfd, 0x35, 0x5b, 0x39, 0xb1, 0x05, 0x32, 0x06, 0x92,
	0xbb, 0xf7, 0xf3, 0x1c, 0xe7, 0x59, 0x68, 0x5a, 0x2e, 0x34, 0xf7, 0x47, 0x61, 0x25, 0x0c, 0xdf,
	0x59, 0x48, 0x6a, 0x46, 0x37, 0xa3, 0xb9, 0x92, 0x0b, 0xc9, 0xd7, 0xb0, 0x81, 0xd7, 0x94, 0xda,
	0xff, 0xae, 0xb4, 0x8e, 0x15, 0xa1, 0x7d, 0xd8, 0x61, 0xb3, 0x99, 0xc6, 0x19, 0x33, 0x38, 0xa5,
	0x42, 0x26, 0x86, 0xc9, 0x00, 0x69, 0x84, 0xd1, 0x29, 0x6a, 0xbf, 0xb3, 0x88, 0xc2, 0xd0, 0x9b,
	0xf8, 0x57, 0xc0, 0x71, 0x86, 0x7b, 0xed, 0x60, 0x83, 0xdf, 0x6b, 0xd0, 0x5b, 0x9a, 0x00, 0xf2,
	0x21, 0x80, 0xe0, 0x34, 0x66, 0xc1, 0x39, 0x9a, 0xc4, 0xe5, 0x39, 0x4d, 0xab, 0xef, 0x4d, 0xd6,
	0x04, 0x7f, 0x93, 0xd6, 0x5d, 0xa4, 0x39, 0x55, 0x81, 0xb1, 0xa0, 0x5a, 0x11, 0xd4, 0x10, 0xfc,
	0xd8, 0x95, 0xc9, 0x63, 0xe8, 0x0a, 0x4e, 0xe7, 0x81, 0x1d, 0x40, 0xae, 0x67, 0x53, 0xbd, 0x80,
	0x76, 0x04, 0xff, 0xc1, 0x76, 0x73, 0xd1, 0x94, 0x10, 0x95, 0x08, 0x77, 0x2a, 0x84, 0xd7, 0x45,
	0xc2, 0x0b, 0x68, 0x0a, 0x4e, 0x79, 0xb0, 0x88, 0xb6, 0x75, 0x73, 0x67, 0x74, 0xa8, 0xf4, 0x5b,
	0xa6, 0xa7, 0x42, 0xce, 0xf6, 0x43, 0x96, 0x24, 0xbb, 0x41, 0xa0, 0xe6, 0xd2, 0x08, 0x39, 0xb3,
	0x27, 0x38, 0x0c, 0xd2, 0x73, 0x3e, 0x49, 0xb9, 0x6c, 0x11, 0xed, 0x15, 0x17, 0xd2, 0x43, 0x16,
	0x89, 0xf0, 0xb2, 0x42, 0x61, 0x8e, 0x32, 0xf8, 0xcb, 0x83, 0xcd, 0x83, 0xff, 0xba, 0x67, 0x95,
	0x73, 0xd7, 0x6f, 0x71, 0xee, 0x5d, 0xe8, 0x56, 0x9f, 0x5d, 0xf2, 0x08, 0x88, 0xba, 0xaa, 0xe5,
	0x8f, 0xab, 0xe7, 0x1e, 0xc1, 0x0d, 0x55, 0x46, 0xcf, 0x93, 0xc1, 0x6f, 0x1e, 0xfc, 0xef, 0xc6,
	0xb1, 0x64, 0xce, 0x70, 0xb7, 0x85, 0x54, 0x23, 0x7f, 0x99, 0x35, 0xec, 0x2e, 0x6c, 0xd9, 0x62,
	0x78, 0x40, 0xe5, 0xdc, 0xa5, 0xbc, 0x56, 0x7c, 0x39, 0x36, 0x78, 0x70, 0xe4, 0xca, 0x95, 0x39,
	0x94, 0x7c, 0xbb, 0x69, 0x0e, 0x25, 0xaf, 0x16, 0x73, 0x18, 0xfc, 0x51, 0x83, 0x6e, 0xd5, 0x92,
	0x5b, 0x6d, 0xb3, 0xbc, 0x85, 0xda, 0x6d, 0xb6, 0xb0, 0xb2, 0x74, 0x0b, 0xf6, 0x8b, 0x20, 0x38,
	0xbd, 0x78, 0xb6, 0x7c, 0xac, 0x4d, 0xc1, 0x7f, 0x7c, 0x96, 0xcb, 0x7d, 0x04, 0xad, 0x14, 0x9a,
	0x29, 0xd6, 0x8b, 0x48, 0xb0, 0xc8, 0x52, 0xbc, 0xca, 0x69, 0xb9, 0xfb, 0x4f, 0x69, 0x79, 0x04,
	0xeb, 0x0b, 0x42, 0x26, 0xbe, 0x5a, 0xc4, 0xb7, 0x33, 0x7c, 0x66, 0xdb, 0xaf, 0x75, 0xd8, 0xb9,
	0xf9, 0x5d, 0x46, 0x86, 0xd0, 0x4a, 0x5f, 0x7e, 0xd9, 0x18, 0xbd, 0xd2, 0x37, 0xce, 0xb5, 0xb2,
	0x49, 0x3e, 0x80, 0xd5, 0xa5, 0x1e, 0xe6, 0x55, 0x72, 0x1f, 0xea, 0xa7, 0x97, 0x06, 0x2b, 0xee,
	0xa5, 0x35, 0xf2, 0x04, 0x36, 0x0c, 0x13, 0x21, 0x9d, 0x6a, 0x15, 0x2f, 0xb7, 0x6f, 0xdd, 0xf6,
	0xbf, 0xd2, 0x2a, 0xce, 0x0f, 0xfa, 0x12, 0xb6, 0x35, 0x33, 0x48, 0x43, 0x11, 0x09, 0x53, 0x26,
	0x96, 0xdc, 0xdc, 0xb4, 0xa8, 0x57, 0x16, 0x54, 0x64, 0x3f, 0x87, 0x5e, 0x95, 0x9d, 0xee, 0xae,
	0x64, 0x2e, 0x29, 0x71, 0xf7, 0xdc, 0x56, 0x1f, 0x43, 0x57, 0xe3, 0xb4, 0xbc, 0x60, 0xc9, 0xe1,
	0x8e, 0xc6, 0x69, 0x71, 0xa9, 0x87, 0xd0, 0x59, 0x10, 0xd2, 0x35, 0x1a, 0x45, 0x78, 0x2b, 0x83,
	0xa7, 0xea, 0x5f, 0x80, 0xcf, 0x2e, 0x50, 0xb3, 0x19, 0xd2, 0xd3, 0x39, 0xe7, 0xa8, 0xa9, 0x0a,
	0x82, 0x79, 0xcc, 0x64, 0x70, 0xe9, 0x3e, 0x97, 0x29, 0xad, 0xef, 0x4d, 0xb6, 0x32, 0xd8, 0x9e,
	0x43, 0x1d, 0xe7, 0x20, 0x2b, 0x10, 0xcc, 0xb5, 0xbb, 0x6c, 0x5c, 0x13, 0x80, 0x92, 0x40, 0x06,
	0xab, 0x0a, 0x7c, 0x06, 0xbd, 0x18, 0xd9, 0xf9, 0x75, 0x76, 0xb3, 0xc8, 0xbe, 0x67, 0x31, 0x4b,
	0xa8, 0x2c, 0x0c, 0x55, 0xe0, 0xbe, 0x73, 0x19, 0x3f, 0x11, 0xef, 0xd1, 0x6f, 0x95, 0xa8, 0x0b,
	0x4c, 0xca, 0x3f, 0x11, 0xef, 0xf1, 0x05, 0x85, 0xed, 0x9f, 0x64, 0xac, 0xab, 0xb7, 0xb5, 0x83,
	0x77, 0x86, 0x6c, 0x8f, 0xbe, 0x99, 0x4b, 0x11, 0xa3, 0x3e, 0x42, 0xf3, 0x56, 0xe9, 0xf3, 0xe4,
	0x04, 0x65, 0xa2, 0x74, 0x92, 0x5d, 0x53, 0x5a, 0xa3, 0xc2, 0x85, 0x71, 0x72, 0x93, 0xca, 0xdf,
	0x01, 0x00, 0x00, 0xff, 0xff, 0xf4, 0xf4, 0xc0, 0x6a, 0x82, 0x0a, 0x00, 0x00,
}
