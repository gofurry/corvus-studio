package storage_test

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	projectapp "github.com/gofurry/corvus-studio/apps/core/internal/application/project"
	projectdomain "github.com/gofurry/corvus-studio/apps/core/internal/domain/project"
	"github.com/gofurry/corvus-studio/apps/core/internal/infrastructure/projectpath"
	"github.com/gofurry/corvus-studio/apps/core/internal/infrastructure/storage"
	"github.com/google/uuid"
)

func TestProjectPersistsAcrossStoreReopen(t *testing.T) {
	ctx := context.Background()
	root := t.TempDir()
	projectDirectory := filepath.Join(root, "game")
	if err := os.Mkdir(projectDirectory, 0o700); err != nil {
		t.Fatalf("create Project directory: %v", err)
	}
	sentinelPath := filepath.Join(projectDirectory, "sentinel.txt")
	if err := os.WriteFile(sentinelPath, []byte("preserve"), 0o600); err != nil {
		t.Fatalf("write sentinel: %v", err)
	}
	databasePath := filepath.Join(root, "data", "corvus.db")
	store := openStore(t, ctx, databasePath)
	repository := newRepository(t, store)
	id := newID(t)
	createdAt := time.Date(2026, time.July, 31, 6, 0, 0, 0, time.UTC)
	service := newService(t, repository, id, createdAt)

	steamAppID := int64(480)
	created, err := service.Create(ctx, projectapp.CreateCommand{
		Name:        "Raven Game",
		Description: "A local-first test Project",
		Location:    projectDirectory,
		SteamAppID:  &steamAppID,
		Language:    "English",
		Stage:       projectdomain.StageDevelopment,
	})
	if err != nil {
		t.Fatalf("create Project: %v", err)
	}
	if created.ID != id {
		t.Fatalf("created ID = %q, want %q", created.ID, id)
	}

	_, err = service.Create(ctx, projectapp.CreateCommand{
		Name:     "Duplicate",
		Location: projectDirectory,
		Language: "English",
		Stage:    projectdomain.StageConcept,
	})
	if !errors.Is(err, projectdomain.ErrLocationConflict) {
		t.Fatalf("duplicate error = %v, want location conflict", err)
	}
	projects, err := service.List(ctx)
	if err != nil {
		t.Fatalf("list Projects: %v", err)
	}
	if len(projects) != 1 {
		t.Fatalf("Project count = %d, want 1", len(projects))
	}

	if err := store.Close(); err != nil {
		t.Fatalf("close first store: %v", err)
	}
	reopened := openStore(t, ctx, databasePath)
	t.Cleanup(func() { _ = reopened.Close() })
	reopenedRepository := newRepository(t, reopened)
	reopenedService := newService(t, reopenedRepository, newID(t), createdAt.Add(time.Hour))
	fetched, err := reopenedService.Get(ctx, id)
	if err != nil {
		t.Fatalf("get reopened Project: %v", err)
	}
	if fetched.Name != created.Name || fetched.Location != created.Location || fetched.CreatedAt != createdAt {
		t.Fatalf("reopened Project = %#v, want %#v", fetched, created)
	}
	content, err := os.ReadFile(sentinelPath)
	if err != nil {
		t.Fatalf("read sentinel: %v", err)
	}
	if string(content) != "preserve" {
		t.Fatalf("Project directory was modified: %q", content)
	}
}

func TestProjectRepositoryReturnsNotFound(t *testing.T) {
	ctx := context.Background()
	store := openStore(t, ctx, filepath.Join(t.TempDir(), "corvus.db"))
	t.Cleanup(func() { _ = store.Close() })
	repository := newRepository(t, store)
	if _, err := repository.Get(ctx, newID(t)); !errors.Is(err, projectdomain.ErrNotFound) {
		t.Fatalf("get missing error = %v, want not found", err)
	}
}

func TestProjectRepositoryListsNewestFirst(t *testing.T) {
	ctx := context.Background()
	root := t.TempDir()
	store := openStore(t, ctx, filepath.Join(root, "corvus.db"))
	t.Cleanup(func() { _ = store.Close() })
	repository := newRepository(t, store)

	olderDirectory := filepath.Join(root, "older")
	newerDirectory := filepath.Join(root, "newer")
	for _, directory := range []string{olderDirectory, newerDirectory} {
		if err := os.Mkdir(directory, 0o700); err != nil {
			t.Fatalf("create Project directory %q: %v", directory, err)
		}
	}
	olderTime := time.Date(2026, time.July, 31, 6, 0, 0, 0, time.UTC)
	newerTime := olderTime.Add(time.Hour)
	olderService := newService(t, repository, newID(t), olderTime)
	newerService := newService(t, repository, newID(t), newerTime)
	for _, input := range []struct {
		service   *projectapp.Service
		name      string
		directory string
	}{
		{service: olderService, name: "Older", directory: olderDirectory},
		{service: newerService, name: "Newer", directory: newerDirectory},
	} {
		if _, err := input.service.Create(ctx, projectapp.CreateCommand{
			Name:     input.name,
			Location: input.directory,
			Language: "English",
			Stage:    projectdomain.StageConcept,
		}); err != nil {
			t.Fatalf("create %s Project: %v", input.name, err)
		}
	}

	projects, err := olderService.List(ctx)
	if err != nil {
		t.Fatalf("list Projects: %v", err)
	}
	if len(projects) != 2 || projects[0].Name != "Newer" || projects[1].Name != "Older" {
		t.Fatalf("Project order = %#v, want Newer then Older", projects)
	}
}

func openStore(t *testing.T, ctx context.Context, path string) *storage.Store {
	t.Helper()
	store, err := storage.Open(ctx, path)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	return store
}

func newRepository(t *testing.T, store *storage.Store) *storage.ProjectRepository {
	t.Helper()
	repository, err := storage.NewProjectRepository(store)
	if err != nil {
		t.Fatalf("new Project repository: %v", err)
	}
	return repository
}

func newService(
	t *testing.T,
	repository *storage.ProjectRepository,
	id projectdomain.ID,
	now time.Time,
) *projectapp.Service {
	t.Helper()
	service, err := projectapp.NewService(
		repository,
		projectpath.Resolver{},
		fixedIDGenerator{id: id},
		fixedClock{now: now},
	)
	if err != nil {
		t.Fatalf("new Project service: %v", err)
	}
	return service
}

func newID(t *testing.T) projectdomain.ID {
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

type fixedIDGenerator struct {
	id projectdomain.ID
}

func (generator fixedIDGenerator) New() (projectdomain.ID, error) {
	return generator.id, nil
}

type fixedClock struct {
	now time.Time
}

func (clock fixedClock) Now() time.Time {
	return clock.now
}
