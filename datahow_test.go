package datahow_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
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
	ctx := context.Background()
	store := datahow.NewInMemory()
	counter := promauto.NewCounter(prometheus.CounterOpts{
		Name: "unique_ip_addresses",
		Help: "No. of unique IP addresses",
	})
	srv := datahow.NewLogServer(store, counter)
	payload := struct {
		IPAddress string `json:"ip"`
	}{
		IPAddress: "83.150.59.250", // we care only about the IP address
	}

	var buf bytes.Buffer

	err := json.NewEncoder(&buf).Encode(payload)
	assertError(t, err, nil)

	req := httptest.NewRequestWithContext(ctx, http.MethodPost, "/logs", &buf)
	res := httptest.NewRecorder()
	srv.ServeHTTP(res, req)
	assertStatusCode(t, res.Code, http.StatusOK)
}
