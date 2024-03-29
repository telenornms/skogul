// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: cmerror.proto

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
// Juniper Error Item information
//
type Error struct {
	// Identifier that uniquely identifies the source of
	// the error.
	// e.g.
	//
	// junos/system/linecard/0/pcie/0/lane/0/pcie_cmerror_uncorrectable_major
	//
	Identifier *string `protobuf:"bytes,1,opt,name=identifier" json:"identifier,omitempty"`
	// Name of the error
	Name *string `protobuf:"bytes,2,opt,name=name" json:"name,omitempty"`
	// Instance id of the associated resource
	ComponentId *uint32 `protobuf:"varint,3,opt,name=component_id,json=componentId" json:"component_id,omitempty"`
	// Fru information
	FruType *string `protobuf:"bytes,4,opt,name=fru_type,json=fruType" json:"fru_type,omitempty"`
	FruSlot *uint32 `protobuf:"varint,5,opt,name=fru_slot,json=fruSlot" json:"fru_slot,omitempty"`
	// Scope,Category,Severity
	// in which this error belong to.
	Scope    *string `protobuf:"bytes,6,opt,name=scope" json:"scope,omitempty"`
	Category *string `protobuf:"bytes,7,opt,name=category" json:"category,omitempty"`
	Severity *string `protobuf:"bytes,8,opt,name=severity" json:"severity,omitempty"`
	// Thresholds and action configured for this
	// error.
	Threshold         *uint32 `protobuf:"varint,9,opt,name=threshold" json:"threshold,omitempty"`
	Limit             *uint32 `protobuf:"varint,10,opt,name=limit" json:"limit,omitempty"`
	RaisingThreshold  *uint32 `protobuf:"varint,11,opt,name=raising_threshold,json=raisingThreshold" json:"raising_threshold,omitempty"`
	ClearingThreshold *uint32 `protobuf:"varint,12,opt,name=clearing_threshold,json=clearingThreshold" json:"clearing_threshold,omitempty"`
	Action            *uint32 `protobuf:"varint,13,opt,name=action" json:"action,omitempty"`
	// local/global/both
	ActionHandlingType *uint32 `protobuf:"varint,14,opt,name=action_handling_type,json=actionHandlingType" json:"action_handling_type,omitempty"`
	// user configured thresholds and limits for this error.
	ConfiguredThreshold   *uint32  `protobuf:"varint,15,opt,name=configured_threshold,json=configuredThreshold" json:"configured_threshold,omitempty"`
	ConfiguredLimit       *uint32  `protobuf:"varint,16,opt,name=configured_limit,json=configuredLimit" json:"configured_limit,omitempty"`
	ConfiguredAction      *uint32  `protobuf:"varint,17,opt,name=configured_action,json=configuredAction" json:"configured_action,omitempty"`
	ConfiguredClearAction *uint32  `protobuf:"varint,18,opt,name=configured_clear_action,json=configuredClearAction" json:"configured_clear_action,omitempty"`
	XXX_NoUnkeyedLiteral  struct{} `json:"-"`
	XXX_unrecognized      []byte   `json:"-"`
	XXX_sizecache         int32    `json:"-"`
}

func (m *Error) Reset()         { *m = Error{} }
func (m *Error) String() string { return proto.CompactTextString(m) }
func (*Error) ProtoMessage()    {}
func (*Error) Descriptor() ([]byte, []int) {
	return fileDescriptor_747f0735808ade43, []int{0}
}
func (m *Error) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Error.Unmarshal(m, b)
}
func (m *Error) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Error.Marshal(b, m, deterministic)
}
func (m *Error) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Error.Merge(m, src)
}
func (m *Error) XXX_Size() int {
	return xxx_messageInfo_Error.Size(m)
}
func (m *Error) XXX_DiscardUnknown() {
	xxx_messageInfo_Error.DiscardUnknown(m)
}

var xxx_messageInfo_Error proto.InternalMessageInfo

func (m *Error) GetIdentifier() string {
	if m != nil && m.Identifier != nil {
		return *m.Identifier
	}
	return ""
}

func (m *Error) GetName() string {
	if m != nil && m.Name != nil {
		return *m.Name
	}
	return ""
}

func (m *Error) GetComponentId() uint32 {
	if m != nil && m.ComponentId != nil {
		return *m.ComponentId
	}
	return 0
}

