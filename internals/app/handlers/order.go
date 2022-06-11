package handlers

import (
	_ "context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nats-io/stan.go"
	"log"
	"net/http"
	"wb_l0/internals/app/models"
	"wb_l0/internals/app/processors"
)

type OrderHandler struct {
	processor *processors.OrdersProcessor
}

//У хэндлера одна задача, принять запрос, декодировать его, и передать их в процессор,
//при возврате даннных из процессора их обернуть в нужный формат и венуть их обратно

func NewOrdersHandler(processor *processors.OrdersProcessor) *OrderHandler {
	handler := new(OrderHandler)
	handler.processor = processor // процессор это наш сервис работающий с данными
	return handler
}

type natsMessage struct {
	OrderNumber string `json:"order_uid"`
}

func (handler *OrderHandler) Create(m *stan.Msg) {
	var msg natsMessage

	// Check if message has field "order_uid"
	if err := json.Unmarshal(m.Data, &msg); err != nil {
		log.Println(err)
		return
	}

	// пишем в модель данные от натс
	order := &models.Order{}
	order.Data = m.Data
	order.OrderNumber = msg.OrderNumber
	// пишем в постгрес и айди в мапу
	if err := handler.processor.CreateOrder(order); err != nil {
		log.Println(err)
	}
}

// отображения полученных данных , для запроса по id
func (handler *OrderHandler) Find(w http.ResponseWriter, r *http.Request) {

	orderNumber := mux.Vars(r) //переменные, объявленные в ресурсах попадают в Vars и могут быть адресованы
	if orderNumber["id"] == "" {
		fmt.Fprint(w, http.StatusText(http.StatusBadGateway))
		return
	}
	id := orderNumber["id"]

	order, err := handler.processor.FindOrderId(id)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "order not found")
		return
	}

	fmt.Fprint(w, string(order))
}
