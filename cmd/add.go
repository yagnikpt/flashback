package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yagnikpt/flashback/internal/app"
	"github.com/yagnikpt/flashback/internal/contentloaders"
)

func NewAddCmd(app *app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "A brief description of your application",
		Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println(cmd.Long)
				return
			}
			// tags := strings.Split(cmd.Flag("tags").Value.String(), ",")

			words := strings.Join(args, " ")
			var metadata map[string]string
			if len(args) == 1 && strings.HasPrefix(words, "http") {
				pageContent, err := contentloaders.GetWebPage(words)
				if err != nil {
					fmt.Println("Error fetching webpage:", err)
					return
				}
				pageContentWithUrl := fmt.Sprintf("URL: %s\n\n%s", words, pageContent)
				metadata, err = app.GenerateMetadataForWebNote(pageContentWithUrl)
				if err != nil {
					fmt.Println("Error generating metadata:", err)
					return
				}
			} else {
				_metadata, err := app.GenerateMetadataForSimpleNote(words)
				if err != nil {
					fmt.Println("Error generating metadata:", err)
					return
				}
				metadata = _metadata
			}
			var finalContent strings.Builder
			userInput := fmt.Sprintf("USER INPUT:\n content: %s\n", words)
			finalContent.WriteString(userInput)
			finalContent.WriteString("METADATA:\n")
			for key, value := range metadata {
				metadataLine := fmt.Sprintf(" %s: %s\n", key, value)
				finalContent.WriteString(metadataLine)
			}
			embeddings, _ := app.GenerateEmbeddingForNote(finalContent.String(), "RETRIEVAL_DOCUMENT")
			err := app.InsertNote(words, "text", metadata, embeddings)
			if err != nil {
				fmt.Println("Error inserting note:", err)
			} else {
				fmt.Println("Note added successfully!")
			}
		},
	}

	cmd.Flags().StringP("tags", "t", "", "Comma separated tags for the record")
	// cmd.Flags().StringP("title", "", "", "Title for the record")

	return cmd
}
