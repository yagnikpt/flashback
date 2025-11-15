/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yagnikpt/flashback/internal/app"
)

func NewRemoveCmd(app *app.App) *cobra.Command {
	removeCmd := &cobra.Command{
		Use:     "remove",
		Aliases: []string{"rm", "delete", "del"},
		Short:   "Remove a note by its ID",
		Long: `Remove a note from the flashback database by providing its unique ID.

Examples:
  flashback remove 12345`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println("Please provide the ID of the note to remove.")
			}
			noteID := args[0]
			err := app.DeleteNoteByID(noteID)
			if err != nil {
				fmt.Println("Error removing note:", err)
				return
			}
			fmt.Println("Note removed successfully.")
		},
	}

	return removeCmd
}
