// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v3.21.12
// source: api/proto/cluster-contract.proto

package raft_cluster_v1

import (
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

type LogLeadResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id         string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	JsonString string `protobuf:"bytes,2,opt,name=jsonString,proto3" json:"jsonString,omitempty"`
}

func (x *LogLeadResponse) Reset() {
	*x = LogLeadResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_cluster_contract_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LogLeadResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LogLeadResponse) ProtoMessage() {}

func (x *LogLeadResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_cluster_contract_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LogLeadResponse.ProtoReflect.Descriptor instead.
func (*LogLeadResponse) Descriptor() ([]byte, []int) {
	return file_api_proto_cluster_contract_proto_rawDescGZIP(), []int{0}
}

func (x *LogLeadResponse) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *LogLeadResponse) GetJsonString() string {
	if x != nil {
		return x.JsonString
	}
	return ""
}

type LogLeadRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id         string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	JsonString string `protobuf:"bytes,2,opt,name=jsonString,proto3" json:"jsonString,omitempty"`
}

func (x *LogLeadRequest) Reset() {
	*x = LogLeadRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_cluster_contract_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LogLeadRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LogLeadRequest) ProtoMessage() {}

func (x *LogLeadRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_cluster_contract_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LogLeadRequest.ProtoReflect.Descriptor instead.
func (*LogLeadRequest) Descriptor() ([]byte, []int) {
	return file_api_proto_cluster_contract_proto_rawDescGZIP(), []int{1}
}

func (x *LogLeadRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *LogLeadRequest) GetJsonString() string {
	if x != nil {
		return x.JsonString
	}
	return ""
}

type LogInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// In the Raft protocol, each node must store a term for each entry in its log
	Term       int64  `protobuf:"varint,2,opt,name=term,proto3" json:"term,omitempty"`
	Index      int64  `protobuf:"varint,3,opt,name=index,proto3" json:"index,omitempty"`
	JsonString string `protobuf:"bytes,4,opt,name=jsonString,proto3" json:"jsonString,omitempty"`
}

func (x *LogInfo) Reset() {
	*x = LogInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_cluster_contract_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LogInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LogInfo) ProtoMessage() {}

func (x *LogInfo) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_cluster_contract_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LogInfo.ProtoReflect.Descriptor instead.
func (*LogInfo) Descriptor() ([]byte, []int) {
	return file_api_proto_cluster_contract_proto_rawDescGZIP(), []int{2}
}

func (x *LogInfo) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *LogInfo) GetTerm() int64 {
	if x != nil {
		return x.Term
	}
	return 0
}

func (x *LogInfo) GetIndex() int64 {
	if x != nil {
		return x.Index
	}
	return 0
}

func (x *LogInfo) GetJsonString() string {
	if x != nil {
		return x.JsonString
	}
	return ""
}

type LogAccept struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Term int64 `protobuf:"varint,2,opt,name=term,proto3" json:"term,omitempty"`
}

func (x *LogAccept) Reset() {
	*x = LogAccept{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_cluster_contract_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LogAccept) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LogAccept) ProtoMessage() {}

func (x *LogAccept) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_cluster_contract_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LogAccept.ProtoReflect.Descriptor instead.
func (*LogAccept) Descriptor() ([]byte, []int) {
	return file_api_proto_cluster_contract_proto_rawDescGZIP(), []int{3}
}

func (x *LogAccept) GetTerm() int64 {
	if x != nil {
		return x.Term
	}
	return 0
}

type LeadInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	IdLeader int64 `protobuf:"varint,1,opt,name=idLeader,proto3" json:"idLeader,omitempty"`
	// we send term in all request and responses to controll sync with all nodes in cluster
	// if follower get a req with lower term - he ignored this req and responsed actual term
	Term         int64      `protobuf:"varint,2,opt,name=term,proto3" json:"term,omitempty"`
	Logs         []*LogInfo `protobuf:"bytes,3,rep,name=logs,proto3" json:"logs,omitempty"`
	NeedToUpdate bool       `protobuf:"varint,4,opt,name=needToUpdate,proto3" json:"needToUpdate,omitempty"`
}

