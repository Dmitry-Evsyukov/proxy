FROM golang:1.21-alpine AS builder

WORKDIR /source

ENV CGO_ENABLED 0
ENV GOOS linux

COPY internal internal
COPY cmd cmd
COPY go.mod go.sum ./
COPY certs certs

RUN go mod download  \
    && go mod tidy  \
    && go build -o proxy cmd/main.go

CMD ["/proxy"]
