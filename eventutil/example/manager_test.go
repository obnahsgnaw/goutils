package example

import (
	"github.com/obnahsgnaw/goutils/eventutil"
	"log"
	"testing"
)

func TestEvent(t *testing.T) {
	b := eventutil.DefaultManger()
	b.Register("test", func(e *eventutil.Event) {
		log.Println("test event:", e.Data)
	})
	e := b.Build("test", 1, 2, 3)
	e.Fire()
}
