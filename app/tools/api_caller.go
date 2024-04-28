package tools

import "github.com/gilperopiola/grpc-gateway-impl/app/core"

type APICaller struct {
	// API Clients go here :)
}

func NewAPICaller() *APICaller {
	return &APICaller{}
}

func (a *APICaller) GetAPICaller() core.APICaller {
	return a
}
