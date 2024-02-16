package middleware

import (
	"net/http"
	"time"

	"moviepin/utils"
)

type spyWriter struct {
	Writer     *http.ResponseWriter
	statusCode int
}

func (s *spyWriter) Header() http.Header {
	return (*s.Writer).Header()
}

func (s *spyWriter) Write(b []byte) (int, error) {
	return (*s.Writer).Write(b)
}

func (s *spyWriter) WriteHeader(statusCode int) {
	s.statusCode = statusCode

	(*s.Writer).WriteHeader(statusCode)
}

func Logger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		spyWriter := &spyWriter{Writer: &w}

		start := time.Now()
		handler.ServeHTTP(spyWriter, r)
		end := time.Since(start)

		utils.Logger.Printf("[%s] %s %s %d - %s", r.Method, r.URL.Path, end.String(), spyWriter.statusCode, http.StatusText(spyWriter.statusCode))
	})
}
