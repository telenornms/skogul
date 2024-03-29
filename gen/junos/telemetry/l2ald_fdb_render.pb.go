// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: l2ald_fdb_render.proto

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

type NetworkInstancesFdb struct {
	NetworkInstance      []*NetworkInstancesFdbNetworkInstanceList `protobuf:"bytes,151,rep,name=network_instance,json=networkInstance" json:"network_instance,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                                  `json:"-"`
	XXX_unrecognized     []byte                                    `json:"-"`
	XXX_sizecache        int32                                     `json:"-"`
}

func (m *NetworkInstancesFdb) Reset()         { *m = NetworkInstancesFdb{} }
func (m *NetworkInstancesFdb) String() string { return proto.CompactTextString(m) }
func (*NetworkInstancesFdb) ProtoMessage()    {}
func (*NetworkInstancesFdb) Descriptor() ([]byte, []int) {
	return fileDescriptor_1194cdb0d4efc56b, []int{0}
}
func (m *NetworkInstancesFdb) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NetworkInstancesFdb.Unmarshal(m, b)
}
func (m *NetworkInstancesFdb) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NetworkInstancesFdb.Marshal(b, m, deterministic)
}
func (m *NetworkInstancesFdb) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NetworkInstancesFdb.Merge(m, src)
}
func (m *NetworkInstancesFdb) XXX_Size() int {
	return xxx_messageInfo_NetworkInstancesFdb.Size(m)
}
func (m *NetworkInstancesFdb) XXX_DiscardUnknown() {
	xxx_messageInfo_NetworkInstancesFdb.DiscardUnknown(m)
}

var xxx_messageInfo_NetworkInstancesFdb proto.InternalMessageInfo

func (m *NetworkInstancesFdb) GetNetworkInstance() []*NetworkInstancesFdbNetworkInstanceList {
	if m != nil {
		return m.NetworkInstance
	}
	return nil
}

type NetworkInstancesFdbNetworkInstanceList struct {
	Name                 *string                                                   `protobuf:"bytes,51,opt,name=name" json:"name,omitempty"`
	MacTableInfo         *NetworkInstancesFdbNetworkInstanceListMacTableInfoType   `protobuf:"bytes,171,opt,name=mac_table_info,json=macTableInfo" json:"mac_table_info,omitempty"`
	MacTable             *NetworkInstancesFdbNetworkInstanceListMacTableType       `protobuf:"bytes,151,opt,name=mac_table,json=macTable" json:"mac_table,omitempty"`
	MacipTableInfo       *NetworkInstancesFdbNetworkInstanceListMacipTableInfoType `protobuf:"bytes,181,opt,name=macip_table_info,json=macipTableInfo" json:"macip_table_info,omitempty"`
	MacipTable           *NetworkInstancesFdbNetworkInstanceListMacipTableType     `protobuf:"bytes,161,opt,name=macip_table,json=macipTable" json:"macip_table,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                                                  `json:"-"`
	XXX_unrecognized     []byte                                                    `json:"-"`
	XXX_sizecache        int32                                                     `json:"-"`
}

func (m *NetworkInstancesFdbNetworkInstanceList) Reset() {
	*m = NetworkInstancesFdbNetworkInstanceList{}
}
func (m *NetworkInstancesFdbNetworkInstanceList) String() string { return proto.CompactTextString(m) }
func (*NetworkInstancesFdbNetworkInstanceList) ProtoMessage()    {}
func (*NetworkInstancesFdbNetworkInstanceList) Descriptor() ([]byte, []int) {
	return fileDescriptor_1194cdb0d4efc56b, []int{0, 0}
}
func (m *NetworkInstancesFdbNetworkInstanceList) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NetworkInstancesFdbNetworkInstanceList.Unmarshal(m, b)
}
func (m *NetworkInstancesFdbNetworkInstanceList) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NetworkInstancesFdbNetworkInstanceList.Marshal(b, m, deterministic)
}
func (m *NetworkInstancesFdbNetworkInstanceList) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NetworkInstancesFdbNetworkInstanceList.Merge(m, src)
}
func (m *NetworkInstancesFdbNetworkInstanceList) XXX_Size() int {
	return xxx_messageInfo_NetworkInstancesFdbNetworkInstanceList.Size(m)
}
func (m *NetworkInstancesFdbNetworkInstanceList) XXX_DiscardUnknown() {
	xxx_messageInfo_NetworkInstancesFdbNetworkInstanceList.DiscardUnknown(m)
}

var xxx_messageInfo_NetworkInstancesFdbNetworkInstanceList proto.InternalMessageInfo

func (m *NetworkInstancesFdbNetworkInstanceList) GetName() string {
	if m != nil && m.Name != nil {
		return *m.Name
	}
	return ""
}

func (m *NetworkInstancesFdbNetworkInstanceList) GetMacTableInfo() *NetworkInstancesFdbNetworkInstanceListMacTableInfoType {
	if m != nil {
		return m.MacTableInfo
	}
	return nil
}

func (m *NetworkInstancesFdbNetworkInstanceList) GetMacTable() *NetworkInstancesFdbNetworkInstanceListMacTableType {
	if m != nil {
		return m.MacTable
	}
	return nil
}

func (m *NetworkInstancesFdbNetworkInstanceList) GetMacipTableInfo() *NetworkInstancesFdbNetworkInstanceListMacipTableInfoType {
	if m != nil {
		return m.MacipTableInfo
	}
	return nil
}

func (m *NetworkInstancesFdbNetworkInstanceList) GetMacipTable() *NetworkInstancesFdbNetworkInstanceListMacipTableType {
	if m != nil {
		return m.MacipTable
	}
	return nil
}

