package lockerutil

import (
	"time"
)

type Builder interface {
	Compete(key string, ttl time.Duration) Locker
}

type Locker interface {
	Error() error
	Exist() bool
	Unlock()
}
