package middleware

import (
	"net/http"
	"os"
)

// 認証が必要な部分に認証をさせる
func AuthrizationBearer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var isucariAPIToken = os.Getenv("AUTH_BEARER")
		switch {
		case r.URL.Path == "/accept":
			next.ServeHTTP(w, r)
		case r.Header.Get("Authorization") == isucariAPIToken:
			next.ServeHTTP(w, r)
		default:
			http.Error(w, "authorization error", http.StatusUnauthorized)
		}
	})
}
