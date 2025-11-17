package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/yagnikpt/flashback/internal/app"
	"github.com/yagnikpt/flashback/internal/utils"
)

func NewListCmd(app *app.App) *cobra.Command {
	listCmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List all stored notes",
		Long: `List all notes stored in the flashback database.
This command displays a compact list of all notes with their IDs and Content.`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			flashbacks, err := app.GetAllNotes(ctx)
			if err != nil {
				fmt.Println("Error retrieving notes:", err)
				return
			}
			output := utils.FormatMultipleNotesCompact(flashbacks)
			fmt.Println(output)
		},
	}

	return listCmd
}
