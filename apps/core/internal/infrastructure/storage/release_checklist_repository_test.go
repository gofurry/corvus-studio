package storage

import (
	"context"
	"errors"
	"path/filepath"
	"testing"
	"time"

	checklistapp "github.com/gofurry/corvus-studio/apps/core/internal/application/checklist"
	releaseapp "github.com/gofurry/corvus-studio/apps/core/internal/application/release"
	checklistdomain "github.com/gofurry/corvus-studio/apps/core/internal/domain/checklist"
	projectdomain "github.com/gofurry/corvus-studio/apps/core/internal/domain/project"
	releasedomain "github.com/gofurry/corvus-studio/apps/core/internal/domain/release"
	generated "github.com/gofurry/corvus-studio/apps/core/internal/infrastructure/storage/sqlc"
	"github.com/gofurry/corvus-studio/apps/core/internal/infrastructure/templatecatalog"
	"github.com/google/uuid"
)

func TestReleaseChecklistWorkspacePersistsAndEnforcesConsistency(t *testing.T) {
	ctx := context.Background()
	databasePath := filepath.Join(t.TempDir(), "corvus.db")
	store := openReleaseTestStore(t, ctx, databasePath)
	projectRepository, releaseRepository, checklistRepository := createReleaseTestRepositories(t, store)
	project := createReleaseTestProject(t, ctx, projectRepository, "Primary")
	catalog, err := templatecatalog.New()
	if err != nil {
		t.Fatal(err)
	}
	clock := &mutableReleaseClock{now: time.Date(2026, time.July, 31, 9, 0, 0, 0, time.UTC)}
	ids := &uuidSequence{}
	releaseService, err := releaseapp.NewService(releaseRepository, projectRepository, catalog, ids, clock)
	if err != nil {
		t.Fatal(err)
	}
	checklistService, err := checklistapp.NewService(checklistRepository, releaseRepository, ids, clock)
	if err != nil {
		t.Fatal(err)
	}

	goal, err := releaseService.Create(ctx, releaseapp.CreateCommand{
		ProjectID: project.ID, GoalType: releasedomain.GoalTypeSteamComingSoon,
		TemplateKey: "steam-coming-soon", TemplateVersion: "1.0.0",
	})
	if err != nil {
		t.Fatalf("create Release workspace: %v", err)
	}
	if goal.Summary.Total != 12 || goal.Summary.RequiredTotal != 7 {
		t.Fatalf("created summary = %#v", goal.Summary)
	}
	if _, err := releaseService.Create(ctx, releaseapp.CreateCommand{
		ProjectID: project.ID, GoalType: releasedomain.GoalTypeSteamComingSoon,
		TemplateKey: "steam-coming-soon", TemplateVersion: "1.0.0",
	}); !errors.Is(err, releasedomain.ErrConflict) {
		t.Fatalf("duplicate create error = %v, want conflict", err)
	}
	if count, err := store.queries.CountChecklistItemsForRelease(ctx, goal.ID.String()); err != nil || count != 12 {
		t.Fatalf("Checklist count after duplicate = %d, %v", count, err)
	}

	items, err := checklistService.List(ctx, goal.ID, checklistdomain.Filter{})
	if err != nil {
		t.Fatal(err)
	}
	corvusSource := checklistdomain.SourceCorvusTemplate
	corvusItems, err := checklistService.List(ctx, goal.ID, checklistdomain.Filter{Source: &corvusSource})
	if err != nil || len(corvusItems) != 4 {
		t.Fatalf("Corvus filter returned %d items, %v", len(corvusItems), err)
	}
	firstRequired := items[0]
	clock.advance(time.Minute)
	firstRequired, err = checklistService.Transition(ctx, firstRequired.ID, checklistdomain.StatusDone)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := checklistService.Transition(ctx, firstRequired.ID, checklistdomain.StatusDone); err != nil {
		t.Fatalf("same-state Checklist retry: %v", err)
	}
	if count, err := store.queries.CountChecklistItemHistory(ctx, firstRequired.ID.String()); err != nil || count != 1 {
		t.Fatalf("Checklist history count = %d, %v", count, err)
	}

	clock.advance(time.Minute)
	goal, err = releaseService.Transition(ctx, goal.ID, releasedomain.StatusPreparing)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := releaseService.Transition(ctx, goal.ID, releasedomain.StatusReadyForReview); !errors.Is(err, releasedomain.ErrNotReady) {
		t.Fatalf("premature review error = %v, want not ready", err)
	}

	for _, item := range items[1:] {
		if item.RequirementLevel != checklistdomain.RequirementRequired {
			continue
		}
		clock.advance(time.Minute)
		if _, err := checklistService.Transition(ctx, item.ID, checklistdomain.StatusDone); err != nil {
			t.Fatalf("complete required item %q: %v", item.Title, err)
		}
	}
	clock.advance(time.Minute)
	goal, err = releaseService.Transition(ctx, goal.ID, releasedomain.StatusReadyForReview)
	if err != nil {
		t.Fatalf("transition ready for review: %v", err)
	}
	if _, err := checklistService.Transition(ctx, firstRequired.ID, checklistdomain.StatusInProgress); !errors.Is(err, checklistdomain.ErrReleaseState) {
		t.Fatalf("reopen locked required item error = %v", err)
	}
	if _, err := checklistService.Create(ctx, checklistapp.CreateCommand{
		ReleaseGoalID: goal.ID, Title: "Late task", Category: checklistdomain.CategoryReview,
		RequirementLevel: checklistdomain.RequirementRecommended,
	}); !errors.Is(err, checklistdomain.ErrReleaseState) {
		t.Fatalf("create task in review error = %v", err)
	}

	clock.advance(time.Minute)
	goal, err = releaseService.Transition(ctx, goal.ID, releasedomain.StatusPreparing)
	if err != nil {
		t.Fatal(err)
	}
	clock.advance(time.Minute)
	custom, err := checklistService.Create(ctx, checklistapp.CreateCommand{
		ReleaseGoalID: goal.ID, Title: "Confirm external review owner", Description: "A user task",
		Requirement: "Record who will review the page", Category: checklistdomain.CategoryReview,
		RequirementLevel: checklistdomain.RequirementRecommended,
	})
	if err != nil {
		t.Fatalf("create user item: %v", err)
	}
	if custom.Source != checklistdomain.SourceUser || custom.SortOrder != 13 {
		t.Fatalf("custom item = %#v", custom)
	}
	if count, err := store.queries.CountReleaseGoalHistory(ctx, goal.ID.String()); err != nil || count != 3 {
		t.Fatalf("Release history count = %d, %v", count, err)
	}

	if err := store.Close(); err != nil {
		t.Fatal(err)
	}
	store = openReleaseTestStore(t, ctx, databasePath)
	t.Cleanup(func() { _ = store.Close() })
	_, releaseRepository, checklistRepository = createReleaseTestRepositories(t, store)
	reopened, err := releaseRepository.Get(ctx, goal.ID)
	if err != nil {
		t.Fatalf("get reopened Release: %v", err)
	}
	if reopened.Status != releasedomain.StatusPreparing || reopened.Summary.Total != 13 || reopened.Summary.RequiredDone != 7 {
		t.Fatalf("reopened Release = %#v", reopened)
	}
	if fetched, err := checklistRepository.Get(ctx, custom.ID); err != nil || fetched.Title != custom.Title {
		t.Fatalf("reopened custom item = %#v, %v", fetched, err)
	}
}

