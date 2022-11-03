#!/bin/sh
set -e

start_redis() {
  echo "start redis"
  redis-server &
}

start_bg()
{    
  echo "start bingo geo service"
  sleep 5s
  exec /app/bg -conf config.yaml
}

start_redis
start_bg