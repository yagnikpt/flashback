package searchnotes

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/yagnikpt/flashback/internal/models"
)

func newDelegate(keys *delegateKeyMap) list.DefaultDelegate {
	d := list.NewDefaultDelegate()

	c := lipgloss.Color("#94a3b8")
	d.Styles.SelectedTitle = d.Styles.SelectedTitle.Foreground(c).Border(lipgloss.ThickBorder(), false, false, false, true).BorderLeftForeground(c)
	d.Styles.SelectedDesc = d.Styles.SelectedTitle

	d.UpdateFunc = func(msg tea.Msg, m *list.Model) tea.Cmd {
		var note models.FlashbackWithMetadata

		if i, ok := m.SelectedItem().(item); ok {
			note = i.full
		} else {
			return nil
		}

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, keys.choose):
				return relayChooseCmd(note)
			}
		}

		return nil
	}

	help := []key.Binding{keys.choose, keys.search}

	d.ShortHelpFunc = func() []key.Binding {
		return help
	}

	d.FullHelpFunc = func() [][]key.Binding {
		return [][]key.Binding{help}
	}

	return d
}

type delegateKeyMap struct {
	choose key.Binding
	search key.Binding
}

// Additional short help entries. This satisfies the help.KeyMap interface and
// is entirely optional.
func (d delegateKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		d.choose,
		d.search,
	}
}

// Additional full help entries. This satisfies the help.KeyMap interface and
// is entirely optional.
func (d delegateKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			d.choose,
			d.search,
		},
	}
}

func newDelegateKeyMap() *delegateKeyMap {
	return &delegateKeyMap{
		choose: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "choose"),
		),
		search: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "search"),
		),
	}
}
