package datahow_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/magiconair/properties/assert"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/rotiroti/datahow"
	"github.com/rotiroti/datahow/uniq"
	"github.com/stretchr/testify/require"
)

func TestHandleLog(t *testing.T) {
	ctx := context.Background()
	store := uniq.NewInMemory()
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
	err = testutil.CollectAndCompare(counter, strings.NewReader(wantMetricFormat))
	require.NoError(t, err)
}
