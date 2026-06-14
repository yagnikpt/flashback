package spinner

import (
	tea "charm.land/bubbletea/v2"
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