func (x *LeadInfo) Reset() {
	*x = LeadInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_cluster_contract_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LeadInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LeadInfo) ProtoMessage() {}

func (x *LeadInfo) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_cluster_contract_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LeadInfo.ProtoReflect.Descriptor instead.
func (*LeadInfo) Descriptor() ([]byte, []int) {
	return file_api_proto_cluster_contract_proto_rawDescGZIP(), []int{4}
}

func (x *LeadInfo) GetIdLeader() int64 {
	if x != nil {
		return x.IdLeader
	}
	return 0
}

func (x *LeadInfo) GetTerm() int64 {
	if x != nil {
		return x.Term
	}
	return 0
}

func (x *LeadInfo) GetLogs() []*LogInfo {
	if x != nil {
		return x.Logs
	}
	return nil
}

func (x *LeadInfo) GetNeedToUpdate() bool {
	if x != nil {
		return x.NeedToUpdate
	}
	return false
}

type LeadAccept struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Term int64 `protobuf:"varint,2,opt,name=term,proto3" json:"term,omitempty"`
}

func (x *LeadAccept) Reset() {
	*x = LeadAccept{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_cluster_contract_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LeadAccept) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LeadAccept) ProtoMessage() {}

func (x *LeadAccept) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_cluster_contract_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LeadAccept.ProtoReflect.Descriptor instead.
func (*LeadAccept) Descriptor() ([]byte, []int) {
	return file_api_proto_cluster_contract_proto_rawDescGZIP(), []int{5}
}

func (x *LeadAccept) GetTerm() int64 {
	if x != nil {
		return x.Term
	}
	return 0
}

type ElectionDecision struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Term int64 `protobuf:"varint,2,opt,name=term,proto3" json:"term,omitempty"`
}

func (x *ElectionDecision) Reset() {
	*x = ElectionDecision{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_cluster_contract_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ElectionDecision) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ElectionDecision) ProtoMessage() {}

func (x *ElectionDecision) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_cluster_contract_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ElectionDecision.ProtoReflect.Descriptor instead.
func (*ElectionDecision) Descriptor() ([]byte, []int) {
	return file_api_proto_cluster_contract_proto_rawDescGZIP(), []int{6}
}

func (x *ElectionDecision) GetTerm() int64 {
	if x != nil {
		return x.Term
	}
	return 0
}

type RequestVoteRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Term         int64 `protobuf:"varint,1,opt,name=term,proto3" json:"term,omitempty"`
	LastLogIndex int64 `protobuf:"varint,2,opt,name=lastLogIndex,proto3" json:"lastLogIndex,omitempty"`
	LastLogTerm  int64 `protobuf:"varint,3,opt,name=lastLogTerm,proto3" json:"lastLogTerm,omitempty"`
	SenderId     int64 `protobuf:"varint,4,opt,name=SenderId,proto3" json:"SenderId,omitempty"`
}

func (x *RequestVoteRequest) Reset() {
	*x = RequestVoteRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_cluster_contract_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RequestVoteRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RequestVoteRequest) ProtoMessage() {}

func (x *RequestVoteRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_cluster_contract_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RequestVoteRequest.ProtoReflect.Descriptor instead.
func (*RequestVoteRequest) Descriptor() ([]byte, []int) {
	return file_api_proto_cluster_contract_proto_rawDescGZIP(), []int{7}
}

func (x *RequestVoteRequest) GetTerm() int64 {
	if x != nil {
		return x.Term
	}
	return 0
}

func (x *RequestVoteRequest) GetLastLogIndex() int64 {
	if x != nil {
		return x.LastLogIndex
	}
	return 0
}

func (x *RequestVoteRequest) GetLastLogTerm() int64 {
	if x != nil {
		return x.LastLogTerm
	}
	return 0
}

func (x *RequestVoteRequest) GetSenderId() int64 {
	if x != nil {
		return x.SenderId
	}
	return 0
}

type RequestVoteResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Term int64 `protobuf:"varint,2,opt,name=term,proto3" json:"term,omitempty"`
}

