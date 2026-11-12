package main

import (
	"fmt"
	"hqdumbcoders/guardian/internal/config"
	"log"
	"net/http"
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

func newRouter(_ config.Config) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthHandler)
	return mux
}

func newServer(cfg config.Config) *http.Server {
	return &http.Server{
		Addr:    ":8000",
		Handler: newRouter(cfg),
	}
}

func main() {
	fmt.Printf("starting: %s %s\n", appName, appVersion)

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	server := newServer(cfg)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
