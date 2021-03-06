// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.24.0-devel
// 	protoc        v3.12.3
// source: rpc/protocol/heartbeat.proto

package protocol

import (
	context "context"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type PingRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AgentTs       int64      `protobuf:"varint,1,opt,name=agentTs,proto3" json:"agentTs,omitempty"`
	AgentIP       string     `protobuf:"bytes,2,opt,name=agentIP,proto3" json:"agentIP,omitempty"`
	AgentID       string     `protobuf:"bytes,3,opt,name=agentID,proto3" json:"agentID,omitempty"`
	ClientVersion string     `protobuf:"bytes,4,opt,name=clientVersion,proto3" json:"clientVersion,omitempty"`
	Agent         *AgentInfo `protobuf:"bytes,5,opt,name=agent,proto3" json:"agent,omitempty"`
}

func (x *PingRequest) Reset() {
	*x = PingRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_protocol_heartbeat_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PingRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PingRequest) ProtoMessage() {}

func (x *PingRequest) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_protocol_heartbeat_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PingRequest.ProtoReflect.Descriptor instead.
func (*PingRequest) Descriptor() ([]byte, []int) {
	return file_rpc_protocol_heartbeat_proto_rawDescGZIP(), []int{0}
}

func (x *PingRequest) GetAgentTs() int64 {
	if x != nil {
		return x.AgentTs
	}
	return 0
}

func (x *PingRequest) GetAgentIP() string {
	if x != nil {
		return x.AgentIP
	}
	return ""
}

func (x *PingRequest) GetAgentID() string {
	if x != nil {
		return x.AgentID
	}
	return ""
}

func (x *PingRequest) GetClientVersion() string {
	if x != nil {
		return x.ClientVersion
	}
	return ""
}

func (x *PingRequest) GetAgent() *AgentInfo {
	if x != nil {
		return x.Agent
	}
	return nil
}

type PongResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ServerTs      int64  `protobuf:"varint,1,opt,name=serverTs,proto3" json:"serverTs,omitempty"`
	ServerVersion string `protobuf:"bytes,2,opt,name=serverVersion,proto3" json:"serverVersion,omitempty"`
}

func (x *PongResponse) Reset() {
	*x = PongResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_protocol_heartbeat_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PongResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PongResponse) ProtoMessage() {}

func (x *PongResponse) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_protocol_heartbeat_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PongResponse.ProtoReflect.Descriptor instead.
func (*PongResponse) Descriptor() ([]byte, []int) {
	return file_rpc_protocol_heartbeat_proto_rawDescGZIP(), []int{1}
}

func (x *PongResponse) GetServerTs() int64 {
	if x != nil {
		return x.ServerTs
	}
	return 0
}

func (x *PongResponse) GetServerVersion() string {
	if x != nil {
		return x.ServerVersion
	}
	return ""
}

type AgentInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Listen        string                  `protobuf:"bytes,1,opt,name=listen,proto3" json:"listen,omitempty"`
	LogPath       string                  `protobuf:"bytes,2,opt,name=logPath,proto3" json:"logPath,omitempty"`
	Host          *AgentInfoHost          `protobuf:"bytes,3,opt,name=host,proto3" json:"host,omitempty"`
	Load          *AgentInfoLoad          `protobuf:"bytes,4,opt,name=load,proto3" json:"load,omitempty"`
	MemorySwap    *AgentInfoMemorySwap    `protobuf:"bytes,5,opt,name=memorySwap,proto3" json:"memorySwap,omitempty"`
	MemoryVirtual *AgentInfoMemoryVirtual `protobuf:"bytes,6,opt,name=memoryVirtual,proto3" json:"memoryVirtual,omitempty"`
}

func (x *AgentInfo) Reset() {
	*x = AgentInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_protocol_heartbeat_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AgentInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AgentInfo) ProtoMessage() {}

func (x *AgentInfo) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_protocol_heartbeat_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AgentInfo.ProtoReflect.Descriptor instead.
func (*AgentInfo) Descriptor() ([]byte, []int) {
	return file_rpc_protocol_heartbeat_proto_rawDescGZIP(), []int{2}
}

func (x *AgentInfo) GetListen() string {
	if x != nil {
		return x.Listen
	}
	return ""
}

