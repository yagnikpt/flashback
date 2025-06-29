package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/tursodatabase/go-libsql"

	tea "github.com/charmbracelet/bubbletea/v2"

	"github.com/yagnik-patel-47/flashback/internal/app"
	"github.com/yagnik-patel-47/flashback/internal/migration"
)

func main() {
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()

	db, err := sql.Open("libsql", "file:"+"test.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = migration.Migrate(db)
	if err != nil {
		log.Fatal(err)
	}

	p := tea.NewProgram(app.InitModel(db), tea.WithAltScreen(), tea.WithKeyboardEnhancements())
	_, err = p.Run()
	if err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
