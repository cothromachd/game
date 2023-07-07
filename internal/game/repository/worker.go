package repository

import (
	"context"
	"github.com/cothromachd/game/internal/game/models"
)

func (s Storage) GetWorkerInfo(ctx context.Context, workerID int) (models.GetAndCreateWorkerInfo, error) {
	worker := models.GetAndCreateWorkerInfo{}
	err := s.pgPool.QueryRow(ctx, "SELECT max_weight, is_drunk, fatigue, salary FROM workers WHERE user_id = $1;", workerID).Scan(&worker.MaxWeight, &worker.IsDrunk, &worker.Fatigue, &worker.Salary)
	if err != nil {
		return models.GetAndCreateWorkerInfo{}, err
	}

	return worker, nil
}

func (s Storage) GetCompletedWorkerOrders(ctx context.Context, workerID int) ([]models.GetCompletedWorkerOrdersResponse, error) {
	rows, err := s.pgPool.Query(ctx, "SELECT id, name, weight, worker_id FROM orders WHERE worker_id = $1;", workerID)
	if err != nil {
		return nil, err
	}

	var orders []models.GetCompletedWorkerOrdersResponse
	for rows.Next() {
		order := models.GetCompletedWorkerOrdersResponse{}
		err := rows.Scan(&order.ID, &order.Name, &order.Weight, &order.WorkerID)
		if err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}

	return orders, nil
}

func (s Storage) UpdateWorker(ctx context.Context, workerID int, worker models.GetAndCreateWorkerInfo) error {
	_, err := s.pgPool.Exec(ctx, "UPDATE workers SET fatigue = $1 WHERE user_id = $2;", worker.Fatigue, workerID)

	return err
}
