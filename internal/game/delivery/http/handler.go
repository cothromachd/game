package delivery

import (
	"context"

	"github.com/cothromachd/game/internal/game/models"

	"github.com/gofiber/fiber/v2"
)

type CustomerService interface {
	GetCustomerInfo(ctx context.Context, customerID int) (models.GetCustomerInfoResponse, error)
	GetAvaliableOrdersForCustomer(ctx context.Context, customerID int) ([]models.GetAvaliableOrdersForCustomerResponse, error)
	StartGame(ctx context.Context, customerID int, workerIDs []int, orderID int) (bool, error)
}

type WorkerService interface {
	GetWorkerInfo(ctx context.Context, workerID int) (models.GetAndCreateWorkerInfo, error)
	GetCompletedWorkerOrders(ctx context.Context, workerID int) ([]models.GetCompletedWorkerOrdersResponse, error)
}

type OrderService interface {
	CreateRandomOrders(ctx context.Context, customerUsername string) error
}
type GameService interface {
	CustomerService
	WorkerService
	OrderService
}

type Handler struct {
	salt string

	gameService GameService
}

func NewHandler(app *fiber.App, gameService GameService) *fiber.App {
	h := Handler{
		gameService: gameService,
	}

	app.Get("/me", h.authMiddleware, h.GetInfo)
	app.Get("/tasks", h.authMiddleware, h.GetOrders)
	app.Post("/public/tasks", h.CreateRandomOrders)
	app.Post("/start", h.authMiddleware, h.StartGame)

	return app
}
