package proxy

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
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
