FROM golang:1.4
COPY . /go/src/github.com/jfrazelle/weather
RUN go get -v github.com/jfrazelle/weather
ENTRYPOINT [ "weather" ]
