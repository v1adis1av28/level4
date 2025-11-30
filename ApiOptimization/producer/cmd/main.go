package main

import (
	"log"
	"os"
	"os/signal"
	"producer/internal/service"
	"syscall"
	"time"
)

func main() {
	go func() {
		time.Sleep(10 * time.Second)
		service.StartSender()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop

	log.Println("service stopped")
}
