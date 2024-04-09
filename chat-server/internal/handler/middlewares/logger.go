package middlewares

import (
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"time"
)

func MyLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		rw := &responseWriter{WrapResponseWriter: ww}

		defer func() {
			status := rw.Status()
			timeStr := time.Now().Format(time.RFC3339)

			fields := logrus.Fields{
				"MESSAGE": strings.TrimSpace(rw.message),
				"TIME":    timeStr,
				"STATUS":  status,
				"METHOD":  r.Method,
				"URL":     r.URL.String(),
				"LATENCY": time.Since(start),
			}

			if status >= http.StatusInternalServerError {
				logrus.WithFields(fields).Error("Server error occurred.")
			} else if status >= http.StatusBadRequest && status < http.StatusInternalServerError {
				logrus.WithFields(fields).Error("Bad Request.")
			} else {
				logrus.WithFields(fields).Info("Request processed successfully.")
			}
		}()

		next.ServeHTTP(rw, r)
	})
}

type responseWriter struct {
	middleware.WrapResponseWriter
	message string
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	rw.message = string(b)

	return rw.WrapResponseWriter.Write(b)
}
