package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	projectdomain "github.com/gofurry/corvus-studio/apps/core/internal/domain/project"
	generated "github.com/gofurry/corvus-studio/apps/core/internal/infrastructure/storage/sqlc"
	"github.com/google/uuid"
	sqliteDriver "modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

type ProjectRepository struct {
	queries generated.Querier
}

func NewProjectRepository(store *Store) (*ProjectRepository, error) {
	if store == nil || store.queries == nil {
		return nil, errors.New("storage is not open")
	}
	return &ProjectRepository{queries: store.queries}, nil
}

func (repository *ProjectRepository) Create(
	ctx context.Context,
	entity projectdomain.Project,
) error {
	if err := entity.Validate(); err != nil {
		return err
	}
	err := repository.queries.CreateProject(ctx, generated.CreateProjectParams{
		ID:          entity.ID.String(),
		Name:        entity.Name,
		Description: entity.Description,
		Location:    entity.Location,
		LocationKey: entity.LocationKey,
		SteamAppID:  nullableInt64(entity.SteamAppID),
		Language:    entity.Language,
		Stage:       string(entity.Stage),
		Status:      string(entity.Status),
		CreatedAt:   entity.CreatedAt.UTC().Format(time.RFC3339Nano),
		UpdatedAt:   entity.UpdatedAt.UTC().Format(time.RFC3339Nano),
	})
	if isUniqueConstraint(err) {
		return fmt.Errorf("%w: %s", projectdomain.ErrLocationConflict, entity.Location)
	}
	if err != nil {
		return fmt.Errorf("insert Project: %w", err)
	}
	return nil
}

func (repository *ProjectRepository) Get(
	ctx context.Context,
	id projectdomain.ID,
) (projectdomain.Project, error) {
	row, err := repository.queries.GetProject(ctx, id.String())
	if errors.Is(err, sql.ErrNoRows) {
		return projectdomain.Project{}, projectdomain.ErrNotFound
	}
	if err != nil {
		return projectdomain.Project{}, fmt.Errorf("select Project: %w", err)
	}
	return projectFromRow(row)
}

func (repository *ProjectRepository) List(ctx context.Context) ([]projectdomain.Project, error) {
	rows, err := repository.queries.ListProjects(ctx)
	if err != nil {
		return nil, fmt.Errorf("list Projects: %w", err)
	}
	projects := make([]projectdomain.Project, 0, len(rows))
	for _, row := range rows {
		entity, err := projectFromRow(row)
		if err != nil {
			return nil, err
		}
		projects = append(projects, entity)
	}
	return projects, nil
}

func projectFromRow(row generated.Project) (projectdomain.Project, error) {
	parsedID, err := uuid.Parse(row.ID)
	if err != nil || parsedID.Version() != 7 {
		return projectdomain.Project{}, fmt.Errorf("decode stored Project ID %q: expected UUIDv7", row.ID)
	}
	createdAt, err := time.Parse(time.RFC3339Nano, row.CreatedAt)
	if err != nil {
		return projectdomain.Project{}, fmt.Errorf("decode Project created_at: %w", err)
	}
	updatedAt, err := time.Parse(time.RFC3339Nano, row.UpdatedAt)
	if err != nil {
		return projectdomain.Project{}, fmt.Errorf("decode Project updated_at: %w", err)
	}
	stage := projectdomain.Stage(row.Stage)
	if !stage.Valid() {
		return projectdomain.Project{}, fmt.Errorf("decode stored Project stage %q", row.Stage)
	}
	status := projectdomain.Status(row.Status)
	if !status.Valid() {
		return projectdomain.Project{}, fmt.Errorf("decode stored Project status %q", row.Status)
	}
	return projectdomain.Project{
		ID:          projectdomain.ID(parsedID.String()),
		Name:        row.Name,
		Description: row.Description,
		Location:    row.Location,
		LocationKey: row.LocationKey,
		SteamAppID:  pointerFromNullInt64(row.SteamAppID),
		Language:    row.Language,
		Stage:       stage,
		Status:      status,
		CreatedAt:   createdAt.UTC(),
		UpdatedAt:   updatedAt.UTC(),
	}, nil
}

func nullableInt64(value *int64) sql.NullInt64 {
	if value == nil {
		return sql.NullInt64{}
	}
	return sql.NullInt64{Int64: *value, Valid: true}
}

func pointerFromNullInt64(value sql.NullInt64) *int64 {
	if !value.Valid {
		return nil
	}
	copy := value.Int64
	return &copy
}

func isUniqueConstraint(err error) bool {
	var sqliteError *sqliteDriver.Error
	return errors.As(err, &sqliteError) && sqliteError.Code() == sqlite3.SQLITE_CONSTRAINT_UNIQUE
}
