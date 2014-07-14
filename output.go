package main

import (
	"fmt"
	"github.com/mitchellh/colorstring"
	"time"
)

type UnitMeasures struct {
	Degrees          string
	Speed            string
	Pressure         string
	PressureMultiple float64
}

var (
	UnitFormats map[string]UnitMeasures = map[string]UnitMeasures{
		"metric": UnitMeasures{
			Degrees:          "°C",
			Speed:            "m/s",
			Pressure:         "hPa",
			PressureMultiple: 1,
		},
		"imperial": UnitMeasures{
			Degrees:          "°F",
			Speed:            "mph",
			Pressure:         "inHg",
			PressureMultiple: 0.0295,
		},
	}
	Directions []string = []string{
		"N", "NNE", "NE", "ENE", "E", "ESE", "SE", "SSE", "S", "SSW", "SW", "WSW", "W", "WNW", "NW", "NNW",
	}
)

func printWeather(forecast Forecast, geolocation GeoLocation, units string) {
	unitsFormat := UnitFormats[units]

	now := int64(time.Now().Unix())

	var sky, icon string
	isNight := false
	if len(forecast.Weather) > 0 {
		sky = cleanString(forecast.Weather[0].Main)

		if now >= forecast.Sys.Sunset || now < forecast.Sys.Sunrise {
			isNight = true
		}

		icon, _ = getIcon(sky, isNight)
	}

	fmt.Println(colorstring.Color("[blue]" + icon))
	fmt.Printf("\nCurrent weather in %s, %s for %s\n", geolocation.City, geolocation.Region, epochFormat(forecast.Date))

	if sky != "" && sky != "clear" {
		fmt.Printf("We see %s in the sky\n\n", sky)
	}

	fmt.Printf("The temperature is %v%s, with a high of %v%s & a low of %v%s\n\n", forecast.Temperature.Temperature, unitsFormat.Degrees, forecast.Temperature.TemperatureMax, unitsFormat.Degrees, forecast.Temperature.TemperatureMin, unitsFormat.Degrees)

	if forecast.Temperature.Humidity > 20 {
		fmt.Printf("Ick! The humidity is %v%s\n", forecast.Temperature.Humidity, "%")
	} else if forecast.Temperature.Humidity > 0 {
		fmt.Printf("The humidity is %v%s\n", forecast.Temperature.Humidity, "%")

	}

	fmt.Printf("The wind speed is %v %s %v\n", forecast.Wind.Speed, unitsFormat.Speed, getWindDetails(forecast.Wind.Degrees))

	fmt.Printf("The pressure is %v %s\n", Round(forecast.Temperature.Pressure*unitsFormat.PressureMultiple, 2), unitsFormat.Pressure)

	fmt.Printf("Sunrise is %v\n", epochFormat(forecast.Sys.Sunrise))
	fmt.Printf("Sunset is %v\n", epochFormat(forecast.Sys.Sunset))
}

func printDailyWeather(dailyForecast DailyForecast, geolocation GeoLocation, units string) {
	unitsFormat := UnitFormats[units]

	fmt.Printf("\nWeather in %s, %s for the next %v days\n\n", geolocation.City, geolocation.Region, dailyForecast.Count)

	for _, forecast := range dailyForecast.Forecasts {
		var sky string
		for _, weather := range forecast.Weather {
			sky = cleanString(weather.Main)
			if sky != "clear" && sky != "clouds" {
				sky += "y"
			} else if sky != "clouds" {
				sky = "cloudy"
			}

			break
		}

		fmt.Printf("On %s\n", epochFormatDate(forecast.Date))

		fmt.Printf("We are predicting %s skies.\n", sky)

		fmt.Printf("\tTemperature: %v%s\n", forecast.Temperature.Day, unitsFormat.Degrees)
		fmt.Printf("\tHigh: %v%s\n", forecast.Temperature.Max, unitsFormat.Degrees)
		fmt.Printf("\tLow: %v%s\n", forecast.Temperature.Min, unitsFormat.Degrees)
		fmt.Printf("\tMorning: %v%s\n", forecast.Temperature.Morning, unitsFormat.Degrees)
		fmt.Printf("\tNight: %v%s\n", forecast.Temperature.Night, unitsFormat.Degrees)

		if forecast.Rain > 0 {
			fmt.Printf("\n\tPrecipitation: %v", forecast.Rain)
		}
		if forecast.Clouds > 0 {
			fmt.Printf("\n\tCloudiness: %v%s\n", forecast.Clouds, "%")
		}

		fmt.Printf("\n\tHumidity: %v%s\n", forecast.Humidity, "%")

		fmt.Printf("\tWind: %v %s %v\n", forecast.Speed, unitsFormat.Speed, getWindDetails(forecast.Degrees))

		fmt.Printf("\tPressure: %v %s\n\n", Round(forecast.Pressure*unitsFormat.PressureMultiple, 2), unitsFormat.Pressure)
	}
}
