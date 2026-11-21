package main

import (
	"hqdumbcoders/guardian/internal/config"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealthHandler(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		path       string
		wantStatus int
		wantBody   string
		wantAllow  string
	}{
		{
			name:       "GET returns OK",
			method:     "GET",
			path:       "/health",
			wantStatus: http.StatusOK,
			wantBody:   "OK",
			wantAllow:  "",
		},
		{
			name:       "POST returns Not Allowed",
			method:     "POST",
			path:       "/health",
			wantStatus: http.StatusMethodNotAllowed,
			wantBody:   "405 Not Allowed",
			wantAllow:  "GET",
		},
		{
			name:       "PUT returns Not Allowed",
			method:     "PUT",
			path:       "/health",
			wantStatus: http.StatusMethodNotAllowed,
			wantBody:   "405 Not Allowed",
			wantAllow:  "GET",
		},
		{
			name:       "PATCH returns Not Allowed",
			method:     "PATCH",
			path:       "/health",
			wantStatus: http.StatusMethodNotAllowed,
			wantBody:   "405 Not Allowed",
			wantAllow:  "GET",
		},
		{
			name:       "DELETE returns Not Allowed",
			method:     "DELETE",
			path:       "/health",
			wantStatus: http.StatusMethodNotAllowed,
			wantBody:   "405 Not Allowed",
			wantAllow:  "GET",
		},
		{
			name:       "unknown path",
			method:     "GET",
			path:       "/unknown",
			wantStatus: http.StatusNotFound,
			wantBody:   "404 page not found\n",
		},
	}

	cfg := config.Config{
		ErdtreeURL: "http://localhost:8001",
		GladosURL:  "http://localhost:8002",
	}

	mux, err := newRouter(cfg)
	if err != nil {
		t.Fatalf("failed to create router: %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, tt.path, nil)
			recorder := httptest.NewRecorder()

			mux.ServeHTTP(recorder, request)

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

func TestNewServer(t *testing.T) {
	cfg := config.Config{
		ErdtreeURL: "http://localhost:8001",
		GladosURL:  "http://localhost:8002",
	}

	serv, err := newServer(cfg)
	if err != nil {
		t.Fatalf("failed to create server: %v", err)
	}

	request := httptest.NewRequest("GET", "/health", nil)
	recorder := httptest.NewRecorder()

	if serv.Handler == nil {
		t.Fatal("expected Handler to be non-nil")
	}

	serv.Handler.ServeHTTP(recorder, request)

	if serv.Addr != ":8000" {
		t.Errorf("expected: %s, received: %s", ":8000", serv.Addr)
	}

	if recorder.Code != http.StatusOK {
		t.Errorf("expected: %d, received: %d", http.StatusOK, recorder.Code)
	}

}

func TestParseUpstreamURL(t *testing.T) {
	tests := []struct {
		name    string
		rawURL  string
		wantErr bool
	}{
		{"valid HTTP", "http://glados:8002", false},
		{"valid HTTPS", "https://example.com", false},
		{"invalid scheme", "ftp://glados:8002", true},
		{"missing scheme", "glados:8002", true},
		{"empty host", "http:///users", true},
		{"malformed URL", "http://%zz", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			target, err := parseUpstreamURL(tt.rawURL)

			if (err != nil) != tt.wantErr {
				t.Errorf("expected error: %t, got error: %t (%v)",
					tt.wantErr, err != nil, err)
			}

			if !tt.wantErr && target == nil {
				t.Fatal("expected non-nil URL, got nil")
			}

			if !tt.wantErr && target.String() != tt.rawURL {
				t.Errorf("expected: %q, received: %q", tt.rawURL, target.String())
			}

		})
	}

}

func TestNewRouterInvalidGladosURL(t *testing.T) {
	cfg := config.Config{
		ErdtreeURL: "http://erdtree:8001",
		GladosURL:  "ftp://glados:8002",
	}

	router, err := newRouter(cfg)
	if err == nil {
		t.Fatal("expected non-nil error, got nil")
	}

	if router != nil {
		t.Fatalf("expected nil router, got %v", router)
	}

	if !strings.Contains(err.Error(), "glados URL:") {
		t.Errorf("expected: %s, received: %q", "glados URL: ...", err.Error())
	}

}

func TestNewRouterInvalidErdtreeURL(t *testing.T) {
	cfg := config.Config{
		ErdtreeURL: "ftp://erdtree:8001",
		GladosURL:  "http://glados:8002",
	}

	router, err := newRouter(cfg)
	if err == nil {
		t.Fatal("expected non-nil error, got nil")
	}

	if router != nil {
		t.Fatalf("expected nil router, got %v", router)
	}

	if !strings.Contains(err.Error(), "erdtree URL:") {
		t.Errorf("expected: %s, received: %q", "erdtree URL: ...", err.Error())
	}

}

func TestNewRouterGetProductByID(t *testing.T) {
	erdtree := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}

		if r.URL.Path != "/api/v1/products/42" {
			t.Errorf("expected /api/v1/products/42, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("erdtree response"))
	}))
	defer erdtree.Close()

	cfg := config.Config{
		ErdtreeURL: erdtree.URL,
		GladosURL:  "http://glados:8002",
	}

	router, err := newRouter(cfg)
	if err != nil {
		t.Fatal(err)
	}

	request := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/products/42",
		nil,
	)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Errorf("expected: %d, received: %d", http.StatusOK, recorder.Code)
	}

	if recorder.Body.String() != "erdtree response" {
		t.Errorf("expected: erdtree response, received: %s", recorder.Body.String())
	}

}

