package api_clients

import (
	"encoding/json"
	"net/http"

	"github.com/gilperopiola/god"
	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/models"
)

type weatherAPIClient struct{}

func newWeatherAPIClient() core.WeatherAPI {
	return &weatherAPIClient{}
}

// T0D0 Lat lon
func (cli *weatherAPIClient) GetCurrentWeather(ctx god.Ctx, lat, lon float64) (models.GetWeatherResponse, error) {
	resp, err := http.Get("https://api.openweathermap.org/data/2.5/weather?lat=44.34&lon=10.99&appid=f4ecb7e7e30e9c1a3219d1236a63303a")
	if err != nil {
		return models.GetWeatherResponse{}, err
	}
	defer resp.Body.Close()

	var getWeatherResponse models.GetWeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&getWeatherResponse); err != nil {
		return models.GetWeatherResponse{}, err
	}

	return getWeatherResponse, nil
}
