# syntax=docker/dockerfile:1

## Build
FROM ghcr.io/internet24/golang:1.20.1-bullseye AS build

WORKDIR /app

COPY . .
RUN go mod tidy
RUN go build -o shadowsocks24

RUN ./third_party/outline-ss-server.sh

RUN tar -zcf third_party.tar.gz third_party
RUN tar -zcf web.tar.gz web

## Deploy
FROM ghcr.io/internet24/debian:bullseye-slim

WORKDIR /app

COPY --from=build /app/shadowsocks24 shadowsocks24
COPY --from=build /app/configs/config.json configs/config.json
COPY --from=build /app/assets/prometheus/configs/prometheus.yml storage/prometheus/configs/prometheus.yml
COPY --from=build /app/storage/database/.gitignore storage/database/.gitignore
COPY --from=build /app/storage/prometheus/data/.gitignore storage/prometheus/data/.gitignore
COPY --from=build /app/storage/shadowsocks/.gitignore storage/shadowsocks/.gitignore
COPY --from=build /app/third_party.tar.gz third_party.tar.gz
COPY --from=build /app/web.tar.gz web.tar.gz

RUN tar -xvf third_party.tar.gz
RUN tar -xvf web.tar.gz
RUN rm third_party.tar.gz
RUN rm web.tar.gz

EXPOSE 80

ENTRYPOINT ["./shadowsocks24", "start", "-c", "configs/config.json"]
