package frequencer

import (
	"github.com/obnahsgnaw/goutils/limitutil/limiter"
	"time"
)

type Frequency struct {
	storage  *limiter.Storage
	prefix   string
	interval int64
}

func New(s *limiter.Storage, prefix string, interval int64) *Frequency {
	return &Frequency{
		storage:  s,
		prefix:   prefix,
		interval: interval,
	}
}

// SetStorage 设置存储
func (f *Frequency) SetStorage(storage *limiter.Storage) {
	f.storage = storage
}

func (f *Frequency) Attempt(target string) (intervalLeft int64, err error) {
	if f.storage == nil || f.interval < 0 {
		return
	}
	var lastTime int
	if lastTime, err = f.storage.Get(f.target(target)); err != nil {
		return
	}
	if lastTime == 0 {
		return
	}
	intervalSpace := time.Now().Unix() - int64(lastTime)
	if intervalSpace < f.interval {
		intervalLeft = f.interval - intervalSpace
	}
	return
}

func (f *Frequency) Hit(target string) error {
	if f.storage == nil || f.interval < 0 {
		return nil
	}
	return f.storage.Set(f.target(target), int(time.Now().Unix()), time.Second*time.Duration(f.interval))
}

func (f *Frequency) target(key string) string {
	return f.prefix + "frequency:" + key
}
