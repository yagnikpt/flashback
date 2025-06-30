package apiprompt

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/yagnik-patel-47/flashback/internal/components/textarea"
)

var mainViewStyles = lipgloss.NewStyle().Margin(1, 2)

type Model struct {
	textarea textarea.Model
	Output   string
}

func NewModel() Model {
	textarea := textarea.NewModel()
	textarea.SetPlaceholder("Enter your Gemini API key...")

	return Model{
		textarea: textarea,
		Output:   "",
	}
}

func (m Model) Init() tea.Cmd {
	return m.textarea.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "enter":
			if m.textarea.Value() != "" {
				m.Output = m.textarea.Value()
				return m, tea.Quit
			}
		}
	}

	var cmd tea.Cmd
	m.textarea, cmd = m.textarea.Update(msg)

	return m, cmd
}

func (m Model) View() string {
	header := "⚡Flashback"
	var tui strings.Builder
	tui.WriteString(header + "\n\n" + m.textarea.View() + "\n\n" + "Get your Gemini API key from https://aistudio.google.com/apikey")
	return tui.String()
}
