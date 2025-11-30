package kafka

import (
	"bytes"
	"context"
	"demo/internal/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	kafka "github.com/segmentio/kafka-go"
)

type KafkaInfo struct {
	Topic          string
	BrokkerAddress string
	GroupId        string
}

func NewKafka(ki *KafkaInfo) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{ki.BrokkerAddress},
		Topic:   ki.Topic,
		GroupID: ki.GroupId,
	})
	defer reader.Close()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	readMessages(ctx, reader)
}

func readMessages(ctx context.Context, r *kafka.Reader) {
	for {
		msg, err := r.ReadMessage(ctx)
		if err != nil {
			log.Fatal(err)
		}
		var order models.Order
		err = json.Unmarshal(msg.Value, &order)
		if err != nil {
			log.Printf("Error unmarshaling order: %v", err)
			continue
		}
		err = sendOrderToAPI(&order)
		if err != nil {
			log.Printf("Failed to send order to API: %v", err)
			continue
		}

		log.Printf("Order %s processed and sent to API", order.OrderUID)
		if err := r.CommitMessages(ctx, msg); err != nil {
			log.Printf("Failed to commit message: %v", err)
		}
	}
}

func sendOrderToAPI(order *models.Order) error {
	jsonData, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("failed to marshal order: %v", err)
	}

	resp, err := http.Post("http://localhost:8080/order", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to send POST request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
