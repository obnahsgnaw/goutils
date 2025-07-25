package limiter

import (
	"github.com/obnahsgnaw/goutils/cacheutil"
	"strconv"
	"sync"
	"time"
)

type Storage struct {
	c    func() cacheutil.Cache
	ins  cacheutil.Cache
	once sync.Once
}

func NewStorage(c func() cacheutil.Cache) *Storage {
	return &Storage{
		c: c,
	}
}

func (r *Storage) cache() cacheutil.Cache {
	r.once.Do(func() {
		r.ins = r.c()
	})
	return r.ins
}

func (r *Storage) Exists(key string) (bool, error) {
	_, hit, err := r.cache().Cached(key)

	return hit, err
}

func (r *Storage) Del(key string) error {
	if err := r.cache().Remove(key); err != nil {
		return err
	}
	return nil
}

func (r *Storage) Set(key string, val int, ttl time.Duration) error {
	return r.cache().Cache(key, strconv.Itoa(val), ttl)
}

func (r *Storage) Get(key string) (int, error) {
	val, hit, err := r.cache().Cached(key)
	if err != nil {
		return 0, err
	}
	if !hit {
		return 0, err
	}

	return strconv.Atoi(val)
}
