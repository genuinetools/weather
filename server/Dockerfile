FROM progrium/busybox
MAINTAINER Jessica Frazelle <jess@docker.com>

ADD https://jesss.s3.amazonaws.com/binaries/weather-server /usr/local/bin/weather-server

RUN chmod +x /usr/local/bin/weather-server

ENTRYPOINT [ "/usr/local/bin/weather-server" ]
