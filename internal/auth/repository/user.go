package repository

import (
	"context"
	"github.com/cothromachd/game/internal/auth/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"math"
	"math/rand"
	"time"
)

type Storage struct {
	pgPool *pgxpool.Pool
}

func NewStorage(pool *pgxpool.Pool) Storage {
	return Storage{pool}
}

func (s Storage) Close() {
	s.pgPool.Close()
}

func (s Storage) CreateUser(ctx context.Context, user models.User) error {
	tx, err := s.pgPool.Begin(ctx)
	if err != nil {
		return err
	}

	var userID int
	err = tx.QueryRow(ctx, "INSERT INTO users (username, password, role) VALUES ('$1', '$2', '$3') RETURNING id;",
		user.Username, user.Password, user.Role).Scan(&userID)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	if user.Role == models.CustomerRole {
		_, err = tx.Exec(ctx, "INSERT INTO customers (user_id, start_capital) VALUES ($1, $2);", userID,
			r.Intn(100000-10000+1)+10000)
		if err != nil {
			tx.Rollback(ctx)
			return err
		}
	} else if user.Role == models.WorkerRole {
		_, err = tx.Exec(ctx, "INSERT INTO workers (user_id, max_weight, is_drunk, fatigue, salary) VALUES ($1, $2, $3, $4, $5);",
			userID, r.Intn(10-5+1)+5, r.Intn(2) == 1, math.Round((r.Float64())*100+1), r.Intn(30000-10000+1)+10000)
		if err != nil {
			tx.Rollback(ctx)
			return err
		}
	}

	return tx.Commit(ctx)
}

func (s Storage) GetByCredentials(ctx context.Context, username, password string) (models.User, error) {
	var user models.User
	err := s.pgPool.QueryRow(ctx, "SELECT id, username, password, role FROM users WHERE username = $1 AND password = $2;",
		username, password).Scan(&user.ID, &user.Username, &user.Password, &user.Role)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}
