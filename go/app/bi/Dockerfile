FROM golang:1.18-alpine3.15 as builder
ARG VERSION
COPY go/. /build
WORKDIR /build/app/bi
RUN mkdir -p bin/ && \
  apk add build-base && \
	CGO_ENABLED=1 GOOS=linux GOFLAGS=-buildvcs=false go build -ldflags "-s -w -X main.Version=$VERSION" -o ./bin/bi ./cmd/server && \
	cp configs/config.yaml ./bin/config.yaml

FROM alpine:3.15
WORKDIR /app
COPY --from=builder /build/app/bi/bin . 
EXPOSE 8001
EXPOSE 9001
CMD ["/app/bi", "-conf", "/app/config.yaml"]