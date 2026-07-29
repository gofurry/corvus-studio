package webstage

import (
	"os"
	"path/filepath"
	"testing"
)

func TestStageCopiesAndAtomicallyReplacesAssets(t *testing.T) {
	root := t.TempDir()
	source := filepath.Join(root, "source")
	destination := filepath.Join(root, "destination")
	writeAsset(t, source, "index.html", "first")
	writeAsset(t, source, "assets/app.js", "app")
	writeAsset(t, destination, "stale.js", "stale")

	if err := Stage(source, destination); err != nil {
		t.Fatalf("stage: %v", err)
	}
	assertAsset(t, destination, "index.html", "first")
	assertAsset(t, destination, "assets/app.js", "app")
	if _, err := os.Stat(filepath.Join(destination, "stale.js")); !os.IsNotExist(err) {
		t.Fatalf("stale asset remains, stat error = %v", err)
	}
}

func TestStagePreservesDestinationWhenSourceIsInvalid(t *testing.T) {
	root := t.TempDir()
	source := filepath.Join(root, "source")
	destination := filepath.Join(root, "destination")
	writeAsset(t, source, "asset.js", "missing index")
	writeAsset(t, destination, "index.html", "existing")

	if err := Stage(source, destination); err == nil {
		t.Fatal("stage invalid source returned nil error")
	}
	assertAsset(t, destination, "index.html", "existing")
}

func writeAsset(t *testing.T, root, relativePath, content string) {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(relativePath))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("create asset directory: %v", err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write asset: %v", err)
	}
}

func assertAsset(t *testing.T, root, relativePath, want string) {
	t.Helper()
	content, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(relativePath)))
	if err != nil {
		t.Fatalf("read asset %s: %v", relativePath, err)
	}
	if string(content) != want {
		t.Fatalf("asset %s = %q, want %q", relativePath, content, want)
	}
}
