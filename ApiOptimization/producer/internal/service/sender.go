package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"producer/internal/models"
	"time"

	"github.com/segmentio/kafka-go"
)

func StartSender() {
	topic := "orders"
	brokerAddress := "kafka:9092"

	writer := &kafka.Writer{
		Addr:         kafka.TCP(brokerAddress),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}
	defer writer.Close()

	for {
		order := generateOrder()
		jsonData, err := json.Marshal(order)
		if err != nil {
			log.Printf("Marshaling error")
			time.Sleep(30 * time.Second)
			continue
		}

		err = writer.WriteMessages(context.Background(), kafka.Message{
			Key:   []byte(order.OrderUID),
			Value: jsonData,
		})

		if err != nil {
			log.Printf("Sending message error")
		} else {
			log.Printf("Message send succesfully!")
		}

		time.Sleep(30 * time.Second)
	}
}

func generateOrder() models.Order {
	id := rand.Intn(1000) + 1
	orderUID := fmt.Sprintf("order-%d-%d", id, rand.Intn(9999))
	trackNumber := fmt.Sprintf("WBIL%d%04d", id, rand.Intn(9999))

	itemCount := rand.Intn(3) + 1
	var items []models.Item

	for i := 0; i < itemCount; i++ {
		chrtID := 9934000 + rand.Intn(1000)
		price := float64(rand.Intn(2000) + 500)
		sale := float64(rand.Intn(50))
		totalPrice := price * (100 - sale) / 100

		items = append(items, models.Item{
			ChrtID:      int64(chrtID),
			TrackNumber: trackNumber,
			Price:       price,
			RID:         fmt.Sprintf("ab%010d", rand.Intn(9999999999)),
			Name:        randomProduct(),
			Sale:        sale,
			Size:        randomSize(),
			TotalPrice:  totalPrice,
			NmID:        int64(2389212 + rand.Intn(1000)),
			Brand:       randomBrand(),
			Status:      202,
		})
	}

	order := models.Order{
		OrderUID:    orderUID,
		TrackNumber: trackNumber,
		Entry:       "WBIL",
		Delivery: models.Delivery{
			Name:    fmt.Sprintf("User %s", randomName()),
			Phone:   fmt.Sprintf("+7%010d", rand.Intn(9999999999)),
			Zip:     fmt.Sprintf("101%d", rand.Intn(9000)),
			City:    randomCity(),
			Address: fmt.Sprintf("str. test", rand.Intn(100)+1, rand.Intn(200)+1),
			Region:  "Test region",
			Email:   fmt.Sprintf("test%d@example.com", rand.Intn(10000)),
		},
		Payment: models.Payment{
			Transaction:  orderUID,
			RequestID:    "",
			Currency:     "USD",
			Provider:     "wbpay",
			Amount:       0,
			PaymentDt:    time.Now().Unix(),
			Bank:         "alfa",
			DeliveryCost: 300 + float64(rand.Intn(1000)),
			GoodsTotal:   0,
			CustomFee:    0,
		},
		Items:             items,
		Locale:            "ru",
		InternalSignature: "",
		CustomerID:        fmt.Sprintf("cust-%d", rand.Intn(5000)),
		DeliveryService:   "meest",
		ShardKey:          "9",
		SmID:              rand.Intn(10),
		DateCreated:       time.Now(),
		OofShard:          "1",
	}
	var goodsTotal float64
	for _, item := range items {
		goodsTotal += item.TotalPrice
	}
	order.Payment.GoodsTotal = goodsTotal
	order.Payment.Amount = goodsTotal + order.Payment.DeliveryCost

	return order
}

func randomProduct() string {
	products := []string{"ps4", "keyboard", "glasses", "jacket", "gloves", "bball"}
	return products[rand.Intn(len(products))]
}

func randomBrand() string {
	brands := []string{"Nike", "Adidas", "Puma", "Zara", "H&M", "Uniqlo"}
	return brands[rand.Intn(len(brands))]
}

func randomCity() string {
	cities := []string{"MSC", "SPB", "KZN", "EKB", "NN", "VLG"}
	return cities[rand.Intn(len(cities))]
}

func randomName() string {
	names := []string{"Alice", "Bob", "Charlie", "Diana", "Eve", "Frank"}
	return names[rand.Intn(len(names))]
}

func randomSize() string {
	sizes := []string{"S", "M", "L", "XL", "XXL"}
	return sizes[rand.Intn(len(sizes))]
}
