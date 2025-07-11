-- +goose Up
CREATE TABLE IF NOT EXISTS chunks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    content TEXT NOT NULL,
    note_id INTEGER NOT NULL,
    chunk_number INTEGER NOT NULL,
    FOREIGN KEY (note_id) REFERENCES notes (id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE IF EXISTS chunks;
