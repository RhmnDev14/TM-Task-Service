package middleware

import (
	"net/http"
	"strings"

	"task-service/internal/infrastructure"
)

type Middleware struct {
	jwt infrastructure.JWTService
}

func NewMiddleware(jwt infrastructure.JWTService) *Middleware {
	return &Middleware{jwt: jwt}
}

func (m *Middleware) PrivateMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "authorization header is required", http.StatusUnauthorized)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			http.Error(w, "authorization header format must be Bearer {token}", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]

		ctx := r.Context()
		newCtx, _, err := m.jwt.ValidateToken(ctx, tokenString)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r.WithContext(newCtx))
	})
}
