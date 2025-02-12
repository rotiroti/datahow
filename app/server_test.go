package app_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rotiroti/datahow/app"
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

func (m *MockCounter) Count() int64 {
	args := m.Called()

	return args.Get(0).(int64)
}

type MockMetricCollector struct {
	mock.Mock
}

func (m *MockMetricCollector) Inc() {
	m.Called()
}

func TestHandleLogBadRequest(t *testing.T) {
	ctx := context.Background()
	srv := app.NewServer(nil, nil)
	req := httptest.NewRequestWithContext(ctx, http.MethodPost, "/logs", nil)
	res := httptest.NewRecorder()
	srv.ServeHTTP(res, req)
	assert.Equal(t, http.StatusBadRequest, res.Code)
}

func TestHandleLog(t *testing.T) {
	tests := []struct {
		name        string
		ip          string
		wantStatus  int
		isDuplicate bool
		shouldInc   bool
	}{
		{
			name:        "unique IP",
			ip:          "83.150.59.250",
			isDuplicate: false,
			wantStatus:  http.StatusOK,
			shouldInc:   true,
		},
		{
			name:        "duplicate IP",
			ip:          "83.150.59.250",
			isDuplicate: true,
			wantStatus:  http.StatusOK,
			shouldInc:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			mockCounter := new(MockCounter)
			mockCounter.On("Add", tt.ip).Return(tt.isDuplicate)

			mockMetricsCollector := new(MockMetricCollector)
			if !tt.isDuplicate {
				mockMetricsCollector.On("Inc").Once()
			}

			srv := app.NewServer(mockCounter, mockMetricsCollector)

			reqBody := struct {
				IPAddress string `json:"ip"`
			}{
				IPAddress: tt.ip,
			}

			var buf bytes.Buffer
			err := json.NewEncoder(&buf).Encode(reqBody)
			require.NoError(t, err)

			req := httptest.NewRequestWithContext(ctx, http.MethodPost, "/logs", &buf)
			res := httptest.NewRecorder()
			srv.ServeHTTP(res, req)

			assert.Equal(t, tt.wantStatus, res.Code)
			mockCounter.AssertExpectations(t)
			mockMetricsCollector.AssertExpectations(t)
		})
	}
}
