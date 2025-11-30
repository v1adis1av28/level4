package utils

import "demo/internal/models"

var min_track_number_length int = 12
var min_orderUUID_length int = 12
var min_entry_length int = 5
var min_transaction_length = 12

func ValidateOrder(ord *models.Order) bool {
	// if len(ord.OrderUID) < min_orderUUID_length || len(ord.TrackNumber) < min_track_number_length || len(ord.Entry) < min_entry_length || len(ord.Payment.Transaction) < min_transaction_length {
	// 	return false
	// }
	if ord.Payment.Amount < 0 || ord.Payment.DeliveryCost < 0 || ord.Payment.GoodsTotal < 0 || ord.Payment.CustomFee < 0 {
		return false
	}
	//todo добавить валидацию на остальную часть модели
	return true
}
