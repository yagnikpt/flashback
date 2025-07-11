package textarea

import (
	"github.com/charmbracelet/bubbles/v2/key"
	"github.com/charmbracelet/bubbles/v2/textarea"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
)

const (
	defaultHeight    = 5
	defaultWidth     = 40
	defaultCharLimit = 0
)

type (
	errMsg error
)

type Model struct {
	textarea   textarea.Model
	err        errMsg
	OutputChan chan string
}

func (m Model) Init() tea.Cmd {
	return textarea.Blink
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.textarea.SetWidth(msg.Width - 4)
	case tea.KeyPressMsg:
		switch msg.String() {
		case "enter":
			if m.textarea.Value() != "" {
				m.OutputChan <- m.textarea.Value()
				m.textarea.SetValue("")
			}
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	var cmd tea.Cmd
	m.textarea, cmd = m.textarea.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return m.textarea.View()
}

func (m Model) Value() string {
	return m.textarea.Value()
}

func (m *Model) Blur() {
	m.textarea.Blur()
}

func (m *Model) Focus() {
	m.textarea.Focus()
}

func (m *Model) SetPlaceholder(str string) {
	m.textarea.Placeholder = str
}

func (m *Model) SetHeight(height int) {
	m.textarea.SetHeight(height)
}

func NewModel() Model {
	model := textarea.New()
	model.Placeholder = "Enter the note..."

	model.Focus()

	model.Prompt = "┃ "
	model.CharLimit = defaultCharLimit

	model.SetWidth(defaultWidth)
	model.SetHeight(defaultHeight)
	model.KeyMap.InsertNewline = key.NewBinding(
		key.WithKeys("shift+enter"),
		key.WithHelp("shift+enter", "insert newline"),
	)

	// Remove cursor line styling
	model.Styles.Focused.CursorLine = lipgloss.NewStyle()

	model.ShowLineNumbers = false

	return Model{
		textarea:   model,
		err:        nil,
		OutputChan: make(chan string),
	}
}
