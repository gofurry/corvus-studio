package template

import (
	"errors"
	"fmt"
	"strings"
	"time"

	checklistdomain "github.com/gofurry/corvus-studio/apps/core/internal/domain/checklist"
	releasedomain "github.com/gofurry/corvus-studio/apps/core/internal/domain/release"
)

var (
	ErrValidation = errors.New("release template validation failed")
	ErrNotFound   = errors.New("release template not found")
)

type Item struct {
	Key              string                           `json:"key"`
	Title            string                           `json:"title"`
	Description      string                           `json:"description"`
	Requirement      string                           `json:"requirement"`
	Category         checklistdomain.Category         `json:"category"`
	RequirementLevel checklistdomain.RequirementLevel `json:"requirement_level"`
	Source           checklistdomain.Source           `json:"source"`
	SourceReference  string                           `json:"source_reference"`
	SortOrder        int64                            `json:"sort_order"`
}

type Definition struct {
	Key           string                 `json:"key"`
	Version       string                 `json:"version"`
	SchemaVersion int                    `json:"schema_version"`
	GoalType      releasedomain.GoalType `json:"goal_type"`
	Name          string                 `json:"name"`
	Description   string                 `json:"description"`
	ReviewedAt    string                 `json:"reviewed_at"`
	Items         []Item                 `json:"items"`
}

func (definition Definition) Validate() error {
	if strings.TrimSpace(definition.Key) == "" || strings.TrimSpace(definition.Version) == "" {
		return fmt.Errorf("%w: key and version are required", ErrValidation)
	}
	if definition.SchemaVersion != 1 {
		return fmt.Errorf("%w: schema_version must be 1", ErrValidation)
	}
	if !definition.GoalType.Valid() || strings.TrimSpace(definition.Name) == "" {
		return fmt.Errorf("%w: goal type and name are required", ErrValidation)
	}
	if _, err := time.Parse(time.DateOnly, definition.ReviewedAt); err != nil {
		return fmt.Errorf("%w: reviewed_at must be YYYY-MM-DD", ErrValidation)
	}
	if len(definition.Items) == 0 {
		return fmt.Errorf("%w: at least one item is required", ErrValidation)
	}
	keys := make(map[string]struct{}, len(definition.Items))
	orders := make(map[int64]struct{}, len(definition.Items))
	for _, item := range definition.Items {
		if strings.TrimSpace(item.Key) == "" || strings.TrimSpace(item.Title) == "" {
			return fmt.Errorf("%w: item key and title are required", ErrValidation)
		}
		if !item.Category.Valid() || !item.RequirementLevel.Valid() ||
			(item.Source != checklistdomain.SourcePlatformTemplate && item.Source != checklistdomain.SourceCorvusTemplate) {
			return fmt.Errorf("%w: item %q has invalid classification", ErrValidation, item.Key)
		}
		if strings.TrimSpace(item.SourceReference) == "" || item.SortOrder < 1 {
			return fmt.Errorf("%w: item %q has invalid source or order", ErrValidation, item.Key)
		}
		if _, exists := keys[item.Key]; exists {
			return fmt.Errorf("%w: duplicate item key %q", ErrValidation, item.Key)
		}
		if _, exists := orders[item.SortOrder]; exists {
			return fmt.Errorf("%w: duplicate item order %d", ErrValidation, item.SortOrder)
		}
		keys[item.Key] = struct{}{}
		orders[item.SortOrder] = struct{}{}
	}
	return nil
}
