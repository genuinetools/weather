package geolocate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	geoipURI   string = "http://www.telize.com/geoip"
	geocodeURI string = "https://geocode.jessfraz.com/geocode"
)

// Geolocation response from www.telize.com/geoip
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
type Geolocation struct {
	AreaCode      string  `json:"area_code"`
	Asn           string  `json:"asn"`
	City          string  `json:"city"`
	ContinentCode string  `json:"continent_code"`
	Country       string  `json:"country"`
	CountryCode   string  `json:"country_code"`
	CountryCode3  string  `json:"country_code3"`
	DMACode       string  `json:"dma_code"`
	IP            string  `json:"ip"`
	Isp           string  `json:"isp"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
	PostalCode    string  `json:"postal_code"`
	Region        string  `json:"region"`
	RegionCode    string  `json:"region_code"`
	Timezone      string  `json:"timezone"`
}

// Autolocate gets the requesters geolocation based off their IP address
func Autolocate() (geolocation Geolocation, err error) {
	// send the request
	resp, err := http.Get(geoipURI)
	if err != nil {
		return geolocation, err
	}
	defer resp.Body.Close()

	// decode the body
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&geolocation); err != nil {
		return geolocation, fmt.Errorf("Decoding autolocate response failed: %v", err)
	}

	return geolocation, nil
}

// Locate gets the geolocation data of a location that is passed as a string
func Locate(location string) (geolocation Geolocation, err error) {
	// create json data
	jsonByte, err := json.Marshal(map[string]string{"location": location})
	if err != nil {
		return geolocation, fmt.Errorf("Marshaling location json failed: %v", err)
	}

	// send the request
	req, err := http.NewRequest("POST", geocodeURI, bytes.NewReader(jsonByte))
	if err != nil {
		return geolocation, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return geolocation, err
	}
	defer resp.Body.Close()

	// decode the body
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&geolocation); err != nil {
		return geolocation, fmt.Errorf("Decoding geolocate response failed: %v", err)
	}

	return geolocation, nil
}
