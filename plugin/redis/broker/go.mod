module github.com/micro-community/micro/plugin/redis/broker/v3

go 1.15

require (
	github.com/go-redis/redis/v8 v8.8.3
	github.com/golang/protobuf v1.4.3
	github.com/google/uuid v1.1.2
	github.com/micro-community/micro/v3 v3.2.1
	github.com/oxtoacart/bpool v0.0.0-20190530202638-03653db5a59c
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.7.0
)

replace github.com/micro-community/micro/v3 => ../../..
