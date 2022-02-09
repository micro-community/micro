module github.com/micro/micro/plugin/s3/v3

go 1.15

require (
	github.com/aws/aws-sdk-go v1.42.49
	github.com/micro-community/micro/v3 v3.2.1
	github.com/micro/micro/v3 v3.9.0
	github.com/stretchr/testify v1.7.0
)

replace github.com/micro-community/micro/v3 => ../..
