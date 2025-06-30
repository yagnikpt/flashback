package app

import (
	"database/sql"
	"strings"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/muesli/reflow/wordwrap"

	"github.com/yagnik-patel-47/flashback/internal/components/help"
	"github.com/yagnik-patel-47/flashback/internal/components/notelist"
	"github.com/yagnik-patel-47/flashback/internal/components/spinner"
	"github.com/yagnik-patel-47/flashback/internal/components/textarea"
	"github.com/yagnik-patel-47/flashback/internal/config"

	"github.com/yagnik-patel-47/flashback/internal/notes"
)

var (
	primaryBackground = lipgloss.NewStyle().Background(lipgloss.Color("#5b6b6d")).Padding(0, 1)
	dangerBackground  = lipgloss.NewStyle().Background(lipgloss.Color("#9f1239")).Padding(0, 1)
	successStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("#86AF80"))
	mainViewStyles    = lipgloss.NewStyle().Margin(1, 2)
)

type Model struct {
	mode         string
	textarea     textarea.Model
	spinner      spinner.Model
	notelist     notelist.Model
	help         help.Model
	notes        notes.Store
	width        int
	height       int
	showFeedback bool
	loading      bool
	output       string
	config       config.Config
}

func InitModel(db *sql.DB, config config.Config) Model {
	return Model{
		mode:     "note",
		textarea: textarea.NewModel(),
		spinner:  spinner.NewModel(),
		notelist: notelist.NewModel(),
		help:     help.NewModel(),
		output:   "",
		notes:    *notes.NewStore(db, config.APIKey),
		width:    0,
		height:   0,
		config:   config,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(tea.SetWindowTitle("flashback"), m.textarea.Init(), m.spinner.Init(), readInput(m.textarea.OutputChan), readDeleteInput(m, m.notelist.OutputChan))
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

		cmds = append(cmds, readInput(m.textarea.OutputChan))

	case notesMsg:
		notes := []notes.Note(msg)
		m.notelist.SetItems(notes)

	case deleteNoteMsg:
		success := bool(msg)
		if success {
			cmds = append(cmds, getNotesCmd(m))
		}
		cmds = append(cmds, readDeleteInput(m, m.notelist.OutputChan))

	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "tab":
			m.showFeedback = false
			if !m.loading {
				switch m.mode {
				case "note":
					m.mode = "recall"
					m.help.SetMode(m.mode)
					m.textarea.SetPlaceholder("Enter query to recall...")
				case "recall":
					cmds = append(cmds, getNotesCmd(m))
					m.mode = "delete"
					m.help.SetMode(m.mode)
				case "delete":
					m.mode = "note"
					m.help.SetMode(m.mode)
					m.textarea.SetPlaceholder("Enter the note...")
				}
			}
		case "alt+?":
			if m.config.ShowHelp {
				m.config.ShowHelp = false
			} else {
				m.config.ShowHelp = true
			}
			cmds = append(cmds, saveConfigCmd(m))
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	if m.mode == "note" || m.mode == "recall" {
		m.textarea, cmd = m.textarea.Update(msg)
		cmds = append(cmds, cmd)
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
	} else {
		m.notelist, cmd = m.notelist.Update(msg)
		cmds = append(cmds, cmd)
	}
	if m.config.ShowHelp {
		m.help, cmd = m.help.Update(msg)
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

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
		wrappedOutput := wordwrap.String(m.output, m.width-4)
		if m.mode == "note" {
			tui.WriteString(successStyle.Render(wrappedOutput) + "\n\n")
		} else {
			tui.WriteString(wrappedOutput + "\n\n")
		}
	}

	// if m.height != 0 {
	// voidLines := m.height - 2 - strings.Count(tui.String(), "\n") - strings.Count(m.help.View(), "\n") - 1

	if m.config.ShowHelp {
		tui.WriteString(strings.Repeat("\n", 5) + m.help.View())
	}
	// }

	return mainViewStyles.Render(tui.String())
}
