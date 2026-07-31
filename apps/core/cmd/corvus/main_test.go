package main

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestRunShowsCobraHelp(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	if err := run(context.Background(), []string{"--help"}, &stdout, &stderr); err != nil {
		t.Fatalf("run help: %v", err)
	}

	output := stdout.String()
	if !strings.Contains(output, "Corvus Studio Core") || !strings.Contains(output, "serve") {
		t.Fatalf("help output missing Core or serve command:\n%s", output)
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
}

func TestRunRejectsUnknownCommand(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	err := run(context.Background(), []string{"unknown"}, &stdout, &stderr)
	if err == nil {
		t.Fatal("run unknown command returned nil error")
	}
}
