-- +goose Up
CREATE TABLE IF NOT EXISTS metadata (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    flashback_id TEXT NOT NULL,
    key TEXT NOT NULL,
    value TEXT,
    source TEXT DEFAULT 'system',   -- e.g. "opengraph", "gemini", "user"
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (flashback_id) REFERENCES flashbacks(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE IF EXISTS metadata;
