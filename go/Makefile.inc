VERSION=$(shell git describe --tags --always)
APP_NAME=$(shell basename $(CURDIR))
INTERNAL_PROTO_FILES=$(shell find internal -name *.proto)
API_PROTO_FILES=$(shell cd ../../api/$(APP_NAME) && find . -name *.proto)
DOCKER_IMAGE=huangyingting/$(APP_NAME):latest

.PHONY: api
# generate api code
api: 
	 cd ../../api/$(APP_NAME) && protoc --proto_path=. \
           --proto_path=../../third_party \
           --go_out=paths=source_relative:. \
           --go-grpc_out=paths=source_relative:. \
 	         --go-http_out=paths=source_relative:. \
           --go-errors_out=paths=source_relative:. \
 	         --openapi_out==paths=source_relative:. \
           $(API_PROTO_FILES)
	 $(API_POST_PROCESS)

.PHONY: internal
# generate internal proto struct
internal:
	protoc --proto_path=. \
           --proto_path=../../pkg/rabbitmq \
           --proto_path=../../third_party \
           --go_out=paths=source_relative:. \
           $(INTERNAL_PROTO_FILES)

.PHONY: build
# build
build:
	mkdir -p bin/ && \
	go build -ldflags "-s -w -X main.Version=$(VERSION)" -o ./bin/$(APP_NAME) ./cmd/server && \
	cp configs/config.yaml ./bin/config.yaml

.PHONY: run
run:
	cd cmd/server/ && go run .

.PHONY: docker
docker: api internal
	cd ../../.. && docker build -f go/app/$(APP_NAME)/Dockerfile --build-arg VERSION=$(VERSION) -t $(DOCKER_IMAGE) .

.PHONY: all
# generate all
all: api internal build

# show help
help:
	@echo ''
	@echo 'Usage:'
	@echo ' make [target]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help