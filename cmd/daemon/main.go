package main

import (
	"database/sql"
	_ "embed"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"github.com/gen2brain/beeep"
	_ "github.com/tursodatabase/go-libsql"
	"github.com/yagnik-patel-47/flashback/internal/migration"
	"github.com/yagnik-patel-47/flashback/internal/utils"
)

var activeTimers = make(map[int]*time.Timer)
var timersMutex sync.Mutex

func main() {
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
		loadAndSchedulePendingNotes(db)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
}

func loadAndSchedulePendingNotes(db *sql.DB) {
	query := `SELECT id, title, time FROM notifications WHERE datetime(time) >= datetime('now') AND datetime(time) < datetime('now', '+2 days')`

	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	type Notification struct {
		ID    int       `json:"id"`
		Title string    `json:"title"`
		Time  time.Time `json:"time"`
	}

	var notifications []Notification

	for rows.Next() {
		notification := Notification{}
		err = rows.Scan(&notification.ID, &notification.Title, &notification.Time)
		if err != nil {
			log.Println(err)
		}
		notifications = append(notifications, notification)
	}
	if err = rows.Err(); err != nil {
		log.Println(err)
	}

	timersMutex.Lock()
	defer timersMutex.Unlock()

	for id := range activeTimers {
		found := false
		for _, notif := range notifications {
			if id == notif.ID {
				found = true
				break
			}
		}
		if !found {
			if timer, exists := activeTimers[id]; exists && timer != nil {
				timer.Stop()
				delete(activeTimers, id)
			}
		}
	}

	for _, notification := range notifications {
		// Skip if a timer is already active for this notification
		if _, exists := activeTimers[notification.ID]; exists {
			continue
		}

		timeDiff := time.Until(notification.Time)
		if timeDiff > 0 {
			timer := time.AfterFunc(timeDiff, func() {
				sendNotification("Reminder", notification.Title)
				// Remove from active timers after it fires
				timersMutex.Lock()
				delete(activeTimers, notification.ID)
				timersMutex.Unlock()
			})
			activeTimers[notification.ID] = timer
		}
	}
}

//go:embed icon.png
var icon []byte

func sendNotification(title, message string) {
	beeep.AppName = "Flashback"
	err := beeep.Notify(title, message, icon)
	if err != nil {
		log.Println("Error sending notification:", err)
	}
}
