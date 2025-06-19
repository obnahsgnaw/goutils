package static

import (
	"testing"
	"time"
)

func TestStatic(t *testing.T) {
	s := New()

	_ = s.Cache("a", "ok", time.Second*5)

	time.Sleep(time.Second * 2)

	v, hit, _ := s.Cached("a")
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
