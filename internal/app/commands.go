package app

import (
	"log"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/yagnik-patel-47/flashback/internal/config"
	"github.com/yagnik-patel-47/flashback/internal/notes"
	"github.com/yagnik-patel-47/flashback/internal/utils"
)

type inputMsg string
type noteAddedMsg bool
type recallMsg string
type notesMsg []notes.Note
type deleteNoteMsg bool

func readInput(ch <-chan string) tea.Cmd {
	return func() tea.Msg {
		return inputMsg(<-ch)
	}
}

func readDeleteInput(m Model, ch <-chan notes.Note) tea.Cmd {
	return func() tea.Msg {
		note := <-ch
		err := m.notes.DeleteNote(note.ID)
		if err != nil {
			log.Println("Error creating note:", err)
			return deleteNoteMsg(false)
		}
		return deleteNoteMsg(true)
	}
}

func addNoteCmd(m Model, note string) tea.Cmd {
	return func() tea.Msg {
		err := m.notes.CreateNote(note)
		if err != nil {
			log.Println("Error creating note:", err)
			return noteAddedMsg(false)
		}
		return noteAddedMsg(true)
	}
}

func recallCmd(m Model, query string) tea.Cmd {
	return func() tea.Msg {
		content, err := m.notes.Recall(query)
		if err != nil {
			log.Println("Error creating note:", err)
			return recallMsg("")
		}
		return recallMsg(content)
	}
}

func getNotesCmd(m Model) tea.Cmd {
	return func() tea.Msg {
		notes, err := m.notes.GetAllNotes()
		if err != nil {
			log.Println("Error getting notes:", err)
			return notesMsg(nil)
		}
		return notesMsg(notes)
	}
}

func saveConfigCmd(m Model) tea.Cmd {
	return func() tea.Msg {
		configDir, _ := utils.GetConfigDir()
		configFile := filepath.Join(configDir, "config.toml")
		err := config.SaveConfig(configFile, m.config)
		if err != nil {
			log.Println("Error saving config:", err)
			// return saveConfigMsg(false)
		}
		return nil
	}
}
