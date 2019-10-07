// Code generated by protoc-gen-go. DO NOT EDIT.
// source: cmerror_data.proto

package gen

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

//
// Juniper Error Item information
//
type ErrorData struct {
	// Identifier that uniquely identifies the source of
	// the error.
	Identifier  *string `protobuf:"bytes,1,opt,name=identifier" json:"identifier,omitempty"`
	Count       *uint64 `protobuf:"varint,2,opt,name=count" json:"count,omitempty"`
	LastUpdated *uint64 `protobuf:"varint,3,opt,name=last_updated,json=lastUpdated" json:"last_updated,omitempty"`
	IsEnabled   *bool   `protobuf:"varint,4,opt,name=is_enabled,json=isEnabled" json:"is_enabled,omitempty"`
	// Additional Metadata for error processing
	ModuleId    *uint32 `protobuf:"varint,5,opt,name=module_id,json=moduleId" json:"module_id,omitempty"`
	ComponentId *uint32 `protobuf:"varint,6,opt,name=component_id,json=componentId" json:"component_id,omitempty"`
	ErrorCode   *uint32 `protobuf:"varint,7,opt,name=error_code,json=errorCode" json:"error_code,omitempty"`
	// Additional stats for each of the error
	OccurCount        *uint32 `protobuf:"varint,8,opt,name=occur_count,json=occurCount" json:"occur_count,omitempty"`
	ClearedCount      *uint32 `protobuf:"varint,9,opt,name=cleared_count,json=clearedCount" json:"cleared_count,omitempty"`
	LastClearedAt     *uint64 `protobuf:"varint,10,opt,name=last_cleared_at,json=lastClearedAt" json:"last_cleared_at,omitempty"`
	ActionCount       *uint32 `protobuf:"varint,11,opt,name=action_count,json=actionCount" json:"action_count,omitempty"`
	LastActionTakenAt *uint64 `protobuf:"varint,12,opt,name=last_action_taken_at,json=lastActionTakenAt" json:"last_action_taken_at,omitempty"`
	// Fru information
	FruType *string `protobuf:"bytes,13,opt,name=fru_type,json=fruType" json:"fru_type,omitempty"`
	FruSlot *uint32 `protobuf:"varint,14,opt,name=fru_slot,json=fruSlot" json:"fru_slot,omitempty"`
	// Help information regarding the error.
	Description          *string  `protobuf:"bytes,15,opt,name=description" json:"description,omitempty"`
	Help                 *string  `protobuf:"bytes,16,opt,name=help" json:"help,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ErrorData) Reset()         { *m = ErrorData{} }
func (m *ErrorData) String() string { return proto.CompactTextString(m) }
func (*ErrorData) ProtoMessage()    {}
func (*ErrorData) Descriptor() ([]byte, []int) {
	return fileDescriptor_cmerror_data_773ade7b8b3b1225, []int{0}
}
func (m *ErrorData) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ErrorData.Unmarshal(m, b)
}
func (m *ErrorData) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ErrorData.Marshal(b, m, deterministic)
}
func (dst *ErrorData) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ErrorData.Merge(dst, src)
}
func (m *ErrorData) XXX_Size() int {
	return xxx_messageInfo_ErrorData.Size(m)
}
func (m *ErrorData) XXX_DiscardUnknown() {
	xxx_messageInfo_ErrorData.DiscardUnknown(m)
}

var xxx_messageInfo_ErrorData proto.InternalMessageInfo

func (m *ErrorData) GetIdentifier() string {
	if m != nil && m.Identifier != nil {
		return *m.Identifier
	}
	return ""
}

func (m *ErrorData) GetCount() uint64 {
	if m != nil && m.Count != nil {
		return *m.Count
	}
	return 0
}

func (m *ErrorData) GetLastUpdated() uint64 {
	if m != nil && m.LastUpdated != nil {
		return *m.LastUpdated
	}
	return 0
}

func (m *ErrorData) GetIsEnabled() bool {
	if m != nil && m.IsEnabled != nil {
		return *m.IsEnabled
	}
	return false
}

func (m *ErrorData) GetModuleId() uint32 {
	if m != nil && m.ModuleId != nil {
		return *m.ModuleId
	}
	return 0
}

func (m *ErrorData) GetComponentId() uint32 {
	if m != nil && m.ComponentId != nil {
		return *m.ComponentId
	}
	return 0
}

func (m *ErrorData) GetErrorCode() uint32 {
	if m != nil && m.ErrorCode != nil {
		return *m.ErrorCode
	}
	return 0
}

func (m *ErrorData) GetOccurCount() uint32 {
	if m != nil && m.OccurCount != nil {
		return *m.OccurCount
	}
	return 0
}

func (m *ErrorData) GetClearedCount() uint32 {
	if m != nil && m.ClearedCount != nil {
		return *m.ClearedCount
	}
	return 0
}

func (m *ErrorData) GetLastClearedAt() uint64 {
	if m != nil && m.LastClearedAt != nil {
		return *m.LastClearedAt
	}
	return 0
}

func (m *ErrorData) GetActionCount() uint32 {
	if m != nil && m.ActionCount != nil {
		return *m.ActionCount
	}
	return 0
}

func (m *ErrorData) GetLastActionTakenAt() uint64 {
	if m != nil && m.LastActionTakenAt != nil {
		return *m.LastActionTakenAt
	}
	return 0
}

func (m *ErrorData) GetFruType() string {
	if m != nil && m.FruType != nil {
		return *m.FruType
	}
	return ""
}

func (m *ErrorData) GetFruSlot() uint32 {
	if m != nil && m.FruSlot != nil {
		return *m.FruSlot
	}
	return 0
}

func (m *ErrorData) GetDescription() string {
	if m != nil && m.Description != nil {
		return *m.Description
	}
	return ""
}

func (m *ErrorData) GetHelp() string {
	if m != nil && m.Help != nil {
		return *m.Help
	}
	return ""
}

//
// Top-level CmerrorData message
//
type CmerrorData struct {
	// collection of error items
	ErrorItem            []*ErrorData `protobuf:"bytes,1,rep,name=error_item,json=errorItem" json:"error_item,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *CmerrorData) Reset()         { *m = CmerrorData{} }
func (m *CmerrorData) String() string { return proto.CompactTextString(m) }
func (*CmerrorData) ProtoMessage()    {}
func (*CmerrorData) Descriptor() ([]byte, []int) {
	return fileDescriptor_cmerror_data_773ade7b8b3b1225, []int{1}
}
func (m *CmerrorData) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CmerrorData.Unmarshal(m, b)
}
func (m *CmerrorData) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CmerrorData.Marshal(b, m, deterministic)
}
func (dst *CmerrorData) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CmerrorData.Merge(dst, src)
}
func (m *CmerrorData) XXX_Size() int {
	return xxx_messageInfo_CmerrorData.Size(m)
}
func (m *CmerrorData) XXX_DiscardUnknown() {
	xxx_messageInfo_CmerrorData.DiscardUnknown(m)
}