func (x *RequestVoteResponse) Reset() {
	*x = RequestVoteResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_cluster_contract_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RequestVoteResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RequestVoteResponse) ProtoMessage() {}

func (x *RequestVoteResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_cluster_contract_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RequestVoteResponse.ProtoReflect.Descriptor instead.
func (*RequestVoteResponse) Descriptor() ([]byte, []int) {
	return file_api_proto_cluster_contract_proto_rawDescGZIP(), []int{8}
}

func (x *RequestVoteResponse) GetTerm() int64 {
	if x != nil {
		return x.Term
	}
	return 0
}

type HeartBeatRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Term     int64 `protobuf:"varint,1,opt,name=term,proto3" json:"term,omitempty"`
	LeaderId int64 `protobuf:"varint,2,opt,name=leaderId,proto3" json:"leaderId,omitempty"`
	// index - increment value ([index]logs)
	PrevLogIndex int64 `protobuf:"varint,3,opt,name=prevLogIndex,proto3" json:"prevLogIndex,omitempty"`
	PrevLogTerm  int64 `protobuf:"varint,4,opt,name=prevLogTerm,proto3" json:"prevLogTerm,omitempty"`
}

func (x *HeartBeatRequest) Reset() {
	*x = HeartBeatRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_cluster_contract_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HeartBeatRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HeartBeatRequest) ProtoMessage() {}

func (x *HeartBeatRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_cluster_contract_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HeartBeatRequest.ProtoReflect.Descriptor instead.
func (*HeartBeatRequest) Descriptor() ([]byte, []int) {
	return file_api_proto_cluster_contract_proto_rawDescGZIP(), []int{9}
}

func (x *HeartBeatRequest) GetTerm() int64 {
	if x != nil {
		return x.Term
	}
	return 0
}

func (x *HeartBeatRequest) GetLeaderId() int64 {
	if x != nil {
		return x.LeaderId
	}
	return 0
}

func (x *HeartBeatRequest) GetPrevLogIndex() int64 {
	if x != nil {
		return x.PrevLogIndex
	}
	return 0
}

func (x *HeartBeatRequest) GetPrevLogTerm() int64 {
	if x != nil {
		return x.PrevLogTerm
	}
	return 0
}

type HeartBeatResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Term int64 `protobuf:"varint,1,opt,name=term,proto3" json:"term,omitempty"`
}

func (x *HeartBeatResponse) Reset() {
	*x = HeartBeatResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_cluster_contract_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HeartBeatResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HeartBeatResponse) ProtoMessage() {}

func (x *HeartBeatResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_cluster_contract_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HeartBeatResponse.ProtoReflect.Descriptor instead.
func (*HeartBeatResponse) Descriptor() ([]byte, []int) {
	return file_api_proto_cluster_contract_proto_rawDescGZIP(), []int{10}
}

func (x *HeartBeatResponse) GetTerm() int64 {
	if x != nil {
		return x.Term
	}
	return 0
}

type Empty struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Empty) Reset() {
	*x = Empty{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_cluster_contract_proto_msgTypes[11]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Empty) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Empty) ProtoMessage() {}

func (x *Empty) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_cluster_contract_proto_msgTypes[11]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Empty.ProtoReflect.Descriptor instead.
func (*Empty) Descriptor() ([]byte, []int) {
	return file_api_proto_cluster_contract_proto_rawDescGZIP(), []int{11}
}

var File_api_proto_cluster_contract_proto protoreflect.FileDescriptor

