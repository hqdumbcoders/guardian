package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

func New(target *url.URL) http.Handler {
	return httputil.NewSingleHostReverseProxy(target)
}
