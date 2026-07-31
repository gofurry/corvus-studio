package httpserver

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/gofurry/corvus-studio/apps/core/internal/infrastructure/directorypicker"
	"go.uber.org/zap"
)

func TestSelectDirectoryHandlerReturnsSelectionAndCancellation(t *testing.T) {
	tests := []struct {
		name         string
		picker       *fakeDirectoryPicker
		wantSelected bool
		wantPath     *string
	}{
		{
			name:         "selected",
			picker:       &fakeDirectoryPicker{path: `C:\Games\Raven`, selected: true},
			wantSelected: true,
			wantPath:     pointerTo(`C:\Games\Raven`),
		},
		{
			name:         "cancelled",
			picker:       &fakeDirectoryPicker{},
			wantSelected: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			server := newDirectoryPickerTestServer(t, test.picker)
			recorder := serveJSON(
				server,
				http.MethodPost,
				"/api/v1/system/select-directory",
				`{"purpose":"project_location"}`,
			)
			if recorder.Code != http.StatusOK {
				t.Fatalf("status = %d, want 200; body=%s", recorder.Code, recorder.Body)
			}
			var response struct {
				Path     *string `json:"path"`
				Selected bool    `json:"selected"`
			}
			if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
				t.Fatalf("decode response: %v", err)
			}
			if response.Selected != test.wantSelected || !equalOptionalString(response.Path, test.wantPath) {
				t.Fatalf("response = %#v, want selected=%t path=%v", response, test.wantSelected, test.wantPath)
			}
			if test.picker.calls != 1 {
				t.Fatalf("picker calls = %d, want 1", test.picker.calls)
			}
			assertOpenAPIExchange(
				t,
				http.MethodPost,
				"/api/v1/system/select-directory",
				`{"purpose":"project_location"}`,
				recorder,
			)
		})
	}
}

func TestSelectDirectoryHandlerRejectsMalformedRequests(t *testing.T) {
	tests := []string{
		"",
		`{"purpose":"unsupported"}`,
		`{"purpose":"project_location","unknown":true}`,
		`{"purpose":"project_location"} {}`,
	}
	for _, body := range tests {
		picker := &fakeDirectoryPicker{}
		server := newDirectoryPickerTestServer(t, picker)
		recorder := serveJSON(server, http.MethodPost, "/api/v1/system/select-directory", body)
		assertErrorResponse(t, recorder, http.StatusBadRequest, "validation_failed")
		if picker.calls != 0 {
			t.Fatalf("picker called for malformed body %q", body)
		}
	}
}

func TestSelectDirectoryHandlerReportsUnavailablePicker(t *testing.T) {
	picker := &fakeDirectoryPicker{err: directorypicker.ErrUnavailable}
	server := newDirectoryPickerTestServer(t, picker)
	body := `{"purpose":"project_location"}`
	recorder := serveJSON(server, http.MethodPost, "/api/v1/system/select-directory", body)

	assertErrorResponse(t, recorder, http.StatusServiceUnavailable, "directory_picker_unavailable")
	assertOpenAPIExchange(t, http.MethodPost, "/api/v1/system/select-directory", body, recorder)
}

func TestSelectDirectoryHandlerHidesUnexpectedErrors(t *testing.T) {
	picker := &fakeDirectoryPicker{err: errors.New("unexpected picker failure")}
	server := newDirectoryPickerTestServer(t, picker)
	recorder := serveJSON(
		server,
		http.MethodPost,
		"/api/v1/system/select-directory",
		`{"purpose":"project_location"}`,
	)

	assertErrorResponse(t, recorder, http.StatusInternalServerError, "internal_error")
}

type fakeDirectoryPicker struct {
	path     string
	selected bool
	err      error
	calls    int
}

func (picker *fakeDirectoryPicker) Select(context.Context) (string, bool, error) {
	picker.calls++
	return picker.path, picker.selected, picker.err
}

func newDirectoryPickerTestServer(t *testing.T, picker DirectoryPicker) *Server {
	t.Helper()
	server, err := New(Dependencies{
		Health:          fakeHealth{version: 2},
		DirectoryPicker: picker,
		Projects:        emptyProjects{},
		Logger:          zap.NewNop(),
	})
	if err != nil {
		t.Fatalf("new server: %v", err)
	}
	return server
}

func pointerTo(value string) *string {
	return &value
}

func equalOptionalString(left *string, right *string) bool {
	if left == nil || right == nil {
		return left == nil && right == nil
	}
	return *left == *right
}
