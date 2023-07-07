package delivery

import (
	"context"
	"fmt"
	"github.com/cothromachd/game/internal/game/models"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func (h *Handler) GetInfo(ctx *fiber.Ctx) error {
	userIDHeader, userRole := ctx.Get("id"), ctx.Get("role")
	userID, err := strconv.Atoi(userIDHeader)
	if err != nil {
		logError("GetInfo", err)
		return err
	}

	if userRole == models.CustomerRole {
		customerInfo, err := h.gameService.GetCustomerInfo(context.Background(), userID)
		if err != nil {
			logError("GetInfo", err)
			return err
		}

		return ctx.JSON(customerInfo)
	} else if userRole == models.WorkerRole {
		workerInfo, err := h.gameService.GetWorkerInfo(context.Background(), userID)
		if err != nil {
			logError("GetInfo", err)
			return err
		}

		return ctx.JSON(workerInfo)
	}

	return ctx.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("unknown user role: %s", userRole))
}
