package notifications

import (
	"database/sql"
	_ "embed"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	_ "github.com/tursodatabase/turso-go"
)

type NotificationStore struct {
	activeTimers map[int]*time.Timer
	timersMutex  sync.Mutex
}

func StartNotificationService(db *sql.DB) {
	store := &NotificationStore{
		activeTimers: make(map[int]*time.Timer),
	}

	loadAndSchedulePendingNotes(db, store)

	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		loadAndSchedulePendingNotes(db, store)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
}
