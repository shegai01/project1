package internal

import (
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
)

func MWRecoverPanic(logger *slog.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					logger.Error("panic recovered", "err", err, "path", r.URL.Path, "method", r.Method)
					initContentType(w)
					http.Error(w, "internal server error", http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
