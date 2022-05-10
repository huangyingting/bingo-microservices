# build react spa
FROM node:16-alpine AS react-build
ARG VERSION
WORKDIR /build
COPY js/bf/package*.json ./
RUN npm i -g npm
RUN npm install --silent
COPY js/bf/ ./
RUN npm run build

# build bs
FROM golang:1.18-alpine as go-build
COPY go/. /build
WORKDIR /build/app/bs
RUN mkdir -p bin/ && \
  apk add build-base && \
	CGO_ENABLED=1 GOOS=linux GOFLAGS=-buildvcs=false go build -ldflags "-s -w -X main.Version=$VERSION" -o ./bin/bs ./cmd/server && \
	cp configs/config.yaml ./bin/config.yaml

# create docker image
FROM alpine:3.15
WORKDIR /app 
COPY --from=react-build /build/build website
COPY --from=go-build /build/app/bs/bin . 
ENV PORT 8080
EXPOSE 8080
CMD ["/app/bs", "-conf", "/app/config.yaml"]