package weather

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gilperopiola/god"
	"github.com/gilperopiola/grpc-gateway-impl/app/clients/apis/apimodels"
	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/shared/logs"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/shared/utils"
)

var _ core.WeatherAPI = &WeatherAPI{}

type WeatherAPI struct {
	httpClient *http.Client
}

func NewAPI(httpClient *http.Client) core.WeatherAPI {
	return &WeatherAPI{httpClient}
}

func (api *WeatherAPI) GetCurrentWeather(ctx god.Ctx, lat, lon float64) (*apimodels.GetWeatherResponse, error) {

	// Prepare URL.
	url := fmt.Sprintf("https://api.weathermap.org/data/2.5/weather?lat=%.2f&lon=%.2f&appid=%s", lat, lon, "f4ecb7e7e30e9c1a3219d1236a63303a")

	// Send request.
	status, respBody, err := utils.GET(ctx, url, nil, "", api.httpClient)
	if err != nil {
		return nil, logs.LogUnexpected(err)
	}

	if status != http.StatusOK {
		return nil, fmt.Errorf("status code %d on GET %s: %s", status, url, respBody)
	}

	var out apimodels.GetWeatherResponse
	if err := json.Unmarshal(respBody, &out); err != nil {
		return nil, fmt.Errorf("error unmarshalling GetWeatherResponse: %w", err)
	}

	return &out, nil
}
