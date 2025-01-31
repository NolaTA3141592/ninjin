FROM golang:1.17.1-alpine3.13 as builder

WORKDIR /app/ninjin
COPY . .

RUN go mod download

RUN go build main.go