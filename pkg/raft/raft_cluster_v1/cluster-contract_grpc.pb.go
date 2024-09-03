// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.21.12
// source: api/proto/cluster-contract.proto

package raft_cluster_v1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	ClusterNode_SetElectionTimeout_FullMethodName = "/raft_v1.ClusterNode/SetElectionTimeout"
	ClusterNode_LoadLog_FullMethodName            = "/raft_v1.ClusterNode/LoadLog"
	ClusterNode_SetLeader_FullMethodName          = "/raft_v1.ClusterNode/SetLeader"
	ClusterNode_ReciveHeartBeat_FullMethodName    = "/raft_v1.ClusterNode/ReciveHeartBeat"
	ClusterNode_Append_FullMethodName             = "/raft_v1.ClusterNode/Append"
	ClusterNode_UpdateLogs_FullMethodName         = "/raft_v1.ClusterNode/UpdateLogs"
	ClusterNode_StartElection_FullMethodName      = "/raft_v1.ClusterNode/StartElection"
	ClusterNode_RequestVote_FullMethodName        = "/raft_v1.ClusterNode/RequestVote"
)

// ClusterNodeClient is the client API for ClusterNode service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ClusterNodeClient interface {
	SetElectionTimeout(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error)
	// Follower
	LoadLog(ctx context.Context, in *LogInfo, opts ...grpc.CallOption) (*LogAccept, error)
	SetLeader(ctx context.Context, in *LeadInfo, opts ...grpc.CallOption) (*LeadAccept, error)
	ReciveHeartBeat(ctx context.Context, in *HeartBeatRequest, opts ...grpc.CallOption) (*HeartBeatResponse, error)
	// Lead
	Append(ctx context.Context, in *LogLeadRequest, opts ...grpc.CallOption) (*Empty, error)
	UpdateLogs(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*SyncLog, error)
	// Candidate
	StartElection(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*ElectionDecision, error)
	RequestVote(ctx context.Context, in *RequestVoteRequest, opts ...grpc.CallOption) (*RequestVoteResponse, error)
}

type clusterNodeClient struct {
	cc grpc.ClientConnInterface
}

func NewClusterNodeClient(cc grpc.ClientConnInterface) ClusterNodeClient {
	return &clusterNodeClient{cc}
}

func (c *clusterNodeClient) SetElectionTimeout(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Empty)
	err := c.cc.Invoke(ctx, ClusterNode_SetElectionTimeout_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clusterNodeClient) LoadLog(ctx context.Context, in *LogInfo, opts ...grpc.CallOption) (*LogAccept, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(LogAccept)
	err := c.cc.Invoke(ctx, ClusterNode_LoadLog_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clusterNodeClient) SetLeader(ctx context.Context, in *LeadInfo, opts ...grpc.CallOption) (*LeadAccept, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(LeadAccept)
	err := c.cc.Invoke(ctx, ClusterNode_SetLeader_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clusterNodeClient) ReciveHeartBeat(ctx context.Context, in *HeartBeatRequest, opts ...grpc.CallOption) (*HeartBeatResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(HeartBeatResponse)
	err := c.cc.Invoke(ctx, ClusterNode_ReciveHeartBeat_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clusterNodeClient) Append(ctx context.Context, in *LogLeadRequest, opts ...grpc.CallOption) (*Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Empty)
	err := c.cc.Invoke(ctx, ClusterNode_Append_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clusterNodeClient) UpdateLogs(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*SyncLog, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SyncLog)
	err := c.cc.Invoke(ctx, ClusterNode_UpdateLogs_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clusterNodeClient) StartElection(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*ElectionDecision, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ElectionDecision)
	err := c.cc.Invoke(ctx, ClusterNode_StartElection_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clusterNodeClient) RequestVote(ctx context.Context, in *RequestVoteRequest, opts ...grpc.CallOption) (*RequestVoteResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RequestVoteResponse)
	err := c.cc.Invoke(ctx, ClusterNode_RequestVote_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ClusterNodeServer is the server API for ClusterNode service.
// All implementations must embed UnimplementedClusterNodeServer
// for forward compatibility.
type ClusterNodeServer interface {
	SetElectionTimeout(context.Context, *Empty) (*Empty, error)
	// Follower
	LoadLog(context.Context, *LogInfo) (*LogAccept, error)
	SetLeader(context.Context, *LeadInfo) (*LeadAccept, error)
	ReciveHeartBeat(context.Context, *HeartBeatRequest) (*HeartBeatResponse, error)
	// Lead
	Append(context.Context, *LogLeadRequest) (*Empty, error)
	UpdateLogs(context.Context, *Empty) (*SyncLog, error)
	// Candidate
	StartElection(context.Context, *Empty) (*ElectionDecision, error)
	RequestVote(context.Context, *RequestVoteRequest) (*RequestVoteResponse, error)
	mustEmbedUnimplementedClusterNodeServer()
}

// UnimplementedClusterNodeServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedClusterNodeServer struct{}

func (UnimplementedClusterNodeServer) SetElectionTimeout(context.Context, *Empty) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetElectionTimeout not implemented")
}
func (UnimplementedClusterNodeServer) LoadLog(context.Context, *LogInfo) (*LogAccept, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LoadLog not implemented")
}
func (UnimplementedClusterNodeServer) SetLeader(context.Context, *LeadInfo) (*LeadAccept, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetLeader not implemented")
}
func (UnimplementedClusterNodeServer) ReciveHeartBeat(context.Context, *HeartBeatRequest) (*HeartBeatResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReciveHeartBeat not implemented")
}
func (UnimplementedClusterNodeServer) Append(context.Context, *LogLeadRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Append not implemented")
}
func (UnimplementedClusterNodeServer) UpdateLogs(context.Context, *Empty) (*SyncLog, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateLogs not implemented")
}
func (UnimplementedClusterNodeServer) StartElection(context.Context, *Empty) (*ElectionDecision, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StartElection not implemented")
}
func (UnimplementedClusterNodeServer) RequestVote(context.Context, *RequestVoteRequest) (*RequestVoteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RequestVote not implemented")
}
func (UnimplementedClusterNodeServer) mustEmbedUnimplementedClusterNodeServer() {}
func (UnimplementedClusterNodeServer) testEmbeddedByValue()                     {}

// UnsafeClusterNodeServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ClusterNodeServer will
// result in compilation errors.
type UnsafeClusterNodeServer interface {
	mustEmbedUnimplementedClusterNodeServer()
}

func RegisterClusterNodeServer(s grpc.ServiceRegistrar, srv ClusterNodeServer) {
	// If the following call pancis, it indicates UnimplementedClusterNodeServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&ClusterNode_ServiceDesc, srv)
}

func _ClusterNode_SetElectionTimeout_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClusterNodeServer).SetElectionTimeout(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ClusterNode_SetElectionTimeout_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClusterNodeServer).SetElectionTimeout(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _ClusterNode_LoadLog_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LogInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClusterNodeServer).LoadLog(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ClusterNode_LoadLog_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClusterNodeServer).LoadLog(ctx, req.(*LogInfo))
	}
	return interceptor(ctx, in, info, handler)
}

