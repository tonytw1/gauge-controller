FROM golang:1.20

WORKDIR /go/src/app
COPY gauges .
COPY client/dist/ client/dist

CMD ["gauges"]
