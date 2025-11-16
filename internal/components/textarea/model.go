package textarea

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	defaultHeight    = 5
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
	return tea.Batch(textarea.Blink, setWidthCmd())
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case setWidthMsg:
		width := int(msg)
		m.textarea.SetWidth(width - 2)

	case tea.WindowSizeMsg:
		m.textarea.SetWidth(msg.Width - 2)

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

func (m *Model) SetValue(value string) {
	m.textarea.SetValue(value)
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

func (m Model) Focused() bool {
	return m.textarea.Focused()
}

func NewModel() Model {
	model := textarea.New()
	model.Placeholder = "Enter the note..."

	model.Focus()

	model.Prompt = "┃ "
	model.CharLimit = defaultCharLimit
	model.FocusedStyle.Prompt = lipgloss.NewStyle().Foreground(lipgloss.Color("#94a3b8"))
	model.BlurredStyle.Prompt = lipgloss.NewStyle().Foreground(lipgloss.Color("#525252"))

	model.SetHeight(defaultHeight)
	model.KeyMap.InsertNewline = key.NewBinding(
		key.WithKeys("ctrl+enter"),
		key.WithHelp("ctrl+enter", "insert newline"),
	)

	model.ShowLineNumbers = false

	return Model{
		textarea:   model,
		err:        nil,
		OutputChan: make(chan string),
	}
}
