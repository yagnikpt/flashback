package spinner

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
)

func Run(status <-chan string, altScreen bool) {
	model := NewModel(status)
	model.SetAltScreen(altScreen)
	p := tea.NewProgram(model)
	res, err := p.Run()
	if err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
	model = res.(Model)
	model.SetDisplayText("")
}
