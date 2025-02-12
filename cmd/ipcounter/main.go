package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rotiroti/datahow"
	"github.com/rotiroti/datahow/uniq"
)

func run(ctx context.Context) error {
	ipsCounter := promauto.NewCounter(prometheus.CounterOpts{
		Name: "unique_ip_addresses",
		Help: "No. of unique IP addresses",
	})

	// Configure API Log server
	hs := uniq.NewHSet()
	logServer := datahow.NewLogServer(hs, ipsCounter)

	// TODO: create a function to configure the HTTP log server
	httpLogServer := http.Server{
		Addr:              ":5000",
		Handler:           logServer,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      5 * time.Second,
		ReadHeaderTimeout: 3 * time.Second,
		IdleTimeout:       30 * time.Second,
	}

	// Configure metrics server
	metricsServer := http.NewServeMux()
	metricsServer.Handle("/metrics", promhttp.Handler())

	// TODO: create a function to configure the HTTP metrics server
	httpMetricsServer := &http.Server{
		Addr:              ":9102",
		Handler:           metricsServer,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      5 * time.Second,
		ReadHeaderTimeout: 3 * time.Second,
		IdleTimeout:       30 * time.Second,
	}

	// Create a context that listens for interrupt signals
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()

	logServerErrChan := make(chan error, 1)
	metricsServerErrChan := make(chan error, 1)

	go func() {
		slog.Info("Starting HTTP log server...")

		if err := httpLogServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("HTTP log server failed to listen and serve", "error", err)
			logServerErrChan <- err
		}
	}()

	go func() {
		slog.Info("Starting HTTP metrics server...")

		if err := httpMetricsServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("HTTP metrics server failed to listen and serve", "error", err)
			metricsServerErrChan <- err
		}
	}()

	select {
	case <-ctx.Done():
	case err := <-logServerErrChan:
		return fmt.Errorf("HTTP log server startup failed: %w", err)
	case err := <-metricsServerErrChan:
		return fmt.Errorf("HTTP metrics server startup failed: %w", err)
	}

	ctxShutdown, stopShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer stopShutdown()

	if err := httpLogServer.Shutdown(ctxShutdown); err != nil {
		return fmt.Errorf("HTTP log server shutdown failed: %w", err)
	}

	if err := httpMetricsServer.Shutdown(ctxShutdown); err != nil {
		return fmt.Errorf("HTTP metrics server shutdown failed: %w", err)
	}

	return nil
}

func main() {
	ctx := context.Background()

	if err := run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "datahow: %s\n", err)
		os.Exit(1)
	}
}
