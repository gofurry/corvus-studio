-- name: CreateReleaseGoal :exec
INSERT INTO release_goals (
    id, project_id, goal_type, title, status, template_key, template_version, created_at, updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: GetReleaseGoal :one
SELECT
    release_goals.id,
    release_goals.project_id,
    release_goals.goal_type,
    release_goals.title,
    release_goals.status,
    release_goals.template_key,
    release_goals.template_version,
    release_goals.created_at,
    release_goals.updated_at,
    COUNT(checklist_items.id) AS checklist_total,
    CAST(COALESCE(SUM(CASE WHEN checklist_items.status = 'done' THEN 1 ELSE 0 END), 0) AS INTEGER) AS checklist_done,
    CAST(COALESCE(SUM(CASE WHEN checklist_items.status = 'blocked' THEN 1 ELSE 0 END), 0) AS INTEGER) AS checklist_blocked,
    CAST(COALESCE(SUM(CASE WHEN checklist_items.requirement_level = 'required' THEN 1 ELSE 0 END), 0) AS INTEGER) AS required_total,
    CAST(COALESCE(SUM(CASE WHEN checklist_items.requirement_level = 'required' AND checklist_items.status = 'done' THEN 1 ELSE 0 END), 0) AS INTEGER) AS required_done
FROM release_goals
LEFT JOIN checklist_items ON checklist_items.release_goal_id = release_goals.id
WHERE release_goals.id = ?
GROUP BY release_goals.id
LIMIT 1;

-- name: ListReleaseGoalsByProject :many
SELECT
    release_goals.id,
    release_goals.project_id,
    release_goals.goal_type,
    release_goals.title,
    release_goals.status,
    release_goals.template_key,
    release_goals.template_version,
    release_goals.created_at,
    release_goals.updated_at,
    COUNT(checklist_items.id) AS checklist_total,
    CAST(COALESCE(SUM(CASE WHEN checklist_items.status = 'done' THEN 1 ELSE 0 END), 0) AS INTEGER) AS checklist_done,
    CAST(COALESCE(SUM(CASE WHEN checklist_items.status = 'blocked' THEN 1 ELSE 0 END), 0) AS INTEGER) AS checklist_blocked,
    CAST(COALESCE(SUM(CASE WHEN checklist_items.requirement_level = 'required' THEN 1 ELSE 0 END), 0) AS INTEGER) AS required_total,
    CAST(COALESCE(SUM(CASE WHEN checklist_items.requirement_level = 'required' AND checklist_items.status = 'done' THEN 1 ELSE 0 END), 0) AS INTEGER) AS required_done
FROM release_goals
LEFT JOIN checklist_items ON checklist_items.release_goal_id = release_goals.id
WHERE release_goals.project_id = ?
GROUP BY release_goals.id
ORDER BY release_goals.created_at DESC, release_goals.id DESC;

-- name: UpdateReleaseGoalStatus :execrows
UPDATE release_goals
SET status = ?, updated_at = ?
WHERE id = ? AND status = ?;

-- name: CreateReleaseGoalStatusHistory :exec
INSERT INTO release_goal_status_history (release_goal_id, from_status, to_status, changed_at)
VALUES (?, ?, ?, ?);

-- name: CountReleaseGoalHistory :one
SELECT COUNT(*) FROM release_goal_status_history WHERE release_goal_id = ?;

-- name: CountReleaseGoalsByProjectAndType :one
SELECT COUNT(*) FROM release_goals WHERE project_id = ? AND goal_type = ?;
