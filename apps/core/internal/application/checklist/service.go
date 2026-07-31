package checklist

import (
	"context"
	"errors"
	"fmt"
	"time"

	checklistdomain "github.com/gofurry/corvus-studio/apps/core/internal/domain/checklist"
	releasedomain "github.com/gofurry/corvus-studio/apps/core/internal/domain/release"
	"github.com/google/uuid"
)

type Repository interface {
	Create(context.Context, checklistdomain.Item) error
	Get(context.Context, checklistdomain.ID) (checklistdomain.Item, error)
	List(context.Context, releasedomain.ID, checklistdomain.Filter) ([]checklistdomain.Item, error)
	NextSortOrder(context.Context, releasedomain.ID) (int64, error)
	UpdateStatus(context.Context, checklistdomain.Item, checklistdomain.Item) error
}

type ReleaseReader interface {
	Get(context.Context, releasedomain.ID) (releasedomain.Goal, error)
}

type IDGenerator interface{ New() (uuid.UUID, error) }
type Clock interface{ Now() time.Time }

type Service struct {
	repository  Repository
	releases    ReleaseReader
	idGenerator IDGenerator
	clock       Clock
}

func NewService(repository Repository, releases ReleaseReader, idGenerator IDGenerator, clock Clock) (*Service, error) {
	if repository == nil || releases == nil || idGenerator == nil || clock == nil {
		return nil, errors.New("checklist service dependencies are required")
	}
	return &Service{repository: repository, releases: releases, idGenerator: idGenerator, clock: clock}, nil
}

type CreateCommand struct {
	ReleaseGoalID    releasedomain.ID
	Title            string
	Description      string
	Requirement      string
	Category         checklistdomain.Category
	RequirementLevel checklistdomain.RequirementLevel
}

func (service *Service) Create(ctx context.Context, command CreateCommand) (checklistdomain.Item, error) {
	goal, err := service.releases.Get(ctx, command.ReleaseGoalID)
	if err != nil {
		return checklistdomain.Item{}, err
	}
	if goal.Status.LocksRequiredChecklist() {
		return checklistdomain.Item{}, fmt.Errorf("%w: move the Release Goal back to preparation first", checklistdomain.ErrReleaseState)
	}
	sortOrder, err := service.repository.NextSortOrder(ctx, command.ReleaseGoalID)
	if err != nil {
		return checklistdomain.Item{}, err
	}
	value, err := service.idGenerator.New()
	if err != nil {
		return checklistdomain.Item{}, fmt.Errorf("generate Checklist item ID: %w", err)
	}
	id, err := checklistdomain.IDFromUUID(value)
	if err != nil {
		return checklistdomain.Item{}, err
	}
	item, err := checklistdomain.New(checklistdomain.NewInput{
		ID: id, ReleaseGoalID: command.ReleaseGoalID, Title: command.Title,
		Description: command.Description, Requirement: command.Requirement,
		Category: command.Category, RequirementLevel: command.RequirementLevel,
		Source: checklistdomain.SourceUser, SourceReference: "User-created task",
		SortOrder: sortOrder, Now: service.clock.Now(),
	})
	if err != nil {
		return checklistdomain.Item{}, err
	}
	if err := service.repository.Create(ctx, item); err != nil {
		return checklistdomain.Item{}, err
	}
	return item, nil
}

func (service *Service) Get(ctx context.Context, id checklistdomain.ID) (checklistdomain.Item, error) {
	return service.repository.Get(ctx, id)
}

func (service *Service) List(
	ctx context.Context,
	releaseID releasedomain.ID,
	filter checklistdomain.Filter,
) ([]checklistdomain.Item, error) {
	if err := filter.Validate(); err != nil {
		return nil, err
	}
	if _, err := service.releases.Get(ctx, releaseID); err != nil {
		return nil, err
	}
	return service.repository.List(ctx, releaseID, filter)
}

func (service *Service) Transition(
	ctx context.Context,
	id checklistdomain.ID,
	to checklistdomain.Status,
) (checklistdomain.Item, error) {
	current, err := service.repository.Get(ctx, id)
	if err != nil {
		return checklistdomain.Item{}, err
	}
	goal, err := service.releases.Get(ctx, current.ReleaseGoalID)
	if err != nil {
		return checklistdomain.Item{}, err
	}
	updated, changed, err := current.Transition(to, goal.Status, service.clock.Now())
	if err != nil {
		return checklistdomain.Item{}, err
	}
	if !changed {
		return current, nil
	}
	if err := service.repository.UpdateStatus(ctx, current, updated); err != nil {
		return checklistdomain.Item{}, err
	}
	return updated, nil
}

type UUIDv7Generator struct{}

func (UUIDv7Generator) New() (uuid.UUID, error) { return uuid.NewV7() }

type SystemClock struct{}

func (SystemClock) Now() time.Time { return time.Now().UTC() }
