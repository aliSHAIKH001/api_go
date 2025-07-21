package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/aliSHAIKH001/api_go/internal/store"
	"github.com/aliSHAIKH001/api_go/internal/tokens"
	"github.com/aliSHAIKH001/api_go/internal/utils"
)

type UserMiddleware struct {
	UserStore store.UserStore
}

// Custom string type to avoid context key collisions
type contextKey string

const userContextKey = contextKey("user")

func SetUser(r *http.Request, user *store.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

func GetUser(r *http.Request) *store.User {
	user, ok := r.Context().Value(userContextKey).(*store.User)
	if !ok {
		// If we set the user, we should get it or panic.
		panic("missing user in context")
	}
	return user
}

// A wrapper for existing http handler functions.
func (um *UserMiddleware) Authenticate(next http.Handler) http.Handler {

	// http.HandlerFunc is an adapter that takes a normal function and turns in into something that
	// implements the http.Handler interface.
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Vary tells caches what request headers affect the response
		w.Header().Add("Vary", "Authorization")
		authHeader := r.Header.Get("Authorization")

		// If the Authorization header is not set, we can assume the user is anonymous.
		if authHeader == "" {
			r = SetUser(r, store.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}

		// If the Authorization header is set
		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "invalid authorization header"})
			return
		}

		token := headerParts[1]
		user, err := um.UserStore.GetUserToken(tokens.ScopeAuth, token)
		if err != nil {
			utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "invalid token"})
			return
		}
		if user == nil {
			utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "token expired or invalid"})
			return
		}

		r = SetUser(r, user)
		next.ServeHTTP(w, r)
		return
	})

}
