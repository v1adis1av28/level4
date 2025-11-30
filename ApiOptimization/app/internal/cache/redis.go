package cache

import (
	"context"
	"demo/internal/models"
	"encoding/json"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient(addr, password string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})
}

func GetOrderFromCache(uuid string, rc *redis.Client) *models.Order {
	cacheData, err := rc.Get(context.Background(), uuid).Result()
	//Если ответ равен нил то значитв кэше ничего нет надо сделать запрос из бд
	if err == redis.Nil {
		return nil
	}
	var order models.Order
	err = json.Unmarshal([]byte(cacheData), &order)
	if err != nil {
		log.Printf("Error on unmarshaling cache data to dto with uuid: %s", uuid)
		return nil
	}
	return &order
}

func SetCache(uuid string, order *models.Order, rc *redis.Client) error {
	cacheExp := 5 * time.Minute
	data, err := json.Marshal(order)
	if err != nil {
		log.Printf("Error marshaling order with uuid %s: %v", uuid, err)
		return err
	}
	err = rc.Set(context.Background(), uuid, data, cacheExp).Err()
	if err != nil {
		panic(err)
	}

	log.Printf("New cache item add with uuid: %s", uuid)
	return nil
}
