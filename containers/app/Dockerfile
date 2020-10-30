FROM golang:1.15.0

ARG CGO_ENABLED=0
ARG GOOS=linux
ARG GOARCH=amd64

WORKDIR /go/src/github.com/Tatsuemon/isucon9-qualify-shipment

COPY go.mod ./
RUN go mod download

COPY . .
RUN go build \
    -ldflags '-s -w'

RUN go get github.com/pilu/fresh