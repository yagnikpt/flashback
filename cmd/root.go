package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yagnikpt/flashback/internal/app"
)

func NewRootCmd(app *app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "flashback",
		Short: "A CLI tool for managing personal notes with semantic search",
		Long: `Flashback is a command-line application for storing and retrieving personal notes using semantic search powered by embeddings.

It supports adding notes from text or web URLs, generating metadata, and performing similarity-based searches to help you recall information efficiently.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(cmd.Long)
		},
	}

	cmd.AddCommand(NewAddCmd(app))
	cmd.AddCommand(NewSearchCmd(app))
	cmd.AddCommand(NewListCmd(app))
	cmd.AddCommand(NewRemoveCmd(app))
	cmd.AddCommand(NewShowCmd(app))

	return cmd
}

func Execute(app *app.App) {
	rootCmd := NewRootCmd(app)

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
