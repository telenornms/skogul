// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: jkdsd_cpu_oc.proto

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

type ComponentsDebug struct {
	Component            []*ComponentsDebugComponentList `protobuf:"bytes,151,rep,name=component" json:"component,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                        `json:"-"`
	XXX_unrecognized     []byte                          `json:"-"`
	XXX_sizecache        int32                           `json:"-"`
}

func (m *ComponentsDebug) Reset()         { *m = ComponentsDebug{} }
func (m *ComponentsDebug) String() string { return proto.CompactTextString(m) }
func (*ComponentsDebug) ProtoMessage()    {}
func (*ComponentsDebug) Descriptor() ([]byte, []int) {
	return fileDescriptor_904bd24d5b27d4fa, []int{0}
}
func (m *ComponentsDebug) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ComponentsDebug.Unmarshal(m, b)
}
func (m *ComponentsDebug) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ComponentsDebug.Marshal(b, m, deterministic)
}
func (m *ComponentsDebug) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ComponentsDebug.Merge(m, src)
}
func (m *ComponentsDebug) XXX_Size() int {
	return xxx_messageInfo_ComponentsDebug.Size(m)
}
func (m *ComponentsDebug) XXX_DiscardUnknown() {
	xxx_messageInfo_ComponentsDebug.DiscardUnknown(m)
}

var xxx_messageInfo_ComponentsDebug proto.InternalMessageInfo

func (m *ComponentsDebug) GetComponent() []*ComponentsDebugComponentList {
	if m != nil {
		return m.Component
	}
	return nil
}

type ComponentsDebugComponentList struct {
	Name                 *string                              `protobuf:"bytes,51,opt,name=name" json:"name,omitempty"`
	Cpu                  *ComponentsDebugComponentListCpuType `protobuf:"bytes,151,opt,name=cpu" json:"cpu,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                             `json:"-"`
	XXX_unrecognized     []byte                               `json:"-"`
	XXX_sizecache        int32                                `json:"-"`
}

func (m *ComponentsDebugComponentList) Reset()         { *m = ComponentsDebugComponentList{} }
func (m *ComponentsDebugComponentList) String() string { return proto.CompactTextString(m) }
func (*ComponentsDebugComponentList) ProtoMessage()    {}
func (*ComponentsDebugComponentList) Descriptor() ([]byte, []int) {
	return fileDescriptor_904bd24d5b27d4fa, []int{0, 0}
}
func (m *ComponentsDebugComponentList) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ComponentsDebugComponentList.Unmarshal(m, b)
}
func (m *ComponentsDebugComponentList) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ComponentsDebugComponentList.Marshal(b, m, deterministic)
}
func (m *ComponentsDebugComponentList) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ComponentsDebugComponentList.Merge(m, src)
}
func (m *ComponentsDebugComponentList) XXX_Size() int {
	return xxx_messageInfo_ComponentsDebugComponentList.Size(m)
}
func (m *ComponentsDebugComponentList) XXX_DiscardUnknown() {
	xxx_messageInfo_ComponentsDebugComponentList.DiscardUnknown(m)
}

var xxx_messageInfo_ComponentsDebugComponentList proto.InternalMessageInfo

func (m *ComponentsDebugComponentList) GetName() string {
	if m != nil && m.Name != nil {
		return *m.Name
	}
	return ""
}

func (m *ComponentsDebugComponentList) GetCpu() *ComponentsDebugComponentListCpuType {
	if m != nil {
		return m.Cpu
	}
	return nil
}

type ComponentsDebugComponentListCpuType struct {
	Utilization          *ComponentsDebugComponentListCpuTypeUtilizationType `protobuf:"bytes,151,opt,name=utilization" json:"utilization,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                                            `json:"-"`
	XXX_unrecognized     []byte                                              `json:"-"`
	XXX_sizecache        int32                                               `json:"-"`
}

func (m *ComponentsDebugComponentListCpuType) Reset()         { *m = ComponentsDebugComponentListCpuType{} }
func (m *ComponentsDebugComponentListCpuType) String() string { return proto.CompactTextString(m) }
func (*ComponentsDebugComponentListCpuType) ProtoMessage()    {}
func (*ComponentsDebugComponentListCpuType) Descriptor() ([]byte, []int) {
	return fileDescriptor_904bd24d5b27d4fa, []int{0, 0, 0}
}
func (m *ComponentsDebugComponentListCpuType) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ComponentsDebugComponentListCpuType.Unmarshal(m, b)
}
func (m *ComponentsDebugComponentListCpuType) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ComponentsDebugComponentListCpuType.Marshal(b, m, deterministic)
}
func (m *ComponentsDebugComponentListCpuType) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ComponentsDebugComponentListCpuType.Merge(m, src)
}
func (m *ComponentsDebugComponentListCpuType) XXX_Size() int {
	return xxx_messageInfo_ComponentsDebugComponentListCpuType.Size(m)
}
func (m *ComponentsDebugComponentListCpuType) XXX_DiscardUnknown() {
	xxx_messageInfo_ComponentsDebugComponentListCpuType.DiscardUnknown(m)
}

var xxx_messageInfo_ComponentsDebugComponentListCpuType proto.InternalMessageInfo

func (m *ComponentsDebugComponentListCpuType) GetUtilization() *ComponentsDebugComponentListCpuTypeUtilizationType {
	if m != nil {
		return m.Utilization
	}
	return nil
}

type ComponentsDebugComponentListCpuTypeUtilizationType struct {
	State                *ComponentsDebugComponentListCpuTypeUtilizationTypeStateType `protobuf:"bytes,151,opt,name=state" json:"state,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                                                     `json:"-"`
	XXX_unrecognized     []byte                                                       `json:"-"`
	XXX_sizecache        int32                                                        `json:"-"`
}

