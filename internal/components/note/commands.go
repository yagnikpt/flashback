package note

import (
	"log"

	tea "github.com/charmbracelet/bubbletea/v2"
)

type noteAddedMsg bool
type statusMsg string

func addNoteCmd(m Model, note string) tea.Cmd {
	return func() tea.Msg {
		err := m.store.Notes.CreateNote(note)
		if err != nil {
			log.Println("Error creating note:", err)
			return noteAddedMsg(false)
		}
		return noteAddedMsg(true)
	}
}

func readStatusText(ch <-chan string) tea.Cmd {
	return func() tea.Msg {
		status := <-ch
		return statusMsg(status)
	}
}
