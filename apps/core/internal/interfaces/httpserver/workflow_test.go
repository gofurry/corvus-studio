package httpserver

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
	"time"

	checklistapp "github.com/gofurry/corvus-studio/apps/core/internal/application/checklist"
	releaseapp "github.com/gofurry/corvus-studio/apps/core/internal/application/release"
	checklistdomain "github.com/gofurry/corvus-studio/apps/core/internal/domain/checklist"
	projectdomain "github.com/gofurry/corvus-studio/apps/core/internal/domain/project"
	releasedomain "github.com/gofurry/corvus-studio/apps/core/internal/domain/release"
	templatedomain "github.com/gofurry/corvus-studio/apps/core/internal/domain/template"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func TestReleaseHandlersFollowOpenAPIContract(t *testing.T) {
	goal := workflowGoal(t)
	template := workflowTemplate()
	releases := &fakeReleaseService{
		createResult: goal, getResult: goal, listResult: []releasedomain.Goal{goal},
		transitionResult: withReleaseStatus(goal, releasedomain.StatusPreparing), templateResult: template,
	}
	server := newWorkflowTestServer(t, releases, &fakeChecklistService{})

	templateRecorder := serveJSON(server, http.MethodGet, "/api/v1/release-templates/steam-coming-soon", "")
	assertStatus(t, templateRecorder.Code, http.StatusOK, templateRecorder.Body.String())
	assertOpenAPIExchange(t, http.MethodGet, "/api/v1/release-templates/steam-coming-soon", "", templateRecorder)

	createBody := `{"project_id":"` + goal.ProjectID.String() + `","goal_type":"steam_coming_soon","template_key":"steam-coming-soon","template_version":"1.0.0"}`
	createRecorder := serveJSON(server, http.MethodPost, "/api/v1/releases", createBody)
	assertStatus(t, createRecorder.Code, http.StatusCreated, createRecorder.Body.String())
	if createRecorder.Header().Get("Location") != "/api/v1/releases/"+goal.ID.String() {
		t.Fatalf("Location = %q", createRecorder.Header().Get("Location"))
	}
	if releases.createCommand.ProjectID != goal.ProjectID || releases.createCommand.TemplateVersion != "1.0.0" {
		t.Fatalf("create command = %#v", releases.createCommand)
	}
	assertOpenAPIExchange(t, http.MethodPost, "/api/v1/releases", createBody, createRecorder)

	listPath := "/api/v1/projects/" + goal.ProjectID.String() + "/releases"
	listRecorder := serveJSON(server, http.MethodGet, listPath, "")
	assertStatus(t, listRecorder.Code, http.StatusOK, listRecorder.Body.String())
	assertOpenAPIExchange(t, http.MethodGet, listPath, "", listRecorder)

	getPath := "/api/v1/releases/" + goal.ID.String()
	getRecorder := serveJSON(server, http.MethodGet, getPath, "")
	assertStatus(t, getRecorder.Code, http.StatusOK, getRecorder.Body.String())
	assertOpenAPIExchange(t, http.MethodGet, getPath, "", getRecorder)

	transitionPath := getPath + "/transition"
	transitionBody := `{"status":"preparing"}`
	transitionRecorder := serveJSON(server, http.MethodPost, transitionPath, transitionBody)
	assertStatus(t, transitionRecorder.Code, http.StatusOK, transitionRecorder.Body.String())
	if releases.transitionStatus != releasedomain.StatusPreparing {
		t.Fatalf("transition status = %q", releases.transitionStatus)
	}
	assertOpenAPIExchange(t, http.MethodPost, transitionPath, transitionBody, transitionRecorder)
}

func TestChecklistHandlersFollowOpenAPIContract(t *testing.T) {
	goal := workflowGoal(t)
	item := workflowItem(t, goal.ID)
	checklist := &fakeChecklistService{createResult: item, getResult: item, listResult: []checklistdomain.Item{item}, transitionResult: withChecklistStatus(item, checklistdomain.StatusInProgress)}
	server := newWorkflowTestServer(t, &fakeReleaseService{getResult: goal}, checklist)

	listPath := "/api/v1/releases/" + goal.ID.String() + "/checklist?status=not_started&source=user&category=review"
	listRecorder := serveJSON(server, http.MethodGet, listPath, "")
	assertStatus(t, listRecorder.Code, http.StatusOK, listRecorder.Body.String())
	if checklist.listFilter.Status == nil || *checklist.listFilter.Status != checklistdomain.StatusNotStarted ||
		checklist.listFilter.Source == nil || *checklist.listFilter.Source != checklistdomain.SourceUser ||
		checklist.listFilter.Category == nil || *checklist.listFilter.Category != checklistdomain.CategoryReview {
		t.Fatalf("list filter = %#v", checklist.listFilter)
	}
	assertOpenAPIExchange(t, http.MethodGet, listPath, "", listRecorder)

	createBody := `{"release_goal_id":"` + goal.ID.String() + `","title":"User review","description":"Ask another developer","requirement":"Record the result","category":"review","requirement_level":"recommended"}`
	createRecorder := serveJSON(server, http.MethodPost, "/api/v1/checklist", createBody)
	assertStatus(t, createRecorder.Code, http.StatusCreated, createRecorder.Body.String())
	if createRecorder.Header().Get("Location") != "/api/v1/checklist/"+item.ID.String() {
		t.Fatalf("Location = %q", createRecorder.Header().Get("Location"))
	}
	assertOpenAPIExchange(t, http.MethodPost, "/api/v1/checklist", createBody, createRecorder)

	getPath := "/api/v1/checklist/" + item.ID.String()
	getRecorder := serveJSON(server, http.MethodGet, getPath, "")
	assertStatus(t, getRecorder.Code, http.StatusOK, getRecorder.Body.String())
	assertOpenAPIExchange(t, http.MethodGet, getPath, "", getRecorder)

	transitionPath := getPath + "/transition"
	transitionBody := `{"status":"in_progress"}`
	transitionRecorder := serveJSON(server, http.MethodPost, transitionPath, transitionBody)
	assertStatus(t, transitionRecorder.Code, http.StatusOK, transitionRecorder.Body.String())
	assertOpenAPIExchange(t, http.MethodPost, transitionPath, transitionBody, transitionRecorder)
}

