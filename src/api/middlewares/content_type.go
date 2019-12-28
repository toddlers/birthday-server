package middlewares

import (
	"net/http"

	"github.com/sureshk/birthday-server/src/api/logger"
)

//type MiddlewareFunc func(http.Handler) http.Handler
func AddContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		logger.Logger(w, r)
		next.ServeHTTP(w, r)
	})
}
