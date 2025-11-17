package searchnotes

import (
	"context"
	"log"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/yagnikpt/flashback/internal/models"
	"golang.org/x/term"
)

type searchResultsMsg []models.FlashbackWithMetadata
type relayChooseMsg models.FlashbackWithMetadata
type dimensionsMsg struct {
	width  int
	height int
}

func searchNotesCmd(m Model, query string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		embeddings, _ := m.app.GenerateEmbeddingForNote(ctx, query, "RETRIEVAL_QUERY")
		flashbacks, err := m.app.RetrieveNotesBySimilarity(ctx, embeddings)
		if err != nil {
			log.Fatal("Error retrieving notes:", err)
		}
		return searchResultsMsg(flashbacks)
	}
}

func relayChooseCmd(note models.FlashbackWithMetadata) tea.Cmd {
	return func() tea.Msg {
		return relayChooseMsg(note)
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
