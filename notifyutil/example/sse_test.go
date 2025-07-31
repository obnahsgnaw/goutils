package example

import (
	"github.com/obnahsgnaw/goutils/httputil/sse"
	"github.com/obnahsgnaw/goutils/notifyutil"
	"github.com/obnahsgnaw/goutils/notifyutil/channels"
	"log"
	"net/http"
	"testing"
)

func TestSse(t *testing.T) {
	s := notifyutil.DefaultManager()
	s.RegisterChannel(channels.SseGroupChannel, func() notifyutil.Channel {
		return channels.NewSseGroupChannel(sse.DefaultManager())
	})
	s.RegisterChannel(channels.SseTagChannel, func() notifyutil.Channel {
		return channels.NewSseTagChannel(sse.DefaultManager())
	})

	http.HandleFunc("/tag", func(w http.ResponseWriter, r *http.Request) {
		c := sse.DefaultManager().NewClient()
		c.AddTag("tag1")
		c.ServeHTTP(w, r)
	})
	http.HandleFunc("/group", func(w http.ResponseWriter, r *http.Request) {
		c := sse.DefaultManager().NewClient()
		c.JoinGroup("group1")
		c.ServeHTTP(w, r)
	})
	http.HandleFunc("/msg", func(w http.ResponseWriter, r *http.Request) {
		NewSseNotification().Notify()
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
