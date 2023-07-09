package delivery

import "github.com/gofiber/fiber/v2"

func (h *Handler) healthCheck(ctx *fiber.Ctx) error {
	return ctx.SendStatus(fiber.StatusOK)
}
