package gracefulshutdown

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func GracefulShutdown(callback func()) {
	go func(callback func()) {
		termChan := make(chan os.Signal, 1)
		signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)

		<-termChan
		log.Println("TERMINATING")
		callback()
	}(callback)
}
