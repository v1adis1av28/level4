package service

import (
	"demo/internal/cache"
	"demo/internal/models"
	"demo/internal/repository"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type OrderService struct {
	OrderRepository *repository.OrderRepository
	Redis           *redis.Client
}

// зщдесь делать проверку на наличии в кэше записи
func (o *OrderService) GetOrderByUUID(uuid string) (*models.Order, error) {
	//Проверка что в кэше нет инстанса по запросу
	// если нил то делаем вызов в репозиторий и оттуда достаем инстанс ордера и делаем сет в кэш
	//Если не нил то просто возвращаем значение из checkCache
	cacheVal := cache.GetOrderFromCache(uuid, o.Redis)
	if cacheVal == nil {
		order, err := o.OrderRepository.GetOrderByUUID(uuid)
		if err != nil {
			return nil, fmt.Errorf("error on getting order by uuid :%s", uuid)
		}
		cache.SetCache(uuid, order, o.Redis)
		return order, nil
	} else {
		return cacheVal, nil
	}
}

func (o *OrderService) NewOrder(order *models.Order) error {
	return o.OrderRepository.NewOrder(order)
}

func NewOrderService(or *repository.OrderRepository, r *redis.Client) *OrderService {
	return &OrderService{OrderRepository: or, Redis: r}
}
