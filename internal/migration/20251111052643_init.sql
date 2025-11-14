-- +goose Up
CREATE TABLE IF NOT EXISTS flashbacks (
    id TEXT PRIMARY KEY,
    title TEXT,
    content TEXT NOT NULL,
    type TEXT CHECK(type IN ('url', 'text')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE IF NOT EXISTS embeddings (
    flashback_id TEXT,
    vector F32_BLOB(768) NOT NULL,
    FOREIGN KEY (flashback_id) REFERENCES flashbacks(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE IF EXISTS embeddings;
DROP TABLE IF EXISTS flashbacks;
