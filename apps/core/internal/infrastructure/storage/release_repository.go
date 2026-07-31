package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	checklistdomain "github.com/gofurry/corvus-studio/apps/core/internal/domain/checklist"
	projectdomain "github.com/gofurry/corvus-studio/apps/core/internal/domain/project"
	releasedomain "github.com/gofurry/corvus-studio/apps/core/internal/domain/release"
	generated "github.com/gofurry/corvus-studio/apps/core/internal/infrastructure/storage/sqlc"
)

type ReleaseRepository struct {
	db      *sql.DB
	queries *generated.Queries
}

func NewReleaseRepository(store *Store) (*ReleaseRepository, error) {
	if store == nil || store.db == nil || store.queries == nil {
		return nil, errors.New("storage is not open")
	}
	return &ReleaseRepository{db: store.db, queries: store.queries}, nil
}

func (repository *ReleaseRepository) CreateWorkspace(
	ctx context.Context,
	goal releasedomain.Goal,
	items []checklistdomain.Item,
) error {
	if err := goal.Validate(); err != nil {
		return err
	}
	if len(items) == 0 {
		return releasedomain.NewValidationError("checklist", "must contain at least one item")
	}
	for _, item := range items {
		if err := item.Validate(); err != nil {
			return err
		}
		if item.ReleaseGoalID != goal.ID {
			return releasedomain.NewValidationError("checklist", "item belongs to another Release Goal")
		}
	}

	tx, err := repository.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin Release workspace transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }()
	queries := repository.queries.WithTx(tx)
	if err := queries.CreateReleaseGoal(ctx, releaseGoalParams(goal)); err != nil {
		if isUniqueConstraint(err) {
			return releasedomain.ErrConflict
		}
		return fmt.Errorf("insert Release Goal: %w", err)
	}
	for _, item := range items {
		if err := queries.CreateChecklistItem(ctx, checklistItemParams(item)); err != nil {
			return fmt.Errorf("insert generated Checklist item %q: %w", item.ID, err)
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit Release workspace: %w", err)
	}
	return nil
}

func (repository *ReleaseRepository) Get(ctx context.Context, id releasedomain.ID) (releasedomain.Goal, error) {
	row, err := repository.queries.GetReleaseGoal(ctx, id.String())
	if errors.Is(err, sql.ErrNoRows) {
		return releasedomain.Goal{}, releasedomain.ErrNotFound
	}
	if err != nil {
		return releasedomain.Goal{}, fmt.Errorf("select Release Goal: %w", err)
	}
	return releaseGoalFromValues(
		row.ID, row.ProjectID, row.GoalType, row.Title, row.Status, row.TemplateKey, row.TemplateVersion,
		row.CreatedAt, row.UpdatedAt, row.ChecklistTotal, row.ChecklistDone, row.ChecklistBlocked,
		row.RequiredTotal, row.RequiredDone,
	)
}

func (repository *ReleaseRepository) ListByProject(
	ctx context.Context,
	projectID projectdomain.ID,
) ([]releasedomain.Goal, error) {
	rows, err := repository.queries.ListReleaseGoalsByProject(ctx, projectID.String())
	if err != nil {
		return nil, fmt.Errorf("list Release Goals: %w", err)
	}
	goals := make([]releasedomain.Goal, 0, len(rows))
	for _, row := range rows {
		goal, err := releaseGoalFromValues(
			row.ID, row.ProjectID, row.GoalType, row.Title, row.Status, row.TemplateKey, row.TemplateVersion,
			row.CreatedAt, row.UpdatedAt, row.ChecklistTotal, row.ChecklistDone, row.ChecklistBlocked,
			row.RequiredTotal, row.RequiredDone,
		)
		if err != nil {
			return nil, err
		}
		goals = append(goals, goal)
	}
	return goals, nil
}

func (repository *ReleaseRepository) UpdateStatus(
	ctx context.Context,
	before releasedomain.Goal,
	after releasedomain.Goal,
) error {
	if before.ID != after.ID || before.Status == after.Status {
		return releasedomain.ErrStateConflict
	}
	if err := after.Validate(); err != nil {
		return err
	}
	tx, err := repository.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin Release transition: %w", err)
	}
	defer func() { _ = tx.Rollback() }()
	queries := repository.queries.WithTx(tx)
	rows, err := queries.UpdateReleaseGoalStatus(ctx, generated.UpdateReleaseGoalStatusParams{
		Status: string(after.Status), UpdatedAt: after.UpdatedAt.UTC().Format(time.RFC3339Nano),
		ID: before.ID.String(), Status_2: string(before.Status),
	})
	if err != nil {
		return fmt.Errorf("update Release Goal status: %w", err)
	}
	if rows != 1 {
		return releasedomain.ErrStateConflict
	}
	if err := queries.CreateReleaseGoalStatusHistory(ctx, generated.CreateReleaseGoalStatusHistoryParams{
		ReleaseGoalID: before.ID.String(), FromStatus: string(before.Status), ToStatus: string(after.Status),
		ChangedAt: after.UpdatedAt.UTC().Format(time.RFC3339Nano),
	}); err != nil {
		return fmt.Errorf("insert Release Goal history: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit Release transition: %w", err)
	}
	return nil
}

func releaseGoalParams(goal releasedomain.Goal) generated.CreateReleaseGoalParams {
	return generated.CreateReleaseGoalParams{
		ID: goal.ID.String(), ProjectID: goal.ProjectID.String(), GoalType: string(goal.GoalType),
		Title: goal.Title, Status: string(goal.Status), TemplateKey: goal.TemplateKey,
		TemplateVersion: goal.TemplateVersion, CreatedAt: goal.CreatedAt.UTC().Format(time.RFC3339Nano),
		UpdatedAt: goal.UpdatedAt.UTC().Format(time.RFC3339Nano),
	}
}

func releaseGoalFromValues(
	id string,
	projectID string,
	goalType string,
	title string,
	status string,
	templateKey string,
	templateVersion string,
	createdAtValue string,
	updatedAtValue string,
	total int64,
	done int64,
	blocked int64,
	requiredTotal int64,
	requiredDone int64,
) (releasedomain.Goal, error) {
	parsedID, err := releasedomain.ParseID(id)
	if err != nil {
		return releasedomain.Goal{}, fmt.Errorf("decode stored Release Goal ID: %w", err)
	}
	parsedProjectID, err := projectdomain.ParseID(projectID)
	if err != nil {
		return releasedomain.Goal{}, fmt.Errorf("decode stored Release Project ID: %w", err)
	}
	createdAt, err := time.Parse(time.RFC3339Nano, createdAtValue)
	if err != nil {
		return releasedomain.Goal{}, fmt.Errorf("decode Release created_at: %w", err)
	}
	updatedAt, err := time.Parse(time.RFC3339Nano, updatedAtValue)
	if err != nil {
		return releasedomain.Goal{}, fmt.Errorf("decode Release updated_at: %w", err)
	}
	goal := releasedomain.Goal{
		ID: parsedID, ProjectID: parsedProjectID, GoalType: releasedomain.GoalType(goalType), Title: title,
		Status: releasedomain.Status(status), TemplateKey: templateKey, TemplateVersion: templateVersion,
		Summary:   releasedomain.ChecklistSummary{Total: total, Done: done, Blocked: blocked, RequiredTotal: requiredTotal, RequiredDone: requiredDone},
		CreatedAt: createdAt.UTC(), UpdatedAt: updatedAt.UTC(),
	}
	if err := goal.Validate(); err != nil {
		return releasedomain.Goal{}, fmt.Errorf("decode stored Release Goal: %w", err)
	}
	return goal, nil
}
