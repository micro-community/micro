package template

var (
	Makefile = `
GOPATH:=$(shell go env GOPATH)
.PHONY: init
init:
	go get -u google.golang.org/protobuf/proto
	go get -u google.golang.org/protobuf/protoc-gen-go
	go install github.com/micro-community/micro/v3/cmd/protoc-gen-micro
	go install github.com/micro-community/micro/v3/cmd/protoc-gen-openapi

.PHONY: api
api:
	protoc --openapi_out=. --proto_path=. proto/{{.Alias}}.proto

.PHONY: proto
proto:
	protoc --proto_path=. --micro_out=. --go_out=:. proto/{{.Alias}}.proto
	
.PHONY: build
build:
	go build -o {{.Alias}} *.go

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: docker
docker:
	docker build . -t {{.Alias}}:latest
`

	GenerateFile = `package main
//go:generate make proto
`
)
