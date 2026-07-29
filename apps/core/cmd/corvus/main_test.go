package main

import (
	"bytes"
	"testing"
)

func TestRunPrintsBootstrapMessage(t *testing.T) {
	t.Parallel()

	var output bytes.Buffer
	if err := run(&output); err != nil {
		t.Fatalf("run() error = %v", err)
	}

	want := bootstrapMessage + "\n"
	if got := output.String(); got != want {
		t.Fatalf("run() output = %q, want %q", got, want)
	}
}
