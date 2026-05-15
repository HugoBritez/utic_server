package middleware

import (
	"context"
	"net/http"
	"os"
)

type contextKey string

const APIKeyContextKey contextKey = "api_key"

// APIKey validates the X-API-Key header against the API_KEY env var.
func APIKey() func(http.Handler) http.Handler {
	expectedKey := os.Getenv("API_KEY")

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := r.Header.Get("X-API-Key")
			if key == "" {
				http.Error(w, `{"error":"missing X-API-Key header"}`, http.StatusUnauthorized)
				return
			}

			if key != expectedKey {
				http.Error(w, `{"error":"invalid API key"}`, http.StatusForbidden)
				return
			}

			ctx := context.WithValue(r.Context(), APIKeyContextKey, key)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetAPIKeyFromContext retrieves the validated API key from the request context.
func GetAPIKeyFromContext(ctx context.Context) (string, bool) {
	key, ok := ctx.Value(APIKeyContextKey).(string)
	return key, ok
}
