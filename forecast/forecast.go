package forecast

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// response from https://api.forecast.io/forecast/
// comes back like:
// {
//     "alerts": [
//         {
//             "description": "...BLOWING DUST ADVISORY REMAINS IN EFFECT UNTIL 11 PM MST THIS\nEVENING...\nA BLOWING DUST ADVISORY REMAINS IN EFFECT UNTIL 11 PM MST THIS\nEVENING.\n* AFFECTED AREA...NORTHWEST AND NORTH-CENTRAL PINAL COUNTY...\nEXTENDING NORTH INTO EAST CENTRAL MARICOPA COUNTY...INCLUDING\nCHANDLER...GILBERT...MESA...TEMPE...AND THE GREATER PHOENIX\nAREA\n* TIMING...UNTIL 11 PM.\n* WINDS...SOUTH-SOUTHEAST 30 MPH WITH GUSTS UP TO 45 MPH\nPOSSIBLE.\n* VISIBILITY...3 MILES IN BLOWING DUST...LOCALLY DOWN TO AS LOW\nAS ONE MILE.\n* IMPACTS...SUDDEN DROPS IN VISIBILITIES DUE TO BLOWING DUST\nWILL CREATE HAZARDOUS TRAVEL CONDITIONS...ESPECIALLY AT NIGHT.\nTRAVELERS NEED TO BE READY FOR RAPIDLY CHANGING ROAD\nCONDITIONS.\n",
//             "expires": 1405317600,
//             "time": 1405313940,
//             "title": "Blowing Dust Advisory for Maricopa, AZ",
//             "uri": "http://alerts.weather.gov/cap/wwacapget.php?x=AZ1251601B31CC.BlowingDustAdvisory.1251601B68E0AZ.PSRNPWPSR.5957bb654b0ef12b8351fd6d518e951e"
//         }
//     ],
//     "currently": {
//         "apparentTemperature": 91.72,
//         "cloudCover": 0.92,
//         "dewPoint": 64.97,
//         "humidity": 0.44,
//         "icon": "partly-cloudy-night",
//         "nearestStormDistance": 0,
//         "nearestStormBearing": 0,
//         "ozone": 295.8,
//         "precipIntensity": 0,
//         "precipProbability": 0,
//         "pressure": 1010.85,
//         "summary": "Mostly Cloudy",
//         "temperature": 89.72,
//         "time": 1405315610,
//         "visibility": 9.99,
//         "windBearing": 208,
//         "windSpeed": 6.48
//     },
//     "daily": {
//         "data": [
//             {
//                 "apparentTemperatureMax": 107.25,
//                 "apparentTemperatureMaxTime": 1405288800,
//                 "apparentTemperatureMin": 86.87,
//                 "apparentTemperatureMinTime": 1405252800,
//                 "cloudCover": 0.68,
//                 "dewPoint": 62.06,
//                 "humidity": 0.35,
//                 "icon": "rain",
//                 "moonPhase": 0.55,
//                 "ozone": 299.93,
//                 "precipIntensity": 0.0047,
//                 "precipIntensityMax": 0.0308,
//                 "precipIntensityMaxTime": 1405310400,
//                 "precipProbability": 0.95,
//                 "precipType": "rain",
//                 "pressure": 1009.52,
//                 "summary": "Light rain starting in the afternoon.",
//                 "sunriseTime": 1405254533,
//                 "sunsetTime": 1405305631,
//                 "temperatureMax": 105.3,
//                 "temperatureMaxTime": 1405288800,
//                 "temperatureMin": 86.14,
//                 "temperatureMinTime": 1405252800,
//                 "time": 1405234800,
//                 "visibility": 9.88,
//                 "windBearing": 149,
//                 "windSpeed": 2.95
//             }...
//         ],
//         "icon": "rain",
//         "summary": "Light rain today and tomorrow, with temperatures rising to 109°F on Thursday."
//     },
//     "flags": {
//         "units": "us"
//     },
//     "hourly": {
//         "data": [
//             {
//                 "apparentTemperature": 92.39,
//                 "cloudCover": 0.91,
//                 "dewPoint": 64.12,
//                 "humidity": 0.41,
//                 "icon": "rain",
//                 "ozone": 295.87,
//                 "precipIntensity": 0.0178,
//                 "precipProbability": 0.52,
//                 "precipType": "rain",
//                 "pressure": 1010.44,
//                 "summary": "Light Rain",
//                 "temperature": 90.77,
//                 "time": 1405314000,
//                 "visibility": 10,
//                 "windBearing": 195,
//                 "windSpeed": 6.06
//             }...
//         ],
//         "icon": "rain",
//         "summary": "Drizzle starting tomorrow afternoon, continuing until tomorrow evening."
//     },
//     "latitude": 33.4962205,
//     "longitude": -111.9641728,
//     "minutely": {
//         "data": [
//             {
//                 "precipIntensity": 0,
//                 "precipProbability": 0,
//                 "time": 1405315920
//             },
//             {
//                 "precipIntensity": 0.0027,
//                 "precipIntensityError": 0.0002,
//                 "precipProbability": 0.01,
//                 "precipType": "rain",
//                 "time": 1405315980
//             }...
//         ],
//         "icon": "partly-cloudy-night",
//         "summary": "Mostly cloudy for the hour."
//     },
//     "offset": -7,
//     "timezone": "America/Phoenix"
// }

