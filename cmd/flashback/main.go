package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/tursodatabase/go-libsql"

	tea "github.com/charmbracelet/bubbletea/v2"

	"github.com/yagnikpt/flashback/internal/app"
	"github.com/yagnikpt/flashback/internal/components/apiprompt"
	"github.com/yagnikpt/flashback/internal/config"
	"github.com/yagnikpt/flashback/internal/global"
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

	db, err := sql.Open("libsql", "file:"+filepath.Join(dataDir, "flashback.db"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, _ = db.Exec("PRAGMA journal_mode=WAL;")
	_, _ = db.Exec("PRAGMA busy_timeout=5000;")

	err = migration.Migrate(db)
	if err != nil {
		log.Fatal(err)
	}

	if cfg.APIKey == "" {
		p := tea.NewProgram(apiprompt.NewModel(), tea.WithAltScreen())
		res, err := p.Run()
		if err != nil {
			fmt.Printf("Alas, there's been an error: %v", err)
			os.Exit(1)
		}
		model := res.(apiprompt.Model)
		apiKey := model.Output
		cfg.APIKey = apiKey
		err = config.SaveConfig(configFile, cfg)
		if err != nil {
			fmt.Println("Error saving config:", err)
			os.Exit(1)
		}
	}

	if len(cfg.APIKey) != 0 {
		global.InitStore(db, cfg)
		p := tea.NewProgram(app.InitModel(), tea.WithAltScreen(), tea.WithKeyboardEnhancements())
		_, err = p.Run()
		if err != nil {
			fmt.Printf("Alas, there's been an error: %v", err)
			os.Exit(1)
		}
	}
}
