package recall

import (
	"log"

	tea "github.com/charmbracelet/bubbletea/v2"
)

type recallMsg string

func recallCmd(m Model, query string) tea.Cmd {
	return func() tea.Msg {
		if m.store == nil {
			log.Println("Error: store is nil")
			return recallMsg("")
		}
		if m.store.Notes == nil {
			log.Println("Error: store.Notes is nil")
			return recallMsg("")
		}
		content, err := m.store.Notes.Recall(query)
		if err != nil {
			log.Println("Error recalling note:", err)
			return recallMsg("")
		}
		return recallMsg(content)
	}
}
