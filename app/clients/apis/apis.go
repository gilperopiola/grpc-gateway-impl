package apis

import (
	"crypto/tls"
	"net/http"
	"time"

	"github.com/gilperopiola/grpc-gateway-impl/app/clients/apis/gpt"
	"github.com/gilperopiola/grpc-gateway-impl/app/clients/apis/weather"
	"github.com/gilperopiola/grpc-gateway-impl/app/core"
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
		Timeout: time.Duration(1) * time.Millisecond / 10000,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
}