func (m *Error) GetFruType() string {
	if m != nil && m.FruType != nil {
		return *m.FruType
	}
	return ""
}

func (m *Error) GetFruSlot() uint32 {
	if m != nil && m.FruSlot != nil {
		return *m.FruSlot
	}
	return 0
}

func (m *Error) GetScope() string {
	if m != nil && m.Scope != nil {
		return *m.Scope
	}
	return ""
}

func (m *Error) GetCategory() string {
	if m != nil && m.Category != nil {
		return *m.Category
	}
	return ""
}

func (m *Error) GetSeverity() string {
	if m != nil && m.Severity != nil {
		return *m.Severity
	}
	return ""
}

func (m *Error) GetThreshold() uint32 {
	if m != nil && m.Threshold != nil {
		return *m.Threshold
	}
	return 0
}

func (m *Error) GetLimit() uint32 {
	if m != nil && m.Limit != nil {
		return *m.Limit
	}
	return 0
}

func (m *Error) GetRaisingThreshold() uint32 {
	if m != nil && m.RaisingThreshold != nil {
		return *m.RaisingThreshold
	}
	return 0
}

func (m *Error) GetClearingThreshold() uint32 {
	if m != nil && m.ClearingThreshold != nil {
		return *m.ClearingThreshold
	}
	return 0
}

func (m *Error) GetAction() uint32 {
	if m != nil && m.Action != nil {
		return *m.Action
	}
	return 0
}

func (m *Error) GetActionHandlingType() uint32 {
	if m != nil && m.ActionHandlingType != nil {
		return *m.ActionHandlingType
	}
	return 0
}

func (m *Error) GetConfiguredThreshold() uint32 {
	if m != nil && m.ConfiguredThreshold != nil {
		return *m.ConfiguredThreshold
	}
	return 0
}

func (m *Error) GetConfiguredLimit() uint32 {
	if m != nil && m.ConfiguredLimit != nil {
		return *m.ConfiguredLimit
	}
	return 0
}

func (m *Error) GetConfiguredAction() uint32 {
	if m != nil && m.ConfiguredAction != nil {
		return *m.ConfiguredAction
	}
	return 0
}

func (m *Error) GetConfiguredClearAction() uint32 {
	if m != nil && m.ConfiguredClearAction != nil {
		return *m.ConfiguredClearAction
	}
	return 0
}

