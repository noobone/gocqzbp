FROM golang:alpine3.15 AS build

ENV BUILD_PATH /go/src/github.com/noobone/cqhttp

RUN mkdir -p $BUILD_PATH 
WORKDIR $BUILD_PATH
COPY . $BUILD_PATH

FROM alpine:latest

ENV BUILD_PATH /go/src/github.com/noobone/cqhttp

COPY --from=build $BUILD_PATH /usr/bin/cqhttp

RUN chmod +x /usr/bin/cqhttp \
  && apk update \
  && apk add --no-cache ffmpeg

WORKDIR /data

CMD cqhttp
