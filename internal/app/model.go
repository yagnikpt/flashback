package app

type Flashback struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Type      string `json:"type"`
	CreatedAt string `json:"created_at"`
}

type FlashbackWithMetadata struct {
	Flashback
	Metadata map[string]string `json:"metadata"`
}

type Embedding struct {
	FlashbackID string
	Vector      []float64
}

type NoteMetadata struct {
	Tags string `json:"tags"`
	Tldr string `json:"tldr"`
}

type WebMetadata struct {
	ID          int
	FlashbackID string
	Key         string
	Value       string
	Source      string
	CreatedAt   string
}
