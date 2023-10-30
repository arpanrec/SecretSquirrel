FROM golang:1.21.3 AS build-stage

WORKDIR /app

COPY . .

RUN go mod download
RUN go build -o secureserver ./main.go

FROM debian:12

WORKDIR /

COPY --from=build-stage /app/secureserver /secureserver

ENTRYPOINT ["/secureserver"]
