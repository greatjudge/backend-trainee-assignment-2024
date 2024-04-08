package middleware

import (
	"context"
	"net/http"

	usermodels "banner/internal/models/user"
)

type userKeyT string

const UserKey userKeyT = "user key"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := usermodels.User{IsAdmin: true}

		ctx := context.WithValue(r.Context(), UserKey, user) // TODO change, add logic

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
