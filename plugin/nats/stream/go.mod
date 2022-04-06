module github.com/micro-community/micro/plugin/nats/stream/v3

go 1.16

require (
	github.com/google/uuid v1.3.0
	github.com/micro-community/micro/v3 v3.3.1
	github.com/nats-io/jwt v1.2.2 // indirect
	github.com/nats-io/nats-streaming-server v0.19.0 // indirect
	github.com/nats-io/nats.go v1.13.0
	github.com/nats-io/stan.go v0.10.2
	github.com/pkg/errors v0.9.1
	golang.org/x/crypto v0.0.0-20220208233918-bba287dce954 // indirect
)

replace github.com/micro-community/micro/v3 => ../../..