// Forecast contains the information returned from the server
// when requesting the forecast.
type Forecast struct {
	Alerts    []Alert       `json:"alerts"`
	Currently Weather       `json:"currently"`
	Code      int           `json:"code"`
	Daily     TimeDelimited `json:"daily"`
	Error     string        `json:"error"`
	Flags     Flags         `json:"flags"`
	Hourly    TimeDelimited `json:"hourly"`
	Latitude  float64       `json:"latitude"`
	Longitude float64       `json:"longitude"`
	Offset    float64       `json:"offset"`
	Timezone  string        `json:"timezone"`
}

// Alert contains any weather alerts happening at the location.
type Alert struct {
	Description string `json:"description"`
	Expires     int64  `json:"expires"`
	Time        int64  `json:"time"`
	Title       string `json:"title"`
	URI         string `json:"uri"`
}

// Flags describes the flags on a forecast.
type Flags struct {
	Units string `json:"units"`
}

// Weather describes details about the weather for the location.
type Weather struct {
	ApparentTemperature        float64 `json:"apparentTemperature"`
	ApparentTemperatureMax     float64 `json:"apparentTemperatureMax"`
	ApparentTemperatureMaxTime int64   `json:"apparentTemperatureMaxTime"`
	ApparentTemperatureMin     float64 `json:"apparentTemperatureMin"`
	ApparentTemperatureMinTime int64   `json:"apparentTemperatureMinTime"`
	CloudCover                 float64 `json:"cloudCover"`
	DewPoint                   float64 `json:"dewPoint"`
	Humidity                   float64 `json:"humidity"`
	Icon                       string  `json:"icon"`
	NearestStormDistance       float64 `json:"nearestStormDistance"`
	NearestStormBearing        float64 `json:"nearestStormBearing"`
	Ozone                      float64 `json:"ozone"`
	PrecipIntensity            float64 `json:"precipIntensity"`
	PrecipIntensityMax         float64 `json:"precipIntensityMax"`
	PrecipIntensityMaxTime     int64   `json:"precipIntensityMaxTime"`
	PrecipProbability          float64 `json:"precipProbability"`
	PrecipType                 string  `json:"precipType"`
	Pressure                   float64 `json:"pressure"`
	Summary                    string  `json:"summary"`
	SunriseTime                int64   `json:"sunriseTime"`
	SunsetTime                 int64   `json:"sunsetTime"`
	Temperature                float64 `json:"temperature"`
	TemperatureMax             float64 `json:"temperatureMax"`
	TemperatureMaxTime         int64   `json:"temperatureMaxTime"`
	TemperatureMin             float64 `json:"temperatureMin"`
	TemperatureMinTime         int64   `json:"temperatureMinTime"`
	Time                       int64   `json:"time"`
	Visibility                 float64 `json:"visibility"`
	WindBearing                float64 `json:"windBearing"`
	WindSpeed                  float64 `json:"windSpeed"`
}

// TimeDelimited describes the data for the time series.
type TimeDelimited struct {
	Data    []Weather `json:"data"`
	Icon    string    `json:"icon"`
	Summary string    `json:"summary"`
}

// Request describes the request posted to the forecast api.
type Request struct {
	Latitude  float64  `json:"lat"`
	Longitude float64  `json:"lng"`
	Units     string   `json:"units"`
	Exclude   []string `json:"exclude"`
}

// Get performs a request to get the forecast data for a location.
func Get(uri string, data Request) (forecast Forecast, err error) {
	// create json data
	jsonByte, err := json.Marshal(data)
	if err != nil {
		return forecast, fmt.Errorf("Marshaling forecast json failed: %v", err)
	}

	// send the request
	req, err := http.NewRequest("POST", uri, bytes.NewReader(jsonByte))
	if err != nil {
		return forecast, nil
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return forecast, fmt.Errorf("Http request to %s failed: %s", req.URL, err.Error())
	}
	defer resp.Body.Close()

	// decode the body
	dec := json.NewDecoder(resp.Body)
	if err = dec.Decode(&forecast); err != nil {
		return forecast, fmt.Errorf("Decoding forecast response failed: %v", err)
	}

	if forecast.Error != "" {
		return forecast, fmt.Errorf("Forecast API response error: %s", forecast.Error)
	}

	return forecast, nil
}
