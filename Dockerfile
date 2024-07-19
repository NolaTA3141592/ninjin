FROM golang:1.17.1-alpine3.13 as builder

WORKDIR /app/ninjin
COPY . .

RUN go mod download

RUN go build main.go

FROM alpine:3.13 AS release

LABEL maintainer="sohosai"
WORKDIR /app/ninjin

COPY --from=builder /app/ninjin/main .
COPY entrypoint.sh .

RUN chmod +x main
RUN chmod +x entrypoint.sh

EXPOSE 3000

ENTRYPOINT ["/app/ninjin/entrypoint.sh"]