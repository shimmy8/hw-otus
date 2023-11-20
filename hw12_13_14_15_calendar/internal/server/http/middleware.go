package internalhttp

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

type StatusRecorder struct {
	http.ResponseWriter
	Status int
}

func (r *StatusRecorder) WriteHeader(status int) {
	r.Status = status
	r.ResponseWriter.WriteHeader(status)
}

func loggingMiddleware(h http.Handler, logger Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		recorder := &StatusRecorder{
			ResponseWriter: w,
			Status:         0,
		}
		startTime := time.Now()

		h.ServeHTTP(recorder, r)

		var msgBuilder strings.Builder
		for _, part := range []string{
			r.RemoteAddr,
			r.Method,
			r.URL.String(),
			r.Proto,
			strconv.Itoa(recorder.Status),
			strconv.Itoa(int(r.ContentLength)),
			r.Header.Get("User-Agent"),
			strconv.Itoa(int(time.Since(startTime).Microseconds())),
		} {
			msgBuilder.WriteString(part)
			msgBuilder.WriteString(" ")
		}

		logger.Info(msgBuilder.String())
	})
}
