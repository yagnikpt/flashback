package app

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/lithammer/shortuuid/v4"
	"github.com/yagnikpt/flashback/internal/models"
)

func (app *App) InsertNote(ctx context.Context, content, dataType string, metadata map[string]string, embeddings []float32) error {
	id := shortuuid.New()
	tx, err := app.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	insertQuery := `INSERT INTO flashbacks (id, content, type) VALUES (?, ?, ?)`
	_, err = tx.Exec(insertQuery, id, content, dataType)
	if err != nil {
		return err
	}

	insertMetadataQuery := `INSERT INTO metadata (flashback_id, key, value) VALUES (?, ?, ?)`
	for key, value := range metadata {
		_, err := tx.Exec(insertMetadataQuery, id, key, value)
		if err != nil {
			return err
		}
	}

	embeddingsData, err := json.Marshal(embeddings)
	if err != nil {
		return err
	}
	insertEmbeddingQuery := `INSERT INTO embeddings (flashback_id, vector) VALUES (?, vector32(?))`
	_, err = tx.Exec(insertEmbeddingQuery, id, string(embeddingsData))
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (app *App) RetrieveNotesBySimilarity(ctx context.Context, vector []float32) ([]models.FlashbackWithMetadata, error) {
	embeddings, err := json.Marshal(vector)
	if err != nil {
		return nil, err
	}

	query := `
    SELECT f.id, f.content, f.type, f.created_at, m.key, m.value
    FROM flashbacks f
    JOIN embeddings e ON f.id = e.flashback_id
    LEFT JOIN metadata m ON f.id = m.flashback_id
    WHERE vector_distance_cos(e.vector, vector32(?)) < 0.40
    ORDER BY vector_distance_cos(e.vector, vector32(?)) ASC LIMIT 20
    `

	rows, err := app.DB.QueryContext(ctx, query, string(embeddings), string(embeddings))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	flashbacks := []models.FlashbackWithMetadata{}
	idIndex := make(map[string]int)
	for rows.Next() {
		var id, content, dataType, createdAt string
		var key, value sql.NullString
		if err := rows.Scan(&id, &content, &dataType, &createdAt, &key, &value); err != nil {
			return nil, err
		}

		idx, exists := idIndex[id]
		if !exists {
			idx = len(flashbacks)
			idIndex[id] = idx
			flashbacks = append(flashbacks, models.FlashbackWithMetadata{
				Flashback: models.Flashback{
					ID:        id,
					Content:   content,
					Type:      dataType,
					CreatedAt: createdAt,
				},
				Metadata: make(map[string]string),
			})
		}
		if key.Valid && key.String != "" {
			flashbacks[idx].Metadata[key.String] = value.String
		}
	}

	return flashbacks, nil
}

func (app *App) GetAllNotes(ctx context.Context) ([]models.FlashbackWithMetadata, error) {
	query := `
    SELECT f.id, f.content, f.type, f.created_at, m.key, m.value
    FROM flashbacks f
    LEFT JOIN metadata m ON f.id = m.flashback_id
    ORDER BY f.created_at DESC
    `

	rows, err := app.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	flashbacks := []models.FlashbackWithMetadata{}
	idIndex := make(map[string]int)
	for rows.Next() {
		var id, content, dataType, createdAt string
		var key, value sql.NullString
		if err := rows.Scan(&id, &content, &dataType, &createdAt, &key, &value); err != nil {
			return nil, err
		}

		idx, exists := idIndex[id]
		if !exists {
			idx = len(flashbacks)
			idIndex[id] = idx
			flashbacks = append(flashbacks, models.FlashbackWithMetadata{
				Flashback: models.Flashback{
					ID:        id,
					Content:   content,
					Type:      dataType,
					CreatedAt: createdAt,
				},
				Metadata: make(map[string]string),
			})
		}
		if key.Valid && key.String != "" {
			flashbacks[idx].Metadata[key.String] = value.String
		}
	}

	return flashbacks, nil
}

func (app *App) GetNoteByID(ctx context.Context, id string) (models.FlashbackWithMetadata, error) {
	query := `
    SELECT f.id, f.content, f.type, f.created_at, m.key, m.value
    FROM flashbacks f
    LEFT JOIN metadata m ON f.id = m.flashback_id
    WHERE f.id = ?
    `

	rows, err := app.DB.QueryContext(ctx, query, id)
	if err != nil {
		return models.FlashbackWithMetadata{}, err
	}
	defer rows.Close()

	var flashback models.FlashbackWithMetadata
	flashback.Metadata = make(map[string]string)
	found := false

	for rows.Next() {
		var key, value sql.NullString
		if err := rows.Scan(&flashback.ID, &flashback.Content, &flashback.Type, &flashback.CreatedAt, &key, &value); err != nil {
			return models.FlashbackWithMetadata{}, err
		}
		found = true
		if key.Valid && key.String != "" {
			flashback.Metadata[key.String] = value.String
		}
	}

	if !found {
		return models.FlashbackWithMetadata{}, sql.ErrNoRows
	}

	return flashback, nil
}

func (app *App) DeleteNoteByID(ctx context.Context, id string) error {
	deleteQuery := `DELETE FROM flashbacks WHERE id = ?`
	_, err := app.DB.ExecContext(ctx, deleteQuery, id)
	if err != nil {
		return err
	}
	return nil
}
