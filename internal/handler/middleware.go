package handler

import (
	"hw3/pkg/token"
	"net/http"
	"strings"
)

func AuthCheck(secret string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth == "" {
			writeError(w, "Token is none", http.StatusUnauthorized)
			return
		}
		t, ok := strings.CutPrefix(auth, "Bearer ")
		if !ok {
			writeError(w, "Token is none", http.StatusUnauthorized)
			return
		}
		_, err := token.ParseToken(t, secret)
		if err != nil {

			writeError(w, "wrong token", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)

	})
}
