package apikeyinput

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/yagnikpt/flashback/internal/config"
)

func Run(configFile string, cfg config.Config) {
	p := tea.NewProgram(NewModel())
	res, err := p.Run()
	if err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
	model := res.(Model)
	apiKey := model.Output
	cfg.APIKey = apiKey
	err = config.SaveConfig(configFile, cfg)
	if err != nil {
		fmt.Println("Error saving config:", err)
		os.Exit(1)
	}
}
