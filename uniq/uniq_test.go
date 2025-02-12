package uniq_test

import (
	"sync"
	"testing"

	"github.com/rotiroti/datahow/uniq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHSetAdd(t *testing.T) {
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
			s := uniq.NewHSet()
			if tt.name == "dup IP" {
				_ = s.Add(tt.ip)
			}

			got := s.Add(tt.ip)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestHSetAddConcurrent(t *testing.T) {
	s := uniq.NewHSet()
	ip := "83.150.59.250"
	require.False(t, s.Add(ip))

	var wg sync.WaitGroup

	maxGoroutines := 1000
	wg.Add(maxGoroutines)

	for range maxGoroutines {
		go func() {
			defer wg.Done()

			assert.True(t, s.Add(ip))
		}()
	}

	wg.Wait()
}

func TestHSetCount(t *testing.T) {
	s := uniq.NewHSet()
	require.False(t, s.Add("83.150.59.250"))
	assert.Equal(t, int64(1), s.Count())
}

func TestHSetCountConcurrent(t *testing.T) {
	s := uniq.NewHSet()
	ip := "83.150.59.250"
	require.False(t, s.Add(ip))

	var wg sync.WaitGroup

	maxGoroutines := 1000
	wg.Add(maxGoroutines)

	for range maxGoroutines {
		go func() {
			defer wg.Done()

			assert.Equal(t, int64(1), s.Count())
		}()
	}

	wg.Wait()
}