type GlobalErrorConfiguration struct {
	// configuration bucket identifier
	Scope    *string `protobuf:"bytes,1,opt,name=scope" json:"scope,omitempty"`
	Category *string `protobuf:"bytes,2,opt,name=category" json:"category,omitempty"`
	Severity *string `protobuf:"bytes,3,opt,name=severity" json:"severity,omitempty"`
	// configured parameters for this bucket.
	Threshold            *uint32  `protobuf:"varint,4,opt,name=threshold" json:"threshold,omitempty"`
	Action               *uint32  `protobuf:"varint,5,opt,name=action" json:"action,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GlobalErrorConfiguration) Reset()         { *m = GlobalErrorConfiguration{} }
func (m *GlobalErrorConfiguration) String() string { return proto.CompactTextString(m) }
func (*GlobalErrorConfiguration) ProtoMessage()    {}
func (*GlobalErrorConfiguration) Descriptor() ([]byte, []int) {
	return fileDescriptor_747f0735808ade43, []int{1}
}
func (m *GlobalErrorConfiguration) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GlobalErrorConfiguration.Unmarshal(m, b)
}
func (m *GlobalErrorConfiguration) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GlobalErrorConfiguration.Marshal(b, m, deterministic)
}
func (m *GlobalErrorConfiguration) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GlobalErrorConfiguration.Merge(m, src)
}
func (m *GlobalErrorConfiguration) XXX_Size() int {
	return xxx_messageInfo_GlobalErrorConfiguration.Size(m)
}
func (m *GlobalErrorConfiguration) XXX_DiscardUnknown() {
	xxx_messageInfo_GlobalErrorConfiguration.DiscardUnknown(m)
}

var xxx_messageInfo_GlobalErrorConfiguration proto.InternalMessageInfo

func (m *GlobalErrorConfiguration) GetScope() string {
	if m != nil && m.Scope != nil {
		return *m.Scope
	}
	return ""
}

func (m *GlobalErrorConfiguration) GetCategory() string {
	if m != nil && m.Category != nil {
		return *m.Category
	}
	return ""
}

func (m *GlobalErrorConfiguration) GetSeverity() string {
	if m != nil && m.Severity != nil {
		return *m.Severity
	}
	return ""
}

func (m *GlobalErrorConfiguration) GetThreshold() uint32 {
	if m != nil && m.Threshold != nil {
		return *m.Threshold
	}
	return 0
}

func (m *GlobalErrorConfiguration) GetAction() uint32 {
	if m != nil && m.Action != nil {
		return *m.Action
	}
	return 0
}

//
// Top-level Cmerror message
//
type Cmerror struct {
	// collection of error items
	ErrorItem []*Error `protobuf:"bytes,1,rep,name=error_item,json=errorItem" json:"error_item,omitempty"`
	// last configuration change for cmerror.
	LastConfigurationChange *uint64 `protobuf:"varint,2,opt,name=last_configuration_change,json=lastConfigurationChange" json:"last_configuration_change,omitempty"`
	// This will toggle at start of every wrap cycle
	WrapCycleMarker *bool `protobuf:"varint,3,opt,name=wrap_cycle_marker,json=wrapCycleMarker" json:"wrap_cycle_marker,omitempty"`
	// Fru slot identifier
	FruSlot *uint32 `protobuf:"varint,4,opt,name=fru_slot,json=fruSlot" json:"fru_slot,omitempty"`
	FruType *string `protobuf:"bytes,5,opt,name=fru_type,json=fruType" json:"fru_type,omitempty"`
	// Collection of global configuration items
	GlobalConfigurationItem []*GlobalErrorConfiguration `protobuf:"bytes,6,rep,name=global_configuration_item,json=globalConfigurationItem" json:"global_configuration_item,omitempty"`
	XXX_NoUnkeyedLiteral    struct{}                    `json:"-"`
	XXX_unrecognized        []byte                      `json:"-"`
	XXX_sizecache           int32                       `json:"-"`
}

func (m *Cmerror) Reset()         { *m = Cmerror{} }
func (m *Cmerror) String() string { return proto.CompactTextString(m) }
func (*Cmerror) ProtoMessage()    {}
func (*Cmerror) Descriptor() ([]byte, []int) {
	return fileDescriptor_747f0735808ade43, []int{2}
}
func (m *Cmerror) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Cmerror.Unmarshal(m, b)
}
func (m *Cmerror) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Cmerror.Marshal(b, m, deterministic)
}
func (m *Cmerror) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Cmerror.Merge(m, src)
}
func (m *Cmerror) XXX_Size() int {
	return xxx_messageInfo_Cmerror.Size(m)
}
func (m *Cmerror) XXX_DiscardUnknown() {
	xxx_messageInfo_Cmerror.DiscardUnknown(m)
}

var xxx_messageInfo_Cmerror proto.InternalMessageInfo

func (m *Cmerror) GetErrorItem() []*Error {
	if m != nil {
		return m.ErrorItem
	}
	return nil
}

func (m *Cmerror) GetLastConfigurationChange() uint64 {
	if m != nil && m.LastConfigurationChange != nil {
		return *m.LastConfigurationChange
	}
	return 0
}

func (m *Cmerror) GetWrapCycleMarker() bool {
	if m != nil && m.WrapCycleMarker != nil {
		return *m.WrapCycleMarker
	}
	return false
}

func (m *Cmerror) GetFruSlot() uint32 {
	if m != nil && m.FruSlot != nil {
		return *m.FruSlot
	}
	return 0
}

func (m *Cmerror) GetFruType() string {
	if m != nil && m.FruType != nil {
		return *m.FruType
	}
	return ""
}

func (m *Cmerror) GetGlobalConfigurationItem() []*GlobalErrorConfiguration {
	if m != nil {
		return m.GlobalConfigurationItem
	}
	return nil
}

var E_JnprCmerrorExt = &proto.ExtensionDesc{
	ExtendedType:  (*JuniperNetworksSensors)(nil),
	ExtensionType: (*Cmerror)(nil),
	Field:         20,
	Name:          "jnpr_cmerror_ext",
	Tag:           "bytes,20,opt,name=jnpr_cmerror_ext",
	Filename:      "cmerror.proto",
}

func init() {
	proto.RegisterType((*Error)(nil), "Error")
	proto.RegisterType((*GlobalErrorConfiguration)(nil), "GlobalErrorConfiguration")
	proto.RegisterType((*Cmerror)(nil), "Cmerror")
	proto.RegisterExtension(E_JnprCmerrorExt)
}

func init() { proto.RegisterFile("cmerror.proto", fileDescriptor_747f0735808ade43) }

var fileDescriptor_747f0735808ade43 = []byte{
	// 593 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x94, 0x4d, 0x6f, 0xd3, 0x4c,
	0x10, 0xc7, 0xe5, 0x34, 0x4e, 0x93, 0x49, 0x5f, 0x92, 0x6d, 0x9e, 0x27, 0x9b, 0x8a, 0x43, 0xa8,
	0x84, 0x14, 0x40, 0x44, 0xd0, 0x03, 0x87, 0x9e, 0x80, 0xa8, 0x82, 0xf2, 0x76, 0x70, 0xcb, 0xd9,
	0x32, 0xce, 0x24, 0x59, 0x6a, 0xef, 0x5a, 0xeb, 0x0d, 0x6d, 0xae, 0x1c, 0xf8, 0x18, 0xdc, 0xf8,
	0x9e, 0x68, 0x67, 0x9d, 0xd8, 0xa9, 0xd4, 0xde, 0x3c, 0xf3, 0xfb, 0xff, 0xbd, 0x3b, 0xe3, 0x19,
	0xc3, 0x7e, 0x9c, 0xa2, 0xd6, 0x4a, 0x8f, 0x33, 0xad, 0x8c, 0x3a, 0x3e, 0x32, 0x98, 0x60, 0x8a,
	0x46, 0xaf, 0x42, 0xa3, 0x32, 0x97, 0x3c, 0xf9, 0xed, 0x83, 0x7f, 0x6e, 0x45, 0xec, 0x09, 0x80,
	0x98, 0xa2, 0x34, 0x62, 0x26, 0x50, 0x73, 0x6f, 0xe8, 0x8d, 0x5a, 0xef, 0xfc, 0x5f, 0x6f, 0x6a,
	0x4d, 0x2f, 0xa8, 0x00, 0xc6, 0xa0, 0x2e, 0xa3, 0x14, 0x79, 0xcd, 0x0a, 0x02, 0x7a, 0x66, 0x8f,
	0x61, 0x2f, 0x56, 0x69, 0xa6, 0x24, 0x4a, 0x13, 0x8a, 0x29, 0xdf, 0x19, 0x7a, 0xa3, 0xfd, 0xa0,
	0xbd, 0xc9, 0x5d, 0x4c, 0xd9, 0x00, 0x9a, 0x33, 0xbd, 0x0c, 0xcd, 0x2a, 0x43, 0x5e, 0x27, 0xeb,
	0xee, 0x4c, 0x2f, 0xaf, 0x56, 0x19, 0xae, 0x51, 0x9e, 0x28, 0xc3, 0x7d, 0x72, 0x5a, 0x74, 0x99,
	0x28, 0xc3, 0x7a, 0xe0, 0xe7, 0xb1, 0xca, 0x90, 0x37, 0xc8, 0xe2, 0x02, 0x76, 0x0c, 0xcd, 0x38,
	0x32, 0x38, 0x57, 0x7a, 0xc5, 0x77, 0x09, 0x6c, 0x62, 0xcb, 0x72, 0xfc, 0x89, 0x5a, 0x98, 0x15,
	0x6f, 0x3a, 0xb6, 0x8e, 0xd9, 0x23, 0x68, 0x99, 0x85, 0xc6, 0x7c, 0xa1, 0x92, 0x29, 0x6f, 0xd1,
	0x49, 0x65, 0xc2, 0x9e, 0x95, 0x88, 0x54, 0x18, 0x0e, 0x44, 0x5c, 0xc0, 0x9e, 0x43, 0x57, 0x47,
	0x22, 0x17, 0x72, 0x1e, 0x96, 0xde, 0x36, 0x29, 0x3a, 0x05, 0xb8, 0xda, 0xbc, 0xe2, 0x05, 0xb0,
	0x38, 0xc1, 0x48, 0x6f, 0xab, 0xf7, 0x48, 0xdd, 0x5d, 0x93, 0x52, 0xfe, 0x3f, 0x34, 0xa2, 0xd8,
	0x08, 0x25, 0xf9, 0x3e, 0x49, 0x8a, 0x88, 0xbd, 0x84, 0x9e, 0x7b, 0x0a, 0x17, 0x91, 0x9c, 0x26,
	0xf4, 0x36, 0xdb, 0xb7, 0x03, 0x52, 0x31, 0xc7, 0x3e, 0x14, 0x88, 0x5a, 0xf8, 0x0a, 0x7a, 0xb1,
	0x92, 0x33, 0x31, 0x5f, 0x6a, 0x9c, 0x56, 0x8e, 0x3e, 0x24, 0xc7, 0x51, 0xc9, 0xca, 0xc3, 0x9f,
	0x42, 0xa7, 0x62, 0x71, 0x95, 0x77, 0x48, 0x7e, 0x58, 0xe6, 0x3f, 0xaf, 0x7b, 0x50, 0x91, 0x16,
	0x57, 0xee, 0xba, 0x1e, 0x94, 0xe0, 0xad, 0xbb, 0xfc, 0x6b, 0xe8, 0x57, 0xc4, 0x54, 0xf4, 0xda,
	0xc2, 0xc8, 0xf2, 0x5f, 0x89, 0x27, 0x96, 0x3a, 0xdf, 0xc9, 0x1f, 0x0f, 0xf8, 0xfb, 0x44, 0x7d,
	0x8f, 0x12, 0x1a, 0xc7, 0x49, 0x21, 0x8a, 0xe8, 0xa5, 0x9b, 0x39, 0xf0, 0xee, 0x9b, 0x83, 0xda,
	0x03, 0x73, 0xb0, 0xf3, 0xd0, 0x1c, 0xd4, 0xef, 0xce, 0x41, 0xf9, 0x55, 0xfc, 0xea, 0x57, 0x39,
	0xf9, 0x5b, 0x83, 0xdd, 0x89, 0x5b, 0x28, 0xbb, 0x2b, 0xf4, 0x10, 0x0a, 0x83, 0x29, 0xf7, 0x86,
	0x3b, 0xa3, 0xf6, 0x69, 0x63, 0x4c, 0x17, 0x0f, 0x5a, 0x44, 0x2e, 0x0c, 0xa6, 0xec, 0x0c, 0x06,
	0x49, 0x94, 0x9b, 0x30, 0xae, 0x16, 0x13, 0xc6, 0x8b, 0x48, 0xce, 0xdd, 0x02, 0xd5, 0x83, 0xbe,
	0x15, 0x6c, 0x15, 0x3b, 0x21, 0xcc, 0x9e, 0x41, 0xf7, 0x46, 0x47, 0x59, 0x18, 0xaf, 0xe2, 0x04,
	0xc3, 0x34, 0xd2, 0xd7, 0xa8, 0xa9, 0x92, 0x66, 0x70, 0x68, 0xc1, 0xc4, 0xe6, 0xbf, 0x50, 0x7a,
	0x6b, 0x83, 0xea, 0xdb, 0x1b, 0x54, 0xdd, 0x3b, 0x7f, 0x7b, 0xef, 0xbe, 0xc1, 0x60, 0x4e, 0x0d,
	0xbf, 0x73, 0x3f, 0xaa, 0xa9, 0x41, 0x35, 0x0d, 0xc6, 0xf7, 0x7d, 0x92, 0xa0, 0xef, 0xbc, 0x5b,
	0x49, 0x5b, 0xf4, 0xd9, 0x27, 0xe8, 0xfc, 0x90, 0x99, 0x0e, 0x8b, 0x9f, 0x4f, 0x88, 0xb7, 0x86,
	0xf5, 0xc7, 0x1f, 0x97, 0x52, 0x64, 0xa8, 0xbf, 0xa2, 0xb9, 0x51, 0xfa, 0x3a, 0xbf, 0x44, 0x99,
	0x2b, 0x9d, 0xf3, 0xde, 0xd0, 0x1b, 0xb5, 0x4f, 0x9b, 0xe3, 0xa2, 0xb1, 0xc1, 0x81, 0xb5, 0x16,
	0xc1, 0xf9, 0xad, 0xf9, 0x17, 0x00, 0x00, 0xff, 0xff, 0x67, 0x78, 0xc6, 0xc8, 0xc3, 0x04, 0x00,
	0x00,
}
