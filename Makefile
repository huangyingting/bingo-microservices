PKG_PROTO_FILES=$(shell find go/pkg -name *.proto)

.PHONY: go-pkg
go-pkg:
	protoc --proto_path=./go/pkg \
	       --proto_path=./go/third_party \
 	       --go_out=paths=source_relative:./go/pkg \
	       $(PKG_PROTO_FILES)

.PHONY: go-api
go-api:
	find go/app -mindepth 1 -maxdepth 1 -type d -print | xargs -L 1 bash -c 'cd "$$0" && pwd && $(MAKE) api'

.PHONY: go-internal
go-internal: go-pkg
	find go/app -mindepth 1 -maxdepth 1 -type d -print | xargs -L 1 bash -c 'cd "$$0" && pwd && $(MAKE) internal'

.PHONY: go-docker
go-docker: go-pkg go-api go-internal
	find go/app -mindepth 1 -maxdepth 1 -type d -print | xargs -L 1 bash -c 'cd "$$0" && pwd && $(MAKE) docker'

.PHONY: py-api
py-api:
	find python -mindepth 1 -maxdepth 1 -type d -print | xargs -L 1 bash -c 'cd "$$0" && pwd && $(MAKE) api'

.PHONY: py-docker
py-docker: py-api
	find python -mindepth 1 -maxdepth 1 -type d -print | xargs -L 1 bash -c 'cd "$$0" && pwd && $(MAKE) docker'

.PHONY: internal
api: go-internal

.PHONY: api
api: go-api py-api

.PHONY: docker
docker: go-docker py-docker

.PHONY: push
push: docker
	docker push huangyingting/bg
	docker push huangyingting/be
	docker push huangyingting/bi
	docker push huangyingting/bs

.PHONY: up
up:
	docker-compose up -d

.PHONY: down
down:
	docker-compose down
