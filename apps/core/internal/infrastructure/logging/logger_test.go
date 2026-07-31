package logging

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/gofurry/corvus-studio/apps/core/internal/infrastructure/config"
	"go.uber.org/zap"
)

type failingSyncWriter struct {
	bytes.Buffer
	syncCalls int
}

func (w *failingSyncWriter) Sync() error {
	w.syncCalls++
	return errors.New("console sync must not be called")
}

func TestNewWritesStructuredFileLog(t *testing.T) {
	logPath := filepath.Join(t.TempDir(), "nested", "corvus.log")
	handle, err := New(config.LoggingConfig{
		Path:       logPath,
		Level:      "info",
		MaxSizeMB:  1,
		MaxBackups: 1,
		MaxAgeDays: 1,
		Compress:   true,
	})
	if err != nil {
		t.Fatalf("new logger: %v", err)
	}

	handle.Logger.Info("runtime ready", zap.String("component", "test"))
	if err := handle.Close(); err != nil {
		t.Fatalf("close logger: %v", err)
	}

	content, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("read log: %v", err)
	}
	var entry map[string]any
	if err := json.Unmarshal(content, &entry); err != nil {
		t.Fatalf("decode JSON log %q: %v", string(content), err)
	}
	if entry["msg"] != "runtime ready" || entry["component"] != "test" {
		t.Fatalf("unexpected log entry: %#v", entry)
	}
}

func TestCloseDoesNotSyncConsole(t *testing.T) {
	console := &failingSyncWriter{}
	handle, err := newWithConsole(config.LoggingConfig{
		Path:       filepath.Join(t.TempDir(), "corvus.log"),
		Level:      "info",
		MaxSizeMB:  1,
		MaxBackups: 1,
		MaxAgeDays: 1,
		Compress:   true,
	}, console)
	if err != nil {
		t.Fatalf("new logger: %v", err)
	}

	handle.Logger.Info("runtime ready")
	if err := handle.Close(); err != nil {
		t.Fatalf("close logger: %v", err)
	}
	if console.syncCalls != 0 {
		t.Fatalf("console sync calls = %d, want 0", console.syncCalls)
	}
	if console.Len() == 0 {
		t.Fatal("console received no log output")
	}
}

func TestNewRejectsInvalidLevel(t *testing.T) {
	_, err := New(config.LoggingConfig{
		Path:      filepath.Join(t.TempDir(), "corvus.log"),
		Level:     "verbose",
		MaxSizeMB: 1,
	})
	if err == nil {
		t.Fatal("new logger with invalid level returned nil error")
	}
}
