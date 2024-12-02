package main

import (
	"github.com/obnahsgnaw/goutils/runtimeutil"
	"log"
)

func main() {
	defer runtimeutil.HandleRecover(func(err, stack string) {
		log.Print("panic: ", err)
		log.Print("panic stack: ", stack)
	})

	panic("this is test")
}
