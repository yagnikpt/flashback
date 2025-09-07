package help

import (
	"github.com/charmbracelet/bubbles/v2/help"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/yagnikpt/flashback/internal/global"
)

type Model struct {
	createKeys createKeyMap
	recallKeys recallKeyMap
	deleteKeys deleteKeyMap
	createHelp help.Model
	deleteHelp help.Model
	recallHelp help.Model
	store      *global.Store
}

func NewModel() Model {
	store := global.GetStore()
	createHelp := help.New()
	deleteHelp := help.New()
	recallHelp := help.New()

	createHelp.ShowAll = true
	deleteHelp.ShowAll = true
	recallHelp.ShowAll = true

	return Model{
		createKeys: createKeys,
		recallKeys: recallKeys,
		deleteKeys: deleteKeys,
		createHelp: createHelp,
		deleteHelp: deleteHelp,
		recallHelp: recallHelp,
		store:      store,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.createHelp.Width = msg.Width
		m.deleteHelp.Width = msg.Width
		m.recallHelp.Width = msg.Width
	}

	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch m.store.Mode {
	case global.StateNote:
		m.createHelp, cmd = m.createHelp.Update(msg)
		cmds = append(cmds, cmd)
	case global.StateDelete:
		m.deleteHelp, cmd = m.deleteHelp.Update(msg)
		cmds = append(cmds, cmd)
	case global.StateRecall:
		m.recallHelp, cmd = m.recallHelp.Update(msg)
		cmds = append(cmds, cmd)
	default:
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	switch m.store.Mode {
	case global.StateNote:
		return m.createHelp.View(m.createKeys)
	case global.StateDelete:
		return m.deleteHelp.View(m.deleteKeys)
	case global.StateRecall:
		return m.recallHelp.View(m.recallKeys)
	default:
		return m.createHelp.View(m.createKeys)
	}
}
