package middlewares

import "net/http"

func SetMiddlewareJSON(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Contenty-Type", "application/json")
		next(w, r)
	}
}
