package utils

import (
	"encoding/json"
	"os"
	"strings"

	"golang.org/x/term"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
	"github.com/yagnikpt/flashback/internal/models"
)

var (
	// blockStyles = lipgloss.NewStyle().Padding(0, 2, 0, 2)
	keyStyles = lipgloss.NewStyle().Bold(true)
)

var ignoreList = map[string]bool{
	"image_main": true,
}

func FormatSingleNote(note models.FlashbackWithMetadata) string {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		panic(err)
	}

	result := keyStyles.Render("\nID: ") + note.ID + "\n" + keyStyles.Render("Content: ")

	if note.Type != "url" {
		result += "\n"
		content := wordwrap.String(note.Content, width-4)
		result += content
	} else {
		result += note.Content
	}

	result += "\n\nMetadata:\n"
	for key, value := range note.Metadata {
		if ignoreList[key] {
			continue
		}
		if key == "tags" {
			var tags []string
			err := json.Unmarshal([]byte(value), &tags)
			if err != nil {
				continue
			}
			value = stringJoin(tags, ", ")
		}
		if key == "image" {
			result += "  " + keyStyles.Render(key) + ": " + value + "\n"
			continue
		}
		value = wordwrap.String(value, width-len(key)-8)
		value = strings.ReplaceAll(value, "\n", "\n"+strings.Repeat(" ", 4+len(key)))
		result += "  " + keyStyles.Render(key) + ": " + value + "\n"
	}
	// result = wordwrap.String(result, width)
	return result
}

func FormatSingleNoteCompact(note models.FlashbackWithMetadata) string {
	result := "\nID: " + note.ID + "\nContent: " + note.Content + "\n"
	paddedResult := strings.ReplaceAll(result, "\n", "\n  ")
	return paddedResult
}

func FormatMultipleNotes(notes []models.FlashbackWithMetadata) string {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		panic(err)
	}

	result := ""
	for _, note := range notes {
		result += FormatSingleNote(note)
		result += "\n  " + strings.Repeat("-", width-4) + "  \n"
	}
	return result
}

func FormatMultipleNotesCompact(notes []models.FlashbackWithMetadata) string {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		panic(err)
	}

	header := "ID" + strings.Repeat(" ", 24) + "Content\n\n"
	result := header
	for _, note := range notes {
		wrappedContent := wordwrap.String(note.Content, width-26)
		indentedContent := strings.ReplaceAll(wrappedContent, "\n", "\n"+strings.Repeat(" ", 26))
		result += keyStyles.Render(note.ID) + "    " + indentedContent + "\n\n"
	}
	return result
}

func stringJoin(arr []string, sep string) string {
	result := ""
	for i, str := range arr {
		result += str
		if i < len(arr)-1 {
			result += sep
		}
	}
	return result
}
