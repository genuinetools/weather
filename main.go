package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/jfrazelle/weather/forecast"
	"github.com/jfrazelle/weather/geolocate"
	"github.com/jfrazelle/weather/version"
	"github.com/mitchellh/colorstring"
)

const (
	forecastURI string = "https://geocode.jessfraz.com/forecast"
)

var (
	location     string
	units        string
	days         int
	ignoreAlerts bool
	vrsn         bool

	geolocation geolocate.Geolocation
)

func printError(err error) {
	fmt.Println(colorstring.Color("[red]" + err.Error()))
	os.Exit(1)
}

func init() {
	// parse flags
	flag.BoolVar(&vrsn, "version", false, "print version and exit")
	flag.BoolVar(&vrsn, "v", false, "print version and exit (shorthand)")
	flag.StringVar(&location, "location", "", "Location to get the weather")
	flag.StringVar(&location, "l", "", "Location to get the weather (shorthand)")
	flag.StringVar(&units, "units", "auto", "System of units")
	flag.StringVar(&units, "u", "auto", "System of units (shorthand)")
	flag.IntVar(&days, "days", 0, "No. of days to get forecast")
	flag.IntVar(&days, "d", 0, "No. of days to get forecast (shorthand)")
	flag.BoolVar(&ignoreAlerts, "ignore-alerts", false, "Ignore alerts in weather output")
	flag.Parse()
}

//go:generate go run icons/generate/generate.go

func main() {
	if vrsn {
		fmt.Printf("weather version %s, build %s", version.VERSION, version.GITCOMMIT)
		return
	}

	var err error
	if location == "" {
		// auto locate them if we are not given a location
		geolocation, err = geolocate.Autolocate()
		if err != nil {
			printError(err)
		}
	} else {
		// get geolocation data for the given location
		geolocation, err = geolocate.Locate(location)
		if err != nil {
			printError(err)
		}
	}

	data := forecast.Request{
		Latitude:  geolocation.Latitude,
		Longitude: geolocation.Longitude,
		Units:     units,
		Exclude:   []string{"hourly", "minutely"},
	}

	fc, err := forecast.Get(forecastURI, data)
	if err != nil {
		printError(err)
	}

	if err := forecast.PrintCurrent(fc, geolocation, ignoreAlerts); err != nil {
		printError(err)
	}

	if days > 1 {
		forecast.PrintDaily(fc, days)
	}
}
