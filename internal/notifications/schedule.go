package notifications

import (
	"database/sql"
	_ "embed"
	"log"
	"time"

	"github.com/gen2brain/beeep"
)

func loadAndSchedulePendingNotes(db *sql.DB, store *NotificationStore) {
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

	store.timersMutex.Lock()
	defer store.timersMutex.Unlock()

	for id := range store.activeTimers {
		found := false
		for _, notif := range notifications {
			if id == notif.ID {
				found = true
				break
			}
		}
		if !found {
			if timer, exists := store.activeTimers[id]; exists && timer != nil {
				timer.Stop()
				delete(store.activeTimers, id)
			}
		}
	}

	for _, notification := range notifications {
		if _, exists := store.activeTimers[notification.ID]; exists {
			continue
		}

		timeDiff := time.Until(notification.Time)
		if timeDiff > 0 && timeDiff <= 48*time.Hour {
			timer := time.AfterFunc(timeDiff, func() {
				sendNotification("Reminder", notification.Title)
				store.timersMutex.Lock()
				delete(store.activeTimers, notification.ID)
				store.timersMutex.Unlock()
			})
			store.activeTimers[notification.ID] = timer
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
