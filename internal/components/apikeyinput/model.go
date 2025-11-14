package apikeyinput

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/yagnikpt/flashback/internal/components/textarea"
)

var mainViewStyles = lipgloss.NewStyle().Margin(1, 2)

type Model struct {
	textarea textarea.Model
	Output   string
}

func NewModel() Model {
	textarea := textarea.NewModel()
	textarea.SetPlaceholder("Enter your Gemini API key...")
	textarea.SetHeight(1)

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
	case tea.KeyMsg:
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
	var tui strings.Builder
	tui.WriteString("\n" + m.textarea.View() + "\n\n" + "Get your Gemini API key from https://aistudio.google.com/api-keys")
	return tui.String()
}
