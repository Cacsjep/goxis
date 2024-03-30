package acapapp

import (
	"os"
	"os/signal"
	"syscall"
)

func SignalHandler(handler func()) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT, syscall.SIGABRT)

	go func() {
		<-sigs
		handler()
		os.Exit(0)
	}()
}
