package migrations

import (
	"database/sql"
	"embed"
	"github.com/pressly/goose/v3"
)

//go:embed *.sql
var embedMigrations embed.FS

func Migrate(dsn string) {
	var db *sql.DB
	// setup database
	db, err := goose.OpenDBWithDriver("pgx", dsn)
	if err != nil {
		panic(err)
	}

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}

	if err := goose.Up(db, "."); err != nil {
		panic(err)
	}
}
