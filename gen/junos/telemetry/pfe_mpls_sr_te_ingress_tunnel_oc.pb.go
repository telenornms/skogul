// Code generated by protoc-gen-go. DO NOT EDIT.
// source: pfe_mpls_sr_te_ingress_tunnel_oc.proto

package telemetry

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

type MplsPfeMplsSrTeIngressTunnel struct {
	SignalingProtocols   *MplsPfeMplsSrTeIngressTunnelSignalingProtocolsType `protobuf:"bytes,151,opt,name=signaling_protocols,json=signalingProtocols" json:"signaling_protocols,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                                            `json:"-"`
	XXX_unrecognized     []byte                                              `json:"-"`
	XXX_sizecache        int32                                               `json:"-"`
}

func (m *MplsPfeMplsSrTeIngressTunnel) Reset()         { *m = MplsPfeMplsSrTeIngressTunnel{} }
func (m *MplsPfeMplsSrTeIngressTunnel) String() string { return proto.CompactTextString(m) }
func (*MplsPfeMplsSrTeIngressTunnel) ProtoMessage()    {}
func (*MplsPfeMplsSrTeIngressTunnel) Descriptor() ([]byte, []int) {
	return fileDescriptor_a4be3eadfa4d9513, []int{0}
}

func (m *MplsPfeMplsSrTeIngressTunnel) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MplsPfeMplsSrTeIngressTunnel.Unmarshal(m, b)
}
func (m *MplsPfeMplsSrTeIngressTunnel) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MplsPfeMplsSrTeIngressTunnel.Marshal(b, m, deterministic)
}
func (m *MplsPfeMplsSrTeIngressTunnel) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MplsPfeMplsSrTeIngressTunnel.Merge(m, src)
}
func (m *MplsPfeMplsSrTeIngressTunnel) XXX_Size() int {
	return xxx_messageInfo_MplsPfeMplsSrTeIngressTunnel.Size(m)
}
func (m *MplsPfeMplsSrTeIngressTunnel) XXX_DiscardUnknown() {
	xxx_messageInfo_MplsPfeMplsSrTeIngressTunnel.DiscardUnknown(m)
}

var xxx_messageInfo_MplsPfeMplsSrTeIngressTunnel proto.InternalMessageInfo

func (m *MplsPfeMplsSrTeIngressTunnel) GetSignalingProtocols() *MplsPfeMplsSrTeIngressTunnelSignalingProtocolsType {
	if m != nil {
		return m.SignalingProtocols
	}
	return nil
}

type MplsPfeMplsSrTeIngressTunnelSignalingProtocolsType struct {
	SegmentRouting       *MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingType `protobuf:"bytes,151,opt,name=segment_routing,json=segmentRouting" json:"segment_routing,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                                                              `json:"-"`
	XXX_unrecognized     []byte                                                                `json:"-"`
	XXX_sizecache        int32                                                                 `json:"-"`
}

func (m *MplsPfeMplsSrTeIngressTunnelSignalingProtocolsType) Reset() {
	*m = MplsPfeMplsSrTeIngressTunnelSignalingProtocolsType{}
}
func (m *MplsPfeMplsSrTeIngressTunnelSignalingProtocolsType) String() string {
	return proto.CompactTextString(m)
}
func (*MplsPfeMplsSrTeIngressTunnelSignalingProtocolsType) ProtoMessage() {}
func (*MplsPfeMplsSrTeIngressTunnelSignalingProtocolsType) Descriptor() ([]byte, []int) {
	return fileDescriptor_a4be3eadfa4d9513, []int{0, 0}
}

