package requestid

import (
	"context"
	"crypto/rand"
	"net/http"
)

type correlationKey string

const correlationIDKey correlationKey = "CorrelationID"

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		corID := rand.Text()
		ctx := context.WithValue(r.Context(), correlationIDKey, corID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func FromContext(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(correlationIDKey).(string)
	return id, ok
}
