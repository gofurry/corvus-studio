package httpserver

import (
	"context"
	"encoding/json"
	"errors"
	"io/fs"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/fstest"
	"time"

	checklistapp "github.com/gofurry/corvus-studio/apps/core/internal/application/checklist"
	projectapp "github.com/gofurry/corvus-studio/apps/core/internal/application/project"
	releaseapp "github.com/gofurry/corvus-studio/apps/core/internal/application/release"
	checklistdomain "github.com/gofurry/corvus-studio/apps/core/internal/domain/checklist"
	projectdomain "github.com/gofurry/corvus-studio/apps/core/internal/domain/project"
	releasedomain "github.com/gofurry/corvus-studio/apps/core/internal/domain/release"
	templatedomain "github.com/gofurry/corvus-studio/apps/core/internal/domain/template"
	"go.uber.org/zap"
)

type fakeHealth struct {
	err     error
	version int64
}

func (f fakeHealth) Ping(context.Context) error { return f.err }
func (f fakeHealth) SchemaVersion() int64       { return f.version }

type emptyProjects struct{}

func (emptyProjects) Create(context.Context, projectapp.CreateCommand) (projectdomain.Project, error) {
	return projectdomain.Project{}, errors.New("not implemented in this test")
}

func (emptyProjects) Get(context.Context, projectdomain.ID) (projectdomain.Project, error) {
	return projectdomain.Project{}, projectdomain.ErrNotFound
}

func (emptyProjects) List(context.Context) ([]projectdomain.Project, error) {
	return []projectdomain.Project{}, nil
}

type cancelledDirectoryPicker struct{}

func (cancelledDirectoryPicker) Select(context.Context) (string, bool, error) {
	return "", false, nil
}

type emptyReleases struct{}

func (emptyReleases) Create(context.Context, releaseapp.CreateCommand) (releasedomain.Goal, error) {
	return releasedomain.Goal{}, errors.New("not implemented in this test")
}
func (emptyReleases) Get(context.Context, releasedomain.ID) (releasedomain.Goal, error) {
	return releasedomain.Goal{}, releasedomain.ErrNotFound
}
func (emptyReleases) ListByProject(context.Context, projectdomain.ID) ([]releasedomain.Goal, error) {
	return []releasedomain.Goal{}, nil
}
func (emptyReleases) Transition(context.Context, releasedomain.ID, releasedomain.Status) (releasedomain.Goal, error) {
	return releasedomain.Goal{}, errors.New("not implemented in this test")
}
func (emptyReleases) LatestTemplate(string) (templatedomain.Definition, error) {
	return templatedomain.Definition{}, templatedomain.ErrNotFound
}

type emptyChecklist struct{}

func (emptyChecklist) Create(context.Context, checklistapp.CreateCommand) (checklistdomain.Item, error) {
	return checklistdomain.Item{}, errors.New("not implemented in this test")
}
func (emptyChecklist) Get(context.Context, checklistdomain.ID) (checklistdomain.Item, error) {
	return checklistdomain.Item{}, checklistdomain.ErrNotFound
}
func (emptyChecklist) List(context.Context, releasedomain.ID, checklistdomain.Filter) ([]checklistdomain.Item, error) {
	return []checklistdomain.Item{}, nil
}
func (emptyChecklist) Transition(context.Context, checklistdomain.ID, checklistdomain.Status) (checklistdomain.Item, error) {
	return checklistdomain.Item{}, errors.New("not implemented in this test")
}

func TestHealthHandlerReportsReady(t *testing.T) {
	server, err := New(testDependencies(fakeHealth{version: 1}, nil))
	if err != nil {
		t.Fatalf("new server: %v", err)
	}

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	server.Handler().ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", recorder.Code)
	}
	var response map[string]any
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if response["status"] != "ok" || response["database"] != "ok" || response["schema_version"] != float64(1) {
		t.Fatalf("unexpected health response: %#v", response)
	}
}

func TestHealthHandlerReportsUnavailable(t *testing.T) {
	server, err := New(testDependencies(fakeHealth{err: errors.New("database offline"), version: 1}, nil))
	if err != nil {
		t.Fatalf("new server: %v", err)
	}

	recorder := httptest.NewRecorder()
	server.Handler().ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/healthz", nil))
	if recorder.Code != http.StatusServiceUnavailable {
		t.Fatalf("status = %d, want 503", recorder.Code)
	}
}

func TestSPAHandlerServesAssetsAndFallback(t *testing.T) {
	assets := fstest.MapFS{
		"index.html":        &fstest.MapFile{Data: []byte("<h1>Corvus</h1>")},
		"assets/app.js":     &fstest.MapFile{Data: []byte("console.log('corvus')")},
		"assets/ignored.js": &fstest.MapFile{Data: []byte("ignored")},
	}
	server, err := New(testDependencies(fakeHealth{version: 1}, assets))
	if err != nil {
		t.Fatalf("new server: %v", err)
	}

	assetRecorder := httptest.NewRecorder()
	server.Handler().ServeHTTP(assetRecorder, httptest.NewRequest(http.MethodGet, "/assets/app.js", nil))
	if assetRecorder.Code != http.StatusOK || assetRecorder.Header().Get("Cache-Control") != "public, max-age=31536000, immutable" {
		t.Fatalf("asset response status=%d cache=%q", assetRecorder.Code, assetRecorder.Header().Get("Cache-Control"))
	}

	fallbackRecorder := httptest.NewRecorder()
	server.Handler().ServeHTTP(fallbackRecorder, httptest.NewRequest(http.MethodGet, "/nested/example", nil))
	if fallbackRecorder.Code != http.StatusOK || fallbackRecorder.Body.String() != "<h1>Corvus</h1>" {
		t.Fatalf("fallback response status=%d body=%q", fallbackRecorder.Code, fallbackRecorder.Body.String())
	}
	if fallbackRecorder.Header().Get("Cache-Control") != "no-cache" {
		t.Fatalf("fallback cache = %q", fallbackRecorder.Header().Get("Cache-Control"))
	}
}

func TestNewRejectsAssetsWithoutIndex(t *testing.T) {
	_, err := New(testDependencies(
		fakeHealth{},
		fstest.MapFS{"asset.js": &fstest.MapFile{Data: []byte("x")}},
	))
	if err == nil {
		t.Fatal("new server without index returned nil error")
	}
}

func TestServeStopsWhenContextIsCancelled(t *testing.T) {
	server, err := New(testDependencies(fakeHealth{version: 1}, nil))
	if err != nil {
		t.Fatalf("new server: %v", err)
	}
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- server.Serve(ctx, listener) }()

	waitForHealthy(t, "http://"+listener.Addr().String()+"/healthz")
	cancel()
	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("serve after cancel: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("server did not stop after cancellation")
	}
}

func waitForHealthy(t *testing.T, endpoint string) {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		response, err := http.Get(endpoint) //nolint:gosec,noctx
		if err == nil {
			_ = response.Body.Close()
			if response.StatusCode == http.StatusOK {
				return
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("health endpoint %s did not become ready", endpoint)
}

var _ fs.FS = fstest.MapFS{}

func testDependencies(health HealthChecker, assets fs.FS) Dependencies {
	return Dependencies{
		Health:          health,
		DirectoryPicker: cancelledDirectoryPicker{},
		Projects:        emptyProjects{},
		Releases:        emptyReleases{},
		Checklist:       emptyChecklist{},
		Logger:          zap.NewNop(),
		WebAssets:       assets,
	}
}
