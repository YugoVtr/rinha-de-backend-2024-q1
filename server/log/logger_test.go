package log_test

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yugovtr/rinha-de-backend-2024-q1/server/log"
)

func HelperSetDefaultLogger(t *testing.T, writer io.Writer) {
	t.Helper()
	handlerOptions := &slog.HandlerOptions{
		ReplaceAttr: func(_ []string, a slog.Attr) (s slog.Attr) {
			if a.Key == slog.TimeKey {
				return s
			}
			return a
		},
	}
	defer slog.SetDefault(slog.Default())
	slog.SetDefault(slog.New(slog.NewTextHandler(writer, handlerOptions)))
}

func TestLoggingMiddleware(t *testing.T) {
	testCases := []struct {
		statusCode  int
		body        string
		expectedLog string
	}{
		{
			statusCode:  http.StatusOK,
			body:        "test",
			expectedLog: "level=INFO msg=request method=GET path=/test status=200 body=\"\"\n",
		},
		{
			statusCode:  http.StatusNotFound,
			body:        "not found",
			expectedLog: "level=INFO msg=request method=GET path=/test status=404 body=\"not found\"\n",
		},
	}
	for _, testCases := range testCases {
		t.Run(fmt.Sprint(testCases.statusCode), func(t *testing.T) {
			writer := &bytes.Buffer{}
			HelperSetDefaultLogger(t, writer)
			handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(testCases.statusCode)
				w.Write([]byte(testCases.body))
			})
			req, err := http.NewRequest("GET", "/test", nil) //nolint: noctx
			if err != nil {
				t.Fatal(err)
			}
			rr := httptest.NewRecorder()
			loggingMiddleware := log.LoggingMiddleware(handler)
			loggingMiddleware.ServeHTTP(rr, req)
			assert.Equal(t, testCases.expectedLog, writer.String())
		})
	}
}
