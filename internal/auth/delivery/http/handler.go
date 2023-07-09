package delivery

import (
	"context"

	"github.com/cothromachd/game/internal/auth/models"
	"github.com/gofiber/fiber/v2"
)

type UserService interface {
	RegisterUser(ctx context.Context, req models.RegisterUserRequest) error
	LoginUser(ctx context.Context, req models.LoginUserRequest) (string, error)
}

type Handler struct {
	userService UserService
}

func NewHandler(app *fiber.App, userService UserService) *fiber.App {
	h := Handler{
		userService: userService,
	}

	app.Get("/login", h.loginUser)
	app.Post("/register", h.registerUser)

	return app
}
