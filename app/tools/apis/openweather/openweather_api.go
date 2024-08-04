package openweather

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gilperopiola/god"
	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/tools/apis/apimodels"
)

var _ core.OpenWeatherAPI = &OpenWeatherAPI{}

type OpenWeatherAPI struct{}

func NewAPI() core.OpenWeatherAPI {
	return &OpenWeatherAPI{}
}

func (api *OpenWeatherAPI) GetCurrentWeather(ctx god.Ctx, lat, lon float64) (*apimodels.GetCurrentWeatherResponse, error) {

	// Prepare.
	var url = func(lat, lon float64) string {
		return fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?lat=%f.2&lon=%f.2&appid=f4ecb7e7e30e9c1a3219d1236a63303a", lat, lon)
	}

	// Act.
	resp, err := http.Get(url(lat, lon))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Parse.
	var out apimodels.GetCurrentWeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}

	return &out, nil
}
