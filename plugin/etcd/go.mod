module github.com/micro-community/micro/plugin/etcd/v3

go 1.16

require (
	github.com/coreos/go-systemd v0.0.0-20191104093116-d3cd4ed1dbcf // indirect
	github.com/coreos/pkg v0.0.0-20180928190104-399ea9e2e55f // indirect
	github.com/micro-community/micro/v3 v3.2.1
	github.com/mitchellh/hashstructure v1.1.0
	go.etcd.io/bbolt v1.3.6
	go.etcd.io/etcd v0.0.0-20191023171146-3cf2f69b5738
	go.uber.org/zap v1.21.0
	google.golang.org/genproto v0.0.0-20220208230804-65c12eb4c068 // indirect
)

replace github.com/micro-community/micro/v3 => ../..

replace github.com/coreos/bbolt => go.etcd.io/bbolt v1.3.5

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
