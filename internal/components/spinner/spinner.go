package spinner

import (
	"fmt"

	"github.com/charmbracelet/bubbles/v2/spinner"
	tea "github.com/charmbracelet/bubbletea/v2"
)

type Model struct {
	spinner     spinner.Model
	displayText string
}

func (m Model) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var (
		cmd tea.Cmd
	)

	m.spinner, cmd = m.spinner.Update(msg)

	return m, cmd
}

func (m Model) View() string {
	if m.displayText != "" {
		str := fmt.Sprintf("%s %s\n", m.spinner.View(), m.displayText)
		return str
	}
	return m.spinner.View() + "\n"
}

func (m *Model) SetDisplayText(text string) {
	m.displayText = text
}

func NewModel() Model {
	model := spinner.New()
	model.Spinner = spinner.MiniDot

	return Model{
		spinner:     model,
		displayText: "",
	}
}
