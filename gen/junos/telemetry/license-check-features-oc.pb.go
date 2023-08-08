// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: license-check-features-oc.proto

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

type SystemLicense struct {
	License              *SystemLicenseLicenseType `protobuf:"bytes,151,opt,name=license" json:"license,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                  `json:"-"`
	XXX_unrecognized     []byte                    `json:"-"`
	XXX_sizecache        int32                     `json:"-"`
}

func (m *SystemLicense) Reset()         { *m = SystemLicense{} }
func (m *SystemLicense) String() string { return proto.CompactTextString(m) }
func (*SystemLicense) ProtoMessage()    {}
func (*SystemLicense) Descriptor() ([]byte, []int) {
	return fileDescriptor_7bfd55ab599e2f32, []int{0}
}
func (m *SystemLicense) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SystemLicense.Unmarshal(m, b)
}
func (m *SystemLicense) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SystemLicense.Marshal(b, m, deterministic)
}
func (m *SystemLicense) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SystemLicense.Merge(m, src)
}
func (m *SystemLicense) XXX_Size() int {
	return xxx_messageInfo_SystemLicense.Size(m)
}
func (m *SystemLicense) XXX_DiscardUnknown() {
	xxx_messageInfo_SystemLicense.DiscardUnknown(m)
}

var xxx_messageInfo_SystemLicense proto.InternalMessageInfo

func (m *SystemLicense) GetLicense() *SystemLicenseLicenseType {
	if m != nil {
		return m.License
	}
	return nil
}

type SystemLicenseLicenseType struct {
	Licenses             *SystemLicenseLicenseTypeLicensesType `protobuf:"bytes,151,opt,name=licenses" json:"licenses,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                              `json:"-"`
	XXX_unrecognized     []byte                                `json:"-"`
	XXX_sizecache        int32                                 `json:"-"`
}

func (m *SystemLicenseLicenseType) Reset()         { *m = SystemLicenseLicenseType{} }
func (m *SystemLicenseLicenseType) String() string { return proto.CompactTextString(m) }
func (*SystemLicenseLicenseType) ProtoMessage()    {}
func (*SystemLicenseLicenseType) Descriptor() ([]byte, []int) {
	return fileDescriptor_7bfd55ab599e2f32, []int{0, 0}
}
func (m *SystemLicenseLicenseType) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SystemLicenseLicenseType.Unmarshal(m, b)
}
func (m *SystemLicenseLicenseType) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SystemLicenseLicenseType.Marshal(b, m, deterministic)
}
func (m *SystemLicenseLicenseType) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SystemLicenseLicenseType.Merge(m, src)
}
func (m *SystemLicenseLicenseType) XXX_Size() int {
	return xxx_messageInfo_SystemLicenseLicenseType.Size(m)
}
func (m *SystemLicenseLicenseType) XXX_DiscardUnknown() {
	xxx_messageInfo_SystemLicenseLicenseType.DiscardUnknown(m)
}

var xxx_messageInfo_SystemLicenseLicenseType proto.InternalMessageInfo

func (m *SystemLicenseLicenseType) GetLicenses() *SystemLicenseLicenseTypeLicensesType {
	if m != nil {
		return m.Licenses
	}
	return nil
}

