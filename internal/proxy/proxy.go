package proxy

import (
	"context"
	"errors"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

func New(target *url.URL) http.Handler {
	reverseProxy := httputil.NewSingleHostReverseProxy(target)

	clone := http.DefaultTransport.(*http.Transport).Clone()
	clone.ResponseHeaderTimeout = 5 * time.Second
	reverseProxy.Transport = clone

	reverseProxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		if errors.Is(err, context.DeadlineExceeded) {
			log.Printf("status: %d upstream timeout", http.StatusGatewayTimeout)
			w.WriteHeader(http.StatusGatewayTimeout)
			return
		}

		log.Printf("status: %d upstream request failed", http.StatusBadGateway)
		w.WriteHeader(http.StatusBadGateway)
	}

	return reverseProxy
}
