package recall

import (
	"log"

	tea "github.com/charmbracelet/bubbletea/v2"
)

type recallMsg string

func recallCmd(m Model, query string) tea.Cmd {
	return func() tea.Msg {
		content, err := m.store.Notes.Recall(query)
		if err != nil {
			log.Println("Error creating note:", err)
			return recallMsg("")
		}
		return recallMsg(content)
	}
}