func (m *ComponentsDebugComponentListCpuTypeUtilizationType) Reset() {
	*m = ComponentsDebugComponentListCpuTypeUtilizationType{}
}
func (m *ComponentsDebugComponentListCpuTypeUtilizationType) String() string {
	return proto.CompactTextString(m)
}
func (*ComponentsDebugComponentListCpuTypeUtilizationType) ProtoMessage() {}
func (*ComponentsDebugComponentListCpuTypeUtilizationType) Descriptor() ([]byte, []int) {
	return fileDescriptor_904bd24d5b27d4fa, []int{0, 0, 0, 0}
}
func (m *ComponentsDebugComponentListCpuTypeUtilizationType) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ComponentsDebugComponentListCpuTypeUtilizationType.Unmarshal(m, b)
}
func (m *ComponentsDebugComponentListCpuTypeUtilizationType) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ComponentsDebugComponentListCpuTypeUtilizationType.Marshal(b, m, deterministic)
}
func (m *ComponentsDebugComponentListCpuTypeUtilizationType) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ComponentsDebugComponentListCpuTypeUtilizationType.Merge(m, src)
}
func (m *ComponentsDebugComponentListCpuTypeUtilizationType) XXX_Size() int {
	return xxx_messageInfo_ComponentsDebugComponentListCpuTypeUtilizationType.Size(m)
}
func (m *ComponentsDebugComponentListCpuTypeUtilizationType) XXX_DiscardUnknown() {
	xxx_messageInfo_ComponentsDebugComponentListCpuTypeUtilizationType.DiscardUnknown(m)
}

var xxx_messageInfo_ComponentsDebugComponentListCpuTypeUtilizationType proto.InternalMessageInfo

func (m *ComponentsDebugComponentListCpuTypeUtilizationType) GetState() *ComponentsDebugComponentListCpuTypeUtilizationTypeStateType {
	if m != nil {
		return m.State
	}
	return nil
}

