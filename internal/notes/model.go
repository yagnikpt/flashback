package notes

type Note struct {
	ID        int    `json:"id"`
	Content   string `json:"content"`
	Title     string `json:"title"`
	CreatedAt string `json:"created_at"`
}
