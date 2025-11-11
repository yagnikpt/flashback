package app

import (
	"database/sql"

	"github.com/yagnikpt/flashback/internal/config"
)

type App struct {
	DB     *sql.DB
	Config config.Config
}

func NewApp(db *sql.DB, config config.Config) *App {
	return &App{
		DB:     db,
		Config: config,
	}
}

// func (app *App) CreateItem(content string) {
// 	// TODO:
// 	// 1 -
// }
