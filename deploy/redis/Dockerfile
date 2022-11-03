# Build a customized bitnami redis image with redisbloom module
FROM redislabs/rebloom:2.2.9 AS rebloom
FROM bitnami/redis:6.2.7
RUN mkdir /opt/bitnami/redis/modules && \
    echo loadmodule /opt/bitnami/redis/modules/redisbloom.so >> /opt/bitnami/redis/etc/redis.conf
COPY --from=rebloom /usr/lib/redis/modules/redisbloom.so* /opt/bitnami/redis/modules
