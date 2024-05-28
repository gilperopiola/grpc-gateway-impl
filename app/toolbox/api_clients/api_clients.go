package api_clients

import "github.com/gilperopiola/grpc-gateway-impl/app/core"

var _ core.InternalAPIs = (*APIClients)(nil)
var _ core.ExternalAPIs = (*APIClients)(nil)

type APIClients struct {
	core.WeatherAPI
}

func NewAPIClients() *APIClients {
	return &APIClients{
		newWeatherAPIClient(),
	}
}
