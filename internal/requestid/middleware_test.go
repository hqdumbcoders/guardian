package requestid

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWithCorrelationID(t *testing.T) {
	tests := []string{}
	handler := Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, ok := FromContext(r.Context())
		if !ok || id == "" {
			t.Fatal("expected not empty correlationIDKey")
		}

		tests = append(tests, id)
	}))

	request := httptest.NewRequest(
		"GET",
		"/api/v1/products",
		nil,
	)
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)

	request = httptest.NewRequest(
		"POST",
		"/api/v1/products",
		nil,
	)
	recorder = httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)

	if len(tests) != 2 {
		t.Fatal("expected 2 requests")
	}

	if tests[0] == tests[1] {
		t.Fatal("two requests received the same ID")
	}
}
