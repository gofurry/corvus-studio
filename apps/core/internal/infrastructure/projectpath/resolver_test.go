package projectpath

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	projectdomain "github.com/gofurry/corvus-studio/apps/core/internal/domain/project"
)

func TestResolverAcceptsExistingDirectoryWithoutWriting(t *testing.T) {
	directory := t.TempDir()
	sentinel := filepath.Join(directory, "sentinel.txt")
	if err := os.WriteFile(sentinel, []byte("preserve"), 0o600); err != nil {
		t.Fatalf("write sentinel: %v", err)
	}

	location, key, err := (Resolver{}).Resolve(directory)
	if err != nil {
		t.Fatalf("resolve directory: %v", err)
	}
	if !filepath.IsAbs(location) || key == "" {
		t.Fatalf("resolved location = %q, key = %q", location, key)
	}
	content, err := os.ReadFile(sentinel)
	if err != nil {
		t.Fatalf("read sentinel: %v", err)
	}
	if string(content) != "preserve" {
		t.Fatalf("sentinel was modified: %q", content)
	}
	entries, err := os.ReadDir(directory)
	if err != nil {
		t.Fatalf("read directory: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("directory entries = %d, want 1", len(entries))
	}
}

func TestResolverRejectsUnsupportedLocations(t *testing.T) {
	file := filepath.Join(t.TempDir(), "game.txt")
	if err := os.WriteFile(file, []byte("game"), 0o600); err != nil {
		t.Fatalf("write file: %v", err)
	}
	tests := []struct {
		name     string
		location string
	}{
		{name: "empty", location: ""},
		{name: "relative", location: filepath.Join("relative", "game")},
		{name: "missing", location: filepath.Join(t.TempDir(), "missing")},
		{name: "file", location: file},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, _, err := (Resolver{}).Resolve(test.location)
			if !errors.Is(err, projectdomain.ErrValidation) {
				t.Fatalf("error = %v, want validation", err)
			}
		})
	}
}
