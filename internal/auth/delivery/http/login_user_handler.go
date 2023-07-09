package delivery

import (
	"context"
	"encoding/json"

	"github.com/cothromachd/game/internal/auth/models"
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) loginUser(ctx *fiber.Ctx) error {
	reqBody := ctx.Body()
	var req models.LoginUserRequest
	if err := json.Unmarshal(reqBody, &req); err != nil {
		logError("loginUser", err)
		return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	if err := req.Validate(); err != nil {
		logError("loginUser", err)
		return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	token, err := h.userService.LoginUser(context.Background(), req)
	if err != nil {
		logError("loginUser", err)
		return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return ctx.JSON(map[string]string{
		"token": token,
	})
}
