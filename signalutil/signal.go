package signalutil

import (
	"os"
	"os/signal"
	"syscall"
)

// Listen Listen os signal, default:syscall.SIGINT  syscall.SIGTERM syscall.SIGQUIT
func Listen(sig ...os.Signal) {
	if len(sig) == 0 {
		sig = append(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	}
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, sig...)
	<-quit
}
