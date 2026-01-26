package internal

import (
	"net/http"
)

func MWRecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				initContentType(w)
				http.Error(w, "internal server", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
