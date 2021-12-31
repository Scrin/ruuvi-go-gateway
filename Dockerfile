FROM golang:1.17-alpine as build

ARG VERSION="unknown-docker"

WORKDIR /go/src/github.com/Scrin/ruuvi-go-gateway/
COPY . ./
RUN go install -v -ldflags "-X github.com/Scrin/ruuvi-go-gateway/common/version.Version=${VERSION}" ./cmd/ruuvi-go-gateway

FROM alpine

COPY --from=build /go/bin/ruuvi-go-gateway /usr/local/bin/ruuvi-go-gateway

CMD ["ruuvi-go-gateway"]
