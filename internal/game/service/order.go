package service

import (
	"context"
	"github.com/cothromachd/game/internal/game/models"
)

type OrderRepository interface {
	CreateRandomOrders(ctx context.Context, customerUsername string) error
	GetOrder(ctx context.Context, orderID int) (models.GetOrderResponse, error)
	UpdateOrderAsDone(ctx context.Context, workerID, orderID int) error
}

func (gs GameService) CreateRandomOrders(ctx context.Context, customerUsername string) error {
	return gs.gameRepo.CreateRandomOrders(ctx, customerUsername)
}
