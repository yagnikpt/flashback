/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/yagnikpt/flashback/internal/app"
	"github.com/yagnikpt/flashback/internal/utils"
)

func NewSearchCmd(app *app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "search",
		Aliases: []string{"s"},
		Short:   "Search notes using semantic similarity",
		Long: `Search for notes in the flashback database using semantic similarity. Provide a query string, and the tool will find notes with similar meanings based on embeddings.

Usage:
  flashback search [query]

Examples:
  flashback search "machine learning concepts"
  flashback search "buy groceries"`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println(cmd.Long)
				return
			}

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			words := strings.Join(args, " ")
			embeddings, _ := app.GenerateEmbeddingForNote(ctx, words, "RETRIEVAL_QUERY")
			flashbacks, err := app.RetrieveNotesBySimilarity(ctx, embeddings)
			if err != nil {
				fmt.Println("Error retrieving notes:", err)
			}
			output := utils.FormatMultipleNotesCompact(flashbacks)
			fmt.Println(output)
		},
	}

	return cmd
}
