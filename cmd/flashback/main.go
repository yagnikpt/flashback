package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/tursodatabase/turso-go"

	tea "github.com/charmbracelet/bubbletea/v2"

	"github.com/yagnikpt/flashback/internal/app"
	"github.com/yagnikpt/flashback/internal/components/apiprompt"
	"github.com/yagnikpt/flashback/internal/config"
	"github.com/yagnikpt/flashback/internal/global"
	"github.com/yagnikpt/flashback/internal/migration"
	"github.com/yagnikpt/flashback/internal/notifications"
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

	if len(os.Args) > 1 {
		command := os.Args[1]
		switch command {
		case "-h", "--help", "help":
			helpText := `Flashback - A flashcard application powered by spaced repetition and AI.` + "\n\n" + `Usage:
  flashback [command]

Available Commands:
  help            Show this help message
  notifications   Start the notification service

If no command is provided, the main application will start.`
			fmt.Println(helpText)
		case "notifications":
			notifications.StartNotificationService(db)
		}
	} else {
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
		} else {
			global.InitStore(db, cfg)
			p := tea.NewProgram(app.InitModel(), tea.WithAltScreen(), tea.WithKeyboardEnhancements())
			_, err = p.Run()
			if err != nil {
				fmt.Printf("Alas, there's been an error: %v", err)
				os.Exit(1)
			}
		}
	}
}
