include ../../Makefile.inc

.PHONY: up
up:
	docker run -d --rm --name jaeger -p 16686:16686 -p 14268:14268 jaegertracing/all-in-one:1.39 2> /dev/null || true
	docker run -d --rm --name bg -p 8000:8000 -e BG_JAEGER_ADDR=http://host.docker.internal:14268/api/traces huangyingting/bg 2> /dev/null || true

.PHONY: down
down:
	docker stop bg 2> /dev/null || true
	docker stop jaeger 2> /dev/null || true