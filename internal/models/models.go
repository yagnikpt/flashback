package models

type Flashback struct {
	ID        string `json:"id"`
	Content   string `json:"content"`
	Type      string `json:"type"`
	CreatedAt string `json:"created_at"`
}

type FlashbackWithMetadata struct {
	Flashback
	Metadata map[string]string `json:"metadata"`
}
