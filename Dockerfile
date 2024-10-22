# syntax=docker/dockerfile:1

# Build the application from source
FROM golang:1.19 AS build-stage

WORKDIR /be

ENV GOPROXY=https://goproxy.io,direct

COPY go.mod go.sum ./
RUN go mod download

COPY ./ ./

RUN CGO_ENABLED=0 GOOS=linux go build -o ./puxing-be ./cmd/app/main.go

RUN CGO_ENABLED=0 GOOS=linux go build -o ./worker ./cmd/worker/main.go

FROM alpine:3.19.1 AS app

WORKDIR /

COPY --from=build-stage ./be/puxing-be ./app
COPY .env /

EXPOSE 8000

ENTRYPOINT ["/app"]

FROM alpine:3.19.1 AS worker

WORKDIR /

COPY --from=build-stage ./be/worker ./worker
COPY .env /

ENTRYPOINT ["/worker"]