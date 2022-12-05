ARG GO_VERSION=1.16.6

FROM golang:${GO_VERSION}-alpine AS builder

RUN go env -w GOPROXY=direct
RUN apk add --no-cache git
RUN apk --no-cache add ca-certificates && update-ca-certificates

WORKDIR /src

COPY ./go.mod ./go.sum ./
RUN go mod download


COPY contact-service contact-service
COPY database database
COPY events events
COPY means means
COPY models models
COPY nginx nginx
COPY pusher-service pusher-service
COPY query-service query-service
COPY repository repository
COPY search search

RUN go install ./...

FROM alpine:3.11
WORKDIR /usr/bin
COPY --from=builder /go/bin .