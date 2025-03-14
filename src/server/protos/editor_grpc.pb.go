// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.2
// source: protos/editor.proto

package editor

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
	Editor_FetchUpdates_FullMethodName = "/Editor/FetchUpdates"
	Editor_PushOps_FullMethodName      = "/Editor/PushOps"
)

// EditorClient is the client API for Editor service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type EditorClient interface {
	FetchUpdates(ctx context.Context, in *FetchUpdatesReq, opts ...grpc.CallOption) (*FetchUpdatesReply, error)
	PushOps(ctx context.Context, in *Ops, opts ...grpc.CallOption) (*PushOpsReply, error)
}

type editorClient struct {
	cc grpc.ClientConnInterface
}

func NewEditorClient(cc grpc.ClientConnInterface) EditorClient {
	return &editorClient{cc}
}

func (c *editorClient) FetchUpdates(ctx context.Context, in *FetchUpdatesReq, opts ...grpc.CallOption) (*FetchUpdatesReply, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(FetchUpdatesReply)
	err := c.cc.Invoke(ctx, Editor_FetchUpdates_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *editorClient) PushOps(ctx context.Context, in *Ops, opts ...grpc.CallOption) (*PushOpsReply, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(PushOpsReply)
	err := c.cc.Invoke(ctx, Editor_PushOps_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// EditorServer is the server API for Editor service.
// All implementations must embed UnimplementedEditorServer
// for forward compatibility.
type EditorServer interface {
	FetchUpdates(context.Context, *FetchUpdatesReq) (*FetchUpdatesReply, error)
	PushOps(context.Context, *Ops) (*PushOpsReply, error)
	mustEmbedUnimplementedEditorServer()
}

// UnimplementedEditorServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedEditorServer struct{}

func (UnimplementedEditorServer) FetchUpdates(context.Context, *FetchUpdatesReq) (*FetchUpdatesReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FetchUpdates not implemented")
}
func (UnimplementedEditorServer) PushOps(context.Context, *Ops) (*PushOpsReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PushOps not implemented")
}
func (UnimplementedEditorServer) mustEmbedUnimplementedEditorServer() {}
func (UnimplementedEditorServer) testEmbeddedByValue()                {}

// UnsafeEditorServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to EditorServer will
// result in compilation errors.
type UnsafeEditorServer interface {
	mustEmbedUnimplementedEditorServer()
}

func RegisterEditorServer(s grpc.ServiceRegistrar, srv EditorServer) {
	// If the following call pancis, it indicates UnimplementedEditorServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Editor_ServiceDesc, srv)
}

func _Editor_FetchUpdates_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FetchUpdatesReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EditorServer).FetchUpdates(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Editor_FetchUpdates_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EditorServer).FetchUpdates(ctx, req.(*FetchUpdatesReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Editor_PushOps_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Ops)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EditorServer).PushOps(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Editor_PushOps_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EditorServer).PushOps(ctx, req.(*Ops))
	}
	return interceptor(ctx, in, info, handler)
}

// Editor_ServiceDesc is the grpc.ServiceDesc for Editor service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Editor_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "Editor",
	HandlerType: (*EditorServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "FetchUpdates",
			Handler:    _Editor_FetchUpdates_Handler,
		},
		{
			MethodName: "PushOps",
			Handler:    _Editor_PushOps_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "protos/editor.proto",
}
