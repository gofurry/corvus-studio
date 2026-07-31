package projectpath

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"

	projectdomain "github.com/gofurry/corvus-studio/apps/core/internal/domain/project"
)

type Resolver struct{}

func (Resolver) Resolve(value string) (string, string, error) {
	input := strings.TrimSpace(value)
	if input == "" || !filepath.IsAbs(input) {
		return "", "", projectdomain.NewValidationError(
			"location",
			"must reference an existing absolute directory",
		)
	}

	resolved, err := filepath.EvalSymlinks(filepath.Clean(input))
	if err != nil {
		return "", "", projectdomain.NewValidationError(
			"location",
			"must reference an existing absolute directory",
		)
	}
	resolved = filepath.Clean(resolved)
	info, err := os.Stat(resolved)
	if err != nil || !info.IsDir() {
		return "", "", projectdomain.NewValidationError(
			"location",
			"must reference an existing absolute directory",
		)
	}

	key := resolved
	if runtime.GOOS == "windows" {
		key = strings.ToLower(key)
	}
	return resolved, key, nil
}