type NetworkInstancesFdbNetworkInstanceListMacTableInfoType struct {
	Learning             *bool    `protobuf:"varint,71,opt,name=learning" json:"learning,omitempty"`
	AgingTime            *uint32  `protobuf:"varint,72,opt,name=aging_time,json=agingTime" json:"aging_time,omitempty"`
	TableSize            *uint32  `protobuf:"varint,73,opt,name=table_size,json=tableSize" json:"table_size,omitempty"`
	NumLocalEntries      *uint32  `protobuf:"varint,74,opt,name=num_local_entries,json=numLocalEntries" json:"num_local_entries,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *NetworkInstancesFdbNetworkInstanceListMacTableInfoType) Reset() {
	*m = NetworkInstancesFdbNetworkInstanceListMacTableInfoType{}
}
func (m *NetworkInstancesFdbNetworkInstanceListMacTableInfoType) String() string {
	return proto.CompactTextString(m)
}
func (*NetworkInstancesFdbNetworkInstanceListMacTableInfoType) ProtoMessage() {}
func (*NetworkInstancesFdbNetworkInstanceListMacTableInfoType) Descriptor() ([]byte, []int) {
	return fileDescriptor_1194cdb0d4efc56b, []int{0, 0, 0}
}
func (m *NetworkInstancesFdbNetworkInstanceListMacTableInfoType) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NetworkInstancesFdbNetworkInstanceListMacTableInfoType.Unmarshal(m, b)
}
func (m *NetworkInstancesFdbNetworkInstanceListMacTableInfoType) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NetworkInstancesFdbNetworkInstanceListMacTableInfoType.Marshal(b, m, deterministic)
}
func (m *NetworkInstancesFdbNetworkInstanceListMacTableInfoType) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NetworkInstancesFdbNetworkInstanceListMacTableInfoType.Merge(m, src)
}
func (m *NetworkInstancesFdbNetworkInstanceListMacTableInfoType) XXX_Size() int {
	return xxx_messageInfo_NetworkInstancesFdbNetworkInstanceListMacTableInfoType.Size(m)
}
func (m *NetworkInstancesFdbNetworkInstanceListMacTableInfoType) XXX_DiscardUnknown() {
	xxx_messageInfo_NetworkInstancesFdbNetworkInstanceListMacTableInfoType.DiscardUnknown(m)
}

var xxx_messageInfo_NetworkInstancesFdbNetworkInstanceListMacTableInfoType proto.InternalMessageInfo

func (m *NetworkInstancesFdbNetworkInstanceListMacTableInfoType) GetLearning() bool {
	if m != nil && m.Learning != nil {
		return *m.Learning
	}
	return false
}

func (m *NetworkInstancesFdbNetworkInstanceListMacTableInfoType) GetAgingTime() uint32 {
	if m != nil && m.AgingTime != nil {
		return *m.AgingTime
	}
	return 0
}

func (m *NetworkInstancesFdbNetworkInstanceListMacTableInfoType) GetTableSize() uint32 {
	if m != nil && m.TableSize != nil {
		return *m.TableSize
	}
	return 0
}

func (m *NetworkInstancesFdbNetworkInstanceListMacTableInfoType) GetNumLocalEntries() uint32 {
	if m != nil && m.NumLocalEntries != nil {
		return *m.NumLocalEntries
	}
	return 0
}

type NetworkInstancesFdbNetworkInstanceListMacTableType struct {
	Entries              *NetworkInstancesFdbNetworkInstanceListMacTableTypeEntriesType `protobuf:"bytes,151,opt,name=entries" json:"entries,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                                                       `json:"-"`
	XXX_unrecognized     []byte                                                         `json:"-"`
	XXX_sizecache        int32                                                          `json:"-"`
}

func (m *NetworkInstancesFdbNetworkInstanceListMacTableType) Reset() {
	*m = NetworkInstancesFdbNetworkInstanceListMacTableType{}
}
func (m *NetworkInstancesFdbNetworkInstanceListMacTableType) String() string {
	return proto.CompactTextString(m)
}
func (*NetworkInstancesFdbNetworkInstanceListMacTableType) ProtoMessage() {}
func (*NetworkInstancesFdbNetworkInstanceListMacTableType) Descriptor() ([]byte, []int) {
	return fileDescriptor_1194cdb0d4efc56b, []int{0, 0, 1}
}
func (m *NetworkInstancesFdbNetworkInstanceListMacTableType) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NetworkInstancesFdbNetworkInstanceListMacTableType.Unmarshal(m, b)
}
func (m *NetworkInstancesFdbNetworkInstanceListMacTableType) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NetworkInstancesFdbNetworkInstanceListMacTableType.Marshal(b, m, deterministic)
}
func (m *NetworkInstancesFdbNetworkInstanceListMacTableType) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NetworkInstancesFdbNetworkInstanceListMacTableType.Merge(m, src)
}
func (m *NetworkInstancesFdbNetworkInstanceListMacTableType) XXX_Size() int {
	return xxx_messageInfo_NetworkInstancesFdbNetworkInstanceListMacTableType.Size(m)
}
func (m *NetworkInstancesFdbNetworkInstanceListMacTableType) XXX_DiscardUnknown() {
	xxx_messageInfo_NetworkInstancesFdbNetworkInstanceListMacTableType.DiscardUnknown(m)
}

var xxx_messageInfo_NetworkInstancesFdbNetworkInstanceListMacTableType proto.InternalMessageInfo

func (m *NetworkInstancesFdbNetworkInstanceListMacTableType) GetEntries() *NetworkInstancesFdbNetworkInstanceListMacTableTypeEntriesType {
	if m != nil {
		return m.Entries
	}
	return nil
}