func (x *AgentInfo) GetLogPath() string {
	if x != nil {
		return x.LogPath
	}
	return ""
}

func (x *AgentInfo) GetHost() *AgentInfoHost {
	if x != nil {
		return x.Host
	}
	return nil
}

func (x *AgentInfo) GetLoad() *AgentInfoLoad {
	if x != nil {
		return x.Load
	}
	return nil
}

func (x *AgentInfo) GetMemorySwap() *AgentInfoMemorySwap {
	if x != nil {
		return x.MemorySwap
	}
	return nil
}

func (x *AgentInfo) GetMemoryVirtual() *AgentInfoMemoryVirtual {
	if x != nil {
		return x.MemoryVirtual
	}
	return nil
}

type AgentInfoLoad struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Load1  float64 `protobuf:"fixed64,1,opt,name=load1,proto3" json:"load1,omitempty"`
	Load5  float64 `protobuf:"fixed64,2,opt,name=load5,proto3" json:"load5,omitempty"`
	Load15 float64 `protobuf:"fixed64,3,opt,name=load15,proto3" json:"load15,omitempty"`
}

func (x *AgentInfoLoad) Reset() {
	*x = AgentInfoLoad{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_protocol_heartbeat_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AgentInfoLoad) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AgentInfoLoad) ProtoMessage() {}

func (x *AgentInfoLoad) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_protocol_heartbeat_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AgentInfoLoad.ProtoReflect.Descriptor instead.
func (*AgentInfoLoad) Descriptor() ([]byte, []int) {
	return file_rpc_protocol_heartbeat_proto_rawDescGZIP(), []int{3}
}

func (x *AgentInfoLoad) GetLoad1() float64 {
	if x != nil {
		return x.Load1
	}
	return 0
}

func (x *AgentInfoLoad) GetLoad5() float64 {
	if x != nil {
		return x.Load5
	}
	return 0
}

func (x *AgentInfoLoad) GetLoad15() float64 {
	if x != nil {
		return x.Load15
	}
	return 0
}

type AgentInfoHost struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Hostname        string `protobuf:"bytes,1,opt,name=hostname,proto3" json:"hostname,omitempty"`
	Uptime          int64  `protobuf:"varint,2,opt,name=uptime,proto3" json:"uptime,omitempty"`
	BootTime        int64  `protobuf:"varint,3,opt,name=bootTime,proto3" json:"bootTime,omitempty"`
	Procs           int64  `protobuf:"varint,4,opt,name=procs,proto3" json:"procs,omitempty"`
	Os              string `protobuf:"bytes,5,opt,name=os,proto3" json:"os,omitempty"`
	Platform        string `protobuf:"bytes,6,opt,name=platform,proto3" json:"platform,omitempty"`
	PlatformFamily  string `protobuf:"bytes,7,opt,name=platformFamily,proto3" json:"platformFamily,omitempty"`
	PlatformVersion string `protobuf:"bytes,8,opt,name=platformVersion,proto3" json:"platformVersion,omitempty"`
	KernelVersion   string `protobuf:"bytes,9,opt,name=kernelVersion,proto3" json:"kernelVersion,omitempty"`
	KernelArch      string `protobuf:"bytes,10,opt,name=kernelArch,proto3" json:"kernelArch,omitempty"`
}

func (x *AgentInfoHost) Reset() {
	*x = AgentInfoHost{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_protocol_heartbeat_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AgentInfoHost) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AgentInfoHost) ProtoMessage() {}

func (x *AgentInfoHost) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_protocol_heartbeat_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AgentInfoHost.ProtoReflect.Descriptor instead.
func (*AgentInfoHost) Descriptor() ([]byte, []int) {
	return file_rpc_protocol_heartbeat_proto_rawDescGZIP(), []int{4}
}

func (x *AgentInfoHost) GetHostname() string {
	if x != nil {
		return x.Hostname
	}
	return ""
}

func (x *AgentInfoHost) GetUptime() int64 {
	if x != nil {
		return x.Uptime
	}
	return 0
}

func (x *AgentInfoHost) GetBootTime() int64 {
	if x != nil {
		return x.BootTime
	}
	return 0
}

func (x *AgentInfoHost) GetProcs() int64 {
	if x != nil {
		return x.Procs
	}
	return 0
}

