package limiter

import (
	"github.com/obnahsgnaw/goutils/cacheutil"
	"github.com/obnahsgnaw/goutils/cacheutil/static"
)

var l *Limiter

func init() {
	l = New(NewStorage(func() cacheutil.Cache {
		return static.New()
	}), "dft")
}

// SetStorage 设置默认存储
func SetStorage(storage *Storage) {
	l.SetStorage(storage)
}

// Default 返回默认 limiter
func Default() *Limiter {
	return l
}

// Attempt 快捷调用方式
func Attempt(item Item) (bool, error) {
	return l.Attempt(item)
}

// Hit 快捷调用方式
func Hit(item Item) error {
	return l.Hit(item)
}
