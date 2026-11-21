package proxy

import (
	"context"
	"errors"
	"hqdumbcoders/guardian/internal/requestid"
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
		id, ok := requestid.FromContext(r.Context())
		if errors.Is(err, context.DeadlineExceeded) {
			if !ok || id == "" {
				log.Printf("status: %d upstream timeout", http.StatusGatewayTimeout)
			} else {
				log.Printf("correlation_id=%s, status: %d upstream timeout", id, http.StatusGatewayTimeout)
			}

			w.WriteHeader(http.StatusGatewayTimeout)
			return
		}

		if !ok || id == "" {
			log.Printf("status: %d upstream request failed", http.StatusBadGateway)
		} else {
			log.Printf("correlation_id=%s, status: %d upstream request failed", id, http.StatusBadGateway)
		}

		w.WriteHeader(http.StatusBadGateway)
	}

	return reverseProxy
}