type NetworkInstancesFdbNetworkInstanceListMacTableTypeEntriesType struct {
	Entry                []*NetworkInstancesFdbNetworkInstanceListMacTableTypeEntriesTypeEntryList `protobuf:"bytes,151,rep,name=entry" json:"entry,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                                                                  `json:"-"`
	XXX_unrecognized     []byte                                                                    `json:"-"`
	XXX_sizecache        int32                                                                     `json:"-"`
}

func (m *NetworkInstancesFdbNetworkInstanceListMacTableTypeEntriesType) Reset() {
	*m = NetworkInstancesFdbNetworkInstanceListMacTableTypeEntriesType{}
}
func (m *NetworkInstancesFdbNetworkInstanceListMacTableTypeEntriesType) String() string {
	return proto.CompactTextString(m)
}
func (*NetworkInstancesFdbNetworkInstanceListMacTableTypeEntriesType) ProtoMessage() {}
func (*NetworkInstancesFdbNetworkInstanceListMacTableTypeEntriesType) Descriptor() ([]byte, []int) {
	return fileDescriptor_1194cdb0d4efc56b, []int{0, 0, 1, 0}
}
func (m *NetworkInstancesFdbNetworkInstanceListMacTableTypeEntriesType) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NetworkInstancesFdbNetworkInstanceListMacTableTypeEntriesType.Unmarshal(m, b)
}
func (m *NetworkInstancesFdbNetworkInstanceListMacTableTypeEntriesType) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NetworkInstancesFdbNetworkInstanceListMacTableTypeEntriesType.Marshal(b, m, deterministic)
}
func (m *NetworkInstancesFdbNetworkInstanceListMacTableTypeEntriesType) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NetworkInstancesFdbNetworkInstanceListMacTableTypeEntriesType.Merge(m, src)
}
func (m *NetworkInstancesFdbNetworkInstanceListMacTableTypeEntriesType) XXX_Size() int {
	return xxx_messageInfo_NetworkInstancesFdbNetworkInstanceListMacTableTypeEntriesType.Size(m)
}
func (m *NetworkInstancesFdbNetworkInstanceListMacTableTypeEntriesType) XXX_DiscardUnknown() {
	xxx_messageInfo_NetworkInstancesFdbNetworkInstanceListMacTableTypeEntriesType.DiscardUnknown(m)
}

var xxx_messageInfo_NetworkInstancesFdbNetworkInstanceListMacTableTypeEntriesType proto.InternalMessageInfo

func (m *NetworkInstancesFdbNetworkInstanceListMacTableTypeEntriesType) GetEntry() []*NetworkInstancesFdbNetworkInstanceListMacTableTypeEntriesTypeEntryList {
	if m != nil {
		return m.Entry
	}
	return nil
}

type NetworkInstancesFdbNetworkInstanceListMacTableTypeEntriesTypeEntryList struct {
	MacAddress           *string  `protobuf:"bytes,51,opt,name=mac_address,json=macAddress" json:"mac_address,omitempty"`
	VlanId               *uint32  `protobuf:"varint,52,opt,name=vlan_id,json=vlanId" json:"vlan_id,omitempty"`
	Vni                  *uint32  `protobuf:"varint,53,opt,name=vni" json:"vni,omitempty"`
	VlanName             *string  `protobuf:"bytes,54,opt,name=vlan_name,json=vlanName" json:"vlan_name,omitempty"`
	Interface            *string  `protobuf:"bytes,55,opt,name=interface" json:"interface,omitempty"`
	EntryType            *string  `protobuf:"bytes,56,opt,name=entry_type,json=entryType" json:"entry_type,omitempty"`
	EventType            *string  `protobuf:"bytes,57,opt,name=event_type,json=eventType" json:"event_type,omitempty"`
	EthernetTagId        *uint32  `protobuf:"varint,58,opt,name=ethernet_tag_id,json=ethernetTagId" json:"ethernet_tag_id,omitempty"`
	ActiveSource         *string  `protobuf:"bytes,59,opt,name=active_source,json=activeSource" json:"active_source,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *NetworkInstancesFdbNetworkInstanceListMacTableTypeEntriesTypeEntryList) Reset() {
	*m = NetworkInstancesFdbNetworkInstanceListMacTableTypeEntriesTypeEntryList{}
}
func (m *NetworkInstancesFdbNetworkInstanceListMacTableTypeEntriesTypeEntryList) String() string {
	return proto.CompactTextString(m)
}
func (*NetworkInstancesFdbNetworkInstanceListMacTableTypeEntriesTypeEntryList) ProtoMessage() {}
func (*NetworkInstancesFdbNetworkInstanceListMacTableTypeEntriesTypeEntryList) Descriptor() ([]byte, []int) {
	return fileDescriptor_1194cdb0d4efc56b, []int{0, 0, 1, 0, 0}
}
func (m *NetworkInstancesFdbNetworkInstanceListMacTableTypeEntriesTypeEntryList) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NetworkInstancesFdbNetworkInstanceListMacTableTypeEntriesTypeEntryList.Unmarshal(m, b)
}
func (m *NetworkInstancesFdbNetworkInstanceListMacTableTypeEntriesTypeEntryList) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NetworkInstancesFdbNetworkInstanceListMacTableTypeEntriesTypeEntryList.Marshal(b, m, deterministic)
}
func (m *NetworkInstancesFdbNetworkInstanceListMacTableTypeEntriesTypeEntryList) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NetworkInstancesFdbNetworkInstanceListMacTableTypeEntriesTypeEntryList.Merge(m, src)
}
func (m *NetworkInstancesFdbNetworkInstanceListMacTableTypeEntriesTypeEntryList) XXX_Size() int {
	return xxx_messageInfo_NetworkInstancesFdbNetworkInstanceListMacTableTypeEntriesTypeEntryList.Size(m)
}
func (m *NetworkInstancesFdbNetworkInstanceListMacTableTypeEntriesTypeEntryList) XXX_DiscardUnknown() {
	xxx_messageInfo_NetworkInstancesFdbNetworkInstanceListMacTableTypeEntriesTypeEntryList.DiscardUnknown(m)
}

var xxx_messageInfo_NetworkInstancesFdbNetworkInstanceListMacTableTypeEntriesTypeEntryList proto.InternalMessageInfo

func (m *NetworkInstancesFdbNetworkInstanceListMacTableTypeEntriesTypeEntryList) GetMacAddress() string {
	if m != nil && m.MacAddress != nil {
		return *m.MacAddress
	}
	return ""
}

func (m *NetworkInstancesFdbNetworkInstanceListMacTableTypeEntriesTypeEntryList) GetVlanId() uint32 {
	if m != nil && m.VlanId != nil {
		return *m.VlanId
	}
	return 0
}