func TestReleaseWorkspaceRollsBackWhenChecklistInsertFails(t *testing.T) {
	ctx := context.Background()
	store := openReleaseTestStore(t, ctx, filepath.Join(t.TempDir(), "corvus.db"))
	t.Cleanup(func() { _ = store.Close() })
	projectRepository, releaseRepository, _ := createReleaseTestRepositories(t, store)
	project := createReleaseTestProject(t, ctx, projectRepository, "Rollback")
	now := time.Date(2026, time.July, 31, 10, 0, 0, 0, time.UTC)
	releaseID := mustReleaseID(t)
	goal, err := releasedomain.New(releasedomain.NewInput{
		ID: releaseID, ProjectID: project.ID, GoalType: releasedomain.GoalTypeSteamComingSoon,
		Title: "Steam Coming Soon", TemplateKey: "steam-coming-soon", TemplateVersion: "1.0.0", Now: now,
	})
	if err != nil {
		t.Fatal(err)
	}
	itemID := mustChecklistID(t)
	key, templateKey, templateVersion := "one", "steam-coming-soon", "1.0.0"
	item, err := checklistdomain.New(checklistdomain.NewInput{
		ID: itemID, ReleaseGoalID: releaseID, Title: "Atomic item", Category: checklistdomain.CategorySetup,
		RequirementLevel: checklistdomain.RequirementRequired, Source: checklistdomain.SourcePlatformTemplate,
		SourceReference: "https://partner.steamgames.com/", TemplateItemKey: &key,
		TemplateKey: &templateKey, TemplateVersion: &templateVersion, SortOrder: 1, Now: now,
	})
	if err != nil {
		t.Fatal(err)
	}
	duplicate := item
	duplicate.SortOrder = 2
	if err := releaseRepository.CreateWorkspace(ctx, goal, []checklistdomain.Item{item, duplicate}); err == nil {
		t.Fatal("CreateWorkspace succeeded with duplicate Checklist ID")
	}
	goalCount, err := store.queries.CountReleaseGoalsByProjectAndType(ctx, generatedCountReleaseParams(project.ID))
	if err != nil || goalCount != 0 {
		t.Fatalf("Release rows after rollback = %d, %v", goalCount, err)
	}
	itemCount, err := store.queries.CountChecklistItemsForRelease(ctx, releaseID.String())
	if err != nil || itemCount != 0 {
		t.Fatalf("Checklist rows after rollback = %d, %v", itemCount, err)
	}
}

