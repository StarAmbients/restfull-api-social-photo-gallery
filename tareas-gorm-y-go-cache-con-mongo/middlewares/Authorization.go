package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/starambients/tareas-gorm-y-go-cache-con-mongo/data"
	"github.com/starambients/tareas-gorm-y-go-cache-con-mongo/lib"
)

var UserKey = "current_user"

func AuthUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL != nil && strings.Index(r.URL.Path, "/users") == 0 {
			next.ServeHTTP(w, r)
		} else {
			// We extract data from the header for later to validate the token
			// For checking if the user us authorized, we specifically extract info about Authorization,
			// which is part of the header structure, and we save this info in a variable
			auto := r.Header.Get("Authorization")

			// If the header "authorization" is not empty and it contains "Bearer" and a space
			// it means it carries a token
			if len(auto) > 0 && strings.Contains(auto, "Bearer ") {
				// We have a token and we need to extract it
				tokenString := strings.Split(auto, " ")[1]
				// We call our Security package to ask it if this token belongs to some user
				// First of all this package has to gives a user
				userValid := lib.GetUserTokenCache(tokenString, data.GetCacheClient())
				// If we got a valid user (whole structure) in return
				if userValid != nil {
					ctx := context.WithValue(r.Context(), UserKey, userValid)
					newReq := r.WithContext(ctx)
					next.ServeHTTP(w, newReq)
				}
			}else {
				// We do not have a token
				w.WriteHeader(http.StatusUnauthorized)
			}
		}
	})
}