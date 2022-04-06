// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.4
// source: proto/client/client.proto

package client

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// ClientClient is the client API for Client service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ClientClient interface {
	// Call allows a single request to be made
	Call(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error)
	// Stream is a bidirectional stream
	Stream(ctx context.Context, opts ...grpc.CallOption) (Client_StreamClient, error)
	// Publish publishes a message and returns an empty Message
	Publish(ctx context.Context, in *Message, opts ...grpc.CallOption) (*Message, error)
}

type clientClient struct {
	cc grpc.ClientConnInterface
}

func NewClientClient(cc grpc.ClientConnInterface) ClientClient {
	return &clientClient{cc}
}

func (c *clientClient) Call(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/client.Client/Call", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clientClient) Stream(ctx context.Context, opts ...grpc.CallOption) (Client_StreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &Client_ServiceDesc.Streams[0], "/client.Client/Stream", opts...)
	if err != nil {
		return nil, err
	}
	x := &clientStreamClient{stream}
	return x, nil
}

type Client_StreamClient interface {
	Send(*Request) error
	Recv() (*Response, error)
	grpc.ClientStream
}

type clientStreamClient struct {
	grpc.ClientStream
}

func (x *clientStreamClient) Send(m *Request) error {
	return x.ClientStream.SendMsg(m)
}

func (x *clientStreamClient) Recv() (*Response, error) {
	m := new(Response)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *clientClient) Publish(ctx context.Context, in *Message, opts ...grpc.CallOption) (*Message, error) {
	out := new(Message)
	err := c.cc.Invoke(ctx, "/client.Client/Publish", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ClientServer is the server API for Client service.
// All implementations must embed UnimplementedClientServer
// for forward compatibility
type ClientServer interface {
	// Call allows a single request to be made
	Call(context.Context, *Request) (*Response, error)
	// Stream is a bidirectional stream
	Stream(Client_StreamServer) error
	// Publish publishes a message and returns an empty Message
	Publish(context.Context, *Message) (*Message, error)
	mustEmbedUnimplementedClientServer()
}

// UnimplementedClientServer must be embedded to have forward compatible implementations.
type UnimplementedClientServer struct {
}

func (UnimplementedClientServer) Call(context.Context, *Request) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Call not implemented")
}
func (UnimplementedClientServer) Stream(Client_StreamServer) error {
	return status.Errorf(codes.Unimplemented, "method Stream not implemented")
}
func (UnimplementedClientServer) Publish(context.Context, *Message) (*Message, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Publish not implemented")
}
func (UnimplementedClientServer) mustEmbedUnimplementedClientServer() {}

// UnsafeClientServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ClientServer will
// result in compilation errors.
type UnsafeClientServer interface {
	mustEmbedUnimplementedClientServer()
}

func RegisterClientServer(s grpc.ServiceRegistrar, srv ClientServer) {
	s.RegisterService(&Client_ServiceDesc, srv)
}

func _Client_Call_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClientServer).Call(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/client.Client/Call",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClientServer).Call(ctx, req.(*Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _Client_Stream_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ClientServer).Stream(&clientStreamServer{stream})
}

type Client_StreamServer interface {
	Send(*Response) error
	Recv() (*Request, error)
	grpc.ServerStream
}

type clientStreamServer struct {
	grpc.ServerStream
}

func (x *clientStreamServer) Send(m *Response) error {
	return x.ServerStream.SendMsg(m)
}

func (x *clientStreamServer) Recv() (*Request, error) {
	m := new(Request)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Client_Publish_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Message)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClientServer).Publish(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/client.Client/Publish",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClientServer).Publish(ctx, req.(*Message))
	}
	return interceptor(ctx, in, info, handler)
}

// Client_ServiceDesc is the grpc.ServiceDesc for Client service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Client_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "client.Client",
	HandlerType: (*ClientServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Call",
			Handler:    _Client_Call_Handler,
		},
		{
			MethodName: "Publish",
			Handler:    _Client_Publish_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Stream",
			Handler:       _Client_Stream_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "proto/client/client.proto",
}