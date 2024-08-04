package apis

import (
	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/tools/apis/openweather"
)

var _ core.APIs = &APIs{}

type APIs struct {
	core.OpenWeatherAPI
}

func NewAPIs() *APIs {
	return &APIs{
		OpenWeatherAPI: openweather.NewAPI(),
	}
}
