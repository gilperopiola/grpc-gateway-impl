package api_clients

import "github.com/gilperopiola/grpc-gateway-impl/app/core"

var _ core.ExternalAPIs = (*APIs)(nil)

type APIs struct {
	core.WeatherAPI
}

func NewAPIClients() *APIs {
	return &APIs{
		newWeatherAPIClient(),
	}
}
