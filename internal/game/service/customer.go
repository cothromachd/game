package service

import (
	"context"
	"github.com/cothromachd/game/internal/game/models"
)

type CustomerRepository interface {
	GetCustomerInfo(ctx context.Context, customerID int) (models.GetCustomerInfoResponse, error)
	GetAvaliableOrdersForCustomer(ctx context.Context, customerID int) ([]models.GetAvaliableOrdersForCustomerResponse, error)
	GetCustomerAndWorkersStats(ctx context.Context, customerID int, workerIDs []int, orderID int) (int, map[int]models.GetAndCreateWorkerInfo, error)
	UpdateCustomer(ctx context.Context, customerID int, customerStartCapital int) error
}

type GameRepository interface {
	CustomerRepository
	WorkerRepository
	OrderRepository
}

type GameService struct {
	gameRepo GameRepository
}

func NewGame(gameRepo GameRepository) GameService {
	return GameService{gameRepo: gameRepo}
}

func (gs GameService) GetCustomerInfo(ctx context.Context, customerID int) (models.GetCustomerInfoResponse, error) {
	return gs.gameRepo.GetCustomerInfo(ctx, customerID)
}

func (gs GameService) GetAvaliableOrdersForCustomer(ctx context.Context, customerID int) ([]models.GetAvaliableOrdersForCustomerResponse, error) {
	return gs.gameRepo.GetAvaliableOrdersForCustomer(ctx, customerID)
}

func (gs GameService) StartGame(ctx context.Context, customerID int, workerIDs []int, orderID int) (bool, error) {
	customerStartCapital, workers, err := gs.gameRepo.GetCustomerAndWorkersStats(ctx, customerID, workerIDs, orderID)
	if err != nil {
		return false, nil
	}

	order, err := gs.gameRepo.GetOrder(ctx, orderID)
	if err != nil {
		return false, err
	}

	var totalWorkersSalary int
	var totalMaxWeight float64
	for _, worker := range workers {
		totalWorkersSalary += worker.Salary

		if worker.IsDrunk {
			totalMaxWeight += float64(worker.MaxWeight) * ((100 - float64(worker.Fatigue)) / 100) * (float64(worker.Fatigue+50) / 100)
		} else {
			totalMaxWeight += float64(worker.MaxWeight) * ((100 - float64(worker.Fatigue)) / 100)
		}
	}

	var success bool
	if totalWorkersSalary > customerStartCapital {
		success = false
	} else if float64(order.Weight) > totalMaxWeight {
		success = false
	} else {
		success = true
	}

	err = gs.gameRepo.UpdateCustomer(ctx, customerID, customerStartCapital-totalWorkersSalary)
	if err != nil {
		return false, err
	}

	for workerID, worker := range workers {
		if worker.Fatigue+20 >= 100 {
			worker.Fatigue = 100
			err := gs.gameRepo.UpdateWorker(ctx, workerID, worker)
			if err != nil {
				return false, err
			}
		} else {
			worker.Fatigue = worker.Fatigue + 20
			err := gs.gameRepo.UpdateWorker(ctx, workerID, worker)
			if err != nil {
				return false, err
			}
		}

		if success {
			err := gs.gameRepo.UpdateOrderAsDone(ctx, orderID, workerID)
			if err != nil {
				return false, err
			}
		}
	}

	return success, nil
}
