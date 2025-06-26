package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea/v2"
	app "github.com/yagnik-patel-47/flashback/internal/app"
)

func main() {
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()

	p := tea.NewProgram(app.InitModel(), tea.WithAltScreen(), tea.WithKeyboardEnhancements())
	res, err := p.Run()
	if err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}

	final, _ := res.(app.Model)
	fmt.Print(final.Output)
}
