package apis

import "github.com/gilperopiola/grpc-gateway-impl/app/core"

var _ core.InternalAPIs = &InternalAPIs{}

type InternalAPIs struct {
}

func NewInternalAPIs() *InternalAPIs {
	return &InternalAPIs{}
}
