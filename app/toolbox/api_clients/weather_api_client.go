package api_clients

import (
	"encoding/json"
	"net/http"

	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/other"
)

type weatherAPIClient struct{}

func newWeatherAPIClient() core.WeatherAPI {
	return &weatherAPIClient{}
}

// T0D0 Lat lon
func (cli *weatherAPIClient) GetCurrentWeather(ctx core.Ctx, lat, lon float64) (other.GetWeatherResponse, error) {
	resp, err := http.Get("https://api.openweathermap.org/data/2.5/weather?lat=44.34&lon=10.99&appid=f4ecb7e7e30e9c1a3219d1236a63303a")
	if err != nil {
		return other.GetWeatherResponse{}, err
	}
	defer resp.Body.Close()

	var getWeatherResponse other.GetWeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&getWeatherResponse); err != nil {
		return other.GetWeatherResponse{}, err
	}

	return getWeatherResponse, nil
}
