package middleware

import (
	"context"
	"net/http"

	usermodels "banner/internal/models/user"
)

const headerTokenName = "token"

type userKeyT string

const UserKey userKeyT = "user key"

func AuthMiddleware(userToken string, adminToken string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var user usermodels.User

		token := r.Header.Get(headerTokenName)
		switch token {
		case userToken:
			user = usermodels.User{IsAdmin: false}
		case adminToken:
			user = usermodels.User{IsAdmin: true}
		default:
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserKey, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
