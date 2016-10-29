package geocode

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	geoipURI string = "https://telize.j3ss.co/geoip"
)

// Geocode response from www.telize.com/geoip
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
type Geocode struct {
	AreaCode      string  `json:"area_code"`
	Asn           string  `json:"asn"`
	City          string  `json:"city"`
	ContinentCode string  `json:"continent_code"`
	Country       string  `json:"country"`
	CountryCode   string  `json:"country_code"`
	CountryCode3  string  `json:"country_code3"`
	DMACode       string  `json:"dma_code"`
	Error         string  `json:"error"`
	IP            string  `json:"ip"`
	Isp           string  `json:"isp"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
	PostalCode    string  `json:"postal_code"`
	Region        string  `json:"region"`
	RegionCode    string  `json:"region_code"`
	Timezone      string  `json:"timezone"`
}

// Request describes the request posted to the geocode api.
type Request struct {
	Location string `json:"location"`
}

// Response is the response from the  Google Geocoding API
// it comes back like the following:
// {
//   "results" : [
//      {
//         "address_components" : [
//            {
//               "long_name" : "Manhattan Beach",
//               "short_name" : "Manhattan Beach",
//               "types" : [ "locality", "political" ]
//            },
//            {
//               "long_name" : "Los Angeles County",
//               "short_name" : "Los Angeles County",
//               "types" : [ "administrative_area_level_2", "political" ]
//            },
//            {
//               "long_name" : "California",
//               "short_name" : "CA",
//               "types" : [ "administrative_area_level_1", "political" ]
//            },
//            {
//               "long_name" : "United States",
//               "short_name" : "US",
//               "types" : [ "country", "political" ]
//            }
//         ],
//         "formatted_address" : "Manhattan Beach, CA, USA",
//         "geometry" : {
//            "bounds" : {
//               "northeast" : {
//                  "lat" : 33.906185,
//                  "lng" : -118.3785991
//               },
//               "southwest" : {
//                  "lat" : 33.8728038,
//                  "lng" : -118.4234639
//               }
//            },
//            "location" : {
//               "lat" : 33.8847361,
//               "lng" : -118.4109089
//            },
//            "location_type" : "APPROXIMATE",
//            "viewport" : {
//               "northeast" : {
//                  "lat" : 33.906185,
//                  "lng" : -118.3785991
//               },
//               "southwest" : {
//                  "lat" : 33.8728038,
//                  "lng" : -118.4234639
//               }
//            }
//         },
//         "place_id" : "ChIJL2Ow4sWzwoARIEUV9NRRAwc",
//         "types" : [ "locality", "political" ]
//      }
//   ],
//   "status" : "OK"
// }
type Response struct {
	Results      []Result `json:"results"`
	Status       string   `json:"status"`
	ErrorMessage string   `json:"error_message"`
}

// Result is the result from the Google geocoding API.
type Result struct {
	AddressComponents []Address `json:"address_components"`
	FormattedAddress  string    `json:"formatted_address"`
	Geometry          Geometry  `json:"geometry"`
	PlaceID           string    `json:"place_id"`
	Types             []string  `json:"types"`
}

// Address contains details of a place.
type Address struct {
	LongName  string   `json:"long_name"`
	ShortName string   `json:"short_name"`
	Types     []string `json:"types"`
}

// Geometry contains location data.
type Geometry struct {
	Bounds       map[string]Location `json:"bounds"`
	Location     Location            `json:"location"`
	LocationType string              `json:"location_type"`
	Viewport     map[string]Location `json:"viewport"`
}

// Location contains lat and lng.
type Location struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lng"`
}

// Autolocate gets the requesters geocode response based off their IP address.
func Autolocate() (geocode Geocode, err error) {
	// send the request
	resp, err := http.Get(geoipURI)
	if err != nil {
		return geocode, err
	}
	defer resp.Body.Close()

	// decode the body
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&geocode); err != nil {
		return geocode, fmt.Errorf("Decoding autolocate response failed: %v", err)
	}

	return geocode, nil
}

// IPLocate gets the requesters geocode response based off an IP address.
func IPLocate(ip string) (geocode Geocode, err error) {
	// send the request
	resp, err := http.Get(fmt.Sprintf("%s/%s", geoipURI, ip))
	if err != nil {
		return geocode, err
	}
	defer resp.Body.Close()

	// decode the body
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&geocode); err != nil {
		return geocode, fmt.Errorf("Decoding autolocate response failed: %v", err)
	}

	return geocode, nil
}

// Locate gets the geocode data of a location that is passed as a string.
func Locate(location, server string) (geocode Geocode, err error) {
	// create json data
	jsonByte, err := json.Marshal(Request{Location: location})
	if err != nil {
		return geocode, fmt.Errorf("Marshaling location json failed: %v", err)
	}

	// send the request
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/geocode", server), bytes.NewReader(jsonByte))
	if err != nil {
		return geocode, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return geocode, err
	}
	defer resp.Body.Close()

	// decode the body
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&geocode); err != nil {
		return geocode, fmt.Errorf("Decoding geocode response failed: %v", err)
	}

	// These messages come from our API server
	if geocode.Error != "" {
		return geocode, fmt.Errorf("Geocode API response error: %s", geocode.Error)
	}

	return geocode, nil
}
