package middlewares

import (
	"net/http"
)

//type MiddlewareFunc func(http.Handler) http.Handler
func AddContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
