package runtimeutil

import (
	"fmt"
	"runtime/debug"
)

func HandleRecover(handler func(errMsg, stack string)) {
	if err := recover(); err != nil {
		s := string(debug.Stack())
		e := fmt.Sprintf("%v", err)
		if handler != nil {
			handler(e, s)
		}
	}
}
