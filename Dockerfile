FROM golang:1.20

WORKDIR /opt/gauges
COPY gauges .
COPY config.json .
COPY client/dist/ client/dist

CMD ["/opt/gauges/gauges"]
