package spinner

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/v2/spinner"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/muesli/reflow/wordwrap"
)

type Model struct {
	spinner     spinner.Model
	displayText string
	width       int
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
		label := wordwrap.String(m.displayText, m.width-6)
		label = strings.ReplaceAll(label, "\n", "\n  ")
		str := fmt.Sprintf("%s %s\n", m.spinner.View(), label)
		return str
	}
	return m.spinner.View() + "\n"
}

func (m *Model) SetDisplayText(text string) {
	m.displayText = text
}

func (m *Model) SetWidth(width int) {
	m.width = width
}

func NewModel() Model {
	model := spinner.New()
	model.Spinner = spinner.MiniDot

	return Model{
		spinner:     model,
		displayText: "",
		width:       0,
	}
}
