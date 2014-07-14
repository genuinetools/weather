package main

import (
	"encoding/json"
	"fmt"
	// "io/ioutil"
	"net/http"
	"strconv"
)

// response from api.openweathermap.org/data/2.5/weather
// comes back like:
// {
//   "coord": {
//     "lon": -73.95,
//     "lat": 40.78
//   },
//   "sys": {
//     "message": 0.1898,
//     "country": "US",
//     "sunrise": 1405244175,
//     "sunset": 1405297622
//   },
//   "weather": [
//     {
//       "id": 803,
//       "main": "Clouds",
//       "description": "broken clouds",
//       "icon": "04d"
//     }
//   ],
//   "base": "cmc stations",
//   "main": {
//     "temp": 301.32,
//     "pressure": 1015,
//     "humidity": 69,
//     "temp_min": 299.15,
//     "temp_max": 304.15
//   },
//   "wind": {
//     "speed": 1.5,
//     "deg": 0
//   },
//   "rain": {
//     "3h": 0
//   },
//   "snow": {
//     "3h": 0
//   },
//   "clouds": {
//     "all": 75
//   },
//   "dt": 1405278900,
//   "id": 7250946,
//   "name": "Carnegie Hill",
//   "cod": 200
// }

type Forecast struct {
	Base        string      `json:"base"`
	Date        int64       `json:"dt"`
	Clouds      Clouds      `json:"clouds"`
	Coordinates Coordinates `json:coord"`
	Id          int64       `json:"id"`
	Name        string      `json:"name"`
	Rain        Future      `json:"rain"`
	Snow        Future      `json:"rain"`
	Status      int         `json:"cod"`
	Sys         WeatherSys  `json:"sys"`
	Temperature Temperature `json:"main"`
	Weather     []Weather   `json:"weather"`
	Wind        Wind        `json:"wind"`
}

type City struct {
	Id          int64       `json:"id"`
	Name        string      `json:"name"`
	Coordinates Coordinates `json:coord"`
	Country     string      `json:"country"`
	Population  int64       `json:"population"`
}

type Clouds struct {
	All float64 `json:"all"`
}

type Coordinates struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lon"`
}

type DailyForecast struct {
	City      City           `json:"city"`
	Count     int            `json:"cnt"`
	Message   float64        `json:"message"`
	Status    string         `json:"cod"`
	Forecasts []DailyWeather `json:"list"`
}

type DailyWeather struct {
	Clouds      float64     `json:"clouds"`
	Date        int64       `json:"dt"`
	Degrees     float64     `json:"deg"`
	Humidity    float64     `json:"humidity"`
	Pressure    float64     `json:"pressure"`
	Rain        float64     `json:"rain"`
	Speed       float64     `json:"speed"`
	Temperature Temperature `json:"temp"`
	Weather     []Weather   `json:"weather"`
}

type Future struct {
	ThreeHours float64 `json:"3h"`
}

type Temperature struct {
	Humidity       float64 `json:"humidity"`
	Pressure       float64 `json:"pressure"`
	Temperature    float64 `json:"temp"`
	TemperatureMax float64 `json:"temp_max"`
	TemperatureMin float64 `json:"temp_min"`

	Day     float64 `json:"day"`
	Eve     float64 `json:"eve"`
	Max     float64 `json:"max"`
	Min     float64 `json:"min"`
	Morning float64 `json:"morn"`
	Night   float64 `json:"night"`
}

type Weather struct {
	Id          int64  `json:"id"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Main        string `json:"main"`
}

type WeatherSys struct {
	Country string  `json:"country"`
	Message float64 `json:"message"`
	Sunrise int64   `json:"sunrise"`
	Sunset  int64   `json:"sunset"`
}

type Wind struct {
	Degrees float64 `json:"deg"`
	Speed   float64 `json:"speed"`
}

func getWeatherUri(geolocation GeoLocation, units string, days int) (uri string, err error) {
	if units != "metric" && units != "imperial" {
		return uri, fmt.Errorf("%s is not a valid unit. Valid units include metric & imperial.\n", units)
	}

	lat := strconv.FormatFloat(geolocation.Latitude, 'f', 6, 64)
	long := strconv.FormatFloat(geolocation.Longitude, 'f', 6, 64)

	api := "weather"

	if days > 1 {
		api = "forecast/daily"
	}

	uri = "http://api.openweathermap.org/data/2.5/" + api + "?lat=" + lat + "&lon=" + long + "&units=" + units + "&cnt=" + strconv.Itoa(days)

	return uri, nil
}

func getDailyForecast(geolocation GeoLocation, units string, days int) (forecast DailyForecast, err error) {
	url, err := getWeatherUri(geolocation, units, days)
	if err != nil {
		return forecast, err
	}

	resp, err := http.Get(url)
	defer resp.Body.Close()

	// decode the body
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&forecast)
	resp.Body.Close()

	if err != nil {
		return forecast, fmt.Errorf("Decoding the response from %s failed: %s", url, err)
	}

	return forecast, nil
}

func getForecast(geolocation GeoLocation, units string) (forecast Forecast, err error) {
	url, err := getWeatherUri(geolocation, units, 0)
	if err != nil {
		return forecast, err
	}

	resp, err := http.Get(url)
	defer resp.Body.Close()

	// decode the body
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&forecast)
	resp.Body.Close()

	if err != nil {
		return forecast, fmt.Errorf("Decoding the response from %s failed: %s", url, err)
	}

	return forecast, nil
}