func (m *MplsPfeMplsSrTeIngressTunnelSignalingProtocolsType) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MplsPfeMplsSrTeIngressTunnelSignalingProtocolsType.Unmarshal(m, b)
}
func (m *MplsPfeMplsSrTeIngressTunnelSignalingProtocolsType) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MplsPfeMplsSrTeIngressTunnelSignalingProtocolsType.Marshal(b, m, deterministic)
}
func (m *MplsPfeMplsSrTeIngressTunnelSignalingProtocolsType) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MplsPfeMplsSrTeIngressTunnelSignalingProtocolsType.Merge(m, src)
}
func (m *MplsPfeMplsSrTeIngressTunnelSignalingProtocolsType) XXX_Size() int {
	return xxx_messageInfo_MplsPfeMplsSrTeIngressTunnelSignalingProtocolsType.Size(m)
}
func (m *MplsPfeMplsSrTeIngressTunnelSignalingProtocolsType) XXX_DiscardUnknown() {
	xxx_messageInfo_MplsPfeMplsSrTeIngressTunnelSignalingProtocolsType.DiscardUnknown(m)
}

var xxx_messageInfo_MplsPfeMplsSrTeIngressTunnelSignalingProtocolsType proto.InternalMessageInfo

func (m *MplsPfeMplsSrTeIngressTunnelSignalingProtocolsType) GetSegmentRouting() *MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingType {
	if m != nil {
		return m.SegmentRouting
	}
	return nil
}

type MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingType struct {
	SrTeIngressTunnelPolicies *MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingTypeSrTeIngressTunnelPoliciesType `protobuf:"bytes,151,opt,name=sr_te_ingress_tunnel_policies,json=srTeIngressTunnelPolicies" json:"sr_te_ingress_tunnel_policies,omitempty"`
	XXX_NoUnkeyedLiteral      struct{}                                                                                           `json:"-"`
	XXX_unrecognized          []byte                                                                                             `json:"-"`
	XXX_sizecache             int32                                                                                              `json:"-"`
}

func (m *MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingType) Reset() {
	*m = MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingType{}
}
func (m *MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingType) String() string {
	return proto.CompactTextString(m)
}
func (*MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingType) ProtoMessage() {}
func (*MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingType) Descriptor() ([]byte, []int) {
	return fileDescriptor_a4be3eadfa4d9513, []int{0, 0, 0}
}

func (m *MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingType) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingType.Unmarshal(m, b)
}
func (m *MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingType) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingType.Marshal(b, m, deterministic)
}
func (m *MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingType) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingType.Merge(m, src)
}
func (m *MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingType) XXX_Size() int {
	return xxx_messageInfo_MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingType.Size(m)
}
func (m *MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingType) XXX_DiscardUnknown() {
	xxx_messageInfo_MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingType.DiscardUnknown(m)
}

var xxx_messageInfo_MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingType proto.InternalMessageInfo

func (m *MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingType) GetSrTeIngressTunnelPolicies() *MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingTypeSrTeIngressTunnelPoliciesType {
	if m != nil {
		return m.SrTeIngressTunnelPolicies
	}
	return nil
}

type MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingTypeSrTeIngressTunnelPoliciesType struct {
	SrTeIngressTunnelPolicy []*MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingTypeSrTeIngressTunnelPoliciesTypeSrTeIngressTunnelPolicyList `protobuf:"bytes,151,rep,name=sr_te_ingress_tunnel_policy,json=srTeIngressTunnelPolicy" json:"sr_te_ingress_tunnel_policy,omitempty"`
	XXX_NoUnkeyedLiteral    struct{}                                                                                                                        `json:"-"`
	XXX_unrecognized        []byte                                                                                                                          `json:"-"`
	XXX_sizecache           int32                                                                                                                           `json:"-"`
}

func (m *MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingTypeSrTeIngressTunnelPoliciesType) Reset() {
	*m = MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingTypeSrTeIngressTunnelPoliciesType{}
}
func (m *MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingTypeSrTeIngressTunnelPoliciesType) String() string {
	return proto.CompactTextString(m)
}
func (*MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingTypeSrTeIngressTunnelPoliciesType) ProtoMessage() {
}
func (*MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingTypeSrTeIngressTunnelPoliciesType) Descriptor() ([]byte, []int) {
	return fileDescriptor_a4be3eadfa4d9513, []int{0, 0, 0, 0}
}

