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

func New(target *url.URL, upstreamName string) http.Handler {
	reverseProxy := httputil.NewSingleHostReverseProxy(target)

	clone := http.DefaultTransport.(*http.Transport).Clone()
	clone.ResponseHeaderTimeout = 5 * time.Second
	reverseProxy.Transport = clone

	reverseProxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		id, ok := requestid.FromContext(r.Context())
		if errors.Is(err, context.DeadlineExceeded) {
			if !ok || id == "" {
				log.Printf("upstream=%s, status=%d, method=%s, path=%s, "+
					"upstream timeout", upstreamName, http.StatusGatewayTimeout,
					r.Method, r.URL.EscapedPath())
			} else {
				log.Printf("upstream=%s, status=%d, method=%s, path=%s, "+
					"correlation_id=%s, upstream timeout", upstreamName,
					http.StatusGatewayTimeout, r.Method, r.URL.EscapedPath(), id)
			}

			w.WriteHeader(http.StatusGatewayTimeout)
			return
		}

		if !ok || id == "" {
			log.Printf("upstream=%s, status=%d, method=%s, path=%s, "+
				"upstream request failed", upstreamName, http.StatusBadGateway,
				r.Method, r.URL.EscapedPath())
		} else {
			log.Printf("upstream=%s, status=%d, method=%s, path=%s, "+
				"correlation_id=%s, upstream request failed", upstreamName,
				http.StatusBadGateway, r.Method, r.URL.EscapedPath(), id)
		}

		w.WriteHeader(http.StatusBadGateway)
	}

	reverseProxy.ModifyResponse = func(res *http.Response) error {
		id, ok := requestid.FromContext(res.Request.Context())
		if !ok || id == "" {
			log.Printf("upstream=%s, status=%d, method=%s, path=%s", upstreamName, res.StatusCode,
				res.Request.Method, res.Request.URL.EscapedPath())
		} else {
			log.Printf("upstream=%s, status=%d, method=%s, path=%s, correlation_id=%s",
				upstreamName, res.StatusCode, res.Request.Method, res.Request.URL.EscapedPath(), id)
		}

		return nil
	}

	return reverseProxy
}
