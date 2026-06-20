package spinner

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type Model struct {
	spinner     spinner.Model
	displayText string
	width       int
	status      <-chan string
	altScreen   bool
}

func (m Model) Init() tea.Cmd {
	cmds := []tea.Cmd{}
	cmds = append(cmds, m.spinner.Tick)
	if m.status != nil {
		cmds = append(cmds, getCurrentStatus(m))
	}
	return tea.Batch(cmds...)
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
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	}

	m.spinner, cmd = m.spinner.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() tea.View {
	if m.displayText != "" {
		label := lipgloss.Wrap(m.displayText, m.width-6, " ")
		label = strings.ReplaceAll(label, "\n", "\n  ")
		str := fmt.Sprintf("\n%s %s\n", m.spinner.View(), label)
		v := tea.NewView(str)
		v.AltScreen = m.altScreen
		return v
	}
	v := tea.NewView(m.spinner.View() + "\n")
	v.AltScreen = m.altScreen
	return v
}

func (m *Model) SetDisplayText(text string) {
	m.displayText = text
}

func (m *Model) SetWidth(width int) {
	m.width = width
}

func (m *Model) SetAltScreen(altScreen bool) {
	m.altScreen = altScreen
}

func NewModel(status <-chan string) Model {
	model := spinner.New()
	model.Spinner = spinner.MiniDot

	return Model{
		spinner:     model,
		displayText: "",
		width:       0,
		status:      status,
		altScreen:   false,
	}
}