func TestWorkflowErrorMapping(t *testing.T) {
	goal := workflowGoal(t)
	tests := []struct {
		name       string
		err        error
		wantStatus int
		wantCode   string
	}{
		{name: "template missing", err: templatedomain.ErrNotFound, wantStatus: 404, wantCode: "template_not_found"},
		{name: "release missing", err: releasedomain.ErrNotFound, wantStatus: 404, wantCode: "release_not_found"},
		{name: "release conflict", err: releasedomain.ErrConflict, wantStatus: 409, wantCode: "release_goal_conflict"},
		{name: "release not ready", err: releasedomain.ErrNotReady, wantStatus: 409, wantCode: "release_not_ready"},
		{name: "invalid transition", err: checklistdomain.ErrInvalidTransition, wantStatus: 409, wantCode: "invalid_transition"},
		{name: "release state", err: checklistdomain.ErrReleaseState, wantStatus: 409, wantCode: "release_state_conflict"},
		{name: "internal", err: errors.New("database unavailable"), wantStatus: 500, wantCode: "internal_error"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			releases := &fakeReleaseService{getErr: test.err}
			server := newWorkflowTestServer(t, releases, &fakeChecklistService{})
			path := "/api/v1/releases/" + goal.ID.String()
			recorder := serveJSON(server, http.MethodGet, path, "")
			assertErrorResponse(t, recorder, test.wantStatus, test.wantCode)
			assertOpenAPIExchange(t, http.MethodGet, path, "", recorder)
		})
	}
}

func TestWorkflowHandlersRejectMalformedInput(t *testing.T) {
	goal := workflowGoal(t)
	server := newWorkflowTestServer(t, &fakeReleaseService{}, &fakeChecklistService{})
	tests := []struct{ method, path, body string }{
		{http.MethodPost, "/api/v1/releases", `{"project_id":"` + goal.ProjectID.String() + `","goal_type":"steam_coming_soon","template_key":"steam-coming-soon","template_version":"1.0.0","unknown":true}`},
		{http.MethodPost, "/api/v1/releases/" + goal.ID.String() + "/transition", `{"status":"preparing"} {}`},
		{http.MethodGet, "/api/v1/releases/not-a-uuid", ""},
		{http.MethodGet, "/api/v1/checklist/" + uuid.New().String(), ""},
	}
	for _, test := range tests {
		recorder := serveJSON(server, test.method, test.path, test.body)
		assertErrorResponse(t, recorder, http.StatusBadRequest, "validation_failed")
	}
}

type fakeReleaseService struct {
	createResult     releasedomain.Goal
	createErr        error
	createCommand    releaseapp.CreateCommand
	getResult        releasedomain.Goal
	getErr           error
	listResult       []releasedomain.Goal
	listErr          error
	transitionResult releasedomain.Goal
	transitionErr    error
	transitionStatus releasedomain.Status
	templateResult   templatedomain.Definition
	templateErr      error
}

func (service *fakeReleaseService) Create(_ context.Context, command releaseapp.CreateCommand) (releasedomain.Goal, error) {
	service.createCommand = command
	return service.createResult, service.createErr
}
func (service *fakeReleaseService) Get(context.Context, releasedomain.ID) (releasedomain.Goal, error) {
	return service.getResult, service.getErr
}
func (service *fakeReleaseService) ListByProject(context.Context, projectdomain.ID) ([]releasedomain.Goal, error) {
	if service.listResult == nil && service.listErr == nil {
		return []releasedomain.Goal{}, nil
	}
	return service.listResult, service.listErr
}
func (service *fakeReleaseService) Transition(_ context.Context, _ releasedomain.ID, status releasedomain.Status) (releasedomain.Goal, error) {
	service.transitionStatus = status
	return service.transitionResult, service.transitionErr
}
func (service *fakeReleaseService) LatestTemplate(string) (templatedomain.Definition, error) {
	return service.templateResult, service.templateErr
}

