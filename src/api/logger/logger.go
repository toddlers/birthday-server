package logger

import (
	"net/http"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

func init() {
	// Log as JSON instead of default ASCII formatter
	logrus.SetFormatter(&logrus.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer
	logrus.SetOutput(os.Stdout)

	// calling method as a field
	//	log.SetReportCaller(true)

	// Only log the info severity or above
	logrus.SetLevel(logrus.InfoLevel)
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		logrus.WithFields(logrus.Fields{
			"method":     r.Method,
			"path":       r.RequestURI,
			"ip":         r.RemoteAddr,
			"duration":   time.Since(start),
			"user_agent": r.UserAgent(),
		}).Info()
		next.ServeHTTP(w, r)
	})
}