func (x *AgentInfoHost) GetOs() string {
	if x != nil {
		return x.Os
	}
	return ""
}

func (x *AgentInfoHost) GetPlatform() string {
	if x != nil {
		return x.Platform
	}
	return ""
}

func (x *AgentInfoHost) GetPlatformFamily() string {
	if x != nil {
		return x.PlatformFamily
	}
	return ""
}

func (x *AgentInfoHost) GetPlatformVersion() string {
	if x != nil {
		return x.PlatformVersion
	}
	return ""
}

func (x *AgentInfoHost) GetKernelVersion() string {
	if x != nil {
		return x.KernelVersion
	}
	return ""
}

func (x *AgentInfoHost) GetKernelArch() string {
	if x != nil {
		return x.KernelArch
	}
	return ""
}

type AgentInfoMemorySwap struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Total       int64   `protobuf:"varint,1,opt,name=total,proto3" json:"total,omitempty"`
	Used        int64   `protobuf:"varint,2,opt,name=used,proto3" json:"used,omitempty"`
	Free        int64   `protobuf:"varint,3,opt,name=free,proto3" json:"free,omitempty"`
	UsedPercent float64 `protobuf:"fixed64,4,opt,name=usedPercent,proto3" json:"usedPercent,omitempty"`
}

func (x *AgentInfoMemorySwap) Reset() {
	*x = AgentInfoMemorySwap{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_protocol_heartbeat_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AgentInfoMemorySwap) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AgentInfoMemorySwap) ProtoMessage() {}

func (x *AgentInfoMemorySwap) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_protocol_heartbeat_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AgentInfoMemorySwap.ProtoReflect.Descriptor instead.
func (*AgentInfoMemorySwap) Descriptor() ([]byte, []int) {
	return file_rpc_protocol_heartbeat_proto_rawDescGZIP(), []int{5}
}

func (x *AgentInfoMemorySwap) GetTotal() int64 {
	if x != nil {
		return x.Total
	}
	return 0
}

func (x *AgentInfoMemorySwap) GetUsed() int64 {
	if x != nil {
		return x.Used
	}
	return 0
}

func (x *AgentInfoMemorySwap) GetFree() int64 {
	if x != nil {
		return x.Free
	}
	return 0
}

func (x *AgentInfoMemorySwap) GetUsedPercent() float64 {
	if x != nil {
		return x.UsedPercent
	}
	return 0
}

type AgentInfoMemoryVirtual struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Total       int64   `protobuf:"varint,1,opt,name=total,proto3" json:"total,omitempty"`
	Available   int64   `protobuf:"varint,2,opt,name=available,proto3" json:"available,omitempty"`
	Used        int64   `protobuf:"varint,3,opt,name=used,proto3" json:"used,omitempty"`
	UsedPercent float64 `protobuf:"fixed64,4,opt,name=usedPercent,proto3" json:"usedPercent,omitempty"`
	Free        int64   `protobuf:"varint,5,opt,name=free,proto3" json:"free,omitempty"`
	Buffers     int64   `protobuf:"varint,6,opt,name=buffers,proto3" json:"buffers,omitempty"`
	Cached      int64   `protobuf:"varint,7,opt,name=cached,proto3" json:"cached,omitempty"`
}

func (x *AgentInfoMemoryVirtual) Reset() {
	*x = AgentInfoMemoryVirtual{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_protocol_heartbeat_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AgentInfoMemoryVirtual) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AgentInfoMemoryVirtual) ProtoMessage() {}

func (x *AgentInfoMemoryVirtual) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_protocol_heartbeat_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AgentInfoMemoryVirtual.ProtoReflect.Descriptor instead.
func (*AgentInfoMemoryVirtual) Descriptor() ([]byte, []int) {
	return file_rpc_protocol_heartbeat_proto_rawDescGZIP(), []int{6}
}

func (x *AgentInfoMemoryVirtual) GetTotal() int64 {
	if x != nil {
		return x.Total
	}
	return 0
}

func (x *AgentInfoMemoryVirtual) GetAvailable() int64 {
	if x != nil {
		return x.Available
	}
	return 0
}

func (x *AgentInfoMemoryVirtual) GetUsed() int64 {
	if x != nil {
		return x.Used
	}
	return 0
}

