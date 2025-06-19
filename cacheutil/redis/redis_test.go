package redis

import (
	"github.com/obnahsgnaw/pbhttp/core/application"
	"github.com/obnahsgnaw/pbhttp/core/config"
	"testing"
	"time"
)

func TestRedisCache(t *testing.T) {
	c, err := application.NewRedis(&config.Cache{
		Hostname: "127.0.0.1",
		Port:     6379,
		Password: "20210606123456",
	})
	if err != nil {
		t.Errorf("redis connect failed, err=%s", err.Error())
		return
	}
	s := New(c, func(s string) string {
		return s
	})

	err = s.Cache("a", "ok", time.Second*5)
	if err != nil {
		t.Errorf("cache failed, err=%s", err.Error())
		return
	}

	time.Sleep(time.Second * 2)

	v, hit, err1 := s.Cached("a")
	if err1 != nil {
		t.Errorf("cached failed, err=%s", err1.Error())
		return
	}
	if hit == false {
		t.Errorf("ttl 5 seconds, after 2 seconds hit need true but got false")
		return
	}
	if v != "ok" {
		t.Errorf("need got ok but got %s", v)
		return
	}

	time.Sleep(time.Second * 4)

	v, hit, _ = s.Cached("a")
	if hit == true {
		t.Errorf("ttl 5 seconds, after 5 seconds hit need false but got true")
		return
	}
	if v != "" {
		t.Errorf("need got empty but got %s", v)
		return
	}
}
