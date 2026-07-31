package release_test

import (
	"errors"
	"testing"
	"time"

	projectdomain "github.com/gofurry/corvus-studio/apps/core/internal/domain/project"
	releasedomain "github.com/gofurry/corvus-studio/apps/core/internal/domain/release"
	"github.com/google/uuid"
)

func TestReleaseTransitionGraph(t *testing.T) {
	valid := map[releasedomain.Status][]releasedomain.Status{
		releasedomain.StatusDraft:          {releasedomain.StatusPreparing},
		releasedomain.StatusPreparing:      {releasedomain.StatusNeedsAttention, releasedomain.StatusReadyForReview},
		releasedomain.StatusNeedsAttention: {releasedomain.StatusPreparing},
		releasedomain.StatusReadyForReview: {releasedomain.StatusPreparing, releasedomain.StatusNeedsAttention, releasedomain.StatusReady},
		releasedomain.StatusReady:          {releasedomain.StatusPreparing, releasedomain.StatusNeedsAttention, releasedomain.StatusSubmitted},
		releasedomain.StatusSubmitted:      {releasedomain.StatusNeedsAttention},
	}
	statuses := []releasedomain.Status{
		releasedomain.StatusDraft, releasedomain.StatusPreparing, releasedomain.StatusNeedsAttention,
		releasedomain.StatusReadyForReview, releasedomain.StatusReady, releasedomain.StatusSubmitted,
	}
	for _, from := range statuses {
		for _, to := range statuses {
			want := from == to || contains(valid[from], to)
			if got := releasedomain.CanTransition(from, to); got != want {
				t.Errorf("CanTransition(%s, %s) = %v, want %v", from, to, got, want)
			}
		}
	}
}

func TestReleaseRequiresChecklistBeforeReview(t *testing.T) {
	goal := newGoal(t)
	goal.Status = releasedomain.StatusPreparing
	goal.Summary = releasedomain.ChecklistSummary{RequiredTotal: 2, RequiredDone: 1}
	_, _, err := goal.Transition(releasedomain.StatusReadyForReview, time.Now())
	if !errors.Is(err, releasedomain.ErrNotReady) {
		t.Fatalf("transition error = %v, want not ready", err)
	}

	goal.Summary.RequiredDone = 2
	updated, changed, err := goal.Transition(releasedomain.StatusReadyForReview, time.Now())
	if err != nil || !changed || updated.Status != releasedomain.StatusReadyForReview {
		t.Fatalf("ready transition = %#v, %v, %v", updated, changed, err)
	}
	same, changed, err := updated.Transition(releasedomain.StatusReadyForReview, time.Now().Add(time.Hour))
	if err != nil || changed || same.UpdatedAt != updated.UpdatedAt {
		t.Fatalf("same-state transition changed entity: %#v, %v, %v", same, changed, err)
	}
}

func newGoal(t *testing.T) releasedomain.Goal {
	t.Helper()
	now := time.Date(2026, time.July, 31, 8, 0, 0, 0, time.UTC)
	releaseID, err := releasedomain.IDFromUUID(mustUUIDv7(t))
	if err != nil {
		t.Fatal(err)
	}
	projectID, err := projectdomain.IDFromUUID(mustUUIDv7(t))
	if err != nil {
		t.Fatal(err)
	}
	goal, err := releasedomain.New(releasedomain.NewInput{
		ID: releaseID, ProjectID: projectID, GoalType: releasedomain.GoalTypeSteamComingSoon,
		Title: "Steam Coming Soon", TemplateKey: "steam-coming-soon", TemplateVersion: "1.0.0", Now: now,
	})
	if err != nil {
		t.Fatal(err)
	}
	return goal
}

func contains(values []releasedomain.Status, target releasedomain.Status) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}

func mustUUIDv7(t *testing.T) uuid.UUID {
	t.Helper()
	value, err := uuid.NewV7()
	if err != nil {
		t.Fatal(err)
	}
	return value
}