func (m *NetworkInstancesFdbNetworkInstanceListMacTableTypeEntriesTypeEntryList) GetVni() uint32 {
	if m != nil && m.Vni != nil {
		return *m.Vni
	}
	return 0
}

func (m *NetworkInstancesFdbNetworkInstanceListMacTableTypeEntriesTypeEntryList) GetVlanName() string {
	if m != nil && m.VlanName != nil {
		return *m.VlanName
	}
	return ""
}

func (m *NetworkInstancesFdbNetworkInstanceListMacTableTypeEntriesTypeEntryList) GetInterface() string {
	if m != nil && m.Interface != nil {
		return *m.Interface
	}
	return ""
}

func (m *NetworkInstancesFdbNetworkInstanceListMacTableTypeEntriesTypeEntryList) GetEntryType() string {
	if m != nil && m.EntryType != nil {
		return *m.EntryType
	}
	return ""
}

func (m *NetworkInstancesFdbNetworkInstanceListMacTableTypeEntriesTypeEntryList) GetEventType() string {
	if m != nil && m.EventType != nil {
		return *m.EventType
	}
	return ""
}

func (m *NetworkInstancesFdbNetworkInstanceListMacTableTypeEntriesTypeEntryList) GetEthernetTagId() uint32 {
	if m != nil && m.EthernetTagId != nil {
		return *m.EthernetTagId
	}
	return 0
}

func (m *NetworkInstancesFdbNetworkInstanceListMacTableTypeEntriesTypeEntryList) GetActiveSource() string {
	if m != nil && m.ActiveSource != nil {
		return *m.ActiveSource
	}
	return ""
}

type NetworkInstancesFdbNetworkInstanceListMacipTableInfoType struct {
	Learning             *bool    `protobuf:"varint,81,opt,name=learning" json:"learning,omitempty"`
	AgingTime            *uint32  `protobuf:"varint,82,opt,name=aging_time,json=agingTime" json:"aging_time,omitempty"`
	TableSize            *uint32  `protobuf:"varint,83,opt,name=table_size,json=tableSize" json:"table_size,omitempty"`
	ProxyMacip           *bool    `protobuf:"varint,84,opt,name=proxy_macip,json=proxyMacip" json:"proxy_macip,omitempty"`
	NumLocalEntries      *uint32  `protobuf:"varint,85,opt,name=num_local_entries,json=numLocalEntries" json:"num_local_entries,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *NetworkInstancesFdbNetworkInstanceListMacipTableInfoType) Reset() {
	*m = NetworkInstancesFdbNetworkInstanceListMacipTableInfoType{}
}
func (m *NetworkInstancesFdbNetworkInstanceListMacipTableInfoType) String() string {
	return proto.CompactTextString(m)
}
func (*NetworkInstancesFdbNetworkInstanceListMacipTableInfoType) ProtoMessage() {}
func (*NetworkInstancesFdbNetworkInstanceListMacipTableInfoType) Descriptor() ([]byte, []int) {
	return fileDescriptor_1194cdb0d4efc56b, []int{0, 0, 2}
}
func (m *NetworkInstancesFdbNetworkInstanceListMacipTableInfoType) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NetworkInstancesFdbNetworkInstanceListMacipTableInfoType.Unmarshal(m, b)
}
func (m *NetworkInstancesFdbNetworkInstanceListMacipTableInfoType) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NetworkInstancesFdbNetworkInstanceListMacipTableInfoType.Marshal(b, m, deterministic)
}
func (m *NetworkInstancesFdbNetworkInstanceListMacipTableInfoType) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NetworkInstancesFdbNetworkInstanceListMacipTableInfoType.Merge(m, src)
}
func (m *NetworkInstancesFdbNetworkInstanceListMacipTableInfoType) XXX_Size() int {
	return xxx_messageInfo_NetworkInstancesFdbNetworkInstanceListMacipTableInfoType.Size(m)
}
func (m *NetworkInstancesFdbNetworkInstanceListMacipTableInfoType) XXX_DiscardUnknown() {
	xxx_messageInfo_NetworkInstancesFdbNetworkInstanceListMacipTableInfoType.DiscardUnknown(m)
}

var xxx_messageInfo_NetworkInstancesFdbNetworkInstanceListMacipTableInfoType proto.InternalMessageInfo

func (m *NetworkInstancesFdbNetworkInstanceListMacipTableInfoType) GetLearning() bool {
	if m != nil && m.Learning != nil {
		return *m.Learning
	}
	return false
}

func (m *NetworkInstancesFdbNetworkInstanceListMacipTableInfoType) GetAgingTime() uint32 {
	if m != nil && m.AgingTime != nil {
		return *m.AgingTime
	}
	return 0
}

func (m *NetworkInstancesFdbNetworkInstanceListMacipTableInfoType) GetTableSize() uint32 {
	if m != nil && m.TableSize != nil {
		return *m.TableSize
	}
	return 0
}

func (m *NetworkInstancesFdbNetworkInstanceListMacipTableInfoType) GetProxyMacip() bool {
	if m != nil && m.ProxyMacip != nil {
		return *m.ProxyMacip
	}
	return false
}

func (m *NetworkInstancesFdbNetworkInstanceListMacipTableInfoType) GetNumLocalEntries() uint32 {
	if m != nil && m.NumLocalEntries != nil {
		return *m.NumLocalEntries
	}
	return 0
}

type NetworkInstancesFdbNetworkInstanceListMacipTableType struct {
	Entries              *NetworkInstancesFdbNetworkInstanceListMacipTableTypeEntriesType `protobuf:"bytes,161,opt,name=entries" json:"entries,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                                                         `json:"-"`
	XXX_unrecognized     []byte                                                           `json:"-"`
	XXX_sizecache        int32                                                            `json:"-"`
}

