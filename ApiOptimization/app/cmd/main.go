package main

import (
	"demo/internal/app"
	"demo/internal/cache"
	"demo/internal/config"
	"demo/internal/database"
	"demo/internal/handlers"
	"demo/internal/kafka"
	"demo/internal/repository"
	"demo/internal/service"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	go func() {
		log.Println("pprof on http://localhost:6060")
		_ = http.ListenAndServe(":6060", nil)
	}()
	cfg := config.GetConfig("config/dev.yml")
	db := database.NewDB(cfg.Database.PostgresURL)
	redisClient := cache.NewRedisClient(cfg.Redis.Address, cfg.Redis.Password)
	orderRepository := repository.NewOrderRepository(db.DB_CONN)
	orderService := service.NewOrderService(orderRepository, redisClient)
	orderHandler := handlers.NewOrderHandler(orderService)

	start := time.Now()
	fmt.Printf("done after %v\n", time.Since(start))
	app := app.NewApp(db, orderHandler, cfg)

	go func() {
		kafka.NewKafka(&kafka.KafkaInfo{
			BrokkerAddress: cfg.Kafka.BrokerAddress,
			Topic:          cfg.Kafka.Topic,
			GroupId:        cfg.Kafka.GroupId,
		})
	}()

	go func() {
		app.MustStart()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop

	app.Stop()
	log.Println("Gracefully stopped")
}
