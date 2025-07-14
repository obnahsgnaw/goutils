package limiter

import (
	"github.com/obnahsgnaw/goutils/cacheutil"
	"github.com/obnahsgnaw/goutils/cacheutil/static"
)

var limiter *Limiter

func init() {
	limiter = New(NewStorage(func() cacheutil.Cache {
		return static.New()
	}), "dft")
}

// SetStorage 设置默认存储
func SetStorage(storage *Storage) {
	limiter.SetStorage(storage)
}

// Default 返回默认 limiter
func Default() *Limiter {
	return limiter
}

// Attempt 快捷调用方式
func Attempt(item Item) (bool, error) {
	return limiter.Attempt(item)
}

// Hit 快捷调用方式
func Hit(item Item) error {
	return limiter.Hit(item)
}