type fakeChecklistService struct {
	createResult     checklistdomain.Item
	createErr        error
	getResult        checklistdomain.Item
	getErr           error
	listResult       []checklistdomain.Item
	listErr          error
	listFilter       checklistdomain.Filter
	transitionResult checklistdomain.Item
	transitionErr    error
}

func (service *fakeChecklistService) Create(context.Context, checklistapp.CreateCommand) (checklistdomain.Item, error) {
	return service.createResult, service.createErr
}
func (service *fakeChecklistService) Get(context.Context, checklistdomain.ID) (checklistdomain.Item, error) {
	return service.getResult, service.getErr
}
func (service *fakeChecklistService) List(_ context.Context, _ releasedomain.ID, filter checklistdomain.Filter) ([]checklistdomain.Item, error) {
	service.listFilter = filter
	if service.listResult == nil && service.listErr == nil {
		return []checklistdomain.Item{}, nil
	}
	return service.listResult, service.listErr
}
func (service *fakeChecklistService) Transition(context.Context, checklistdomain.ID, checklistdomain.Status) (checklistdomain.Item, error) {
	return service.transitionResult, service.transitionErr
}

func newWorkflowTestServer(t *testing.T, releases ReleaseService, checklist ChecklistService) *Server {
	t.Helper()
	server, err := New(Dependencies{
		Health: fakeHealth{version: 3}, DirectoryPicker: cancelledDirectoryPicker{}, Projects: emptyProjects{},
		Releases: releases, Checklist: checklist, Logger: zap.NewNop(),
	})
	if err != nil {
		t.Fatal(err)
	}
	return server
}

func workflowGoal(t *testing.T) releasedomain.Goal {
	t.Helper()
	releaseID, _ := releasedomain.IDFromUUID(mustWorkflowUUID(t))
	projectID, _ := projectdomain.IDFromUUID(mustWorkflowUUID(t))
	now := time.Date(2026, time.July, 31, 12, 0, 0, 0, time.UTC)
	goal, err := releasedomain.New(releasedomain.NewInput{
		ID: releaseID, ProjectID: projectID, GoalType: releasedomain.GoalTypeSteamComingSoon,
		Title: "Steam Coming Soon", TemplateKey: "steam-coming-soon", TemplateVersion: "1.0.0", Now: now,
	})
	if err != nil {
		t.Fatal(err)
	}
	goal.Summary = releasedomain.ChecklistSummary{Total: 12, Done: 1, RequiredTotal: 7, RequiredDone: 1}
	return goal
}

func workflowItem(t *testing.T, releaseID releasedomain.ID) checklistdomain.Item {
	t.Helper()
	id, _ := checklistdomain.IDFromUUID(mustWorkflowUUID(t))
	item, err := checklistdomain.New(checklistdomain.NewInput{
		ID: id, ReleaseGoalID: releaseID, Title: "User review", Description: "Ask another developer",
		Requirement: "Record the result", Category: checklistdomain.CategoryReview,
		RequirementLevel: checklistdomain.RequirementRecommended, Source: checklistdomain.SourceUser,
		SourceReference: "User-created task", SortOrder: 13, Now: time.Date(2026, time.July, 31, 12, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatal(err)
	}
	return item
}

func workflowTemplate() templatedomain.Definition {
	return templatedomain.Definition{
		Key: "steam-coming-soon", Version: "1.0.0", SchemaVersion: 1,
		GoalType: releasedomain.GoalTypeSteamComingSoon, Name: "Steam Coming Soon",
		Description: "Prepare a store page", ReviewedAt: "2026-07-31",
		Items: []templatedomain.Item{{
			Key: "store-review", Title: "Review store page", Description: "Check the page",
			Requirement: "Complete review", Category: checklistdomain.CategoryReview,
			RequirementLevel: checklistdomain.RequirementRecommended, Source: checklistdomain.SourceCorvusTemplate,
			SourceReference: "docs/product/Corvus_Studio_Launch_PRD_v0.1.md", SortOrder: 1,
		}},
	}
}

func withReleaseStatus(goal releasedomain.Goal, status releasedomain.Status) releasedomain.Goal {
	goal.Status = status
	return goal
}

func withChecklistStatus(item checklistdomain.Item, status checklistdomain.Status) checklistdomain.Item {
	item.Status = status
	return item
}

func mustWorkflowUUID(t *testing.T) uuid.UUID {
	t.Helper()
	value, err := uuid.NewV7()
	if err != nil {
		t.Fatal(err)
	}
	return value
}

func assertStatus(t *testing.T, got int, want int, body string) {
	t.Helper()
	if got != want {
		var decoded any
		_ = json.Unmarshal([]byte(body), &decoded)
		t.Fatalf("status = %d, want %d; body=%v", got, want, decoded)
	}
}
