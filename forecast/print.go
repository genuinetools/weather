package forecast

import (
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/jfrazelle/weather/geolocate"
	"github.com/mitchellh/colorstring"
)

// UnitMeasures are the location specific terms for weather data
type UnitMeasures struct {
	Degrees       string
	Speed         string
	Length        string
	Precipitation string
}

var (
	// UnitFormats describe each regions UnitMeasures
	UnitFormats = map[string]UnitMeasures{
		"us": UnitMeasures{
			Degrees:       "°F",
			Speed:         "mph",
			Length:        "miles",
			Precipitation: "in/hr",
		},
		"si": UnitMeasures{
			Degrees:       "°C",
			Speed:         "m/s",
			Length:        "kilometers",
			Precipitation: "mm/h",
		},
		"ca": UnitMeasures{
			Degrees:       "°C",
			Speed:         "km/h",
			Length:        "kilometers",
			Precipitation: "mm/h",
		},
		"uk": UnitMeasures{
			Degrees:       "°C",
			Speed:         "mph",
			Length:        "kilometers",
			Precipitation: "mm/h",
		},
	}
	// Directions contain all the combinations of N,S,E,W
	Directions = []string{
		"N", "NNE", "NE", "ENE", "E", "ESE", "SE", "SSE", "S", "SSW", "SW", "WSW", "W", "WNW", "NW", "NNW",
	}
)

func epochFormat(seconds int64) string {
	epochTime := time.Unix(0, seconds*int64(time.Second))
	return epochTime.Format("January 2 at 3:04pm MST")
}

func epochFormatDate(seconds int64) string {
	epochTime := time.Unix(0, seconds*int64(time.Second))
	return epochTime.Format("January 2 (Monday)")
}

func epochFormatTime(seconds int64) string {
	epochTime := time.Unix(0, seconds*int64(time.Second))
	return epochTime.Format("3:04pm MST")
}

func getIcon(icon string) (iconTxt string, err error) {
	color := "blue"

	switch icon {
	case "clear-day":
		color = "yellow"
	case "clear-night":
		color = "light_yellow"
	case "snow":
		color = "white"
	case "wind":
		color = "black"
	case "partly-cloudy-day":
		color = "yellow"
	case "partly-cloudy-night":
		color = "light_yellow"
	case "thunderstorm":
		color = "black"
	case "tornado":
		color = "black"
	}
	uri := "https://s3.j3ss.co/weather/icons/" + icon + ".txt"

	resp, err := http.Get(uri)
	if err != nil {
		return iconTxt, fmt.Errorf("Requesting icon (%s) failed: %s", icon, err)
	}
	defer resp.Body.Close()

	// decode the body
	out, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return iconTxt, fmt.Errorf("Reading response body for icon (%s) failed: %s", icon, err)
	}

	iconTxt = string(out)

	if strings.Contains(iconTxt, "<?xml") {
		return "", fmt.Errorf("No icon found for %s.", icon)
	}

	return colorstring.Color("[" + color + "]" + iconTxt), nil
}

func getBearingDetails(degrees float64) string {
	index := int(math.Mod((degrees+11.25)/22.5, 16))
	return Directions[index]
}

func printCommon(weather Weather, unitsFormat UnitMeasures) error {
	if weather.Humidity > 0 {
		humidity := colorstring.Color(fmt.Sprintf("[white]%v%s", weather.Humidity*100, "%"))
		if weather.Humidity > 0.20 {
			fmt.Printf("Ick! The humidity is %s\n", humidity)
		} else {
			fmt.Printf("The humidity is %s\n", humidity)
		}
	}

	if weather.PrecipIntensity > 0 {
		precInt := colorstring.Color(fmt.Sprintf("[white]%v %s", weather.PrecipIntensity, unitsFormat.Precipitation))
		fmt.Printf("The precipitation intensity of %s is %s\n", colorstring.Color("[white]"+weather.PrecipType), precInt)
	}

	if weather.PrecipProbability > 0 {
		prec := colorstring.Color(fmt.Sprintf("[white]%v%s", weather.PrecipProbability*100, "%"))
		fmt.Printf("The precipitation probability is %s\n", prec)
	}

	if weather.NearestStormDistance > 0 {
		dist := colorstring.Color(fmt.Sprintf("[white]%v %s %v", weather.NearestStormDistance, unitsFormat.Length, getBearingDetails(weather.NearestStormBearing)))
		fmt.Printf("The nearest storm is %s away\n", dist)
	}

	if weather.WindSpeed > 0 {
		wind := colorstring.Color(fmt.Sprintf("[white]%v %s %v", weather.WindSpeed, unitsFormat.Speed, getBearingDetails(weather.WindBearing)))
		fmt.Printf("The wind speed is %s\n", wind)
	}

	if weather.CloudCover > 0 {
		cloudCover := colorstring.Color(fmt.Sprintf("[white]%v%s", weather.CloudCover*100, "%"))
		fmt.Printf("The cloud coverage is %s\n", cloudCover)
	}

	if weather.Visibility < 10 {
		visibility := colorstring.Color(fmt.Sprintf("[white]%v %s", weather.Visibility, unitsFormat.Length))
		fmt.Printf("The visibility is %s\n", visibility)
	}

	if weather.Pressure > 0 {
		pressure := colorstring.Color(fmt.Sprintf("[white]%v %s", weather.Pressure, "mbar"))
		fmt.Printf("The pressure is %s\n", pressure)
	}

	return nil
}