func (m *MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingTypeSrTeIngressTunnelPoliciesType) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingTypeSrTeIngressTunnelPoliciesType.Unmarshal(m, b)
}
func (m *MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingTypeSrTeIngressTunnelPoliciesType) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingTypeSrTeIngressTunnelPoliciesType.Marshal(b, m, deterministic)
}
func (m *MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingTypeSrTeIngressTunnelPoliciesType) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingTypeSrTeIngressTunnelPoliciesType.Merge(m, src)
}
func (m *MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingTypeSrTeIngressTunnelPoliciesType) XXX_Size() int {
	return xxx_messageInfo_MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingTypeSrTeIngressTunnelPoliciesType.Size(m)
}
func (m *MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingTypeSrTeIngressTunnelPoliciesType) XXX_DiscardUnknown() {
	xxx_messageInfo_MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingTypeSrTeIngressTunnelPoliciesType.DiscardUnknown(m)
}

var xxx_messageInfo_MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingTypeSrTeIngressTunnelPoliciesType proto.InternalMessageInfo

func (m *MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingTypeSrTeIngressTunnelPoliciesType) GetSrTeIngressTunnelPolicy() []*MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingTypeSrTeIngressTunnelPoliciesTypeSrTeIngressTunnelPolicyList {
	if m != nil {
		return m.SrTeIngressTunnelPolicy
	}
	return nil
}

type MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingTypeSrTeIngressTunnelPoliciesTypeSrTeIngressTunnelPolicyList struct {
	State                *MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingTypeSrTeIngressTunnelPoliciesTypeSrTeIngressTunnelPolicyListStateType `protobuf:"bytes,151,opt,name=state" json:"state,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                                                                                                                               `json:"-"`
	XXX_unrecognized     []byte                                                                                                                                 `json:"-"`
	XXX_sizecache        int32                                                                                                                                  `json:"-"`
}

func (m *MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingTypeSrTeIngressTunnelPoliciesTypeSrTeIngressTunnelPolicyList) Reset() {
	*m = MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingTypeSrTeIngressTunnelPoliciesTypeSrTeIngressTunnelPolicyList{}
}
func (m *MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingTypeSrTeIngressTunnelPoliciesTypeSrTeIngressTunnelPolicyList) String() string {
	return proto.CompactTextString(m)
}
func (*MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingTypeSrTeIngressTunnelPoliciesTypeSrTeIngressTunnelPolicyList) ProtoMessage() {
}
func (*MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingTypeSrTeIngressTunnelPoliciesTypeSrTeIngressTunnelPolicyList) Descriptor() ([]byte, []int) {
	return fileDescriptor_a4be3eadfa4d9513, []int{0, 0, 0, 0, 0}
}

func (m *MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingTypeSrTeIngressTunnelPoliciesTypeSrTeIngressTunnelPolicyList) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingTypeSrTeIngressTunnelPoliciesTypeSrTeIngressTunnelPolicyList.Unmarshal(m, b)
}
func (m *MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingTypeSrTeIngressTunnelPoliciesTypeSrTeIngressTunnelPolicyList) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingTypeSrTeIngressTunnelPoliciesTypeSrTeIngressTunnelPolicyList.Marshal(b, m, deterministic)
}
func (m *MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingTypeSrTeIngressTunnelPoliciesTypeSrTeIngressTunnelPolicyList) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingTypeSrTeIngressTunnelPoliciesTypeSrTeIngressTunnelPolicyList.Merge(m, src)
}
func (m *MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingTypeSrTeIngressTunnelPoliciesTypeSrTeIngressTunnelPolicyList) XXX_Size() int {
	return xxx_messageInfo_MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingTypeSrTeIngressTunnelPoliciesTypeSrTeIngressTunnelPolicyList.Size(m)
}
func (m *MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingTypeSrTeIngressTunnelPoliciesTypeSrTeIngressTunnelPolicyList) XXX_DiscardUnknown() {
	xxx_messageInfo_MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingTypeSrTeIngressTunnelPoliciesTypeSrTeIngressTunnelPolicyList.DiscardUnknown(m)
}

var xxx_messageInfo_MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingTypeSrTeIngressTunnelPoliciesTypeSrTeIngressTunnelPolicyList proto.InternalMessageInfo

func (m *MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingTypeSrTeIngressTunnelPoliciesTypeSrTeIngressTunnelPolicyList) GetState() *MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingTypeSrTeIngressTunnelPoliciesTypeSrTeIngressTunnelPolicyListStateType {
	if m != nil {
		return m.State
	}
	return nil
}

type MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingTypeSrTeIngressTunnelPoliciesTypeSrTeIngressTunnelPolicyListStateType struct {
	Counters             []*MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingTypeSrTeIngressTunnelPoliciesTypeSrTeIngressTunnelPolicyListStateTypeCountersList `protobuf:"bytes,151,rep,name=counters" json:"counters,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                                                                                                                                             `json:"-"`
	XXX_unrecognized     []byte                                                                                                                                               `json:"-"`
	XXX_sizecache        int32                                                                                                                                                `json:"-"`
}

