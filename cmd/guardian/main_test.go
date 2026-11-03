package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthHandler(t *testing.T) {
	request := httptest.NewRequest("GET", "/health", nil)
	recorder := httptest.NewRecorder()

	healthHandler(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Errorf("expected: %d, received: %d", http.StatusOK, recorder.Code)
	}

	if recorder.Body.String() != "OK" {
		t.Errorf("expected: %s, received: %s", "OK", recorder.Body.String())
	}

}

func TestHealthHandlerMethodNotAllowed(t *testing.T) {
	request := httptest.NewRequest("POST", "/health", nil)
	recorder := httptest.NewRecorder()

	healthHandler(recorder, request)

	if recorder.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected: %d, received: %d", http.StatusMethodNotAllowed, recorder.Code)
	}

	if recorder.Body.String() != "405 Not Allowed" {
		t.Errorf("expected: %s, received: %s", "405 Not Allowed", recorder.Body.String())
	}

	if recorder.Header().Get("Allow") != "GET" {
		t.Errorf("expected: %s, received: %s", "GET", recorder.Header().Get("Allow"))
	}

}
