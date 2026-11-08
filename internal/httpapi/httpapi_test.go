package httpapi_test

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/VlasovArtem/palai/internal/httpapi"
)

func newTestHandler() http.Handler {
	return httpapi.NewHandler(slog.New(slog.DiscardHandler))
}

func TestHealthz(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		method     string
		path       string
		wantStatus int
		wantBody   map[string]string
	}{
		{
			name:       "GET returns ok",
			method:     http.MethodGet,
			path:       "/healthz",
			wantStatus: http.StatusOK,
			wantBody:   map[string]string{"status": "ok"},
		},
		{
			name:       "HEAD is allowed for GET routes",
			method:     http.MethodHead,
			path:       "/healthz",
			wantStatus: http.StatusOK,
		},
		{
			name:       "POST is not allowed",
			method:     http.MethodPost,
			path:       "/healthz",
			wantStatus: http.StatusMethodNotAllowed,
		},
		{
			name:       "unknown path",
			method:     http.MethodGet,
			path:       "/nope",
			wantStatus: http.StatusNotFound,
		},
	}

	h := newTestHandler()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequestWithContext(t.Context(), tt.method, tt.path, nil)
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
			if tt.wantBody == nil {
				return
			}
			if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
				t.Errorf("Content-Type = %q, want application/json", ct)
			}
			var got map[string]string
			if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
				t.Fatalf("decode body %q: %v", rec.Body.String(), err)
			}
			if got["status"] != tt.wantBody["status"] || len(got) != len(tt.wantBody) {
				t.Errorf("body = %v, want %v", got, tt.wantBody)
			}
		})
	}
}

func TestHealthzOverHTTP(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(newTestHandler())
	t.Cleanup(srv.Close)

	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, srv.URL+"/healthz", nil)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	resp, err := srv.Client().Do(req)
	if err != nil {
		t.Fatalf("do request: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}
	if want := `{"status":"ok"}` + "\n"; string(body) != want {
		t.Errorf("body = %q, want %q", body, want)
	}
}
