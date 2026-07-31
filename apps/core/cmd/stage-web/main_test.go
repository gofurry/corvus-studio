package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindRepositoryRoot(t *testing.T) {
	root := t.TempDir()
	for _, name := range []string{"go.work", "pnpm-workspace.yaml"} {
		if err := os.WriteFile(filepath.Join(root, name), []byte("test"), 0o600); err != nil {
			t.Fatalf("write marker: %v", err)
		}
	}
	nested := filepath.Join(root, "apps", "core")
	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatalf("create nested path: %v", err)
	}

	got, err := findRepositoryRoot(nested)
	if err != nil {
		t.Fatalf("find root: %v", err)
	}
	if got != root {
		t.Fatalf("root = %q, want %q", got, root)
	}
}

func TestRunRejectsDestinationOutsideCoreBoundary(t *testing.T) {
	root := t.TempDir()
	for _, name := range []string{"go.work", "pnpm-workspace.yaml"} {
		if err := os.WriteFile(filepath.Join(root, name), []byte("test"), 0o600); err != nil {
			t.Fatalf("write marker: %v", err)
		}
	}
	err := run([]string{"--destination", filepath.Join(root, "outside")}, root)
	if err == nil {
		t.Fatal("run outside destination returned nil error")
	}
}
