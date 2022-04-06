module github.com/micro-community/micro/plugin/cockroach/v3

go 1.18

require (
	github.com/lib/pq v1.8.0
	github.com/micro-community/micro/v3 v3.2.1
	github.com/pkg/errors v0.9.1
)

replace github.com/micro-community/micro/v3 => ../..
