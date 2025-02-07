package datahow_test

import (
	"fmt"
	"sync"
	"testing"

	"github.com/rotiroti/datahow"
)

func assertBool(tb testing.TB, got, want bool) {
	tb.Helper()

	if got != want {
		tb.Errorf("got %v want %v", got, want)
	}
}

func TestExistOrAdd(t *testing.T) {
	tests := []struct {
		name string
		ip   string
		want bool
	}{
		{"new IP", "83.150.59.250", false},
		{"dup IP", "83.150.59.250", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := datahow.New()
			if tt.name == "dup IP" {
				_ = s.ExistOrAdd(tt.ip)
			}

			got := s.ExistOrAdd(tt.ip)
			assertBool(t, got, tt.want)
		})
	}
}

func TestSafeExistsOrAdd(t *testing.T) {
	var wg sync.WaitGroup

	s := datahow.New()
	maxGoroutines := 256
	wg.Add(maxGoroutines)

	for i := range maxGoroutines {
		go func(v int) {
			defer wg.Done()

			got := s.ExistOrAdd(fmt.Sprintf("83.150.59.%d", v))
			assertBool(t, got, false)
		}(i)
	}

	wg.Wait()
}
