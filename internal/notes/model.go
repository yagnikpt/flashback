package notes

import "time"

type Note struct {
	ID        int       `json:"id"`
	Content   string    `json:"content"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
}
