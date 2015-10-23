FROM debian:sid
MAINTAINER Jessica Frazelle <jess@docker.com>

RUN apt-get update && apt-get install -y \
	ca-certificates \
	--no-install-recommends \
	&& rm -rf /var/apt/lists/*

ADD https://jesss.s3.amazonaws.com/binaries/weather-server /usr/local/bin/weather-server

RUN chmod +x /usr/local/bin/weather-server

ENTRYPOINT [ "/usr/local/bin/weather-server" ]
