API_POST_PROCESS=cd ../../api/$(APP_NAME) && ls v1/*.pb.go | xargs -n1 -IX bash -c 'sed s/,omitempty// X > X.tmp && mv X{.tmp,}'
include ../../Makefile.inc

ifeq ($(DB),sqlite)
    DB_DOCKER_CMD = 
		DB_ENV = -e BI_DB_DRIVER=sqlite
endif
ifeq ($(DB),mysql)
    DB_DOCKER_CMD = docker run -d --rm --name mysql -p 3306:3306 -e MYSQL_ROOT_PASSWORD=P@ssw0rd -e MYSQL_USER=bingo -e MYSQL_PASSWORD=P@ssw0rd -e MYSQL_DATABASE=bingo mysql:5.7 2> /dev/null || true
		DB_ENV = -e BI_DB_DRIVER=mysql -e BI_DB_HOST=host.docker.internal -e BI_DB_PORT=3306 -e BI_DB_DATABASE=bingo -e BI_DB_USERNAME=bingo -e BI_DB_PASSWORD=P@ssw0rd
endif
ifeq ($(DB),postgres)
    DB_DOCKER_CMD = docker run -d --rm --name postgres -p 5432:5432 -e POSTGRES_PASSWORD=P@ssw0rd -e POSTGRES_USER=bingo -e POSTGRES_DB=bingo postgres:14.2-alpine 2> /dev/null || true
		DB_ENV = -e BI_DB_DRIVER=postgres -e BI_DB_HOST=host.docker.internal -e BI_DB_PORT=5432 -e BI_DB_DATABASE=bingo -e BI_DB_USERNAME=bingo -e BI_DB_PASSWORD=P@ssw0rd
endif
ifeq ($(DB),sqlserver)
    DB_DOCKER_CMD = docker run -d --rm --name sqlserver -p 1433:1433 -e ACCEPT_EULA=Y -e SA_PASSWORD=P@ssw0rd mcr.microsoft.com/mssql/server:2017-latest  2> /dev/null || true && \
	                  docker exec -it sqlserver /opt/mssql-tools/bin/sqlcmd -S localhost -U sa -P P@ssw0rd -Q "CREATE DATABASE bingo"
		DB_ENV = -e BI_DB_DRIVER=sqlserver -e BI_DB_HOST=host.docker.internal -e BI_DB_PORT=1433 -e BI_DB_DATABASE=bingo -e BI_DB_USERNAME=sa -e BI_DB_PASSWORD=P@ssw0rd
endif
ifeq ($(DB),mongo)
    DB_DOCKER_CMD = docker run -d --rm --name mongo -p 27017:27017 -e MONGO_INITDB_ROOT_USERNAME=bingo -e MONGO_INITDB_ROOT_PASSWORD=P@ssw0rd -e MONGO_INITDB_DATABASE=bingo mongo:5.0.6 2> /dev/null || true
		DB_ENV = -e BI_DB_DRIVER=mongo -e BI_DB_HOST=host.docker.internal -e BI_DB_PORT=27017 -e BI_DB_DATABASE=bingo -e BI_DB_USERNAME=bingo -e BI_DB_PASSWORD=P@ssw0rd
endif

.PHONY: up
up:
	$(DB_DOCKER_CMD)
	docker run -d --rm --name jaeger -p 16686:16686 -p 14268:14268 jaegertracing/all-in-one:1.34 2> /dev/null || true
	docker run -d --rm --name etcd -p 2379:2379 -e ALLOW_NONE_AUTHENTICATION=yes bitnami/etcd:3.5.3 2> /dev/null || true
	docker run -d --rm --name redis -p 6379:6379 redislabs/rebloom:2.2.9 2> /dev/null || true
	docker run -d --rm --name es -p 9200:9200 -p 9300:9300 -e discovery.type=single-node -e xpack.security.enabled=false -e ES_JAVA_OPTS="-Xms512m -Xmx512m" docker.elastic.co/elasticsearch/elasticsearch:7.17.2 2> /dev/null || true
	docker run -d --rm --name rabbitmq -p 15672:15672 -p 5672:5672 -e RABBITMQ_DEFAULT_USER= -e RABBITMQ_DEFAULT_PASS= rabbitmq:3.9-management-alpine 2> /dev/null || true
	docker run -d --rm --name gowitness -p 7171:7171 huangyingting/gowitness gowitness server --address 0.0.0.0:7171 2> /dev/null || true
	docker run -d --rm --name bg -p 8000:8000 -e BG_JAEGER_ADDR=http://host.docker.internal:14268/api/traces huangyingting/bg 2> /dev/null || true
	sleep 10s
	docker run -d --rm --name bi -p 8001:8001 $(DB_ENV) -e BI_AMQP_URI=amqp://host.docker.internal:5672/ -e BI_REMOTE_BG_HTTP_ADDR=host.docker.internal:8000 -e BI_JAEGER_ADDR=http://host.docker.internal:14268/api/traces huangyingting/bi 2> /dev/null || true
	docker run -d --rm --name be -p 8002:8002 -e BE_JAEGER_ADDR=http://host.docker.internal:14268/api/traces huangyingting/be 2> /dev/null || true
	rm -f configs/config.yaml
	cp configs/config.$(DB) configs/config.yaml
	cd cmd/server && go run .


.PHONY: down
down:
	docker stop be 2> /dev/null || true
	docker stop bi 2> /dev/null || true
	docker stop bg 2> /dev/null || true
	docker stop gowitness 2> /dev/null | true
	docker stop rabbitmq 2> /dev/null | true
	docker stop es 2> /dev/null | true
	docker stop redis 2> /dev/null || true
	docker stop etcd 2> /dev/null || true
	docker stop $(DB) 2> /dev/null || true
	docker stop jaeger 2> /dev/null || true
	rm -f configs/config.yaml
	cp configs/config.sqlite configs/config.yaml	