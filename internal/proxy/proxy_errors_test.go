package proxy

import (
	"bytes"
	"context"
	"hqdumbcoders/guardian/internal/requestid"
	"log"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"strings"
	"testing"
	"time"
)

func TestReverseProxyUnavailableUpstream(t *testing.T) {
	var (
		tests     string
		logBuffer bytes.Buffer
	)

	server := httptest.NewServer(http.NotFoundHandler())
	server.Close()

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
	if !strings.Contains(logBuffer.String(), tests) {
		t.Errorf("expected: %v, received: %v", tests, logBuffer.String())
	}

	if recorder.Code != http.StatusBadGateway {
		t.Errorf("expected: %d, received: %d", http.StatusBadGateway, recorder.Code)
	}

}

func TestReverseProxyUnavailableUpstreamWithoutCorrelationID(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
	},
	))
	server.Close()

	target, err := url.Parse(server.URL)
	if err != nil {
		t.Fatal(err)
	}

	proxy := New(target)

	request := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/products",
		nil,
	)

	recorder := httptest.NewRecorder()
	proxy.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadGateway {
		t.Errorf("expected: %d, received: %d", http.StatusBadGateway, recorder.Code)
	}

}

func TestReverseProxyUpstreamTimeout(t *testing.T) {
	var (
		tests     string
		logBuffer bytes.Buffer
	)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-r.Context().Done()
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
		http.MethodPost,
		"/api/v1/products",
		nil,
	)

	ctx, cancel := context.WithTimeout(request.Context(), 100*time.Millisecond)
	request = request.WithContext(ctx)
	defer cancel()

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)

	if !strings.Contains(logBuffer.String(), tests) {
		t.Errorf("expected: %v, received: %v", tests, logBuffer.String())
	}

	if recorder.Code != http.StatusGatewayTimeout {
		t.Errorf("expected: %d, received: %d", http.StatusGatewayTimeout, recorder.Code)
	}
}

func TestReverseProxyResponseHeaderTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-r.Context().Done()
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
		nil,
	)
	recorder := httptest.NewRecorder()

	p := proxy.(*httputil.ReverseProxy)
	transport := p.Transport.(*http.Transport)

	if transport.ResponseHeaderTimeout <= 0 {
		t.Fatal("expected a positive ResponseHeaderTimeout")
	}

	transport.ResponseHeaderTimeout = 100 * time.Millisecond

	proxy.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusGatewayTimeout {
		t.Errorf("expected: %d, received: %d", http.StatusGatewayTimeout, recorder.Code)
	}
}