var xxx_messageInfo_CmerrorData proto.InternalMessageInfo

func (m *CmerrorData) GetErrorItem() []*ErrorData {
	if m != nil {
		return m.ErrorItem
	}
	return nil
}

var E_JnprCmerrorDataExt = &proto.ExtensionDesc{
	ExtendedType:  (*JuniperNetworksSensors)(nil),
	ExtensionType: (*CmerrorData)(nil),
	Field:         21,
	Name:          "jnpr_cmerror_data_ext",
	Tag:           "bytes,21,opt,name=jnpr_cmerror_data_ext,json=jnprCmerrorDataExt",
	Filename:      "cmerror_data.proto",
}

func init() {
	proto.RegisterType((*ErrorData)(nil), "ErrorData")
	proto.RegisterType((*CmerrorData)(nil), "CmerrorData")
	proto.RegisterExtension(E_JnprCmerrorDataExt)
}

func init() { proto.RegisterFile("cmerror_data.proto", fileDescriptor_cmerror_data_773ade7b8b3b1225) }

var fileDescriptor_cmerror_data_773ade7b8b3b1225 = []byte{
	// 464 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x4c, 0x92, 0x41, 0x6f, 0xd3, 0x30,
	0x14, 0xc7, 0x95, 0xad, 0x63, 0xcd, 0x4b, 0xcb, 0x86, 0x61, 0xc2, 0x30, 0x21, 0x85, 0x49, 0xa0,
	0x80, 0x44, 0x0f, 0x3b, 0x20, 0xc4, 0x89, 0x52, 0x7a, 0x28, 0x07, 0x0e, 0xd9, 0x90, 0xb8, 0x45,
	0x26, 0x7e, 0x15, 0x66, 0x89, 0x6d, 0x39, 0x2f, 0x62, 0xbd, 0xf2, 0xe9, 0xf8, 0x58, 0xc8, 0x76,
	0x5a, 0xe5, 0x56, 0xfd, 0x7f, 0xbf, 0xfc, 0xab, 0xf7, 0xfc, 0x80, 0xd5, 0x2d, 0x3a, 0x67, 0x5c,
	0x25, 0x05, 0x89, 0x85, 0x75, 0x86, 0xcc, 0xf3, 0xc7, 0x84, 0x0d, 0xb6, 0x48, 0x6e, 0x57, 0x91,
	0xb1, 0x31, 0xbc, 0xfa, 0x37, 0x81, 0x74, 0xed, 0xcd, 0x2f, 0x82, 0x04, 0x7b, 0x05, 0xa0, 0x24,
	0x6a, 0x52, 0x5b, 0x85, 0x8e, 0x27, 0x79, 0x52, 0xa4, 0x9f, 0x4f, 0xfe, 0x7e, 0x3a, 0x9a, 0x26,
	0xe5, 0x08, 0xb0, 0x4b, 0x38, 0xa9, 0x4d, 0xaf, 0x89, 0x1f, 0xe5, 0x49, 0x31, 0x09, 0x06, 0x4f,
	0xca, 0x98, 0xb1, 0x02, 0x66, 0x8d, 0xe8, 0xa8, 0xea, 0xad, 0x14, 0x84, 0x92, 0x1f, 0x1f, 0x9c,
	0xf3, 0xa4, 0xcc, 0x3c, 0xfa, 0x1e, 0x09, 0x7b, 0x01, 0xa0, 0xba, 0x0a, 0xb5, 0xf8, 0xd9, 0xa0,
	0xe4, 0x93, 0x3c, 0x29, 0xa6, 0x65, 0xaa, 0xba, 0x75, 0x0c, 0xd8, 0x25, 0xa4, 0xad, 0x91, 0x7d,
	0x83, 0x95, 0x92, 0xfc, 0x24, 0x4f, 0x8a, 0x79, 0x39, 0x8d, 0xc1, 0x46, 0xb2, 0x97, 0x30, 0xab,
	0x4d, 0x6b, 0x8d, 0x46, 0x4d, 0x9e, 0x3f, 0x08, 0x3c, 0x3b, 0x64, 0x9b, 0x50, 0x1f, 0x77, 0x50,
	0x1b, 0x89, 0xfc, 0x34, 0x08, 0x69, 0x48, 0x56, 0x46, 0x22, 0x7b, 0x0d, 0x99, 0xa9, 0xeb, 0xde,
	0x63, 0x3f, 0xca, 0xd4, 0xf3, 0xfd, 0x28, 0x10, 0xc8, 0x2a, 0xcc, 0xf3, 0x16, 0xe6, 0x75, 0x83,
	0xc2, 0xa1, 0x1c, 0xcc, 0x74, 0x6c, 0xce, 0x06, 0x16, 0xdd, 0x77, 0x70, 0x16, 0x66, 0xdf, 0x7f,
	0x20, 0x88, 0xc3, 0x78, 0xfc, 0xb9, 0xa7, 0xab, 0x08, 0x97, 0x61, 0x55, 0xa2, 0x26, 0x65, 0xf4,
	0xd0, 0x9c, 0x8d, 0x9b, 0xb3, 0x88, 0x62, 0xf1, 0x7b, 0x78, 0x12, 0x8a, 0x07, 0x9d, 0xc4, 0x1d,
	0x6a, 0xdf, 0x3e, 0x1b, 0xb7, 0x3f, 0xf2, 0xca, 0x32, 0x18, 0xb7, 0x5e, 0x58, 0x12, 0x7b, 0x06,
	0xd3, 0xad, 0xeb, 0x2b, 0xda, 0x59, 0xe4, 0x73, 0xff, 0x9c, 0xe5, 0xe9, 0xd6, 0xf5, 0xb7, 0x3b,
	0x8b, 0x7b, 0xd4, 0x35, 0x86, 0xf8, 0xc3, 0xb0, 0x1c, 0x8f, 0x6e, 0x1a, 0x43, 0x2c, 0x87, 0x4c,
	0x62, 0x57, 0x3b, 0x65, 0x7d, 0x17, 0x3f, 0x0b, 0x1f, 0x8e, 0x23, 0xc6, 0x60, 0xf2, 0x0b, 0x1b,
	0xcb, 0xcf, 0x03, 0x0a, 0xbf, 0xaf, 0x3e, 0x40, 0xb6, 0x8a, 0x57, 0x17, 0x6e, 0xe9, 0xcd, 0x7e,
	0xfd, 0x8a, 0xb0, 0xe5, 0x49, 0x7e, 0x5c, 0x64, 0xd7, 0xb0, 0x38, 0xdc, 0xda, 0xf0, 0x14, 0x1b,
	0xc2, 0xf6, 0xe3, 0x0f, 0xb8, 0xf8, 0xad, 0xad, 0xab, 0xc6, 0x47, 0x5b, 0xe1, 0x3d, 0xb1, 0xa7,
	0x8b, 0xaf, 0xbd, 0x56, 0x16, 0xdd, 0x37, 0xa4, 0x3f, 0xc6, 0xdd, 0x75, 0x37, 0xa8, 0x3b, 0xe3,
	0x3a, 0x7e, 0x91, 0x27, 0x45, 0x76, 0x3d, 0x5b, 0x8c, 0xfe, 0xb0, 0x64, 0xbe, 0x63, 0x14, 0xac,
	0xef, 0xe9, 0x7f, 0x00, 0x00, 0x00, 0xff, 0xff, 0x6d, 0x05, 0xf1, 0x39, 0x08, 0x03, 0x00, 0x00,
}
