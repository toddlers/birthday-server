package logger

import (
	"net/http"
	"strconv"
	"os"
	"strings"
	"time"
	"sync"
	"github.com/felixge/httpsnoop"
	"github.com/sirupsen/logrus"
	"https://github.com/kjk/siser"
)

var mmuLogHTTP sync.Mutex

// LogReqInfo describes info about HTTP request
type HTTPReqInfo struct {
	// GET etc.
	method  string
	uri     string
	referer string
	ipaddr  string
	// response code, like 200, 404
	code int
	// number of bytes of the response sent
	size int64
	// how long did it take to
	duration  time.Duration
	userAgent string
}

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

// Request.RemoteAddress contains port, which we want to remove i.e.:
//"[::1]:1234" => "[::1]"
func ipAddrFromRemoteAddr(s string) string {
	idx := strings.LastIndex(s, ":")
	if idx == -1 {
		return s
	}
	return s[:idx]
}

// requestGetRemoteAddress returns ip address of the client making the request,
// taking into account http proxies

func requestGetRemoteAddress(r *http.Request) string {

	hdr := r.Header
	hdrRealIP := hdr.Get("X-Real-Ip")
	hdrForwardedFor := hdr.Get("X-Forwarded-For")
	if hdrRealIP == "" && hdrForwardedFor == "" {
		return ipAddrFromRemoteAddr(r.RemoteAddr)
	}
	if hdrForwardedFor != "" {
		// X-Forwarded-For is potentially a list of addresses separated with ","
		parts := strings.Split(hdrForwardedFor, ",")
		for i, p := range parts {
			parts[i] = strings.TrimSpace(p)
		}
		return parts[0]
	}
	return hdrRealIP
}

func logHTTPReq(ri *HTTPReqInfo){
	var rec siser.Record
	rec.Name = "httplog"
	rec.Append("method", ri.method)
	rec.Append("uri", ri.uri)
	if ri.referer != ""{
		rec.Append("referer", ri.referer)
	}
	rec.Append("ipaddr", ri.ipaddr)
	rec.Append("code", strconv.Itoa(ri.code))
	rec.Append("size", strconv.FormatInt(ri.size, 10))
	durMs := ri.duration / time.Millisecond
	rec.Append("duration", strconv.FormatInt(int64(durMs),10))
	mmuLogHTTP.Lock()
	defer mmuLogHTTP.Unlock()
	_, _ = httpLogSiser.WriteRecord(&rec)

}

func Logger(h http.Handler) http.Handler {

	fn := func(w http.ResponseWriter, r *http.Request) {
		ri := &HTTPReqInfo{
			method:    r.Method,
			uri:       r.URL.String(),
			referer:   r.Header.Get("Referer"),
			userAgent: r.Header.Get("User-Agent"),
		}
		ri.ipaddr = requestGetRemoteAddress(r)

		// this runs handler h and captures information about
		// HTTP Request
		m := httpsnoop.CaptureMetrics(h,w,r)
		ri.code = m.Code
		ri.size = m.BytesWritten
		ri.duration = m.Duration
		logHTTPReq(ri)
	}
	return http.HandlerFunc(fn)
}
// func Logger(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		start := time.Now()
// 		logrus.WithFields(logrus.Fields{
// 			"method":     r.Method,
// 			"path":       r.RequestURI,
// 			"ip":         r.RemoteAddr,
// 			"duration":   time.Since(start),
// 			"user_agent": r.UserAgent(),
// 		}).Info()
// 		next.ServeHTTP(w, r)
// 	})
// }
