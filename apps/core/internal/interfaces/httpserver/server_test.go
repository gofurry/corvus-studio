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

	"go.uber.org/zap"
)

type fakeHealth struct {
	err     error
	version int64
}

func (f fakeHealth) Ping(context.Context) error { return f.err }
func (f fakeHealth) SchemaVersion() int64       { return f.version }

func TestHealthHandlerReportsReady(t *testing.T) {
	server, err := New(fakeHealth{version: 1}, zap.NewNop(), nil)
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
	server, err := New(fakeHealth{err: errors.New("database offline"), version: 1}, zap.NewNop(), nil)
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
	server, err := New(fakeHealth{version: 1}, zap.NewNop(), assets)
	if err != nil {
		t.Fatalf("new server: %v", err)
	}

	assetRecorder := httptest.NewRecorder()
	server.Handler().ServeHTTP(assetRecorder, httptest.NewRequest(http.MethodGet, "/assets/app.js", nil))
	if assetRecorder.Code != http.StatusOK || assetRecorder.Header().Get("Cache-Control") != "public, max-age=31536000, immutable" {
		t.Fatalf("asset response status=%d cache=%q", assetRecorder.Code, assetRecorder.Header().Get("Cache-Control"))
	}

	fallbackRecorder := httptest.NewRecorder()
	server.Handler().ServeHTTP(fallbackRecorder, httptest.NewRequest(http.MethodGet, "/projects/example", nil))
	if fallbackRecorder.Code != http.StatusOK || fallbackRecorder.Body.String() != "<h1>Corvus</h1>" {
		t.Fatalf("fallback response status=%d body=%q", fallbackRecorder.Code, fallbackRecorder.Body.String())
	}
	if fallbackRecorder.Header().Get("Cache-Control") != "no-cache" {
		t.Fatalf("fallback cache = %q", fallbackRecorder.Header().Get("Cache-Control"))
	}
}

func TestNewRejectsAssetsWithoutIndex(t *testing.T) {
	_, err := New(fakeHealth{}, zap.NewNop(), fstest.MapFS{"asset.js": &fstest.MapFile{Data: []byte("x")}})
	if err == nil {
		t.Fatal("new server without index returned nil error")
	}
}

func TestServeStopsWhenContextIsCancelled(t *testing.T) {
	server, err := New(fakeHealth{version: 1}, zap.NewNop(), nil)
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
