// Command palai runs the Palai HTTP API server.
package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/VlasovArtem/palai/internal/config"
	"github.com/VlasovArtem/palai/internal/httpapi"
)

func main() {
	// JSON from the first line, even before config is loaded.
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := run(ctx); err != nil {
		slog.Error("palai exited with error", slog.Any("error", err))
		stop()
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("run: %w", err)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: cfg.LogLevel}))
	slog.SetDefault(logger)

	srv := &http.Server{
		Addr:              cfg.Addr,
		Handler:           httpapi.NewHandler(logger),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       120 * time.Second,
		ErrorLog:          slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	serveErr := make(chan error, 1)
	go func() {
		logger.InfoContext(ctx, "http server listening", slog.String("addr", cfg.Addr))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serveErr <- fmt.Errorf("listen and serve: %w", err)
		}
		close(serveErr)
	}()

	select {
	case err := <-serveErr:
		return err
	case <-ctx.Done():
	}

	logger.InfoContext(ctx, "shutting down", slog.String("timeout", cfg.ShutdownTimeout.String()))
	// The parent ctx is already cancelled; shutdown needs its own deadline.
	shutdownCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), cfg.ShutdownTimeout)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown: %w", err)
	}
	logger.InfoContext(shutdownCtx, "shutdown complete")
	return nil
}
