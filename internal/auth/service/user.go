package service

import (
	"context"
	"github.com/cothromachd/game/internal/auth/models"
	"github.com/golang-jwt/jwt/v5"
	"strconv"
	"time"
)

type PasswordHasher interface {
	Hash(password string) (string, error)
}

type UserRepository interface {
	CreateUser(ctx context.Context, user models.User) error
	GetByCredentials(ctx context.Context, username, password string) (models.User, error)
	Close()
}

type User struct {
	repo   UserRepository
	hasher PasswordHasher

	hmacSercet []byte
}

func NewUser(repo UserRepository, hasher PasswordHasher, secret []byte) *User {
	return &User{
		repo:       repo,
		hasher:     hasher,
		hmacSercet: secret,
	}
}

func (s *User) RegisterUser(ctx context.Context, req models.RegisterUserRequest) error {
	password, err := s.hasher.Hash(req.Password)
	if err != nil {
		return err
	}

	return s.repo.CreateUser(ctx, models.User{
		Username: req.Username,
		Password: password,
		Role:     req.Role,
	})
}

func (s *User) LoginUser(ctx context.Context, req models.LoginUserRequest) (string, error) {
	password, err := s.hasher.Hash(req.Password)
	if err != nil {
		return "", err
	}

	user, err := s.repo.GetByCredentials(ctx, req.Username, password)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   strconv.Itoa(user.ID),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
		Issuer:    user.Role,
	})

	return token.SignedString(s.hmacSercet)
}

/*
func (s *User) ParseToken(ctx context.Context, token string) (int64, string, error) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return s.hmacSercet, nil
	})
	if err != nil {
		return 0, "", err
	}

	if !t.Valid {
		return 0, "", errors.New("invalid token")
	}

	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return 0, "", errors.New("invalid claims")
	}

	expTime, ok := claims["exp"].(float64)
	if !ok {
		return 0, "", errors.New("invalid claims")
	}
	if float64(time.Now().Unix()) > expTime {
		return 0, "", errors.New("token is expired")
	}

	subject, ok := claims["sub"].(string)
	if !ok {
		return 0, "", errors.New("invalid subject")
	}

	id, err := strconv.Atoi(subject)
	if err != nil {
		return 0, "", errors.New("invalid subject")
	}

	role, ok := claims["iss"].(string)
	if !ok {
		return 0, "", errors.New("invalid subject")
	}

	return int64(id), role, nil
}
*/