var file_api_proto_cluster_contract_proto_rawDesc = []byte{
	0x0a, 0x20, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x6c, 0x75, 0x73,
	0x74, 0x65, 0x72, 0x2d, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x07, 0x72, 0x61, 0x66, 0x74, 0x5f, 0x76, 0x31, 0x22, 0x41, 0x0a, 0x0f, 0x4c,
	0x6f, 0x67, 0x4c, 0x65, 0x61, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x0e,
	0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1e,
	0x0a, 0x0a, 0x6a, 0x73, 0x6f, 0x6e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0a, 0x6a, 0x73, 0x6f, 0x6e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x22, 0x40,
	0x0a, 0x0e, 0x4c, 0x6f, 0x67, 0x4c, 0x65, 0x61, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64,
	0x12, 0x1e, 0x0a, 0x0a, 0x6a, 0x73, 0x6f, 0x6e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x6a, 0x73, 0x6f, 0x6e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67,
	0x22, 0x63, 0x0a, 0x07, 0x4c, 0x6f, 0x67, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x0e, 0x0a, 0x02, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x74,
	0x65, 0x72, 0x6d, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x74, 0x65, 0x72, 0x6d, 0x12,
	0x14, 0x0a, 0x05, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05,
	0x69, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x1e, 0x0a, 0x0a, 0x6a, 0x73, 0x6f, 0x6e, 0x53, 0x74, 0x72,
	0x69, 0x6e, 0x67, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x6a, 0x73, 0x6f, 0x6e, 0x53,
	0x74, 0x72, 0x69, 0x6e, 0x67, 0x22, 0x1f, 0x0a, 0x09, 0x4c, 0x6f, 0x67, 0x41, 0x63, 0x63, 0x65,
	0x70, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x65, 0x72, 0x6d, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x04, 0x74, 0x65, 0x72, 0x6d, 0x22, 0x84, 0x01, 0x0a, 0x08, 0x4c, 0x65, 0x61, 0x64, 0x49,
	0x6e, 0x66, 0x6f, 0x12, 0x1a, 0x0a, 0x08, 0x69, 0x64, 0x4c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x69, 0x64, 0x4c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x12,
	0x12, 0x0a, 0x04, 0x74, 0x65, 0x72, 0x6d, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x74,
	0x65, 0x72, 0x6d, 0x12, 0x24, 0x0a, 0x04, 0x6c, 0x6f, 0x67, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x10, 0x2e, 0x72, 0x61, 0x66, 0x74, 0x5f, 0x76, 0x31, 0x2e, 0x4c, 0x6f, 0x67, 0x49,
	0x6e, 0x66, 0x6f, 0x52, 0x04, 0x6c, 0x6f, 0x67, 0x73, 0x12, 0x22, 0x0a, 0x0c, 0x6e, 0x65, 0x65,
	0x64, 0x54, 0x6f, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x08, 0x52,
	0x0c, 0x6e, 0x65, 0x65, 0x64, 0x54, 0x6f, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x22, 0x20, 0x0a,
	0x0a, 0x4c, 0x65, 0x61, 0x64, 0x41, 0x63, 0x63, 0x65, 0x70, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x74,
	0x65, 0x72, 0x6d, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x74, 0x65, 0x72, 0x6d, 0x22,
	0x26, 0x0a, 0x10, 0x45, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x44, 0x65, 0x63, 0x69, 0x73,
	0x69, 0x6f, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x65, 0x72, 0x6d, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x04, 0x74, 0x65, 0x72, 0x6d, 0x22, 0x8a, 0x01, 0x0a, 0x12, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x56, 0x6f, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12,
	0x0a, 0x04, 0x74, 0x65, 0x72, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x74, 0x65,
	0x72, 0x6d, 0x12, 0x22, 0x0a, 0x0c, 0x6c, 0x61, 0x73, 0x74, 0x4c, 0x6f, 0x67, 0x49, 0x6e, 0x64,
	0x65, 0x78, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0c, 0x6c, 0x61, 0x73, 0x74, 0x4c, 0x6f,
	0x67, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x20, 0x0a, 0x0b, 0x6c, 0x61, 0x73, 0x74, 0x4c, 0x6f,
	0x67, 0x54, 0x65, 0x72, 0x6d, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0b, 0x6c, 0x61, 0x73,
	0x74, 0x4c, 0x6f, 0x67, 0x54, 0x65, 0x72, 0x6d, 0x12, 0x1a, 0x0a, 0x08, 0x53, 0x65, 0x6e, 0x64,
	0x65, 0x72, 0x49, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x53, 0x65, 0x6e, 0x64,
	0x65, 0x72, 0x49, 0x64, 0x22, 0x29, 0x0a, 0x13, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x56,
	0x6f, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x74,
	0x65, 0x72, 0x6d, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x74, 0x65, 0x72, 0x6d, 0x22,
	0x88, 0x01, 0x0a, 0x10, 0x48, 0x65, 0x61, 0x72, 0x74, 0x42, 0x65, 0x61, 0x74, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x65, 0x72, 0x6d, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x04, 0x74, 0x65, 0x72, 0x6d, 0x12, 0x1a, 0x0a, 0x08, 0x6c, 0x65, 0x61, 0x64,
	0x65, 0x72, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x6c, 0x65, 0x61, 0x64,
	0x65, 0x72, 0x49, 0x64, 0x12, 0x22, 0x0a, 0x0c, 0x70, 0x72, 0x65, 0x76, 0x4c, 0x6f, 0x67, 0x49,
	0x6e, 0x64, 0x65, 0x78, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0c, 0x70, 0x72, 0x65, 0x76,
	0x4c, 0x6f, 0x67, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x20, 0x0a, 0x0b, 0x70, 0x72, 0x65, 0x76,
	0x4c, 0x6f, 0x67, 0x54, 0x65, 0x72, 0x6d, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0b, 0x70,
	0x72, 0x65, 0x76, 0x4c, 0x6f, 0x67, 0x54, 0x65, 0x72, 0x6d, 0x22, 0x27, 0x0a, 0x11, 0x48, 0x65,
	0x61, 0x72, 0x74, 0x42, 0x65, 0x61, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x12, 0x0a, 0x04, 0x74, 0x65, 0x72, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x74,
	0x65, 0x72, 0x6d, 0x22, 0x07, 0x0a, 0x05, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x32, 0xcc, 0x04, 0x0a,
	0x0b, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x4e, 0x6f, 0x64, 0x65, 0x12, 0x34, 0x0a, 0x12,
	0x53, 0x65, 0x74, 0x45, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x69, 0x6d, 0x65, 0x6f,
	0x75, 0x74, 0x12, 0x0e, 0x2e, 0x72, 0x61, 0x66, 0x74, 0x5f, 0x76, 0x31, 0x2e, 0x45, 0x6d, 0x70,
	0x74, 0x79, 0x1a, 0x0e, 0x2e, 0x72, 0x61, 0x66, 0x74, 0x5f, 0x76, 0x31, 0x2e, 0x45, 0x6d, 0x70,
	0x74, 0x79, 0x12, 0x2f, 0x0a, 0x07, 0x4c, 0x6f, 0x61, 0x64, 0x4c, 0x6f, 0x67, 0x12, 0x10, 0x2e,
	0x72, 0x61, 0x66, 0x74, 0x5f, 0x76, 0x31, 0x2e, 0x4c, 0x6f, 0x67, 0x49, 0x6e, 0x66, 0x6f, 0x1a,
	0x12, 0x2e, 0x72, 0x61, 0x66, 0x74, 0x5f, 0x76, 0x31, 0x2e, 0x4c, 0x6f, 0x67, 0x41, 0x63, 0x63,
	0x65, 0x70, 0x74, 0x12, 0x31, 0x0a, 0x09, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x4c, 0x6f, 0x67,
	0x12, 0x10, 0x2e, 0x72, 0x61, 0x66, 0x74, 0x5f, 0x76, 0x31, 0x2e, 0x4c, 0x6f, 0x67, 0x49, 0x6e,
	0x66, 0x6f, 0x1a, 0x12, 0x2e, 0x72, 0x61, 0x66, 0x74, 0x5f, 0x76, 0x31, 0x2e, 0x4c, 0x6f, 0x67,
	0x41, 0x63, 0x63, 0x65, 0x70, 0x74, 0x12, 0x33, 0x0a, 0x09, 0x53, 0x65, 0x74, 0x4c, 0x65, 0x61,
	0x64, 0x65, 0x72, 0x12, 0x11, 0x2e, 0x72, 0x61, 0x66, 0x74, 0x5f, 0x76, 0x31, 0x2e, 0x4c, 0x65,
	0x61, 0x64, 0x49, 0x6e, 0x66, 0x6f, 0x1a, 0x13, 0x2e, 0x72, 0x61, 0x66, 0x74, 0x5f, 0x76, 0x31,
	0x2e, 0x4c, 0x65, 0x61, 0x64, 0x41, 0x63, 0x63, 0x65, 0x70, 0x74, 0x12, 0x48, 0x0a, 0x0f, 0x52,
	0x65, 0x63, 0x69, 0x76, 0x65, 0x48, 0x65, 0x61, 0x72, 0x74, 0x42, 0x65, 0x61, 0x74, 0x12, 0x19,
	0x2e, 0x72, 0x61, 0x66, 0x74, 0x5f, 0x76, 0x31, 0x2e, 0x48, 0x65, 0x61, 0x72, 0x74, 0x42, 0x65,
	0x61, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x72, 0x61, 0x66, 0x74,
	0x5f, 0x76, 0x31, 0x2e, 0x48, 0x65, 0x61, 0x72, 0x74, 0x42, 0x65, 0x61, 0x74, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x31, 0x0a, 0x06, 0x41, 0x70, 0x70, 0x65, 0x6e, 0x64, 0x12,
	0x17, 0x2e, 0x72, 0x61, 0x66, 0x74, 0x5f, 0x76, 0x31, 0x2e, 0x4c, 0x6f, 0x67, 0x4c, 0x65, 0x61,
	0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0e, 0x2e, 0x72, 0x61, 0x66, 0x74, 0x5f,
	0x76, 0x31, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x12, 0x31, 0x0a, 0x06, 0x44, 0x65, 0x6c, 0x65,
	0x74, 0x65, 0x12, 0x17, 0x2e, 0x72, 0x61, 0x66, 0x74, 0x5f, 0x76, 0x31, 0x2e, 0x4c, 0x6f, 0x67,
	0x4c, 0x65, 0x61, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0e, 0x2e, 0x72, 0x61,
	0x66, 0x74, 0x5f, 0x76, 0x31, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x12, 0x38, 0x0a, 0x03, 0x47,
	0x65, 0x74, 0x12, 0x17, 0x2e, 0x72, 0x61, 0x66, 0x74, 0x5f, 0x76, 0x31, 0x2e, 0x4c, 0x6f, 0x67,
	0x4c, 0x65, 0x61, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x18, 0x2e, 0x72, 0x61,
	0x66, 0x74, 0x5f, 0x76, 0x31, 0x2e, 0x4c, 0x6f, 0x67, 0x4c, 0x65, 0x61, 0x64, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3a, 0x0a, 0x0d, 0x53, 0x74, 0x61, 0x72, 0x74, 0x45, 0x6c,
	0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x0e, 0x2e, 0x72, 0x61, 0x66, 0x74, 0x5f, 0x76, 0x31,
	0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x19, 0x2e, 0x72, 0x61, 0x66, 0x74, 0x5f, 0x76, 0x31,
	0x2e, 0x45, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x44, 0x65, 0x63, 0x69, 0x73, 0x69, 0x6f,
	0x6e, 0x12, 0x48, 0x0a, 0x0b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x56, 0x6f, 0x74, 0x65,
	0x12, 0x1b, 0x2e, 0x72, 0x61, 0x66, 0x74, 0x5f, 0x76, 0x31, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x56, 0x6f, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1c, 0x2e,
	0x72, 0x61, 0x66, 0x74, 0x5f, 0x76, 0x31, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x56,
	0x6f, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x1a, 0x5a, 0x18, 0x70,
	0x6b, 0x67, 0x2f, 0x72, 0x61, 0x66, 0x74, 0x2f, 0x72, 0x61, 0x66, 0x74, 0x5f, 0x63, 0x6c, 0x75,
	0x73, 0x74, 0x65, 0x72, 0x5f, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_proto_cluster_contract_proto_rawDescOnce sync.Once
	file_api_proto_cluster_contract_proto_rawDescData = file_api_proto_cluster_contract_proto_rawDesc
)

