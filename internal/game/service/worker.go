package service

import (
	"context"
	"github.com/cothromachd/game/internal/game/models"
)

type WorkerRepository interface {
	GetWorkerInfo(ctx context.Context, workerID int) (models.GetAndCreateWorkerInfo, error)
	GetCompletedWorkerOrders(ctx context.Context, workerID int) ([]models.GetCompletedWorkerOrdersResponse, error)
	UpdateWorker(ctx context.Context, workerID int, worker models.GetAndCreateWorkerInfo) error
}

func (gs GameService) GetWorkerInfo(ctx context.Context, workerID int) (models.GetAndCreateWorkerInfo, error) {
	return gs.gameRepo.GetWorkerInfo(ctx, workerID)
}
func (gs GameService) GetCompletedWorkerOrders(ctx context.Context, workerID int) ([]models.GetCompletedWorkerOrdersResponse, error) {
	return gs.gameRepo.GetCompletedWorkerOrders(ctx, workerID)
}
