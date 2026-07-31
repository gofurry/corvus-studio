-- +goose Up
CREATE TABLE release_goals (
    id TEXT PRIMARY KEY NOT NULL CHECK (length(id) = 36),
    project_id TEXT NOT NULL REFERENCES projects (id) ON DELETE RESTRICT,
    goal_type TEXT NOT NULL CHECK (goal_type IN ('steam_coming_soon')),
    title TEXT NOT NULL CHECK (length(trim(title)) BETWEEN 1 AND 160),
    status TEXT NOT NULL CHECK (
        status IN ('draft', 'preparing', 'needs_attention', 'ready_for_review', 'ready', 'submitted')
    ),
    template_key TEXT NOT NULL CHECK (length(trim(template_key)) > 0),
    template_version TEXT NOT NULL CHECK (length(trim(template_version)) > 0),
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    UNIQUE (project_id, goal_type)
) STRICT;

CREATE INDEX release_goals_project_created_at
    ON release_goals (project_id, created_at DESC, id DESC);

CREATE TABLE checklist_items (
    id TEXT PRIMARY KEY NOT NULL CHECK (length(id) = 36),
    release_goal_id TEXT NOT NULL REFERENCES release_goals (id) ON DELETE CASCADE,
    title TEXT NOT NULL CHECK (length(trim(title)) BETWEEN 1 AND 160),
    description TEXT NOT NULL DEFAULT '' CHECK (length(description) <= 2000),
    requirement TEXT NOT NULL DEFAULT '' CHECK (length(requirement) <= 2000),
    category TEXT NOT NULL CHECK (
        category IN ('setup', 'store_copy', 'branding', 'media', 'compliance', 'timeline', 'positioning', 'localization', 'review')
    ),
    requirement_level TEXT NOT NULL CHECK (requirement_level IN ('required', 'recommended')),
    source TEXT NOT NULL CHECK (source IN ('platform_template', 'corvus_template', 'user', 'agent')),
    source_reference TEXT NOT NULL CHECK (length(trim(source_reference)) > 0),
    template_item_key TEXT,
    template_key TEXT,
    template_version TEXT,
    status TEXT NOT NULL CHECK (
        status IN ('not_started', 'in_progress', 'needs_review', 'done', 'blocked', 'not_applicable')
    ),
    sort_order INTEGER NOT NULL CHECK (sort_order > 0),
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    CHECK (requirement_level != 'required' OR status != 'not_applicable'),
    CHECK (
        (source IN ('platform_template', 'corvus_template')
            AND template_item_key IS NOT NULL
            AND template_key IS NOT NULL
            AND template_version IS NOT NULL)
        OR (source = 'user'
            AND template_item_key IS NULL
            AND template_key IS NULL
            AND template_version IS NULL)
        OR source = 'agent'
    ),
    UNIQUE (release_goal_id, template_item_key)
) STRICT;

CREATE INDEX checklist_items_release_order
    ON checklist_items (release_goal_id, sort_order, id);

CREATE INDEX checklist_items_release_filters
    ON checklist_items (release_goal_id, status, source, category);

CREATE TABLE release_goal_status_history (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    release_goal_id TEXT NOT NULL REFERENCES release_goals (id) ON DELETE CASCADE,
    from_status TEXT NOT NULL,
    to_status TEXT NOT NULL,
    changed_at TEXT NOT NULL,
    CHECK (from_status != to_status)
) STRICT;

CREATE INDEX release_goal_status_history_release
    ON release_goal_status_history (release_goal_id, id);

CREATE TABLE checklist_item_status_history (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    checklist_item_id TEXT NOT NULL REFERENCES checklist_items (id) ON DELETE CASCADE,
    from_status TEXT NOT NULL,
    to_status TEXT NOT NULL,
    changed_at TEXT NOT NULL,
    CHECK (from_status != to_status)
) STRICT;

CREATE INDEX checklist_item_status_history_item
    ON checklist_item_status_history (checklist_item_id, id);

-- +goose Down
DROP INDEX checklist_item_status_history_item;
DROP TABLE checklist_item_status_history;
DROP INDEX release_goal_status_history_release;
DROP TABLE release_goal_status_history;
DROP INDEX checklist_items_release_filters;
DROP INDEX checklist_items_release_order;
DROP TABLE checklist_items;
DROP INDEX release_goals_project_created_at;
DROP TABLE release_goals;
