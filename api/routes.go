package api

import (
	"github.com/gorilla/mux"
	"wb_l0/internals/app/handlers"
)

// папка апи отдельно для удобной генерации api (если понадобится)
func CreateRoutes(orderHandler *handlers.OrderHandler) *mux.Router {
	r := mux.NewRouter() //создадим роутер для обработки путей, он же будет основным роутером для нашего сервера
	r.HandleFunc("/orders/{id}", orderHandler.Find).Methods("GET")

	return r
}
