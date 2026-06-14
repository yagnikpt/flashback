package tui

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/yagnikpt/flashback/internal/app"
	"github.com/yagnikpt/flashback/internal/components/insertnote"
	"github.com/yagnikpt/flashback/internal/components/notelist"
	"github.com/yagnikpt/flashback/internal/components/searchnotes"
)

type Model struct {
	app         *app.App
	active      Screen
	notelist    notelist.Model
	insertnote  insertnote.Model
	searchnotes searchnotes.Model
}

type Screen int

const (
	screenListNotes Screen = iota
	screenInsertNote
	screenSearchNotes
)

func NewModel(app *app.App) Model {
	return Model{
		app:         app,
		active:      screenListNotes,
		notelist:    notelist.NewModel(app),
		insertnote:  insertnote.NewModel(app),
		searchnotes: searchnotes.NewModel(app),
	}
}

func (m Model) Init() tea.Cmd {
	return m.notelist.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "tab":
			m.active = (m.active + 1) % 3
			switch m.active {
			case screenListNotes:
				cmd = m.notelist.Init()
				m.notelist.ResetView()
			case screenInsertNote:
				cmd = m.insertnote.Init()
				m.insertnote.ResetView()
			case screenSearchNotes:
				cmd = m.searchnotes.Init()
				m.searchnotes.ResetView()
			}
			return m, cmd
		case "shift+tab":
			m.active = ((m.active-1)%3 + 3) % 3
			switch m.active {
			case screenListNotes:
				cmd = m.notelist.Init()
				m.notelist.ResetView()
			case screenInsertNote:
				cmd = m.insertnote.Init()
				m.insertnote.ResetView()
			case screenSearchNotes:
				cmd = m.searchnotes.Init()
				m.searchnotes.ResetView()
			}
			return m, cmd
		}
	}

	switch m.active {
	case screenListNotes:
		newNotelist, cmd := m.notelist.Update(msg)
		m.notelist = newNotelist.(notelist.Model)
		cmds = append(cmds, cmd)
	case screenInsertNote:
		newInsertnote, cmd := m.insertnote.Update(msg)
		m.insertnote = newInsertnote.(insertnote.Model)
		cmds = append(cmds, cmd)
	case screenSearchNotes:
		newSearchnotes, cmd := m.searchnotes.Update(msg)
		m.searchnotes = newSearchnotes.(searchnotes.Model)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

var (
	tabStyles       = lipgloss.NewStyle().Padding(0, 1).Margin(0, 1).Render
	activeTabStyles = lipgloss.NewStyle().Padding(0, 1).Margin(0, 1).Background(lipgloss.Blue).Foreground(lipgloss.White).Bold(true).Render
)

func (m Model) View() tea.View {
	views := []string{"Manage Notes", "Add Note", "Search Notes"}
	var builder strings.Builder
	for i, v := range views {
		if Screen(i) == m.active {
			builder.WriteString(activeTabStyles(v))
		} else {
			builder.WriteString(tabStyles(v))
		}
	}
	builder.WriteString("\n\n")
	switch m.active {
	case screenListNotes:
		builder.WriteString(m.notelist.View().Content)
	case screenInsertNote:
		builder.WriteString(m.insertnote.View().Content)
	case screenSearchNotes:
		builder.WriteString(m.searchnotes.View().Content)
	}

	v := tea.NewView(builder.String())
	v.AltScreen = true
	v.WindowTitle = "flashback"
	return v
}
