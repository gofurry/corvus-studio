package release

import (
	"context"
	"errors"
	"fmt"
	"time"

	checklistdomain "github.com/gofurry/corvus-studio/apps/core/internal/domain/checklist"
	projectdomain "github.com/gofurry/corvus-studio/apps/core/internal/domain/project"
	releasedomain "github.com/gofurry/corvus-studio/apps/core/internal/domain/release"
	templatedomain "github.com/gofurry/corvus-studio/apps/core/internal/domain/template"
	"github.com/google/uuid"
)

type Repository interface {
	CreateWorkspace(context.Context, releasedomain.Goal, []checklistdomain.Item) error
	Get(context.Context, releasedomain.ID) (releasedomain.Goal, error)
	ListByProject(context.Context, projectdomain.ID) ([]releasedomain.Goal, error)
	UpdateStatus(context.Context, releasedomain.Goal, releasedomain.Goal) error
}

type ProjectReader interface {
	Get(context.Context, projectdomain.ID) (projectdomain.Project, error)
}

type TemplateCatalog interface {
	Get(string, string) (templatedomain.Definition, error)
	Latest(string) (templatedomain.Definition, error)
}

type IDGenerator interface{ New() (uuid.UUID, error) }
type Clock interface{ Now() time.Time }

type Service struct {
	repository  Repository
	projects    ProjectReader
	templates   TemplateCatalog
	idGenerator IDGenerator
	clock       Clock
}

func NewService(
	repository Repository,
	projects ProjectReader,
	templates TemplateCatalog,
	idGenerator IDGenerator,
	clock Clock,
) (*Service, error) {
	if repository == nil || projects == nil || templates == nil || idGenerator == nil || clock == nil {
		return nil, errors.New("release service dependencies are required")
	}
	return &Service{repository: repository, projects: projects, templates: templates, idGenerator: idGenerator, clock: clock}, nil
}

type CreateCommand struct {
	ProjectID       projectdomain.ID
	GoalType        releasedomain.GoalType
	TemplateKey     string
	TemplateVersion string
}

func (service *Service) Create(ctx context.Context, command CreateCommand) (releasedomain.Goal, error) {
	project, err := service.projects.Get(ctx, command.ProjectID)
	if err != nil {
		return releasedomain.Goal{}, err
	}
	if project.Status != projectdomain.StatusActive {
		return releasedomain.Goal{}, fmt.Errorf("%w: Project is not active", releasedomain.ErrStateConflict)
	}
	definition, err := service.templates.Get(command.TemplateKey, command.TemplateVersion)
	if err != nil {
		return releasedomain.Goal{}, err
	}
	if command.GoalType != definition.GoalType {
		return releasedomain.Goal{}, releasedomain.NewValidationError("goal_type", "does not match the selected template")
	}
	now := service.clock.Now().UTC()
	releaseUUID, err := service.idGenerator.New()
	if err != nil {
		return releasedomain.Goal{}, fmt.Errorf("generate Release Goal ID: %w", err)
	}
	releaseID, err := releasedomain.IDFromUUID(releaseUUID)
	if err != nil {
		return releasedomain.Goal{}, err
	}
	goal, err := releasedomain.New(releasedomain.NewInput{
		ID: releaseID, ProjectID: command.ProjectID, GoalType: command.GoalType, Title: definition.Name,
		TemplateKey: definition.Key, TemplateVersion: definition.Version, Now: now,
	})
	if err != nil {
		return releasedomain.Goal{}, err
	}

	items := make([]checklistdomain.Item, 0, len(definition.Items))
	for _, definitionItem := range definition.Items {
		itemUUID, err := service.idGenerator.New()
		if err != nil {
			return releasedomain.Goal{}, fmt.Errorf("generate Checklist item ID: %w", err)
		}
		itemID, err := checklistdomain.IDFromUUID(itemUUID)
		if err != nil {
			return releasedomain.Goal{}, err
		}
		templateItemKey, templateKey, templateVersion := definitionItem.Key, definition.Key, definition.Version
		item, err := checklistdomain.New(checklistdomain.NewInput{
			ID: itemID, ReleaseGoalID: releaseID, Title: definitionItem.Title,
			Description: definitionItem.Description, Requirement: definitionItem.Requirement,
			Category: definitionItem.Category, RequirementLevel: definitionItem.RequirementLevel,
			Source: definitionItem.Source, SourceReference: definitionItem.SourceReference,
			TemplateItemKey: &templateItemKey, TemplateKey: &templateKey, TemplateVersion: &templateVersion,
			SortOrder: definitionItem.SortOrder, Now: now,
		})
		if err != nil {
			return releasedomain.Goal{}, err
		}
		items = append(items, item)
		goal.Summary.Total++
		if item.RequirementLevel == checklistdomain.RequirementRequired {
			goal.Summary.RequiredTotal++
		}
	}
	if err := service.repository.CreateWorkspace(ctx, goal, items); err != nil {
		return releasedomain.Goal{}, err
	}
	return goal, nil
}

func (service *Service) Get(ctx context.Context, id releasedomain.ID) (releasedomain.Goal, error) {
	return service.repository.Get(ctx, id)
}

func (service *Service) ListByProject(ctx context.Context, projectID projectdomain.ID) ([]releasedomain.Goal, error) {
	if _, err := service.projects.Get(ctx, projectID); err != nil {
		return nil, err
	}
	return service.repository.ListByProject(ctx, projectID)
}

func (service *Service) Transition(
	ctx context.Context,
	id releasedomain.ID,
	to releasedomain.Status,
) (releasedomain.Goal, error) {
	current, err := service.repository.Get(ctx, id)
	if err != nil {
		return releasedomain.Goal{}, err
	}
	updated, changed, err := current.Transition(to, service.clock.Now())
	if err != nil {
		return releasedomain.Goal{}, err
	}
	if !changed {
		return current, nil
	}
	if err := service.repository.UpdateStatus(ctx, current, updated); err != nil {
		return releasedomain.Goal{}, err
	}
	return updated, nil
}

func (service *Service) LatestTemplate(key string) (templatedomain.Definition, error) {
	return service.templates.Latest(key)
}

type UUIDv7Generator struct{}

func (UUIDv7Generator) New() (uuid.UUID, error) { return uuid.NewV7() }

type SystemClock struct{}

func (SystemClock) Now() time.Time { return time.Now().UTC() }
