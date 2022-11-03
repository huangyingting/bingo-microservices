docker run -d --rm --name redis -p 6379:6379 redis:6.2.7-alpine3.15
go run .
rm -f ../../docker/dump.rdb
docker cp redis:/data/dump.rdb ../../docker/dump.rdb