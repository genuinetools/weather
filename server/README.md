weather server
==============

API Server for `weather` command line tool. Connects to the [Google Geocode
API](https://developers.google.com/maps/documentation/geocoding/intro)
and [darksky.net API](https://darksky.net/dev/docs).

### Usage

```bash
$ ./weather-server --help
Usage of ./weather-server:
  --cert string
        path to ssl certificate
  --key string
        path to ssl key
  -darksky-apikey string
        Key for darksky API
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
    --darksky-apikey "YOUR_DARKSKY.NET_APIKEY"
```
