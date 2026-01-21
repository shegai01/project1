package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/shegai01/project1/internal"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelError,
	}))

	r := mux.NewRouter()
	r.Use(internal.MWRecoverPanic)

	h := internal.New(logger, r, &internal.Storage{})

	r.HandleFunc("/getstatus", h.GetStatus)
	r.HandleFunc("/getbyid", h.GetbyID)

	log.Fatal(http.ListenAndServe(":8080", r))

}
