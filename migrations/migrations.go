package migrations

import (
	"database/sql"
	"embed"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/pressly/goose/v3"
)

//go:embed *.sql
var embedMigrations embed.FS

func Migrate(dsn string) error {
	var db *sql.DB

	// setup database
	db, err := goose.OpenDBWithDriver("pgx", dsn)
	if err != nil {
		return err
	}

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	if err := goose.Up(db, "."); err != nil {
		return err
	}

	return nil
}
