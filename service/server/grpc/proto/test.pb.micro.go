// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: server/grpc/proto/test.proto

package test

import (
	fmt "fmt"
	math "math"

	proto "google.golang.org/protobuf/proto"
	_ "google.golang.org/genproto/googleapis/api/annotations"

	context "context"

	api "github.com/micro-community/micro/v3/service/api"

	client "github.com/micro-community/micro/v3/service/client"

	server "github.com/micro-community/micro/v3/service/server"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// Reference imports to suppress errors if they are not otherwise used.
var _ api.Endpoint
var _ context.Context
var _ client.Option
var _ server.Option

// Api Endpoints for Test service

func NewTestEndpoints() []*api.Endpoint {
	return []*api.Endpoint{
		&api.Endpoint{
			Name:    "Test.Call",
			Path:    []string{"/api/v0/test/call/{uuid}"},
			Method:  []string{"POST"},
			Body:    "*",
			Handler: "rpc",
		},
		&api.Endpoint{
			Name:    "Test.CallPcre",
			Path:    []string{"^/api/v0/test/call/pcre/?$"},
			Method:  []string{"POST"},
			Body:    "*",
			Handler: "rpc",
		},
		&api.Endpoint{
			Name:    "Test.CallPcreInvalid",
			Path:    []string{"^/api/v0/test/call/pcre/invalid/?"},
			Method:  []string{"POST"},
			Body:    "*",
			Handler: "rpc",
		},
	}
}

// Client API for Test service

type TestService interface {
	Call(ctx context.Context, in *Request, opts ...client.CallOption) (*Response, error)
	CallPcre(ctx context.Context, in *Request, opts ...client.CallOption) (*Response, error)
	CallPcreInvalid(ctx context.Context, in *Request, opts ...client.CallOption) (*Response, error)
}

type testService struct {
	c    client.Client
	name string
}

func NewTestService(name string, c client.Client) TestService {
	return &testService{
		c:    c,
		name: name,
	}
}

func (c *testService) Call(ctx context.Context, in *Request, opts ...client.CallOption) (*Response, error) {
	req := c.c.NewRequest(c.name, "Test.Call", in)
	out := new(Response)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *testService) CallPcre(ctx context.Context, in *Request, opts ...client.CallOption) (*Response, error) {
	req := c.c.NewRequest(c.name, "Test.CallPcre", in)
	out := new(Response)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *testService) CallPcreInvalid(ctx context.Context, in *Request, opts ...client.CallOption) (*Response, error) {
	req := c.c.NewRequest(c.name, "Test.CallPcreInvalid", in)
	out := new(Response)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Test service

type TestHandler interface {
	Call(context.Context, *Request, *Response) error
	CallPcre(context.Context, *Request, *Response) error
	CallPcreInvalid(context.Context, *Request, *Response) error
}

func RegisterTestHandler(s server.Server, hdlr TestHandler, opts ...server.HandlerOption) error {
	type test interface {
		Call(ctx context.Context, in *Request, out *Response) error
		CallPcre(ctx context.Context, in *Request, out *Response) error
		CallPcreInvalid(ctx context.Context, in *Request, out *Response) error
	}
	type Test struct {
		test
	}
	h := &testHandler{hdlr}
	opts = append(opts, api.WithEndpoint(&api.Endpoint{
		Name:    "Test.Call",
		Path:    []string{"/api/v0/test/call/{uuid}"},
		Method:  []string{"POST"},
		Body:    "*",
		Handler: "rpc",
	}))
	opts = append(opts, api.WithEndpoint(&api.Endpoint{
		Name:    "Test.CallPcre",
		Path:    []string{"^/api/v0/test/call/pcre/?$"},
		Method:  []string{"POST"},
		Body:    "*",
		Handler: "rpc",
	}))
	opts = append(opts, api.WithEndpoint(&api.Endpoint{
		Name:    "Test.CallPcreInvalid",
		Path:    []string{"^/api/v0/test/call/pcre/invalid/?"},
		Method:  []string{"POST"},
		Body:    "*",
		Handler: "rpc",
	}))
	return s.Handle(s.NewHandler(&Test{h}, opts...))
}

type testHandler struct {
	TestHandler
}

func (h *testHandler) Call(ctx context.Context, in *Request, out *Response) error {
	return h.TestHandler.Call(ctx, in, out)
}

func (h *testHandler) CallPcre(ctx context.Context, in *Request, out *Response) error {
	return h.TestHandler.CallPcre(ctx, in, out)
}

func (h *testHandler) CallPcreInvalid(ctx context.Context, in *Request, out *Response) error {
	return h.TestHandler.CallPcreInvalid(ctx, in, out)
}
