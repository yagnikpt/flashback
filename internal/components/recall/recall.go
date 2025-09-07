package recall

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/muesli/reflow/wordwrap"
	"github.com/yagnikpt/flashback/internal/components/spinner"
	"github.com/yagnikpt/flashback/internal/components/textarea"
	"github.com/yagnikpt/flashback/internal/global"
)

type Model struct {
	textarea textarea.Model
	spinner  spinner.Model
	store    *global.Store
	output   string
}

func NewModel() Model {
	store := global.GetStore()
	textarea := textarea.NewModel()
	textarea.SetHeight(3)
	textarea.SetPlaceholder("Enter query to recall...")
	return Model{
		textarea: textarea,
		spinner:  spinner.NewModel(),
		output:   "",
		store:    store,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.textarea.Init(), m.spinner.Init())
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case recallMsg:
		m.store.Loading = false
		content := string(msg)
		if content != "" {
			m.store.ShowFeedback = true
			m.output = content
		} else {
			m.store.ShowFeedback = true
			m.output = "No notes found."
		}
		m.textarea.Focus()
		m.spinner.SetDisplayText("")

	case tea.KeyPressMsg:
		switch msg.String() {
		case "enter":
			if m.store.Loading {
				return m, nil
			}
			if m.textarea.Value() != "" {
				m.store.Loading = true
				cmds = append(cmds, recallCmd(m, m.textarea.Value()))
				m.textarea.SetValue("")
				m.spinner.SetDisplayText("Recalling note...")
			}
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

	if m.store.Loading {
		tui.WriteString(m.spinner.View() + "\n\n")
	} else {
		tui.WriteString(m.textarea.View() + "\n\n")
	}

	if m.store.ShowFeedback {
		wrappedOutput := wordwrap.String(m.output, m.store.Width-4)
		tui.WriteString(wrappedOutput + "\n\n")
	}

	return tui.String()
}
