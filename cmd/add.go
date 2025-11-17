package cmd

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/yagnikpt/flashback/internal/app"
	"github.com/yagnikpt/flashback/internal/components/spinner"
	"github.com/yagnikpt/flashback/internal/contentloaders"
)

func NewAddCmd(app *app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "add",
		Aliases: []string{"a"},
		Short:   "Add a new note from text or URL",
		Long: `Add a new note to the flashback database. Provide text directly or a URL to fetch and store webpage content. The tool automatically generates metadata and embeddings for semantic search.

Examples:
  flashback add Remember to buy groceries
  flashback add https://example.com/useful-article
`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println(cmd.Long)
				return
			}

			statusChan := make(chan string)
			errorChan := make(chan error)

			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			defer cancel()

			go func() {
				words := strings.Join(args, " ")
				var metadata map[string]string
				var noteType string
				if len(args) == 1 && strings.HasPrefix(words, "http") {
					statusChan <- "Fetching webpage content..."
					pageContent, err := contentloaders.GetWebPage(ctx, words)
					if err != nil {
						errorChan <- err
						return
					}
					pageContentWithUrl := fmt.Sprintf("URL: %s\n\n%s", words, pageContent)
					statusChan <- "Generating metadata for webpage..."
					metadata, err = app.GenerateMetadataForWebNote(ctx, pageContentWithUrl)
					if err != nil {
						errorChan <- err
						return
					}
					noteType = "url"
				} else {
					statusChan <- "Generating metadata for note..."
					_metadata, err := app.GenerateMetadataForSimpleNote(ctx, words)
					if err != nil {
						errorChan <- err
						return
					}
					metadata = _metadata
					noteType = "text"
				}
				var finalContent strings.Builder
				userInput := fmt.Sprintf("USER INPUT:\n content: %s\n", words)
				finalContent.WriteString(userInput)
				finalContent.WriteString("METADATA:\n")
				for key, value := range metadata {
					metadataLine := fmt.Sprintf(" %s: %s\n", key, value)
					finalContent.WriteString(metadataLine)
				}
				statusChan <- "Saving the note..."
				embeddings, err := app.GenerateEmbeddingForNote(ctx, finalContent.String(), "RETRIEVAL_DOCUMENT")
				if err != nil {
					errorChan <- err
					return
				}
				err = app.InsertNote(ctx, words, noteType, metadata, embeddings)
				if err != nil {
					errorChan <- err
				} else {
					close(statusChan)
					close(errorChan)
				}
			}()

			go func() {
				err, ok := <-errorChan
				if !ok {
					fmt.Print("\033[A\033[2K")
					fmt.Println("Note added successfully!")
					return
				}
				fmt.Printf("Error: %v\n", err)
				close(statusChan)
				close(errorChan)
			}()

			spinner.Run(statusChan, false)
		},
	}

	cmd.Flags().StringP("tags", "t", "", "Comma separated tags for the record")

	return cmd
}
