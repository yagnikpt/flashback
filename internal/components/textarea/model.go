package textarea

import (
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/textarea"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
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

func (m Model) View() tea.View {
	return tea.NewView(m.textarea.View())
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
	model.KeyMap.InsertNewline = key.NewBinding(
		key.WithKeys("shift+enter"),
		key.WithHelp("shift+enter", "newline"),
	)
	model.Placeholder = "Enter the note..."

	model.Focus()

	model.Prompt = "┃ "
	model.CharLimit = defaultCharLimit

	styles := model.Styles()
	styles.Focused.Prompt = lipgloss.NewStyle().Foreground(lipgloss.Color("4"))
	styles.Blurred.Prompt = lipgloss.NewStyle().Foreground(lipgloss.Color("#525252"))
	styles.Focused.CursorLine = lipgloss.NewStyle()
	model.SetStyles(styles)

	model.SetHeight(defaultHeight)

	model.ShowLineNumbers = false

	return Model{
		textarea:   model,
		err:        nil,
		OutputChan: make(chan string),
	}
}
