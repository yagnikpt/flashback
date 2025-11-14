package app

import (
	"context"
	"encoding/json"

	"github.com/lithammer/shortuuid/v4"
)

func (app *App) InsertNote(content, dataType string, metadata map[string]string, embeddings []float32) error {
	id := shortuuid.New()
	ctx := context.Background()
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

func (app *App) RetrieveNotesBySimilarity(vector []float32) ([]FlashbackWithMetadata, error) {
	embeddings, err := json.Marshal(vector)
	if err != nil {
		return nil, err
	}

	// 0.40 distance threshold can be applied later if needed
	// WHERE vector_distance_cos(e.vector, vector32(?)) < 0.40
	query := `
    SELECT f.id, f.content, f.type, f.created_at, m.key, m.value
    FROM flashbacks f
    JOIN embeddings e ON f.id = e.flashback_id
    LEFT JOIN metadata m ON f.id = m.flashback_id
    WHERE vector_distance_cos(e.vector, vector32(?)) < 0.40
    ORDER BY vector_distance_cos(e.vector, vector32(?)) ASC LIMIT 20
    `

	rows, err := app.DB.Query(query, string(embeddings), string(embeddings))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// cols, _ := rows.Columns()
	// vals := make([]any, len(cols))
	// raw := make([][]byte, len(cols))

	// for i := range vals {
	// 	vals[i] = &raw[i]
	// }

	// for rows.Next() {
	// 	if err := rows.Scan(vals...); err != nil {
	// 		return nil, err
	// 	}

	// 	for i, col := range cols {
	// 		fmt.Printf("%s=%s ", col, string(raw[i]))
	// 	}
	// 	fmt.Println()
	// }

	flashbacks := []FlashbackWithMetadata{}
	idIndex := make(map[string]int)
	for rows.Next() {
		var id, content, dataType, createdAt, key, value string
		if err := rows.Scan(&id, &content, &dataType, &createdAt, &key, &value); err != nil {
			return nil, err
		}

		idx, exists := idIndex[id]
		if !exists {
			idx = len(flashbacks)
			idIndex[id] = idx
			flashbacks = append(flashbacks, FlashbackWithMetadata{
				Flashback: Flashback{
					ID:        id,
					Content:   content,
					Type:      dataType,
					CreatedAt: createdAt,
					Title:     "",
				},
				Metadata: make(map[string]string),
			})
		}
		if key != "" {
			flashbacks[idx].Metadata[key] = value
		}
	}

	return flashbacks, nil
}
