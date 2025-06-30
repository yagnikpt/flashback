package notes

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"google.golang.org/genai"
)

type Store struct {
	db    *sql.DB
	genai *genai.Client
}

func NewStore(db *sql.DB, apiKey string) *Store {
	client, err := genai.NewClient(context.Background(), &genai.ClientConfig{
		APIKey: apiKey,
	})
	if err != nil {
		log.Fatal(err)
	}

	return &Store{
		db:    db,
		genai: client,
	}
}

func (s *Store) CreateNote(content string) error {
	ctx := context.Background()

	query := `INSERT INTO notes (title, content) VALUES (?, ?) RETURNING id`
	var noteID int
	err := s.db.QueryRow(query, "", content).Scan(&noteID)
	if err != nil {
		return err
	}

	contents := []*genai.Content{
		genai.NewContentFromText(content, genai.RoleUser),
	}
	result, err := s.genai.Models.EmbedContent(ctx,
		"gemini-embedding-exp-03-07",
		contents,
		&genai.EmbedContentConfig{
			TaskType:             "RETRIEVAL_DOCUMENT",
			OutputDimensionality: func(i int32) *int32 { return &i }(768),
		},
	)
	if err != nil {
		return err
	}

	embeddings, err := json.Marshal(result.Embeddings[0].Values)
	if err != nil {
		return err
	}

	_, err = s.db.Exec("INSERT INTO embeddings (note_id, vector) VALUES (?, vector32(?))", noteID, string(embeddings))
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) Recall(userQuery string) (string, error) {
	ctx := context.Background()
	contents := []*genai.Content{
		genai.NewContentFromText(userQuery, genai.RoleUser),
	}
	result, err := s.genai.Models.EmbedContent(ctx,
		"gemini-embedding-exp-03-07",
		contents,
		&genai.EmbedContentConfig{
			TaskType:             "RETRIEVAL_QUERY",
			OutputDimensionality: func(i int32) *int32 { return &i }(768),
		},
	)
	if err != nil {
		return "", err
	}

	embeddings, err := json.Marshal(result.Embeddings[0].Values)
	if err != nil {
		return "", err
	}

	query := `SELECT notes.id, notes.title, notes.content, notes.created_at FROM notes JOIN embeddings ON notes.id = embeddings.note_id ORDER BY
       vector_distance_cos(embeddings.vector, vector32(?))
    ASC LIMIT 3`

	rows, err := s.db.Query(query, string(embeddings), string(embeddings))
	if err != nil {
		return "", err
	}
	defer rows.Close()

	var fetchedNotes []Note

	for rows.Next() {
		note := Note{}
		err = rows.Scan(&note.ID, &note.Title, &note.Content, &note.CreatedAt)
		if err != nil {
			return "", err
		}
		fetchedNotes = append(fetchedNotes, note)
	}
	if err = rows.Err(); err != nil {
		return "", err
	}

	notesContext := strings.Builder{}
	notesContext.WriteString("\n\nUser notes:\n")
	loc, _ := time.LoadLocation("Local")
	for _, note := range fetchedNotes {
		notesContext.WriteString(fmt.Sprintf("- %s - timestamp: %s\n", note.Content, note.CreatedAt.In(loc).Format("2006-01-02 15:04")))
	}

	finalInput := "User query: " + userQuery + notesContext.String()
	// log.Println(finalInput)

	data, err := os.ReadFile("internal/notes/system_prompt.txt")
	if err != nil {
		return "", err
	}
	config := &genai.GenerateContentConfig{
		SystemInstruction: genai.NewContentFromText(string(data), genai.RoleUser),
	}

	response, err := s.genai.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash",
		genai.Text(finalInput),
		config,
	)

	if err != nil {
		return "", err
	}

	log.Println(response.Text())

	return response.Text(), nil
	// return "", nil
}

func (s *Store) GetAllNotes() (notes []Note, e error) {
	rows, err := s.db.Query(`SELECT id, title, content, created_at FROM notes ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		note := Note{}
		err = rows.Scan(&note.ID, &note.Title, &note.Content, &note.CreatedAt)
		if err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return notes, nil
}

func (s *Store) DeleteNote(id int) error {
	_, err := s.db.Exec("DELETE FROM notes WHERE id = ?", id)
	return err
}
