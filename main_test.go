package main

import (
	"errors"
	"testing"
)

func TestRun(t *testing.T) {
	if err := run(); !errors.Is(err, ErrNotImplemeted) {
		t.Error("expected error, got nil")
	}
}