func (m *NetworkInstancesFdbNetworkInstanceListMacipTableType) Reset() {
	*m = NetworkInstancesFdbNetworkInstanceListMacipTableType{}
}
func (m *NetworkInstancesFdbNetworkInstanceListMacipTableType) String() string {
	return proto.CompactTextString(m)
}
func (*NetworkInstancesFdbNetworkInstanceListMacipTableType) ProtoMessage() {}
func (*NetworkInstancesFdbNetworkInstanceListMacipTableType) Descriptor() ([]byte, []int) {
	return fileDescriptor_1194cdb0d4efc56b, []int{0, 0, 3}
}
func (m *NetworkInstancesFdbNetworkInstanceListMacipTableType) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NetworkInstancesFdbNetworkInstanceListMacipTableType.Unmarshal(m, b)
}
func (m *NetworkInstancesFdbNetworkInstanceListMacipTableType) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NetworkInstancesFdbNetworkInstanceListMacipTableType.Marshal(b, m, deterministic)
}
func (m *NetworkInstancesFdbNetworkInstanceListMacipTableType) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NetworkInstancesFdbNetworkInstanceListMacipTableType.Merge(m, src)
}
func (m *NetworkInstancesFdbNetworkInstanceListMacipTableType) XXX_Size() int {
	return xxx_messageInfo_NetworkInstancesFdbNetworkInstanceListMacipTableType.Size(m)
}
func (m *NetworkInstancesFdbNetworkInstanceListMacipTableType) XXX_DiscardUnknown() {
	xxx_messageInfo_NetworkInstancesFdbNetworkInstanceListMacipTableType.DiscardUnknown(m)
}

var xxx_messageInfo_NetworkInstancesFdbNetworkInstanceListMacipTableType proto.InternalMessageInfo

func (m *NetworkInstancesFdbNetworkInstanceListMacipTableType) GetEntries() *NetworkInstancesFdbNetworkInstanceListMacipTableTypeEntriesType {
	if m != nil {
		return m.Entries
	}
	return nil
}

type NetworkInstancesFdbNetworkInstanceListMacipTableTypeEntriesType struct {
	Entry                []*NetworkInstancesFdbNetworkInstanceListMacipTableTypeEntriesTypeEntryList `protobuf:"bytes,161,rep,name=entry" json:"entry,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                                                                    `json:"-"`
	XXX_unrecognized     []byte                                                                      `json:"-"`
	XXX_sizecache        int32                                                                       `json:"-"`
}

func (m *NetworkInstancesFdbNetworkInstanceListMacipTableTypeEntriesType) Reset() {
	*m = NetworkInstancesFdbNetworkInstanceListMacipTableTypeEntriesType{}
}
func (m *NetworkInstancesFdbNetworkInstanceListMacipTableTypeEntriesType) String() string {
	return proto.CompactTextString(m)
}
func (*NetworkInstancesFdbNetworkInstanceListMacipTableTypeEntriesType) ProtoMessage() {}
func (*NetworkInstancesFdbNetworkInstanceListMacipTableTypeEntriesType) Descriptor() ([]byte, []int) {
	return fileDescriptor_1194cdb0d4efc56b, []int{0, 0, 3, 0}
}
func (m *NetworkInstancesFdbNetworkInstanceListMacipTableTypeEntriesType) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NetworkInstancesFdbNetworkInstanceListMacipTableTypeEntriesType.Unmarshal(m, b)
}
func (m *NetworkInstancesFdbNetworkInstanceListMacipTableTypeEntriesType) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NetworkInstancesFdbNetworkInstanceListMacipTableTypeEntriesType.Marshal(b, m, deterministic)
}
func (m *NetworkInstancesFdbNetworkInstanceListMacipTableTypeEntriesType) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NetworkInstancesFdbNetworkInstanceListMacipTableTypeEntriesType.Merge(m, src)
}
func (m *NetworkInstancesFdbNetworkInstanceListMacipTableTypeEntriesType) XXX_Size() int {
	return xxx_messageInfo_NetworkInstancesFdbNetworkInstanceListMacipTableTypeEntriesType.Size(m)
}
func (m *NetworkInstancesFdbNetworkInstanceListMacipTableTypeEntriesType) XXX_DiscardUnknown() {
	xxx_messageInfo_NetworkInstancesFdbNetworkInstanceListMacipTableTypeEntriesType.DiscardUnknown(m)
}

var xxx_messageInfo_NetworkInstancesFdbNetworkInstanceListMacipTableTypeEntriesType proto.InternalMessageInfo

func (m *NetworkInstancesFdbNetworkInstanceListMacipTableTypeEntriesType) GetEntry() []*NetworkInstancesFdbNetworkInstanceListMacipTableTypeEntriesTypeEntryList {
	if m != nil {
		return m.Entry
	}
	return nil
}

type NetworkInstancesFdbNetworkInstanceListMacipTableTypeEntriesTypeEntryList struct {
	IpAddress            *string  `protobuf:"bytes,51,opt,name=ip_address,json=ipAddress" json:"ip_address,omitempty"`
	MacAddress           *string  `protobuf:"bytes,52,opt,name=mac_address,json=macAddress" json:"mac_address,omitempty"`
	VlanId               *uint32  `protobuf:"varint,53,opt,name=vlan_id,json=vlanId" json:"vlan_id,omitempty"`
	Vni                  *uint32  `protobuf:"varint,54,opt,name=vni" json:"vni,omitempty"`
	VlanName             *string  `protobuf:"bytes,55,opt,name=vlan_name,json=vlanName" json:"vlan_name,omitempty"`
	Interface            *string  `protobuf:"bytes,56,opt,name=interface" json:"interface,omitempty"`
	L3Interface          *string  `protobuf:"bytes,57,opt,name=l3_interface,json=l3Interface" json:"l3_interface,omitempty"`
	EntryType            *string  `protobuf:"bytes,58,opt,name=entry_type,json=entryType" json:"entry_type,omitempty"`
	EventType            *string  `protobuf:"bytes,59,opt,name=event_type,json=eventType" json:"event_type,omitempty"`
	EthernetTagId        *uint32  `protobuf:"varint,60,opt,name=ethernet_tag_id,json=ethernetTagId" json:"ethernet_tag_id,omitempty"`
	ActiveSource         *string  `protobuf:"bytes,61,opt,name=active_source,json=activeSource" json:"active_source,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *NetworkInstancesFdbNetworkInstanceListMacipTableTypeEntriesTypeEntryList) Reset() {
	*m = NetworkInstancesFdbNetworkInstanceListMacipTableTypeEntriesTypeEntryList{}
}
func (m *NetworkInstancesFdbNetworkInstanceListMacipTableTypeEntriesTypeEntryList) String() string {
	return proto.CompactTextString(m)
}
func (*NetworkInstancesFdbNetworkInstanceListMacipTableTypeEntriesTypeEntryList) ProtoMessage() {}
func (*NetworkInstancesFdbNetworkInstanceListMacipTableTypeEntriesTypeEntryList) Descriptor() ([]byte, []int) {
	return fileDescriptor_1194cdb0d4efc56b, []int{0, 0, 3, 0, 0}
}
func (m *NetworkInstancesFdbNetworkInstanceListMacipTableTypeEntriesTypeEntryList) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NetworkInstancesFdbNetworkInstanceListMacipTableTypeEntriesTypeEntryList.Unmarshal(m, b)
}
func (m *NetworkInstancesFdbNetworkInstanceListMacipTableTypeEntriesTypeEntryList) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NetworkInstancesFdbNetworkInstanceListMacipTableTypeEntriesTypeEntryList.Marshal(b, m, deterministic)
}
func (m *NetworkInstancesFdbNetworkInstanceListMacipTableTypeEntriesTypeEntryList) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NetworkInstancesFdbNetworkInstanceListMacipTableTypeEntriesTypeEntryList.Merge(m, src)
}
func (m *NetworkInstancesFdbNetworkInstanceListMacipTableTypeEntriesTypeEntryList) XXX_Size() int {
	return xxx_messageInfo_NetworkInstancesFdbNetworkInstanceListMacipTableTypeEntriesTypeEntryList.Size(m)
}
func (m *NetworkInstancesFdbNetworkInstanceListMacipTableTypeEntriesTypeEntryList) XXX_DiscardUnknown() {
	xxx_messageInfo_NetworkInstancesFdbNetworkInstanceListMacipTableTypeEntriesTypeEntryList.DiscardUnknown(m)
}

