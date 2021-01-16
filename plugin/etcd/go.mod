module github.com/micro-community/micro/plugin/etcd/v3

go 1.15

require (
	github.com/micro-community/micro/v3 v3.0.4
	github.com/mitchellh/hashstructure v1.0.0
	go.etcd.io/bbolt v1.3.5
	go.uber.org/zap v1.16.0
)

replace github.com/micro-community/micro/v3 => ../..

replace github.com/coreos/bbolt => go.etcd.io/bbolt v1.3.5
