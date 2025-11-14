/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yagnikpt/flashback/internal/app"
)

func NewSearchCmd(app *app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "search",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
 and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println(cmd.Long)
				return
			}
			words := strings.Join(args, " ")
			embeddings, _ := app.GenerateEmbeddingForNote(words, "RETRIEVAL_QUERY")
			flashbacks, err := app.RetrieveNotesBySimilarity(embeddings)
			if err != nil {
				fmt.Println("Error retrieving notes:", err)
			}
			fmt.Println(flashbacks)
		},
	}

	return cmd
}