func (m *MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingTypeSrTeIngressTunnelPoliciesTypeSrTeIngressTunnelPolicyListStateType) Reset() {
	*m = MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingTypeSrTeIngressTunnelPoliciesTypeSrTeIngressTunnelPolicyListStateType{}
}
func (m *MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingTypeSrTeIngressTunnelPoliciesTypeSrTeIngressTunnelPolicyListStateType) String() string {
	return proto.CompactTextString(m)
}
func (*MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingTypeSrTeIngressTunnelPoliciesTypeSrTeIngressTunnelPolicyListStateType) ProtoMessage() {
}
func (*MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingTypeSrTeIngressTunnelPoliciesTypeSrTeIngressTunnelPolicyListStateType) Descriptor() ([]byte, []int) {
	return fileDescriptor_a4be3eadfa4d9513, []int{0, 0, 0, 0, 0, 0}
}

func (m *MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingTypeSrTeIngressTunnelPoliciesTypeSrTeIngressTunnelPolicyListStateType) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingTypeSrTeIngressTunnelPoliciesTypeSrTeIngressTunnelPolicyListStateType.Unmarshal(m, b)
}
func (m *MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingTypeSrTeIngressTunnelPoliciesTypeSrTeIngressTunnelPolicyListStateType) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingTypeSrTeIngressTunnelPoliciesTypeSrTeIngressTunnelPolicyListStateType.Marshal(b, m, deterministic)
}
func (m *MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingTypeSrTeIngressTunnelPoliciesTypeSrTeIngressTunnelPolicyListStateType) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingTypeSrTeIngressTunnelPoliciesTypeSrTeIngressTunnelPolicyListStateType.Merge(m, src)
}
func (m *MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingTypeSrTeIngressTunnelPoliciesTypeSrTeIngressTunnelPolicyListStateType) XXX_Size() int {
	return xxx_messageInfo_MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingTypeSrTeIngressTunnelPoliciesTypeSrTeIngressTunnelPolicyListStateType.Size(m)
}
func (m *MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingTypeSrTeIngressTunnelPoliciesTypeSrTeIngressTunnelPolicyListStateType) XXX_DiscardUnknown() {
	xxx_messageInfo_MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingTypeSrTeIngressTunnelPoliciesTypeSrTeIngressTunnelPolicyListStateType.DiscardUnknown(m)
}

var xxx_messageInfo_MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingTypeSrTeIngressTunnelPoliciesTypeSrTeIngressTunnelPolicyListStateType proto.InternalMessageInfo

