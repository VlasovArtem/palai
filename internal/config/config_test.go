package config_test

import (
	"log/slog"
	"testing"
	"time"

	"github.com/VlasovArtem/palai/internal/config"
)

func lookupFrom(env map[string]string) config.LookupFunc {
	return func(key string) (string, bool) {
		v, ok := env[key]
		return v, ok
	}
}

func TestLoadFrom(t *testing.T) {
	t.Parallel()

	defaults := config.Config{
		Addr:            config.DefaultAddr,
		LogLevel:        config.DefaultLogLevel,
		ShutdownTimeout: config.DefaultShutdownTimeout,
	}

	tests := []struct {
		name    string
		env     map[string]string
		want    config.Config
		wantErr bool
	}{
		{
			name: "defaults when unset",
			env:  nil,
			want: defaults,
		},
		{
			name: "empty values fall back to defaults",
			env:  map[string]string{"PALAI_ADDR": "  ", "PALAI_LOG_LEVEL": "", "PALAI_SHUTDOWN_TIMEOUT": ""},
			want: defaults,
		},
		{
			name: "all values set",
			env: map[string]string{
				"PALAI_ADDR":             "127.0.0.1:9090",
				"PALAI_LOG_LEVEL":        "debug",
				"PALAI_SHUTDOWN_TIMEOUT": "30s",
			},
			want: config.Config{Addr: "127.0.0.1:9090", LogLevel: slog.LevelDebug, ShutdownTimeout: 30 * time.Second},
		},
		{
			name: "log level is case-insensitive",
			env:  map[string]string{"PALAI_LOG_LEVEL": "WARN"},
			want: config.Config{Addr: config.DefaultAddr, LogLevel: slog.LevelWarn, ShutdownTimeout: config.DefaultShutdownTimeout},
		},
		{
			name: "unprefixed variables are ignored",
			env:  map[string]string{"ADDR": ":1"},
			want: defaults,
		},
		{
			name:    "invalid log level",
			env:     map[string]string{"PALAI_LOG_LEVEL": "loud"},
			wantErr: true,
		},
		{
			name:    "invalid shutdown timeout",
			env:     map[string]string{"PALAI_SHUTDOWN_TIMEOUT": "soon"},
			wantErr: true,
		},
		{
			name:    "non-positive shutdown timeout",
			env:     map[string]string{"PALAI_SHUTDOWN_TIMEOUT": "0s"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := config.LoadFrom(lookupFrom(tt.env))
			if tt.wantErr {
				if err == nil {
					t.Fatalf("LoadFrom() error = nil, want error (got config %+v)", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("LoadFrom() unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("LoadFrom() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestLoadReadsProcessEnv(t *testing.T) {
	t.Setenv("PALAI_ADDR", ":7070")

	got, err := config.Load()
	if err != nil {
		t.Fatalf("Load() unexpected error: %v", err)
	}
	if got.Addr != ":7070" {
		t.Errorf("Load().Addr = %q, want %q", got.Addr, ":7070")
	}
}
