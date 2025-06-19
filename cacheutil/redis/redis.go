package redis

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"time"
)

type Cache struct {
	client       *redis.Client
	timeout      time.Duration
	keyFormatter func(string) string
}

// IsRedisKeyNotExists redis key 不存在判断
func isRedisKeyNotExists(err error) bool {
	return errors.Is(err, redis.Nil)
}

func New(c *redis.Client, keyFormatter func(string) string) *Cache {
	return &Cache{client: c, timeout: time.Second * 5, keyFormatter: keyFormatter}
}

func (s *Cache) Cache(key, val string, ttl time.Duration) (err error) {
	ctx, cl := context.WithTimeout(context.Background(), s.timeout)
	defer cl()
	res := s.client.Set(ctx, s.initKey(key), val, ttl)
	return res.Err()
}

func (s *Cache) Cached(key string) (val string, hit bool, err error) {
	ctx, cl := context.WithTimeout(context.Background(), s.timeout)
	defer cl()
	res := s.client.Get(ctx, s.initKey(key))
	if res.Err() != nil {
		err = res.Err()
		if isRedisKeyNotExists(res.Err()) {
			err = nil
		}
		return
	}
	hit = true
	val = res.Val()
	return
}

func (s *Cache) Remove(key string) (err error) {
	ctx, cl := context.WithTimeout(context.Background(), s.timeout)
	defer cl()
	res := s.client.Del(ctx, s.initKey(key))
	if res.Err() != nil {
		err = res.Err()
		if isRedisKeyNotExists(res.Err()) {
			err = nil
		}
	}
	return
}

func (s *Cache) initKey(key string) string {
	if s.keyFormatter != nil {
		return s.keyFormatter(key)
	}

	return key
}

type RedisMapCache struct {
	client       *redis.Client
	timeout      time.Duration
	keyFormatter func(string) string
}

func NewMap(c *redis.Client, keyFormatter func(string) string) *RedisMapCache {
	return &RedisMapCache{client: c, timeout: time.Second * 5, keyFormatter: keyFormatter}
}

func (s *RedisMapCache) Cache(key string, attrs map[string]interface{}, ttl time.Duration) (err error) {
	if len(attrs) == 0 {
		return nil
	}
	ctx, cl := context.WithTimeout(context.Background(), s.timeout)
	defer cl()
	res := s.client.HSet(ctx, s.initKey(key), attrs)
	if res.Err() != nil {
		err = res.Err()
		return
	}
	res1 := s.client.Expire(ctx, s.initKey(key), ttl)
	return res1.Err()
}

func (s *RedisMapCache) Cached(key string) (val map[string]string, hit bool, err error) {
	ctx, cl := context.WithTimeout(context.Background(), s.timeout)
	defer cl()
	res := s.client.HGetAll(ctx, s.initKey(key))
	if res.Err() != nil {
		err = res.Err()
		if isRedisKeyNotExists(res.Err()) {
			err = nil
		}
		return
	}
	hit = true
	val = res.Val()
	return
}

func (s *RedisMapCache) CacheAttr(key, attr, val string) (err error) {
	ctx, cl := context.WithTimeout(context.Background(), s.timeout)
	defer cl()
	res := s.client.HSet(ctx, s.initKey(key), attr, val)
	if res.Err() != nil {
		err = res.Err()
		return
	}
	return nil
}

func (s *RedisMapCache) CachedAttr(key, attr string) (val string, hit bool, err error) {
	ctx, cl := context.WithTimeout(context.Background(), s.timeout)
	defer cl()
	res := s.client.HGet(ctx, s.initKey(key), attr)
	if res.Err() != nil {
		err = res.Err()
		if isRedisKeyNotExists(res.Err()) {
			err = nil
		}
		return
	}
	hit = true
	val = res.Val()
	return
}

func (s *RedisMapCache) CountAttr(key string) (num int64, err error) {
	ctx, cl := context.WithTimeout(context.Background(), s.timeout)
	defer cl()
	res := s.client.HLen(ctx, s.initKey(key))
	if res.Err() != nil {
		err = res.Err()
		if isRedisKeyNotExists(res.Err()) {
			err = nil
		}
		return
	}
	num = res.Val()
	return
}

func (s *RedisMapCache) RemoveAttr(key, attr string) (err error) {
	ctx, cl := context.WithTimeout(context.Background(), s.timeout)
	defer cl()
	res := s.client.HDel(ctx, s.initKey(key), attr)
	if res.Err() != nil {
		err = res.Err()
		if isRedisKeyNotExists(res.Err()) {
			err = nil
		}
		return
	}
	res1 := s.client.HLen(ctx, s.initKey(key))
	if res1.Err() != nil {
		return res1.Err()
	}
	if res1.Val() == 0 {
		return s.Remove(key)
	}
	return
}

func (s *RedisMapCache) Remove(key string) (err error) {
	ctx, cl := context.WithTimeout(context.Background(), s.timeout)
	defer cl()
	res := s.client.Del(ctx, s.initKey(key))
	if res.Err() != nil {
		err = res.Err()
		if isRedisKeyNotExists(res.Err()) {
			err = nil
		}
	}
	return
}

func (s *RedisMapCache) Expire(key string, ttl time.Duration) (err error) {
	ctx, cl := context.WithTimeout(context.Background(), s.timeout)
	defer cl()
	res1 := s.client.Expire(ctx, s.initKey(key), ttl)
	return res1.Err()
}

func (s *RedisMapCache) Exist(key string) (hit bool, err error) {
	ctx, cl := context.WithTimeout(context.Background(), s.timeout)
	defer cl()
	res1 := s.client.Exists(ctx, s.initKey(key))
	return res1.Val() > 0, res1.Err()
}

func (s *RedisMapCache) initKey(key string) string {
	if s.keyFormatter != nil {
		return s.keyFormatter(key)
	}

	return key
}
