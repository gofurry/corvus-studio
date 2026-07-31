-- +goose Up
CREATE TABLE corvus_runtime_metadata (
    key TEXT PRIMARY KEY NOT NULL CHECK (length(key) > 0),
    value TEXT NOT NULL,
    updated_at TEXT NOT NULL
) STRICT;

-- +goose Down
DROP TABLE corvus_runtime_metadata;
