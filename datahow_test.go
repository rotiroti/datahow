package datahow_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/rotiroti/datahow"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockCounter struct {
	mock.Mock
}

func (m *MockCounter) Add(ip string) bool {
	args := m.Called(ip)

	return args.Bool(0)
}

func TestHandleLogBadRequest(t *testing.T) {
	ctx := context.Background()
	srv := datahow.NewLogServer(nil, nil)
	req := httptest.NewRequestWithContext(ctx, http.MethodPost, "/logs", nil)
	res := httptest.NewRecorder()
	srv.ServeHTTP(res, req)
	assert.Equal(t, http.StatusBadRequest, res.Code)
}

func TestHandleLog(t *testing.T) {
	// NOTE: The current test case verifies the actual increment of the Prometheus counter.
	// According to the Prometheus testutil documentation, it is more robust and faithful to the concept of unit tests
	// to use mock implementations of the prometheus.Counter and prometheus.Registerer interfaces that simply assert
	// that the Add or Register methods have been called with the expected arguments.
	// Reference: https://pkg.go.dev/github.com/prometheus/client_golang/prometheus/testutil
	ctx := context.Background()
	mockCounter := new(MockCounter)
	mockCounter.On("Add", "83.150.59.250").Return(false)
	promCounter := promauto.NewCounter(prometheus.CounterOpts{
		Name: "unique_ip_addresses",
		Help: "No. of unique IP addresses",
	})

	srv := datahow.NewLogServer(mockCounter, promCounter)
	payload := struct {
		IPAddress string `json:"ip"`
	}{
		IPAddress: "83.150.59.250", // we care only about the IP address
	}

	var buf bytes.Buffer

	err := json.NewEncoder(&buf).Encode(payload)
	require.NoError(t, err)

	req := httptest.NewRequestWithContext(ctx, http.MethodPost, "/logs", &buf)
	res := httptest.NewRecorder()
	srv.ServeHTTP(res, req)
	assert.Equal(t, http.StatusOK, res.Code)

	// Define the expected metric format
	wantMetricFormat := `
		# HELP unique_ip_addresses No. of unique IP addresses
		# TYPE unique_ip_addresses counter
		unique_ip_addresses 1
	`
	err = testutil.CollectAndCompare(promCounter, strings.NewReader(wantMetricFormat))
	require.NoError(t, err)
}
