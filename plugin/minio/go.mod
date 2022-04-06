module github.com/micro-community/micro/plugin/minio/v3

go 1.16

require (
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/micro-community/micro/v3 v3.2.1
	github.com/minio/md5-simd v1.1.2 // indirect
	github.com/minio/minio-go/v7 v7.0.21
	github.com/minio/sha256-simd v1.0.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/rs/xid v1.3.0 // indirect
	github.com/stretchr/testify v1.7.0
	golang.org/x/crypto v0.0.0-20220208233918-bba287dce954 // indirect
	gopkg.in/ini.v1 v1.66.3 // indirect
)

replace github.com/micro-community/micro/v3 => ../..
