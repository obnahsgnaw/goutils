package local

import (
	"github.com/obnahsgnaw/goutils/lockerutil"
	"sync"
	"time"
)

type Builder struct {
}

func New() *Builder {
	return &Builder{}
}

func (b *Builder) Compete(key string, timeout time.Duration) lockerutil.Locker {
	return newLocker(key, timeout)
}

type Locker struct {
	l sync.Mutex
}

func newLocker(_ string, _ time.Duration) *Locker {
	s := &Locker{}
	s.l.Lock()
	return s
}

func (l *Locker) Error() error {
	return nil
}

func (l *Locker) Exist() bool {
	return false
}

func (l *Locker) Unlock() {
	l.l.Unlock()
}
