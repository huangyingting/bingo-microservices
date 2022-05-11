VERSION=$(shell git describe --tags --always)
APP_NAME=$(shell basename $(CURDIR))
DOCKER_IMAGE=huangyingting/$(APP_NAME):latest
API_PROTO_FILES=$(shell cd ../../go/api/$(APP_NAME) && find . -name *.proto)

.PHONY: docker
docker: api
	cd ../.. && docker build -f python/$(APP_NAME)/Dockerfile --build-arg VERSION=$(VERSION) -t $(DOCKER_IMAGE) .

.PHONY: api
# generate api code
api:
	 cd ../../go/api/$(APP_NAME) && protoc --proto_path=. \
           --proto_path=../../third_party \
           --go_out=paths=source_relative:. \
           --go-grpc_out=paths=source_relative:. \
 	         --go-http_out=paths=source_relative:. \
           --go-errors_out=paths=source_relative:. \
 	         --openapi_out==paths=source_relative:. \
           $(API_PROTO_FILES)
