package migration

import (
	"database/sql"
	"embed"

	"github.com/pressly/goose/v3"
)

//go:embed *.sql
var sqlMigrations embed.FS

func Migrate(db *sql.DB) error {
	goose.SetBaseFS(sqlMigrations)
	goose.SetLogger(goose.NopLogger())

	if err := goose.SetDialect("turso"); err != nil {
		return err
	}

	if err := goose.Up(db, "."); err != nil {
		return err
	}
	return nil
}
