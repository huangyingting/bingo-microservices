# build
FROM golang:1.19.3-alpine3.16 as builder
ARG VERSION
COPY go/. /build
WORKDIR /build/app/bg
RUN mkdir -p bin/ && \
	GOOS=linux GOFLAGS=-buildvcs=false go build -ldflags "-s -w -X main.Version=$VERSION" -o ./bin/bg ./cmd/server && \
	cp configs/config.yaml ./bin/config.yaml
# create docker image
FROM redis:6.2.7-alpine3.16
WORKDIR /app
COPY go/app/bg/docker/start.sh .
COPY go/app/bg/docker/dump.rdb .
COPY --from=builder /build/app/bg/bin . 
RUN chmod 755 /app/start.sh
EXPOSE 8000
EXPOSE 9000
ENTRYPOINT ["/app/start.sh"]