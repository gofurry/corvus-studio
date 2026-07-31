package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	checklistdomain "github.com/gofurry/corvus-studio/apps/core/internal/domain/checklist"
	releasedomain "github.com/gofurry/corvus-studio/apps/core/internal/domain/release"
	generated "github.com/gofurry/corvus-studio/apps/core/internal/infrastructure/storage/sqlc"
)

type ChecklistRepository struct {
	db      *sql.DB
	queries *generated.Queries
}

func NewChecklistRepository(store *Store) (*ChecklistRepository, error) {
	if store == nil || store.db == nil || store.queries == nil {
		return nil, errors.New("storage is not open")
	}
	return &ChecklistRepository{db: store.db, queries: store.queries}, nil
}

func (repository *ChecklistRepository) Create(ctx context.Context, item checklistdomain.Item) error {
	if err := item.Validate(); err != nil {
		return err
	}
	if err := repository.queries.CreateChecklistItem(ctx, checklistItemParams(item)); err != nil {
		return fmt.Errorf("insert Checklist item: %w", err)
	}
	return nil
}

func (repository *ChecklistRepository) Get(
	ctx context.Context,
	id checklistdomain.ID,
) (checklistdomain.Item, error) {
	row, err := repository.queries.GetChecklistItem(ctx, id.String())
	if errors.Is(err, sql.ErrNoRows) {
		return checklistdomain.Item{}, checklistdomain.ErrNotFound
	}
	if err != nil {
		return checklistdomain.Item{}, fmt.Errorf("select Checklist item: %w", err)
	}
	return checklistItemFromRow(row)
}

func (repository *ChecklistRepository) List(
	ctx context.Context,
	releaseID releasedomain.ID,
	filter checklistdomain.Filter,
) ([]checklistdomain.Item, error) {
	if err := filter.Validate(); err != nil {
		return nil, err
	}
	rows, err := repository.queries.ListChecklistItems(ctx, generated.ListChecklistItemsParams{
		ReleaseGoalID: releaseID.String(), StatusFilter: filterValue(filter.Status),
		SourceFilter: filterValue(filter.Source), CategoryFilter: filterValue(filter.Category),
	})
	if err != nil {
		return nil, fmt.Errorf("list Checklist items: %w", err)
	}
	items := make([]checklistdomain.Item, 0, len(rows))
	for _, row := range rows {
		item, err := checklistItemFromRow(row)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (repository *ChecklistRepository) NextSortOrder(ctx context.Context, releaseID releasedomain.ID) (int64, error) {
	order, err := repository.queries.NextChecklistSortOrder(ctx, releaseID.String())
	if err != nil {
		return 0, fmt.Errorf("calculate Checklist order: %w", err)
	}
	return order, nil
}

func (repository *ChecklistRepository) UpdateStatus(
	ctx context.Context,
	before checklistdomain.Item,
	after checklistdomain.Item,
) error {
	if before.ID != after.ID || before.Status == after.Status {
		return checklistdomain.ErrReleaseState
	}
	if err := after.Validate(); err != nil {
		return err
	}
	tx, err := repository.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin Checklist transition: %w", err)
	}
	defer func() { _ = tx.Rollback() }()
	queries := repository.queries.WithTx(tx)
	rows, err := queries.UpdateChecklistItemStatus(ctx, generated.UpdateChecklistItemStatusParams{
		Status: string(after.Status), UpdatedAt: after.UpdatedAt.UTC().Format(time.RFC3339Nano),
		ID: before.ID.String(), Status_2: string(before.Status),
	})
	if err != nil {
		return fmt.Errorf("update Checklist status: %w", err)
	}
	if rows != 1 {
		return checklistdomain.ErrReleaseState
	}
	if err := queries.CreateChecklistItemStatusHistory(ctx, generated.CreateChecklistItemStatusHistoryParams{
		ChecklistItemID: before.ID.String(), FromStatus: string(before.Status), ToStatus: string(after.Status),
		ChangedAt: after.UpdatedAt.UTC().Format(time.RFC3339Nano),
	}); err != nil {
		return fmt.Errorf("insert Checklist history: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit Checklist transition: %w", err)
	}
	return nil
}

func checklistItemParams(item checklistdomain.Item) generated.CreateChecklistItemParams {
	return generated.CreateChecklistItemParams{
		ID: item.ID.String(), ReleaseGoalID: item.ReleaseGoalID.String(), Title: item.Title,
		Description: item.Description, Requirement: item.Requirement, Category: string(item.Category),
		RequirementLevel: string(item.RequirementLevel), Source: string(item.Source),
		SourceReference: item.SourceReference, TemplateItemKey: nullableString(item.TemplateItemKey),
		TemplateKey: nullableString(item.TemplateKey), TemplateVersion: nullableString(item.TemplateVersion),
		Status: string(item.Status), SortOrder: item.SortOrder,
		CreatedAt: item.CreatedAt.UTC().Format(time.RFC3339Nano), UpdatedAt: item.UpdatedAt.UTC().Format(time.RFC3339Nano),
	}
}

func checklistItemFromRow(row generated.ChecklistItem) (checklistdomain.Item, error) {
	id, err := checklistdomain.ParseID(row.ID)
	if err != nil {
		return checklistdomain.Item{}, fmt.Errorf("decode stored Checklist ID: %w", err)
	}
	releaseID, err := releasedomain.ParseID(row.ReleaseGoalID)
	if err != nil {
		return checklistdomain.Item{}, fmt.Errorf("decode stored Checklist Release ID: %w", err)
	}
	createdAt, err := time.Parse(time.RFC3339Nano, row.CreatedAt)
	if err != nil {
		return checklistdomain.Item{}, fmt.Errorf("decode Checklist created_at: %w", err)
	}
	updatedAt, err := time.Parse(time.RFC3339Nano, row.UpdatedAt)
	if err != nil {
		return checklistdomain.Item{}, fmt.Errorf("decode Checklist updated_at: %w", err)
	}
	item := checklistdomain.Item{
		ID: id, ReleaseGoalID: releaseID, Title: row.Title, Description: row.Description,
		Requirement: row.Requirement, Category: checklistdomain.Category(row.Category),
		RequirementLevel: checklistdomain.RequirementLevel(row.RequirementLevel), Source: checklistdomain.Source(row.Source),
		SourceReference: row.SourceReference, TemplateItemKey: pointerFromNullString(row.TemplateItemKey),
		TemplateKey: pointerFromNullString(row.TemplateKey), TemplateVersion: pointerFromNullString(row.TemplateVersion),
		Status: checklistdomain.Status(row.Status), SortOrder: row.SortOrder,
		CreatedAt: createdAt.UTC(), UpdatedAt: updatedAt.UTC(),
	}
	if err := item.Validate(); err != nil {
		return checklistdomain.Item{}, fmt.Errorf("decode stored Checklist item: %w", err)
	}
	return item, nil
}

func nullableString(value *string) sql.NullString {
	if value == nil {
		return sql.NullString{}
	}
	return sql.NullString{String: *value, Valid: true}
}

func pointerFromNullString(value sql.NullString) *string {
	if !value.Valid {
		return nil
	}
	copy := value.String
	return &copy
}

func filterValue[T ~string](value *T) interface{} {
	if value == nil {
		return ""
	}
	return string(*value)
}
