package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/mitchellh/colorstring"
)

func createRequest(uri, method string, data interface{}) (req *http.Request, err error) {
	// create json data
	jsonByte, err := json.Marshal(data)
	if err != nil {
		return req, fmt.Errorf("Marshaling JSON for %s to %s failed: %s", method, uri, err.Error())
	}

	// send the request
	req, err = http.NewRequest("POST", uri, bytes.NewReader(jsonByte))
	if err != nil {
		return req, fmt.Errorf("Creating the %s request to %s failed: %s", method, uri, err.Error())
	}
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

func epochFormat(seconds int64) string {
	epochTime := time.Unix(0, seconds*int64(time.Second))
	return epochTime.Format("January 2 at 3:04pm MST")
}

func epochFormatDate(seconds int64) string {
	epochTime := time.Unix(0, seconds*int64(time.Second))
	return epochTime.Format("January 2")
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
	uri := "https://jesss.s3.amazonaws.com/weather/icons/" + icon + ".txt"

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

func printError(err error) {
	fmt.Println(colorstring.Color("[red]" + err.Error()))
}

func Round(x float64, prec int) float64 {
	var rounder float64
	pow := math.Pow(10, float64(prec))
	intermed := x * pow
	_, frac := math.Modf(intermed)
	if frac >= 0.5 {
		rounder = math.Ceil(intermed)
	} else {
		rounder = math.Floor(intermed)
	}

	return rounder / pow
}
