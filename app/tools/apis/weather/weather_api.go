package weather

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gilperopiola/god"
	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/tools/apis/apimodels"
)

var _ core.WeatherAPI = &WeatherAPI{}

type WeatherAPI struct {
	httpClient *http.Client
	getFn      func(ctx context.Context, client *http.Client, url string, urlParams map[string]string) (int, []byte, error)
}

func NewAPI(getFn func(ctx context.Context, client *http.Client, url string, urlParams map[string]string) (int, []byte, error)) core.WeatherAPI {
	return &WeatherAPI{
		&http.Client{Timeout: 90},
		getFn,
	}
}

func (api *WeatherAPI) GetCurrentWeather(ctx god.Ctx, lat, lon float64) (*apimodels.GetWeatherResponse, error) {

	// Prepare.
	url := fmt.Sprintf("/weather?lat=%.2f&lon=%.2f&appid=%s", lat, lon, "f4ecb7e7e30e9c1a3219d1236a63303a")

	// Act.
	status, respBody, err := api.getFn(ctx, api.httpClient, url, nil)
	if err != nil {
		return nil, core.LogUnexpected(fmt.Errorf("error on GET %s: %w", url, err))
	}

	if status != http.StatusOK {
		return nil, fmt.Errorf("status code %d on GET %s: %s", status, url, respBody)
	}

	var out apimodels.GetWeatherResponse
	if err := json.Unmarshal(respBody, &out); err != nil {
		return nil, err
	}

	return &out, nil
}
