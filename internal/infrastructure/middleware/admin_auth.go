package middleware

import (
	"crypto/sha256"
	"crypto/subtle"
	"net/http"
	"os"
)

func AdminAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok {
			w.Header().Set("WWW-Authenticate", `Basic realm="admin"`)
			http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
			return
		}

		expectedUser := os.Getenv("ADMIN_USER")
		expectedPass := os.Getenv("ADMIN_PASSWORD")

		userHash := sha256.Sum256([]byte(user))
		passHash := sha256.Sum256([]byte(pass))
		expectedUserHash := sha256.Sum256([]byte(expectedUser))
		expectedPassHash := sha256.Sum256([]byte(expectedPass))

		if subtle.ConstantTimeCompare(userHash[:], expectedUserHash[:]) != 1 ||
			subtle.ConstantTimeCompare(passHash[:], expectedPassHash[:]) != 1 {
			w.Header().Set("WWW-Authenticate", `Basic realm="admin"`)
			http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
