package apikeyinput

import (
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
)

type Model struct {
	input  textinput.Model
	Output string
}

func NewModel() Model {
	input := textinput.New()
	input.SetWidth(50)
	input.Placeholder = "Enter your Gemini API key"
	input.EchoMode = textinput.EchoPassword
	input.Focus()

	return Model{
		input:  input,
		Output: "",
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "enter":
			if m.input.Value() != "" {
				m.Output = m.input.Value()
				return m, tea.Quit
			}
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)

	return m, cmd
}

func (m Model) View() tea.View {
	var tui strings.Builder
	tui.WriteString("\n")
	tui.WriteString(m.input.View())
	tui.WriteString("\n\n")
	tui.WriteString("Get your Gemini API key from https://aistudio.google.com/api-keys")
	tui.WriteString("\n\n")
	return tea.NewView(tui.String())
}