type SystemLicenseLicenseTypeLicensesType struct {
	License              []*SystemLicenseLicenseTypeLicensesTypeLicenseList `protobuf:"bytes,151,rep,name=license" json:"license,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                                           `json:"-"`
	XXX_unrecognized     []byte                                             `json:"-"`
	XXX_sizecache        int32                                              `json:"-"`
}

func (m *SystemLicenseLicenseTypeLicensesType) Reset()         { *m = SystemLicenseLicenseTypeLicensesType{} }
func (m *SystemLicenseLicenseTypeLicensesType) String() string { return proto.CompactTextString(m) }
func (*SystemLicenseLicenseTypeLicensesType) ProtoMessage()    {}
func (*SystemLicenseLicenseTypeLicensesType) Descriptor() ([]byte, []int) {
	return fileDescriptor_7bfd55ab599e2f32, []int{0, 0, 0}
}
func (m *SystemLicenseLicenseTypeLicensesType) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SystemLicenseLicenseTypeLicensesType.Unmarshal(m, b)
}
func (m *SystemLicenseLicenseTypeLicensesType) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SystemLicenseLicenseTypeLicensesType.Marshal(b, m, deterministic)
}
func (m *SystemLicenseLicenseTypeLicensesType) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SystemLicenseLicenseTypeLicensesType.Merge(m, src)
}
func (m *SystemLicenseLicenseTypeLicensesType) XXX_Size() int {
	return xxx_messageInfo_SystemLicenseLicenseTypeLicensesType.Size(m)
}
func (m *SystemLicenseLicenseTypeLicensesType) XXX_DiscardUnknown() {
	xxx_messageInfo_SystemLicenseLicenseTypeLicensesType.DiscardUnknown(m)
}

var xxx_messageInfo_SystemLicenseLicenseTypeLicensesType proto.InternalMessageInfo

func (m *SystemLicenseLicenseTypeLicensesType) GetLicense() []*SystemLicenseLicenseTypeLicensesTypeLicenseList {
	if m != nil {
		return m.License
	}
	return nil
}

type SystemLicenseLicenseTypeLicensesTypeLicenseList struct {
	LicenseId            *string                                                   `protobuf:"bytes,151,opt,name=license_id,json=licenseId" json:"license_id,omitempty"`
	State                *SystemLicenseLicenseTypeLicensesTypeLicenseListStateType `protobuf:"bytes,152,opt,name=state" json:"state,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                                                  `json:"-"`
	XXX_unrecognized     []byte                                                    `json:"-"`
	XXX_sizecache        int32                                                     `json:"-"`
}

func (m *SystemLicenseLicenseTypeLicensesTypeLicenseList) Reset() {
	*m = SystemLicenseLicenseTypeLicensesTypeLicenseList{}
}
func (m *SystemLicenseLicenseTypeLicensesTypeLicenseList) String() string {
	return proto.CompactTextString(m)
}
func (*SystemLicenseLicenseTypeLicensesTypeLicenseList) ProtoMessage() {}
func (*SystemLicenseLicenseTypeLicensesTypeLicenseList) Descriptor() ([]byte, []int) {
	return fileDescriptor_7bfd55ab599e2f32, []int{0, 0, 0, 0}
}
func (m *SystemLicenseLicenseTypeLicensesTypeLicenseList) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SystemLicenseLicenseTypeLicensesTypeLicenseList.Unmarshal(m, b)
}
func (m *SystemLicenseLicenseTypeLicensesTypeLicenseList) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SystemLicenseLicenseTypeLicensesTypeLicenseList.Marshal(b, m, deterministic)
}
func (m *SystemLicenseLicenseTypeLicensesTypeLicenseList) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SystemLicenseLicenseTypeLicensesTypeLicenseList.Merge(m, src)
}
func (m *SystemLicenseLicenseTypeLicensesTypeLicenseList) XXX_Size() int {
	return xxx_messageInfo_SystemLicenseLicenseTypeLicensesTypeLicenseList.Size(m)
}
func (m *SystemLicenseLicenseTypeLicensesTypeLicenseList) XXX_DiscardUnknown() {
	xxx_messageInfo_SystemLicenseLicenseTypeLicensesTypeLicenseList.DiscardUnknown(m)
}

var xxx_messageInfo_SystemLicenseLicenseTypeLicensesTypeLicenseList proto.InternalMessageInfo

func (m *SystemLicenseLicenseTypeLicensesTypeLicenseList) GetLicenseId() string {
	if m != nil && m.LicenseId != nil {
		return *m.LicenseId
	}
	return ""
}