func (m *MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingTypeSrTeIngressTunnelPoliciesTypeSrTeIngressTunnelPolicyListStateType) GetCounters() []*MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingTypeSrTeIngressTunnelPoliciesTypeSrTeIngressTunnelPolicyListStateTypeCountersList {
	if m != nil {
		return m.Counters
	}
	return nil
}

type MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingTypeSrTeIngressTunnelPoliciesTypeSrTeIngressTunnelPolicyListStateTypeCountersList struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingTypeSrTeIngressTunnelPoliciesTypeSrTeIngressTunnelPolicyListStateTypeCountersList) Reset() {
	*m = MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingTypeSrTeIngressTunnelPoliciesTypeSrTeIngressTunnelPolicyListStateTypeCountersList{}
}
func (m *MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingTypeSrTeIngressTunnelPoliciesTypeSrTeIngressTunnelPolicyListStateTypeCountersList) String() string {
	return proto.CompactTextString(m)
}
func (*MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingTypeSrTeIngressTunnelPoliciesTypeSrTeIngressTunnelPolicyListStateTypeCountersList) ProtoMessage() {
}
func (*MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingTypeSrTeIngressTunnelPoliciesTypeSrTeIngressTunnelPolicyListStateTypeCountersList) Descriptor() ([]byte, []int) {
	return fileDescriptor_a4be3eadfa4d9513, []int{0, 0, 0, 0, 0, 0, 0}
}

func (m *MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingTypeSrTeIngressTunnelPoliciesTypeSrTeIngressTunnelPolicyListStateTypeCountersList) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingTypeSrTeIngressTunnelPoliciesTypeSrTeIngressTunnelPolicyListStateTypeCountersList.Unmarshal(m, b)
}
func (m *MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingTypeSrTeIngressTunnelPoliciesTypeSrTeIngressTunnelPolicyListStateTypeCountersList) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingTypeSrTeIngressTunnelPoliciesTypeSrTeIngressTunnelPolicyListStateTypeCountersList.Marshal(b, m, deterministic)
}
func (m *MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingTypeSrTeIngressTunnelPoliciesTypeSrTeIngressTunnelPolicyListStateTypeCountersList) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingTypeSrTeIngressTunnelPoliciesTypeSrTeIngressTunnelPolicyListStateTypeCountersList.Merge(m, src)
}
func (m *MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingTypeSrTeIngressTunnelPoliciesTypeSrTeIngressTunnelPolicyListStateTypeCountersList) XXX_Size() int {
	return xxx_messageInfo_MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingTypeSrTeIngressTunnelPoliciesTypeSrTeIngressTunnelPolicyListStateTypeCountersList.Size(m)
}
func (m *MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingTypeSrTeIngressTunnelPoliciesTypeSrTeIngressTunnelPolicyListStateTypeCountersList) XXX_DiscardUnknown() {
	xxx_messageInfo_MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingTypeSrTeIngressTunnelPoliciesTypeSrTeIngressTunnelPolicyListStateTypeCountersList.DiscardUnknown(m)
}

var xxx_messageInfo_MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingTypeSrTeIngressTunnelPoliciesTypeSrTeIngressTunnelPolicyListStateTypeCountersList proto.InternalMessageInfo

var E_JnprMplsPfeMplsSrTeIngressTunnelExt = &proto.ExtensionDesc{
	ExtendedType:  (*JuniperNetworksSensors)(nil),
	ExtensionType: (*MplsPfeMplsSrTeIngressTunnel)(nil),
	Field:         84,
	Name:          "jnpr_mpls_pfe_mpls_sr_te_ingress_tunnel_ext",
	Tag:           "bytes,84,opt,name=jnpr_mpls_pfe_mpls_sr_te_ingress_tunnel_ext",
	Filename:      "pfe_mpls_sr_te_ingress_tunnel_oc.proto",
}

