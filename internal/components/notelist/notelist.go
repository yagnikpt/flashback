package notelist

import (
	"time"

	"github.com/charmbracelet/bubbles/v2/list"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
)

type item struct {
	content   string
	createdAt time.Time
}

func (i item) Content() string     { return i.content }
func (i item) CreatedAt() string   { return i.createdAt.Format(time.RFC3339) }
func (i item) FilterValue() string { return i.content }

type Model struct {
	list     list.Model
	docStyle lipgloss.Style
}

func NewModel() Model {
	items := []list.Item{
		item{content: "Note 1", createdAt: time.Now()},
		item{content: "Note 2", createdAt: time.Now()},
		item{content: "Note 3", createdAt: time.Now()},
	}
	m := Model{list: list.New(items, list.NewDefaultDelegate(), 0, 0)}
	m.list.Title = "Remove notes"

	m.docStyle = lipgloss.NewStyle().Margin(1, 2)
	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	// case tea.KeyMsg:
	// if msg.String() == "ctrl+c" {
	// 	return m, tea.Quit
	// }
	case tea.WindowSizeMsg:
		h, v := m.docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return m.docStyle.Render(m.list.View())
}
