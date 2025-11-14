package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/tursodatabase/turso-go"

	"github.com/yagnikpt/flashback/cmd"
	"github.com/yagnikpt/flashback/internal/app"
	"github.com/yagnikpt/flashback/internal/components/apikeyinput"
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

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	logFile := filepath.Join(dataDir, "debug.log")
	fLog, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	log.SetOutput(fLog)
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer fLog.Close()

	db, err := sql.Open("turso", filepath.Join(dataDir, "flashback.db"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = migration.Migrate(db)
	if err != nil {
		log.Fatal(err)
	}

	if cfg.APIKey == "" {
		apikeyinput.Run(configFile, cfg)
		cfg, _ = config.LoadConfig(configFile)
	}

	app := app.NewApp(db, cfg)
	cmd.Execute(app)
}
