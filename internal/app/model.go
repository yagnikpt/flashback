package app

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"

	"github.com/yagnikpt/flashback/internal/components/help"
	"github.com/yagnikpt/flashback/internal/components/note"
	"github.com/yagnikpt/flashback/internal/components/notelist"
	"github.com/yagnikpt/flashback/internal/components/recall"
	"github.com/yagnikpt/flashback/internal/global"
)

var (
	primaryBackground   = lipgloss.NewStyle().Background(lipgloss.Color("#475569")).Padding(0, 1)
	secondaryBackground = lipgloss.NewStyle().Background(lipgloss.Color("#475569")).Padding(0, 1)
	dangerBackground    = lipgloss.NewStyle().Background(lipgloss.Color("#9f1239")).Padding(0, 1)
	successStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("#86AF80"))
	mainViewStyles      = lipgloss.NewStyle().Margin(1, 2)
)

type Model struct {
	store    *global.Store
	note     note.Model
	recall   recall.Model
	notelist notelist.Model
	help     help.Model
}

func InitModel() Model {
	store := global.GetStore()

	return Model{
		store:    store,
		note:     note.NewModel(),
		recall:   recall.NewModel(),
		notelist: notelist.NewModel(),
		help:     help.NewModel(),
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(tea.SetWindowTitle("flashback"), m.note.Init(), m.help.Init())
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "tab":
			switch m.store.Mode {
			case "note":
				cmds = append(cmds, m.recall.Init())
				m.store.Mode = "recall"
			case "recall":
				m.store.Mode = "delete"
				cmds = append(cmds, m.notelist.Init())
			case "delete":
				cmds = append(cmds, m.note.Init())
				m.store.Mode = "note"
			case "edit":
				cmds = append(cmds, m.note.Init())
				m.store.Mode = "note"
			}
		case "alt+?":
			if m.store.Config.ShowHelp {
				m.store.Config.ShowHelp = false
			} else {
				m.store.Config.ShowHelp = true
			}
			cmds = append(cmds, saveConfigCmd(m))
		}

	case tea.WindowSizeMsg:
		m.store.Width = msg.Width
		m.store.Height = msg.Height
	}

	switch m.store.Mode {
	case "note":
		m.note, cmd = m.note.Update(msg)
		cmds = append(cmds, cmd)
	case "recall":
		m.recall, cmd = m.recall.Update(msg)
		cmds = append(cmds, cmd)
		// case "delete":
		// case "edit":
		// 	m.notelist, cmd = m.notelist.Update(msg)
		// 	cmds = append(cmds, cmd)
	}
	m.notelist, cmd = m.notelist.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	header := "⚡Flashback"
	modeIndicator := "Mode:" + " "
	switch m.store.Mode {
	case "note":
		modeIndicator += primaryBackground.Render(m.store.Mode)
	case "recall":
		modeIndicator += secondaryBackground.Render(m.store.Mode)
	case "delete":
		modeIndicator += dangerBackground.Render(m.store.Mode)
	}

	var tui strings.Builder

	tui.WriteString(header + "\n\n" + modeIndicator + "\n\n")
	switch m.store.Mode {
	case "note":
		tui.WriteString(m.note.View() + "\n\n")
	case "recall":
		tui.WriteString(m.recall.View() + "\n\n")
	case "delete":
		tui.WriteString(m.notelist.View() + "\n\n")
	case "edit":
		tui.WriteString(m.notelist.View() + "\n\n")
	}

	// if m.height != 0 {
	// voidLines := m.height - 2 - strings.Count(tui.String(), "\n") - strings.Count(m.help.View(), "\n") - 1

	if m.store.Config.ShowHelp {
		tui.WriteString(strings.Repeat("\n", 5) + m.help.View())
	}
	// }

	return mainViewStyles.Render(tui.String())
}
