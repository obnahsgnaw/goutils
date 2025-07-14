package limiter

import (
	"github.com/obnahsgnaw/goutils/cacheutil"
	"github.com/obnahsgnaw/goutils/cacheutil/static"
	"testing"
	"time"
)

func TestLimiter_Attempt(t *testing.T) {
	c := static.New()
	s := NewStorage(func() cacheutil.Cache {
		return c
	})
	l := New(s, "")
	testItem := Item{
		Key: "xxx",
		Max: func() int {
			return 3
		},
		Ttl: func() time.Duration {
			return 60 * time.Second
		},
	}
	err := l.Hit(testItem)
	if err != nil {
		t.Error(err)
		return
	}
	used, left, err := l.UsedAndLeft(testItem)
	if err != nil {
		t.Error(err)
		return
	}
	if used != 1 {
		t.Error("used err", used)
		return
	}
	if left != 2 {
		t.Error("left err", left)
		return
	}

	enable, err := l.Attempt(testItem)
	if err != nil {
		t.Error(err)
		return
	}
	if enable != true {
		t.Error("attempt not true")
		return
	}

	err = l.Hit(testItem)
	if err != nil {
		t.Error(err)
		return
	}

	err = l.Hit(testItem)
	if err != nil {
		t.Error(err)
		return
	}

	used, left, err = l.UsedAndLeft(testItem)
	if err != nil {
		t.Error(err)
		return
	}
	if used != 3 {
		t.Error("used err1", used)
		return
	}
	if left != 0 {
		t.Error("left err1", left)
		return
	}

	enable, err = l.Attempt(testItem)
	if err != nil {
		t.Error(err)
		return
	}
	if enable != false {
		t.Error("attempt not false")
		return
	}
}
