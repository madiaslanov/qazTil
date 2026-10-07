package httpapi

import (
	"context"
	"net/http"
	"strings"

	"github.com/madiaslanov/qazTil/internal/domain"
)

type ctxKey int

const userIDKey ctxKey = iota

// RequireUser rejects a request without a valid Bearer token and puts the
// learner id into the request context for the handler behind it.
func (a *API) RequireUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw, ok := bearerToken(r)
		if !ok {
			writeErr(w, domain.ErrUnauthorized)
			return
		}
		userID, err := a.tokens.Parse(raw)
		if err != nil {
			writeErr(w, domain.ErrUnauthorized)
			return
		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), userIDKey, userID)))
	})
}

// UserID returns the learner id that RequireUser put into the context.
func UserID(ctx context.Context) (int64, bool) {
	userID, ok := ctx.Value(userIDKey).(int64)
	return userID, ok
}

func bearerToken(r *http.Request) (string, bool) {
	const prefix = "bearer "
	header := strings.TrimSpace(r.Header.Get("Authorization"))
	if len(header) <= len(prefix) || !strings.EqualFold(header[:len(prefix)], prefix) {
		return "", false
	}
	token := strings.TrimSpace(header[len(prefix):])
	return token, token != ""
}
