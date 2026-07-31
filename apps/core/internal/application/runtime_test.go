package application

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/gofurry/corvus-studio/apps/core/internal/infrastructure/config"
)

func TestRunStartsHealthyRuntimeAndStops(t *testing.T) {
	port := reservePort(t)
	dataDir := t.TempDir()
	cfg := config.Config{
		Runtime: config.RuntimeConfig{DataDir: dataDir},
		Server:  config.ServerConfig{Host: config.DefaultHost, Port: port},
		Storage: config.StorageConfig{Path: filepath.Join(dataDir, "corvus.db")},
		Logging: config.LoggingConfig{
			Path:       filepath.Join(dataDir, "logs", "corvus.log"),
			Level:      "info",
			MaxSizeMB:  1,
			MaxBackups: 1,
			MaxAgeDays: 1,
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- Run(ctx, cfg) }()

	endpoint := "http://" + cfg.Server.Address() + "/healthz"
	response := waitForRuntimeHealth(t, endpoint)
	if response.Status != "ok" || response.Database != "ok" || response.SchemaVersion != 2 {
		t.Fatalf("unexpected health response: %#v", response)
	}
	cancel()
	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("runtime after cancellation: %v", err)
		}
	case <-time.After(10 * time.Second):
		t.Fatal("runtime did not stop after cancellation")
	}

	for _, path := range []string{cfg.Storage.Path, cfg.Logging.Path} {
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("expected runtime file %q: %v", path, err)
		}
	}
}

type healthResponse struct {
	Status        string `json:"status"`
	Database      string `json:"database"`
	SchemaVersion int64  `json:"schema_version"`
}

func waitForRuntimeHealth(t *testing.T, endpoint string) healthResponse {
	t.Helper()
	client := &http.Client{Timeout: time.Second}
	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		response, err := client.Get(endpoint) //nolint:noctx
		if err == nil {
			var health healthResponse
			decodeErr := json.NewDecoder(response.Body).Decode(&health)
			_ = response.Body.Close()
			if response.StatusCode == http.StatusOK && decodeErr == nil {
				return health
			}
		}
		time.Sleep(20 * time.Millisecond)
	}
	t.Fatalf("runtime health endpoint %s did not become ready", endpoint)
	return healthResponse{}
}

func reservePort(t *testing.T) int {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("reserve port: %v", err)
	}
	_, portText, err := net.SplitHostPort(listener.Addr().String())
	if err != nil {
		_ = listener.Close()
		t.Fatalf("split reserved address: %v", err)
	}
	if err := listener.Close(); err != nil {
		t.Fatalf("release reserved port: %v", err)
	}
	port, err := strconv.Atoi(portText)
	if err != nil {
		t.Fatalf("parse reserved port: %v", err)
	}
	return port
}
