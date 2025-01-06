package middleware

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5/request"
	"github.com/ziliscite/messaging-app/pkg/must"
	"github.com/ziliscite/messaging-app/pkg/token"
	"net/http"
	"os"
)

type contextKey string

const (
	AccessToken  contextKey = "accessToken"
	UserIDKey    contextKey = "userID"
	UserEmailKey contextKey = "userEmail"
)

var ErrInvalidToken = errors.New("invalid token type")

// AuthMiddleware is a middleware that checks if the user is authenticated using JWT bearer token
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bearer, err := request.BearerExtractor{}.ExtractToken(r)
		if err != nil {
			http.Error(w, ErrInvalidToken.Error(), http.StatusForbidden)
			return
		}

		id, email, err := token.Validate(bearer, must.MustEnv(os.Getenv("JWT_SECRET")))
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, id)
		ctx = context.WithValue(ctx, UserEmailKey, email)
		ctx = context.WithValue(ctx, AccessToken, bearer)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
