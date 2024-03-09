package middleware

import (
	"net/http"
	"time"

	"moviepin/utils"
)

func Logger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		handler.ServeHTTP(w, r)
		end := time.Since(start)

		utils.Logger.Printf("[%s] %s %s", r.Method, r.URL.Path, end.String())
	})
}