func generatedCountReleaseParams(projectID projectdomain.ID) generated.CountReleaseGoalsByProjectAndTypeParams {
	return generated.CountReleaseGoalsByProjectAndTypeParams{ProjectID: projectID.String(), GoalType: string(releasedomain.GoalTypeSteamComingSoon)}
}

func openReleaseTestStore(t *testing.T, ctx context.Context, path string) *Store {
	t.Helper()
	store, err := Open(ctx, path)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	return store
}

func createReleaseTestRepositories(
	t *testing.T,
	store *Store,
) (*ProjectRepository, *ReleaseRepository, *ChecklistRepository) {
	t.Helper()
	projects, err := NewProjectRepository(store)
	if err != nil {
		t.Fatal(err)
	}
	releases, err := NewReleaseRepository(store)
	if err != nil {
		t.Fatal(err)
	}
	checklist, err := NewChecklistRepository(store)
	if err != nil {
		t.Fatal(err)
	}
	return projects, releases, checklist
}

func createReleaseTestProject(
	t *testing.T,
	ctx context.Context,
	repository *ProjectRepository,
	name string,
) projectdomain.Project {
	t.Helper()
	now := time.Date(2026, time.July, 31, 8, 0, 0, 0, time.UTC)
	id, err := projectdomain.IDFromUUID(mustUUIDv7(t))
	if err != nil {
		t.Fatal(err)
	}
	entity, err := projectdomain.New(projectdomain.NewInput{
		ID: id, Name: name, Location: filepath.Join(t.TempDir(), name), LocationKey: "key-" + id.String(),
		Language: "English", Stage: projectdomain.StageReleasePreparation, Now: now,
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := repository.Create(ctx, entity); err != nil {
		t.Fatal(err)
	}
	return entity
}

func mustReleaseID(t *testing.T) releasedomain.ID {
	t.Helper()
	id, err := releasedomain.IDFromUUID(mustUUIDv7(t))
	if err != nil {
		t.Fatal(err)
	}
	return id
}

func mustChecklistID(t *testing.T) checklistdomain.ID {
	t.Helper()
	id, err := checklistdomain.IDFromUUID(mustUUIDv7(t))
	if err != nil {
		t.Fatal(err)
	}
	return id
}

func mustUUIDv7(t *testing.T) uuid.UUID {
	t.Helper()
	value, err := uuid.NewV7()
	if err != nil {
		t.Fatal(err)
	}
	return value
}

type uuidSequence struct{}

func (*uuidSequence) New() (uuid.UUID, error) { return uuid.NewV7() }

type mutableReleaseClock struct{ now time.Time }

func (clock *mutableReleaseClock) Now() time.Time { return clock.now }
func (clock *mutableReleaseClock) advance(delta time.Duration) {
	clock.now = clock.now.Add(delta)
}
