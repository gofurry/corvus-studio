package httpserver

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	legacyrouter "github.com/getkin/kin-openapi/routers/legacy"
	projectapp "github.com/gofurry/corvus-studio/apps/core/internal/application/project"
	projectdomain "github.com/gofurry/corvus-studio/apps/core/internal/domain/project"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func TestCreateProjectHandlerReturnsCreatedProject(t *testing.T) {
	entity := testProject(t)
	projects := &fakeProjectService{createResult: entity}
	server := newProjectTestServer(t, projects)
	body := `{
		"name":"Raven Game",
		"description":"Demo",
		"location":"C:\\Games\\Raven",
		"steam_app_id":480,
		"language":"English",
		"stage":"development"
	}`

	recorder := serveJSON(server, http.MethodPost, "/api/v1/projects", body)
	if recorder.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201; body=%s", recorder.Code, recorder.Body)
	}
	if recorder.Header().Get("Location") != "/api/v1/projects/"+entity.ID.String() {
		t.Fatalf("Location = %q", recorder.Header().Get("Location"))
	}
	if projects.createCommand.Name != "Raven Game" ||
		projects.createCommand.Stage != projectdomain.StageDevelopment ||
		projects.createCommand.SteamAppID == nil ||
		*projects.createCommand.SteamAppID != 480 {
		t.Fatalf("create command = %#v", projects.createCommand)
	}
	assertOpenAPIExchange(t, http.MethodPost, "/api/v1/projects", body, recorder)
}

func TestListAndGetProjectHandlers(t *testing.T) {
	entity := testProject(t)
	projects := &fakeProjectService{
		listResult: []projectdomain.Project{entity},
		getResult:  entity,
	}
	server := newProjectTestServer(t, projects)

	listRecorder := serveJSON(server, http.MethodGet, "/api/v1/projects", "")
	if listRecorder.Code != http.StatusOK {
		t.Fatalf("list status = %d, want 200", listRecorder.Code)
	}
	var listResponse struct {
		Items []json.RawMessage `json:"items"`
	}
	if err := json.Unmarshal(listRecorder.Body.Bytes(), &listResponse); err != nil {
		t.Fatalf("decode list: %v", err)
	}
	if len(listResponse.Items) != 1 {
		t.Fatalf("list count = %d, want 1", len(listResponse.Items))
	}
	assertOpenAPIExchange(t, http.MethodGet, "/api/v1/projects", "", listRecorder)

	path := "/api/v1/projects/" + entity.ID.String()
	getRecorder := serveJSON(server, http.MethodGet, path, "")
	if getRecorder.Code != http.StatusOK {
		t.Fatalf("get status = %d, want 200", getRecorder.Code)
	}
	if projects.getID != entity.ID {
		t.Fatalf("get ID = %q, want %q", projects.getID, entity.ID)
	}
	assertOpenAPIExchange(t, http.MethodGet, path, "", getRecorder)
}

func TestProjectHandlersMapErrors(t *testing.T) {
	validID := testProject(t).ID.String()
	tests := []struct {
		name       string
		method     string
		path       string
		body       string
		projects   *fakeProjectService
		wantStatus int
		wantCode   string
	}{
		{
			name:       "validation",
			method:     http.MethodPost,
			path:       "/api/v1/projects",
			body:       `{"name":"Raven","location":"C:\\Missing","language":"English","stage":"concept"}`,
			projects:   &fakeProjectService{createErr: projectdomain.NewValidationError("location", "must exist")},
			wantStatus: http.StatusBadRequest,
			wantCode:   "validation_failed",
		},
		{
			name:       "location conflict",
			method:     http.MethodPost,
			path:       "/api/v1/projects",
			body:       `{"name":"Raven","location":"C:\\Games\\Raven","language":"English","stage":"concept"}`,
			projects:   &fakeProjectService{createErr: projectdomain.ErrLocationConflict},
			wantStatus: http.StatusConflict,
			wantCode:   "project_location_conflict",
		},
		{
			name:       "not found",
			method:     http.MethodGet,
			path:       "/api/v1/projects/" + validID,
			projects:   &fakeProjectService{getErr: projectdomain.ErrNotFound},
			wantStatus: http.StatusNotFound,
			wantCode:   "project_not_found",
		},
		{
			name:       "internal",
			method:     http.MethodGet,
			path:       "/api/v1/projects",
			projects:   &fakeProjectService{listErr: errors.New("database unavailable")},
			wantStatus: http.StatusInternalServerError,
			wantCode:   "internal_error",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			server := newProjectTestServer(t, test.projects)
			recorder := serveJSON(server, test.method, test.path, test.body)
			assertErrorResponse(t, recorder, test.wantStatus, test.wantCode)
			assertOpenAPIExchange(t, test.method, test.path, test.body, recorder)
		})
	}
}