func (x *AgentInfoMemoryVirtual) GetUsedPercent() float64 {
	if x != nil {
		return x.UsedPercent
	}
	return 0
}

func (x *AgentInfoMemoryVirtual) GetFree() int64 {
	if x != nil {
		return x.Free
	}
	return 0
}

func (x *AgentInfoMemoryVirtual) GetBuffers() int64 {
	if x != nil {
		return x.Buffers
	}
	return 0
}

func (x *AgentInfoMemoryVirtual) GetCached() int64 {
	if x != nil {
		return x.Cached
	}
	return 0
}

var File_rpc_protocol_heartbeat_proto protoreflect.FileDescriptor

var file_rpc_protocol_heartbeat_proto_rawDesc = []byte{
	0x0a, 0x1c, 0x72, 0x70, 0x63, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x2f, 0x68,
	0x65, 0x61, 0x72, 0x74, 0x62, 0x65, 0x61, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x22, 0xac, 0x01, 0x0a, 0x0b, 0x50, 0x69, 0x6e,
	0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x61, 0x67, 0x65, 0x6e,
	0x74, 0x54, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x61, 0x67, 0x65, 0x6e, 0x74,
	0x54, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x49, 0x50, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x07, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x49, 0x50, 0x12, 0x18, 0x0a, 0x07,
	0x61, 0x67, 0x65, 0x6e, 0x74, 0x49, 0x44, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x61,
	0x67, 0x65, 0x6e, 0x74, 0x49, 0x44, 0x12, 0x24, 0x0a, 0x0d, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74,
	0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x63,
	0x6c, 0x69, 0x65, 0x6e, 0x74, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x29, 0x0a, 0x05,
	0x61, 0x67, 0x65, 0x6e, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x2e, 0x41, 0x67, 0x65, 0x6e, 0x74, 0x49, 0x6e, 0x66, 0x6f,
	0x52, 0x05, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x22, 0x50, 0x0a, 0x0c, 0x50, 0x6f, 0x6e, 0x67, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x73, 0x65, 0x72, 0x76, 0x65,
	0x72, 0x54, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x73, 0x65, 0x72, 0x76, 0x65,
	0x72, 0x54, 0x73, 0x12, 0x24, 0x0a, 0x0d, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x56, 0x65, 0x72,
	0x73, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x73, 0x65, 0x72, 0x76,
	0x65, 0x72, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x22, 0x9e, 0x02, 0x0a, 0x09, 0x41, 0x67,
	0x65, 0x6e, 0x74, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x16, 0x0a, 0x06, 0x6c, 0x69, 0x73, 0x74, 0x65,
	0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6c, 0x69, 0x73, 0x74, 0x65, 0x6e, 0x12,
	0x18, 0x0a, 0x07, 0x6c, 0x6f, 0x67, 0x50, 0x61, 0x74, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x07, 0x6c, 0x6f, 0x67, 0x50, 0x61, 0x74, 0x68, 0x12, 0x2b, 0x0a, 0x04, 0x68, 0x6f, 0x73,
	0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63,
	0x6f, 0x6c, 0x2e, 0x41, 0x67, 0x65, 0x6e, 0x74, 0x49, 0x6e, 0x66, 0x6f, 0x48, 0x6f, 0x73, 0x74,
	0x52, 0x04, 0x68, 0x6f, 0x73, 0x74, 0x12, 0x2b, 0x0a, 0x04, 0x6c, 0x6f, 0x61, 0x64, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x2e,
	0x41, 0x67, 0x65, 0x6e, 0x74, 0x49, 0x6e, 0x66, 0x6f, 0x4c, 0x6f, 0x61, 0x64, 0x52, 0x04, 0x6c,
	0x6f, 0x61, 0x64, 0x12, 0x3d, 0x0a, 0x0a, 0x6d, 0x65, 0x6d, 0x6f, 0x72, 0x79, 0x53, 0x77, 0x61,
	0x70, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63,
	0x6f, 0x6c, 0x2e, 0x41, 0x67, 0x65, 0x6e, 0x74, 0x49, 0x6e, 0x66, 0x6f, 0x4d, 0x65, 0x6d, 0x6f,
	0x72, 0x79, 0x53, 0x77, 0x61, 0x70, 0x52, 0x0a, 0x6d, 0x65, 0x6d, 0x6f, 0x72, 0x79, 0x53, 0x77,
	0x61, 0x70, 0x12, 0x46, 0x0a, 0x0d, 0x6d, 0x65, 0x6d, 0x6f, 0x72, 0x79, 0x56, 0x69, 0x72, 0x74,
	0x75, 0x61, 0x6c, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x20, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x63, 0x6f, 0x6c, 0x2e, 0x41, 0x67, 0x65, 0x6e, 0x74, 0x49, 0x6e, 0x66, 0x6f, 0x4d, 0x65,
	0x6d, 0x6f, 0x72, 0x79, 0x56, 0x69, 0x72, 0x74, 0x75, 0x61, 0x6c, 0x52, 0x0d, 0x6d, 0x65, 0x6d,
	0x6f, 0x72, 0x79, 0x56, 0x69, 0x72, 0x74, 0x75, 0x61, 0x6c, 0x22, 0x53, 0x0a, 0x0d, 0x41, 0x67,
	0x65, 0x6e, 0x74, 0x49, 0x6e, 0x66, 0x6f, 0x4c, 0x6f, 0x61, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x6c,
	0x6f, 0x61, 0x64, 0x31, 0x18, 0x01, 0x20, 0x01, 0x28, 0x01, 0x52, 0x05, 0x6c, 0x6f, 0x61, 0x64,
	0x31, 0x12, 0x14, 0x0a, 0x05, 0x6c, 0x6f, 0x61, 0x64, 0x35, 0x18, 0x02, 0x20, 0x01, 0x28, 0x01,
	0x52, 0x05, 0x6c, 0x6f, 0x61, 0x64, 0x35, 0x12, 0x16, 0x0a, 0x06, 0x6c, 0x6f, 0x61, 0x64, 0x31,
	0x35, 0x18, 0x03, 0x20, 0x01, 0x28, 0x01, 0x52, 0x06, 0x6c, 0x6f, 0x61, 0x64, 0x31, 0x35, 0x22,
	0xb9, 0x02, 0x0a, 0x0d, 0x41, 0x67, 0x65, 0x6e, 0x74, 0x49, 0x6e, 0x66, 0x6f, 0x48, 0x6f, 0x73,
	0x74, 0x12, 0x1a, 0x0a, 0x08, 0x68, 0x6f, 0x73, 0x74, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x08, 0x68, 0x6f, 0x73, 0x74, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x16, 0x0a,
	0x06, 0x75, 0x70, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x75,
	0x70, 0x74, 0x69, 0x6d, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x62, 0x6f, 0x6f, 0x74, 0x54, 0x69, 0x6d,
	0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x62, 0x6f, 0x6f, 0x74, 0x54, 0x69, 0x6d,
	0x65, 0x12, 0x14, 0x0a, 0x05, 0x70, 0x72, 0x6f, 0x63, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x05, 0x70, 0x72, 0x6f, 0x63, 0x73, 0x12, 0x0e, 0x0a, 0x02, 0x6f, 0x73, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x02, 0x6f, 0x73, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x6c, 0x61, 0x74, 0x66,
	0x6f, 0x72, 0x6d, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x70, 0x6c, 0x61, 0x74, 0x66,
	0x6f, 0x72, 0x6d, 0x12, 0x26, 0x0a, 0x0e, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x46,
	0x61, 0x6d, 0x69, 0x6c, 0x79, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x70, 0x6c, 0x61,
	0x74, 0x66, 0x6f, 0x72, 0x6d, 0x46, 0x61, 0x6d, 0x69, 0x6c, 0x79, 0x12, 0x28, 0x0a, 0x0f, 0x70,
	0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x08,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x56, 0x65,
	0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x24, 0x0a, 0x0d, 0x6b, 0x65, 0x72, 0x6e, 0x65, 0x6c, 0x56,
	0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x6b, 0x65,
	0x72, 0x6e, 0x65, 0x6c, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x1e, 0x0a, 0x0a, 0x6b,
	0x65, 0x72, 0x6e, 0x65, 0x6c, 0x41, 0x72, 0x63, 0x68, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0a, 0x6b, 0x65, 0x72, 0x6e, 0x65, 0x6c, 0x41, 0x72, 0x63, 0x68, 0x22, 0x75, 0x0a, 0x13, 0x41,
	0x67, 0x65, 0x6e, 0x74, 0x49, 0x6e, 0x66, 0x6f, 0x4d, 0x65, 0x6d, 0x6f, 0x72, 0x79, 0x53, 0x77,
	0x61, 0x70, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x05, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x12, 0x12, 0x0a, 0x04, 0x75, 0x73, 0x65, 0x64,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x75, 0x73, 0x65, 0x64, 0x12, 0x12, 0x0a, 0x04,
	0x66, 0x72, 0x65, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x66, 0x72, 0x65, 0x65,
	0x12, 0x20, 0x0a, 0x0b, 0x75, 0x73, 0x65, 0x64, 0x50, 0x65, 0x72, 0x63, 0x65, 0x6e, 0x74, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x01, 0x52, 0x0b, 0x75, 0x73, 0x65, 0x64, 0x50, 0x65, 0x72, 0x63, 0x65,
	0x6e, 0x74, 0x22, 0xc8, 0x01, 0x0a, 0x16, 0x41, 0x67, 0x65, 0x6e, 0x74, 0x49, 0x6e, 0x66, 0x6f,
	0x4d, 0x65, 0x6d, 0x6f, 0x72, 0x79, 0x56, 0x69, 0x72, 0x74, 0x75, 0x61, 0x6c, 0x12, 0x14, 0x0a,
	0x05, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x74, 0x6f,
	0x74, 0x61, 0x6c, 0x12, 0x1c, 0x0a, 0x09, 0x61, 0x76, 0x61, 0x69, 0x6c, 0x61, 0x62, 0x6c, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x61, 0x76, 0x61, 0x69, 0x6c, 0x61, 0x62, 0x6c,
	0x65, 0x12, 0x12, 0x0a, 0x04, 0x75, 0x73, 0x65, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x04, 0x75, 0x73, 0x65, 0x64, 0x12, 0x20, 0x0a, 0x0b, 0x75, 0x73, 0x65, 0x64, 0x50, 0x65, 0x72,
	0x63, 0x65, 0x6e, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x01, 0x52, 0x0b, 0x75, 0x73, 0x65, 0x64,
	0x50, 0x65, 0x72, 0x63, 0x65, 0x6e, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x66, 0x72, 0x65, 0x65, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x66, 0x72, 0x65, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x62,
	0x75, 0x66, 0x66, 0x65, 0x72, 0x73, 0x18, 0x06, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x62, 0x75,
	0x66, 0x66, 0x65, 0x72, 0x73, 0x12, 0x16, 0x0a, 0x06, 0x63, 0x61, 0x63, 0x68, 0x65, 0x64, 0x18,
	0x07, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x63, 0x61, 0x63, 0x68, 0x65, 0x64, 0x32, 0x44, 0x0a,
	0x09, 0x48, 0x65, 0x61, 0x72, 0x74, 0x62, 0x65, 0x61, 0x74, 0x12, 0x37, 0x0a, 0x04, 0x50, 0x69,
	0x6e, 0x67, 0x12, 0x15, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x2e, 0x50, 0x69,
	0x6e, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x63, 0x6f, 0x6c, 0x2e, 0x50, 0x6f, 0x6e, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x22, 0x00, 0x42, 0x29, 0x0a, 0x19, 0x63, 0x63, 0x2e, 0x61, 0x69, 0x63, 0x6f, 0x64, 0x65,
	0x2e, 0x61, 0x64, 0x61, 0x6e, 0x6f, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c,
	0x5a, 0x0c, 0x72, 0x70, 0x63, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_rpc_protocol_heartbeat_proto_rawDescOnce sync.Once
	file_rpc_protocol_heartbeat_proto_rawDescData = file_rpc_protocol_heartbeat_proto_rawDesc
)

