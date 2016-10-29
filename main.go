package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"github.com/jessfraz/weather/forecast"
	"github.com/jessfraz/weather/geocode"
	"github.com/jessfraz/weather/version"
	"github.com/mitchellh/colorstring"
)

const (
	defaultServerURI string = "https://geocode.jessfraz.com"
)

var (
	location     string
	units        string
	days         int
	ignoreAlerts bool
	server       string
	vrsn         bool
	client       bool
	geo geocode.Geocode
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
		if client == false {
			// auto locate them if we are not given a location or ssh data
			geo, err = geocode.Autolocate()
			if err != nil {
				printError(err)
			}
		} else {
			var ssh_conn string
			ssh_conn = os.Getenv("SSH_CONNECTION")
			if len(ssh_conn) > 0 {
				var ipports []string
				ipports = strings.Split(ssh_conn, " ")
				geo, err = geocode.Iplocate(ipports[0])
				if err != nil {
					printError(err)
				}
			} else {
				// auto locate them
				geo, err = geocode.Autolocate()
				if err != nil {
					printError(err)
				}
			}
		}
	} else {
		// get geolocation data for the given location
		geo, err = geocode.Locate(location, server)
		if err != nil {
			printError(err)
		}
	}

	data := forecast.Request{
		Latitude:  geo.Latitude,
		Longitude: geo.Longitude,
		Units:     units,
		Exclude:   []string{"hourly", "minutely"},
	}

	fc, err := forecast.Get(fmt.Sprintf("%s/forecast", server), data)
	if err != nil {
		printError(err)
	}

	if err := forecast.PrintCurrent(fc, geo, ignoreAlerts); err != nil {
		printError(err)
	}

	if days > 1 {
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
