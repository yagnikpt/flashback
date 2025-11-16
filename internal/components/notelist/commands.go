package notelist

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/yagnikpt/flashback/internal/models"
)

type getAllNotesMsg []models.FlashbackWithMetadata
type deleteNoteMsg bool
type chosenNoteMsg models.FlashbackWithMetadata
type relayChooseMsg string
type relayDeleteMsg string

func getAllNotesCmd(m Model) tea.Cmd {
	return func() tea.Msg {
		notes, err := m.app.GetAllNotes()
		if err != nil {
			log.Fatal(err)
		}
		return getAllNotesMsg(notes)
	}
}

func chooseNoteCmd(m Model, noteID string) tea.Cmd {
	return func() tea.Msg {
		note, err := m.app.GetNoteByID(noteID)
		if err != nil {
			log.Println("Error finding note:", err)
			return chosenNoteMsg(note)
		}
		return chosenNoteMsg(note)
	}
}

func deleteNoteCmd(m Model, noteID string) tea.Cmd {
	return func() tea.Msg {
		err := m.app.DeleteNoteByID(noteID)
		if err != nil {
			log.Println("Error deleting note:", err)
			return deleteNoteMsg(false)
		}
		return deleteNoteMsg(true)
	}
}

func relayChooseCmd(id string) tea.Cmd {
	return func() tea.Msg {
		return relayChooseMsg(id)
	}
}

func relayDeleteCmd(id string) tea.Cmd {
	return func() tea.Msg {
		return relayDeleteMsg(id)
	}
}
