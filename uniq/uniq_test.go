package uniq_test

import (
	"sync"
	"testing"

	"github.com/rotiroti/datahow/uniq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
			s := uniq.NewInMemory()
			if tt.name == "dup IP" {
				_ = s.ExistOrAdd(tt.ip)
			}

			got := s.ExistOrAdd(tt.ip)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSafeExistsOrAdd(t *testing.T) {
	s := uniq.NewInMemory()
	ip := "83.150.59.250"
	require.False(t, s.ExistOrAdd(ip))

	var wg sync.WaitGroup

	maxGoroutines := 1000
	wg.Add(maxGoroutines)

	for range maxGoroutines {
		go func() {
			defer wg.Done()

			assert.True(t, s.ExistOrAdd(ip))
		}()
	}

	wg.Wait()
}
