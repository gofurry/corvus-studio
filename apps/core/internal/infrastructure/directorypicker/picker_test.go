package directorypicker

import (
	"context"
	"errors"
	"path/filepath"
	"testing"
)

func TestPickerSelectsExistingAbsoluteDirectory(t *testing.T) {
	directory := t.TempDir()
	picker := testPicker(
		"windows",
		func(context.Context, string, ...string) ([]byte, error) {
			return []byte(directory + "\r\n"), nil
		},
	)

	selectedPath, selected, err := picker.Select(t.Context())
	if err != nil {
		t.Fatalf("select directory: %v", err)
	}
	if !selected || selectedPath != filepath.Clean(directory) {
		t.Fatalf("selection = (%q, %t), want %q", selectedPath, selected, directory)
	}
}

func TestPickerReportsCancellation(t *testing.T) {
	tests := []struct {
		name     string
		platform string
		run      runCommandFunc
	}{
		{
			name:     "marker",
			platform: "windows",
			run: func(context.Context, string, ...string) ([]byte, error) {
				return []byte(cancelledMarker), nil
			},
		},
		{
			name:     "Linux exit code",
			platform: "linux",
			run: func(context.Context, string, ...string) ([]byte, error) {
				return nil, fakeExitError{code: 1}
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			picker := testPicker(test.platform, test.run)
			selectedPath, selected, err := picker.Select(t.Context())
			if err != nil {
				t.Fatalf("cancel selection: %v", err)
			}
			if selected || selectedPath != "" {
				t.Fatalf("selection = (%q, %t), want cancelled", selectedPath, selected)
			}
		})
	}
}

func TestPickerRejectsUnavailableCommandAndInvalidOutput(t *testing.T) {
	unavailable := &Picker{
		platform: "linux",
		lookPath: func(string) (string, error) {
			return "", errors.New("missing")
		},
		run: func(context.Context, string, ...string) ([]byte, error) {
			t.Fatal("run called without a picker command")
			return nil, nil
		},
	}
	if _, _, err := unavailable.Select(t.Context()); !errors.Is(err, ErrUnavailable) {
		t.Fatalf("unavailable error = %v", err)
	}

	invalid := testPicker(
		"windows",
		func(context.Context, string, ...string) ([]byte, error) {
			return []byte("relative/path"), nil
		},
	)
	if _, _, err := invalid.Select(t.Context()); !errors.Is(err, ErrUnavailable) {
		t.Fatalf("invalid output error = %v", err)
	}
}

func testPicker(platform string, run runCommandFunc) *Picker {
	return &Picker{
		platform: platform,
		lookPath: func(name string) (string, error) {
			return name, nil
		},
		run: run,
	}
}

type fakeExitError struct {
	code int
}

func (err fakeExitError) Error() string { return "command exited" }
func (err fakeExitError) ExitCode() int { return err.code }
