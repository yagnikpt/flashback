package tui

import (
	"log"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/yagnikpt/flashback/internal/config"
	"github.com/yagnikpt/flashback/internal/utils"
)

func saveConfigCmd(m Model) tea.Cmd {
	return func() tea.Msg {
		configDir, _ := utils.GetConfigDir()
		configFile := filepath.Join(configDir, "config.toml")
		err := config.SaveConfig(configFile, m.store.Config)
		if err != nil {
			log.Println("Error saving config:", err)
			// return saveConfigMsg(false)
		}
		return nil
	}
}
