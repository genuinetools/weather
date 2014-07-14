package main

import (
	"flag"
	"os"
)

type ForecastRequest struct {
	Latitude  float64  `json:"lat"`
	Longitude float64  `json:"lng"`
	Units     string   `json:"units"`
	Exclude   []string `json:"exclude"`
}

func main() {
	var location string
	var units string
	var days int
	var ignoreAlerts bool

	// parse flags
	flag.StringVar(&location, "location", "", "Location to get the weather")
	flag.StringVar(&location, "l", "", "Location to get the weather (shorthand)")
	flag.StringVar(&units, "units", "auto", "System of units")
	flag.StringVar(&units, "u", "auto", "System of units (shorthand)")
	flag.IntVar(&days, "days", 0, "No. of days to get forecast")
	flag.IntVar(&days, "d", 0, "No. of days to get forecast (shorthand)")
	flag.BoolVar(&ignoreAlerts, "ignore-alerts", false, "Ignore alerts in weather output")
	flag.Parse()

	geolocation, err := locate(location)
	if err != nil {
		printError(err)
		os.Exit(1)
	}

	data := ForecastRequest{
		Latitude:  geolocation.Latitude,
		Longitude: geolocation.Longitude,
		Units:     units,
		Exclude:   []string{"hourly", "minutely"},
	}

	forecast, err := getForecast(data)
	if err != nil {
		printError(err)
		os.Exit(1)
	}

	printCurrentWeather(forecast, geolocation, ignoreAlerts)

	if days > 1 {
		printDailyWeather(forecast, days)
	}
}
