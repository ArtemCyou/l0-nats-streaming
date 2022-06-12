package processors

import (
	"encoding/json"
	//"errors"
	"github.com/nats-io/stan.go"
	"log"
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

type natsMessage struct {
	OrderNumber string `json:"order_uid"`
}

func (processor *OrdersProcessor) CreateOrder(m *stan.Msg)  {
	//if order.OrderNumber == "" {
	//	return errors.New("order number not be empty")
	//}

	var msg natsMessage
	if err := json.Unmarshal(m.Data, &msg); err != nil {
		log.Println(err)
		return
	}

	// пишем в модель данные от натс
	order := &models.Order{}
	order.Data = m.Data
	order.OrderNumber = msg.OrderNumber
	// пишем в постгрес и айди в мапу

	err := processor.storage.CreateOrder(order)
	if err != nil {
		log.Println(err)
	}
}

func (processor *OrdersProcessor) FindOrderId(id string) ([]byte, error) {
	order, err := processor.storage.Order(id)
	return order, err
}
