package repository

import (
	"context"
	"github.com/cothromachd/game/internal/game/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	pgPool *pgxpool.Pool
}

func NewStorage(pool *pgxpool.Pool) Storage {
	return Storage{pool}
}

func (s Storage) Close() {
	s.pgPool.Close()
}

func (s Storage) GetCustomerInfo(ctx context.Context, customerID int) (models.GetCustomerInfoResponse, error) {
	tx, err := s.pgPool.Begin(ctx)
	if err != nil {
		return models.GetCustomerInfoResponse{}, err
	}

	avaliableAttributes := models.GetCustomerInfoResponse{}
	err = tx.QueryRow(ctx, "SELECT start_capital FROM customers WHERE user_id = $1;", customerID).Scan(&avaliableAttributes.CustomerStartCapital)
	if err != nil {
		tx.Rollback(ctx)
		return models.GetCustomerInfoResponse{}, err
	}

	rows, err := tx.Query(ctx, "SELECT user_id, max_weight, is_drunk, fatigue, salary FROM workers;")
	for rows.Next() {
		worker := models.Worker{}
		err = rows.Scan(&worker.UserID, &worker.MaxWeight, &worker.IsDrunk, &worker.Fatigue, &worker.Salary)
		if err != nil {
			tx.Rollback(ctx)
			return models.GetCustomerInfoResponse{}, err
		}

		avaliableAttributes.Workers = append(avaliableAttributes.Workers, worker)
	}

	if err = tx.Commit(ctx); err != nil {
		return models.GetCustomerInfoResponse{}, err
	}

	return avaliableAttributes, nil
}

func (s Storage) GetAvaliableOrdersForCustomer(ctx context.Context, customerID int) ([]models.GetAvaliableOrdersForCustomerResponse, error) {
	rows, err := s.pgPool.Query(ctx, "SELECT id, name, weight FROM orders WHERE customer_id = $1 AND worker_id IS NULL;", customerID)
	if err != nil {
		return nil, err
	}

	var orders []models.GetAvaliableOrdersForCustomerResponse
	for rows.Next() {
		order := models.GetAvaliableOrdersForCustomerResponse{}
		err = rows.Scan(&order.ID, &order.Name, &order.Weight)
		if err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}

	return orders, nil
}

func (s Storage) GetCustomerAndWorkersStats(ctx context.Context, customerID int, workerIDs []int, orderID int) (int, map[int]models.GetAndCreateWorkerInfo, error) {
	tx, err := s.pgPool.Begin(ctx)
	if err != nil {
		return 0, nil, err
	}

	var customerStartCapital int
	err = tx.QueryRow(ctx, "SELECT start_capital FROM customers WHERE user_id = $1;", customerID).Scan(&customerStartCapital)
	if err != nil {
		tx.Rollback(ctx)
		return 0, nil, err
	}

	workers := make(map[int]models.GetAndCreateWorkerInfo, len(workerIDs))
	for _, workerID := range workerIDs {
		var worker models.GetAndCreateWorkerInfo
		err = tx.QueryRow(ctx, "SELECT max_weight, is_drunk, fatigue, salary FROM workers WHERE user_id = $1;", workerID).
			Scan(&worker.MaxWeight, &worker.IsDrunk, &worker.Fatigue, &worker.Salary)
		if err != nil {
			tx.Rollback(ctx)
			return 0, nil, err
		}

		workers[workerID] = worker
	}

	err = tx.Commit(ctx)
	if err != nil {
		tx.Rollback(ctx)
		return 0, nil, err
	}

	return customerStartCapital, workers, nil
}

func (s Storage) UpdateCustomer(ctx context.Context, customerID int, customerStartCapital int) error {
	_, err := s.pgPool.Exec(ctx, "UPDATE customers SET start_capital = $1 WHERE user_id = $2;", customerStartCapital, customerID)
	if err != nil {
		return err
	}

	return nil
}
