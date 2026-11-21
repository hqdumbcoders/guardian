package proxy

import (
	"bytes"
	"hqdumbcoders/guardian/internal/requestid"
	"io"
	"log"
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

func TestReverseProxyLogsUpstreamResponse(t *testing.T) {
	var (
		tests     string
		logBuffer bytes.Buffer
		found     bool
	)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("erdtree ok"))
	},
	))
	defer server.Close()

	target, err := url.Parse(server.URL)
	if err != nil {
		t.Fatal(err)
	}

	proxy := New(target)

	handler := requestid.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, ok := requestid.FromContext(r.Context())
		if !ok || id == "" {
			t.Fatal("expected not empty correlationIDKey")
		}

		tests = id
		proxy.ServeHTTP(w, r)
	}))
	loger := log.Writer()
	log.SetOutput(&logBuffer)
	t.Cleanup(func() { log.SetOutput(loger) })

	request := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/products",
		nil,
	)

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Errorf("expected: %d, received: %d", http.StatusOK, recorder.Code)
	}

	if recorder.Body.String() != "erdtree ok" {
		t.Errorf("expected: %s, received: %s", "erdtree ok", recorder.Body.String())
	}

	lines := strings.Split(logBuffer.String(), "\n")
	for _, line := range lines {
		if strings.Contains(line, tests) && strings.Contains(line, "status: 200") {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("expected ID: %v,expected status: %s, received: %v", tests, "status: 200", lines)
	}

}
