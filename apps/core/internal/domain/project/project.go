package project

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	MaxNameLength        = 120
	MaxDescriptionLength = 2000
	MaxLanguageLength    = 64
	MaxSteamAppID        = int64(4294967295)
)

var (
	ErrValidation       = errors.New("project validation failed")
	ErrNotFound         = errors.New("project not found")
	ErrLocationConflict = errors.New("project location conflict")
)

type ID string

func ParseID(value string) (ID, error) {
	parsed, err := uuid.Parse(value)
	if err != nil || parsed.Version() != 7 {
		return "", NewValidationError("project_id", "must be a UUIDv7")
	}
	return ID(parsed.String()), nil
}

func IDFromUUID(value uuid.UUID) (ID, error) {
	return ParseID(value.String())
}

func (id ID) String() string {
	return string(id)
}

type Stage string

const (
	StageConcept            Stage = "concept"
	StageDevelopment        Stage = "development"
	StageReleasePreparation Stage = "release_preparation"
	StageReleased           Stage = "released"
)

func (stage Stage) Valid() bool {
	switch stage {
	case StageConcept, StageDevelopment, StageReleasePreparation, StageReleased:
		return true
	default:
		return false
	}
}

type Status string

const (
	StatusActive   Status = "active"
	StatusArchived Status = "archived"
)

func (status Status) Valid() bool {
	return status == StatusActive || status == StatusArchived
}

type Project struct {
	ID          ID
	Name        string
	Description string
	Location    string
	LocationKey string
	SteamAppID  *int64
	Language    string
	Stage       Stage
	Status      Status
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type NewInput struct {
	ID          ID
	Name        string
	Description string
	Location    string
	LocationKey string
	SteamAppID  *int64
	Language    string
	Stage       Stage
	Now         time.Time
}

func New(input NewInput) (Project, error) {
	createdAt := input.Now.UTC()
	entity := Project{
		ID:          input.ID,
		Name:        strings.TrimSpace(input.Name),
		Description: strings.TrimSpace(input.Description),
		Location:    input.Location,
		LocationKey: input.LocationKey,
		SteamAppID:  cloneInt64(input.SteamAppID),
		Language:    strings.TrimSpace(input.Language),
		Stage:       input.Stage,
		Status:      StatusActive,
		CreatedAt:   createdAt,
		UpdatedAt:   createdAt,
	}
	if err := entity.Validate(); err != nil {
		return Project{}, err
	}
	return entity, nil
}

func (project Project) Validate() error {
	if _, err := ParseID(project.ID.String()); err != nil {
		return err
	}
	if project.Name == "" || len([]rune(project.Name)) > MaxNameLength {
		return NewValidationError("name", "must contain between 1 and 120 characters")
	}
	if len([]rune(project.Description)) > MaxDescriptionLength {
		return NewValidationError("description", "must not exceed 2000 characters")
	}
	if project.Location == "" {
		return NewValidationError("location", "must reference an existing absolute directory")
	}
	if project.LocationKey == "" {
		return NewValidationError("location", "could not be normalized")
	}
	if project.SteamAppID != nil && (*project.SteamAppID < 1 || *project.SteamAppID > MaxSteamAppID) {
		return NewValidationError("steam_app_id", "must be a positive 32-bit integer")
	}
	if project.Language == "" || len([]rune(project.Language)) > MaxLanguageLength {
		return NewValidationError("language", "must contain between 1 and 64 characters")
	}
	if !project.Stage.Valid() {
		return NewValidationError("stage", "must be a supported Project stage")
	}
	if !project.Status.Valid() {
		return NewValidationError("status", "must be a supported Project status")
	}
	if project.CreatedAt.IsZero() || project.UpdatedAt.IsZero() {
		return NewValidationError("timestamps", "must not be zero")
	}
	return nil
}

type ValidationError struct {
	Field   string
	Message string
}

func NewValidationError(field string, message string) *ValidationError {
	return &ValidationError{Field: field, Message: message}
}

func (err *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", err.Field, err.Message)
}

func (err *ValidationError) Unwrap() error {
	return ErrValidation
}

func cloneInt64(value *int64) *int64 {
	if value == nil {
		return nil
	}
	copy := *value
	return &copy
}
