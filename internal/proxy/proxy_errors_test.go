package proxy

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func TestReverseProxyUnavailableUpstream(t *testing.T) {
	server := httptest.NewServer(http.NotFoundHandler())
	server.Close()

	target, err := url.Parse(server.URL)
	if err != nil {
		t.Fatal(err)
	}

	proxy := New(target)

	request := httptest.NewRequest(
		http.MethodGet,
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

	ctx, cancel := context.WithTimeout(request.Context(), 100*time.Millisecond)
	request = request.WithContext(ctx)
	defer cancel()

	recorder := httptest.NewRecorder()
	proxy.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusGatewayTimeout {
		t.Errorf("expected: %d, received: %d", http.StatusGatewayTimeout, recorder.Code)
	}
}
