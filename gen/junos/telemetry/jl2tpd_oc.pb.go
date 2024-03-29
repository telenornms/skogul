// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: jl2tpd_oc.proto

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

type JunosL2Tp struct {
	System               *JunosL2TpSystemType `protobuf:"bytes,151,opt,name=system" json:"system,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *JunosL2Tp) Reset()         { *m = JunosL2Tp{} }
func (m *JunosL2Tp) String() string { return proto.CompactTextString(m) }
func (*JunosL2Tp) ProtoMessage()    {}
func (*JunosL2Tp) Descriptor() ([]byte, []int) {
	return fileDescriptor_fe9259f956296c37, []int{0}
}
func (m *JunosL2Tp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_JunosL2Tp.Unmarshal(m, b)
}
func (m *JunosL2Tp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_JunosL2Tp.Marshal(b, m, deterministic)
}
func (m *JunosL2Tp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_JunosL2Tp.Merge(m, src)
}
func (m *JunosL2Tp) XXX_Size() int {
	return xxx_messageInfo_JunosL2Tp.Size(m)
}
func (m *JunosL2Tp) XXX_DiscardUnknown() {
	xxx_messageInfo_JunosL2Tp.DiscardUnknown(m)
}

var xxx_messageInfo_JunosL2Tp proto.InternalMessageInfo

func (m *JunosL2Tp) GetSystem() *JunosL2TpSystemType {
	if m != nil {
		return m.System
	}
	return nil
}

type JunosL2TpSystemType struct {
	SubscriberManagement *JunosL2TpSystemTypeSubscriberManagementType `protobuf:"bytes,151,opt,name=subscriber_management,json=subscriberManagement" json:"subscriber_management,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                                     `json:"-"`
	XXX_unrecognized     []byte                                       `json:"-"`
	XXX_sizecache        int32                                        `json:"-"`
}

func (m *JunosL2TpSystemType) Reset()         { *m = JunosL2TpSystemType{} }
func (m *JunosL2TpSystemType) String() string { return proto.CompactTextString(m) }
func (*JunosL2TpSystemType) ProtoMessage()    {}
func (*JunosL2TpSystemType) Descriptor() ([]byte, []int) {
	return fileDescriptor_fe9259f956296c37, []int{0, 0}
}
func (m *JunosL2TpSystemType) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_JunosL2TpSystemType.Unmarshal(m, b)
}
func (m *JunosL2TpSystemType) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_JunosL2TpSystemType.Marshal(b, m, deterministic)
}
func (m *JunosL2TpSystemType) XXX_Merge(src proto.Message) {
	xxx_messageInfo_JunosL2TpSystemType.Merge(m, src)
}
func (m *JunosL2TpSystemType) XXX_Size() int {
	return xxx_messageInfo_JunosL2TpSystemType.Size(m)
}
func (m *JunosL2TpSystemType) XXX_DiscardUnknown() {
	xxx_messageInfo_JunosL2TpSystemType.DiscardUnknown(m)
}

var xxx_messageInfo_JunosL2TpSystemType proto.InternalMessageInfo

func (m *JunosL2TpSystemType) GetSubscriberManagement() *JunosL2TpSystemTypeSubscriberManagementType {
	if m != nil {
		return m.SubscriberManagement
	}
	return nil
}