func file_api_proto_cluster_contract_proto_rawDescGZIP() []byte {
	file_api_proto_cluster_contract_proto_rawDescOnce.Do(func() {
		file_api_proto_cluster_contract_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_proto_cluster_contract_proto_rawDescData)
	})
	return file_api_proto_cluster_contract_proto_rawDescData
}

var file_api_proto_cluster_contract_proto_msgTypes = make([]protoimpl.MessageInfo, 12)
var file_api_proto_cluster_contract_proto_goTypes = []any{
	(*LogLeadResponse)(nil),     // 0: raft_v1.LogLeadResponse
	(*LogLeadRequest)(nil),      // 1: raft_v1.LogLeadRequest
	(*LogInfo)(nil),             // 2: raft_v1.LogInfo
	(*LogAccept)(nil),           // 3: raft_v1.LogAccept
	(*LeadInfo)(nil),            // 4: raft_v1.LeadInfo
	(*LeadAccept)(nil),          // 5: raft_v1.LeadAccept
	(*ElectionDecision)(nil),    // 6: raft_v1.ElectionDecision
	(*RequestVoteRequest)(nil),  // 7: raft_v1.RequestVoteRequest
	(*RequestVoteResponse)(nil), // 8: raft_v1.RequestVoteResponse
	(*HeartBeatRequest)(nil),    // 9: raft_v1.HeartBeatRequest
	(*HeartBeatResponse)(nil),   // 10: raft_v1.HeartBeatResponse
	(*Empty)(nil),               // 11: raft_v1.Empty
}
var file_api_proto_cluster_contract_proto_depIdxs = []int32{
	2,  // 0: raft_v1.LeadInfo.logs:type_name -> raft_v1.LogInfo
	11, // 1: raft_v1.ClusterNode.SetElectionTimeout:input_type -> raft_v1.Empty
	2,  // 2: raft_v1.ClusterNode.LoadLog:input_type -> raft_v1.LogInfo
	2,  // 3: raft_v1.ClusterNode.DeleteLog:input_type -> raft_v1.LogInfo
	4,  // 4: raft_v1.ClusterNode.SetLeader:input_type -> raft_v1.LeadInfo
	9,  // 5: raft_v1.ClusterNode.ReciveHeartBeat:input_type -> raft_v1.HeartBeatRequest
	1,  // 6: raft_v1.ClusterNode.Append:input_type -> raft_v1.LogLeadRequest
	1,  // 7: raft_v1.ClusterNode.Delete:input_type -> raft_v1.LogLeadRequest
	1,  // 8: raft_v1.ClusterNode.Get:input_type -> raft_v1.LogLeadRequest
	11, // 9: raft_v1.ClusterNode.StartElection:input_type -> raft_v1.Empty
	7,  // 10: raft_v1.ClusterNode.RequestVote:input_type -> raft_v1.RequestVoteRequest
	11, // 11: raft_v1.ClusterNode.SetElectionTimeout:output_type -> raft_v1.Empty
	3,  // 12: raft_v1.ClusterNode.LoadLog:output_type -> raft_v1.LogAccept
	3,  // 13: raft_v1.ClusterNode.DeleteLog:output_type -> raft_v1.LogAccept
	5,  // 14: raft_v1.ClusterNode.SetLeader:output_type -> raft_v1.LeadAccept
	10, // 15: raft_v1.ClusterNode.ReciveHeartBeat:output_type -> raft_v1.HeartBeatResponse
	11, // 16: raft_v1.ClusterNode.Append:output_type -> raft_v1.Empty
	11, // 17: raft_v1.ClusterNode.Delete:output_type -> raft_v1.Empty
	0,  // 18: raft_v1.ClusterNode.Get:output_type -> raft_v1.LogLeadResponse
	6,  // 19: raft_v1.ClusterNode.StartElection:output_type -> raft_v1.ElectionDecision
	8,  // 20: raft_v1.ClusterNode.RequestVote:output_type -> raft_v1.RequestVoteResponse
	11, // [11:21] is the sub-list for method output_type
	1,  // [1:11] is the sub-list for method input_type
	1,  // [1:1] is the sub-list for extension type_name
	1,  // [1:1] is the sub-list for extension extendee
	0,  // [0:1] is the sub-list for field type_name
}

