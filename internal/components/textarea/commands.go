package textarea

import (
	"os"

	tea "charm.land/bubbletea/v2"
	"golang.org/x/term"
)

type setWidthMsg int

func setWidthCmd() tea.Cmd {
	return func() tea.Msg {
		width, _, err := term.GetSize(int(os.Stdout.Fd()))
		if err != nil {
			panic(err)
		}
		return setWidthMsg(width)
	}
}
