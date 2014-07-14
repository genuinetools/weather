# weather

Get weather via the commandline.

### Installation

```bash
$ go get github.com/jfrazelle/weather
```

### Usage

- **`--location, -l`:** Your address, can be in the format of just a zipcode or a city, state, or the full address. **defaults to auto locating you based off your ip**
- **`--units, -u`:** The unit system to use. **defaults to `imperial`**, other option is `metric`
- **`--days, -d`:** Days of weather to retrieve. **defaults to the current weather, ie. 0 or 1**

### Examples

```bash
# get the weather in Manhattan Beach, CA
$ weather -l "Manhattan Beach, CA"

# change the units to metric
$ weather -l "Paris, France" -u metric

# get three days forecast for NY
$ weather -l 10028 -d 3

# or you can autolocate yourself
$ weather -d 3

# or if you want the current weather in your current location
$ weather
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
# Current weather in New York, New York for July 13 at 7:56pm EDT
# We see clouds in the sky
#
# The temperature is 80.33°F, with a high of 84.2°F & a low of 73.4°F
#
# Ick! The humidity is 78%
# The wind speed is 9.94 mph SSE
# The pressure is 29.88 inHg
# Sunrise is July 14 at 5:37am EDT
# Sunset is July 14 at 8:26pm EDT
```