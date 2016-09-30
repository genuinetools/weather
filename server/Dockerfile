FROM alpine:latest
MAINTAINER Jessica Frazelle <jess@linux.com>

ENV PATH /go/bin:/usr/local/go/bin:$PATH
ENV GOPATH /go

RUN	apk add --no-cache \
	ca-certificates

RUN set -x \
	&& apk add --no-cache --virtual .build-deps \
		go \
		git \
		gcc \
		libc-dev \
		libgcc \
	&& go get -v github.com/jessfraz/weather \
	&& cd /go/src/github.com/jessfraz/weather \
	&& go build -o /usr/bin/weather-server ./server \
	&& apk del .build-deps \
	&& rm -rf /go \
	&& echo "Build complete."


ENTRYPOINT [ "weather-server" ]
