# Stage 1: Build binary
FROM golang:1.16-alpine as builder 

WORKDIR /build
COPY . .

RUN apk add --no-cache build-base

RUN go get ./... \
 && go test -v \
 && go build -o cf-ddns

# Stage 2: Release image
FROM alpine:3
WORKDIR /app
COPY --from=builder /build/cf-ddns /app/