func file_rpc_protocol_heartbeat_proto_rawDescGZIP() []byte {
	file_rpc_protocol_heartbeat_proto_rawDescOnce.Do(func() {
		file_rpc_protocol_heartbeat_proto_rawDescData = protoimpl.X.CompressGZIP(file_rpc_protocol_heartbeat_proto_rawDescData)
	})
	return file_rpc_protocol_heartbeat_proto_rawDescData
}

var file_rpc_protocol_heartbeat_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_rpc_protocol_heartbeat_proto_goTypes = []interface{}{
	(*PingRequest)(nil),            // 0: protocol.PingRequest
	(*PongResponse)(nil),           // 1: protocol.PongResponse
	(*AgentInfo)(nil),              // 2: protocol.AgentInfo
	(*AgentInfoLoad)(nil),          // 3: protocol.AgentInfoLoad
	(*AgentInfoHost)(nil),          // 4: protocol.AgentInfoHost
	(*AgentInfoMemorySwap)(nil),    // 5: protocol.AgentInfoMemorySwap
	(*AgentInfoMemoryVirtual)(nil), // 6: protocol.AgentInfoMemoryVirtual
}
var file_rpc_protocol_heartbeat_proto_depIdxs = []int32{
	2, // 0: protocol.PingRequest.agent:type_name -> protocol.AgentInfo
	4, // 1: protocol.AgentInfo.host:type_name -> protocol.AgentInfoHost
	3, // 2: protocol.AgentInfo.load:type_name -> protocol.AgentInfoLoad
	5, // 3: protocol.AgentInfo.memorySwap:type_name -> protocol.AgentInfoMemorySwap
	6, // 4: protocol.AgentInfo.memoryVirtual:type_name -> protocol.AgentInfoMemoryVirtual
	0, // 5: protocol.Heartbeat.Ping:input_type -> protocol.PingRequest
	1, // 6: protocol.Heartbeat.Ping:output_type -> protocol.PongResponse
	6, // [6:7] is the sub-list for method output_type
	5, // [5:6] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_rpc_protocol_heartbeat_proto_init() }
func file_rpc_protocol_heartbeat_proto_init() {
	if File_rpc_protocol_heartbeat_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_rpc_protocol_heartbeat_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PingRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_rpc_protocol_heartbeat_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PongResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_rpc_protocol_heartbeat_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AgentInfo); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_rpc_protocol_heartbeat_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AgentInfoLoad); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_rpc_protocol_heartbeat_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AgentInfoHost); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_rpc_protocol_heartbeat_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AgentInfoMemorySwap); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_rpc_protocol_heartbeat_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AgentInfoMemoryVirtual); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_rpc_protocol_heartbeat_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_rpc_protocol_heartbeat_proto_goTypes,
		DependencyIndexes: file_rpc_protocol_heartbeat_proto_depIdxs,
		MessageInfos:      file_rpc_protocol_heartbeat_proto_msgTypes,
	}.Build()
	File_rpc_protocol_heartbeat_proto = out.File
	file_rpc_protocol_heartbeat_proto_rawDesc = nil
	file_rpc_protocol_heartbeat_proto_goTypes = nil
	file_rpc_protocol_heartbeat_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// HeartbeatClient is the client API for Heartbeat service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type HeartbeatClient interface {
	Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PongResponse, error)
}

