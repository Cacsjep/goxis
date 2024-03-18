package acap

import (
	"os"
	"os/signal"
	"syscall"
)

func SignalHandler(handler func()) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		<-sigs
		handler()
		os.Exit(0)
	}()
}
