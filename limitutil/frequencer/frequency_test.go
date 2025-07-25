package frequencer

import (
	"github.com/obnahsgnaw/goutils/cacheutil"
	"github.com/obnahsgnaw/goutils/cacheutil/static"
	"github.com/obnahsgnaw/goutils/limitutil/limiter"
	"log"
	"testing"
	"time"
)

func TestNewFrequency(t *testing.T) {
	f := New(limiter.NewStorage(func() cacheutil.Cache {
		return static.New()
	}), "test", 60)

	interval, err := f.Attempt("a")
	if err != nil {
		t.Fatal(err)
		return
	}
	if interval != 0 {
		t.Fatal("need interval 0 but not")
		return
	}
	err = f.Hit("a")
	if err != nil {
		t.Fatal(err)
		return
	}
	interval, err = f.Attempt("a")
	if err != nil {
		t.Fatal(err)
		return
	}
	if interval == 0 {
		t.Fatal("need interval big 0 but not", interval)
	}
	time.Sleep(time.Second * 5)
	interval, err = f.Attempt("a")
	if err != nil {
		t.Fatal(err)
		return
	}
	if interval == 0 {
		t.Fatal("need interval big 0 but not", interval)
	}
	log.Println(interval)
}
