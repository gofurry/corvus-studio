-- name: CreateChecklistItem :exec
INSERT INTO checklist_items (
    id, release_goal_id, title, description, requirement, category, requirement_level,
    source, source_reference, template_item_key, template_key, template_version,
    status, sort_order, created_at, updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: GetChecklistItem :one
SELECT
    id, release_goal_id, title, description, requirement, category, requirement_level,
    source, source_reference, template_item_key, template_key, template_version,
    status, sort_order, created_at, updated_at
FROM checklist_items
WHERE id = ?
LIMIT 1;

-- name: ListChecklistItems :many
SELECT
    id, release_goal_id, title, description, requirement, category, requirement_level,
    source, source_reference, template_item_key, template_key, template_version,
    status, sort_order, created_at, updated_at
FROM checklist_items
WHERE release_goal_id = ?
    AND (sqlc.arg(status_filter) = '' OR status = sqlc.arg(status_filter))
    AND (sqlc.arg(source_filter) = '' OR source = sqlc.arg(source_filter))
    AND (sqlc.arg(category_filter) = '' OR category = sqlc.arg(category_filter))
ORDER BY sort_order, id;

-- name: NextChecklistSortOrder :one
SELECT COALESCE(MAX(sort_order), 0) + 1
FROM checklist_items
WHERE release_goal_id = ?;

-- name: UpdateChecklistItemStatus :execrows
UPDATE checklist_items
SET status = ?, updated_at = ?
WHERE id = ? AND status = ?;

-- name: CreateChecklistItemStatusHistory :exec
INSERT INTO checklist_item_status_history (checklist_item_id, from_status, to_status, changed_at)
VALUES (?, ?, ?, ?);

-- name: CountChecklistItemHistory :one
SELECT COUNT(*) FROM checklist_item_status_history WHERE checklist_item_id = ?;

-- name: CountChecklistItemsForRelease :one
SELECT COUNT(*) FROM checklist_items WHERE release_goal_id = ?;