var xxx_messageInfo_NetworkInstancesFdbNetworkInstanceListMacipTableTypeEntriesTypeEntryList proto.InternalMessageInfo

func (m *NetworkInstancesFdbNetworkInstanceListMacipTableTypeEntriesTypeEntryList) GetIpAddress() string {
	if m != nil && m.IpAddress != nil {
		return *m.IpAddress
	}
	return ""
}

func (m *NetworkInstancesFdbNetworkInstanceListMacipTableTypeEntriesTypeEntryList) GetMacAddress() string {
	if m != nil && m.MacAddress != nil {
		return *m.MacAddress
	}
	return ""
}

func (m *NetworkInstancesFdbNetworkInstanceListMacipTableTypeEntriesTypeEntryList) GetVlanId() uint32 {
	if m != nil && m.VlanId != nil {
		return *m.VlanId
	}
	return 0
}

func (m *NetworkInstancesFdbNetworkInstanceListMacipTableTypeEntriesTypeEntryList) GetVni() uint32 {
	if m != nil && m.Vni != nil {
		return *m.Vni
	}
	return 0
}

func (m *NetworkInstancesFdbNetworkInstanceListMacipTableTypeEntriesTypeEntryList) GetVlanName() string {
	if m != nil && m.VlanName != nil {
		return *m.VlanName
	}
	return ""
}

func (m *NetworkInstancesFdbNetworkInstanceListMacipTableTypeEntriesTypeEntryList) GetInterface() string {
	if m != nil && m.Interface != nil {
		return *m.Interface
	}
	return ""
}

func (m *NetworkInstancesFdbNetworkInstanceListMacipTableTypeEntriesTypeEntryList) GetL3Interface() string {
	if m != nil && m.L3Interface != nil {
		return *m.L3Interface
	}
	return ""
}

func (m *NetworkInstancesFdbNetworkInstanceListMacipTableTypeEntriesTypeEntryList) GetEntryType() string {
	if m != nil && m.EntryType != nil {
		return *m.EntryType
	}
	return ""
}

func (m *NetworkInstancesFdbNetworkInstanceListMacipTableTypeEntriesTypeEntryList) GetEventType() string {
	if m != nil && m.EventType != nil {
		return *m.EventType
	}
	return ""
}

func (m *NetworkInstancesFdbNetworkInstanceListMacipTableTypeEntriesTypeEntryList) GetEthernetTagId() uint32 {
	if m != nil && m.EthernetTagId != nil {
		return *m.EthernetTagId
	}
	return 0
}

func (m *NetworkInstancesFdbNetworkInstanceListMacipTableTypeEntriesTypeEntryList) GetActiveSource() string {
	if m != nil && m.ActiveSource != nil {
		return *m.ActiveSource
	}
	return ""
}

var E_JnprNetworkInstancesFdbExt = &proto.ExtensionDesc{
	ExtendedType:  (*JuniperNetworksSensors)(nil),
	ExtensionType: (*NetworkInstancesFdb)(nil),
	Field:         114,
	Name:          "jnpr_network_instances_fdb_ext",
	Tag:           "bytes,114,opt,name=jnpr_network_instances_fdb_ext",
	Filename:      "l2ald_fdb_render.proto",
}

