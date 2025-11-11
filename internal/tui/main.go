package tui

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/tursodatabase/turso-go"

	tea "github.com/charmbracelet/bubbletea/v2"

	"github.com/yagnikpt/flashback/internal/components/apiprompt"
	"github.com/yagnikpt/flashback/internal/config"
	"github.com/yagnikpt/flashback/internal/global"
	"github.com/yagnikpt/flashback/internal/utils"
)

func Run(db *sql.DB, cfg config.Config) {
	configDir, err := utils.GetConfigDir()
	if err != nil {
		fmt.Println("Error getting local data directory:", err)
		os.Exit(1)
	}
	configFile := filepath.Join(configDir, "config.toml")

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
		p := tea.NewProgram(InitModel(), tea.WithAltScreen(), tea.WithKeyboardEnhancements())
		_, err := p.Run()
		if err != nil {
			fmt.Printf("Alas, there's been an error: %v", err)
			os.Exit(1)
		}
	}
}
