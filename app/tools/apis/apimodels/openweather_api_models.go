package apimodels

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*    - OpenWeatherMap API Models -    */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

type (
	GetCurrentWeatherResponse struct {
		Coord      Coord     `json:"coord"`
		Weather    []Weather `json:"weather"`
		Base       string    `json:"base"`
		Main       Main      `json:"main"`
		Visibility int       `json:"visibility"`
		Wind       Wind      `json:"wind"`
		Rain       Rain      `json:"rain"`
		Clouds     Clouds    `json:"clouds"`
		Dt         int       `json:"dt"`
		Sys        Sys       `json:"sys"`
		Timezone   int       `json:"timezone"`
		ID         int       `json:"id"`
		Name       string    `json:"name"`
		Cod        int       `json:"cod"`
	}
	Coord struct {
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
	}
	Weather struct {
		ID          int    `json:"id"`
		Main        string `json:"main"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	}
	Main struct {
		Temp      float64 `json:"temp"`
		FeelsLike float64 `json:"feels_like"`
		TempMin   float64 `json:"temp_min"`
		TempMax   float64 `json:"temp_max"`
		Pressure  int     `json:"pressure"`
		Humidity  int     `json:"humidity"`
	}
	Wind struct {
		Speed float64 `json:"speed"`
		Deg   int     `json:"deg"`
	}
	Rain struct {
		OneH float64 `json:"1h"`
	}
	Clouds struct {
		All int `json:"all"`
	}
	Sys struct {
		Type    int    `json:"type"`
		ID      int    `json:"id"`
		Country string `json:"country"`
		Sunrise int    `json:"sunrise"`
		Sunset  int    `json:"sunset"`
	}
)