func (m *SystemLicenseLicenseTypeLicensesTypeLicenseList) GetState() *SystemLicenseLicenseTypeLicensesTypeLicenseListStateType {
	if m != nil {
		return m.State
	}
	return nil
}

type SystemLicenseLicenseTypeLicensesTypeLicenseListStateType struct {
	LicenseId            *string  `protobuf:"bytes,51,opt,name=license_id,json=licenseId" json:"license_id,omitempty"`
	LicenseData          *string  `protobuf:"bytes,52,opt,name=license_data,json=licenseData" json:"license_data,omitempty"`
	Active               *bool    `protobuf:"varint,53,opt,name=active" json:"active,omitempty"`
	Description          *string  `protobuf:"bytes,54,opt,name=description" json:"description,omitempty"`
	ExpirationDate       *uint64  `protobuf:"varint,55,opt,name=expiration_date,json=expirationDate" json:"expiration_date,omitempty"`
	InUse                *bool    `protobuf:"varint,56,opt,name=in_use,json=inUse" json:"in_use,omitempty"`
	Expired              *bool    `protobuf:"varint,57,opt,name=expired" json:"expired,omitempty"`
	Valid                *bool    `protobuf:"varint,58,opt,name=valid" json:"valid,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SystemLicenseLicenseTypeLicensesTypeLicenseListStateType) Reset() {
	*m = SystemLicenseLicenseTypeLicensesTypeLicenseListStateType{}
}
func (m *SystemLicenseLicenseTypeLicensesTypeLicenseListStateType) String() string {
	return proto.CompactTextString(m)
}
func (*SystemLicenseLicenseTypeLicensesTypeLicenseListStateType) ProtoMessage() {}
func (*SystemLicenseLicenseTypeLicensesTypeLicenseListStateType) Descriptor() ([]byte, []int) {
	return fileDescriptor_7bfd55ab599e2f32, []int{0, 0, 0, 0, 0}
}
func (m *SystemLicenseLicenseTypeLicensesTypeLicenseListStateType) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SystemLicenseLicenseTypeLicensesTypeLicenseListStateType.Unmarshal(m, b)
}
func (m *SystemLicenseLicenseTypeLicensesTypeLicenseListStateType) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SystemLicenseLicenseTypeLicensesTypeLicenseListStateType.Marshal(b, m, deterministic)
}
func (m *SystemLicenseLicenseTypeLicensesTypeLicenseListStateType) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SystemLicenseLicenseTypeLicensesTypeLicenseListStateType.Merge(m, src)
}
func (m *SystemLicenseLicenseTypeLicensesTypeLicenseListStateType) XXX_Size() int {
	return xxx_messageInfo_SystemLicenseLicenseTypeLicensesTypeLicenseListStateType.Size(m)
}
func (m *SystemLicenseLicenseTypeLicensesTypeLicenseListStateType) XXX_DiscardUnknown() {
	xxx_messageInfo_SystemLicenseLicenseTypeLicensesTypeLicenseListStateType.DiscardUnknown(m)
}

var xxx_messageInfo_SystemLicenseLicenseTypeLicensesTypeLicenseListStateType proto.InternalMessageInfo

func (m *SystemLicenseLicenseTypeLicensesTypeLicenseListStateType) GetLicenseId() string {
	if m != nil && m.LicenseId != nil {
		return *m.LicenseId
	}
	return ""
}

func (m *SystemLicenseLicenseTypeLicensesTypeLicenseListStateType) GetLicenseData() string {
	if m != nil && m.LicenseData != nil {
		return *m.LicenseData
	}
	return ""
}

func (m *SystemLicenseLicenseTypeLicensesTypeLicenseListStateType) GetActive() bool {
	if m != nil && m.Active != nil {
		return *m.Active
	}
	return false
}

func (m *SystemLicenseLicenseTypeLicensesTypeLicenseListStateType) GetDescription() string {
	if m != nil && m.Description != nil {
		return *m.Description
	}
	return ""
}

func (m *SystemLicenseLicenseTypeLicensesTypeLicenseListStateType) GetExpirationDate() uint64 {
	if m != nil && m.ExpirationDate != nil {
		return *m.ExpirationDate
	}
	return 0
}

func (m *SystemLicenseLicenseTypeLicensesTypeLicenseListStateType) GetInUse() bool {
	if m != nil && m.InUse != nil {
		return *m.InUse
	}
	return false
}

func (m *SystemLicenseLicenseTypeLicensesTypeLicenseListStateType) GetExpired() bool {
	if m != nil && m.Expired != nil {
		return *m.Expired
	}
	return false
}

func (m *SystemLicenseLicenseTypeLicensesTypeLicenseListStateType) GetValid() bool {
	if m != nil && m.Valid != nil {
		return *m.Valid
	}
	return false
}

var E_JnprSystemLicenseExt = &proto.ExtensionDesc{
	ExtendedType:  (*JuniperNetworksSensors)(nil),
	ExtensionType: (*SystemLicense)(nil),
	Field:         108,
	Name:          "jnpr_system_license_ext",
	Tag:           "bytes,108,opt,name=jnpr_system_license_ext",
	Filename:      "license-check-features-oc.proto",
}

func init() {
	proto.RegisterType((*SystemLicense)(nil), "system_license")
	proto.RegisterType((*SystemLicenseLicenseType)(nil), "system_license.license_type")
	proto.RegisterType((*SystemLicenseLicenseTypeLicensesType)(nil), "system_license.license_type.licenses_type")
	proto.RegisterType((*SystemLicenseLicenseTypeLicensesTypeLicenseList)(nil), "system_license.license_type.licenses_type.license_list")
	proto.RegisterType((*SystemLicenseLicenseTypeLicensesTypeLicenseListStateType)(nil), "system_license.license_type.licenses_type.license_list.state_type")
	proto.RegisterExtension(E_JnprSystemLicenseExt)
}

func init() { proto.RegisterFile("license-check-features-oc.proto", fileDescriptor_7bfd55ab599e2f32) }

var fileDescriptor_7bfd55ab599e2f32 = []byte{
	// 415 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x50, 0x4f, 0x8b, 0xd3, 0x40,
	0x14, 0x27, 0x76, 0xb3, 0xdb, 0x7d, 0xd1, 0x5d, 0x18, 0x57, 0x77, 0x08, 0x8a, 0x51, 0x10, 0x83,
	0xd0, 0x1c, 0xd6, 0x3f, 0xd5, 0x3d, 0x49, 0xa9, 0x88, 0x22, 0x1e, 0x52, 0x3c, 0x78, 0x0a, 0x43,
	0xf2, 0xc4, 0xb1, 0x69, 0x12, 0x66, 0x5e, 0x6b, 0x7b, 0xf5, 0xea, 0x07, 0xa8, 0xdf, 0xc9, 0x4f,
	0xe3, 0xcd, 0xa3, 0x64, 0x26, 0x69, 0x1a, 0x0f, 0x22, 0x9e, 0x92, 0xdf, 0x9f, 0xf7, 0x9b, 0xdf,
	0x7b, 0x70, 0x27, 0x97, 0x29, 0x16, 0x1a, 0x47, 0xe9, 0x27, 0x4c, 0xe7, 0xa3, 0x8f, 0x28, 0x68,
	0xa9, 0x50, 0x8f, 0xca, 0x34, 0xaa, 0x54, 0x49, 0xa5, 0x7f, 0x9d, 0x30, 0xc7, 0x05, 0x92, 0xda,
	0x24, 0x54, 0x56, 0x96, 0xbc, 0xb7, 0x75, 0xe1, 0x44, 0x6f, 0x34, 0xe1, 0x22, 0x69, 0xe6, 0xd9,
	0x18, 0x8e, 0x9a, 0x5f, 0xbe, 0x75, 0x02, 0x27, 0xf4, 0x2e, 0x6e, 0x45, 0x7d, 0x4b, 0xd4, 0x7c,
	0x13, 0xda, 0x54, 0x18, 0xb7, 0x6e, 0xff, 0xc7, 0x01, 0x5c, 0xdd, 0x57, 0xd8, 0x2b, 0x18, 0x36,
	0x58, 0xb7, 0x51, 0x0f, 0xff, 0x16, 0xd5, 0x02, 0x6d, 0x83, 0x77, 0xc3, 0xfe, 0xcf, 0x01, 0x5c,
	0xeb, 0x69, 0x2c, 0xee, 0x95, 0x1c, 0x84, 0xde, 0xc5, 0xf8, 0xdf, 0x93, 0x77, 0x52, 0x2e, 0x35,
	0x75, 0xfd, 0xbf, 0x0d, 0xba, 0xfe, 0xb5, 0xc2, 0xee, 0x03, 0xb4, 0x58, 0x66, 0x76, 0x83, 0xe3,
	0x89, 0xfb, 0xf5, 0xc5, 0x95, 0xa1, 0x13, 0x1f, 0x37, 0xca, 0xeb, 0x8c, 0x7d, 0x00, 0x57, 0x93,
	0x20, 0xe4, 0xdf, 0xed, 0x8e, 0x93, 0xff, 0x6c, 0x12, 0x99, 0x14, 0xbb, 0xbb, 0x4d, 0xf4, 0x7f,
	0x39, 0x00, 0x1d, 0xcb, 0x6e, 0xf7, 0x0a, 0x3d, 0xaa, 0xfb, 0xec, 0x17, 0xb9, 0xdb, 0xf5, 0xcf,
	0x04, 0x09, 0xfe, 0xd8, 0x18, 0xbc, 0x86, 0x9b, 0x0a, 0x12, 0xec, 0x26, 0x1c, 0x8a, 0x94, 0xe4,
	0x0a, 0xf9, 0x93, 0xc0, 0x09, 0x87, 0x71, 0x83, 0x58, 0x00, 0x5e, 0x86, 0x3a, 0x55, 0xb2, 0x22,
	0x59, 0x16, 0xfc, 0xa9, 0x9d, 0xdc, 0xa3, 0xd8, 0x03, 0x38, 0xc5, 0x75, 0x25, 0x95, 0xa8, 0x51,
	0x9d, 0x8f, 0x7c, 0x1c, 0x38, 0xe1, 0x41, 0x7c, 0xd2, 0xd1, 0x53, 0x41, 0xc8, 0x6e, 0xc0, 0xa1,
	0x2c, 0x92, 0xa5, 0x46, 0xfe, 0xcc, 0x3c, 0xe1, 0xca, 0xe2, 0xbd, 0x46, 0xc6, 0xe1, 0xc8, 0x18,
	0x31, 0xe3, 0xcf, 0x0d, 0xdf, 0x42, 0x76, 0x06, 0xee, 0x4a, 0xe4, 0x32, 0xe3, 0x97, 0xd6, 0x6f,
	0xc0, 0x65, 0x02, 0xe7, 0x9f, 0x8b, 0x4a, 0x25, 0xfd, 0x5b, 0x26, 0xb8, 0x26, 0x76, 0x1e, 0xbd,
	0x59, 0x16, 0xb2, 0x42, 0xf5, 0x0e, 0xe9, 0x4b, 0xa9, 0xe6, 0x7a, 0x86, 0x85, 0x2e, 0x95, 0xe6,
	0xb9, 0xb9, 0xff, 0xe9, 0x1f, 0xf7, 0x8f, 0xcf, 0xea, 0xa0, 0x99, 0xe1, 0xde, 0x5a, 0xea, 0xe5,
	0x9a, 0x7e, 0x07, 0x00, 0x00, 0xff, 0xff, 0x93, 0xa0, 0x7b, 0x4a, 0x31, 0x03, 0x00, 0x00,
}