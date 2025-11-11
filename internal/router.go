package internal

import (
	"github.com/gorilla/mux"
)

func RegRouter() *mux.Router {
	nr := mux.NewRouter()
	return nr
}
