package main

import (
	"database/sql"
	_ "embed"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/gen2brain/beeep"
	_ "github.com/tursodatabase/go-libsql"
	"github.com/yagnik-patel-47/flashback/internal/migration"
	"github.com/yagnik-patel-47/flashback/internal/utils"
)

//go:embed icon.png
var icon []byte

func main() {
	fmt.Println("Flashback Daemon")

	dataDir, err := utils.GetLocalDataDir()
	if err != nil {
		os.Exit(1)
	}

	logFile, err := os.OpenFile(filepath.Join(dataDir, "debug.log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		os.Exit(1)
	}
	defer logFile.Close()

	log.SetOutput(logFile)

	db, err := sql.Open("libsql", "file:"+filepath.Join(dataDir, "flashback.db"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = migration.Migrate(db)
	if err != nil {
		log.Fatal(err)
	}

	loadAndSchedulePendingNotes(db)

	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		log.Println("Checking for new pending notifications...")
		loadAndSchedulePendingNotes(db)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
}

func loadAndSchedulePendingNotes(db *sql.DB) {
	query := `SELECT title, time FROM notifications WHERE datetime(time) >= datetime('now') AND datetime(time) < datetime('now', '+2 days')`

	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	type Notification struct {
		Title string    `json:"title"`
		Time  time.Time `json:"time"`
	}

	var notifications []Notification

	for rows.Next() {
		notification := Notification{}
		err = rows.Scan(&notification.Title, &notification.Time)
		if err != nil {
			log.Println(err)
		}
		notifications = append(notifications, notification)
	}
	if err = rows.Err(); err != nil {
		log.Println(err)
	}

	for _, notification := range notifications {
		timeDiff := time.Until(notification.Time)
		fmt.Println(notification.Title, timeDiff)
		time.AfterFunc(timeDiff, func() {
			sendNotification("Reminder", notification.Title)
		})
	}
}

func sendNotification(title, message string) {
	beeep.AppName = "Flashback"
	err := beeep.Notify(title, message, icon)
	if err != nil {
		fmt.Println("Error sending notification:", err)
	}
}
