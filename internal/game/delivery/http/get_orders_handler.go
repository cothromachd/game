package delivery

import (
	"context"
	"fmt"
	"github.com/cothromachd/game/internal/game/models"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func (h *Handler) GetOrders(ctx *fiber.Ctx) error {
	userIDHeader, userRole := ctx.Get("id"), ctx.Get("role")
	userID, err := strconv.Atoi(userIDHeader)
	if err != nil {
		logError("GetOrders", err)
		return err
	}

	if userRole == models.CustomerRole {
		customerOrders, err := h.gameService.GetAvaliableOrdersForCustomer(context.Background(), userID)
		if err != nil {
			logError("GetOrders", err)
			return err
		}

		return ctx.JSON(customerOrders)
	} else if userRole == models.WorkerRole {
		workerOrders, err := h.gameService.GetCompletedWorkerOrders(context.Background(), userID)
		if err != nil {
			logError("GetOrders", err)
			return err
		}

		return ctx.JSON(workerOrders)
	}

	return ctx.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("unknown user role: %s", userRole))
}
