package static

import (
	"fmt"
	"sync"
	"time"
)

type Cache struct {
	data sync.Map //map[string]cacheItem
}
type cacheItem struct {
	Val       string
	Ttl       time.Duration
	CreatedAt time.Time
}

func New() *Cache {
	return &Cache{}
}

func (s *Cache) Cache(key, val string, ttl time.Duration) (err error) {
	item := &cacheItem{
		Val:       val,
		Ttl:       ttl,
		CreatedAt: time.Now(),
	}
	s.data.Store(key, item)
	return
}

func (s *Cache) Cached(key string) (val string, hit bool, err error) {
	var item *cacheItem
	var v any
	v, hit = s.data.Load(key)
	if !hit {
		return
	}
	item = v.(*cacheItem)
	if item.Ttl > 0 && item.CreatedAt.Add(item.Ttl).Before(time.Now()) {
		hit = false
		s.data.Delete(key)
		return
	}

	val = item.Val
	return
}

func (s *Cache) Remove(key string) (err error) {
	if _, ok := s.data.Load(key); ok {
		s.data.Delete(key)
	}
	return nil
}

type StaticMapCache struct {
	data sync.Map // map[string]mapCacheItem
}
type mapCacheItem struct {
	Val       sync.Map //map[string]string
	Ttl       time.Duration
	CreatedAt time.Time
}

func NewMap() *StaticMapCache {
	return &StaticMapCache{}
}

func (s *StaticMapCache) Cache(key string, attrs map[string]interface{}, ttl time.Duration) (err error) {
	item := &mapCacheItem{
		Ttl:       ttl,
		CreatedAt: time.Now(),
	}
	for k, v := range attrs {
		item.Val.Store(k, fmt.Sprint(v))
	}
	s.data.Store(key, item)
	return
}

func (s *StaticMapCache) Cached(key string) (val map[string]string, hit bool, err error) {
	var item *mapCacheItem
	var v any
	val = make(map[string]string)
	v, hit = s.data.Load(key)
	if !hit {
		return
	}
	item = v.(*mapCacheItem)
	if item.Ttl > 0 && item.CreatedAt.Add(item.Ttl).Before(time.Now()) {
		hit = false
		s.data.Delete(key)
		return
	}

	item.Val.Range(func(key, value any) bool {
		kk := key.(string)
		vv := value.(string)
		val[kk] = vv
		return true
	})
	return
}

func (s *StaticMapCache) CacheAttr(key, attr, val string) (err error) {
	var item *mapCacheItem
	v, hit := s.data.Load(key)
	if !hit {
		return
	}
	item = v.(*mapCacheItem)
	if item.Ttl > 0 && item.CreatedAt.Add(item.Ttl).Before(time.Now()) {
		hit = false
		s.data.Delete(key)
		return
	}

	item.Val.Store(attr, val)
	return
}

func (s *StaticMapCache) CachedAttr(key, attr string) (val string, hit bool, err error) {
	var item *mapCacheItem
	var v any
	v, hit = s.data.Load(key)
	if !hit {
		return
	}
	item = v.(*mapCacheItem)
	if item.Ttl > 0 && item.CreatedAt.Add(item.Ttl).Before(time.Now()) {
		hit = false
		s.data.Delete(key)
		return
	}

	var value any
	value, hit = item.Val.Load(attr)
	if hit {
		val = value.(string)
	}
	return
}

func (s *StaticMapCache) CountAttr(key string) (num int64, err error) {
	if v, ok := s.data.Load(key); ok {
		item := v.(*mapCacheItem)
		item.Val.Range(func(key, value any) bool {
			num++
			return true
		})
	}
	return
}

func (s *StaticMapCache) RemoveAttr(key, attr string) (err error) {
	if v, ok := s.data.Load(key); ok {
		item := v.(*mapCacheItem)
		item.Val.Delete(attr)
	}
	return nil
}

func (s *StaticMapCache) Remove(key string) (err error) {
	if _, ok := s.data.Load(key); ok {
		s.data.Delete(key)
	}
	return nil
}

func (s *StaticMapCache) Expire(key string, ttl time.Duration) (err error) {
	if v, ok := s.data.Load(key); ok {
		item := v.(*mapCacheItem)
		item.Ttl = ttl
		s.data.Store(key, item)
	}
	return nil
}

func (s *StaticMapCache) Exist(key string) (hit bool, err error) {
	_, ok := s.data.Load(key)
	return ok, nil
}
