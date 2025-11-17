package notelist

import (
	"context"
	"log"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/yagnikpt/flashback/internal/models"
	"golang.org/x/term"
)

type getAllNotesMsg []models.FlashbackWithMetadata
type deleteNoteMsg bool
type chosenNoteMsg models.FlashbackWithMetadata
type relayChooseMsg string
type relayDeleteMsg string
type dimensionsMsg struct {
	width  int
	height int
}

func getAllNotesCmd(m Model) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		notes, err := m.app.GetAllNotes(ctx)
		if err != nil {
			log.Fatal(err)
		}
		return getAllNotesMsg(notes)
	}
}

func chooseNoteCmd(m Model, noteID string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		note, err := m.app.GetNoteByID(ctx, noteID)
		if err != nil {
			log.Println("Error finding note:", err)
			return chosenNoteMsg(note)
		}
		return chosenNoteMsg(note)
	}
}

func deleteNoteCmd(m Model, noteID string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := m.app.DeleteNoteByID(ctx, noteID)
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

func getDimensionsCmd() tea.Cmd {
	return func() tea.Msg {
		width, height, err := term.GetSize(int(os.Stdout.Fd()))
		if err != nil {
			panic(err)
		}
		return dimensionsMsg{
			width:  width,
			height: height,
		}
	}
}
