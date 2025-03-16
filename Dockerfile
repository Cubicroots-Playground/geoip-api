FROM golang:1.24.1-alpine3.21 as builder
ARG VERSION="development"

WORKDIR /run

COPY ./ ./
RUN go mod download
RUN go build -o /run ./

FROM alpine:3.21
COPY --from=builder /run/geoip-api /run/
WORKDIR /run

CMD ["/run/geoip-api"]