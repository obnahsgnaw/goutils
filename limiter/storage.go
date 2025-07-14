package limiter

import (
	"github.com/obnahsgnaw/goutils/cacheutil"
	"strconv"
	"time"
)

type Storage struct {
	c func() cacheutil.Cache
}

func NewStorage(c func() cacheutil.Cache) *Storage {
	return &Storage{
		c: c,
	}
}

func (r *Storage) Exists(key string) (bool, error) {
	_, hit, err := r.c().Cached(key)

	return hit, err
}

func (r *Storage) Del(key string) error {
	if err := r.c().Remove(key); err != nil {
		return err
	}
	return nil
}

func (r *Storage) Set(key string, val int, ttl time.Duration) error {
	return r.c().Cache(key, strconv.Itoa(val), ttl)
}

func (r *Storage) Get(key string) (int, error) {
	val, hit, err := r.c().Cached(key)
	if err != nil {
		return 0, err
	}
	if !hit {
		return 0, err
	}

	return strconv.Atoi(val)
}
