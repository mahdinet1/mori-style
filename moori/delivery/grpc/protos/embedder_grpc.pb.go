// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.12.4
// source: embedder.proto

package __

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
	Embedder_ReturnTextVector_FullMethodName  = "/embedder.Embedder/ReturnTextVector"
	Embedder_ReturnImageVector_FullMethodName = "/embedder.Embedder/ReturnImageVector"
)

// EmbedderClient is the client API for Embedder service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type EmbedderClient interface {
	ReturnTextVector(ctx context.Context, in *TextToVectorRequest, opts ...grpc.CallOption) (*TextToVectorReply, error)
	ReturnImageVector(ctx context.Context, in *ImageVectorRequest, opts ...grpc.CallOption) (*ImageVectorReply, error)
}

type embedderClient struct {
	cc grpc.ClientConnInterface
}

func NewEmbedderClient(cc grpc.ClientConnInterface) EmbedderClient {
	return &embedderClient{cc}
}

func (c *embedderClient) ReturnTextVector(ctx context.Context, in *TextToVectorRequest, opts ...grpc.CallOption) (*TextToVectorReply, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(TextToVectorReply)
	err := c.cc.Invoke(ctx, Embedder_ReturnTextVector_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *embedderClient) ReturnImageVector(ctx context.Context, in *ImageVectorRequest, opts ...grpc.CallOption) (*ImageVectorReply, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ImageVectorReply)
	err := c.cc.Invoke(ctx, Embedder_ReturnImageVector_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// EmbedderServer is the server API for Embedder service.
// All implementations must embed UnimplementedEmbedderServer
// for forward compatibility.
type EmbedderServer interface {
	ReturnTextVector(context.Context, *TextToVectorRequest) (*TextToVectorReply, error)
	ReturnImageVector(context.Context, *ImageVectorRequest) (*ImageVectorReply, error)
	mustEmbedUnimplementedEmbedderServer()
}

// UnimplementedEmbedderServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedEmbedderServer struct{}

func (UnimplementedEmbedderServer) ReturnTextVector(context.Context, *TextToVectorRequest) (*TextToVectorReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReturnTextVector not implemented")
}
func (UnimplementedEmbedderServer) ReturnImageVector(context.Context, *ImageVectorRequest) (*ImageVectorReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReturnImageVector not implemented")
}
func (UnimplementedEmbedderServer) mustEmbedUnimplementedEmbedderServer() {}
func (UnimplementedEmbedderServer) testEmbeddedByValue()                  {}

// UnsafeEmbedderServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to EmbedderServer will
// result in compilation errors.
type UnsafeEmbedderServer interface {
	mustEmbedUnimplementedEmbedderServer()
}

func RegisterEmbedderServer(s grpc.ServiceRegistrar, srv EmbedderServer) {
	// If the following call pancis, it indicates UnimplementedEmbedderServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Embedder_ServiceDesc, srv)
}

func _Embedder_ReturnTextVector_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TextToVectorRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EmbedderServer).ReturnTextVector(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Embedder_ReturnTextVector_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EmbedderServer).ReturnTextVector(ctx, req.(*TextToVectorRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Embedder_ReturnImageVector_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ImageVectorRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EmbedderServer).ReturnImageVector(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Embedder_ReturnImageVector_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EmbedderServer).ReturnImageVector(ctx, req.(*ImageVectorRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Embedder_ServiceDesc is the grpc.ServiceDesc for Embedder service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Embedder_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "embedder.Embedder",
	HandlerType: (*EmbedderServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ReturnTextVector",
			Handler:    _Embedder_ReturnTextVector_Handler,
		},
		{
			MethodName: "ReturnImageVector",
			Handler:    _Embedder_ReturnImageVector_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "embedder.proto",
}