func TestNewRouterProductRouteRejections(t *testing.T) {
	tests := []struct {
		method     string
		path       string
		wantStatus int
	}{
		{"POST", "/api/v1/products/42", http.StatusMethodNotAllowed},
		{"GET", "/api/v1/products/42/details", http.StatusNotFound},
	}

	cfg := config.Config{
		ErdtreeURL: "http://erdtree:8001",
		GladosURL:  "http://glados:8002",
	}

	router, err := newRouter(cfg)
	if err != nil {
		t.Fatal(err)
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, tt.path, nil)
			recorder := httptest.NewRecorder()

			router.ServeHTTP(recorder, request)

			if recorder.Code != tt.wantStatus {
				t.Errorf("method: %s, expected: %d, received: %d", tt.method, tt.wantStatus, recorder.Code)
			}

		})
	}
}

func TestNewRouterCreateProduct(t *testing.T) {
	erdtree := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

		w.WriteHeader(http.StatusCreated)
	}))
	defer erdtree.Close()

	cfg := config.Config{
		ErdtreeURL: erdtree.URL,
		GladosURL:  "http://glados:8002",
	}

	router, err := newRouter(cfg)
	if err != nil {
		t.Fatal(err)
	}

	request := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/products",
		strings.NewReader(`{"name":"Massage"}`),
	)
	request.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusCreated {
		t.Errorf("expected: %d, received: %d", http.StatusCreated, recorder.Code)
	}

}

func TestNewRouterDeleteProduct(t *testing.T) {
	erdtree := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}

		if r.URL.Path != "/api/v1/products/42" {
			t.Errorf("expected /api/v1/products/42, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer erdtree.Close()

	cfg := config.Config{
		ErdtreeURL: erdtree.URL,
		GladosURL:  "http://glados:8002",
	}

	router, err := newRouter(cfg)
	if err != nil {
		t.Fatal(err)
	}

	request := httptest.NewRequest(
		http.MethodDelete,
		"/api/v1/products/42",
		nil,
	)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusNoContent {
		t.Errorf("expected: %d, received: %d", http.StatusNoContent, recorder.Code)
	}

	if recorder.Body.String() != "" {
		t.Errorf("expected: empty body, received: %s", recorder.Body.String())
	}

}

func TestNewRouterUpdateProduct(t *testing.T) {
	tests := []struct {
		method string
		body   string
	}{
		{"PUT", `{"name":"Deep Tissue Massage"}`},
		{"PATCH", `{"name":"Relaxing Massage"}`},
	}

	for _, tt := range tests {
		t.Run(tt.method, func(t *testing.T) {
			erdtree := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != tt.method {
					t.Errorf("expected %s, got %s", tt.method, r.Method)
				}

				if r.URL.Path != "/api/v1/products/42" {
					t.Errorf("expected /api/v1/products/42, got %s", r.URL.Path)
				}

				if r.Header.Get("Content-Type") != "application/json" {
					t.Errorf("expected application/json, got %s", r.Header.Get("Content-Type"))
				}

				body, err := io.ReadAll(r.Body)
				if err != nil {
					t.Error(err)
					return
				}

				if string(body) != tt.body {
					t.Errorf("expected %s, got %s", tt.body, string(body))
				}

				w.WriteHeader(http.StatusOK)
			}))
			defer erdtree.Close()

			cfg := config.Config{
				ErdtreeURL: erdtree.URL,
				GladosURL:  "http://glados:8002",
			}

			router, err := newRouter(cfg)
			if err != nil {
				t.Fatal(err)
			}

			request := httptest.NewRequest(
				tt.method,
				"/api/v1/products/42",
				strings.NewReader(tt.body),
			)
			request.Header.Set("Content-Type", "application/json")
			recorder := httptest.NewRecorder()

			router.ServeHTTP(recorder, request)

			if recorder.Code != http.StatusOK {
				t.Errorf("method: %s, expected: %d, received: %d", tt.method, http.StatusOK, recorder.Code)
			}

		})
	}
}

