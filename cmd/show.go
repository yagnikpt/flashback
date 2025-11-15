package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yagnikpt/flashback/internal/app"
	"github.com/yagnikpt/flashback/internal/utils"
)

func NewShowCmd(app *app.App) *cobra.Command {
	showCmd := &cobra.Command{
		Use:     "show",
		Aliases: []string{"sh"},
		Short:   "Display details of a specific note",
		Long: `This command shows the content, metadata, and other details of the specified note.

Examples:
  flashback show 3C5uPKK4yvGZ3qUMJoCcdv`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println("Please provide the ID of the note to show.")
			}
			noteID := args[0]
			flashback, err := app.GetNoteByID(noteID)
			if err != nil {
				fmt.Println("Error retrieving note:", err)
				return
			}
			output := utils.FormatSingleNote(flashback)
			fmt.Println(output)
		},
	}

	return showCmd
}
