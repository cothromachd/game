package delivery

import (
	"context"
	"encoding/json"

	"github.com/cothromachd/game/internal/auth/models"
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) registerUser(ctx *fiber.Ctx) error {
	reqBody := ctx.Body()

	var req models.RegisterUserRequest
	if err := json.Unmarshal(reqBody, &req); err != nil {
		logError("registerUser", err)
		return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	if err := req.Validate(); err != nil {
		logError("registerUser", err)
		return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	if req.Role != models.CustomerRole && req.Role != models.WorkerRole {
		return ctx.Status(fiber.StatusBadRequest).SendString("invalid role, choose between customer and worker")
	}

	err := h.userService.RegisterUser(context.Background(), req)
	if err != nil {
		logError("registerUser", err)
		return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return ctx.SendStatus(fiber.StatusOK)
}
