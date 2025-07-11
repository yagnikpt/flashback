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
	"github.com/yagnik-patel-47/flashback/internal/notes"
	"github.com/yagnik-patel-47/flashback/internal/utils"
)

type Model struct {
	items      []notes.CombinedNote
	cursor     int
	paginator  paginator.Model
	width      int
	height     int
	OutputChan chan notes.CombinedNote
}

func (m *Model) SetDimensions(width, height int) {
	m.width = width
	m.height = height
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
	Left: key.NewBinding(
		key.WithKeys("h", "left", "pgup"),
		key.WithHelp("←/h/pgdn", "move left"),
	),
	Right: key.NewBinding(
		key.WithKeys("l", "right", "pgdn"),
		key.WithHelp("→/l/pgup", "move right"),
	),
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd

	if len(m.items) == 0 {
		m.cursor = 0
		m.paginator.Page = 0
		return m, cmd
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
				m.OutputChan <- currentPageItems[m.cursor]
			}
		}
	}

	m.paginator, cmd = m.paginator.Update(msg)
	return m, cmd
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
		content := wordwrap.String(item.Content, m.width-6)
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

func (m *Model) SetItems(items []notes.CombinedNote) {
	m.items = items
	m.paginator.SetTotalPages(len(items))
}

func NewModel() Model {
	pagi := paginator.New()
	pagi.PerPage = 5
	pagi.Type = paginator.Dots
	pagi.ActiveDot = lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(252)).Render("•")
	pagi.InactiveDot = lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(238)).Render("•")
	pagi.SetTotalPages(0)

	return Model{
		items:      []notes.CombinedNote{},
		cursor:     0,
		paginator:  pagi,
		OutputChan: make(chan notes.CombinedNote),
		width:      0,
		height:     0,
	}
}
