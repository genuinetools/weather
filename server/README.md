weather server
==============

API Server for `weather` command line tool. Connects to the [Google Geocode
API](https://developers.google.com/maps/documentation/geocoding/intro)
and [forecast.io API](https://developer.forecast.io/docs/v2).

### Usage

```bash
$ ./weather-server --help
Usage of ./weather-server:
  --cert string
        path to ssl certificate
  --key string
        path to ssl key
  -forecast-apikey string
        Key for forecast.io API
  -geocode-apikey string
        Key for Google Maps Geocode API
  -p string
        port for server to run on (default "1234")
```

### Run with Docker

```
$ docker run --restart always -d \
    --name weather-server \
    -p 1234:1234 \
    jess/weather-server \
    --geocode-apikey "YOUR_GOOGLE_GEOCODING_APIKEY" \
    --forecast-apikey "YOUR_FORECAST.IO_APIKEY"
```
