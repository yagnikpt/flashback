package note

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/muesli/reflow/wordwrap"
	"github.com/yagnikpt/flashback/internal/components/spinner"
	"github.com/yagnikpt/flashback/internal/components/textarea"
	"github.com/yagnikpt/flashback/internal/global"
)

type Model struct {
	textarea     textarea.Model
	spinner      spinner.Model
	store        *global.Store
	loading      bool
	showFeedback bool
	output       string
}

func NewModel() Model {
	store := global.GetStore()
	textarea := textarea.NewModel()
	textarea.SetHeight(5)
	textarea.SetPlaceholder("Enter the note...")
	return Model{
		textarea:     textarea,
		spinner:      spinner.NewModel(),
		output:       "",
		store:        store,
		loading:      false,
		showFeedback: false,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.textarea.Init(), m.spinner.Init())
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
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
	case statusMsg:
		status := string(msg)
		m.spinner.SetDisplayText(status)
		cmds = append(cmds, readStatusText(m.store.Notes.StatusChan))
	case tea.KeyPressMsg:
		switch msg.String() {
		case "enter":
			if m.loading {
				return m, nil
			}
			if m.textarea.Value() != "" {
				cmds = append(cmds, readStatusText(m.store.Notes.StatusChan), addNoteCmd(m, m.textarea.Value()))
				m.textarea.SetValue("")
			}
			m.loading = true
			m.spinner.SetDisplayText("Creating note...")
			cmds = append(cmds, addNoteCmd(m, m.textarea.Value()))
		}
	}

	m.textarea, cmd = m.textarea.Update(msg)
	cmds = append(cmds, cmd)
	m.spinner, cmd = m.spinner.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	var tui strings.Builder

	if m.loading {
		tui.WriteString(m.spinner.View() + "\n\n")
	} else {
		tui.WriteString(m.textarea.View() + "\n\n")
	}

	if m.showFeedback {
		wrappedOutput := wordwrap.String(m.output, m.store.Width-4)
		tui.WriteString(wrappedOutput + "\n\n")
	}

	return tui.String()
}
