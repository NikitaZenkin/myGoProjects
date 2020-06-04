package middleware

import (
	"context"
	"net/http"
	"rclone/pkg/session"
)

func Auth(sm *session.SessionsManager, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sess, err := sm.CheckSession(r)
		if err == nil {
			ctx := context.WithValue(r.Context(), "sessionKey", sess)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
