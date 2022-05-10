package cache

import (
	"context"
	"fmt"
	"time"

	"bingo/app/bs/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
)

type Redis struct {
	rdb *redis.Client
	ttl uint32
	h   *log.Helper
}

func NewRedis(c *conf.Cache, h *log.Helper) (*Redis, error) {
	var rdb *redis.Client = nil
	if c.Addr != "" {
		h.Debugf("redis addr: %s", c.Addr)
		rdb = redis.NewClient(&redis.Options{
			Addr: c.Addr,
		})
		if rdb.Ping(context.Background()).Err() == nil {
			return &Redis{rdb: rdb, ttl: c.CacheTtl, h: h}, nil
		}
	}
	h.Errorf("incorrect address or redis isn't alive: %s", c.Addr)
	return nil, fmt.Errorf("incorrect address or redis isn't alive: %s", c.Addr)
}

func (cs *Redis) Get(key string, val interface{}) error {
	err := cs.rdb.Get(context.Background(), key).Scan(val)
	if err != nil {
		cs.h.Debugf("get error: %v", err)
	}
	return err
}

func (cs *Redis) Set(key string, val interface{}) error {
	err := cs.rdb.Set(context.Background(), key, val, time.Duration(cs.ttl)*time.Second).Err()
	if err != nil {
		cs.h.Debugf("set error: %v", err)
	}
	return err
}

func (cs *Redis) Delete(key string) error {
	err := cs.rdb.Del(context.Background(), key).Err()
	if err != nil {
		cs.h.Debugf("delete error: %v", err)
	}
	return err
}

func (cs *Redis) BFAdd(key string, val string) (bool, error) {
	added, err := cs.rdb.Do(context.Background(), "BF.ADD", key, val).Bool()
	if err != nil {
		cs.h.Debugf("bf.add error: %v", err)
	}
	return added, err
}

func (cs *Redis) BFExists(key string, val string) (bool, error) {
	exists, err := cs.rdb.Do(context.Background(), "BF.EXISTS", key, val).Bool()
	if err != nil {
		cs.h.Debugf("bf.exists error: %v", err)
	}
	return exists, err
}
