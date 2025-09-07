package global

import (
	"database/sql"
	"sync"

	"github.com/yagnikpt/flashback/internal/config"
	"github.com/yagnikpt/flashback/internal/notes"
)

type Store struct {
	Mode         ModeState
	Width        int
	Height       int
	Loading      bool
	ShowFeedback bool
	Notes        *notes.Store
	Config       config.Config
}

var instance *Store
var once sync.Once

func InitStore(db *sql.DB, config config.Config) *Store {
	once.Do(func() {
		instance = &Store{Mode: StateNote, Width: 0, Height: 0, Loading: false, ShowFeedback: false, Notes: notes.NewStore(db, config.APIKey), Config: config}
	})
	return instance
}

func GetStore() *Store {
	return instance
}