type heartbeatClient struct {
	cc grpc.ClientConnInterface
}

func NewHeartbeatClient(cc grpc.ClientConnInterface) HeartbeatClient {
	return &heartbeatClient{cc}
}

func (c *heartbeatClient) Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PongResponse, error) {
	out := new(PongResponse)
	err := c.cc.Invoke(ctx, "/protocol.Heartbeat/Ping", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// HeartbeatServer is the server API for Heartbeat service.
type HeartbeatServer interface {
	Ping(context.Context, *PingRequest) (*PongResponse, error)
}

// UnimplementedHeartbeatServer can be embedded to have forward compatible implementations.
type UnimplementedHeartbeatServer struct {
}

func (*UnimplementedHeartbeatServer) Ping(context.Context, *PingRequest) (*PongResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}

func RegisterHeartbeatServer(s *grpc.Server, srv HeartbeatServer) {
	s.RegisterService(&_Heartbeat_serviceDesc, srv)
}

func _Heartbeat_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HeartbeatServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protocol.Heartbeat/Ping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HeartbeatServer).Ping(ctx, req.(*PingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Heartbeat_serviceDesc = grpc.ServiceDesc{
	ServiceName: "protocol.Heartbeat",
	HandlerType: (*HeartbeatServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Ping",
			Handler:    _Heartbeat_Ping_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "rpc/protocol/heartbeat.proto",
}