func TestProjectHandlersRejectMalformedInput(t *testing.T) {
	validBodyPrefix := `{"name":"Raven","location":"C:\\Games\\Raven","language":"English","stage":"concept"`
	tests := []struct {
		name   string
		method string
		path   string
		body   string
	}{
		{name: "unknown field", method: http.MethodPost, path: "/api/v1/projects", body: validBodyPrefix + `,"unknown":true}`},
		{name: "trailing JSON", method: http.MethodPost, path: "/api/v1/projects", body: validBodyPrefix + `} {}`},
		{name: "empty body", method: http.MethodPost, path: "/api/v1/projects"},
		{name: "invalid UUID", method: http.MethodGet, path: "/api/v1/projects/not-a-uuid"},
		{name: "non-v7 UUID", method: http.MethodGet, path: "/api/v1/projects/" + uuid.New().String()},
		{
			name:   "body too large",
			method: http.MethodPost,
			path:   "/api/v1/projects",
			body:   `{"name":"` + strings.Repeat("x", maxProjectRequestBody) + `"}`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			projects := &fakeProjectService{}
			server := newProjectTestServer(t, projects)
			recorder := serveJSON(server, test.method, test.path, test.body)
			assertErrorResponse(t, recorder, http.StatusBadRequest, "validation_failed")
			if projects.createCalls != 0 || projects.getCalls != 0 {
				t.Fatalf("service called for malformed input: create=%d get=%d", projects.createCalls, projects.getCalls)
			}
		})
	}
}

type fakeProjectService struct {
	createResult  projectdomain.Project
	createErr     error
	createCommand projectapp.CreateCommand
	createCalls   int
	listResult    []projectdomain.Project
	listErr       error
	getResult     projectdomain.Project
	getErr        error
	getID         projectdomain.ID
	getCalls      int
}

func (service *fakeProjectService) Create(
	_ context.Context,
	command projectapp.CreateCommand,
) (projectdomain.Project, error) {
	service.createCalls++
	service.createCommand = command
	return service.createResult, service.createErr
}

func (service *fakeProjectService) Get(
	_ context.Context,
	id projectdomain.ID,
) (projectdomain.Project, error) {
	service.getCalls++
	service.getID = id
	return service.getResult, service.getErr
}

func (service *fakeProjectService) List(context.Context) ([]projectdomain.Project, error) {
	if service.listResult == nil && service.listErr == nil {
		return []projectdomain.Project{}, nil
	}
	return service.listResult, service.listErr
}

func newProjectTestServer(t *testing.T, projects ProjectService) *Server {
	t.Helper()
	server, err := New(Dependencies{
		Health:   fakeHealth{version: 2},
		Projects: projects,
		Logger:   zap.NewNop(),
	})
	if err != nil {
		t.Fatalf("new server: %v", err)
	}
	return server
}

