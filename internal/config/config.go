// Package config loads Palai configuration from environment variables.
package config

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"
)

// Prefix is prepended to every environment variable Palai reads.
const Prefix = "PALAI_"

// Defaults used when a variable is unset or empty.
const (
	DefaultAddr            = ":8080"
	DefaultLogLevel        = slog.LevelInfo
	DefaultShutdownTimeout = 10 * time.Second
)

// Config is the runtime configuration of the service.
type Config struct {
	// Addr is the HTTP listen address (PALAI_ADDR).
	Addr string
	// LogLevel is the minimum slog level (PALAI_LOG_LEVEL: debug, info, warn, error).
	LogLevel slog.Level
	// ShutdownTimeout bounds graceful shutdown (PALAI_SHUTDOWN_TIMEOUT, Go duration).
	ShutdownTimeout time.Duration
}

// LookupFunc matches the signature of os.LookupEnv.
type LookupFunc func(key string) (string, bool)

// Load reads the configuration from the process environment.
func Load() (Config, error) {
	return LoadFrom(os.LookupEnv)
}

// LoadFrom reads the configuration using lookup, which makes it testable
// without touching the process environment.
func LoadFrom(lookup LookupFunc) (Config, error) {
	cfg := Config{
		Addr:            DefaultAddr,
		LogLevel:        DefaultLogLevel,
		ShutdownTimeout: DefaultShutdownTimeout,
	}

	var errs []error

	if v, ok := get(lookup, "ADDR"); ok {
		cfg.Addr = v
	}

	if v, ok := get(lookup, "LOG_LEVEL"); ok {
		if err := cfg.LogLevel.UnmarshalText([]byte(v)); err != nil {
			errs = append(errs, fmt.Errorf("%sLOG_LEVEL: %w", Prefix, err))
		}
	}

	if v, ok := get(lookup, "SHUTDOWN_TIMEOUT"); ok {
		d, err := time.ParseDuration(v)
		switch {
		case err != nil:
			errs = append(errs, fmt.Errorf("%sSHUTDOWN_TIMEOUT: %w", Prefix, err))
		case d <= 0:
			errs = append(errs, fmt.Errorf("%sSHUTDOWN_TIMEOUT: must be positive, got %s", Prefix, d))
		default:
			cfg.ShutdownTimeout = d
		}
	}

	if err := errors.Join(errs...); err != nil {
		return Config{}, fmt.Errorf("load config: %w", err)
	}
	return cfg, nil
}

// get returns the trimmed value of PALAI_<name>; empty values count as unset.
func get(lookup LookupFunc, name string) (string, bool) {
	v, ok := lookup(Prefix + name)
	v = strings.TrimSpace(v)
	return v, ok && v != ""
}
