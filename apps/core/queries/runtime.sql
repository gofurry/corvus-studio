-- name: Ping :one
SELECT 1 AS ok;

-- name: GetRuntimeMetadata :one
SELECT value
FROM corvus_runtime_metadata
WHERE key = ?
LIMIT 1;

-- name: UpsertRuntimeMetadata :exec
INSERT INTO corvus_runtime_metadata (key, value, updated_at)
VALUES (?, ?, ?)
ON CONFLICT (key) DO UPDATE SET
    value = excluded.value,
    updated_at = excluded.updated_at;
