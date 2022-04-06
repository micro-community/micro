module github.com/micro-community/micro/plugin/prometheus/v3

go 1.16

require (
	github.com/micro-community/micro/v3 v3.2.1
	github.com/prometheus/client_golang v1.12.1
	github.com/stretchr/testify v1.7.0
)

replace github.com/micro-community/micro/v3 => ../..
