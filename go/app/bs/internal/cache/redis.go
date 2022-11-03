package cache

import (
	"context"
	"time"

	"bingo/app/bs/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
)

type Redis struct {
	rc               *redis.Client
	rcc              *redis.ClusterClient
	sentinel_enabled bool
	ttl              uint32
	h                *log.Helper
}

func NewRedis(c *conf.Cache, h *log.Helper) (*Redis, error) {
	var rc *redis.Client = nil
	var err error = nil
	if !c.SentinelEnabled && c.Addr != "" {
		h.Debugf("redis addr: %s", c.Addr)

		rc = redis.NewClient(&redis.Options{
			Addr:     c.Addr,
			Username: c.Username,
			Password: c.Password,
		})
		if err = rc.Ping(context.Background()).Err(); err == nil {
			return &Redis{rc: rc, rcc: nil, sentinel_enabled: false, ttl: c.CacheTtl, h: h}, nil
		}
	}

	if c.SentinelEnabled && len(c.SentinelAddrs) > 0 {
		rcc := redis.NewFailoverClusterClient(&redis.FailoverOptions{
			MasterName:       c.SentinelMasterSet,
			SentinelAddrs:    c.SentinelAddrs,
			SentinelUsername: c.SentinelUsername,
			SentinelPassword: c.SentinelPassword,
			Username:         c.Username,
			Password:         c.Password,
			RouteRandomly:    true,
		})
		if err = rcc.Ping(context.Background()).Err(); err == nil {
			return &Redis{rc: nil, rcc: rcc, sentinel_enabled: true, ttl: c.CacheTtl, h: h}, nil
		}
	}

	h.Errorf("new redis client error: %v", err)
	return nil, err
}

func (cs *Redis) Get(key string, val interface{}) error {
	var err error = nil
	if cs.sentinel_enabled {
		err = cs.rcc.Get(context.Background(), key).Scan(val)
	} else {
		err = cs.rc.Get(context.Background(), key).Scan(val)
	}
	if err != nil {
		cs.h.Debugf("get error: %v", err)
	}
	return err
}

func (cs *Redis) Set(key string, val interface{}) error {
	var err error = nil
	if cs.sentinel_enabled {
		err = cs.rcc.Set(context.Background(), key, val, time.Duration(cs.ttl)*time.Second).Err()
	} else {
		err = cs.rc.Set(context.Background(), key, val, time.Duration(cs.ttl)*time.Second).Err()
	}
	if err != nil {
		cs.h.Debugf("set error: %v", err)
	}
	return err
}

func (cs *Redis) Delete(key string) error {
	var err error = nil
	if cs.sentinel_enabled {
		err = cs.rcc.Del(context.Background(), key).Err()
	} else {
		err = cs.rc.Del(context.Background(), key).Err()
	}
	if err != nil {
		cs.h.Debugf("delete error: %v", err)
	}
	return err
}

func (cs *Redis) BFAdd(key string, val string) (bool, error) {
	var err error = nil
	var added bool = false
	if cs.sentinel_enabled {
		added, err = cs.rcc.Do(context.Background(), "BF.ADD", key, val).Bool()
	} else {
		added, err = cs.rc.Do(context.Background(), "BF.ADD", key, val).Bool()
	}
	if err != nil {
		cs.h.Debugf("bf.add error: %v", err)
	}
	return added, err
}

func (cs *Redis) BFExists(key string, val string) (bool, error) {
	var err error = nil
	var exists bool = false
	if cs.sentinel_enabled {
		exists, err = cs.rcc.Do(context.Background(), "BF.EXISTS", key, val).Bool()
	} else {
		exists, err = cs.rc.Do(context.Background(), "BF.EXISTS", key, val).Bool()
	}
	if err != nil {
		cs.h.Debugf("bf.exists error: %v", err)
	}
	return exists, err
}
