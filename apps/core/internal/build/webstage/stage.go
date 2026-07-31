package webstage

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

func Stage(source, destination string) error {
	sourcePath, err := filepath.Abs(source)
	if err != nil {
		return fmt.Errorf("resolve Web source: %w", err)
	}
	destinationPath, err := filepath.Abs(destination)
	if err != nil {
		return fmt.Errorf("resolve Web destination: %w", err)
	}
	if filepath.Clean(sourcePath) == filepath.Clean(destinationPath) {
		return errors.New("web source and destination must differ")
	}
	if _, err := os.Stat(filepath.Join(sourcePath, "index.html")); err != nil {
		return fmt.Errorf("web source does not contain index.html: %w", err)
	}

	parent := filepath.Dir(destinationPath)
	if err := os.MkdirAll(parent, 0o755); err != nil {
		return fmt.Errorf("create Web staging parent: %w", err)
	}
	temporaryPath, err := os.MkdirTemp(parent, ".dist-stage-")
	if err != nil {
		return fmt.Errorf("create temporary Web stage: %w", err)
	}
	defer func() { _ = os.RemoveAll(temporaryPath) }()

	if err := copyTree(sourcePath, temporaryPath); err != nil {
		return err
	}

	backupPath := ""
	if _, err := os.Stat(destinationPath); err == nil {
		backupPath, err = unusedSiblingPath(parent, ".dist-backup-")
		if err != nil {
			return err
		}
		if err := os.Rename(destinationPath, backupPath); err != nil {
			return fmt.Errorf("preserve previous Web stage: %w", err)
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("inspect Web destination: %w", err)
	}

	if err := os.Rename(temporaryPath, destinationPath); err != nil {
		if backupPath != "" {
			_ = os.Rename(backupPath, destinationPath)
		}
		return fmt.Errorf("activate Web stage: %w", err)
	}
	if backupPath != "" {
		if err := os.RemoveAll(backupPath); err != nil {
			return fmt.Errorf("remove previous generated Web stage: %w", err)
		}
	}
	return nil
}

func copyTree(source, destination string) error {
	return fs.WalkDir(os.DirFS(source), ".", func(relativePath string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if relativePath == "." {
			return nil
		}
		if entry.Type()&fs.ModeSymlink != 0 {
			return fmt.Errorf("web stage does not accept symbolic links: %s", relativePath)
		}
		target := filepath.Join(destination, filepath.FromSlash(relativePath))
		if entry.IsDir() {
			if err := os.MkdirAll(target, 0o755); err != nil {
				return fmt.Errorf("create staged directory %q: %w", relativePath, err)
			}
			return nil
		}
		if !entry.Type().IsRegular() {
			return fmt.Errorf("web stage accepts regular files only: %s", relativePath)
		}
		if err := copyFile(filepath.Join(source, filepath.FromSlash(relativePath)), target); err != nil {
			return fmt.Errorf("copy Web asset %q: %w", relativePath, err)
		}
		return nil
	})
}

func copyFile(source, destination string) error {
	input, err := os.Open(source)
	if err != nil {
		return err
	}
	defer func() { _ = input.Close() }()
	output, err := os.OpenFile(destination, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		return err
	}
	if _, err := io.Copy(output, input); err != nil {
		_ = output.Close()
		return err
	}
	return output.Close()
}

func unusedSiblingPath(parent, pattern string) (string, error) {
	path, err := os.MkdirTemp(parent, pattern)
	if err != nil {
		return "", fmt.Errorf("reserve generated backup path: %w", err)
	}
	if err := os.Remove(path); err != nil {
		return "", fmt.Errorf("release generated backup path: %w", err)
	}
	return path, nil
}
