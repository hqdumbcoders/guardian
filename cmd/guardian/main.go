package main

import (
	"fmt"
	"hqdumbcoders/guardian/internal/config"
	"hqdumbcoders/guardian/internal/proxy"
	"log"
	"net/http"
	"net/url"
)

const (
	appName    = "guardian"
	appVersion = "0.1.0"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	} else {
		w.Header().Set("Allow", "GET")
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("405 Not Allowed"))
	}
}

func parseUpstreamURL(rawURL string) (*url.URL, error) {
	target, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("parse upstream URL: %w", err)
	}

	if target.Scheme != "http" && target.Scheme != "https" {
		return nil, fmt.Errorf("expected: %s, received: %s", "http or https", target.Scheme)
	}

	if target.Host == "" {
		return nil, fmt.Errorf("upstream URL: empty host")
	}

	return target, nil
}

func newRouter(cfg config.Config) (*http.ServeMux, error) {
	gladosURL, err := parseUpstreamURL(cfg.GladosURL)
	if err != nil {
		return nil, fmt.Errorf("glados URL: %w", err)
	}

	erdtreeURL, err := parseUpstreamURL(cfg.ErdtreeURL)
	if err != nil {
		return nil, fmt.Errorf("erdtree URL: %w", err)
	}

	gladosProxy := proxy.New(gladosURL)
	erdtreeProxy := proxy.New(erdtreeURL)

	mux := http.NewServeMux()

	mux.Handle("POST /api/v1/auth/login", gladosProxy)

	mux.Handle("POST /api/v1/products", erdtreeProxy)
	mux.Handle("GET /api/v1/products", erdtreeProxy)
	mux.Handle("GET /api/v1/products/{id}", erdtreeProxy)
	mux.Handle("PUT /api/v1/products/{id}", erdtreeProxy)
	mux.Handle("PATCH /api/v1/products/{id}", erdtreeProxy)
	mux.Handle("DELETE /api/v1/products/{id}", erdtreeProxy)

	mux.Handle("POST /api/v1/slots", erdtreeProxy)
	mux.Handle("GET /api/v1/slots", erdtreeProxy)
	mux.Handle("GET /api/v1/slots/{slot_id}", erdtreeProxy)
	mux.Handle("PATCH /api/v1/slots/{slot_id}", erdtreeProxy)
	mux.Handle("DELETE /api/v1/slots/{slot_id}", erdtreeProxy)

	mux.Handle("POST /api/v1/bookings", erdtreeProxy)
	mux.Handle("GET /api/v1/bookings/{id}", erdtreeProxy)
	mux.Handle("DELETE /api/v1/bookings/{id}", erdtreeProxy)

	mux.HandleFunc("/health", healthHandler)
	return mux, nil
}

func newServer(cfg config.Config) (*http.Server, error) {
	handler, err := newRouter(cfg)
	if err != nil {
		return nil, fmt.Errorf("create router: %w", err)
	}

	return &http.Server{
		Addr:    ":8000",
		Handler: handler,
	}, nil
}

func main() {
	fmt.Printf("starting: %s %s\n", appName, appVersion)

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	server, err := newServer(cfg)
	if err != nil {
		log.Fatal(err)
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
