package delivery

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

var secret = os.Getenv("JWT_SECRET")

func (h *Handler) authMiddleware(ctx *fiber.Ctx) error {
	token, err := getTokenFromRequest(ctx)
	if err != nil {
		return err
	}

	parser := jwt.Parser{ValidMethods: []string{jwt.SigningMethodHS256.Alg()}}
	t, err := parser.ParseWithClaims(token, &jwt.StandardClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	claims, ok := t.Claims.(*jwt.StandardClaims)
	if !ok {
		logError("authMiddleware", err)
		return errors.New("invalid claims")
	}

	err = claims.Valid()
	if err != nil {
		logError("authMiddleware", err)
		return errors.New("invalid token")
	}
	if !t.Valid {
		logError("authMiddleware", err)
		return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	exp := claims.ExpiresAt
	if exp == 0 {
		logError("authMiddleware", err)
		return ctx.Status(fiber.StatusBadRequest).SendString("token expiration is not set")
	} else if !(exp >= time.Now().Unix()) {
		logError("authMiddleware", err)
		return ctx.Status(fiber.StatusBadRequest).SendString("login token expired, get new one")
	}

	subject := claims.Subject
	if subject == "" {
		logError("authMiddleware", err)
		return ctx.Status(fiber.StatusBadRequest).SendString("userID is not set")
	}

	id, err := strconv.Atoi(subject)
	if err != nil {
		logError("authMiddleware", err)
		return ctx.Status(fiber.StatusInternalServerError).SendString("can't convert id to integer")
	}
	ctx.Locals("id", id)

	role := claims.Issuer
	if role == "" {
		logError("authMiddleware", err)
		return ctx.Status(fiber.StatusBadRequest).SendString("user role is not set")
	}

	ctx.Locals("role", role)
	return ctx.Next()
}

func getTokenFromRequest(ctx *fiber.Ctx) (string, error) {
	header := ctx.Get("Authorization")
	if header == "" {
		return "", errors.New("empty auth header")
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return "", errors.New("invalid auth header")
	}

	if len(headerParts[1]) == 0 {
		return "", errors.New("token is empty")
	}
	return headerParts[1], nil
}
