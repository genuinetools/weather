package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/genuinetools/pkg/cli"
	"github.com/genuinetools/weather/forecast"
	"github.com/genuinetools/weather/geocode"
	"github.com/genuinetools/weather/version"
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
	hideIcon     bool
	noForecast   bool
	jsonOut      bool
	server       string
	client       bool

	geo geocode.Geocode
)

//go:generate go run icons/generate.go

func main() {
	// Create a new cli program.
	p := cli.NewProgram()
	p.Name = "weather"
	p.Description = "Weather forecast via the command line"

	// Set the GitCommit and Version.
	p.GitCommit = version.GITCOMMIT
	p.Version = version.VERSION

	// Build the list of available commands.
	p.Commands = []cli.Command{
		&serverCommand{},
	}

	// Setup the global flags.
	p.FlagSet = flag.NewFlagSet("global", flag.ExitOnError)
	p.FlagSet.StringVar(&location, "location", "", "Location to get the weather")
	p.FlagSet.StringVar(&location, "l", "", "Location to get the weather (shorthand)")

	p.FlagSet.BoolVar(&client, "client", false, "Get location for the ssh client")
	p.FlagSet.BoolVar(&client, "c", false, "Get location for the ssh client (shorthand)")

	p.FlagSet.StringVar(&units, "units", "auto", "System of units (e.g. auto, us, si, ca, uk2)")
	p.FlagSet.StringVar(&units, "u", "auto", "System of units (shorthand) (e.g. auto, us, si, ca, uk2)")

	p.FlagSet.StringVar(&server, "server", defaultServerURI, "Weather API server uri")
	p.FlagSet.StringVar(&server, "s", defaultServerURI, "Weather API server uri (shorthand)")

	p.FlagSet.IntVar(&days, "days", 0, "No. of days to get forecast")
	p.FlagSet.IntVar(&days, "d", 0, "No. of days to get forecast (shorthand)")

	p.FlagSet.BoolVar(&ignoreAlerts, "ignore-alerts", false, "Ignore alerts in weather output")
	p.FlagSet.BoolVar(&hideIcon, "hide-icon", false, "Hide the weather icons from being output")
	p.FlagSet.BoolVar(&noForecast, "no-forecast", false, "Hide the forecast for the next 16 hours")

	p.FlagSet.BoolVar(&jsonOut, "json", false, "Prints the raw JSON API response")

	// Set the before function.
	p.Before = func(ctx context.Context) error {
		if len(server) < 1 {
			return errors.New("please enter a weather API server uri or leave blank to use the default")
		}

		return nil
	}

	// Set the main program action.
	p.Action = func(ctx context.Context, args []string) error {
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
					printError(errors.New("latitude and longitude could not be determined from your IP so the weather will not be accurate\nTry: weather -l <your_zipcode> OR weather -l \"your city, state\""))
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
			printError(errors.New("latitude and longitude could not be determined so the weather will not be accurate"))
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

		if jsonOut {
			jsn, err := json.Marshal(&fc)
			if err != nil {
				printError(err)
			}

			fmt.Println(string(jsn))
			return nil
		}

		if err := forecast.PrintCurrent(fc, geo, ignoreAlerts, hideIcon); err != nil {
			printError(err)
		}

		if days > 0 {
			forecast.PrintDaily(fc, days)
		}

		return nil
	}

	// Run our program.
	p.Run()
}

func printError(err error) {
	fmt.Println(colorstring.Color("[red]" + err.Error()))
	os.Exit(1)
}
