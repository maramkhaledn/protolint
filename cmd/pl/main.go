package main

import (
	"os"

	"github.com/maramkhaledn/protolint/internal/cmd"
)

// DEPRECATED: Use cmd/protolint. See https://github.com/maramkhaledn/protolint/issues/20.
func main() {
	// Initialize the lint runner
	cmd.Initialize()

	os.Exit(int(
		cmd.Do(
			os.Args[1:],
			os.Stdout,
			os.Stderr,
		),
	))
}
