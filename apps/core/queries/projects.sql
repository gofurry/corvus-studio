-- name: CreateProject :exec
INSERT INTO projects (
    id,
    name,
    description,
    location,
    location_key,
    steam_app_id,
    language,
    stage,
    status,
    created_at,
    updated_at
)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: GetProject :one
SELECT
    id,
    name,
    description,
    location,
    location_key,
    steam_app_id,
    language,
    stage,
    status,
    created_at,
    updated_at
FROM projects
WHERE id = ?
LIMIT 1;

-- name: ListProjects :many
SELECT
    id,
    name,
    description,
    location,
    location_key,
    steam_app_id,
    language,
    stage,
    status,
    created_at,
    updated_at
FROM projects
ORDER BY created_at DESC, id DESC;
