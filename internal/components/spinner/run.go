package spinner

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func Run(status <-chan string, altScreen bool) {
	var options []tea.ProgramOption
	if altScreen {
		options = append(options, tea.WithAltScreen())
	}
	model := NewModel(status)
	p := tea.NewProgram(model, options...)
	res, err := p.Run()
	if err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
	model = res.(Model)
	model.SetDisplayText("")
}
