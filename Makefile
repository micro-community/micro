NAME=micro
IMAGE_NAME=micro-community/$(NAME)
GOPATH:=$(shell go env GOPATH)
GOROOT:=$(shell go env GOROOT)
empty :=
space := $(empty) $(empty)
GIT_COMMIT=$(shell git rev-parse --short HEAD)
GIT_TAG=$(shell git describe --abbrev=0 --tags --always --match "v*")
GIT_IMPORT=github.com/micro-community/micro/v3/cmd
CGO_ENABLED=0
BUILD_DATE=$(shell date +%s)
LDFLAGS=-X $(GIT_IMPORT).BuildDate=$(BUILD_DATE) -X $(GIT_IMPORT).GitCommit=$(GIT_COMMIT) -X $(GIT_IMPORT).GitTag=$(GIT_TAG)
IMAGE_TAG=$(GIT_TAG)-$(GIT_COMMIT)
PROTO_FLAGS=--go_opt=paths=source_relative --go-grpc_opt=paths=source_relative --micro_opt=paths=source_relative


PROTO_PATH=/usr/local/include:.
SRC_DIR=/Users/aiot/codes/micro-community/micro

all: build

.PHONY: api
api:
	find proto/ -name '*.proto' -exec protoc --proto_path=$(PROTO_PATH) --openapi_out=${SRC_DIR} {} \;

vendor:
	go mod vendor

build:
	go build -a -installsuffix cgo -ldflags "-s -w ${LDFLAGS}" -o $(NAME)

docker:
	docker buildx build --platform linux/amd64 --platform linux/arm64 --tag $(IMAGE_NAME):$(IMAGE_TAG) --tag $(IMAGE_NAME):latest --push .

.PHONY: proto
proto:
	find proto/ -name '*.proto' -exec protoc --proto_path=$(PROTO_PATH) $(PROTO_FLAGS) --micro_out=. --go_out=. --go-grpc_out=. {} \;


vet:
	go vet ./...

test: vet
	go test -v ./...

clean:
	rm -rf ./micro

.PHONY: build clean vet test docker
