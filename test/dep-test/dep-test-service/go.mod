module dep-test-service

go 1.15

replace dependency => ../

replace github.com/coreos/bbolt => go.etcd.io/bbolt v1.3.5

require (
	dependency v0.0.0-00010101000000-000000000000
	github.com/fsnotify/fsnotify v1.4.9 // indirect
	github.com/golang/protobuf v1.4.3
	github.com/kr/text v0.2.0 // indirect
	github.com/micro-community/micro/v3 v3.0.4
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	google.golang.org/grpc v1.34.0
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
)

replace github.com/micro-community/micro/v3 => ../../..
