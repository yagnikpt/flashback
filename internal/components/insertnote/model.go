package insertnote

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/yagnikpt/flashback/internal/app"
	"github.com/yagnikpt/flashback/internal/components/spinner"
	"github.com/yagnikpt/flashback/internal/components/textarea"
)

type Model struct {
	app          *app.App
	textarea     textarea.Model
	spinner      spinner.Model
	showFeedback bool
	feedbackMsg  string
	isLoading    bool
	statusChan   chan string
}

func (m *Model) ResetView() {
	m.showFeedback = false
	m.isLoading = false
	m.feedbackMsg = ""
}

func NewModel(app *app.App) Model {
	t := textarea.NewModel()
	statusChan := make(chan string)
	s := spinner.NewModel(statusChan)
	s.SetDisplayText("Creating the note...")

	return Model{
		app:          app,
		textarea:     t,
		spinner:      s,
		showFeedback: false,
		feedbackMsg:  "",
		isLoading:    false,
		statusChan:   statusChan,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.textarea.Init(), m.spinner.Init(), getHeightCmd())
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case addNoteMsg:
		res := addNoteMsg(msg)
		if !res.success {
			m.feedbackMsg = res.err.Error()
		} else {
			m.feedbackMsg = "Note successfully added."
		}
		m.isLoading = false
		m.showFeedback = true
	case heightMsg:
		height := int(msg)
		if height > 24 {
			m.textarea.SetHeight(16)
		} else {
			height = int(75.0*float32(height)) / 100
			m.textarea.SetHeight(height - 4)
		}
	case tea.WindowSizeMsg:
		height := msg.Height
		if height > 24 {
			m.textarea.SetHeight(16)
		} else {
			height = int(75.0*float32(height)) / 100
			m.textarea.SetHeight(height - 4)
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			m.showFeedback = false
			noteContent := m.textarea.Value()
			if noteContent != "" {
				m.isLoading = true
				cmds = append(cmds, addNoteCmd(m, noteContent))
				m.textarea.SetValue("")
				cmds = append(cmds, m.spinner.Init())
			}
		}
	}

	if m.isLoading {
		newSpinner, cmd := m.spinner.Update(msg)
		m.spinner = newSpinner.(spinner.Model)
		cmds = append(cmds, cmd)
	}
	newTextarea, cmd := m.textarea.Update(msg)
	m.textarea = newTextarea.(textarea.Model)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

var (
	docStyles = lipgloss.NewStyle().Margin(1, 1).Render
)

func (m Model) View() string {
	var builder strings.Builder
	if m.isLoading {
		builder.WriteString(m.spinner.View())
	} else {
		builder.WriteString(m.textarea.View())
	}
	if m.showFeedback {
		builder.WriteString("\n\n" + m.feedbackMsg)
	}
	return docStyles(builder.String())
}
