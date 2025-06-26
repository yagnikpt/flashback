package model

import (
	"log"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
	ta "github.com/yagnik-patel-47/flashback/internal/components/textarea"
)

type Model struct {
	mode     string
	textarea ta.Model
	Output   string
}

func InitModel() Model {
	return Model{
		mode:     "input",
		textarea: ta.NewModel(),
	}
}

type inputMsg string

func readInput(ch <-chan string) tea.Cmd {
	return func() tea.Msg {
		return inputMsg(<-ch)
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(tea.SetWindowTitle("flashback"), m.textarea.Init(), readInput(m.textarea.InputChan))
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case inputMsg:
		m.Output = string(msg)
		log.Println(m.Output)
		cmds = append(cmds, readInput(m.textarea.InputChan))
		// return m, tea.Batch(cmds...)
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
			// case "enter":
			// 	m.Output = m.textarea.Value()
			// 	return m, tea.Quit
		}
	}

	m.textarea, cmd = m.textarea.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	header := "⚡Flashback"
	mainViewStyles := lipgloss.NewStyle().Margin(1, 2)
	return mainViewStyles.Render(header + "\n\n" + m.textarea.View())
}
