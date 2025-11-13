package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/Numeez/go-zenith/internal/store"
	"github.com/Numeez/go-zenith/internal/tokens"
	"github.com/Numeez/go-zenith/internal/utils"
)

type contextKey string

const userContextKey = contextKey("userKey")

type UserMiddleware struct {
	UserStore store.UserStore
}

func SetUser(r *http.Request, user *store.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

func GetUser(r *http.Request) *store.User {
	user, ok := r.Context().Value(userContextKey).(*store.User)
	if !ok {
		//very aggresive can be thought through
		panic("missing user in request")
	}
	return user
}

func (um *UserMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			r = SetUser(r, store.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}
		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			_ = utils.WriteJson(w, http.StatusBadRequest, utils.Envelope{"error": "invalid authorization header "})
			return
		}
		token := headerParts[1]
		user, err := um.UserStore.GetUserToken(tokens.ScopeAuth, token)
		if err != nil {
			_ = utils.WriteJson(w, http.StatusUnauthorized, utils.Envelope{"error": "invalid authorization token"})
			return

		}
		if user == nil {
			_ = utils.WriteJson(w, http.StatusUnauthorized, utils.Envelope{"error": "token expired"})
			return
		}
		r = SetUser(r, user)
		next.ServeHTTP(w, r)

	})

}

func (um *UserMiddleware) RequireUser(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := GetUser(r)
		if user.IsAnonymous() {
			_ = utils.WriteJson(w, http.StatusUnauthorized, utils.Envelope{"error": "you must be logged into access this"})
			return
		}
		next.ServeHTTP(w, r)
	})
}
