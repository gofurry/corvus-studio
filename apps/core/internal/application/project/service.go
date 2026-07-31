package project

import (
	"context"
	"errors"
	"fmt"
	"time"

	projectdomain "github.com/gofurry/corvus-studio/apps/core/internal/domain/project"
	"github.com/google/uuid"
)

type Repository interface {
	Create(context.Context, projectdomain.Project) error
	Get(context.Context, projectdomain.ID) (projectdomain.Project, error)
	List(context.Context) ([]projectdomain.Project, error)
}

type DirectoryResolver interface {
	Resolve(string) (location string, locationKey string, err error)
}

type IDGenerator interface {
	New() (projectdomain.ID, error)
}

type Clock interface {
	Now() time.Time
}

type Service struct {
	repository        Repository
	directoryResolver DirectoryResolver
	idGenerator       IDGenerator
	clock             Clock
}

func NewService(
	repository Repository,
	directoryResolver DirectoryResolver,
	idGenerator IDGenerator,
	clock Clock,
) (*Service, error) {
	if repository == nil || directoryResolver == nil || idGenerator == nil || clock == nil {
		return nil, errors.New("project service dependencies are required")
	}
	return &Service{
		repository:        repository,
		directoryResolver: directoryResolver,
		idGenerator:       idGenerator,
		clock:             clock,
	}, nil
}

type CreateCommand struct {
	Name        string
	Description string
	Location    string
	SteamAppID  *int64
	Language    string
	Stage       projectdomain.Stage
}

func (service *Service) Create(ctx context.Context, command CreateCommand) (projectdomain.Project, error) {
	location, locationKey, err := service.directoryResolver.Resolve(command.Location)
	if err != nil {
		return projectdomain.Project{}, err
	}
	id, err := service.idGenerator.New()
	if err != nil {
		return projectdomain.Project{}, fmt.Errorf("generate Project ID: %w", err)
	}
	entity, err := projectdomain.New(projectdomain.NewInput{
		ID:          id,
		Name:        command.Name,
		Description: command.Description,
		Location:    location,
		LocationKey: locationKey,
		SteamAppID:  command.SteamAppID,
		Language:    command.Language,
		Stage:       command.Stage,
		Now:         service.clock.Now(),
	})
	if err != nil {
		return projectdomain.Project{}, err
	}
	if err := service.repository.Create(ctx, entity); err != nil {
		return projectdomain.Project{}, err
	}
	return entity, nil
}

func (service *Service) Get(ctx context.Context, id projectdomain.ID) (projectdomain.Project, error) {
	return service.repository.Get(ctx, id)
}

func (service *Service) List(ctx context.Context) ([]projectdomain.Project, error) {
	return service.repository.List(ctx)
}

type UUIDv7Generator struct{}

func (UUIDv7Generator) New() (projectdomain.ID, error) {
	value, err := uuid.NewV7()
	if err != nil {
		return "", err
	}
	return projectdomain.IDFromUUID(value)
}

type SystemClock struct{}

func (SystemClock) Now() time.Time {
	return time.Now().UTC()
}
