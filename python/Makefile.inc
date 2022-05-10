VERSION=$(shell git describe --tags --always)
APP_NAME=$(shell basename $(CURDIR))
DOCKER_IMAGE=huangyingting/$(APP_NAME):latest

.PHONY: docker
docker:
	cd ../.. && docker build -f python/$(APP_NAME)/Dockerfile --build-arg VERSION=$(VERSION) -t $(DOCKER_IMAGE) .
	