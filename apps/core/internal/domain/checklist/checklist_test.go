package checklist_test

import (
	"errors"
	"testing"
	"time"

	checklistdomain "github.com/gofurry/corvus-studio/apps/core/internal/domain/checklist"
	releasedomain "github.com/gofurry/corvus-studio/apps/core/internal/domain/release"
	"github.com/google/uuid"
)

func TestRequiredChecklistTransitionRules(t *testing.T) {
	item := newItem(t, checklistdomain.RequirementRequired)
	_, _, err := item.Transition(checklistdomain.StatusNotApplicable, releasedomain.StatusPreparing, time.Now())
	if !errors.Is(err, checklistdomain.ErrValidation) {
		t.Fatalf("required not-applicable error = %v, want validation", err)
	}

	done, changed, err := item.Transition(checklistdomain.StatusDone, releasedomain.StatusPreparing, time.Now())
	if err != nil || !changed || done.Status != checklistdomain.StatusDone {
		t.Fatalf("done transition = %#v, %v, %v", done, changed, err)
	}
	_, _, err = done.Transition(checklistdomain.StatusInProgress, releasedomain.StatusReady, time.Now())
	if !errors.Is(err, checklistdomain.ErrReleaseState) {
		t.Fatalf("locked transition error = %v, want release state", err)
	}
}

func TestRecommendedChecklistCanBeNotApplicableAndRetryIsIdempotent(t *testing.T) {
	item := newItem(t, checklistdomain.RequirementRecommended)
	updated, changed, err := item.Transition(checklistdomain.StatusNotApplicable, releasedomain.StatusPreparing, time.Now())
	if err != nil || !changed {
		t.Fatalf("not-applicable transition: %v, %v", changed, err)
	}
	same, changed, err := updated.Transition(checklistdomain.StatusNotApplicable, releasedomain.StatusPreparing, time.Now().Add(time.Hour))
	if err != nil || changed || same.UpdatedAt != updated.UpdatedAt {
		t.Fatalf("same transition changed item: %#v, %v, %v", same, changed, err)
	}
}

func newItem(t *testing.T, level checklistdomain.RequirementLevel) checklistdomain.Item {
	t.Helper()
	releaseID, err := releasedomain.IDFromUUID(mustUUIDv7(t))
	if err != nil {
		t.Fatal(err)
	}
	itemID, err := checklistdomain.IDFromUUID(mustUUIDv7(t))
	if err != nil {
		t.Fatal(err)
	}
	item, err := checklistdomain.New(checklistdomain.NewInput{
		ID: itemID, ReleaseGoalID: releaseID, Title: "Task", Category: checklistdomain.CategoryReview,
		RequirementLevel: level, Source: checklistdomain.SourceUser, SourceReference: "User-created task",
		SortOrder: 1, Now: time.Date(2026, time.July, 31, 8, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatal(err)
	}
	return item
}

func mustUUIDv7(t *testing.T) uuid.UUID {
	t.Helper()
	value, err := uuid.NewV7()
	if err != nil {
		t.Fatal(err)
	}
	return value
}
