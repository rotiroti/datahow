package datahow_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/rotiroti/datahow"
)

func assertStatusCode(tb testing.TB, got, want int) {
	tb.Helper()

	if got != want {
		tb.Errorf("status code returned %d; expected %d", got, want)
	}
}

func assertBool(tb testing.TB, got, want bool) {
	tb.Helper()

	if got != want {
		tb.Errorf("got %t; expected %t", got, want)
	}
}

func assertError(tb testing.TB, got, want error) {
	tb.Helper()

	if !errors.Is(got, want) {
		tb.Errorf("error %v; expected %v", got, want)
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
			s := datahow.NewInMemory()
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

	s := datahow.NewInMemory()
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

func TestHandleLog(t *testing.T) {
	tests := []struct {
		name    string
		payload string
		want    int
	}{
		{
			name:    "empty IP address",
			payload: `{}`,
			want:    http.StatusUnprocessableEntity,
		},
		{
			name:    "new IP addrees",
			payload: `{"timestamp": "2020-06-24T15:27:00.123456Z", "ip": "83.150.59.250", "url": "https://datahow.ch"}`,
			want:    http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := datahow.NewInMemory()
			srv := datahow.NewLogServer(store)
			ctx := context.Background()
			req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/logs", bytes.NewBufferString(tt.payload))
			assertError(t, err, nil)

			res := httptest.NewRecorder()
			srv.ServeHTTP(res, req)
			assertStatusCode(t, res.Code, tt.want)
		})
	}
}
