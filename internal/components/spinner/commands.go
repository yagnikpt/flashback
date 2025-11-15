package spinner

import (
	tea "github.com/charmbracelet/bubbletea"
)

type statusMsg string
type closedMsg struct{}

func getCurrentStatus(m Model) tea.Cmd {
	return func() tea.Msg {
		status, ok := <-m.status
		if !ok {
			return closedMsg{}
		}
		return statusMsg(status)
	}
}
