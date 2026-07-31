package checklist

import (
	"errors"
	"fmt"
	"strings"
	"time"

	releasedomain "github.com/gofurry/corvus-studio/apps/core/internal/domain/release"
	"github.com/google/uuid"
)

const (
	MaxTitleLength       = 160
	MaxDescriptionLength = 2000
	MaxRequirementLength = 2000
)

var (
	ErrValidation        = errors.New("checklist validation failed")
	ErrNotFound          = errors.New("checklist item not found")
	ErrInvalidTransition = errors.New("invalid checklist transition")
	ErrReleaseState      = errors.New("release state conflict")
)

type ID string

func ParseID(value string) (ID, error) {
	parsed, err := uuid.Parse(value)
	if err != nil || parsed.Version() != 7 {
		return "", NewValidationError("item_id", "must be a UUIDv7")
	}
	return ID(parsed.String()), nil
}

func IDFromUUID(value uuid.UUID) (ID, error) { return ParseID(value.String()) }
func (id ID) String() string                 { return string(id) }

type Status string

const (
	StatusNotStarted    Status = "not_started"
	StatusInProgress    Status = "in_progress"
	StatusNeedsReview   Status = "needs_review"
	StatusDone          Status = "done"
	StatusBlocked       Status = "blocked"
	StatusNotApplicable Status = "not_applicable"
)

func (status Status) Valid() bool {
	switch status {
	case StatusNotStarted, StatusInProgress, StatusNeedsReview, StatusDone, StatusBlocked, StatusNotApplicable:
		return true
	default:
		return false
	}
}

func CanTransition(from Status, to Status) bool {
	if from == to {
		return true
	}
	switch from {
	case StatusNotStarted:
		return to == StatusInProgress || to == StatusBlocked || to == StatusDone || to == StatusNotApplicable
	case StatusInProgress:
		return to == StatusNotStarted || to == StatusNeedsReview || to == StatusBlocked || to == StatusDone || to == StatusNotApplicable
	case StatusNeedsReview:
		return to == StatusInProgress || to == StatusBlocked || to == StatusDone || to == StatusNotApplicable
	case StatusBlocked:
		return to == StatusNotStarted || to == StatusInProgress || to == StatusDone || to == StatusNotApplicable
	case StatusDone:
		return to == StatusInProgress
	case StatusNotApplicable:
		return to == StatusNotStarted
	default:
		return false
	}
}

type Source string

const (
	SourcePlatformTemplate Source = "platform_template"
	SourceCorvusTemplate   Source = "corvus_template"
	SourceUser             Source = "user"
	SourceAgent            Source = "agent"
)

func (source Source) Valid() bool {
	switch source {
	case SourcePlatformTemplate, SourceCorvusTemplate, SourceUser, SourceAgent:
		return true
	default:
		return false
	}
}

type RequirementLevel string

const (
	RequirementRequired    RequirementLevel = "required"
	RequirementRecommended RequirementLevel = "recommended"
)

func (level RequirementLevel) Valid() bool {
	return level == RequirementRequired || level == RequirementRecommended
}

type Category string

const (
	CategorySetup        Category = "setup"
	CategoryStoreCopy    Category = "store_copy"
	CategoryBranding     Category = "branding"
	CategoryMedia        Category = "media"
	CategoryCompliance   Category = "compliance"
	CategoryTimeline     Category = "timeline"
	CategoryPositioning  Category = "positioning"
	CategoryLocalization Category = "localization"
	CategoryReview       Category = "review"
)

func (category Category) Valid() bool {
	switch category {
	case CategorySetup, CategoryStoreCopy, CategoryBranding, CategoryMedia, CategoryCompliance,
		CategoryTimeline, CategoryPositioning, CategoryLocalization, CategoryReview:
		return true
	default:
		return false
	}
}

