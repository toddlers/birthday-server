package middlewares

import (
	"net/http"
	"net/http/httputil"

	log "github.com/sirupsen/logrus"
)

func DebugRequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dump, _ := httputil.DumpRequest(r, false)
		log.Debugf("Request: \n%+s", dump)
		next.ServeHTTP(w, r)
	})
}
