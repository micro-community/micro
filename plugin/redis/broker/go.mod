module github.com/micro-community/micro/plugin/redis/broker/v3

go 1.18

require (
	github.com/go-redis/redis/v8 v8.11.4
	github.com/google/uuid v1.3.0
	github.com/micro-community/micro/v3 v3.3.1
	github.com/oxtoacart/bpool v0.0.0-20190530202638-03653db5a59c
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.7.0
	google.golang.org/protobuf v1.27.1
)

replace github.com/micro-community/micro/v3 => ../../..
