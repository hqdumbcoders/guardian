package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthHandler(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		wantStatus int
		wantBody   string
		wantAllow  string
	}{
		{
			name:       "GET returns OK",
			method:     "GET",
			wantStatus: http.StatusOK,
			wantBody:   "OK",
			wantAllow:  "",
		},
		{
			name:       "POST returns Not Allowed",
			method:     "POST",
			wantStatus: http.StatusMethodNotAllowed,
			wantBody:   "405 Not Allowed",
			wantAllow:  "GET",
		},
		{
			name:       "PUT returns Not Allowed",
			method:     "PUT",
			wantStatus: http.StatusMethodNotAllowed,
			wantBody:   "405 Not Allowed",
			wantAllow:  "GET",
		},
		{
			name:       "PATCH returns Not Allowed",
			method:     "PATCH",
			wantStatus: http.StatusMethodNotAllowed,
			wantBody:   "405 Not Allowed",
			wantAllow:  "GET",
		},
		{
			name:       "DELETE returns Not Allowed",
			method:     "DELETE",
			wantStatus: http.StatusMethodNotAllowed,
			wantBody:   "405 Not Allowed",
			wantAllow:  "GET",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, "/health", nil)
			recorder := httptest.NewRecorder()

			healthHandler(recorder, request)

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
