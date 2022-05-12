include ../../Makefile.inc
ifeq ($(DB),sqlite)
    DB_DOCKER_CMD =
endif
ifeq ($(DB),mysql)
    DB_DOCKER_CMD = docker run -d --rm --name mysql -p 3306:3306 -e MYSQL_ROOT_PASSWORD=P@ssw0rd -e MYSQL_USER=bingo -e MYSQL_PASSWORD=P@ssw0rd -e MYSQL_DATABASE=bingo mysql:5.7 2> /dev/null || true
endif
ifeq ($(DB),postgres)
    DB_DOCKER_CMD = docker run -d --rm --name postgres -p 5432:5432 -e POSTGRES_PASSWORD=P@ssw0rd -e POSTGRES_USER=bingo -e POSTGRES_DB=bingo postgres:14.2-alpine 2> /dev/null || true
endif
ifeq ($(DB),sqlserver)
    DB_DOCKER_CMD = docker run -d --rm --name sqlserver -p 1433:1433 -e ACCEPT_EULA=Y -e SA_PASSWORD=P@ssw0rd mcr.microsoft.com/mssql/server:2017-latest  2> /dev/null || true && \
	                  docker exec -it sqlserver /opt/mssql-tools/bin/sqlcmd -S localhost -U sa -P P@ssw0rd -Q "CREATE DATABASE bingo"
endif
ifeq ($(DB),mongo)
    DB_DOCKER_CMD = docker run -d --rm --name mongo -p 27017:27017 -e MONGO_INITDB_ROOT_USERNAME=bingo -e MONGO_INITDB_ROOT_PASSWORD=P@ssw0rd -e MONGO_INITDB_DATABASE=bingo mongo:5.0.6 2> /dev/null || true
endif

.PHONY: up
up:
	rm configs/config.yaml
	cp configs/config.$(DB) configs/config.yaml
	docker run -d --rm --name rabbitmq -p 15672:15672 -p 5672:5672 -e RABBITMQ_DEFAULT_USER= -e RABBITMQ_DEFAULT_PASS= rabbitmq:3.9-management-alpine 2> /dev/null || true
	$(DB_DOCKER_CMD)

.PHONY: down
down:
	docker stop rabbitmq 2> /dev/null || true
	docker stop $(DB) 2> /dev/null || true