func init() { file_api_proto_cluster_contract_proto_init() }
func file_api_proto_cluster_contract_proto_init() {
	if File_api_proto_cluster_contract_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_proto_cluster_contract_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*LogLeadResponse); i {
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
		file_api_proto_cluster_contract_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*LogLeadRequest); i {
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
		file_api_proto_cluster_contract_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*LogInfo); i {
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
		file_api_proto_cluster_contract_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*LogAccept); i {
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
		file_api_proto_cluster_contract_proto_msgTypes[4].Exporter = func(v any, i int) any {
			switch v := v.(*LeadInfo); i {
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
		file_api_proto_cluster_contract_proto_msgTypes[5].Exporter = func(v any, i int) any {
			switch v := v.(*LeadAccept); i {
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
		file_api_proto_cluster_contract_proto_msgTypes[6].Exporter = func(v any, i int) any {
			switch v := v.(*ElectionDecision); i {
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
		file_api_proto_cluster_contract_proto_msgTypes[7].Exporter = func(v any, i int) any {
			switch v := v.(*RequestVoteRequest); i {
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
		file_api_proto_cluster_contract_proto_msgTypes[8].Exporter = func(v any, i int) any {
			switch v := v.(*RequestVoteResponse); i {
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
		file_api_proto_cluster_contract_proto_msgTypes[9].Exporter = func(v any, i int) any {
			switch v := v.(*HeartBeatRequest); i {
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
		file_api_proto_cluster_contract_proto_msgTypes[10].Exporter = func(v any, i int) any {
			switch v := v.(*HeartBeatResponse); i {
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
		file_api_proto_cluster_contract_proto_msgTypes[11].Exporter = func(v any, i int) any {
			switch v := v.(*Empty); i {
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
			RawDescriptor: file_api_proto_cluster_contract_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   12,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_proto_cluster_contract_proto_goTypes,
		DependencyIndexes: file_api_proto_cluster_contract_proto_depIdxs,
		MessageInfos:      file_api_proto_cluster_contract_proto_msgTypes,
	}.Build()
	File_api_proto_cluster_contract_proto = out.File
	file_api_proto_cluster_contract_proto_rawDesc = nil
	file_api_proto_cluster_contract_proto_goTypes = nil
	file_api_proto_cluster_contract_proto_depIdxs = nil
}
