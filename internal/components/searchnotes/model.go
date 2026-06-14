package searchnotes

import (
	"strings"
	"time"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/dustin/go-humanize"
	"github.com/yagnikpt/flashback/internal/app"
	"github.com/yagnikpt/flashback/internal/components/spinner"
	"github.com/yagnikpt/flashback/internal/components/textarea"
	"github.com/yagnikpt/flashback/internal/models"
	"github.com/yagnikpt/flashback/internal/utils"
)

type Model struct {
	app          *app.App
	textarea     textarea.Model
	spinner      spinner.Model
	list         list.Model
	showFeedback bool
	isLoading    bool
	showingNote  bool
	activeNote   models.FlashbackWithMetadata
}

type item struct {
	title, desc string
	full        models.FlashbackWithMetadata
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

func NewModel(app *app.App) Model {
	t := textarea.NewModel()
	t.SetHeight(3)
	t.SetPlaceholder("Search the notes...")
	s := spinner.NewModel(nil)
	s.SetDisplayText("Searching the notes...")

	items := make([]list.Item, 0)
	d := newDelegate(newDelegateKeyMap())
	l := list.New(items, d, 0, 0)
	l.SetShowTitle(false)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)

	return Model{
		app:          app,
		textarea:     t,
		spinner:      s,
		list:         l,
		showFeedback: false,
		isLoading:    false,
		showingNote:  false,
		activeNote:   models.FlashbackWithMetadata{},
	}
}

func (m *Model) ResetView() {
	m.showFeedback = false
	m.isLoading = false
	m.activeNote = models.FlashbackWithMetadata{}
	m.showingNote = false
	m.textarea.Focus()
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.textarea.Init(), m.spinner.Init(), getDimensionsCmd())
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case searchResultsMsg:
		notes := []models.FlashbackWithMetadata(msg)
		items := make([]list.Item, len(notes))
		for i := range items {
			t, _ := time.Parse(time.RFC3339, notes[i].CreatedAt)
			items[i] = item{
				full:  notes[i],
				title: notes[i].Content,
				desc:  humanize.Time(t),
			}
		}
		m.list.SetItems(items)
		m.isLoading = false
		m.showFeedback = true
		if len(notes) > 0 {
			m.textarea.Blur()
		}

	case relayChooseMsg:
		note := models.FlashbackWithMetadata(msg)
		m.activeNote = note
		m.showingNote = true

	case dimensionsMsg:
		dims := dimensionsMsg(msg)
		m.list.SetSize(dims.width, dims.height-8)

	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height-8)

	case tea.KeyPressMsg:
		switch msg.String() {
		case "esc":
			if m.showingNote {
				m.showingNote = false
				m.activeNote = models.FlashbackWithMetadata{}
			} else if m.showFeedback {
				m.textarea.Focus()
				m.showFeedback = false
			}
			return m, nil
		case "enter":
			query := m.textarea.Value()
			if m.textarea.Focused() && query != "" {
				m.showFeedback = false
				m.isLoading = true
				cmds = append(cmds, searchNotesCmd(m, query))
				cmds = append(cmds, m.spinner.Init())
			}
		}
	}

	if m.isLoading {
		newSpinner, cmd := m.spinner.Update(msg)
		m.spinner = newSpinner.(spinner.Model)
		cmds = append(cmds, cmd)
	}
	newTextarea, cmd := m.textarea.Update(msg)
	m.textarea = newTextarea.(textarea.Model)
	cmds = append(cmds, cmd)
	if m.showFeedback && !m.textarea.Focused() {
		m.list, cmd = m.list.Update(msg)
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

var (
	docStyles = lipgloss.NewStyle().Margin(1, 1).Render
)

func (m Model) View() tea.View {
	var builder strings.Builder
	if m.showingNote {
		return tea.NewView(docStyles(utils.FormatSingleNoteForTUI(m.activeNote)))
	}
	if m.isLoading {
		builder.WriteString(m.spinner.View().Content)
	} else {
		builder.WriteString(m.textarea.View().Content)
	}
	if m.showFeedback {
		builder.WriteString("\n\n" + m.list.View())
	}
	return tea.NewView(docStyles(builder.String()))
}
