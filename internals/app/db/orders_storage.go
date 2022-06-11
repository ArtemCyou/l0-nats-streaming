package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"sync"
	"wb_l0/internals/app/models"
)

type OrdersStorage struct {
	databasePool *pgxpool.Pool
	m  sync.Map
}

func NewOrdersStorage (pool *pgxpool.Pool) *OrdersStorage  {
	storage := new(OrdersStorage)
	storage.databasePool = pool
	return storage
}

// записал в бд и мапку
func (storage *OrdersStorage) CreateOrder(order *models.Order)  error {
	query := "INSERT INTO orders (order_number, order_data) VALUES ($1, $2) RETURNING id"
	if err := storage.databasePool.QueryRow(context.Background(), query, order.OrderNumber, order.Data).Scan(&order.Id); err != nil {
		return fmt.Errorf("error adding order: %w", err)
	}

	storage.m.Store(order.OrderNumber, order.Data)

	return nil
}

//обновляем кеш из бд
func (storage *OrdersStorage) UploadCache(ctx context.Context) error {
	rows, err := storage.databasePool.Query(ctx, `SELECT * FROM orders ORDER BY id`)
	if err != nil {
		return fmt.Errorf("ошибка получения \"данных\" из бд: %w", err)
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		err := rows.Scan(&order.Id, &order.OrderNumber, &order.Data)
		if err != nil {
			log.Fatal(err)
		}
		orders = append(orders, order)
		storage.m.Store(order.OrderNumber, order.Data)
	}

	return nil
}

//ищем данные по айди
func (storage *OrdersStorage) Order(orderNumber string) ([]byte, error) {
	val, ok := storage.m.Load(orderNumber)
	if ok == false {
		log.Println("данный order отсутствует в кеше")
		return nil, errors.New("order не найден")
	}
	value := val.([]byte)
	return value, nil
}