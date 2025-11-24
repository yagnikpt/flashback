package insertnote

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/yagnikpt/flashback/internal/contentloaders"
	"golang.org/x/term"
)

type heightMsg int
type addNoteMsg struct {
	success bool
	err     error
}

func addNoteCmd(m Model, content string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
		defer cancel()

		var metadata map[string]string
		var noteType string
		if strings.HasPrefix(content, "http") {
			m.statusChan <- "Fetching webpage content..."
			noteType = "link"
			pageContent, err := contentloaders.GetWebPage(ctx, content)
			if err != nil {
				return addNoteMsg{
					success: false,
					err:     err,
				}
			}
			pageContentWithUrl := fmt.Sprintf("URL: %s\n\n%s", content, pageContent)
			m.statusChan <- "Generating metadata for webpage..."
			metadata, err = m.app.GenerateMetadataForWebNote(ctx, pageContentWithUrl)
			if err != nil {
				return addNoteMsg{
					success: false,
					err:     err,
				}
			}
		} else {
			m.statusChan <- "Generating metadata for note..."
			_metadata, err := m.app.GenerateMetadataForSimpleNote(ctx, content)
			if err != nil {
				return addNoteMsg{
					success: false,
					err:     err,
				}
			}
			metadata = _metadata
			noteType = "text"
		}
		var finalContent strings.Builder
		userInput := fmt.Sprintf("USER INPUT:\n content: %s\n", content)
		finalContent.WriteString(userInput)
		finalContent.WriteString("METADATA:\n")
		for key, value := range metadata {
			metadataLine := fmt.Sprintf(" %s: %s\n", key, value)
			finalContent.WriteString(metadataLine)
		}
		m.statusChan <- "Saving the note..."
		embeddings, err := m.app.GenerateEmbeddingForNote(ctx, finalContent.String(), "RETRIEVAL_DOCUMENT")
		if err != nil {
			return addNoteMsg{
				success: false,
				err:     err,
			}
		}
		err = m.app.InsertNote(ctx, content, noteType, metadata, embeddings)
		if err != nil {
			return addNoteMsg{
				success: false,
				err:     err,
			}
		}
		return addNoteMsg{
			success: true,
			err:     nil,
		}
	}
}

func getHeightCmd() tea.Cmd {
	return func() tea.Msg {
		_, height, err := term.GetSize(int(os.Stdout.Fd()))
		if err != nil {
			panic(err)
		}
		return heightMsg(height)
	}
}
