package project

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestNewNormalizesAndValidatesProject(t *testing.T) {
	id := mustUUIDv7(t)
	steamAppID := int64(480)
	now := time.Date(2026, time.July, 31, 12, 30, 0, 0, time.FixedZone("test", 8*60*60))

	entity, err := New(NewInput{
		ID:          id,
		Name:        "  Raven Game  ",
		Description: "  First build  ",
		Location:    `C:\Games\Raven`,
		LocationKey: `c:\games\raven`,
		SteamAppID:  &steamAppID,
		Language:    "  English  ",
		Stage:       StageDevelopment,
		Now:         now,
	})
	if err != nil {
		t.Fatalf("new Project: %v", err)
	}
	if entity.Name != "Raven Game" || entity.Description != "First build" || entity.Language != "English" {
		t.Fatalf("text was not normalized: %#v", entity)
	}
	if entity.Status != StatusActive {
		t.Fatalf("status = %q, want active", entity.Status)
	}
	if entity.CreatedAt.Location() != time.UTC || entity.UpdatedAt.Location() != time.UTC {
		t.Fatalf("timestamps are not UTC: %#v", entity)
	}
	if entity.SteamAppID == nil || *entity.SteamAppID != steamAppID {
		t.Fatalf("Steam AppID = %#v", entity.SteamAppID)
	}
}

func TestNewAcceptsEverySupportedStage(t *testing.T) {
	for _, stage := range []Stage{
		StageConcept,
		StageDevelopment,
		StageReleasePreparation,
		StageReleased,
	} {
		t.Run(string(stage), func(t *testing.T) {
			_, err := New(validInput(t, stage))
			if err != nil {
				t.Fatalf("new Project at stage %q: %v", stage, err)
			}
		})
	}
}

func TestNewRejectsInvalidFields(t *testing.T) {
	validSteamAppID := int64(480)
	tests := []struct {
		name   string
		mutate func(*NewInput)
		field  string
	}{
		{name: "name empty", mutate: func(input *NewInput) { input.Name = "  " }, field: "name"},
		{name: "name too long", mutate: func(input *NewInput) { input.Name = strings.Repeat("n", 121) }, field: "name"},
		{name: "description too long", mutate: func(input *NewInput) { input.Description = strings.Repeat("d", 2001) }, field: "description"},
		{name: "location empty", mutate: func(input *NewInput) { input.Location = "" }, field: "location"},
		{name: "location key empty", mutate: func(input *NewInput) { input.LocationKey = "" }, field: "location"},
		{name: "Steam AppID zero", mutate: func(input *NewInput) { value := int64(0); input.SteamAppID = &value }, field: "steam_app_id"},
		{name: "Steam AppID too large", mutate: func(input *NewInput) { value := MaxSteamAppID + 1; input.SteamAppID = &value }, field: "steam_app_id"},
		{name: "language empty", mutate: func(input *NewInput) { input.Language = "" }, field: "language"},
		{name: "language too long", mutate: func(input *NewInput) { input.Language = strings.Repeat("l", 65) }, field: "language"},
		{name: "stage invalid", mutate: func(input *NewInput) { input.Stage = "unknown" }, field: "stage"},
		{name: "time zero", mutate: func(input *NewInput) { input.Now = time.Time{} }, field: "timestamps"},
		{name: "ID invalid", mutate: func(input *NewInput) { input.ID = "not-a-uuid" }, field: "project_id"},
		{name: "valid optional Steam AppID", mutate: func(input *NewInput) { input.SteamAppID = &validSteamAppID }, field: ""},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			input := validInput(t, StageConcept)
			test.mutate(&input)
			_, err := New(input)
			if test.field == "" {
				if err != nil {
					t.Fatalf("unexpected validation error: %v", err)
				}
				return
			}
			if !errors.Is(err, ErrValidation) {
				t.Fatalf("error = %v, want validation error", err)
			}
			var validationError *ValidationError
			if !errors.As(err, &validationError) || validationError.Field != test.field {
				t.Fatalf("error = %#v, want field %q", err, test.field)
			}
		})
	}
}

func TestParseIDRequiresUUIDv7(t *testing.T) {
	if _, err := ParseID(uuid.New().String()); !errors.Is(err, ErrValidation) {
		t.Fatalf("UUIDv4 error = %v, want validation error", err)
	}
	id := mustUUIDv7(t)
	parsed, err := ParseID(id.String())
	if err != nil {
		t.Fatalf("parse UUIDv7: %v", err)
	}
	if parsed != id {
		t.Fatalf("parsed ID = %q, want %q", parsed, id)
	}
}

func validInput(t *testing.T, stage Stage) NewInput {
	t.Helper()
	return NewInput{
		ID:          mustUUIDv7(t),
		Name:        "Raven Game",
		Description: "",
		Location:    "/games/raven",
		LocationKey: "/games/raven",
		Language:    "English",
		Stage:       stage,
		Now:         time.Date(2026, time.July, 31, 4, 0, 0, 0, time.UTC),
	}
}

func mustUUIDv7(t *testing.T) ID {
	t.Helper()
	value, err := uuid.NewV7()
	if err != nil {
		t.Fatalf("generate UUIDv7: %v", err)
	}
	id, err := IDFromUUID(value)
	if err != nil {
		t.Fatalf("convert UUIDv7: %v", err)
	}
	return id
}
