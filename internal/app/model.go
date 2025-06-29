package app

import (
	"database/sql"
	"strings"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"

	"github.com/yagnik-patel-47/flashback/internal/components/notelist"
	"github.com/yagnik-patel-47/flashback/internal/components/spinner"
	"github.com/yagnik-patel-47/flashback/internal/components/textarea"

	"github.com/yagnik-patel-47/flashback/internal/notes"
)

type Model struct {
	mode         string
	textarea     textarea.Model
	spinner      spinner.Model
	notelist     notelist.Model
	notes        notes.Store
	showFeedback bool
	loading      bool
	output       string
}

func InitModel(db *sql.DB) Model {
	return Model{
		mode:     "note",
		textarea: textarea.NewModel(),
		spinner:  spinner.NewModel(),
		notelist: notelist.NewModel(),
		output:   "",
		notes:    *notes.NewStore(db),
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(tea.SetWindowTitle("flashback"), m.textarea.Init(), m.spinner.Init(), readInput(m.textarea.InputChan))
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case noteAddedMsg:
		m.loading = false
		status := bool(msg)
		if status {
			m.showFeedback = true
			m.output = "Note saved successfully."
			m.textarea.SetPlaceholder("Add a new note...")
		} else {
			m.showFeedback = true
			m.output = "Error saving note."
			m.textarea.SetPlaceholder("Try again...")
		}
		m.textarea.Focus()
		m.spinner.SetDisplayText("")

	case recallMsg:
		m.loading = false
		content := string(msg)
		m.textarea.Focus()
		m.spinner.SetDisplayText("")
		if content != "" {
			m.showFeedback = true
			m.output = content
		} else {
			m.showFeedback = true
			m.output = "No notes found."
		}

	case inputMsg:
		content := string(msg)
		m.showFeedback = false
		m.loading = true
		m.textarea.Blur()
		switch m.mode {
		case "note":
			m.spinner.SetDisplayText("Creating note...")
			cmds = append(cmds, addNoteCmd(m, content))
		case "recall":
			m.spinner.SetDisplayText("Recalling notes...")
			cmds = append(cmds, recallCmd(m, content))
		}

		cmds = append(cmds, readInput(m.textarea.InputChan))
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			m.showFeedback = false
			if !m.loading {
				switch m.mode {
				case "note":
					m.mode = "recall"
					m.textarea.SetPlaceholder("Enter query to recall...")
				case "recall":
					cmds = append(cmds, getNotesCmd(m))
					m.mode = "delete"
				case "delete":
					m.mode = "note"
					m.textarea.SetPlaceholder("Enter the note...")
				}
			}
		}
	}

	m.textarea, cmd = m.textarea.Update(msg)
	cmds = append(cmds, cmd)
	m.spinner, cmd = m.spinner.Update(msg)
	cmds = append(cmds, cmd)
	m.notelist, cmd = m.notelist.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

var (
	primaryBackground = lipgloss.NewStyle().Background(lipgloss.Color("#5b6b6d")).Padding(0, 1)
	dangerBackground  = lipgloss.NewStyle().Background(lipgloss.Color("#9f1239")).Padding(0, 1)
	successStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("#86AF80"))
	mainViewStyles    = lipgloss.NewStyle().Margin(1, 2)
)

func (m Model) View() string {
	header := "⚡Flashback"
	modeIndicator := "Mode:" + " "
	if m.mode == "delete" {
		modeIndicator += dangerBackground.Render(m.mode)
	} else {
		modeIndicator += primaryBackground.Render(m.mode)
	}

	var tui strings.Builder

	tui.WriteString(header + "\n\n" + modeIndicator + "\n\n")
	if m.loading {
		tui.WriteString(m.spinner.View() + "\n\n")
	} else {
		if m.mode == "delete" {
			tui.WriteString(m.notelist.View() + "\n\n")
		} else {
			tui.WriteString(m.textarea.View() + "\n\n")
		}
	}
	if m.showFeedback {
		if m.mode == "note" {
			tui.WriteString(successStyle.Render(m.output) + "\n\n")
		} else {
			tui.WriteString(m.output + "\n\n")
		}
	}

	return mainViewStyles.Render(tui.String())
}
