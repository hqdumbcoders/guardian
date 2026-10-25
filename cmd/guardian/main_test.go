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
