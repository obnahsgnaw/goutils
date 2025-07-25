package frequencer

import (
	"github.com/obnahsgnaw/goutils/cacheutil"
	"github.com/obnahsgnaw/goutils/cacheutil/static"
	"github.com/obnahsgnaw/goutils/limitutil/limiter"
)

var f *Frequency

func init() {
	f = New(limiter.NewStorage(func() cacheutil.Cache {
		return static.New()
	}), "dft", 60)
}

func Default() *Frequency {
	return f
}

// SetStorage 设置默认存储
func SetStorage(storage *limiter.Storage) {
	f.SetStorage(storage)
}

func Attempt(target string) (intervalLeft int64, err error) {
	return f.Attempt(target)
}

func Hit(target string) error {
	return f.Hit(target)
}
