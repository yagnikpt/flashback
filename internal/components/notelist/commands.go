package notelist

import (
	"log"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/yagnikpt/flashback/internal/notes"
)

type notesMsg []notes.CombinedNote
type deleteNoteMsg bool

func deleteNoteCmd(m Model, noteID int) tea.Cmd {
	return func() tea.Msg {
		err := m.store.Notes.DeleteNote(noteID)
		if err != nil {
			log.Println("Error creating note:", err)
			return deleteNoteMsg(false)
		}
		return deleteNoteMsg(true)
	}
}

func getNotesCmd(m Model) tea.Cmd {
	return func() tea.Msg {
		notes, err := m.store.Notes.GetAllNotes()
		if err != nil {
			log.Println("Error getting notes:", err)
			return notesMsg(nil)
		}
		// log.Println(notes)
		return notesMsg(notes)
	}
}