func _ClusterNode_SetLeader_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LeadInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClusterNodeServer).SetLeader(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ClusterNode_SetLeader_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClusterNodeServer).SetLeader(ctx, req.(*LeadInfo))
	}
	return interceptor(ctx, in, info, handler)
}

func _ClusterNode_ReciveHeartBeat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HeartBeatRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClusterNodeServer).ReciveHeartBeat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ClusterNode_ReciveHeartBeat_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClusterNodeServer).ReciveHeartBeat(ctx, req.(*HeartBeatRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ClusterNode_Append_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LogLeadRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClusterNodeServer).Append(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ClusterNode_Append_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClusterNodeServer).Append(ctx, req.(*LogLeadRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ClusterNode_UpdateLogs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClusterNodeServer).UpdateLogs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ClusterNode_UpdateLogs_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClusterNodeServer).UpdateLogs(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _ClusterNode_StartElection_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClusterNodeServer).StartElection(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ClusterNode_StartElection_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClusterNodeServer).StartElection(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _ClusterNode_RequestVote_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestVoteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClusterNodeServer).RequestVote(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ClusterNode_RequestVote_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClusterNodeServer).RequestVote(ctx, req.(*RequestVoteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ClusterNode_ServiceDesc is the grpc.ServiceDesc for ClusterNode service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ClusterNode_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "raft_v1.ClusterNode",
	HandlerType: (*ClusterNodeServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SetElectionTimeout",
			Handler:    _ClusterNode_SetElectionTimeout_Handler,
		},
		{
			MethodName: "LoadLog",
			Handler:    _ClusterNode_LoadLog_Handler,
		},
		{
			MethodName: "SetLeader",
			Handler:    _ClusterNode_SetLeader_Handler,
		},
		{
			MethodName: "ReciveHeartBeat",
			Handler:    _ClusterNode_ReciveHeartBeat_Handler,
		},
		{
			MethodName: "Append",
			Handler:    _ClusterNode_Append_Handler,
		},
		{
			MethodName: "UpdateLogs",
			Handler:    _ClusterNode_UpdateLogs_Handler,
		},
		{
			MethodName: "StartElection",
			Handler:    _ClusterNode_StartElection_Handler,
		},
		{
			MethodName: "RequestVote",
			Handler:    _ClusterNode_RequestVote_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/proto/cluster-contract.proto",
}