type JunosL2TpSystemTypeSubscriberManagementType struct {
	ClientProtocols      *JunosL2TpSystemTypeSubscriberManagementTypeClientProtocolsType `protobuf:"bytes,151,opt,name=client_protocols,json=clientProtocols" json:"client_protocols,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                                                        `json:"-"`
	XXX_unrecognized     []byte                                                          `json:"-"`
	XXX_sizecache        int32                                                           `json:"-"`
}

func (m *JunosL2TpSystemTypeSubscriberManagementType) Reset() {
	*m = JunosL2TpSystemTypeSubscriberManagementType{}
}
func (m *JunosL2TpSystemTypeSubscriberManagementType) String() string {
	return proto.CompactTextString(m)
}
func (*JunosL2TpSystemTypeSubscriberManagementType) ProtoMessage() {}
func (*JunosL2TpSystemTypeSubscriberManagementType) Descriptor() ([]byte, []int) {
	return fileDescriptor_fe9259f956296c37, []int{0, 0, 0}
}
func (m *JunosL2TpSystemTypeSubscriberManagementType) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_JunosL2TpSystemTypeSubscriberManagementType.Unmarshal(m, b)
}
func (m *JunosL2TpSystemTypeSubscriberManagementType) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_JunosL2TpSystemTypeSubscriberManagementType.Marshal(b, m, deterministic)
}
func (m *JunosL2TpSystemTypeSubscriberManagementType) XXX_Merge(src proto.Message) {
	xxx_messageInfo_JunosL2TpSystemTypeSubscriberManagementType.Merge(m, src)
}
func (m *JunosL2TpSystemTypeSubscriberManagementType) XXX_Size() int {
	return xxx_messageInfo_JunosL2TpSystemTypeSubscriberManagementType.Size(m)
}
func (m *JunosL2TpSystemTypeSubscriberManagementType) XXX_DiscardUnknown() {
	xxx_messageInfo_JunosL2TpSystemTypeSubscriberManagementType.DiscardUnknown(m)
}

var xxx_messageInfo_JunosL2TpSystemTypeSubscriberManagementType proto.InternalMessageInfo

func (m *JunosL2TpSystemTypeSubscriberManagementType) GetClientProtocols() *JunosL2TpSystemTypeSubscriberManagementTypeClientProtocolsType {
	if m != nil {
		return m.ClientProtocols
	}
	return nil
}

type JunosL2TpSystemTypeSubscriberManagementTypeClientProtocolsType struct {
	L2Tp                 *JunosL2TpSystemTypeSubscriberManagementTypeClientProtocolsTypeL2TpType `protobuf:"bytes,151,opt,name=l2tp" json:"l2tp,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                                                                `json:"-"`
	XXX_unrecognized     []byte                                                                  `json:"-"`
	XXX_sizecache        int32                                                                   `json:"-"`
}

func (m *JunosL2TpSystemTypeSubscriberManagementTypeClientProtocolsType) Reset() {
	*m = JunosL2TpSystemTypeSubscriberManagementTypeClientProtocolsType{}
}
func (m *JunosL2TpSystemTypeSubscriberManagementTypeClientProtocolsType) String() string {
	return proto.CompactTextString(m)
}
func (*JunosL2TpSystemTypeSubscriberManagementTypeClientProtocolsType) ProtoMessage() {}
func (*JunosL2TpSystemTypeSubscriberManagementTypeClientProtocolsType) Descriptor() ([]byte, []int) {
	return fileDescriptor_fe9259f956296c37, []int{0, 0, 0, 0}
}
func (m *JunosL2TpSystemTypeSubscriberManagementTypeClientProtocolsType) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_JunosL2TpSystemTypeSubscriberManagementTypeClientProtocolsType.Unmarshal(m, b)
}
func (m *JunosL2TpSystemTypeSubscriberManagementTypeClientProtocolsType) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_JunosL2TpSystemTypeSubscriberManagementTypeClientProtocolsType.Marshal(b, m, deterministic)
}
func (m *JunosL2TpSystemTypeSubscriberManagementTypeClientProtocolsType) XXX_Merge(src proto.Message) {
	xxx_messageInfo_JunosL2TpSystemTypeSubscriberManagementTypeClientProtocolsType.Merge(m, src)
}
func (m *JunosL2TpSystemTypeSubscriberManagementTypeClientProtocolsType) XXX_Size() int {
	return xxx_messageInfo_JunosL2TpSystemTypeSubscriberManagementTypeClientProtocolsType.Size(m)
}
func (m *JunosL2TpSystemTypeSubscriberManagementTypeClientProtocolsType) XXX_DiscardUnknown() {
	xxx_messageInfo_JunosL2TpSystemTypeSubscriberManagementTypeClientProtocolsType.DiscardUnknown(m)
}

