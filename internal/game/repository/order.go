package repository

import (
	"context"
	"github.com/cothromachd/game/internal/game/models"
	"math/rand"
	"time"
)

func (s Storage) GetOrder(ctx context.Context, orderID int) (models.GetOrderResponse, error) {
	var order models.GetOrderResponse
	err := s.pgPool.QueryRow(ctx, "SELECT name, weight, customer_id FROM orders WHERE id = $1;", orderID).
		Scan(&order.Name, &order.Weight, &order.CustomerID)
	if err != nil {
		return models.GetOrderResponse{}, err
	}

	return order, nil
}

func (s Storage) UpdateOrderAsDone(ctx context.Context, orderID, workerID int) error {
	_, err := s.pgPool.Exec(ctx, "UPDATE orders SET worker_id = $1 WHERE id = $2", workerID, orderID)

	return err
}

func (s Storage) CreateRandomOrders(ctx context.Context, customerUsername string) error {
	tx, err := s.pgPool.Begin(ctx)
	if err != nil {
		return err
	}

	var customerID int
	err = tx.QueryRow(ctx, "SELECT id FROM users WHERE username = $1;", customerUsername).Scan(&customerID)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	numOfOrders := r.Intn(8-2+1) + 2
	for i := 0; i < numOfOrders; i++ {
		_, err := tx.Exec(ctx, "INSERT INTO orders (name, weight, customer_id) VALUES ($1, $2, $3);",
			models.GenerateRandomOrderName(), r.Intn(80-10+1)+10, customerID)
		if err != nil {
			tx.Rollback(ctx)
			return err
		}
	}

	return tx.Commit(ctx)
}
