module github.com/micro/micro/plugin/prometheus/v3

go 1.15

require (
	github.com/micro-community/micro/v3 v3.0.0-beta.6
	github.com/prometheus/client_golang v1.7.1
	github.com/stretchr/testify v1.6.1
)

replace github.com/micro-community/micro/v3 => ../..