type Item struct {
	ID               ID
	ReleaseGoalID    releasedomain.ID
	Title            string
	Description      string
	Requirement      string
	Category         Category
	RequirementLevel RequirementLevel
	Source           Source
	SourceReference  string
	TemplateItemKey  *string
	TemplateKey      *string
	TemplateVersion  *string
	Status           Status
	SortOrder        int64
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type NewInput struct {
	ID               ID
	ReleaseGoalID    releasedomain.ID
	Title            string
	Description      string
	Requirement      string
	Category         Category
	RequirementLevel RequirementLevel
	Source           Source
	SourceReference  string
	TemplateItemKey  *string
	TemplateKey      *string
	TemplateVersion  *string
	SortOrder        int64
	Now              time.Time
}

func New(input NewInput) (Item, error) {
	now := input.Now.UTC()
	item := Item{
		ID:               input.ID,
		ReleaseGoalID:    input.ReleaseGoalID,
		Title:            strings.TrimSpace(input.Title),
		Description:      strings.TrimSpace(input.Description),
		Requirement:      strings.TrimSpace(input.Requirement),
		Category:         input.Category,
		RequirementLevel: input.RequirementLevel,
		Source:           input.Source,
		SourceReference:  strings.TrimSpace(input.SourceReference),
		TemplateItemKey:  trimPointer(input.TemplateItemKey),
		TemplateKey:      trimPointer(input.TemplateKey),
		TemplateVersion:  trimPointer(input.TemplateVersion),
		Status:           StatusNotStarted,
		SortOrder:        input.SortOrder,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
	if err := item.Validate(); err != nil {
		return Item{}, err
	}
	return item, nil
}

func (item Item) Validate() error {
	if _, err := ParseID(item.ID.String()); err != nil {
		return err
	}
	if _, err := releasedomain.ParseID(item.ReleaseGoalID.String()); err != nil {
		return NewValidationError("release_goal_id", "must be a UUIDv7")
	}
	if item.Title == "" || len([]rune(item.Title)) > MaxTitleLength {
		return NewValidationError("title", "must contain between 1 and 160 characters")
	}
	if len([]rune(item.Description)) > MaxDescriptionLength {
		return NewValidationError("description", "must not exceed 2000 characters")
	}
	if len([]rune(item.Requirement)) > MaxRequirementLength {
		return NewValidationError("requirement", "must not exceed 2000 characters")
	}
	if !item.Category.Valid() || !item.RequirementLevel.Valid() || !item.Source.Valid() || !item.Status.Valid() {
		return NewValidationError("classification", "category, requirement level, source and status must be supported")
	}
	if item.SourceReference == "" {
		return NewValidationError("source_reference", "must not be empty")
	}
	generated := item.Source == SourcePlatformTemplate || item.Source == SourceCorvusTemplate
	if generated && (blank(item.TemplateItemKey) || blank(item.TemplateKey) || blank(item.TemplateVersion)) {
		return NewValidationError("template", "generated items require template metadata")
	}
	if item.Source == SourceUser && (!blank(item.TemplateItemKey) || !blank(item.TemplateKey) || !blank(item.TemplateVersion)) {
		return NewValidationError("template", "user items must not have template metadata")
	}
	if item.RequirementLevel == RequirementRequired && item.Status == StatusNotApplicable {
		return NewValidationError("status", "required items cannot be not applicable")
	}
	if item.SortOrder < 1 {
		return NewValidationError("sort_order", "must be positive")
	}
	if item.CreatedAt.IsZero() || item.UpdatedAt.IsZero() {
		return NewValidationError("timestamps", "must not be zero")
	}
	return nil
}

func (item Item) Transition(to Status, releaseStatus releasedomain.Status, now time.Time) (Item, bool, error) {
	if !to.Valid() {
		return Item{}, false, NewValidationError("status", "must be a supported Checklist status")
	}
	if item.Status == to {
		return item, false, nil
	}
	if item.RequirementLevel == RequirementRequired && to == StatusNotApplicable {
		return Item{}, false, NewValidationError("status", "required items cannot be not applicable")
	}
	if item.RequirementLevel == RequirementRequired && releaseStatus.LocksRequiredChecklist() {
		return Item{}, false, fmt.Errorf("%w: move the Release Goal back to preparation first", ErrReleaseState)
	}
	if !CanTransition(item.Status, to) {
		return Item{}, false, fmt.Errorf("%w: %s to %s", ErrInvalidTransition, item.Status, to)
	}
	item.Status = to
	item.UpdatedAt = now.UTC()
	return item, true, nil
}

type Filter struct {
	Status   *Status
	Source   *Source
	Category *Category
}

func (filter Filter) Validate() error {
	if filter.Status != nil && !filter.Status.Valid() {
		return NewValidationError("status", "must be a supported Checklist status")
	}
	if filter.Source != nil && !filter.Source.Valid() {
		return NewValidationError("source", "must be a supported Checklist source")
	}
	if filter.Category != nil && !filter.Category.Valid() {
		return NewValidationError("category", "must be a supported Checklist category")
	}
	return nil
}

type ValidationError struct{ Field, Message string }

func NewValidationError(field string, message string) *ValidationError {
	return &ValidationError{Field: field, Message: message}
}
func (err *ValidationError) Error() string { return fmt.Sprintf("%s: %s", err.Field, err.Message) }
func (err *ValidationError) Unwrap() error { return ErrValidation }

func trimPointer(value *string) *string {
	if value == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*value)
	return &trimmed
}

func blank(value *string) bool { return value == nil || *value == "" }
