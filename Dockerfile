FROM golang:latest as build

WORKDIR /app

ARG service
COPY cmd/$service/main.go main.go
COPY pkg/database/ pkg/database/
COPY pkg/server/ pkg/server/
COPY pkg/domain/ pkg/domain/
COPY pkg/utils/ pkg/utils/
COPY pkg/$service/ pkg/$service/
COPY go.mod go.mod

RUN go mod tidy
RUN go build -o main main.go

FROM ubuntu:latest
ENV TZ=Europe/Moscow
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

ENV DEBIAN_FRONTEND=noninteractive

COPY --from=build /app/main main
COPY config.json config.json

CMD ./main