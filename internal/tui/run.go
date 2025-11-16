package tui

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/yagnikpt/flashback/internal/app"
)

func Run(app *app.App) {
	p := tea.NewProgram(NewModel(app), tea.WithAltScreen())
	_, err := p.Run()
	if err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
