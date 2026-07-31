package application

import (
	"bytes"
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
	if response.Status != "ok" || response.Database != "ok" || response.SchemaVersion != 3 {
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

func TestRunPersistsProjectAndReleaseWorkspaceAcrossRestart(t *testing.T) {
	dataDir := t.TempDir()
	projectDirectory := filepath.Join(dataDir, "game")
	if err := os.Mkdir(projectDirectory, 0o700); err != nil {
		t.Fatalf("create Project directory: %v", err)
	}
	port := reservePort(t)
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
	baseURL := "http://" + cfg.Server.Address()

	cancel, done := startRuntime(cfg)
	waitForRuntimeHealth(t, baseURL+"/healthz")
	requestBody, err := json.Marshal(map[string]any{
		"name":        "Persistent Raven",
		"description": "Created before restart",
		"location":    projectDirectory,
		"language":    "English",
		"stage":       "concept",
	})
	if err != nil {
		t.Fatalf("encode create request: %v", err)
	}
	response, err := http.Post( //nolint:noctx
		baseURL+"/api/v1/projects",
		"application/json",
		bytes.NewReader(requestBody),
	)
	if err != nil {
		t.Fatalf("create Project: %v", err)
	}
	var created struct {
		ID string `json:"id"`
	}
	decodeErr := json.NewDecoder(response.Body).Decode(&created)
	_ = response.Body.Close()
	if response.StatusCode != http.StatusCreated || decodeErr != nil || created.ID == "" {
		t.Fatalf(
			"create response status=%d decode=%v ID=%q",
			response.StatusCode,
			decodeErr,
			created.ID,
		)
	}
	releaseRequestBody, err := json.Marshal(map[string]any{
		"project_id":       created.ID,
		"goal_type":        "steam_coming_soon",
		"template_key":     "steam-coming-soon",
		"template_version": "1.0.0",
	})
	if err != nil {
		t.Fatal(err)
	}
	response, err = http.Post( //nolint:noctx
		baseURL+"/api/v1/releases",
		"application/json",
		bytes.NewReader(releaseRequestBody),
	)
	if err != nil {
		t.Fatalf("create Release workspace: %v", err)
	}
	var createdRelease struct {
		ID               string `json:"id"`
		ChecklistSummary struct {
			Total int `json:"total"`
		} `json:"checklist_summary"`
	}
	decodeErr = json.NewDecoder(response.Body).Decode(&createdRelease)
	_ = response.Body.Close()
	if response.StatusCode != http.StatusCreated || decodeErr != nil || createdRelease.ID == "" || createdRelease.ChecklistSummary.Total != 12 {
		t.Fatalf("create Release response status=%d decode=%v value=%#v", response.StatusCode, decodeErr, createdRelease)
	}
	response, err = http.Get(baseURL + "/api/v1/releases/" + createdRelease.ID + "/checklist") //nolint:noctx
	if err != nil {
		t.Fatalf("list Checklist: %v", err)
	}
	var checklistResponse struct {
		Items []struct {
			ID string `json:"id"`
		} `json:"items"`
	}
	decodeErr = json.NewDecoder(response.Body).Decode(&checklistResponse)
	_ = response.Body.Close()
	if response.StatusCode != http.StatusOK || decodeErr != nil || len(checklistResponse.Items) != 12 {
		t.Fatalf("Checklist response status=%d decode=%v value=%#v", response.StatusCode, decodeErr, checklistResponse)
	}
	transitionBody := bytes.NewBufferString(`{"status":"in_progress"}`)
	response, err = http.Post( //nolint:noctx
		baseURL+"/api/v1/checklist/"+checklistResponse.Items[0].ID+"/transition",
		"application/json",
		transitionBody,
	)
	if err != nil {
		t.Fatalf("transition Checklist item: %v", err)
	}
	_ = response.Body.Close()
	if response.StatusCode != http.StatusOK {
		t.Fatalf("transition Checklist status=%d", response.StatusCode)
	}
	stopRuntime(t, cancel, done)

	cancel, done = startRuntime(cfg)
	defer stopRuntime(t, cancel, done)
	waitForRuntimeHealth(t, baseURL+"/healthz")
	response, err = http.Get(baseURL + "/api/v1/projects/" + created.ID) //nolint:noctx
	if err != nil {
		t.Fatalf("get Project after restart: %v", err)
	}
	var reopened struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	decodeErr = json.NewDecoder(response.Body).Decode(&reopened)
	_ = response.Body.Close()
	if response.StatusCode != http.StatusOK || decodeErr != nil {
		t.Fatalf("reopen response status=%d decode=%v", response.StatusCode, decodeErr)
	}
	if reopened.ID != created.ID || reopened.Name != "Persistent Raven" {
		t.Fatalf("reopened Project = %#v", reopened)
	}
	response, err = http.Get(baseURL + "/api/v1/releases/" + createdRelease.ID) //nolint:noctx
	if err != nil {
		t.Fatalf("get Release after restart: %v", err)
	}
	var reopenedRelease struct {
		ID               string `json:"id"`
		ChecklistSummary struct {
			Total int `json:"total"`
		} `json:"checklist_summary"`
	}
	decodeErr = json.NewDecoder(response.Body).Decode(&reopenedRelease)
	_ = response.Body.Close()
	if response.StatusCode != http.StatusOK || decodeErr != nil || reopenedRelease.ID != createdRelease.ID || reopenedRelease.ChecklistSummary.Total != 12 {
		t.Fatalf("reopened Release status=%d decode=%v value=%#v", response.StatusCode, decodeErr, reopenedRelease)
	}
	response, err = http.Get(baseURL + "/api/v1/releases/" + createdRelease.ID + "/checklist?status=in_progress") //nolint:noctx
	if err != nil {
		t.Fatalf("get transitioned Checklist after restart: %v", err)
	}
	decodeErr = json.NewDecoder(response.Body).Decode(&checklistResponse)
	_ = response.Body.Close()
	if response.StatusCode != http.StatusOK || decodeErr != nil || len(checklistResponse.Items) != 1 || checklistResponse.Items[0].ID == "" {
		t.Fatalf("reopened Checklist status=%d decode=%v value=%#v", response.StatusCode, decodeErr, checklistResponse)
	}
}

type healthResponse struct {
	Status        string `json:"status"`
	Database      string `json:"database"`
	SchemaVersion int64  `json:"schema_version"`
}

func startRuntime(cfg config.Config) (context.CancelFunc, <-chan error) {
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- Run(ctx, cfg) }()
	return cancel, done
}

func stopRuntime(t *testing.T, cancel context.CancelFunc, done <-chan error) {
	t.Helper()
	cancel()
	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("runtime after cancellation: %v", err)
		}
	case <-time.After(10 * time.Second):
		t.Fatal("runtime did not stop after cancellation")
	}
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
