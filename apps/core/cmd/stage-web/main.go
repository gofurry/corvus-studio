package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gofurry/corvus-studio/apps/core/internal/build/webstage"
)

func run(args []string, workingDirectory string) error {
	repositoryRoot, err := findRepositoryRoot(workingDirectory)
	if err != nil {
		return err
	}
	defaultSource := filepath.Join(repositoryRoot, "apps", "web", "dist")
	expectedDestination := filepath.Join(repositoryRoot, "apps", "core", "internal", "interfaces", "webui", "dist")

	flags := flag.NewFlagSet("stage-web", flag.ContinueOnError)
	source := flags.String("source", defaultSource, "Vite dist directory")
	destination := flags.String("destination", expectedDestination, "generated in-module Web directory")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() != 0 {
		return fmt.Errorf("unexpected arguments: %v", flags.Args())
	}

	resolvedDestination, err := filepath.Abs(*destination)
	if err != nil {
		return fmt.Errorf("resolve Web destination: %w", err)
	}
	if filepath.Clean(resolvedDestination) != filepath.Clean(expectedDestination) {
		return fmt.Errorf("destination must remain inside the Core Web UI boundary: %s", expectedDestination)
	}
	return webstage.Stage(*source, resolvedDestination)
}

func findRepositoryRoot(start string) (string, error) {
	current, err := filepath.Abs(start)
	if err != nil {
		return "", fmt.Errorf("resolve working directory: %w", err)
	}
	for {
		if fileExists(filepath.Join(current, "go.work")) && fileExists(filepath.Join(current, "pnpm-workspace.yaml")) {
			return current, nil
		}
		parent := filepath.Dir(current)
		if parent == current {
			return "", fmt.Errorf("repository root not found from %q", start)
		}
		current = parent
	}
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func main() {
	workingDirectory, err := os.Getwd()
	if err == nil {
		err = run(os.Args[1:], workingDirectory)
	}
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "stage-web:", err)
		os.Exit(1)
	}
}