func init() {
	proto.RegisterType((*NetworkInstancesFdb)(nil), "network_instances_fdb")
	proto.RegisterType((*NetworkInstancesFdbNetworkInstanceList)(nil), "network_instances_fdb.network_instance_list")
	proto.RegisterType((*NetworkInstancesFdbNetworkInstanceListMacTableInfoType)(nil), "network_instances_fdb.network_instance_list.mac_table_info_type")
	proto.RegisterType((*NetworkInstancesFdbNetworkInstanceListMacTableType)(nil), "network_instances_fdb.network_instance_list.mac_table_type")
	proto.RegisterType((*NetworkInstancesFdbNetworkInstanceListMacTableTypeEntriesType)(nil), "network_instances_fdb.network_instance_list.mac_table_type.entries_type")
	proto.RegisterType((*NetworkInstancesFdbNetworkInstanceListMacTableTypeEntriesTypeEntryList)(nil), "network_instances_fdb.network_instance_list.mac_table_type.entries_type.entry_list")
	proto.RegisterType((*NetworkInstancesFdbNetworkInstanceListMacipTableInfoType)(nil), "network_instances_fdb.network_instance_list.macip_table_info_type")
	proto.RegisterType((*NetworkInstancesFdbNetworkInstanceListMacipTableType)(nil), "network_instances_fdb.network_instance_list.macip_table_type")
	proto.RegisterType((*NetworkInstancesFdbNetworkInstanceListMacipTableTypeEntriesType)(nil), "network_instances_fdb.network_instance_list.macip_table_type.entries_type")
	proto.RegisterType((*NetworkInstancesFdbNetworkInstanceListMacipTableTypeEntriesTypeEntryList)(nil), "network_instances_fdb.network_instance_list.macip_table_type.entries_type.entry_list")
	proto.RegisterExtension(E_JnprNetworkInstancesFdbExt)
}

func init() { proto.RegisterFile("l2ald_fdb_render.proto", fileDescriptor_1194cdb0d4efc56b) }

