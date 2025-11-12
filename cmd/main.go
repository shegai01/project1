package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/shegai01/project1/internal"
)

func main() {
	r := mux.NewRouter()
	h := internal.NewHandler(r)
	r.HandleFunc("/status", h.Status)
	r.HandleFunc("/getbyid", h.GetbyID)
	log.Fatal(http.ListenAndServe(":8080", r))

}
