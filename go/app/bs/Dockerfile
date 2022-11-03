# build react spa
FROM node:16-alpine3.16 AS react-build
ARG VERSION
WORKDIR /build
COPY js/bf/package.json js/bf/yarn.lock .
RUN yarn install --pure-lockfile
COPY js/bf/ .
RUN yarn build

# build bs
FROM golang:1.19.3-alpine3.16 as go-build
COPY go/. /build
WORKDIR /build/app/bs
RUN mkdir -p bin/ && \
  apk add build-base && \
	CGO_ENABLED=1 GOOS=linux GOFLAGS=-buildvcs=false go build -ldflags "-s -w -X main.Version=$VERSION" -o ./bin/bs ./cmd/server && \
	cp configs/config.yaml ./bin/config.yaml

# create docker image
FROM alpine:3.16
WORKDIR /app 
COPY --from=react-build /build/build website
COPY --from=go-build /build/app/bs/bin .
COPY js/bf/env-config.sh .
COPY js/bf/.env .

RUN apk --no-cache add bash sed

ENV PORT 8080
EXPOSE 8080
# CMD ["/app/bs", "-conf", "/app/config.yaml"]
CMD ["/bin/bash", "-c", "/app/env-config.sh /app/.env && mv /app/env-config.js /app/website && /app/bs -conf /app/config.yaml"]
