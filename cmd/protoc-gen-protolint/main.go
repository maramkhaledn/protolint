package main

import (
	"os"

	protoc "github.com/maramkhaledn/protolint/internal/cmd/protocgenprotolint"
)

func main() {
	os.Exit(int(
		protoc.Do(
			os.Args[1:],
			os.Stdin,
			os.Stdout,
			os.Stderr,
		),
	))
}
