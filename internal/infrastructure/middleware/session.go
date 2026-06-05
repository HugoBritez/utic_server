package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"sync"
)

type sessionContextKey string

const sessionKey sessionContextKey = "admin_user"

type SessionStore struct {
	mu    sync.RWMutex
	tokens map[string]string // token -> username
}

func NewSessionStore() *SessionStore {
	return &SessionStore{ tokens: make(map[string]string) }
}

func (s *SessionStore) Create(username string) string {
	b := make([]byte, 32)
	rand.Read(b)
	token := hex.EncodeToString(b)
	s.mu.Lock()
	s.tokens[token] = username
	s.mu.Unlock()
	return token
}

func (s *SessionStore) Get(token string) (string, bool) {
	s.mu.RLock()
	user, ok := s.tokens[token]
	s.mu.RUnlock()
	return user, ok
}

func (s *SessionStore) Delete(token string) {
	s.mu.Lock()
	delete(s.tokens, token)
	s.mu.Unlock()
}

func AdminSession(sessions *SessionStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/admin/login" {
				next.ServeHTTP(w, r)
				return
			}

			c, err := r.Cookie("admin_session")
			if err != nil {
				http.Redirect(w, r, "/admin/login", http.StatusFound)
				return
			}

			user, ok := sessions.Get(c.Value)
			if !ok {
				http.Redirect(w, r, "/admin/login", http.StatusFound)
				return
			}

			ctx := context.WithValue(r.Context(), sessionKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
