package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func MWRecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				initContentType(w)
				response := fmt.Sprintln("incorrect input:id")
				json.NewEncoder(w).Encode(response)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
