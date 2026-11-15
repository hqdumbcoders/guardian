package proxy

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestReverseProxy(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}

		if r.URL.Path != "/api/v1/products" {
			t.Errorf("expected /api/v1/products, got %s", r.URL.Path)
		}

		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected application/json, got %s", r.Header.Get("Content-Type"))
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Error(err)
			return
		}

		if string(body) != `{"name":"Massage"}` {
			t.Errorf("expected %s, got %s", `{"name":"Massage"}`, string(body))
		}

		w.WriteHeader(http.StatusNoContent)
	},
	))
	defer server.Close()

	target, err := url.Parse(server.URL)
	if err != nil {
		t.Fatal(err)
	}

	proxy := New(target)

	request := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/products",
		strings.NewReader(`{"name":"Massage"}`),
	)
	request.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	proxy.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusNoContent {
		t.Errorf("expected: %d, received: %d", http.StatusNoContent, recorder.Code)
	}
}
