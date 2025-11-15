package spinner

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/reflow/wordwrap"
)

type Model struct {
	spinner     spinner.Model
	displayText string
	width       int
	status      <-chan string
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, getCurrentStatus(m))
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case statusMsg:
		m.displayText = string(msg)
		cmds = append(cmds, getCurrentStatus(m))
	case closedMsg:
		return m, tea.Quit
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	}

	m.spinner, cmd = m.spinner.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.displayText != "" {
		label := wordwrap.String(m.displayText, m.width-6)
		label = strings.ReplaceAll(label, "\n", "\n  ")
		str := fmt.Sprintf("\n%s %s\n", m.spinner.View(), label)
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

func NewModel(status <-chan string) Model {
	model := spinner.New()
	model.Spinner = spinner.MiniDot

	return Model{
		spinner:     model,
		displayText: "",
		width:       0,
		status:      status,
	}
}
