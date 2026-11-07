package main

import (
	"hqdumbcoders/guardian/internal/config"
	"net/http"
	"net/http/httptest"
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

	mux := newRouter(cfg)

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

	serv := newServer(cfg)

	request := httptest.NewRequest("GET", "/health", nil)
	recorder := httptest.NewRecorder()

	if serv.Handler == nil {
		t.Fatal("expected Handler to be non-nil")
	}

	serv.Handler.ServeHTTP(recorder, request)

	if serv.Addr != ":8080" {
		t.Errorf("expected: %s, received: %s", ":8080", serv.Addr)
	}

	if recorder.Code != http.StatusOK {
		t.Errorf("expected: %d, received: %d", http.StatusOK, recorder.Code)
	}

}
