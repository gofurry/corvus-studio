package project

import (
	"context"
	"errors"
	"testing"
	"time"

	projectdomain "github.com/gofurry/corvus-studio/apps/core/internal/domain/project"
	"github.com/google/uuid"
)

func TestServiceCreateUsesInjectedPorts(t *testing.T) {
	repository := &recordingRepository{}
	id := mustID(t)
	now := time.Date(2026, time.July, 31, 5, 0, 0, 0, time.UTC)
	service, err := NewService(
		repository,
		fakeDirectoryResolver{location: "/canonical/game", key: "/canonical/game"},
		fakeIDGenerator{id: id},
		fakeClock{now: now},
	)
	if err != nil {
		t.Fatalf("new service: %v", err)
	}

	entity, err := service.Create(context.Background(), CreateCommand{
		Name:        "Raven",
		Description: "Demo",
		Location:    "/input/game",
		Language:    "English",
		Stage:       projectdomain.StageConcept,
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if entity.ID != id || entity.Location != "/canonical/game" || entity.CreatedAt != now {
		t.Fatalf("created entity = %#v", entity)
	}
	if repository.created.ID != id {
		t.Fatalf("repository entity = %#v", repository.created)
	}
}

func TestServiceDoesNotWriteRepositoryOnFailure(t *testing.T) {
	repository := &recordingRepository{}
	service, err := NewService(
		repository,
		fakeDirectoryResolver{location: "/canonical/game", key: "/canonical/game"},
		fakeIDGenerator{id: mustID(t)},
		fakeClock{now: time.Now()},
	)
	if err != nil {
		t.Fatalf("new service: %v", err)
	}

	_, err = service.Create(context.Background(), CreateCommand{
		Name:     "",
		Location: "/input/game",
		Language: "English",
		Stage:    projectdomain.StageConcept,
	})
	if !errors.Is(err, projectdomain.ErrValidation) {
		t.Fatalf("error = %v, want validation", err)
	}
	if repository.createCalls != 0 {
		t.Fatalf("repository create calls = %d, want 0", repository.createCalls)
	}
}

type recordingRepository struct {
	created     projectdomain.Project
	createCalls int
}

func (repository *recordingRepository) Create(_ context.Context, entity projectdomain.Project) error {
	repository.createCalls++
	repository.created = entity
	return nil
}

func (repository *recordingRepository) Get(_ context.Context, _ projectdomain.ID) (projectdomain.Project, error) {
	return projectdomain.Project{}, projectdomain.ErrNotFound
}

func (repository *recordingRepository) List(_ context.Context) ([]projectdomain.Project, error) {
	return []projectdomain.Project{}, nil
}

type fakeDirectoryResolver struct {
	location string
	key      string
	err      error
}

func (resolver fakeDirectoryResolver) Resolve(string) (string, string, error) {
	return resolver.location, resolver.key, resolver.err
}

type fakeIDGenerator struct {
	id  projectdomain.ID
	err error
}

func (generator fakeIDGenerator) New() (projectdomain.ID, error) {
	return generator.id, generator.err
}

type fakeClock struct {
	now time.Time
}

func (clock fakeClock) Now() time.Time {
	return clock.now
}

func mustID(t *testing.T) projectdomain.ID {
	t.Helper()
	value, err := uuid.NewV7()
	if err != nil {
		t.Fatalf("generate UUIDv7: %v", err)
	}
	id, err := projectdomain.IDFromUUID(value)
	if err != nil {
		t.Fatalf("convert UUIDv7: %v", err)
	}
	return id
}
