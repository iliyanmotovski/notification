package gracefulshutdown

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

var termChan = make(chan os.Signal, 1)

func GracefulShutdown(callback func()) {
	go func(callback func()) {
		signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)

		<-termChan
		log.Println("TERMINATING")
		callback()
	}(callback)
}
