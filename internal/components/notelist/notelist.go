package notelist

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/v2/key"
	"github.com/charmbracelet/bubbles/v2/paginator"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/muesli/reflow/wordwrap"
	"github.com/yagnikpt/flashback/internal/global"
	"github.com/yagnikpt/flashback/internal/notes"
	"github.com/yagnikpt/flashback/internal/utils"
)

type Model struct {
	store     *global.Store
	items     []notes.CombinedNote
	cursor    int
	paginator paginator.Model
}

var listContainerStyle = lipgloss.NewStyle()
var lightTextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#5e5e5e"))
var cursorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#f43f5e")).Bold(true)

type KeyMap struct {
	Up    key.Binding
	Down  key.Binding
	Left  key.Binding
	Right key.Binding
}

var listNavigation = KeyMap{
	Up: key.NewBinding(
		key.WithKeys("k", "up"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("j", "down"),
		key.WithHelp("↓/j", "move down"),
	),
}

func (m Model) Init() tea.Cmd {
	return getNotesCmd(m)
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	if len(m.items) == 0 {
		m.cursor = 0
		m.paginator.Page = 0
	}

	if m.paginator.Page*m.paginator.PerPage >= len(m.items) {
		m.paginator.Page = (len(m.items) - 1) / m.paginator.PerPage
	}

	start, end := m.paginator.GetSliceBounds(len(m.items))
	if start >= len(m.items) {
		start = 0
		m.paginator.Page = 0
	}
	if end > len(m.items) {
		end = len(m.items)
	}
	currentPageItems := m.items[start:end]

	if len(currentPageItems) == 0 {
		m.cursor = 0
	} else if m.cursor >= len(currentPageItems) {
		m.cursor = len(currentPageItems) - 1
	}

	switch msg := msg.(type) {
	case notesMsg:
		notes := []notes.CombinedNote(msg)
		m.items = notes
		m.paginator.SetTotalPages(len(notes))

	case deleteNoteMsg:
		success := bool(msg)
		if success {
			cmds = append(cmds, getNotesCmd(m))
		}

	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, listNavigation.Up):
			if m.cursor > 0 {
				m.cursor--
			} else {
				if m.paginator.Page > 0 {
					m.paginator.PrevPage()
					start, end = m.paginator.GetSliceBounds(len(m.items))
					m.cursor = len(m.items[start:end]) - 1
				}
			}
		case key.Matches(msg, listNavigation.Down):
			if m.cursor < len(currentPageItems)-1 {
				m.cursor++
			} else if m.paginator.Page < m.paginator.TotalPages-1 {
				m.paginator.NextPage()
				m.cursor = 0
			}
		}
		switch msg.String() {
		case "enter":
			start, end := m.paginator.GetSliceBounds(len(m.items))
			if start >= len(m.items) {
				start = 0
				m.paginator.Page = 0
			}
			if end > len(m.items) {
				end = len(m.items)
			}
			currentPageItems := m.items[start:end]
			if len(currentPageItems) > 0 {
				cmds = append(cmds, deleteNoteCmd(m, currentPageItems[m.cursor].ID))
			}
		}
	}

	m.paginator, cmd = m.paginator.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

var loc, _ = time.LoadLocation("Local")

func (m Model) View() string {
	if len(m.items) == 0 {
		return "No notes found!"
	}

	var b strings.Builder
	b.WriteString(lightTextStyle.Render(fmt.Sprintf("%d notes", len(m.items))) + "\n\n")
	start, end := m.paginator.GetSliceBounds(len(m.items))
	if start >= len(m.items) {
		start = 0
	}
	if end > len(m.items) {
		end = len(m.items)
	}

	var liView strings.Builder
	displayItems := m.items[start:end]

	for i, item := range displayItems {
		cursor := " "
		if m.cursor == i {
			cursor = cursorStyle.Render(">")
		}
		content := wordwrap.String(item.Content, m.store.Width-6)
		content = strings.ReplaceAll(content, "\n", "\n  ")
		relativeTime := utils.RelativeTime(item.CreatedAt.In(loc))

		liView.WriteString(fmt.Sprintf("%s %s\n  %s\n\n", cursor, content, lightTextStyle.Render(relativeTime)))
	}

	b.WriteString(listContainerStyle.Render(liView.String()))
	if len(m.items) > 5 {
		b.WriteString("  " + m.paginator.View())
	}

	return b.String()
}

func NewModel() Model {
	store := global.GetStore()
	p := paginator.New()
	p.PerPage = 5
	p.Type = paginator.Dots
	p.ActiveDot = lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(252)).Render("•")
	p.InactiveDot = lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(238)).Render("•")
	p.SetTotalPages(0)

	return Model{
		store:     store,
		items:     []notes.CombinedNote{},
		cursor:    0,
		paginator: p,
	}
}
