module github.com/micro-community/micro/plugin/nats/broker/v3

go 1.15

require (
	github.com/micro-community/micro/v3 v3.3.1
	github.com/nats-io/nats.go v1.13.1-0.20220121202836-972a071d373d
	golang.org/x/crypto v0.0.0-20220208233918-bba287dce954 // indirect
	google.golang.org/protobuf v1.27.1
)

replace github.com/micro-community/micro/v3 => ../../..
