package main

import (
	"github.com/obnahsgnaw/goutils/signalutil"
	"log"
)

func main() {
	log.Print("Server started and serving...")
	signalutil.Listen()
	log.Print("Server done")
}
