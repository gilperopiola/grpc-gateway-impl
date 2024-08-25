package apis

import (
	"crypto/tls"
	"net/http"

	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/tools/apis/gpt"
	"github.com/gilperopiola/grpc-gateway-impl/app/tools/apis/weather"
)

var _ core.APIs = &APIs{}

type APIs struct {
	core.ChatGPTAPI
	core.WeatherAPI
}

func NewAPIs(cfg core.APIsCfg) *APIs {

	// OpenAI GPT API
	gptHTTPClient := newAPIHTTPClient()
	gptAPI := gpt.NewAPI(gptHTTPClient, cfg.GPT.APIKey)

	// OpenWeatherMap API
	weatherHTTPClient := newAPIHTTPClient()
	weatherAPI := weather.NewAPI(weatherHTTPClient)

	return &APIs{
		gptAPI,
		weatherAPI,
	}
}

// All HTTP APIs should create their own HTTP Client with this function.
//
// Do we gain anything by having them be injected from here instead of
// being created in each API itself?
func newAPIHTTPClient() *http.Client {
	return &http.Client{
		Timeout: 120,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
}
