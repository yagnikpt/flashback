package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yagnikpt/flashback/internal/app"
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
			tags := strings.Split(cmd.Flag("tags").Value.String(), ",")
			words := strings.Join(args, " ")
			title := cmd.Flag("title").Value.String()
			fmt.Println("add called", tags, title, words)
		},
	}

	cmd.Flags().StringP("tags", "t", "", "Comma separated tags for the record")
	cmd.Flags().StringP("title", "", "", "Title for the record")

	return cmd
}
