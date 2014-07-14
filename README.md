# weather

Weather via the command line.

![Screenshot](screenshot.png)

## Installation

#### Via Go

```bash
$ go get github.com/jfrazelle/weather
```

#### Binaries

- **darwin**
    - [386](https://jesss.s3.amazonaws.com/weather/binaries/darwin/386/weather)
    - [amd64](https://jesss.s3.amazonaws.com/weather/binaries/darwin/amd64/weather)
- **freebsd**
    - [386](https://jesss.s3.amazonaws.com/weather/binaries/freebsd/386/weather)
    - [amd64](https://jesss.s3.amazonaws.com/weather/binaries/freebsd/amd64/weather)
    - [arm](https://jesss.s3.amazonaws.com/weather/binaries/freebsd/arm/weather)
- **linux**
    - [386](https://jesss.s3.amazonaws.com/weather/binaries/linux/386/weather)
    - [amd64](https://jesss.s3.amazonaws.com/weather/binaries/linux/amd64/weather)
    - [arm](https://jesss.s3.amazonaws.com/weather/binaries/linux/arm/weather)
- **netbsd**
    - [386](https://jesss.s3.amazonaws.com/weather/binaries/netbsd/386/weather)
    - [amd64](https://jesss.s3.amazonaws.com/weather/binaries/netbsd/amd64/weather)
    - [arm](https://jesss.s3.amazonaws.com/weather/binaries/netbsd/arm/weather)
- **openbsd**
    - [386](https://jesss.s3.amazonaws.com/weather/binaries/openbsd/386/weather)
    - [amd64](https://jesss.s3.amazonaws.com/weather/binaries/openbsd/amd64/weather)
- **plan9**
    - [386](https://jesss.s3.amazonaws.com/weather/binaries/plan9/386/weather)
- **windows**
    - [386](https://jesss.s3.amazonaws.com/weather/binaries/windows/386/weather.exe)
    - [amd64](https://jesss.s3.amazonaws.com/weather/binaries/windows/amd64/weather.exe)


## Usage

- **`--location, -l`:** Your address, can be in the format of just a zipcode or a city, state, or the full address. **defaults to auto locating you based off your ip**
- **`--units, -u`:** The unit system to use. **defaults to `imperial`**, other option is `metric`
- **`--days, -d`:** Days of weather to retrieve. **defaults to the current weather, ie. 0 or 1**

### Examples

```bash
# get the current weather in your current location
$ weather

# change the units to metric
$ weather -l "Paris, France" -u metric

# get three days forecast for NY
$ weather -l 10028 -d 3

# or you can autolocate and get three days forecast
$ weather -d 3

# get the weather in Manhattan Beach, CA
$ weather -l "Manhattan Beach, CA"
#         .,.
#      ,dKNNX.
#    ,0NNXNNk
#   oXNO':NNx
#  lNNk  .XNX.
#  XNX.   cNN0.
# .NNX     :XNXo.
#  0NN,      lKNNKxc;,,:ldl
#  ,XNK.       'lkKXNNNNNNK.
#   ,KNXo.          .:0NNO.
#     oXNNOl;'...':dKNNK:
#       ,d0XNNNNNNNXOo'
#           .',,,'.
#
# Current weather in Manhattan Beach, California for July 13 at 10:00pm EDT
# The temperature is 79.2°F, with a high of 87.8°F & a low of 75.2°F
#
# Ick! The humidity is 46%
# The wind speed is 3.37 mph SSE
# The pressure is 29.76 inHg
# Sunrise is July 14 at 8:53am EDT
# Sunset is July 14 at 11:05pm EDT
```