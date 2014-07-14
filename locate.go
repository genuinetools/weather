package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// response from www.telize.com/geoip
// comes back like:
// {
//     timezone: "America/New_York",
//     isp: "Road Runner HoldCo LLC",
//     region_code: "NY",
//     country: "United States",
//     dma_code: "0",
//     area_code: "0",
//     region: "New York",
//     ip: "74.71.83.142",
//     asn: "AS11351",
//     continent_code: "NA",
//     city: "New York",
//     postal_code: "10128",
//     longitude: -73.9512,
//     latitude: 40.7805,
//     country_code: "US",
//     country_code3: "USA"
// }

type GeoLocation struct {
	AreaCode      string  `json:"area_code"`
	Asn           string  `json:"asn"`
	City          string  `json:"city"`
	ContinentCode string  `json:"continent_code"`
	Country       string  `json:"country"`
	CountryCode   string  `json:"country_code"`
	CountryCode3  string  `json:"country_code3"`
	DMACode       string  `json:"dma_code"`
	Ip            string  `json:"ip"`
	Isp           string  `json:"isp"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
	PostalCode    string  `json:"postal_code"`
	Region        string  `json:"region"`
	RegionCode    string  `json:"region_code"`
	Timezone      string  `json:"timezone"`
}

func requestLocation(uri string) (geolocation GeoLocation, err error) {
	resp, err := http.Get(uri)
	defer resp.Body.Close()

	// decode the body
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&geolocation)
	resp.Body.Close()

	if err != nil {
		return geolocation, fmt.Errorf("Decoding the response from %s failed: %s", uri, err)
	}

	return geolocation, nil
}

func autolocate() (geolocation GeoLocation, err error) {
	uri := "http://www.telize.com/geoip"

	return requestLocation(uri)
}

func locate(location string) (geolocation GeoLocation, err error) {
	if location == "" {
		return autolocate()
	}

	uri := "http://geocode.jessfraz.com/?location=" + url.QueryEscape(location)
	return requestLocation(uri)
}
