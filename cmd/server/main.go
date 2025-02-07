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
	"github.com/rotiroti/datahow"
)

func run(ctx context.Context) error {
	ipsCounter := promauto.NewCounter(prometheus.CounterOpts{
		Name: "unique_ip_addresses",
		Help: "No. of unique IP addresses",
	})

	// Configure API Log server
	store := datahow.NewInMemory()
	logServer := datahow.NewLogServer(store, ipsCounter)
	httpLogServer := http.Server{
		Addr:              ":5001", // TODO: Pass this as environment variable
		Handler:           logServer,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      5 * time.Second,
		ReadHeaderTimeout: 3 * time.Second,
		IdleTimeout:       30 * time.Second,
	}

	// Create a context that listens for interrupt signals
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()

	logServerErrChan := make(chan error, 1)

	go func() {
		slog.Info("Starting HTTP log server...")

		if err := httpLogServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("HTTP log server failed to listen and serve", "error", err)
			logServerErrChan <- err
		}
	}()

	select {
	case <-ctx.Done():
	case err := <-logServerErrChan:
		return fmt.Errorf("HTTP log server startup failed: %w", err)
	}

	ctxShutdown, stopShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer stopShutdown()

	slog.Info("Shutting down HTTP log server...")

	if err := httpLogServer.Shutdown(ctxShutdown); err != nil {
		return fmt.Errorf("HTTP log server shutdown failed: %w", err)
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
