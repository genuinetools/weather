package main

import (
	"flag"
	"fmt"
	"os"
)

type ForecastRequest struct {
	Latitude  float64  `json:"lat"`
	Longitude float64  `json:"lng"`
	Units     string   `json:"units"`
	Exclude   []string `json:"exclude"`
}

const VERSION = "v0.2.0"

func main() {
	var location string
	var units string
	var days int
	var ignoreAlerts bool
	var version bool

	// parse flags
	flag.BoolVar(&version, "version", false, "print version and exit")
	flag.BoolVar(&version, "v", false, "print version and exit (shorthand)")
	flag.StringVar(&location, "location", "", "Location to get the weather")
	flag.StringVar(&location, "l", "", "Location to get the weather (shorthand)")
	flag.StringVar(&units, "units", "auto", "System of units")
	flag.StringVar(&units, "u", "auto", "System of units (shorthand)")
	flag.IntVar(&days, "days", 0, "No. of days to get forecast")
	flag.IntVar(&days, "d", 0, "No. of days to get forecast (shorthand)")
	flag.BoolVar(&ignoreAlerts, "ignore-alerts", false, "Ignore alerts in weather output")
	flag.Parse()

	if version {
		fmt.Println(VERSION)
		return
	}

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
