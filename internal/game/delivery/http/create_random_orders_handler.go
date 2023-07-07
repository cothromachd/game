package delivery

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) CreateRandomOrders(ctx *fiber.Ctx) error {
	reqBody := ctx.Body()
	m := make(map[string]string)
	err := json.Unmarshal(reqBody, &m)
	if err != nil {
		logError("CreateRandomOrders", err)
		return err
	}

	customerUsername, ok := m["customerUsername"]
	if !ok {
		err := errors.New("empty input")
		logError("CreateRandomOrders", err)
		return err
	}

	err = h.gameService.CreateRandomOrders(context.Background(), customerUsername)
	if err != nil {
		logError("CreateRandomOrders", err)
		return err
	}

	return ctx.SendStatus(fiber.StatusOK)
}