func TestNewRouterListSlots(t *testing.T) {
	erdtree := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}

		if r.URL.Path != "/api/v1/slots" {
			t.Errorf("expected /api/v1/slots, got %s", r.URL.Path)
		}

		if r.URL.RawQuery != "product_id=42" {
			t.Errorf("expected product_id=42, got %s", r.URL.RawQuery)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer erdtree.Close()

	cfg := config.Config{
		ErdtreeURL: erdtree.URL,
		GladosURL:  "http://glados:8002",
	}

	router, err := newRouter(cfg)
	if err != nil {
		t.Fatal(err)
	}

	request := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/slots?product_id=42",
		nil,
	)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Errorf("expected: %d, received: %d", http.StatusOK, recorder.Code)
	}

}

func TestNewRouterGetSlotByID(t *testing.T) {
	erdtree := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}

		if r.URL.Path != "/api/v1/slots/42" {
			t.Errorf("expected /api/v1/slots/42, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id":42}`))
	}))
	defer erdtree.Close()

	cfg := config.Config{
		ErdtreeURL: erdtree.URL,
		GladosURL:  "http://glados:8002",
	}

	router, err := newRouter(cfg)
	if err != nil {
		t.Fatal(err)
	}

	request := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/slots/42",
		nil,
	)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Errorf("expected: %d, received: %d", http.StatusOK, recorder.Code)
	}

	if recorder.Result().Header.Get("Content-Type") != "application/json" {
		t.Errorf("expected application/json, got %s", recorder.Result().Header.Get("Content-Type"))
	}

}

func TestNewRouterCreateSlotsUpstreamError(t *testing.T) {
	erdtree := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}

		if r.URL.Path != "/api/v1/slots" {
			t.Errorf("expected /api/v1/slots, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid slot data"))
	}))
	defer erdtree.Close()

	cfg := config.Config{
		ErdtreeURL: erdtree.URL,
		GladosURL:  "http://glados:8002",
	}

	router, err := newRouter(cfg)
	if err != nil {
		t.Fatal(err)
	}

	request := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/slots",
		strings.NewReader("{}"),
	)
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Errorf("expected: %d, received: %d", http.StatusBadRequest, recorder.Code)
	}

	if recorder.Body.String() != "invalid slot data" {
		t.Errorf("expected: invalid slot data, received: %s", recorder.Body.String())
	}

}

func TestNewRouterSlotByIDMutations(t *testing.T) {
	tests := []string{
		"PATCH",
		"DELETE",
	}

	for _, tt := range tests {
		t.Run(tt, func(t *testing.T) {
			erdtree := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != tt {
					t.Errorf("expected %s, got %s", tt, r.Method)
				}

				if r.URL.Path != "/api/v1/slots/42" {
					t.Errorf("expected /api/v1/slots/42, got %s", r.URL.Path)
				}

				w.Header().Set("X-Test-Upstream", "erdtree")
				w.WriteHeader(http.StatusNoContent)
			}))
			defer erdtree.Close()

			cfg := config.Config{
				ErdtreeURL: erdtree.URL,
				GladosURL:  "http://glados:8002",
			}

			router, err := newRouter(cfg)
			if err != nil {
				t.Fatal(err)
			}

			request := httptest.NewRequest(
				tt,
				"/api/v1/slots/42",
				nil,
			)

			recorder := httptest.NewRecorder()

			router.ServeHTTP(recorder, request)

			if recorder.Code != http.StatusNoContent {
				t.Errorf("method: %s, expected: %d, received: %d", tt, http.StatusNoContent, recorder.Code)
			}

			if recorder.Result().Header.Get("X-Test-Upstream") != "erdtree" {
				t.Errorf("expected X-Test-Upstream: erdtree, got %s", recorder.Result().Header.Get("X-Test-Upstream"))
			}

		})
	}
}

func TestNewRouterBookingRoutes(t *testing.T) {
	tests := []struct {
		method     string
		path       string
		wantStatus int
	}{
		{"POST", "/api/v1/bookings", http.StatusCreated},
		{"GET", "/api/v1/bookings/42", http.StatusOK},
		{"DELETE", "/api/v1/bookings/42", http.StatusNoContent},
	}

	for _, tt := range tests {
		t.Run(tt.method, func(t *testing.T) {
			erdtree := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != tt.method {
					t.Errorf("expected %s, got %s", tt.method, r.Method)
				}

				if r.URL.Path != tt.path {
					t.Errorf("expected %s, got %s", tt.path, r.URL.Path)
				}

				w.Header().Set("X-Test-Upstream", "erdtree")
				w.WriteHeader(tt.wantStatus)
			}))
			defer erdtree.Close()

			cfg := config.Config{
				ErdtreeURL: erdtree.URL,
				GladosURL:  "http://glados:8002",
			}

			router, err := newRouter(cfg)
			if err != nil {
				t.Fatal(err)
			}

			request := httptest.NewRequest(
				tt.method,
				tt.path,
				nil,
			)

			recorder := httptest.NewRecorder()

			router.ServeHTTP(recorder, request)

			if recorder.Code != tt.wantStatus {
				t.Errorf("method: %s, expected: %d, received: %d", tt.method, tt.wantStatus, recorder.Code)
			}

			if recorder.Result().Header.Get("X-Test-Upstream") != "erdtree" {
				t.Errorf("expected X-Test-Upstream: erdtree, got %s", recorder.Result().Header.Get("X-Test-Upstream"))
			}

		})
	}
}

func TestNewRouterGladosPostRoutes(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		path       string
		wantStatus int
	}{
		{"register", "POST", "/api/v1/users", http.StatusCreated},
		{"login", "POST", "/api/v1/auth/login", http.StatusOK},
		{"logout", "POST", "/api/v1/auth/logout", http.StatusNoContent},
		{"refresh", "POST", "/api/v1/auth/refresh", http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			glados := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != tt.method {
					t.Errorf("expected %s, got %s", tt.method, r.Method)
				}

				if r.URL.Path != tt.path {
					t.Errorf("expected %s, got %s", tt.path, r.URL.Path)
				}

				w.Header().Set("X-Test-Upstream", "glados")
				w.WriteHeader(tt.wantStatus)
			}))
			defer glados.Close()

			cfg := config.Config{
				ErdtreeURL: "http://erdtree:8001",
				GladosURL:  glados.URL,
			}

			router, err := newRouter(cfg)
			if err != nil {
				t.Fatal(err)
			}

			request := httptest.NewRequest(
				tt.method,
				tt.path,
				nil,
			)

			recorder := httptest.NewRecorder()

			router.ServeHTTP(recorder, request)

			if recorder.Code != tt.wantStatus {
				t.Errorf("name: %s, expected: %d, received: %d", tt.name, tt.wantStatus, recorder.Code)
			}

			if recorder.Result().Header.Get("X-Test-Upstream") != "glados" {
				t.Errorf("expected X-Test-Upstream: glados, got %s", recorder.Result().Header.Get("X-Test-Upstream"))
			}

		})
	}
}