func init() {
	proto.RegisterType((*MplsPfeMplsSrTeIngressTunnel)(nil), "mpls_pfe_mpls_sr_te_ingress_tunnel")
	proto.RegisterType((*MplsPfeMplsSrTeIngressTunnelSignalingProtocolsType)(nil), "mpls_pfe_mpls_sr_te_ingress_tunnel.signaling_protocols_type")
	proto.RegisterType((*MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingType)(nil), "mpls_pfe_mpls_sr_te_ingress_tunnel.signaling_protocols_type.segment_routing_type")
	proto.RegisterType((*MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingTypeSrTeIngressTunnelPoliciesType)(nil), "mpls_pfe_mpls_sr_te_ingress_tunnel.signaling_protocols_type.segment_routing_type.sr_te_ingress_tunnel_policies_type")
	proto.RegisterType((*MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingTypeSrTeIngressTunnelPoliciesTypeSrTeIngressTunnelPolicyList)(nil), "mpls_pfe_mpls_sr_te_ingress_tunnel.signaling_protocols_type.segment_routing_type.sr_te_ingress_tunnel_policies_type.sr_te_ingress_tunnel_policy_list")
	proto.RegisterType((*MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingTypeSrTeIngressTunnelPoliciesTypeSrTeIngressTunnelPolicyListStateType)(nil), "mpls_pfe_mpls_sr_te_ingress_tunnel.signaling_protocols_type.segment_routing_type.sr_te_ingress_tunnel_policies_type.sr_te_ingress_tunnel_policy_list.state_type")
	proto.RegisterType((*MplsPfeMplsSrTeIngressTunnelSignalingProtocolsTypeSegmentRoutingTypeSrTeIngressTunnelPoliciesTypeSrTeIngressTunnelPolicyListStateTypeCountersList)(nil), "mpls_pfe_mpls_sr_te_ingress_tunnel.signaling_protocols_type.segment_routing_type.sr_te_ingress_tunnel_policies_type.sr_te_ingress_tunnel_policy_list.state_type.counters_list")
	proto.RegisterExtension(E_JnprMplsPfeMplsSrTeIngressTunnelExt)
}

func init() {
	proto.RegisterFile("pfe_mpls_sr_te_ingress_tunnel_oc.proto", fileDescriptor_a4be3eadfa4d9513)
}

