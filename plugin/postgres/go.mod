module github.com/micro-community/micro/plugin/postgres/v3

go 1.18

require (
	github.com/lib/pq v1.10.4
	github.com/micro-community/micro/v3 v3.2.1
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.7.0
)

replace github.com/micro-community/micro/v3 => ../..