func serveJSON(server *Server, method string, path string, body string) *httptest.ResponseRecorder {
	var input *bytes.Reader
	if body == "" {
		input = bytes.NewReader(nil)
	} else {
		input = bytes.NewReader([]byte(body))
	}
	request := httptest.NewRequest(method, path, input)
	if body != "" {
		request.Header.Set("Content-Type", "application/json")
	}
	recorder := httptest.NewRecorder()
	server.Handler().ServeHTTP(recorder, request)
	return recorder
}

func assertErrorResponse(
	t *testing.T,
	recorder *httptest.ResponseRecorder,
	wantStatus int,
	wantCode string,
) {
	t.Helper()
	if recorder.Code != wantStatus {
		t.Fatalf("status = %d, want %d; body=%s", recorder.Code, wantStatus, recorder.Body)
	}
	var response struct {
		Error struct {
			Code        string `json:"code"`
			Message     string `json:"message"`
			Recoverable bool   `json:"recoverable"`
		} `json:"error"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode error response: %v", err)
	}
	if response.Error.Code != wantCode || response.Error.Message == "" {
		t.Fatalf("error response = %#v", response)
	}
	if wantStatus == http.StatusInternalServerError && response.Error.Recoverable {
		t.Fatal("internal error must not be recoverable")
	}
	if wantStatus != http.StatusInternalServerError && !response.Error.Recoverable {
		t.Fatal("client error must be recoverable")
	}
}

func assertOpenAPIExchange(
	t *testing.T,
	method string,
	path string,
	body string,
	recorder *httptest.ResponseRecorder,
) {
	t.Helper()
	loader := openapi3.NewLoader()
	document, err := loader.LoadFromFile(filepath.Join("..", "..", "..", "openapi", "openapi.yaml"))
	if err != nil {
		t.Fatalf("load OpenAPI: %v", err)
	}
	if err := document.Validate(t.Context()); err != nil {
		t.Fatalf("validate OpenAPI: %v", err)
	}
	router, err := legacyrouter.NewRouter(document)
	if err != nil {
		t.Fatalf("create OpenAPI router: %v", err)
	}
	request, err := http.NewRequestWithContext(
		t.Context(),
		method,
		"http://corvus.local"+path,
		strings.NewReader(body),
	)
	if err != nil {
		t.Fatalf("create contract request: %v", err)
	}
	if body != "" {
		request.Header.Set("Content-Type", "application/json")
	}
	route, pathParameters, err := router.FindRoute(request)
	if err != nil {
		t.Fatalf("find OpenAPI route: %v", err)
	}
	requestInput := &openapi3filter.RequestValidationInput{
		Request:    request,
		PathParams: pathParameters,
		Route:      route,
	}
	if err := openapi3filter.ValidateRequest(t.Context(), requestInput); err != nil {
		t.Fatalf("request violates OpenAPI: %v", err)
	}
	responseInput := (&openapi3filter.ResponseValidationInput{
		RequestValidationInput: requestInput,
		Status:                 recorder.Code,
		Header:                 recorder.Header(),
	}).SetBodyBytes(recorder.Body.Bytes())
	if err := openapi3filter.ValidateResponse(t.Context(), responseInput); err != nil {
		t.Fatalf("response violates OpenAPI: %v; body=%s", err, recorder.Body)
	}
}

func testProject(t *testing.T) projectdomain.Project {
	t.Helper()
	value, err := uuid.NewV7()
	if err != nil {
		t.Fatalf("generate UUIDv7: %v", err)
	}
	id, err := projectdomain.IDFromUUID(value)
	if err != nil {
		t.Fatalf("convert UUIDv7: %v", err)
	}
	steamAppID := int64(480)
	now := time.Date(2026, time.July, 31, 8, 0, 0, 0, time.UTC)
	return projectdomain.Project{
		ID:          id,
		Name:        "Raven Game",
		Description: "Demo",
		Location:    `C:\Games\Raven`,
		LocationKey: `c:\games\raven`,
		SteamAppID:  &steamAppID,
		Language:    "English",
		Stage:       projectdomain.StageDevelopment,
		Status:      projectdomain.StatusActive,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}
