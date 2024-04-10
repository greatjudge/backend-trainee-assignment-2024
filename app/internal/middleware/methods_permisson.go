package middleware

import (
	"banner/internal/constants"
	usermodels "banner/internal/models/user"
	"banner/internal/sending"
	"net/http"
)

func OnlyAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value(UserKey).(usermodels.User)
		if !ok {
			sending.SendErrorMsg(w, http.StatusInternalServerError, constants.ErrMsgUserNotFoundInCTX)
			return
		}

		if !user.IsAdmin {
			w.WriteHeader(http.StatusForbidden)
		}

		next.ServeHTTP(w, r)
	})
}