var fileDescriptor_1194cdb0d4efc56b = []byte{
	// 768 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x96, 0x5d, 0x4f, 0x13, 0x4d,
	0x14, 0xc7, 0xb3, 0xbc, 0x3c, 0xb4, 0xa7, 0x85, 0xf6, 0x59, 0x02, 0xec, 0xb3, 0x79, 0x90, 0x8a,
	0x89, 0x69, 0x8c, 0xe9, 0x05, 0xef, 0x82, 0x24, 0x48, 0x82, 0x52, 0xa2, 0x24, 0xb6, 0x25, 0xd1,
	0xab, 0xc9, 0xb0, 0x3b, 0xad, 0x83, 0xbb, 0xb3, 0x9b, 0xd9, 0x69, 0xa5, 0x5c, 0xfa, 0x25, 0x4c,
	0xe4, 0x0a, 0x13, 0xe3, 0x85, 0xd7, 0xde, 0xfa, 0x61, 0xbc, 0xd0, 0xaf, 0x61, 0x66, 0xa6, 0xa5,
	0x2f, 0x2c, 0x0d, 0x6f, 0x77, 0xec, 0xef, 0x9c, 0x3d, 0xff, 0xb3, 0x33, 0xff, 0x3f, 0x29, 0x4c,
	0x7b, 0x0b, 0xd8, 0x73, 0x51, 0xd5, 0x3d, 0x44, 0x9c, 0x30, 0x97, 0xf0, 0x42, 0xc8, 0x03, 0x11,
	0xd8, 0x93, 0x82, 0x78, 0xc4, 0x27, 0x82, 0x37, 0x91, 0x08, 0x42, 0x0d, 0xe7, 0x4f, 0xb3, 0x30,
	0xc5, 0x88, 0xf8, 0x10, 0xf0, 0xf7, 0x88, 0xb2, 0x48, 0x60, 0xe6, 0x90, 0x48, 0xbe, 0x6b, 0xbe,
	0x81, 0x6c, 0x7f, 0xc1, 0xfa, 0x64, 0xe4, 0x86, 0xf3, 0xa9, 0x85, 0xc7, 0x85, 0xd8, 0x57, 0x2e,
	0x50, 0xe4, 0xd1, 0x48, 0x94, 0x32, 0x2d, 0x5c, 0x6c, 0x51, 0xfb, 0x5b, 0xe6, 0xa2, 0xa6, 0x6a,
	0x35, 0xff, 0x83, 0x11, 0x86, 0x7d, 0x62, 0x2d, 0xe6, 0x8c, 0x7c, 0x72, 0x7b, 0xf4, 0xe3, 0xd6,
	0x50, 0xc2, 0x28, 0x29, 0x64, 0xd6, 0x60, 0xc2, 0xc7, 0x0e, 0x12, 0xf8, 0xd0, 0x23, 0x88, 0xb2,
	0x6a, 0x60, 0x7d, 0x37, 0x72, 0x46, 0x3e, 0xb5, 0xb0, 0x75, 0x9d, 0x65, 0x0a, 0xbd, 0x33, 0x90,
	0x68, 0x86, 0xa4, 0x94, 0xf6, 0xb1, 0x53, 0x91, 0xac, 0xc8, 0xaa, 0x81, 0xf9, 0x16, 0x92, 0xe7,
	0x4d, 0xf2, 0x83, 0xa5, 0xc6, 0xc6, 0x0d, 0x35, 0xd4, 0xf8, 0x44, 0x7b, 0xbc, 0xe9, 0x43, 0xd6,
	0xc7, 0x0e, 0x0d, 0xbb, 0xbf, 0xe2, 0x87, 0x56, 0xd8, 0xbe, 0xae, 0x42, 0xcf, 0x14, 0x2d, 0x34,
	0xa1, 0x70, 0xe7, 0x4b, 0x10, 0xa4, 0xba, 0x1a, 0xad, 0x33, 0xad, 0xb4, 0x79, 0x63, 0x25, 0x25,
	0x02, 0x1d, 0x11, 0xfb, 0xd4, 0x80, 0xc9, 0x98, 0x03, 0x35, 0x6d, 0x48, 0x78, 0x04, 0x73, 0x46,
	0x59, 0xcd, 0x7a, 0x91, 0x33, 0xf2, 0x89, 0xd2, 0xf9, 0xb3, 0x39, 0x0b, 0x80, 0x6b, 0x94, 0xd5,
	0x90, 0xa0, 0x3e, 0xb1, 0x76, 0x73, 0x46, 0x7e, 0xbc, 0x94, 0x54, 0xa4, 0x42, 0x7d, 0x22, 0xcb,
	0x7a, 0x5a, 0x44, 0x4f, 0x88, 0x55, 0xd4, 0x65, 0x45, 0xca, 0xf4, 0x84, 0x98, 0x8f, 0xe0, 0x5f,
	0x56, 0xf7, 0x91, 0x17, 0x38, 0xd8, 0x43, 0x84, 0x09, 0x4e, 0x49, 0x64, 0xed, 0xa9, 0xae, 0x0c,
	0xab, 0xfb, 0x2f, 0x25, 0xdf, 0xd1, 0xd8, 0xfe, 0x3c, 0xd2, 0x6d, 0x19, 0xb5, 0x98, 0x03, 0x63,
	0xed, 0x97, 0x5a, 0x37, 0xbb, 0x7b, 0x8b, 0x9b, 0x2d, 0xb4, 0x66, 0xe9, 0x83, 0x69, 0x4f, 0xb6,
	0xbf, 0x0c, 0x43, 0xba, 0xbb, 0x62, 0x1e, 0xc1, 0xa8, 0x7c, 0x6e, 0xb6, 0xe3, 0x53, 0xbe, 0x2b,
	0x4d, 0xf5, 0xd0, 0xd4, 0x29, 0xd3, 0x12, 0xf6, 0xd7, 0x21, 0x80, 0x0e, 0x35, 0xe7, 0x94, 0x05,
	0x10, 0x76, 0x5d, 0x4e, 0xa2, 0x48, 0xe7, 0x4a, 0x5d, 0xe1, 0x33, 0x4d, 0xcc, 0x19, 0x18, 0x6b,
	0x78, 0x98, 0x21, 0xea, 0x5a, 0x4b, 0xea, 0x18, 0xff, 0x91, 0x8f, 0x45, 0xd7, 0xcc, 0xc2, 0x70,
	0x83, 0x51, 0x6b, 0x59, 0x41, 0xf9, 0xa7, 0x39, 0x0f, 0x49, 0xd5, 0xaa, 0x12, 0xba, 0xd2, 0x9d,
	0xd0, 0x84, 0xe4, 0xfb, 0x32, 0xa5, 0xff, 0x43, 0x92, 0x32, 0x41, 0x78, 0x15, 0x3b, 0xc4, 0x5a,
	0x55, 0x6a, 0x1d, 0x20, 0x2f, 0x57, 0xef, 0x26, 0x97, 0xb7, 0xd6, 0x74, 0x59, 0x91, 0x8a, 0x3c,
	0x27, 0x59, 0x6e, 0x10, 0x26, 0x74, 0xf9, 0x49, 0xab, 0x2c, 0x89, 0x2a, 0x3f, 0x84, 0x0c, 0x11,
	0xef, 0x08, 0x67, 0x44, 0x20, 0x81, 0x6b, 0x72, 0xe5, 0x75, 0xb5, 0xdd, 0x78, 0x1b, 0x57, 0x70,
	0xad, 0xe8, 0x9a, 0x0f, 0x60, 0x1c, 0x3b, 0x82, 0x36, 0x08, 0x8a, 0x82, 0x3a, 0x77, 0x88, 0xb5,
	0xa1, 0x26, 0xa5, 0x35, 0x2c, 0x2b, 0x66, 0xff, 0x34, 0x60, 0x2a, 0x36, 0x45, 0x3d, 0xe6, 0x7d,
	0x3d, 0xd0, 0xbc, 0xa5, 0xc1, 0xe6, 0x2d, 0xf7, 0x9b, 0x77, 0x0e, 0x52, 0x21, 0x0f, 0x8e, 0x9b,
	0x48, 0x09, 0x5b, 0x15, 0x35, 0x1c, 0x14, 0x7a, 0x25, 0x49, 0xbc, 0xbb, 0x0f, 0xe2, 0xdd, 0xfd,
	0x6b, 0xa4, 0xf7, 0x9f, 0x89, 0xda, 0x9d, 0x74, 0xfc, 0xdd, 0x4a, 0xfb, 0xde, 0xad, 0xd2, 0x7e,
	0x89, 0xc3, 0x7f, 0xf7, 0x3b, 0xdc, 0x6b, 0x3b, 0xfc, 0x4c, 0x3b, 0xfc, 0xe0, 0xee, 0x54, 0x63,
	0x3c, 0xfe, 0xa7, 0xd7, 0xe3, 0xb3, 0x00, 0x34, 0xec, 0xb3, 0x78, 0x92, 0x86, 0x6d, 0x87, 0xf7,
	0x45, 0x60, 0x69, 0x50, 0x04, 0x96, 0xe3, 0x22, 0xb0, 0x72, 0x49, 0x04, 0x56, 0xaf, 0x10, 0x81,
	0xb5, 0xfe, 0x08, 0xdc, 0x87, 0xb4, 0xb7, 0x88, 0x3a, 0x0d, 0xda, 0xe5, 0x29, 0x6f, 0xb1, 0x78,
	0x49, 0x4a, 0xd6, 0x07, 0xa7, 0x64, 0xe3, 0x0a, 0x29, 0x79, 0x7a, 0xa5, 0x94, 0x6c, 0x5e, 0x4c,
	0xc9, 0x7a, 0x04, 0xf7, 0x8e, 0x58, 0xc8, 0x51, 0xec, 0x6d, 0x22, 0x72, 0x2c, 0xcc, 0x99, 0xc2,
	0x5e, 0x9d, 0xd1, 0x90, 0xf0, 0x7d, 0xdd, 0x12, 0x95, 0x09, 0x8b, 0x02, 0x1e, 0x59, 0x5c, 0xd9,
	0x6f, 0x3a, 0xde, 0x08, 0x25, 0x5b, 0x8e, 0xdd, 0xef, 0xfd, 0x5d, 0x10, 0x3d, 0x77, 0x0f, 0x77,
	0x8e, 0xc5, 0xdf, 0x00, 0x00, 0x00, 0xff, 0xff, 0x5c, 0x0c, 0x97, 0x5b, 0xc0, 0x08, 0x00, 0x00,
}
