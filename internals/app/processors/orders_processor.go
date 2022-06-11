package processors

import (
	"errors"
	"wb_l0/internals/app/db"
	"wb_l0/internals/app/models"
)

type OrdersProcessor struct {
	storage *db.OrdersStorage
}

func NewUserProcessor(storage *db.OrdersStorage) *OrdersProcessor {
	processor := new(OrdersProcessor)
	processor.storage = storage
	return processor
}

func (processor *OrdersProcessor) CreateOrder(order *models.Order) error {
	if order.OrderNumber == "" {
		return errors.New("order number not be empty")
	}
	return processor.storage.CreateOrder(order)
}

func (processor *OrdersProcessor) FindOrderId(id string) ([]byte, error) {
	order, err := processor.storage.Order(id)
	return order, err
}
