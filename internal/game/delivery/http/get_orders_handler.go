package delivery

import (
	"context"
	"fmt"
	"github.com/cothromachd/game/internal/game/models"
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) GetOrders(ctx *fiber.Ctx) error {
	userID, userRole := ctx.Locals("id").(int), ctx.Locals("role").(string)
	if userRole == models.CustomerRole {
		customerOrders, err := h.gameService.GetAvaliableOrdersForCustomer(context.Background(), userID)
		if err != nil {
			logError("GetOrders", err)
			return err
		} else if len(customerOrders) == 0 {
			return ctx.Status(fiber.StatusOK).SendString("No assigned orders")
		}

		return ctx.JSON(customerOrders)
	} else if userRole == models.WorkerRole {
		workerOrders, err := h.gameService.GetCompletedWorkerOrders(context.Background(), userID)
		if err != nil {
			logError("GetOrders", err)
			return err
		} else if len(workerOrders) == 0 {
			return ctx.Status(fiber.StatusOK).SendString("No orders you already done")
		}

		return ctx.JSON(workerOrders)
	}

	return ctx.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("unknown user role: %s", userRole))
}
