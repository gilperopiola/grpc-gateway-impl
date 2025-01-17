package apis

import (
	"crypto/tls"
	"net/http"
	"time"

	"github.com/gilperopiola/grpc-gateway-impl/app/clients/apis/gpt"
	"github.com/gilperopiola/grpc-gateway-impl/app/clients/apis/weather"
	"github.com/gilperopiola/grpc-gateway-impl/app/core"
)

var _ core.APIClients = &APIClients{}

type APIClients struct {
	core.GPTAPI
	core.WeatherAPI
}

func NewAPIs(cfg *core.APIsCfg) *APIClients {
	gptAPIHTTPClient := newAPIHTTPClient()
	weatherAPIHTTPClient := newAPIHTTPClient()

	return &APIClients{
		gpt.NewAPI(gptAPIHTTPClient, cfg.GPT.APIKey),
		weather.NewAPI(weatherAPIHTTPClient),
	}
}

// All HTTP APIs should create their own HTTP Client with this function.
//
// Do we gain anything by having them be injected from here instead of
// being created in each API itself?
func newAPIHTTPClient() *http.Client {
	return &http.Client{
		Timeout: 90 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
}
