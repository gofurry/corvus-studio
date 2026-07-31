-- +goose Up
CREATE TABLE projects (
    id TEXT PRIMARY KEY NOT NULL CHECK (length(id) = 36),
    name TEXT NOT NULL CHECK (length(trim(name)) BETWEEN 1 AND 120),
    description TEXT NOT NULL DEFAULT '' CHECK (length(description) <= 2000),
    location TEXT NOT NULL CHECK (length(location) > 0),
    location_key TEXT NOT NULL CHECK (length(location_key) > 0),
    steam_app_id INTEGER CHECK (steam_app_id BETWEEN 1 AND 4294967295),
    language TEXT NOT NULL CHECK (length(trim(language)) BETWEEN 1 AND 64),
    stage TEXT NOT NULL CHECK (
        stage IN ('concept', 'development', 'release_preparation', 'released')
    ),
    status TEXT NOT NULL CHECK (status IN ('active', 'archived')),
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
) STRICT;

CREATE UNIQUE INDEX projects_location_key_unique
    ON projects (location_key);

CREATE INDEX projects_created_at_desc
    ON projects (created_at DESC, id DESC);

-- +goose Down
DROP INDEX projects_created_at_desc;
DROP INDEX projects_location_key_unique;
DROP TABLE projects;
