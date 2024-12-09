package main

import (
	"github.com/obnahsgnaw/goutils/proxyutil"
	"github.com/obnahsgnaw/goutils/signalutil"
	"log"
	"net/http"
)

func main() {
	px, err := proxyutil.New("http://yt.eflyyzh.com/#/home")
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		px.ServeHTTP(w, r)
	})
	go func() {
		log.Fatal(http.ListenAndServe(":8072", nil))
	}()

	log.Println("proxy server started and serving...")
	signalutil.Listen()
	log.Println("proxy server stopped")
}
