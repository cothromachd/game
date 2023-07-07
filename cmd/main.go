package main

import (
	"context"
	authHandler "github.com/cothromachd/game/internal/auth/delivery/http"
	authRepo "github.com/cothromachd/game/internal/auth/repository"
	authService "github.com/cothromachd/game/internal/auth/service"
	"github.com/cothromachd/game/internal/config"
	gameHandler "github.com/cothromachd/game/internal/game/delivery/http"
	gameRepo "github.com/cothromachd/game/internal/game/repository"
	gameService "github.com/cothromachd/game/internal/game/service"
	hashito "github.com/cothromachd/game/pkg/hash"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	log "github.com/sirupsen/logrus"

	"os"
)

var (
	secret = os.Getenv("JWT_SECRET")
	salt   = os.Getenv("JWT_SALT")
)

const (
	configPath = "../configs/config.yaml"
)

func main() {
	runApp()
}

func runApp() {
	app := fiber.New()
	cfg, err := config.New(configPath)
	if err != nil {
		log.Println(err)
	}

	pgPool, err := pgxpool.New(context.Background(), cfg.DB.Conn)
	if err != nil {
		log.Fatal(err)
	}

	hasher := hashito.NewSHA1Hasher(salt)

	authStorage := authRepo.NewStorage(pgPool)
	authServ := authService.NewUser(authStorage, hasher, []byte(secret))
	app = authHandler.NewHandler(app, authServ)

	gameStorage := gameRepo.NewStorage(pgPool)
	gameServ := gameService.NewGame(gameStorage)
	app = gameHandler.NewHandler(app, gameServ)

	log.Infoln("Server started...")
	log.Fatal(app.Listen(cfg.Srv.Addr))
}
