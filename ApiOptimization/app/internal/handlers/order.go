package handlers

import (
	"demo/internal/models"
	"demo/internal/service"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderService *service.OrderService
}

func NewOrderHandler(or *service.OrderService) *OrderHandler {
	return &OrderHandler{orderService: or}
}

func (oh *OrderHandler) HandleIncomingOrder(c *gin.Context) {
	var order models.Order
	decoder := json.NewDecoder(c.Request.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&order)
	//err := c.ShouldBindJSON(&order)
	if err != nil {
		c.JSON(400, gin.H{"error": "bad json"})
		return
	}
	// err = oh.orderService.NewOrder(&order)
	// if err != nil {
	// 	log.Printf("Error on creating new order %s", err.Error())
	// 	return
	// }
	//log.Printf("Order was succesfully created, order_id: %s", order.OrderUID)

}

func (oh *OrderHandler) GetOrderById(c *gin.Context) {
	uuid := c.Param("id")
	if len(uuid) < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	order := &models.Order{}

	order.CustomerID = "cached_customer_id"
	order.OrderUID = uuid
	order.Delivery = models.Delivery{
		Address: "cached_address",
		Name:    "cached_name",
		Phone:   "cached_phone",
		Email:   "cached_email",
	}
	order.Payment = models.Payment{
		Transaction:  "cached_transaction",
		RequestID:    "cached_request_id",
		Currency:     "cached_currency",
		Provider:     "cached_provider",
		Amount:       100,
		PaymentDt:    1631022245,
		Bank:         "cached_bank",
		DeliveryCost: 0,
		GoodsTotal:   100,
		CustomFee:    0,
	}
	order.Items = []models.Item{
		{
			ChrtID:      9934930,
			TrackNumber: "cached_track_number",
			Price:       100,
			RID:         "cached_rid",
			Name:        "cached_name",
			Sale:        0,
			Size:        "cached_size",
			TotalPrice:  100,
			NmID:        2389212,
			Brand:       "cached_brand",
			Status:      202,
		},
	}

	orderDTO := models.OrderToDTO(order)
	c.Writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(c.Writer).Encode(orderDTO)
	//c.JSON(http.StatusOK, orderDTO) //убираем двойное копирование
}
