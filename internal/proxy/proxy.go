package proxy

import (
	"context"
	"errors"
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
			w.WriteHeader(http.StatusGatewayTimeout)
			return
		}

		w.WriteHeader(http.StatusBadGateway)
	}

	return reverseProxy
}
