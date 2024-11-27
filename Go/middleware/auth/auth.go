//go:build !solution

package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"
)

type User struct {
	Name  string
	Email string
}

func ContextUser(ctx context.Context) (*User, bool) {
	user, ok := ctx.Value(userKey{}).(*User)
	return user, ok
}

var ErrInvalidToken = errors.New("invalid token")

type TokenChecker interface {
	CheckToken(ctx context.Context, token string) (*User, error)
}

func CheckAuth(checker TokenChecker) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			token := parts[1]

			user, err := checker.CheckToken(r.Context(), token)
			if err != nil {
				if errors.Is(err, ErrInvalidToken) {
					http.Error(w, "Unauthorized", http.StatusUnauthorized)
				} else {
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
				return
			}

			ctx := context.WithValue(r.Context(), userKey{}, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

type userKey struct{}
