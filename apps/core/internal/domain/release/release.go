package release

import (
	"errors"
	"fmt"
	"strings"
	"time"

	projectdomain "github.com/gofurry/corvus-studio/apps/core/internal/domain/project"
	"github.com/google/uuid"
)

var (
	ErrValidation        = errors.New("release validation failed")
	ErrNotFound          = errors.New("release goal not found")
	ErrConflict          = errors.New("release goal conflict")
	ErrNotReady          = errors.New("release goal is not ready")
	ErrInvalidTransition = errors.New("invalid release transition")
	ErrStateConflict     = errors.New("release state conflict")
)

type ID string

func ParseID(value string) (ID, error) {
	parsed, err := uuid.Parse(value)
	if err != nil || parsed.Version() != 7 {
		return "", NewValidationError("release_id", "must be a UUIDv7")
	}
	return ID(parsed.String()), nil
}

func IDFromUUID(value uuid.UUID) (ID, error) { return ParseID(value.String()) }
func (id ID) String() string                 { return string(id) }

type GoalType string

const GoalTypeSteamComingSoon GoalType = "steam_coming_soon"

func (goalType GoalType) Valid() bool { return goalType == GoalTypeSteamComingSoon }

type Status string

const (
	StatusDraft          Status = "draft"
	StatusPreparing      Status = "preparing"
	StatusNeedsAttention Status = "needs_attention"
	StatusReadyForReview Status = "ready_for_review"
	StatusReady          Status = "ready"
	StatusSubmitted      Status = "submitted"
)

func (status Status) Valid() bool {
	switch status {
	case StatusDraft, StatusPreparing, StatusNeedsAttention, StatusReadyForReview, StatusReady, StatusSubmitted:
		return true
	default:
		return false
	}
}

func (status Status) LocksRequiredChecklist() bool {
	return status == StatusReadyForReview || status == StatusReady || status == StatusSubmitted
}

func CanTransition(from Status, to Status) bool {
	if from == to {
		return true
	}
	switch from {
	case StatusDraft:
		return to == StatusPreparing
	case StatusPreparing:
		return to == StatusNeedsAttention || to == StatusReadyForReview
	case StatusNeedsAttention:
		return to == StatusPreparing
	case StatusReadyForReview:
		return to == StatusPreparing || to == StatusNeedsAttention || to == StatusReady
	case StatusReady:
		return to == StatusPreparing || to == StatusNeedsAttention || to == StatusSubmitted
	case StatusSubmitted:
		return to == StatusNeedsAttention
	default:
		return false
	}
}

type ChecklistSummary struct {
	Total         int64
	Done          int64
	Blocked       int64
	RequiredTotal int64
	RequiredDone  int64
}

func (summary ChecklistSummary) RequiredComplete() bool {
	return summary.RequiredDone == summary.RequiredTotal
}

type Goal struct {
	ID              ID
	ProjectID       projectdomain.ID
	GoalType        GoalType
	Title           string
	Status          Status
	TemplateKey     string
	TemplateVersion string
	Summary         ChecklistSummary
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type NewInput struct {
	ID              ID
	ProjectID       projectdomain.ID
	GoalType        GoalType
	Title           string
	TemplateKey     string
	TemplateVersion string
	Now             time.Time
}

func New(input NewInput) (Goal, error) {
	now := input.Now.UTC()
	goal := Goal{
		ID:              input.ID,
		ProjectID:       input.ProjectID,
		GoalType:        input.GoalType,
		Title:           strings.TrimSpace(input.Title),
		Status:          StatusDraft,
		TemplateKey:     strings.TrimSpace(input.TemplateKey),
		TemplateVersion: strings.TrimSpace(input.TemplateVersion),
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	if err := goal.Validate(); err != nil {
		return Goal{}, err
	}
	return goal, nil
}

func (goal Goal) Validate() error {
	if _, err := ParseID(goal.ID.String()); err != nil {
		return err
	}
	if _, err := projectdomain.ParseID(goal.ProjectID.String()); err != nil {
		return NewValidationError("project_id", "must be a UUIDv7")
	}
	if !goal.GoalType.Valid() {
		return NewValidationError("goal_type", "must be a supported Release Goal type")
	}
	if goal.Title == "" || len([]rune(goal.Title)) > 160 {
		return NewValidationError("title", "must contain between 1 and 160 characters")
	}
	if !goal.Status.Valid() {
		return NewValidationError("status", "must be a supported Release status")
	}
	if goal.TemplateKey == "" || goal.TemplateVersion == "" {
		return NewValidationError("template", "key and version are required")
	}
	if goal.CreatedAt.IsZero() || goal.UpdatedAt.IsZero() {
		return NewValidationError("timestamps", "must not be zero")
	}
	return nil
}

func (goal Goal) Transition(to Status, now time.Time) (Goal, bool, error) {
	if !to.Valid() {
		return Goal{}, false, NewValidationError("status", "must be a supported Release status")
	}
	if goal.Status == to {
		return goal, false, nil
	}
	if !CanTransition(goal.Status, to) {
		return Goal{}, false, fmt.Errorf("%w: %s to %s", ErrInvalidTransition, goal.Status, to)
	}
	if to == StatusReadyForReview && !goal.Summary.RequiredComplete() {
		return Goal{}, false, fmt.Errorf("%w: required Checklist items remain incomplete", ErrNotReady)
	}
	goal.Status = to
	goal.UpdatedAt = now.UTC()
	return goal, true, nil
}

type ValidationError struct {
	Field   string
	Message string
}

func NewValidationError(field string, message string) *ValidationError {
	return &ValidationError{Field: field, Message: message}
}

func (err *ValidationError) Error() string { return fmt.Sprintf("%s: %s", err.Field, err.Message) }
func (err *ValidationError) Unwrap() error { return ErrValidation }
