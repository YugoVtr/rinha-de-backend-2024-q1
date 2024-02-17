package log

import (
	"bytes"
	"log/slog"
	"net/http"
)

type ResponseWriter struct {
	http.ResponseWriter
	status int
	body   *bytes.Buffer
}

func (rw *ResponseWriter) WriteHeader(status int) {
	rw.ResponseWriter.WriteHeader(status)
	rw.status = status
}

func (rw *ResponseWriter) Write(b []byte) (int, error) {
	if rw.status >= 400 || rw.status < 200 {
		rw.body.Write(b)
	}
	return rw.ResponseWriter.Write(b) //nolint: wrapcheck
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		responseWriter := ResponseWriter{ResponseWriter: w, body: &bytes.Buffer{}}
		next.ServeHTTP(&responseWriter, request)
		slog.Info(
			"request",
			"method", request.Method,
			"path", request.URL.Path,
			"status", responseWriter.status,
			"body", responseWriter.body.String(),
		)
	})
}
