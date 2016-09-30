package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/Sirupsen/logrus"
	"github.com/jessfraz/weather/forecast"
	"github.com/jessfraz/weather/geocode"
)

const (
	forecastAPIURI = "https://api.forecast.io/forecast"
	geocodeAPIURI  = "https://maps.googleapis.com/maps/api/geocode/json"
)

var (
	forecastAPIKey string
	geocodeAPIKey  string

	port     string
	certFile string
	keyFile  string
)

// JSONResponse is a map[string]string
// response from the web server
type JSONResponse map[string]string

// String returns the string representation of the
// JSONResponse object
func (j JSONResponse) String() string {
	str, err := json.MarshalIndent(j, "", "  ")
	if err != nil {
		return fmt.Sprintf(`{
  "error": "%v"
}`, err)
	}

	return string(str)
}

// forecastHandler takes a forecast.Request object
// and passes it to the forecast.io API
func forecastHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var f forecast.Request
	if err := decoder.Decode(&f); err != nil {
		writeError(w, fmt.Sprintf("parsing request body for forecast failed: %v", err))
		return
	}

	// data to send to the API
	exclude, err := json.Marshal(f.Exclude)
	if err != nil {
		writeError(w, fmt.Sprintf("marshal forecast exclude failed: %v", err))
		return
	}
	data := url.Values{"units": {f.Units}, "exclude": {string(exclude)}}

	// request the forecast.io API
	url := fmt.Sprintf("%s/%s/%g,%g?%s", forecastAPIURI, forecastAPIKey, f.Latitude, f.Longitude, data.Encode())
	resp, err := http.Get(url)
	if err != nil {
		writeError(w, fmt.Sprintf("request to %s failed: %v", url, err))
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		writeError(w, fmt.Sprintf("reading response body from %s failed: %v", url, err))
		return
	}

	// write the response from the API to our client
	w.WriteHeader(resp.StatusCode)
	if _, err := w.Write(body); err != nil {
		writeError(w, fmt.Sprintf("writing response from %s failed: %v", url, err))
		return
	}
	return
}

// geocodeHandler takes a geocode.Request object
// and passes it to the Google Geocode API
func geocodeHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var g geocode.Request
	if err := decoder.Decode(&g); err != nil {
		writeError(w, fmt.Sprintf("parsing request body for geocode failed: %v", err))
		return
	}

	if g.Location == "" {
		writeError(w, "Location was not sent.")
		return
	}

	// data to send to the API
	data := url.Values{"address": {g.Location}, "key": {geocodeAPIKey}}

	// request the geocode API
	url := fmt.Sprintf("%s?%s", geocodeAPIURI, data.Encode())
	resp, err := http.Get(url)
	if err != nil {
		writeError(w, fmt.Sprintf("request to %s failed: %v", url, err))
		return
	}
	defer resp.Body.Close()

	decoder = json.NewDecoder(resp.Body)
	var geoResp geocode.Response
	if err := decoder.Decode(&geoResp); err != nil {
		writeError(w, fmt.Sprintf("parsing response body for geocode failed: %v", err))
		return
	}

	// These messages come from Google Geocoding API server
	if geoResp.ErrorMessage != "" {
		writeError(w, fmt.Sprintf("Google Geocode API response error: %s - %s", geoResp.Status, geoResp.ErrorMessage))
		return
	}

	// check if we have results
	if len(geoResp.Results) <= 0 {
		writeError(w, "No results found.")
		return
	}

	result := geoResp.Results[0]

	geo := geocode.Geocode{
		Latitude:  result.Geometry.Location.Latitude,
		Longitude: result.Geometry.Location.Longitude,
	}

	// parse each address for information to add to the geocode struct
	for _, addr := range result.AddressComponents {
		for _, t := range addr.Types {
			switch t {
			case "postal_code":
				geo.PostalCode = addr.LongName
			case "country":
				geo.Country = addr.LongName
				geo.CountryCode = addr.ShortName
				geo.CountryCode3 = addr.ShortName
			case "locality":
				geo.City = addr.LongName
			case "administrative_area_level_1":
				geo.Region = addr.LongName
				geo.RegionCode = addr.ShortName
			}
		}
	}

	// marshal the geo object
	body, err := json.Marshal(geo)
	if err != nil {
		writeError(w, fmt.Sprintf("marshal geo body failed: %v", err))
		return
	}

	// write the response from the API to our client
	w.WriteHeader(resp.StatusCode)
	if _, err := w.Write(body); err != nil {
		writeError(w, fmt.Sprintf("writing response from %s failed: %v", url, err))
		return
	}
	return
}

// failHandler returns not a valid endpoint
func failHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, JSONResponse{
		"error": fmt.Sprintf("Not a valid endpoint: %s", r.URL.Path),
	})
	return
}

// writeError sends an error back to the requester
// and also logrus. the error
func writeError(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, JSONResponse{
		"error": msg,
	})
	logrus.Printf("writing error: %s", msg)
	return
}

func init() {
	flag.StringVar(&forecastAPIKey, "forecast-apikey", "", "Key for forecast.io API")
	flag.StringVar(&geocodeAPIKey, "geocode-apikey", "", "Key for Google Maps Geocode API")

	flag.StringVar(&port, "p", "1234", "port for server to run on")
	flag.StringVar(&certFile, "cert", "", "path to ssl certificate")
	flag.StringVar(&keyFile, "key", "", "path to ssl key")

	flag.Parse()

	if forecastAPIKey == "" {
		logrus.Fatalf("You need to pass a forecast.io API Key")
	}

	if geocodeAPIKey == "" {
		logrus.Fatalf("You need to pass a Google Maps Geocode API Key")
	}
}

func main() {
	// create mux server
	mux := http.NewServeMux()

	mux.HandleFunc("/forecast", forecastHandler) // forecast handler
	mux.HandleFunc("/geocode", geocodeHandler)   // geocode handler
	mux.HandleFunc("/", failHandler)             // everything else fail handler

	// set up the server
	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	logrus.Infof("Starting server on port %q", port)
	if certFile != "" && keyFile != "" {
		logrus.Fatal(server.ListenAndServeTLS(certFile, keyFile))
	} else {
		logrus.Fatal(server.ListenAndServe())
	}
}
