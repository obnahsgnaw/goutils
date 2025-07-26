package main

import (
	"github.com/obnahsgnaw/goutils/httputil/sse"
	"log"
	"net/http"
	"time"
)

func main() {
	s := sse.Default()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		c := s.NewClient()
		g := r.URL.Query().Get("group")
		u := r.URL.Query().Get("user")
		if g != "" {
			c.JoinGroup(g)
		}
		if u != "" {
			c.AddTag(u)
		}
		//c.Provider(func() *sse.Message {
		//	return sse.NewMessage("this is a test message")
		//}, time.Second)
		c.ServeHTTP(w, r)
	})
	go func() {
		for {
			time.Sleep(time.Second * 2)
			s.BroadcastGroup("group1", sse.NewMessage("this is a test group message"))
			s.BroadcastTag("user1", sse.NewMessage("this is a test tag message"))
		}
	}()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
