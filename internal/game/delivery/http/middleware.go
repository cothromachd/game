package delivery

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"strconv"
	"strings"
	"time"
)

func (h *Handler) authMiddleware(ctx *fiber.Ctx) error {
	token, err := getTokenFromRequest(ctx)
	if err != nil {
		return err
	}

	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return os.Getenv("JWT_SECRET"), nil
	})
	if err != nil {
		return err
	}

	if !t.Valid {
		return errors.New("invalid token")
	}

	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return errors.New("invalid claims")
	}

	expTime, ok := claims["exp"].(float64)
	if !ok {
		return errors.New("invalid claims")
	}
	if float64(time.Now().Unix()) > expTime {
		return errors.New("token is expired")
	}

	subject, ok := claims["sub"].(string)
	if !ok {
		return errors.New("invalid subject")
	}

	id, err := strconv.Atoi(subject)
	if err != nil {
		return errors.New("invalid subject")
	}
	ctx.Set("id", strconv.Itoa(id))

	role, ok := claims["iss"].(string)
	if !ok {
		return errors.New("invalid subject")
	}
	ctx.Set("role", role)

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
