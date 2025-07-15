package notes

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/yagnikpt/flashback/internal/contentloaders"
	"github.com/yagnikpt/flashback/internal/utils"
	"google.golang.org/genai"
)

type Store struct {
	db         *sql.DB
	genai      *genai.Client
	StatusChan chan string
}

func NewStore(db *sql.DB, apiKey string) *Store {
	client, err := genai.NewClient(context.Background(), &genai.ClientConfig{
		APIKey: apiKey,
	})
	if err != nil {
		log.Fatal(err)
	}
	statusChan := make(chan string)

	return &Store{
		db:         db,
		genai:      client,
		StatusChan: statusChan,
	}
}

func (s *Store) CreateNote(content string) error {
	ctx := context.Background()

	query := `INSERT INTO notes DEFAULT VALUES RETURNING id`
	var noteID int
	err := s.db.QueryRow(query).Scan(&noteID)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	wg.Add(2)

	errChan := make(chan error, 2)

	go func() {
		defer wg.Done()

		if err != nil {
			errChan <- fmt.Errorf("error: %w", err)
			return
		}
		config := &genai.GenerateContentConfig{
			ResponseMIMEType: "application/json",
			ResponseSchema: &genai.Schema{
				Type: genai.TypeArray,
				Items: &genai.Schema{
					Type: genai.TypeObject,
					Properties: map[string]*genai.Schema{
						"title": {Type: genai.TypeString, Description: "The title to show in notification"},
						"time":  {Type: genai.TypeString, Description: "The time mentioned for notification in format YYYY-MM-DD HH:MM AM/PM"},
					},
					PropertyOrdering: []string{"title", "time"},
				},
			},
		}

		input := "Identify the time given in the provided note for using the time in notification. Current date and time is " + time.Now().Format("2006-01-02 15:04 PM.") + "\n\n" + "Note: " + content
		result, err := s.genai.Models.GenerateContent(
			ctx,
			"gemini-2.5-flash",
			genai.Text(input),
			config,
		)

		if err != nil {
			errChan <- fmt.Errorf("error: %w", err)
			return
		}

		type Notification struct {
			Title string `json:"title"`
			Time  string `json:"time"`
		}
		var notifications []Notification

		err = json.Unmarshal([]byte(result.Text()), &notifications)
		if err != nil {
			errChan <- fmt.Errorf("notification parsing error: %w", err)
			return
		}

		if len(notifications) > 0 {
			s.StatusChan <- "Setting notifications..."
		}

		for _, notification := range notifications {
			parsedTime, err := time.ParseInLocation("2006-01-02 15:04 PM", notification.Time, time.Local)
			if err != nil {
				log.Println("Error parsing time:", err)
				continue
			}
			_, err = s.db.Exec("INSERT INTO notifications (title, time, note_id) VALUES (?, ?, ?)", notification.Title, parsedTime.In(time.UTC), noteID)
			if err != nil {
				log.Println("Error inserting notification:", err)
				continue
			}
		}
	}()

	go func() {
		defer wg.Done()

		var chunks []string

		if strings.Contains(content, "#clipboard") {
			clipboardChunks, err := contentloaders.GetClipboardContent()
			if err != nil {
				errChan <- err
				return
			}
			withClipboardContent := strings.ReplaceAll(content, "#clipboard", strings.Join(clipboardChunks, "\n"))
			newChunks, err := utils.SplitText(withClipboardContent)
			if err != nil {
				errChan <- err
				return
			}
			chunks = append(chunks, newChunks...)
		} else {
			newChunks, err := utils.SplitText(content)
			if err != nil {
				errChan <- err
				return
			}
			chunks = append(chunks, newChunks...)
		}

		searchTerms := utils.ExtractSearchTerms(content)
		for _, url := range searchTerms.Web {
			s.StatusChan <- fmt.Sprintf("Fetching content from %s", url)
			webChunks, err := contentloaders.GetWebpageContent(url)
			if err != nil {
				log.Println("Error fetching web content:", err)
				continue
			}
			chunks = append(chunks, webChunks...)
		}

		for _, file := range searchTerms.Files {
			s.StatusChan <- fmt.Sprintf("Fetching content from %s", file)
			fileChunks, err := contentloaders.GetTextContent(file)
			log.Println(fileChunks)
			if err != nil {
				log.Println("Error fetching file content:", err)
				continue
			}
			chunks = append(chunks, fileChunks...)
		}

		for index, chunk := range chunks {
			query := `INSERT INTO chunks (note_id, content, chunk_number) VALUES (?, ?, ?) RETURNING id`
			var chunkID int
			err := s.db.QueryRow(query, noteID, chunk, index+1).Scan(&chunkID)
			if err != nil {
				errChan <- err
				continue
			}

			contents := []*genai.Content{
				genai.NewContentFromText(chunk, genai.RoleUser),
			}
			result, err := s.genai.Models.EmbedContent(ctx,
				"gemini-embedding-exp-03-07",
				contents,
				&genai.EmbedContentConfig{
					TaskType:             "RETRIEVAL_DOCUMENT",
					OutputDimensionality: genai.Ptr[int32](768),
				},
			)
			if err != nil {
				errChan <- fmt.Errorf("embedding generation error: %w", err)
				continue
			}

			embeddings, err := json.Marshal(result.Embeddings[0].Values)
			if err != nil {
				errChan <- fmt.Errorf("embedding marshaling error: %w", err)
				continue
			}

			_, err = s.db.Exec("INSERT INTO embeddings (chunk_id, vector) VALUES (?, vector32(?))", chunkID, string(embeddings))
			if err != nil {
				errChan <- fmt.Errorf("embedding storage error: %w", err)
				continue
			}
		}
	}()

	wg.Wait()

	close(errChan)

	var firstErr error
	for err := range errChan {
		if firstErr == nil {
			firstErr = err
		}
	}

	return firstErr
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
			OutputDimensionality: genai.Ptr[int32](768),
		},
	)
	if err != nil {
		return "", err
	}

	embeddings, err := json.Marshal(result.Embeddings[0].Values)
	if err != nil {
		return "", err
	}

	query := `SELECT n.id, n.created_at, c.id AS chunk_id, c.content, c.chunk_number
		FROM notes n
		INNER JOIN chunks c ON n.id = c.note_id
		INNER JOIN embeddings e ON c.id = e.chunk_id
		WHERE vector_distance_cos(e.vector, vector32(?)) > 0.3
		ORDER BY vector_distance_cos(e.vector, vector32(?)) ASC`

	rows, err := s.db.Query(query, string(embeddings), string(embeddings))
	if err != nil {
		return "", err
	}
	defer rows.Close()

	var fetchedChunks []CombinedNote

	for rows.Next() {
		note := CombinedNote{}
		err = rows.Scan(&note.ID, &note.CreatedAt, &note.ChunkID, &note.Content, &note.ChunkNumber)
		if err != nil {
			return "", err
		}
		fetchedChunks = append(fetchedChunks, note)
	}
	if err = rows.Err(); err != nil {
		return "", err
	}

	// log.Println(len(fetchedChunks), "records populated")

	// for _, chunk := range fetchedChunks {
	// 	log.Println(chunk.Content)
	// }

	chunksContext := strings.Builder{}
	chunksContext.WriteString("\n\nChunks:\n")
	loc, _ := time.LoadLocation("Local")
	for _, chunk := range fetchedChunks {
		chunksContext.WriteString(fmt.Sprintf("- %s - timestamp: %s\n", chunk.Content, chunk.CreatedAt.In(loc).Format("2006-01-02 15:04")))
	}

	finalInput := "User query: " + userQuery + chunksContext.String()

	prompt, err := os.ReadFile("internal/notes/retrieve_prompt.txt")
	if err != nil {
		return "", err
	}
	config := &genai.GenerateContentConfig{
		SystemInstruction: genai.NewContentFromText(string(prompt), genai.RoleUser),
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

	return response.Text(), nil
	// return "", nil
}

func (s *Store) GetAllNotes() (notes []CombinedNote, e error) {
	rows, err := s.db.Query(`SELECT n.id, n.created_at, c.content, c.id AS chunk_id, c.chunk_number FROM notes n INNER JOIN chunks c ON n.id = c.note_id WHERE c.chunk_number = 1 ORDER BY n.created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		note := CombinedNote{}
		err = rows.Scan(&note.ID, &note.CreatedAt, &note.Content, &note.ChunkID, &note.ChunkNumber)
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
