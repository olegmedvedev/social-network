package server

import (
	"context"
	"net/http"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if authHeader := r.Header.Get("Authorization"); authHeader != "" {
			ctx = context.WithValue(ctx, "Authorization", authHeader)
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
