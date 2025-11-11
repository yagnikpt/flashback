package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/yagnikpt/flashback/cmd"
	"github.com/yagnikpt/flashback/internal/app"
	"github.com/yagnikpt/flashback/internal/config"
	"github.com/yagnikpt/flashback/internal/migration"
	"github.com/yagnikpt/flashback/internal/utils"
)

func main() {
	dataDir, err := utils.GetLocalDataDir()
	if err != nil {
		fmt.Println("Error getting local data directory:", err)
		os.Exit(1)
	}

	configDir, err := utils.GetConfigDir()
	if err != nil {
		fmt.Println("Error getting local data directory:", err)
		os.Exit(1)
	}
	configFile := filepath.Join(configDir, "config.toml")
	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		fmt.Println("Error loading config:", err)
		os.Exit(1)
	}

	f, err := tea.LogToFile(filepath.Join(dataDir, "debug.log"), "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()

	db, err := sql.Open("turso", filepath.Join(dataDir, "flashback.db"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = migration.Migrate(db)
	if err != nil {
		log.Fatal(err)
	}

	app := &app.App{
		DB:     db,
		Config: cfg,
	}
	cmd.Execute(app)
}
