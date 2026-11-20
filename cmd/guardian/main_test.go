package main

import (
	"hqdumbcoders/guardian/internal/config"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealthHandler(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		path       string
		wantStatus int
		wantBody   string
		wantAllow  string
	}{
		{
			name:       "GET returns OK",
			method:     "GET",
			path:       "/health",
			wantStatus: http.StatusOK,
			wantBody:   "OK",
			wantAllow:  "",
		},
		{
			name:       "POST returns Not Allowed",
			method:     "POST",
			path:       "/health",
			wantStatus: http.StatusMethodNotAllowed,
			wantBody:   "405 Not Allowed",
			wantAllow:  "GET",
		},
		{
			name:       "PUT returns Not Allowed",
			method:     "PUT",
			path:       "/health",
			wantStatus: http.StatusMethodNotAllowed,
			wantBody:   "405 Not Allowed",
			wantAllow:  "GET",
		},
		{
			name:       "PATCH returns Not Allowed",
			method:     "PATCH",
			path:       "/health",
			wantStatus: http.StatusMethodNotAllowed,
			wantBody:   "405 Not Allowed",
			wantAllow:  "GET",
		},
		{
			name:       "DELETE returns Not Allowed",
			method:     "DELETE",
			path:       "/health",
			wantStatus: http.StatusMethodNotAllowed,
			wantBody:   "405 Not Allowed",
			wantAllow:  "GET",
		},
		{
			name:       "unknown path",
			method:     "GET",
			path:       "/unknown",
			wantStatus: http.StatusNotFound,
			wantBody:   "404 page not found\n",
		},
	}

	cfg := config.Config{
		ErdtreeURL: "http://localhost:8001",
		GladosURL:  "http://localhost:8002",
	}

	mux, err := newRouter(cfg)
	if err != nil {
		t.Fatalf("failed to create router: %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, tt.path, nil)
			recorder := httptest.NewRecorder()

			mux.ServeHTTP(recorder, request)

			if recorder.Code != tt.wantStatus {
				t.Errorf("method: %s, expected: %d, received: %d", tt.method, tt.wantStatus, recorder.Code)
			}

			if recorder.Body.String() != tt.wantBody {
				t.Errorf("method: %s, expected: %s, received: %s", tt.method, tt.wantBody, recorder.Body.String())
			}

			if recorder.Header().Get("Allow") != tt.wantAllow {
				t.Errorf("method: %s, expected: %s, received: %s", tt.method, tt.wantAllow, recorder.Header().Get("Allow"))
			}
		})

	}

}

func TestNewServer(t *testing.T) {
	cfg := config.Config{
		ErdtreeURL: "http://localhost:8001",
		GladosURL:  "http://localhost:8002",
	}

	serv, err := newServer(cfg)
	if err != nil {
		t.Fatalf("failed to create server: %v", err)
	}

	request := httptest.NewRequest("GET", "/health", nil)
	recorder := httptest.NewRecorder()

	if serv.Handler == nil {
		t.Fatal("expected Handler to be non-nil")
	}

	serv.Handler.ServeHTTP(recorder, request)

	if serv.Addr != ":8000" {
		t.Errorf("expected: %s, received: %s", ":8000", serv.Addr)
	}

	if recorder.Code != http.StatusOK {
		t.Errorf("expected: %d, received: %d", http.StatusOK, recorder.Code)
	}

}

func TestParseUpstreamURL(t *testing.T) {
	tests := []struct {
		name    string
		rawURL  string
		wantErr bool
	}{
		{"valid HTTP", "http://glados:8002", false},
		{"valid HTTPS", "https://example.com", false},
		{"invalid scheme", "ftp://glados:8002", true},
		{"missing scheme", "glados:8002", true},
		{"empty host", "http:///users", true},
		{"malformed URL", "http://%zz", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			target, err := parseUpstreamURL(tt.rawURL)

			if (err != nil) != tt.wantErr {
				t.Errorf("expected error: %t, got error: %t (%v)",
					tt.wantErr, err != nil, err)
			}

			if !tt.wantErr && target == nil {
				t.Fatal("expected non-nil URL, got nil")
			}

			if !tt.wantErr && target.String() != tt.rawURL {
				t.Errorf("expected: %q, received: %q", tt.rawURL, target.String())
			}

		})
	}

}

func TestNewRouterInvalidGladosURL(t *testing.T) {
	cfg := config.Config{
		ErdtreeURL: "http://erdtree:8001",
		GladosURL:  "ftp://glados:8002",
	}

	router, err := newRouter(cfg)
	if err == nil {
		t.Fatal("expected non-nil error, got nil")
	}

	if router != nil {
		t.Fatalf("expected nil router, got %v", router)
	}

	if !strings.Contains(err.Error(), "glados URL:") {
		t.Errorf("expected: %s, received: %q", "glados URL: ...", err.Error())
	}

}

func TestNewRouterInvalidErdtreeURL(t *testing.T) {
	cfg := config.Config{
		ErdtreeURL: "ftp://erdtree:8001",
		GladosURL:  "http://glados:8002",
	}

	router, err := newRouter(cfg)
	if err == nil {
		t.Fatal("expected non-nil error, got nil")
	}

	if router != nil {
		t.Fatalf("expected nil router, got %v", router)
	}

	if !strings.Contains(err.Error(), "erdtree URL:") {
		t.Errorf("expected: %s, received: %q", "erdtree URL: ...", err.Error())
	}

}
