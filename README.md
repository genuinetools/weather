# weather

[![make-all](https://github.com/genuinetools/weather/workflows/make%20all/badge.svg)](https://github.com/genuinetools/weather/actions?query=workflow%3A%22make+all%22)
[![make-image](https://github.com/genuinetools/weather/workflows/make%20image/badge.svg)](https://github.com/genuinetools/weather/actions?query=workflow%3A%22make+image%22)
[![GoDoc](https://img.shields.io/badge/godoc-reference-5272B4.svg?style=for-the-badge)](https://godoc.org/github.com/genuinetools/weather)
[![Github All Releases](https://img.shields.io/github/downloads/genuinetools/weather/total.svg?style=for-the-badge)](https://github.com/genuinetools/weather/releases)

Weather via the command line. Uses the [darksky.net](https://darksky.net) API so it's super accurate. Also includes any current weather alerts in the output.

![Screenshot](screenshot.png)

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**

- [Installation](#installation)
    - [Binaries](#binaries)
    - [Via Go](#via-go)
    - [Via Homebrew](#via-homebrew)
- [Usage](#usage)
  - [Examples](#examples)
- [Running the Server](#running-the-server)
    - [Usage](#usage-1)
    - [Running with Docker](#running-with-docker)
- [Contributing](#contributing)
    - [Makefile Usage](#makefile-usage)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Installation

#### Binaries

For installation instructions from binaries please visit the [Releases Page](https://github.com/genuinetools/weather/releases).

#### Via Go

```console
$ go get github.com/genuinetools/weather
```

#### Via Homebrew

```console
$ brew install darksky-weather
```

## Usage

```console
$ weather -h
weather -  Weather forecast via the command line.

Usage: weather <command>

Flags:

  -c              Get location for the ssh client (shorthand) (default: false)
  -client         Get location for the ssh client (default: false)
  -d              No. of days to get forecast (shorthand) (default: 0)
  -days           No. of days to get forecast (default: 0)
  -hide-icon      Hide the weather icons from being output (default: false)
  -ignore-alerts  Ignore alerts in weather output (default: false)
  -l              Location to get the weather (shorthand) (default: <none>)
  -location       Location to get the weather (default: <none>)
  -no-forecast    Hide the forecast for the next 16 hours (default: false)
  -s              Weather API server uri (shorthand) (default: https://geocode.jessfraz.com)
  -server         Weather API server uri (default: https://geocode.jessfraz.com)
  -u              System of units (shorthand) (e.g. auto, us, si, ca, uk2) (default: auto)
  -units          System of units (e.g. auto, us, si, ca, uk2) (default: auto)

Commands:

  version  Show the version information.
```

### Examples

```bash
# get the current weather in your current location
$ weather

# change the units to metric
$ weather -l "Paris, France" -u si

# it will auto guess the units though so changing
# the location to paris will change the units to `si`
$ weather -l "Paris, France"

# get three days forecast for NY
$ weather -l 10028 -d 3

# or you can autolocate and get three days forecast
$ weather -d 3

# get the weather in Manhattan Beach, CA
# even includes alerts
$ weather -l "Manhattan Beach, CA"
#                             .;odc
#                           ;kXNNNO
#                         .0NNO0NN:
#                        'XNK; dNNl
#                        KNX'  'XNK.
#                       ,NNk    cXNK,
#                       ,NNk     '0NNO:.
#                     .'cXNXl;,.   ,xXNNKOxxxk0Xx
#                 'lOXNNNNNNNNNNXOo'  ':oxkOXNNXc
#               cKNNKd:'.    ..;d0NNKl    ,xXNK,
#        .;:cclKNXd.              .oXNXxOXNNXl
#    .cOXNNNNNNNO.                  .kNNNNNNNXOc.
#   lXNXx;.    .                      .    .;dXNXo
#  ONNd.                                       oXN0.
# dNNo                                          cNNk
# XNN.                                           NNX
# 0NN'                                          .NNK
# ;XN0.                                        .ONNc
#  ;XNXo.                                    .lXNX:
#   .oXNX0dlcclx0Xo.              .oXKxlccldOXNXd.
#      ,lk0KXXK0xKNN0o;..    ..;o0NNKx0KXXX0ko,
#                 'lOXNNNNNNNNNNXOo,
#                     :x0XNNX0x:.
#
#
# Current weather is Partly Cloudy in Manhattan Beach in California for July 14 at 4:14am EDT
# The temperature is 69.2°F, but it feels like 69.2°F
#
# Special Weather Statement for Los Angeles, CA
# ...THREAT OF MONSOONAL THUNDERSTORMS LATE TONIGHT THROUGH WEDNESDAY...
# A STRONG UPPER LEVEL HIGH PRESSURE SYSTEM CURRENTLY CENTERED OVER NEVADA
# WILL BRING INCREASING EAST TO SOUTHEAST FLOW OVER SOUTHERN
# CALIFORNIA. AS A RESULT...A SIGNIFICANT SURGE OF MONSOONAL MOISTURE
# WILL MOVE INTO SOUTHWEST CALIFORNIA LATE TONIGHT THROUGH WEDNESDAY.
# THE GREATEST THREAT OF SHOWERS AND THUNDERSTORMS WILL BE ACROSS THE
# MOUNTAINS AND ANTELOPE VALLEY LATE TONIGHT INTO TUESDAY. DUE TO THE
# EASTERLY UPPER LEVEL FLOW ON MONDAY...THERE WILL ALSO BE A SLIGHT
# CHANCE OF SHOWERS AND THUNDERSTORMS ACROSS MOST COASTAL AND VALLEY
# AREAS.
# THE DEEPER MONSOONAL MOISTURE WILL BRING THE POTENTIAL FOR BRIEF HEAVY
# RAINFALL WITH STORMS THAT DEVELOP ON MONDAY AND TUESDAY...ESPECIALLY
# ACROSS THE MOUNTAINS AND ANTELOPE VALLEY. WHILE STORMS ARE EXPECTED
# TO BE FAST MOVING...THERE WILL BE THE POTENTIAL FOR LOCALIZED FLOODING
# OF ROADWAYS AND ARROYOS. ON TUESDAY...THE THREAT OF THUNDERSTORMS IS
# EXPECTED TO REMAIN CONFINED TO THE MOUNTAINS AND DESERTS. WITH WEAKER
# UPPER LEVEL WINDS ON TUESDAY...STORMS WILL LIKELY MOVE SLOWER. AS A
# RESULT...THERE WILL BE AN INCREASED THREAT OF FLASH FLOODING.
# IT WILL NOT BE AS HOT ACROSS MUCH OF THE REGION TOMORROW DUE TO THE
# INCREASED MOISTURE AND CLOUD COVERAGE...WITH INTERIOR SECTIONS
# GENERALLY REMAINING IN THE 90S. HOWEVER...THERE WILL BE A
# SIGNIFICANT INCREASE IN HUMIDITY ON MONDAY THAT WILL CONTINUE TO
# BRING DISCOMFORT.
# ANYONE PLANNING OUTDOOR ACTIVITIES IN THE MOUNTAINS AND DESERTS
# DURING THE NEXT FEW DAYS SHOULD CAREFULLY MONITOR THE LATEST
# NATIONAL WEATHER SERVICE FORECASTS AND STATEMENTS DUE TO THE
# POTENTIAL HAZARDS ASSOCIATED WITH THUNDERSTORMS.
#             Created: July 13 at 10:50pm EDT
#             Expires: July 14 at 7:00pm EDT
#
# Ick! The humidity is 85%
# The nearest storm is 18 miles NE away
# The wind speed is 3.96 mph SE
# The cloud coverage is 35%
# The visibility is 9.58 miles
# The pressure is 1012.99 mbar
```

## Running the Server

API Server for `weather` command line tool. Connects to the [Google Geocode
API](https://developers.google.com/maps/documentation/geocoding/intro)
and [darksky.net API](https://darksky.net/dev/docs).

#### Usage

```bash
$ weather server -h
Usage: weather server [OPTIONS]

Run a static UI server for a registry.

Flags:

  -cert            path to ssl cert (default: <none>)
  -darksky-apikey  Key for darksky.net API (default: <none>)
  -geocode-apikey  Key for Google Maps Geocode API (default: <none>)
  -key             path to ssl key (default: <none>)
  -port            port for server to run on (default: 1234)
```

#### Running with Docker

```console
$ docker run --restart always -d \
    --name weather-server \
    -p 1234:1234 \
    r.j3ss.co/weather server \
    --geocode-apikey "YOUR_GOOGLE_GEOCODING_APIKEY" \
    --darksky-apikey "YOUR_DARKSKY.NET_APIKEY"
```

## Contributing

Please do!

#### Makefile Usage

```console
$ make help
all                            Runs a clean, build, fmt, lint, test, staticcheck, vet and install
build                          Builds a dynamic executable or package
bump-version                   Bump the version in the version file. Set BUMP to [ patch | major | minor ]
clean                          Cleanup any build binaries or packages
cover                          Runs go test with coverage
cross                          Builds the cross-compiled binaries, creating a clean directory structure (eg. GOOS/GOARCH/binary)
fmt                            Verifies all files have been `gofmt`ed
install                        Installs the executable or package
lint                           Verifies `golint` passes
release                        Builds the cross-compiled binaries, naming them in such a way for release (eg. binary-GOOS-GOARCH)
static                         Builds a static executable
staticcheck                    Verifies `staticcheck` passes
tag                            Create a new git tag to prepare to build a release
test                           Runs the go tests
vet                            Verifies `go vet` passes
```

[![Analytics](https://ga-beacon.appspot.com/UA-29404280-16/weather/README.md)](https://github.com/genuinetools/weather)
