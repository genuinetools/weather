package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"strings"
	"time"
)

func cleanString(s string) string {
	return strings.TrimSpace(strings.ToLower(s))
}

func epochFormat(seconds int64) string {
	epochTime := time.Unix(0, seconds*int64(time.Second))
	return epochTime.Format("January 2 at 3:04pm MST")
}

func epochFormatDate(seconds int64) string {
	epochTime := time.Unix(0, seconds*int64(time.Second))
	return epochTime.Format("January 2")
}

func getIcon(sky string, isNight bool) (icon string, err error) {
	uri := "http://jesss.s3.amazonaws.com/weather/icons/" + sky
	if isNight {
		uri += "_night"
	}

	resp, err := http.Get(uri + ".txt")
	if err != nil {
		// if it's a night icon, try for day
		if isNight {
			return getIcon(sky, false)
		}

		return icon, fmt.Errorf("Requesting icon (%s) failed: %s", sky, err)
	}
	defer resp.Body.Close()

	// decode the body
	out, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// if it's a night icon, try for day
		if isNight {
			return getIcon(sky, false)
		}

		return icon, fmt.Errorf("Reading response body for icon (%s) failed: %s", sky, err)
	}

	icon = string(out)

	// if it's a night icon, try for day
	if (isNight && icon == "") || (isNight && strings.Contains(icon, "<?xml")) {
		return getIcon(sky, false)
	}

	return icon, nil
}

func getWindDetails(degrees float64) (direction string) {
	windDeg := (degrees + 11.25) / 22.5
	directionInt := int(math.Abs(math.Remainder(windDeg, 16)))

	if len(Directions) > directionInt && directionInt >= 0 {
		direction = Directions[directionInt]
	}

	return direction
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
