package help

import "github.com/charmbracelet/bubbles/v2/key"

type createKeyMap struct {
	Help       key.Binding
	Quit       key.Binding
	Create     key.Binding
	ChangeMode key.Binding
	NewLine    key.Binding
}

type recallKeyMap struct {
	Help       key.Binding
	Quit       key.Binding
	Recall     key.Binding
	ChangeMode key.Binding
	NewLine    key.Binding
}

type deleteKeyMap struct {
	Up         key.Binding
	Down       key.Binding
	Left       key.Binding
	Right      key.Binding
	Help       key.Binding
	Quit       key.Binding
	Delete     key.Binding
	ChangeMode key.Binding
}

func (k createKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{}
}

func (k createKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.ChangeMode, k.Create, k.NewLine},
		{k.Help, k.Quit},
	}
}

func (k recallKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{}
}

func (k recallKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.ChangeMode, k.Recall, k.NewLine},
		{k.Help, k.Quit},
	}
}

func (k deleteKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{}
}

func (k deleteKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.ChangeMode, k.Delete},
		{k.Left, k.Down, k.Up, k.Right},
		{k.Help, k.Quit},
	}
}

var createKeys = createKeyMap{
	Help: key.NewBinding(
		key.WithKeys("alt+?"),
		key.WithHelp("alt+?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("esc", "ctrl+c"),
		key.WithHelp("esc,ctrl+c", "quit"),
	),
	Create: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "create note"),
	),
	ChangeMode: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "change mode"),
	),
	NewLine: key.NewBinding(
		key.WithKeys("shift+enter"),
		key.WithHelp("shift+enter", "new line"),
	),
}

var recallKeys = recallKeyMap{
	Help: key.NewBinding(
		key.WithKeys("alt+?"),
		key.WithHelp("alt+?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("esc", "ctrl+c"),
		key.WithHelp("esc,ctrl+c", "quit"),
	),
	Recall: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "recall note"),
	),
	ChangeMode: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "change mode"),
	),
	NewLine: key.NewBinding(
		key.WithKeys("shift+enter"),
		key.WithHelp("shift+enter", "new line"),
	),
}

var deleteKeys = deleteKeyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "move left"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "move right"),
	),
	Help: key.NewBinding(
		key.WithKeys("alt+?"),
		key.WithHelp("alt+?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("esc", "ctrl+c"),
		key.WithHelp("esc,ctrl+c", "quit"),
	),
	Delete: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "delete note"),
	),
	ChangeMode: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "change mode"),
	),
}
