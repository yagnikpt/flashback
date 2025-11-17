package notelist

import (
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dustin/go-humanize"
	"github.com/yagnikpt/flashback/internal/app"
	"github.com/yagnikpt/flashback/internal/models"
	"github.com/yagnikpt/flashback/internal/utils"
)

type Model struct {
	app         *app.App
	list        list.Model
	showingNote bool
	activeNote  models.FlashbackWithMetadata
}

func (m *Model) ResetView() {
	m.activeNote = models.FlashbackWithMetadata{}
	m.showingNote = false
}

type item struct {
	id, title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

func NewModel(app *app.App) Model {
	items := make([]list.Item, 0)
	d := newDelegate(newDelegateKeyMap())
	l := list.New(items, d, 0, 0)
	l.SetShowTitle(false)

	return Model{
		app:         app,
		list:        l,
		showingNote: false,
		activeNote:  models.FlashbackWithMetadata{},
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(getAllNotesCmd(m), getDimensionsCmd())
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case getAllNotesMsg:
		notes := []models.FlashbackWithMetadata(msg)
		items := make([]list.Item, len(notes))
		for i := range items {
			t, _ := time.Parse(time.RFC3339, notes[i].CreatedAt)
			items[i] = item{
				id:    notes[i].ID,
				title: notes[i].Content,
				desc:  humanize.Time(t),
			}
		}
		m.list.SetItems(items)

	case chosenNoteMsg:
		note := models.FlashbackWithMetadata(msg)
		m.activeNote = note
		m.showingNote = true

	// case deleteNoteMsg:

	case relayChooseMsg:
		return m, chooseNoteCmd(m, string(msg))
	case relayDeleteMsg:
		return m, deleteNoteCmd(m, string(msg))

	case dimensionsMsg:
		dims := dimensionsMsg(msg)
		m.list.SetSize(dims.width, dims.height-3)

	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height-3)

	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.showingNote {
				m.showingNote = false
				m.activeNote = models.FlashbackWithMetadata{}
				return m, nil
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)

	return m, cmd
}

var docStyles = lipgloss.NewStyle().Margin(0, 2).Render

func (m Model) View() string {
	if m.showingNote {
		return docStyles(utils.FormatSingleNote(m.activeNote))
	}
	return m.list.View()
}
