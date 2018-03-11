package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/genuinetools/weather/forecast"
	"github.com/genuinetools/weather/geocode"
	"github.com/genuinetools/weather/version"
	"github.com/mitchellh/colorstring"
)

const (
	defaultServerURI string = "https://geocode.genuinetools.com"
)

var (
	location     string
	units        string
	days         int
	ignoreAlerts bool
	hideIcon     bool
	noForecast   bool
	server       string
	vrsn         bool
	client       bool
	geo          geocode.Geocode
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
	flag.BoolVar(&client, "client", false, "Get location for the ssh client")
	flag.BoolVar(&client, "c", false, "Get location for the ssh client (shorthand)")
	flag.StringVar(&units, "units", "auto", "System of units")
	flag.StringVar(&units, "u", "auto", "System of units (shorthand)")
	flag.StringVar(&server, "server", defaultServerURI, "Weather API server uri")
	flag.StringVar(&server, "s", defaultServerURI, "Weather API server uri (shorthand)")
	flag.IntVar(&days, "days", 0, "No. of days to get forecast")
	flag.IntVar(&days, "d", 0, "No. of days to get forecast (shorthand)")
	flag.BoolVar(&ignoreAlerts, "ignore-alerts", false, "Ignore alerts in weather output")
	flag.BoolVar(&hideIcon, "hide-icon", false, "Hide the weather icons from being output")
	flag.BoolVar(&noForecast, "no-forecast", false, "Hide the forecast for the next 16 hours")

	flag.Usage = func() {
		flag.PrintDefaults()
	}

	flag.Parse()

	if server == "" {
		usageAndExit("Please enter a Weather API server uri or leave blank to use the default.", 0)
	}
}

//go:generate go run icons/generate.go

func main() {
	if vrsn {
		fmt.Printf("weather version %s, build %s", version.VERSION, version.GITCOMMIT)
		return
	}

	var err error
	if location == "" {
		sshConn := os.Getenv("SSH_CONNECTION")
		if client && len(sshConn) > 0 {
			// use their ssh connection to locate them
			ipports := strings.Split(sshConn, " ")
			geo, err = geocode.IPLocate(ipports[0])
			if err != nil {
				printError(err)
			}

		} else {
			// auto locate them
			geo, err = geocode.Autolocate()
			if err != nil {
				printError(err)
			}

			if geo.Latitude == 0 || geo.Longitude == 0 {
				printError(errors.New("Latitude and Longitude could not be determined from your IP so the weather will not be accurate\nTry: weather -l <your_zipcode> OR weather -l \"your city, state\""))
			}
		}
	} else {
		// get geolocation data for the given location
		geo, err = geocode.Locate(location, server)
		if err != nil {
			printError(err)
		}
	}

	if geo.Latitude == 0 || geo.Longitude == 0 {
		printError(errors.New("Latitude and Longitude could not be determined so the weather will not be accurate"))
	}

	data := forecast.Request{
		Latitude:  geo.Latitude,
		Longitude: geo.Longitude,
		Units:     units,
		Exclude:   []string{"minutely"},
	}
	if noForecast {
		data.Exclude = append(data.Exclude, "hourly")
	}

	fc, err := forecast.Get(fmt.Sprintf("%s/forecast", server), data)
	if err != nil {
		printError(err)
	}

	if err := forecast.PrintCurrent(fc, geo, ignoreAlerts, hideIcon); err != nil {
		printError(err)
	}

	if days > 0 {
		forecast.PrintDaily(fc, days)
	}
}

func usageAndExit(message string, exitCode int) {
	if message != "" {
		fmt.Fprintf(os.Stderr, message)
		fmt.Fprintf(os.Stderr, "\n\n")
	}
	flag.Usage()
	fmt.Fprintf(os.Stderr, "\n")
	os.Exit(exitCode)
}
