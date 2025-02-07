package main

import (
	"errors"
	"fmt"
	"os"
)

var ErrNotImplemeted = errors.New("not implemented")

func run() error {
	return ErrNotImplemeted
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "datahow: %s\n", err)
		os.Exit(1)
	}
}