// PrintCurrent pretty prints the current forecast data
func PrintCurrent(forecast Forecast, geolocation geolocate.Geolocation, ignoreAlerts bool) error {
	unitsFormat := UnitFormats[forecast.Flags.Units]

	icon, err := getIcon(forecast.Currently.Icon)
	if err != nil {
		return err
	}

	fmt.Println(icon)

	location := colorstring.Color(fmt.Sprintf("[green]%s in %s", geolocation.City, geolocation.Region))
	fmt.Printf("\nCurrent weather is %s in %s for %s\n", colorstring.Color("[cyan]"+forecast.Currently.Summary), location, colorstring.Color("[cyan]"+epochFormat(forecast.Currently.Time)))

	temp := colorstring.Color(fmt.Sprintf("[magenta]%v%s", forecast.Currently.Temperature, unitsFormat.Degrees))
	feelslike := colorstring.Color(fmt.Sprintf("[magenta]%v%s", forecast.Currently.ApparentTemperature, unitsFormat.Degrees))
	fmt.Printf("The temperature is %s, but it feels like %s\n\n", temp, feelslike)

	if !ignoreAlerts {
		for _, alert := range forecast.Alerts {
			if alert.Title != "" {
				fmt.Println(colorstring.Color("[red]" + alert.Title))
			}
			if alert.Description != "" {
				fmt.Print(colorstring.Color("[red]" + alert.Description))
			}
			fmt.Println("\t\t\t" + colorstring.Color("[red]Created: "+epochFormat(alert.Time)))
			fmt.Println("\t\t\t" + colorstring.Color("[red]Expires: "+epochFormat(alert.Expires)) + "\n")
		}
	}

	return printCommon(forecast.Currently, unitsFormat)
}

// PrintDaily pretty prints the daily forecast data
func PrintDaily(forecast Forecast, days int) error {
	unitsFormat := UnitFormats[forecast.Flags.Units]

	fmt.Println(colorstring.Color("\n[white]" + fmt.Sprintf("%v Day Forecast", days)))

	for index, daily := range forecast.Daily.Data {
		// only do the amount of days they request
		if index == days {
			break
		}

		fmt.Println(colorstring.Color("\n[magenta]" + epochFormatDate(daily.Time)))

		tempMax := colorstring.Color(fmt.Sprintf("[blue]%v%s", daily.TemperatureMax, unitsFormat.Degrees))
		tempMin := colorstring.Color(fmt.Sprintf("[blue]%v%s", daily.TemperatureMin, unitsFormat.Degrees))
		feelsLikeMax := colorstring.Color(fmt.Sprintf("[cyan]%v%s", daily.ApparentTemperatureMax, unitsFormat.Degrees))
		feelsLikeMin := colorstring.Color(fmt.Sprintf("[cyan]%v%s", daily.ApparentTemperatureMin, unitsFormat.Degrees))
		fmt.Printf("The temperature high is %s, feels like %s around %s, and low is %s, feels like %s around %s\n\n", tempMax, feelsLikeMax, epochFormatTime(daily.TemperatureMaxTime), tempMin, feelsLikeMin, epochFormatTime(daily.TemperatureMinTime))

		printCommon(daily, unitsFormat)
	}

	return nil
}
