FROM golang:1.10.2-alpine3.7 AS build
RUN apk --no-cache add gcc g++ make ca-certificates
WORKDIR /go/src/github.com/bialas1993/go-notifier

COPY Gopkg.lock Gopkg.toml ./
COPY vendor vendor
COPY util util
COPY event event
COPY db db
COPY search search
COPY schema schema
COPY notify-service notify-service
COPY query-service query-service
COPY pusher-service pusher-service

RUN go install ./...

FROM alpine:3.7
WORKDIR /usr/bin
COPY --from=build /go/bin .