type ComponentsDebugComponentListCpuTypeUtilizationTypeStateType struct {
	Name                 *string  `protobuf:"bytes,51,opt,name=name" json:"name,omitempty"`
	Instant              *uint32  `protobuf:"varint,52,opt,name=instant" json:"instant,omitempty"`
	Avg                  *uint32  `protobuf:"varint,53,opt,name=avg" json:"avg,omitempty"`
	Min                  *uint32  `protobuf:"varint,54,opt,name=min" json:"min,omitempty"`
	Max                  *uint32  `protobuf:"varint,55,opt,name=max" json:"max,omitempty"`
	Interval             *uint64  `protobuf:"varint,56,opt,name=interval" json:"interval,omitempty"`
	MinTime              *uint64  `protobuf:"varint,57,opt,name=min_time,json=minTime" json:"min_time,omitempty"`
	MaxTime              *uint64  `protobuf:"varint,58,opt,name=max_time,json=maxTime" json:"max_time,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ComponentsDebugComponentListCpuTypeUtilizationTypeStateType) Reset() {
	*m = ComponentsDebugComponentListCpuTypeUtilizationTypeStateType{}
}
func (m *ComponentsDebugComponentListCpuTypeUtilizationTypeStateType) String() string {
	return proto.CompactTextString(m)
}
func (*ComponentsDebugComponentListCpuTypeUtilizationTypeStateType) ProtoMessage() {}
func (*ComponentsDebugComponentListCpuTypeUtilizationTypeStateType) Descriptor() ([]byte, []int) {
	return fileDescriptor_904bd24d5b27d4fa, []int{0, 0, 0, 0, 0}
}
func (m *ComponentsDebugComponentListCpuTypeUtilizationTypeStateType) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ComponentsDebugComponentListCpuTypeUtilizationTypeStateType.Unmarshal(m, b)
}
func (m *ComponentsDebugComponentListCpuTypeUtilizationTypeStateType) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ComponentsDebugComponentListCpuTypeUtilizationTypeStateType.Marshal(b, m, deterministic)
}
func (m *ComponentsDebugComponentListCpuTypeUtilizationTypeStateType) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ComponentsDebugComponentListCpuTypeUtilizationTypeStateType.Merge(m, src)
}
func (m *ComponentsDebugComponentListCpuTypeUtilizationTypeStateType) XXX_Size() int {
	return xxx_messageInfo_ComponentsDebugComponentListCpuTypeUtilizationTypeStateType.Size(m)
}
func (m *ComponentsDebugComponentListCpuTypeUtilizationTypeStateType) XXX_DiscardUnknown() {
	xxx_messageInfo_ComponentsDebugComponentListCpuTypeUtilizationTypeStateType.DiscardUnknown(m)
}

var xxx_messageInfo_ComponentsDebugComponentListCpuTypeUtilizationTypeStateType proto.InternalMessageInfo

func (m *ComponentsDebugComponentListCpuTypeUtilizationTypeStateType) GetName() string {
	if m != nil && m.Name != nil {
		return *m.Name
	}
	return ""
}

func (m *ComponentsDebugComponentListCpuTypeUtilizationTypeStateType) GetInstant() uint32 {
	if m != nil && m.Instant != nil {
		return *m.Instant
	}
	return 0
}

func (m *ComponentsDebugComponentListCpuTypeUtilizationTypeStateType) GetAvg() uint32 {
	if m != nil && m.Avg != nil {
		return *m.Avg
	}
	return 0
}

func (m *ComponentsDebugComponentListCpuTypeUtilizationTypeStateType) GetMin() uint32 {
	if m != nil && m.Min != nil {
		return *m.Min
	}
	return 0
}

func (m *ComponentsDebugComponentListCpuTypeUtilizationTypeStateType) GetMax() uint32 {
	if m != nil && m.Max != nil {
		return *m.Max
	}
	return 0
}

func (m *ComponentsDebugComponentListCpuTypeUtilizationTypeStateType) GetInterval() uint64 {
	if m != nil && m.Interval != nil {
		return *m.Interval
	}
	return 0
}

func (m *ComponentsDebugComponentListCpuTypeUtilizationTypeStateType) GetMinTime() uint64 {
	if m != nil && m.MinTime != nil {
		return *m.MinTime
	}
	return 0
}

func (m *ComponentsDebugComponentListCpuTypeUtilizationTypeStateType) GetMaxTime() uint64 {
	if m != nil && m.MaxTime != nil {
		return *m.MaxTime
	}
	return 0
}

var E_JnprComponentsDebugExt = &proto.ExtensionDesc{
	ExtendedType:  (*JuniperNetworksSensors)(nil),
	ExtensionType: (*ComponentsDebug)(nil),
	Field:         107,
	Name:          "jnpr_components_debug_ext",
	Tag:           "bytes,107,opt,name=jnpr_components_debug_ext",
	Filename:      "jkdsd_cpu_oc.proto",
}

func init() {
	proto.RegisterType((*ComponentsDebug)(nil), "components_debug")
	proto.RegisterType((*ComponentsDebugComponentList)(nil), "components_debug.component_list")
	proto.RegisterType((*ComponentsDebugComponentListCpuType)(nil), "components_debug.component_list.cpu_type")
	proto.RegisterType((*ComponentsDebugComponentListCpuTypeUtilizationType)(nil), "components_debug.component_list.cpu_type.utilization_type")
	proto.RegisterType((*ComponentsDebugComponentListCpuTypeUtilizationTypeStateType)(nil), "components_debug.component_list.cpu_type.utilization_type.state_type")
	proto.RegisterExtension(E_JnprComponentsDebugExt)
}

func init() { proto.RegisterFile("jkdsd_cpu_oc.proto", fileDescriptor_904bd24d5b27d4fa) }

var fileDescriptor_904bd24d5b27d4fa = []byte{
	// 394 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x90, 0xbd, 0x8e, 0xd3, 0x40,
	0x14, 0x85, 0xe5, 0x8d, 0x57, 0xeb, 0xbd, 0x11, 0x28, 0x0c, 0x12, 0x4c, 0x5c, 0x59, 0x54, 0xa6,
	0x71, 0xb1, 0xfc, 0x47, 0x48, 0x20, 0x60, 0x1b, 0x0a, 0x0a, 0x43, 0x09, 0x1a, 0x0d, 0xce, 0xd5,
	0x6a, 0x36, 0x9e, 0x1f, 0x79, 0xae, 0x17, 0x2f, 0x25, 0x3d, 0x35, 0x35, 0x0f, 0xc1, 0x0b, 0xf0,
	0x2a, 0xbc, 0x08, 0x1a, 0x3b, 0x4e, 0x42, 0x28, 0x40, 0xdb, 0xcd, 0xfd, 0xce, 0xbd, 0xe7, 0x68,
	0x0e, 0xb0, 0xf3, 0xd5, 0xd2, 0x2f, 0x45, 0xe5, 0x5a, 0x61, 0xab, 0xc2, 0x35, 0x96, 0x6c, 0x7a,
	0x93, 0xb0, 0x46, 0x8d, 0xd4, 0x5c, 0x0a, 0xb2, 0x6e, 0x80, 0x77, 0x7e, 0xc5, 0x30, 0xab, 0xac,
	0x76, 0xd6, 0xa0, 0x21, 0x2f, 0x96, 0xf8, 0xb1, 0x3d, 0x63, 0xcf, 0xe0, 0x78, 0xc3, 0xf8, 0xb7,
	0x28, 0x9b, 0xe4, 0xd3, 0x93, 0xac, 0xd8, 0x5f, 0xdb, 0x02, 0x51, 0x2b, 0x4f, 0xe5, 0xf6, 0x26,
	0xfd, 0x1a, 0xc3, 0xf5, 0x3f, 0x55, 0x36, 0x87, 0xd8, 0x48, 0x8d, 0xfc, 0x5e, 0x16, 0xe5, 0xc7,
	0x2f, 0x0e, 0xbf, 0x3c, 0x3f, 0x48, 0xa2, 0xb2, 0x47, 0xec, 0x29, 0x4c, 0x2a, 0xd7, 0x86, 0xa0,
	0x28, 0x9f, 0x9e, 0xdc, 0xfd, 0x57, 0x50, 0x11, 0x7e, 0x45, 0x97, 0x0e, 0xcb, 0x70, 0x96, 0x7e,
	0x9f, 0x40, 0x32, 0x12, 0xf6, 0x01, 0xa6, 0x2d, 0xa9, 0x5a, 0x7d, 0x96, 0xa4, 0xac, 0x19, 0x2d,
	0x17, 0xff, 0x6d, 0x59, 0xec, 0x5c, 0x0f, 0x19, 0xbb, 0x7e, 0xe9, 0x8f, 0x03, 0x98, 0xed, 0x6f,
	0xb0, 0xf7, 0x70, 0xe8, 0x49, 0x12, 0x8e, 0x69, 0xa7, 0x57, 0x4f, 0x2b, 0x7a, 0xa3, 0x21, 0x78,
	0x30, 0x4d, 0x7f, 0x46, 0x00, 0x5b, 0xca, 0xd8, 0x6e, 0x8d, 0xeb, 0xfe, 0x38, 0x1c, 0x29, 0xe3,
	0x49, 0x1a, 0xe2, 0xf7, 0xb3, 0x28, 0xbf, 0x56, 0x8e, 0x23, 0x9b, 0xc1, 0x44, 0x5e, 0x9c, 0xf1,
	0x07, 0x3d, 0x0d, 0xcf, 0x40, 0xb4, 0x32, 0xfc, 0xe1, 0x40, 0xb4, 0x32, 0x3d, 0x91, 0x1d, 0x7f,
	0xb4, 0x26, 0xb2, 0x63, 0x29, 0x24, 0xca, 0x10, 0x36, 0x17, 0xb2, 0xe6, 0x8f, 0xb3, 0x28, 0x8f,
	0xcb, 0xcd, 0xcc, 0xe6, 0x90, 0x68, 0x65, 0x04, 0x29, 0x8d, 0xfc, 0x49, 0xaf, 0x1d, 0x69, 0x65,
	0xde, 0x29, 0x8d, 0xbd, 0x24, 0xbb, 0x41, 0x5a, 0xac, 0x25, 0xd9, 0x05, 0x69, 0x81, 0x30, 0x3f,
	0x37, 0xae, 0x11, 0xfb, 0xc5, 0x08, 0xec, 0x88, 0xdd, 0x2e, 0x5e, 0xb7, 0x46, 0x39, 0x6c, 0xde,
	0x20, 0x7d, 0xb2, 0xcd, 0xca, 0xbf, 0x45, 0xe3, 0x6d, 0xe3, 0xf9, 0xaa, 0xaf, 0xf3, 0xc6, 0x5f,
	0x75, 0x96, 0xb7, 0x82, 0xd9, 0xcb, 0x0d, 0x7d, 0x15, 0xe0, 0x69, 0x47, 0xbf, 0x03, 0x00, 0x00,
	0xff, 0xff, 0x4e, 0xd3, 0xe0, 0xe0, 0xf6, 0x02, 0x00, 0x00,
}