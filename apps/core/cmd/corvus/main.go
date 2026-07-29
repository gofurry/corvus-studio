package main

import (
	"fmt"
	"io"
	"os"
)

const bootstrapMessage = "Corvus Studio core bootstrap"

func run(output io.Writer) error {
	_, err := fmt.Fprintln(output, bootstrapMessage)
	return err
}

func main() {
	if err := run(os.Stdout); err != nil {
		os.Exit(1)
	}
}
