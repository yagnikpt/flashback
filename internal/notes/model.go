package notes

import "time"

type Note struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"created_at"`
}

type Chunk struct {
	ID          int    `json:"id"`
	NoteID      int    `json:"note_id"`
	Content     string `json:"content"`
	ChunkNumber int    `json:"chunk_number"`
}

type CombinedNote struct {
	ID          int       `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	Content     string    `json:"content"`
	ChunkID     int       `json:"chunk_id"`
	ChunkNumber int       `json:"chunk_number"`
}

type NoteWithChunks struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Chunks    []Chunk   `json:"chunks"`
}
