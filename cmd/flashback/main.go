package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	_ "github.com/tursodatabase/go-libsql"

	tea "github.com/charmbracelet/bubbletea/v2"

	"github.com/yagnik-patel-47/flashback/internal/app"
	"github.com/yagnik-patel-47/flashback/internal/migration"
)

func main() {
	dataDir, err := getLocalDataDir()
	// fmt.Println("Data directory:", dataDir)
	// os.Exit(0)
	if err != nil {
		fmt.Println("Error getting local data directory:", err)
		os.Exit(1)
	}

	f, err := tea.LogToFile(filepath.Join(dataDir, "debug.log"), "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()

	db, err := sql.Open("libsql", "file:"+filepath.Join(dataDir, "flashback.db"))
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

func getLocalDataDir() (string, error) {
	appName := "flashback"

	var dataDir string
	var err error

	switch runtime.GOOS {
	case "darwin":
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Println("Error getting user home directory:", err)
			return "", err
		}
		dataDir = filepath.Join(homeDir, "Library", "Application Support")
	case "linux":
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Println("Error getting user home directory:", err)
			return "", err
		}
		dataDir = filepath.Join(homeDir, ".local", "share")

	case "windows":
		dataDir = os.Getenv("LocalAppData")
		if dataDir == "" {
			return "", fmt.Errorf("error: %%LocalAppData%% environment variable not found")
		}
	default:
		return "", fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	appDataDir := filepath.Join(dataDir, appName)

	err = os.MkdirAll(appDataDir, 0755)
	if err != nil {
		return "", fmt.Errorf("error creating app data directory: %w", err)
	}

	return appDataDir, nil
}
