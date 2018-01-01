FROM golang:alpine as builder
MAINTAINER Jessica Frazelle <jess@linux.com>

ENV PATH /go/bin:/usr/local/go/bin:$PATH
ENV GOPATH /go

RUN	apk add --no-cache \
	ca-certificates

RUN set -x \
	&& apk add --no-cache --virtual .build-deps \
		git \
		gcc \
		libc-dev \
		libgcc \
	&& go get -v github.com/jessfraz/weather \
	&& cd /go/src/github.com/jessfraz/weather \
	&& CGO_ENABLED=0 go build -a -tags netgo -ldflags '-extldflags "-static"' -o /usr/bin/weather-server ./server \
	&& apk del .build-deps \
	&& rm -rf /go \
	&& echo "Build complete."

FROM scratch

COPY --from=builder /usr/bin/weather-server /usr/bin/weather-server
COPY --from=builder /etc/ssl/certs/ /etc/ssl/certs

ENTRYPOINT [ "weather-server" ]
CMD [ "--help" ]
