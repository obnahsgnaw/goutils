package limiter

// Rate limiter
// Usage:
// limiter.SetStorage
// limiter.Attempt()

import (
	"strings"
	"time"
)

type Item struct {
	Key    string               // 限流对象
	Max    func() int           // 最大次数
	Ttl    func() time.Duration // 过期时间
	Global bool                 // 全局 即非用户级别
}

// Limiter 限流器
type Limiter struct {
	storage *Storage
	prefix  string
}

// New 创建一个新的limiter
func New(s *Storage, prefix string) *Limiter {
	if prefix != "" && !strings.HasSuffix(prefix, ":") {
		prefix += ":"
	}
	return &Limiter{
		storage: s,
		prefix:  prefix,
	}
}

// SetStorage 设置存储
func (l *Limiter) SetStorage(storage *Storage) {
	l.storage = storage
}

// Attempt Check limit, max < 0 no limit, max=0 all limit, ttl <=0 no limit
func (l *Limiter) Attempt(item Item) (bool, error) {
	if l.storage == nil || item.Max() < 0 || item.Ttl() <= 0 {
		return true, nil
	}
	if item.Max() == 0 {
		return false, nil
	}
	exist, err := l.storage.Exists(l.timer(item.Key))
	if err != nil {
		return false, err
	}
	if !exist {
		if err = l.Reset(item.Key); err != nil {
			return false, err
		}
		return true, nil
	}

	if _, left, err := l.UsedAndLeft(item); err != nil {
		return false, err
	} else {
		return left > 0, nil
	}
}

// Hit 触发一次限流次数
func (l *Limiter) Hit(item Item) error {
	if l.storage == nil || item.Max() <= 0 || item.Ttl() <= 0 {
		return nil
	}
	itemTimer := l.timer(item.Key)
	exist, err := l.storage.Exists(itemTimer)
	if err != nil {
		return err
	}
	if !exist {
		if err = l.storage.Set(itemTimer, int(time.Now().Add(item.Ttl()).Unix()), item.Ttl()); err != nil {
			return err
		}
	}
	used, _, err := l.UsedAndLeft(item)
	if err != nil {
		return err
	}
	v := used + 1
	if v > item.Max() {
		v = item.Max()
	}
	return l.storage.Set(l.target(item.Key), v, item.Ttl())
}

// UsedAndLeft 已用和剩余次数
func (l *Limiter) UsedAndLeft(item Item) (used int, left int, err error) {
	used = 0
	left = 0
	if l.storage == nil || item.Max() <= 0 || item.Ttl() <= 0 {
		left = l.left(item, used)
		return
	}

	if used, err = l.storage.Get(l.target(item.Key)); err != nil {
		return
	}
	left = l.left(item, used)

	return
}

// AvailableAtAndIn 可用时间和可用的间隔
func (l *Limiter) AvailableAtAndIn(item Item) (at time.Time, in time.Duration, err error) {
	at = time.Now()
	in = 0

	if l.storage == nil || item.Max() < 0 || item.Ttl() <= 0 {
		return
	}
	ts, err := l.storage.Get(l.timer(item.Key))
	if err != nil || ts == 0 {
		return
	}
	_, left, err := l.UsedAndLeft(item)
	if err != nil || left > 0 {
		return
	}
	at = time.Unix(int64(ts), 0)
	in = at.Sub(time.Now())
	if in < 0 {
		in = 0
	}

	return
}

// Reset 重置次数
func (l *Limiter) Reset(key string) error {
	if l.storage == nil {
		return nil
	}
	return l.storage.Del(l.target(key))
}

// Clear 清除限流信息
func (l *Limiter) Clear(key string) error {
	if l.storage == nil {
		return nil
	}
	if err := l.storage.Del(l.timer(key)); err != nil {
		return err
	}

	if err := l.Reset(key); err != nil {
		return err
	}

	return nil
}

// 剩余次数
func (l *Limiter) left(item Item, used int) int {
	if l.storage == nil || item.Max() < 0 || item.Ttl() <= 0 {
		if item.Max() < 0 {
			return 0
		} else {
			return item.Max()
		}
	}
	if item.Max() == 0 {
		return 0
	}
	left := item.Max() - used
	if left < 0 {
		left = 0
	}

	return left
}

// 定时器key
func (l *Limiter) timer(key string) string {
	return l.prefix + "limiters:" + "timer:" + key
}

// 限制对象key
func (l *Limiter) target(key string) string {
	return l.prefix + "limiters:" + "times:" + key
}
