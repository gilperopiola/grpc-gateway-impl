package apis

import "github.com/gilperopiola/grpc-gateway-impl/app/core"

var _ core.ExternalAPIs = &ExternalAPIs{}

type ExternalAPIs struct {
	core.OpenWeatherAPI
}

func NewExternalAPIs() *ExternalAPIs {
	return &ExternalAPIs{
		newWeatherAPIClient(),
	}
}
