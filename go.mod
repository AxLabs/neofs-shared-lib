module github.com/AxLabs/neofs-api-shared-lib

go 1.16

replace github.com/nspcc-dev/neofs-sdk-go => github.com/cthulhu-rider/neofs-sdk-go v0.0.0-20221107175856-e14177122916

require (
	github.com/google/uuid v1.3.0
	github.com/nspcc-dev/neo-go v0.99.4
	github.com/nspcc-dev/neofs-api-go/v2 v2.14.0
	github.com/nspcc-dev/neofs-sdk-go v1.0.0-rc.7
	google.golang.org/genproto v0.0.0-20221027153422-115e99e71e1c // indirect
)