var fileDescriptor_a4be3eadfa4d9513 = []byte{
	// 396 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xcc, 0x52, 0xc1, 0xaa, 0x13, 0x31,
	0x14, 0x25, 0x3c, 0x1e, 0x3e, 0xef, 0x43, 0x1f, 0xe4, 0x09, 0xad, 0x23, 0x42, 0x69, 0x41, 0x0a,
	0x42, 0x16, 0x5d, 0x8a, 0x5b, 0x17, 0x0a, 0x4a, 0x9d, 0x76, 0x1f, 0xca, 0x70, 0x3b, 0x44, 0xd3,
	0x24, 0x24, 0x77, 0xb0, 0x83, 0xdf, 0xa1, 0xee, 0xc4, 0x2f, 0x70, 0x6b, 0xff, 0xc1, 0x82, 0x2b,
	0xff, 0xc2, 0x8f, 0x90, 0x4e, 0x3a, 0x15, 0xcb, 0xb4, 0x23, 0x28, 0xd8, 0xd5, 0xc0, 0xbd, 0xe7,
	0x9c, 0x9c, 0x73, 0xcf, 0xc0, 0x03, 0x37, 0x47, 0xb9, 0x70, 0x3a, 0xc8, 0xe0, 0x25, 0xa1, 0x54,
	0x26, 0xf7, 0x18, 0x82, 0xa4, 0xc2, 0x18, 0xd4, 0xd2, 0x66, 0xc2, 0x79, 0x4b, 0x36, 0xb9, 0x26,
	0xd4, 0xb8, 0x40, 0xf2, 0xa5, 0x24, 0xeb, 0xe2, 0xb0, 0xff, 0xed, 0x26, 0xf4, 0x2b, 0xee, 0x51,
	0x11, 0x6e, 0xe0, 0x3a, 0xa8, 0xdc, 0xcc, 0xb4, 0x32, 0xb9, 0xac, 0x98, 0x99, 0xd5, 0xa1, 0xfb,
	0x81, 0xf5, 0xd8, 0xf0, 0x72, 0xf4, 0x58, 0xb4, 0x4b, 0x88, 0x06, 0xbe, 0xa4, 0xd2, 0x61, 0xca,
	0x77, 0x9b, 0x71, 0xbd, 0x48, 0x56, 0x17, 0xd0, 0x3d, 0x44, 0xe0, 0x6f, 0xe1, 0x2a, 0x60, 0xbe,
	0x40, 0x43, 0xd2, 0xdb, 0x82, 0x94, 0xc9, 0x6b, 0x23, 0x2f, 0xff, 0xc6, 0x88, 0xd8, 0x13, 0x8d,
	0xee, 0x6e, 0x6f, 0xa7, 0x69, 0x1c, 0x26, 0x1f, 0x6f, 0xc0, 0x9d, 0x26, 0x20, 0xff, 0xc2, 0xe0,
	0x7e, 0x63, 0x01, 0xce, 0x6a, 0x95, 0x29, 0xdc, 0x5d, 0x8b, 0xfe, 0xb9, 0x49, 0x71, 0xf4, 0xdd,
	0x98, 0xe3, 0x6e, 0xf0, 0x53, 0x7c, 0x1a, 0x11, 0xd3, 0x0a, 0x30, 0xde, 0xee, 0x93, 0xaf, 0xe7,
	0xd0, 0x6f, 0x57, 0xe0, 0xdf, 0x19, 0xdc, 0x3b, 0x0c, 0x2b, 0x37, 0xf1, 0xce, 0x86, 0x97, 0xa3,
	0x77, 0xec, 0x7f, 0xe4, 0x3b, 0x02, 0x29, 0xa5, 0x56, 0x81, 0xd2, 0x4e, 0xf3, 0x01, 0xca, 0x64,
	0x7d, 0x06, 0xbd, 0x36, 0x36, 0x5f, 0x31, 0x38, 0x0f, 0x34, 0x23, 0xac, 0x5b, 0xfc, 0x74, 0x9a,
	0x31, 0x45, 0x65, 0x32, 0x56, 0x1e, 0x0d, 0x27, 0x3f, 0x18, 0xc0, 0xaf, 0x29, 0x5f, 0x33, 0xb8,
	0xc8, 0x6c, 0x61, 0x08, 0x7d, 0xa8, 0x3b, 0xfb, 0x7c, 0xf2, 0x61, 0x44, 0x6d, 0x39, 0x96, 0xb9,
	0x4b, 0x90, 0x5c, 0xc1, 0xad, 0xdf, 0x56, 0x8f, 0xde, 0x33, 0x78, 0xf8, 0xca, 0x38, 0x2f, 0xdb,
	0x13, 0x49, 0x5c, 0x12, 0xef, 0x88, 0x67, 0x85, 0x51, 0x0e, 0xfd, 0x0b, 0xa4, 0x37, 0xd6, 0xbf,
	0x0e, 0x13, 0x34, 0xc1, 0xfa, 0xd0, 0x9d, 0x56, 0x3d, 0x0f, 0xfe, 0xe0, 0x32, 0xe9, 0x60, 0xf3,
	0xe0, 0x73, 0xa7, 0xc3, 0x78, 0x8e, 0x9b, 0xcf, 0x64, 0xff, 0x57, 0x7b, 0xb2, 0xa4, 0x9f, 0x01,
	0x00, 0x00, 0xff, 0xff, 0x50, 0xd5, 0x6e, 0x09, 0xa8, 0x05, 0x00, 0x00,
}