var xxx_messageInfo_JunosL2TpSystemTypeSubscriberManagementTypeClientProtocolsType proto.InternalMessageInfo

func (m *JunosL2TpSystemTypeSubscriberManagementTypeClientProtocolsType) GetL2Tp() *JunosL2TpSystemTypeSubscriberManagementTypeClientProtocolsTypeL2TpType {
	if m != nil {
		return m.L2Tp
	}
	return nil
}

type JunosL2TpSystemTypeSubscriberManagementTypeClientProtocolsTypeL2TpType struct {
	Summary              *JunosL2TpSystemTypeSubscriberManagementTypeClientProtocolsTypeL2TpTypeSummaryType `protobuf:"bytes,151,opt,name=summary" json:"summary,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                                                                           `json:"-"`
	XXX_unrecognized     []byte                                                                             `json:"-"`
	XXX_sizecache        int32                                                                              `json:"-"`
}

func (m *JunosL2TpSystemTypeSubscriberManagementTypeClientProtocolsTypeL2TpType) Reset() {
	*m = JunosL2TpSystemTypeSubscriberManagementTypeClientProtocolsTypeL2TpType{}
}
func (m *JunosL2TpSystemTypeSubscriberManagementTypeClientProtocolsTypeL2TpType) String() string {
	return proto.CompactTextString(m)
}
func (*JunosL2TpSystemTypeSubscriberManagementTypeClientProtocolsTypeL2TpType) ProtoMessage() {}
func (*JunosL2TpSystemTypeSubscriberManagementTypeClientProtocolsTypeL2TpType) Descriptor() ([]byte, []int) {
	return fileDescriptor_fe9259f956296c37, []int{0, 0, 0, 0, 0}
}
func (m *JunosL2TpSystemTypeSubscriberManagementTypeClientProtocolsTypeL2TpType) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_JunosL2TpSystemTypeSubscriberManagementTypeClientProtocolsTypeL2TpType.Unmarshal(m, b)
}
func (m *JunosL2TpSystemTypeSubscriberManagementTypeClientProtocolsTypeL2TpType) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_JunosL2TpSystemTypeSubscriberManagementTypeClientProtocolsTypeL2TpType.Marshal(b, m, deterministic)
}
func (m *JunosL2TpSystemTypeSubscriberManagementTypeClientProtocolsTypeL2TpType) XXX_Merge(src proto.Message) {
	xxx_messageInfo_JunosL2TpSystemTypeSubscriberManagementTypeClientProtocolsTypeL2TpType.Merge(m, src)
}
func (m *JunosL2TpSystemTypeSubscriberManagementTypeClientProtocolsTypeL2TpType) XXX_Size() int {
	return xxx_messageInfo_JunosL2TpSystemTypeSubscriberManagementTypeClientProtocolsTypeL2TpType.Size(m)
}
func (m *JunosL2TpSystemTypeSubscriberManagementTypeClientProtocolsTypeL2TpType) XXX_DiscardUnknown() {
	xxx_messageInfo_JunosL2TpSystemTypeSubscriberManagementTypeClientProtocolsTypeL2TpType.DiscardUnknown(m)
}

var xxx_messageInfo_JunosL2TpSystemTypeSubscriberManagementTypeClientProtocolsTypeL2TpType proto.InternalMessageInfo

func (m *JunosL2TpSystemTypeSubscriberManagementTypeClientProtocolsTypeL2TpType) GetSummary() *JunosL2TpSystemTypeSubscriberManagementTypeClientProtocolsTypeL2TpTypeSummaryType {
	if m != nil {
		return m.Summary
	}
	return nil
}

type JunosL2TpSystemTypeSubscriberManagementTypeClientProtocolsTypeL2TpTypeSummaryType struct {
	L2TpStatsTotalTunnels             *uint32  `protobuf:"varint,51,opt,name=l2tp_stats_total_tunnels,json=l2tpStatsTotalTunnels" json:"l2tp_stats_total_tunnels,omitempty"`
	L2TpStatsTotalSessions            *uint32  `protobuf:"varint,52,opt,name=l2tp_stats_total_sessions,json=l2tpStatsTotalSessions" json:"l2tp_stats_total_sessions,omitempty"`
	L2TpStatsControlRxPackets         *uint32  `protobuf:"varint,53,opt,name=l2tp_stats_control_rx_packets,json=l2tpStatsControlRxPackets" json:"l2tp_stats_control_rx_packets,omitempty"`
	L2TpStatsControlRxBytes           *uint32  `protobuf:"varint,54,opt,name=l2tp_stats_control_rx_bytes,json=l2tpStatsControlRxBytes" json:"l2tp_stats_control_rx_bytes,omitempty"`
	L2TpStatsControlTxPackets         *uint32  `protobuf:"varint,55,opt,name=l2tp_stats_control_tx_packets,json=l2tpStatsControlTxPackets" json:"l2tp_stats_control_tx_packets,omitempty"`
	L2TpStatsControlTxBytes           *uint32  `protobuf:"varint,56,opt,name=l2tp_stats_control_tx_bytes,json=l2tpStatsControlTxBytes" json:"l2tp_stats_control_tx_bytes,omitempty"`
	L2TpEraTypeIcrqInflightCount      *uint32  `protobuf:"varint,57,opt,name=l2tp_era_type_icrq_inflight_count,json=l2tpEraTypeIcrqInflightCount" json:"l2tp_era_type_icrq_inflight_count,omitempty"`
	L2TpEraTypeIcrqReportedSuccesses  *uint32  `protobuf:"varint,58,opt,name=l2tp_era_type_icrq_reported_successes,json=l2tpEraTypeIcrqReportedSuccesses" json:"l2tp_era_type_icrq_reported_successes,omitempty"`
	L2TpEraTypeIcrqReportedFailures   *uint32  `protobuf:"varint,59,opt,name=l2tp_era_type_icrq_reported_failures,json=l2tpEraTypeIcrqReportedFailures" json:"l2tp_era_type_icrq_reported_failures,omitempty"`
	L2TpEraTypeSccrqInflightCount     *uint32  `protobuf:"varint,60,opt,name=l2tp_era_type_sccrq_inflight_count,json=l2tpEraTypeSccrqInflightCount" json:"l2tp_era_type_sccrq_inflight_count,omitempty"`
	L2TpEraTypeSccrqReportedSuccesses *uint32  `protobuf:"varint,61,opt,name=l2tp_era_type_sccrq_reported_successes,json=l2tpEraTypeSccrqReportedSuccesses" json:"l2tp_era_type_sccrq_reported_successes,omitempty"`
	L2TpEraTypeSccrqReportedFailures  *uint32  `protobuf:"varint,62,opt,name=l2tp_era_type_sccrq_reported_failures,json=l2tpEraTypeSccrqReportedFailures" json:"l2tp_era_type_sccrq_reported_failures,omitempty"`
	XXX_NoUnkeyedLiteral              struct{} `json:"-"`
	XXX_unrecognized                  []byte   `json:"-"`
	XXX_sizecache                     int32    `json:"-"`
}

func (m *JunosL2TpSystemTypeSubscriberManagementTypeClientProtocolsTypeL2TpTypeSummaryType) Reset() {
	*m = JunosL2TpSystemTypeSubscriberManagementTypeClientProtocolsTypeL2TpTypeSummaryType{}
}
func (m *JunosL2TpSystemTypeSubscriberManagementTypeClientProtocolsTypeL2TpTypeSummaryType) String() string {
	return proto.CompactTextString(m)
}
func (*JunosL2TpSystemTypeSubscriberManagementTypeClientProtocolsTypeL2TpTypeSummaryType) ProtoMessage() {
}
func (*JunosL2TpSystemTypeSubscriberManagementTypeClientProtocolsTypeL2TpTypeSummaryType) Descriptor() ([]byte, []int) {
	return fileDescriptor_fe9259f956296c37, []int{0, 0, 0, 0, 0, 0}
}
func (m *JunosL2TpSystemTypeSubscriberManagementTypeClientProtocolsTypeL2TpTypeSummaryType) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_JunosL2TpSystemTypeSubscriberManagementTypeClientProtocolsTypeL2TpTypeSummaryType.Unmarshal(m, b)
}
func (m *JunosL2TpSystemTypeSubscriberManagementTypeClientProtocolsTypeL2TpTypeSummaryType) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_JunosL2TpSystemTypeSubscriberManagementTypeClientProtocolsTypeL2TpTypeSummaryType.Marshal(b, m, deterministic)
}
func (m *JunosL2TpSystemTypeSubscriberManagementTypeClientProtocolsTypeL2TpTypeSummaryType) XXX_Merge(src proto.Message) {
	xxx_messageInfo_JunosL2TpSystemTypeSubscriberManagementTypeClientProtocolsTypeL2TpTypeSummaryType.Merge(m, src)
}
func (m *JunosL2TpSystemTypeSubscriberManagementTypeClientProtocolsTypeL2TpTypeSummaryType) XXX_Size() int {
	return xxx_messageInfo_JunosL2TpSystemTypeSubscriberManagementTypeClientProtocolsTypeL2TpTypeSummaryType.Size(m)
}
func (m *JunosL2TpSystemTypeSubscriberManagementTypeClientProtocolsTypeL2TpTypeSummaryType) XXX_DiscardUnknown() {
	xxx_messageInfo_JunosL2TpSystemTypeSubscriberManagementTypeClientProtocolsTypeL2TpTypeSummaryType.DiscardUnknown(m)
}

var xxx_messageInfo_JunosL2TpSystemTypeSubscriberManagementTypeClientProtocolsTypeL2TpTypeSummaryType proto.InternalMessageInfo

func (m *JunosL2TpSystemTypeSubscriberManagementTypeClientProtocolsTypeL2TpTypeSummaryType) GetL2TpStatsTotalTunnels() uint32 {
	if m != nil && m.L2TpStatsTotalTunnels != nil {
		return *m.L2TpStatsTotalTunnels
	}
	return 0
}

func (m *JunosL2TpSystemTypeSubscriberManagementTypeClientProtocolsTypeL2TpTypeSummaryType) GetL2TpStatsTotalSessions() uint32 {
	if m != nil && m.L2TpStatsTotalSessions != nil {
		return *m.L2TpStatsTotalSessions
	}
	return 0
}

func (m *JunosL2TpSystemTypeSubscriberManagementTypeClientProtocolsTypeL2TpTypeSummaryType) GetL2TpStatsControlRxPackets() uint32 {
	if m != nil && m.L2TpStatsControlRxPackets != nil {
		return *m.L2TpStatsControlRxPackets
	}
	return 0
}

func (m *JunosL2TpSystemTypeSubscriberManagementTypeClientProtocolsTypeL2TpTypeSummaryType) GetL2TpStatsControlRxBytes() uint32 {
	if m != nil && m.L2TpStatsControlRxBytes != nil {
		return *m.L2TpStatsControlRxBytes
	}
	return 0
}

func (m *JunosL2TpSystemTypeSubscriberManagementTypeClientProtocolsTypeL2TpTypeSummaryType) GetL2TpStatsControlTxPackets() uint32 {
	if m != nil && m.L2TpStatsControlTxPackets != nil {
		return *m.L2TpStatsControlTxPackets
	}
	return 0
}

func (m *JunosL2TpSystemTypeSubscriberManagementTypeClientProtocolsTypeL2TpTypeSummaryType) GetL2TpStatsControlTxBytes() uint32 {
	if m != nil && m.L2TpStatsControlTxBytes != nil {
		return *m.L2TpStatsControlTxBytes
	}
	return 0
}

func (m *JunosL2TpSystemTypeSubscriberManagementTypeClientProtocolsTypeL2TpTypeSummaryType) GetL2TpEraTypeIcrqInflightCount() uint32 {
	if m != nil && m.L2TpEraTypeIcrqInflightCount != nil {
		return *m.L2TpEraTypeIcrqInflightCount
	}
	return 0
}

func (m *JunosL2TpSystemTypeSubscriberManagementTypeClientProtocolsTypeL2TpTypeSummaryType) GetL2TpEraTypeIcrqReportedSuccesses() uint32 {
	if m != nil && m.L2TpEraTypeIcrqReportedSuccesses != nil {
		return *m.L2TpEraTypeIcrqReportedSuccesses
	}
	return 0
}

func (m *JunosL2TpSystemTypeSubscriberManagementTypeClientProtocolsTypeL2TpTypeSummaryType) GetL2TpEraTypeIcrqReportedFailures() uint32 {
	if m != nil && m.L2TpEraTypeIcrqReportedFailures != nil {
		return *m.L2TpEraTypeIcrqReportedFailures
	}
	return 0
}

func (m *JunosL2TpSystemTypeSubscriberManagementTypeClientProtocolsTypeL2TpTypeSummaryType) GetL2TpEraTypeSccrqInflightCount() uint32 {
	if m != nil && m.L2TpEraTypeSccrqInflightCount != nil {
		return *m.L2TpEraTypeSccrqInflightCount
	}
	return 0
}

func (m *JunosL2TpSystemTypeSubscriberManagementTypeClientProtocolsTypeL2TpTypeSummaryType) GetL2TpEraTypeSccrqReportedSuccesses() uint32 {
	if m != nil && m.L2TpEraTypeSccrqReportedSuccesses != nil {
		return *m.L2TpEraTypeSccrqReportedSuccesses
	}
	return 0
}

func (m *JunosL2TpSystemTypeSubscriberManagementTypeClientProtocolsTypeL2TpTypeSummaryType) GetL2TpEraTypeSccrqReportedFailures() uint32 {
	if m != nil && m.L2TpEraTypeSccrqReportedFailures != nil {
		return *m.L2TpEraTypeSccrqReportedFailures
	}
	return 0
}

var E_JnprJunosL2TpExt = &proto.ExtensionDesc{
	ExtendedType:  (*JuniperNetworksSensors)(nil),
	ExtensionType: (*JunosL2Tp)(nil),
	Field:         44,
	Name:          "jnpr_junos_l2tp_ext",
	Tag:           "bytes,44,opt,name=jnpr_junos_l2tp_ext",
	Filename:      "jl2tpd_oc.proto",
}

func init() {
	proto.RegisterType((*JunosL2Tp)(nil), "junos_l2tp")
	proto.RegisterType((*JunosL2TpSystemType)(nil), "junos_l2tp.system_type")
	proto.RegisterType((*JunosL2TpSystemTypeSubscriberManagementType)(nil), "junos_l2tp.system_type.subscriber_management_type")
	proto.RegisterType((*JunosL2TpSystemTypeSubscriberManagementTypeClientProtocolsType)(nil), "junos_l2tp.system_type.subscriber_management_type.client_protocols_type")
	proto.RegisterType((*JunosL2TpSystemTypeSubscriberManagementTypeClientProtocolsTypeL2TpType)(nil), "junos_l2tp.system_type.subscriber_management_type.client_protocols_type.l2tp_type")
	proto.RegisterType((*JunosL2TpSystemTypeSubscriberManagementTypeClientProtocolsTypeL2TpTypeSummaryType)(nil), "junos_l2tp.system_type.subscriber_management_type.client_protocols_type.l2tp_type.summary_type")
	proto.RegisterExtension(E_JnprJunosL2TpExt)
}

func init() { proto.RegisterFile("jl2tpd_oc.proto", fileDescriptor_fe9259f956296c37) }

var fileDescriptor_fe9259f956296c37 = []byte{
	// 586 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x93, 0x5d, 0x4f, 0xd4, 0x4e,
	0x14, 0xc6, 0x43, 0xf2, 0x0f, 0xfc, 0x39, 0xab, 0x81, 0x0c, 0x22, 0xa5, 0x4a, 0x04, 0xa2, 0x86,
	0x0b, 0xd3, 0x18, 0x7c, 0x41, 0x10, 0x8d, 0x91, 0xf8, 0x02, 0x11, 0xc5, 0xdd, 0xbd, 0x9f, 0x94,
	0x61, 0x80, 0x42, 0x77, 0xa6, 0xcc, 0x39, 0x8d, 0x5b, 0xbf, 0x88, 0x1f, 0xc5, 0xef, 0xe2, 0xb5,
	0x89, 0xdf, 0xc0, 0x2b, 0x2f, 0xcc, 0x4c, 0xdb, 0x7d, 0x6b, 0x57, 0x63, 0xc2, 0xe5, 0x9e, 0xe7,
	0x79, 0x7e, 0xcf, 0x39, 0xdd, 0x0c, 0xcc, 0x9c, 0xc5, 0xeb, 0x94, 0x1c, 0x71, 0x2d, 0x82, 0xc4,
	0x68, 0xd2, 0xfe, 0x1c, 0xc9, 0x58, 0x76, 0x24, 0x99, 0x8c, 0x93, 0x4e, 0xf2, 0xe1, 0xea, 0xd7,
	0x06, 0xc0, 0x59, 0xaa, 0x34, 0x72, 0xeb, 0x66, 0xf7, 0x61, 0x12, 0x33, 0x24, 0xd9, 0xf1, 0xbe,
	0x4c, 0x2c, 0x4f, 0xac, 0x35, 0xd6, 0x17, 0x82, 0xbe, 0x1a, 0xe4, 0x12, 0xa7, 0x2c, 0x91, 0xcd,
	0xc2, 0xe7, 0xff, 0x04, 0x68, 0x0c, 0xcc, 0xd9, 0x29, 0xcc, 0x63, 0x7a, 0x88, 0xc2, 0x44, 0x87,
	0xd2, 0xf0, 0x4e, 0xa8, 0xc2, 0x13, 0xd9, 0x91, 0x8a, 0x4a, 0xe0, 0xfa, 0x18, 0x60, 0x50, 0x9b,
	0xca, 0xbb, 0xae, 0xf5, 0xb5, 0xfd, 0x9e, 0xe4, 0xff, 0x9a, 0x06, 0x7f, 0x7c, 0x88, 0x11, 0xcc,
	0x8a, 0x38, 0xb2, 0x3f, 0xdd, 0xa5, 0x42, 0xc7, 0x58, 0xee, 0xf0, 0xf6, 0xdf, 0x77, 0x08, 0x46,
	0x59, 0xf9, 0x66, 0x33, 0xf9, 0xf8, 0xa0, 0x9c, 0xfa, 0x3f, 0xfe, 0x87, 0xf9, 0x5a, 0x2b, 0x3b,
	0x81, 0xff, 0x6c, 0x5f, 0xb9, 0x43, 0xf3, 0xb2, 0x76, 0x08, 0x2c, 0x21, 0xdf, 0xc6, 0x15, 0xf8,
	0xdf, 0xa7, 0x60, 0xba, 0x37, 0x63, 0x9f, 0x61, 0x0a, 0xd3, 0x4e, 0x27, 0x34, 0x59, 0xd9, 0xcc,
	0x2f, 0xbf, 0x39, 0x28, 0x2a, 0xf2, 0x35, 0xca, 0x42, 0xff, 0xdb, 0x24, 0x5c, 0x19, 0x54, 0xd8,
	0x06, 0x78, 0x2e, 0x83, 0x14, 0x12, 0x72, 0xd2, 0x14, 0xc6, 0x9c, 0x52, 0xa5, 0x64, 0x8c, 0xde,
	0x83, 0xe5, 0x89, 0xb5, 0xab, 0xcd, 0x79, 0xab, 0xb7, 0xac, 0xdc, 0xb6, 0x6a, 0x3b, 0x17, 0xd9,
	0x26, 0x2c, 0x56, 0x82, 0x28, 0x11, 0x23, 0xad, 0xd0, 0x7b, 0xe8, 0x92, 0xd7, 0x87, 0x93, 0xad,
	0x42, 0x65, 0x2f, 0x60, 0x69, 0x20, 0x2a, 0xb4, 0x22, 0xa3, 0x63, 0x6e, 0xba, 0x3c, 0x09, 0xc5,
	0xb9, 0x24, 0xf4, 0x1e, 0xb9, 0xf8, 0x62, 0x2f, 0xbe, 0x93, 0x5b, 0x9a, 0xdd, 0x83, 0xdc, 0xc0,
	0xb6, 0xe1, 0x46, 0x3d, 0xe1, 0x30, 0x23, 0x89, 0xde, 0x63, 0x97, 0x5f, 0xa8, 0xe6, 0x5f, 0x5a,
	0x79, 0x4c, 0x3f, 0xf5, 0xfb, 0x37, 0xea, 0xfb, 0xdb, 0x7f, 0xe9, 0xa7, 0xb2, 0xff, 0x49, 0x7d,
	0x7f, 0xbb, 0xe8, 0x7f, 0x03, 0x2b, 0x2e, 0x2d, 0x4d, 0xe8, 0xfe, 0x04, 0x1e, 0x09, 0x73, 0xc1,
	0x23, 0x75, 0x1c, 0x47, 0x27, 0xa7, 0xc4, 0x85, 0x4e, 0x15, 0x79, 0x9b, 0x8e, 0x71, 0xd3, 0x1a,
	0x5f, 0x99, 0xb0, 0x9d, 0x25, 0x72, 0x57, 0x98, 0x8b, 0xdd, 0xc2, 0xb4, 0x63, 0x3d, 0xec, 0x03,
	0xdc, 0xa9, 0x01, 0x19, 0x99, 0x68, 0x43, 0xf2, 0x88, 0x63, 0x2a, 0x84, 0x44, 0x94, 0xe8, 0x6d,
	0x39, 0xd8, 0xf2, 0x08, 0xac, 0x59, 0x18, 0x5b, 0xa5, 0x8f, 0xed, 0xc3, 0xed, 0x3f, 0x01, 0x8f,
	0xc3, 0x28, 0x4e, 0x8d, 0x44, 0xef, 0xa9, 0xe3, 0xdd, 0x1a, 0xc3, 0x7b, 0x5d, 0xd8, 0xd8, 0x2e,
	0xac, 0x0e, 0xe3, 0x50, 0xd4, 0x5c, 0xba, 0xed, 0x60, 0x4b, 0x03, 0xb0, 0x96, 0xa8, 0x9c, 0xfa,
	0x11, 0xee, 0xd6, 0xa1, 0x6a, 0x6e, 0x7d, 0xe6, 0x70, 0x2b, 0xa3, 0xb8, 0xea, 0xb1, 0x95, 0xaf,
	0x37, 0x82, 0xec, 0x5d, 0xfb, 0xbc, 0xf2, 0xf5, 0x86, 0x88, 0xe5, 0xb9, 0x5b, 0x2d, 0x98, 0x3b,
	0x53, 0x89, 0xe1, 0xfd, 0xc7, 0xcc, 0x65, 0x97, 0xd8, 0x42, 0xb0, 0x97, 0xaa, 0x28, 0x91, 0xe6,
	0xbd, 0xa4, 0x4f, 0xda, 0x9c, 0x63, 0x4b, 0x2a, 0xd4, 0x06, 0xbd, 0x7b, 0xee, 0xf1, 0x37, 0x06,
	0x1e, 0x7f, 0x73, 0xd6, 0x02, 0xf6, 0xec, 0xef, 0x77, 0xb6, 0xad, 0x4b, 0xbf, 0x03, 0x00, 0x00,
	0xff, 0xff, 0x77, 0x34, 0xfa, 0x2a, 0x35, 0x06, 0x00, 0x00,
}
