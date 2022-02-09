// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: store.proto

package store

import (
	fmt "fmt"
	proto "google.golang.org/protobuf/proto"
	math "math"
)

import (
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

// Api Endpoints for Store service

func NewStoreEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for Store service

type StoreService interface {
	Read(ctx context.Context, in *ReadRequest, opts ...client.CallOption) (*ReadResponse, error)
	Write(ctx context.Context, in *WriteRequest, opts ...client.CallOption) (*WriteResponse, error)
	Delete(ctx context.Context, in *DeleteRequest, opts ...client.CallOption) (*DeleteResponse, error)
	List(ctx context.Context, in *ListRequest, opts ...client.CallOption) (Store_ListService, error)
	Databases(ctx context.Context, in *DatabasesRequest, opts ...client.CallOption) (*DatabasesResponse, error)
	Tables(ctx context.Context, in *TablesRequest, opts ...client.CallOption) (*TablesResponse, error)
}

type storeService struct {
	c    client.Client
	name string
}

func NewStoreService(name string, c client.Client) StoreService {
	return &storeService{
		c:    c,
		name: name,
	}
}

func (c *storeService) Read(ctx context.Context, in *ReadRequest, opts ...client.CallOption) (*ReadResponse, error) {
	req := c.c.NewRequest(c.name, "Store.Read", in)
	out := new(ReadResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *storeService) Write(ctx context.Context, in *WriteRequest, opts ...client.CallOption) (*WriteResponse, error) {
	req := c.c.NewRequest(c.name, "Store.Write", in)
	out := new(WriteResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *storeService) Delete(ctx context.Context, in *DeleteRequest, opts ...client.CallOption) (*DeleteResponse, error) {
	req := c.c.NewRequest(c.name, "Store.Delete", in)
	out := new(DeleteResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *storeService) List(ctx context.Context, in *ListRequest, opts ...client.CallOption) (Store_ListService, error) {
	req := c.c.NewRequest(c.name, "Store.List", &ListRequest{})
	stream, err := c.c.Stream(ctx, req, opts...)
	if err != nil {
		return nil, err
	}
	if err := stream.Send(in); err != nil {
		return nil, err
	}
	return &storeServiceList{stream}, nil
}

type Store_ListService interface {
	Context() context.Context
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
	Recv() (*ListResponse, error)
}

type storeServiceList struct {
	stream client.Stream
}

func (x *storeServiceList) Close() error {
	return x.stream.Close()
}

func (x *storeServiceList) Context() context.Context {
	return x.stream.Context()
}

func (x *storeServiceList) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *storeServiceList) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *storeServiceList) Recv() (*ListResponse, error) {
	m := new(ListResponse)
	err := x.stream.Recv(m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (c *storeService) Databases(ctx context.Context, in *DatabasesRequest, opts ...client.CallOption) (*DatabasesResponse, error) {
	req := c.c.NewRequest(c.name, "Store.Databases", in)
	out := new(DatabasesResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *storeService) Tables(ctx context.Context, in *TablesRequest, opts ...client.CallOption) (*TablesResponse, error) {
	req := c.c.NewRequest(c.name, "Store.Tables", in)
	out := new(TablesResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Store service

type StoreHandler interface {
	Read(context.Context, *ReadRequest, *ReadResponse) error
	Write(context.Context, *WriteRequest, *WriteResponse) error
	Delete(context.Context, *DeleteRequest, *DeleteResponse) error
	List(context.Context, *ListRequest, Store_ListStream) error
	Databases(context.Context, *DatabasesRequest, *DatabasesResponse) error
	Tables(context.Context, *TablesRequest, *TablesResponse) error
}

func RegisterStoreHandler(s server.Server, hdlr StoreHandler, opts ...server.HandlerOption) error {
	type store interface {
		Read(ctx context.Context, in *ReadRequest, out *ReadResponse) error
		Write(ctx context.Context, in *WriteRequest, out *WriteResponse) error
		Delete(ctx context.Context, in *DeleteRequest, out *DeleteResponse) error
		List(ctx context.Context, stream server.Stream) error
		Databases(ctx context.Context, in *DatabasesRequest, out *DatabasesResponse) error
		Tables(ctx context.Context, in *TablesRequest, out *TablesResponse) error
	}
	type Store struct {
		store
	}
	h := &storeHandler{hdlr}
	return s.Handle(s.NewHandler(&Store{h}, opts...))
}

type storeHandler struct {
	StoreHandler
}

func (h *storeHandler) Read(ctx context.Context, in *ReadRequest, out *ReadResponse) error {
	return h.StoreHandler.Read(ctx, in, out)
}

func (h *storeHandler) Write(ctx context.Context, in *WriteRequest, out *WriteResponse) error {
	return h.StoreHandler.Write(ctx, in, out)
}

func (h *storeHandler) Delete(ctx context.Context, in *DeleteRequest, out *DeleteResponse) error {
	return h.StoreHandler.Delete(ctx, in, out)
}

func (h *storeHandler) List(ctx context.Context, stream server.Stream) error {
	m := new(ListRequest)
	if err := stream.Recv(m); err != nil {
		return err
	}
	return h.StoreHandler.List(ctx, m, &storeListStream{stream})
}

type Store_ListStream interface {
	Context() context.Context
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
	Send(*ListResponse) error
}

type storeListStream struct {
	stream server.Stream
}

func (x *storeListStream) Close() error {
	return x.stream.Close()
}

func (x *storeListStream) Context() context.Context {
	return x.stream.Context()
}

func (x *storeListStream) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *storeListStream) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *storeListStream) Send(m *ListResponse) error {
	return x.stream.Send(m)
}

func (h *storeHandler) Databases(ctx context.Context, in *DatabasesRequest, out *DatabasesResponse) error {
	return h.StoreHandler.Databases(ctx, in, out)
}

func (h *storeHandler) Tables(ctx context.Context, in *TablesRequest, out *TablesResponse) error {
	return h.StoreHandler.Tables(ctx, in, out)
}

// Api Endpoints for BlobStore service

func NewBlobStoreEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for BlobStore service

type BlobStoreService interface {
	Read(ctx context.Context, in *BlobReadRequest, opts ...client.CallOption) (BlobStore_ReadService, error)
	Write(ctx context.Context, opts ...client.CallOption) (BlobStore_WriteService, error)
	Delete(ctx context.Context, in *BlobDeleteRequest, opts ...client.CallOption) (*BlobDeleteResponse, error)
}

type blobStoreService struct {
	c    client.Client
	name string
}

func NewBlobStoreService(name string, c client.Client) BlobStoreService {
	return &blobStoreService{
		c:    c,
		name: name,
	}
}

func (c *blobStoreService) Read(ctx context.Context, in *BlobReadRequest, opts ...client.CallOption) (BlobStore_ReadService, error) {
	req := c.c.NewRequest(c.name, "BlobStore.Read", &BlobReadRequest{})
	stream, err := c.c.Stream(ctx, req, opts...)
	if err != nil {
		return nil, err
	}
	if err := stream.Send(in); err != nil {
		return nil, err
	}
	return &blobStoreServiceRead{stream}, nil
}

type BlobStore_ReadService interface {
	Context() context.Context
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
	Recv() (*BlobReadResponse, error)
}

type blobStoreServiceRead struct {
	stream client.Stream
}

func (x *blobStoreServiceRead) Close() error {
	return x.stream.Close()
}

func (x *blobStoreServiceRead) Context() context.Context {
	return x.stream.Context()
}

func (x *blobStoreServiceRead) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *blobStoreServiceRead) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *blobStoreServiceRead) Recv() (*BlobReadResponse, error) {
	m := new(BlobReadResponse)
	err := x.stream.Recv(m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (c *blobStoreService) Write(ctx context.Context, opts ...client.CallOption) (BlobStore_WriteService, error) {
	req := c.c.NewRequest(c.name, "BlobStore.Write", &BlobWriteRequest{})
	stream, err := c.c.Stream(ctx, req, opts...)
	if err != nil {
		return nil, err
	}
	return &blobStoreServiceWrite{stream}, nil
}

type BlobStore_WriteService interface {
	Context() context.Context
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	CloseAndRecv() (*BlobWriteResponse, error)
	Send(*BlobWriteRequest) error
}

type blobStoreServiceWrite struct {
	stream client.Stream
}

func (x *blobStoreServiceWrite) CloseAndRecv() (*BlobWriteResponse, error) {
	if err := x.stream.Close(); err != nil {
		return nil, err
	}
	r := new(BlobWriteResponse)
	err := x.RecvMsg(r)
	return r, err
}

func (x *blobStoreServiceWrite) Context() context.Context {
	return x.stream.Context()
}

func (x *blobStoreServiceWrite) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *blobStoreServiceWrite) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *blobStoreServiceWrite) Send(m *BlobWriteRequest) error {
	return x.stream.Send(m)
}

func (c *blobStoreService) Delete(ctx context.Context, in *BlobDeleteRequest, opts ...client.CallOption) (*BlobDeleteResponse, error) {
	req := c.c.NewRequest(c.name, "BlobStore.Delete", in)
	out := new(BlobDeleteResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for BlobStore service

type BlobStoreHandler interface {
	Read(context.Context, *BlobReadRequest, BlobStore_ReadStream) error
	Write(context.Context, BlobStore_WriteStream) error
	Delete(context.Context, *BlobDeleteRequest, *BlobDeleteResponse) error
}

func RegisterBlobStoreHandler(s server.Server, hdlr BlobStoreHandler, opts ...server.HandlerOption) error {
	type blobStore interface {
		Read(ctx context.Context, stream server.Stream) error
		Write(ctx context.Context, stream server.Stream) error
		Delete(ctx context.Context, in *BlobDeleteRequest, out *BlobDeleteResponse) error
	}
	type BlobStore struct {
		blobStore
	}
	h := &blobStoreHandler{hdlr}
	return s.Handle(s.NewHandler(&BlobStore{h}, opts...))
}

type blobStoreHandler struct {
	BlobStoreHandler
}

func (h *blobStoreHandler) Read(ctx context.Context, stream server.Stream) error {
	m := new(BlobReadRequest)
	if err := stream.Recv(m); err != nil {
		return err
	}
	return h.BlobStoreHandler.Read(ctx, m, &blobStoreReadStream{stream})
}

type BlobStore_ReadStream interface {
	Context() context.Context
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
	Send(*BlobReadResponse) error
}

type blobStoreReadStream struct {
	stream server.Stream
}

func (x *blobStoreReadStream) Close() error {
	return x.stream.Close()
}

func (x *blobStoreReadStream) Context() context.Context {
	return x.stream.Context()
}

func (x *blobStoreReadStream) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *blobStoreReadStream) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *blobStoreReadStream) Send(m *BlobReadResponse) error {
	return x.stream.Send(m)
}

func (h *blobStoreHandler) Write(ctx context.Context, stream server.Stream) error {
	return h.BlobStoreHandler.Write(ctx, &blobStoreWriteStream{stream})
}

type BlobStore_WriteStream interface {
	Context() context.Context
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	SendAndClose(*BlobWriteResponse) error
	Recv() (*BlobWriteRequest, error)
}

type blobStoreWriteStream struct {
	stream server.Stream
}

func (x *blobStoreWriteStream) SendAndClose(in *BlobWriteResponse) error {
	if err := x.SendMsg(in); err != nil {
		return err
	}
	return x.stream.Close()
}

func (x *blobStoreWriteStream) Context() context.Context {
	return x.stream.Context()
}

func (x *blobStoreWriteStream) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *blobStoreWriteStream) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *blobStoreWriteStream) Recv() (*BlobWriteRequest, error) {
	m := new(BlobWriteRequest)
	if err := x.stream.Recv(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (h *blobStoreHandler) Delete(ctx context.Context, in *BlobDeleteRequest, out *BlobDeleteResponse) error {
	return h.BlobStoreHandler.Delete(ctx, in, out)
}
