package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	var location string
	var units string
	var days int

	// parse flags
	flag.StringVar(&location, "location", "", "Location to get the weather")
	flag.StringVar(&location, "l", "", "Location to get the weather (shorthand)")
	flag.StringVar(&units, "units", "imperial", "System of units")
	flag.StringVar(&units, "u", "imperial", "System of units (shorthand)")
	flag.IntVar(&days, "days", 0, "No. of days to get forecast")
	flag.IntVar(&days, "d", 0, "No. of days to get forecast (shorthand)")
	flag.Parse()

	geolocation, err := locate(location)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	if days <= 1 {
		forecast, err := getForecast(geolocation, units)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		printWeather(forecast, geolocation, units)
	} else {
		dailyForecast, err := getDailyForecast(geolocation, units, days)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		printDailyWeather(dailyForecast, geolocation, units)
	}
}
