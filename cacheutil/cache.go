package cacheutil

import "time"

// Cache 通用缓存接口
type Cache interface {
	Cache(key, val string, ttl time.Duration) (err error)
	Cached(key string) (val string, hit bool, err error)
	Remove(key string) (err error)
}

type MapCache interface {
	Cache(key string, attrs map[string]interface{}, ttl time.Duration) (err error)
	Cached(key string) (val map[string]string, hit bool, err error)
	CacheAttr(key, attr, val string) (err error)
	CachedAttr(key, attr string) (val string, hit bool, err error)
	CountAttr(key string) (num int64, err error)
	Remove(key string) (err error)
	RemoveAttr(key, attr string) (err error)
	Expire(key string, ttl time.Duration) (err error)
	Exist(key string) (hit bool, err error)
}
