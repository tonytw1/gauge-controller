FROM golang:1.20

WORKDIR /go/src/app
COPY . .

RUN go get -d -t -v ./...
RUN go install -v ./...
RUN go test -v
RUN go build -v

CMD ["gauges"]
