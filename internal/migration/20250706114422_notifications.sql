-- +goose Up
CREATE TABLE IF NOT EXISTS notifications (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    note_id INTEGER NOT NULL,
    title TEXT NOT NULL,
    time TIMESTAMP NOT NULL,
    FOREIGN KEY (note_id) REFERENCES notes (id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE IF EXISTS notifications;
