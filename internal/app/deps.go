package app

import (
	"context"
	"database/sql"

	"github.com/yagnikpt/flashback/internal/config"
	"google.golang.org/genai"
)

type App struct {
	DB     *sql.DB
	Gemini *genai.Client
	Config config.Config
}

func NewApp(db *sql.DB, config config.Config) *App {
	client, err := genai.NewClient(context.Background(), &genai.ClientConfig{
		APIKey: config.APIKey,
	})
	if err != nil {
		panic(err)
	}

	return &App{
		DB:     db,
		Gemini: client,
		Config: config,
	}
}
