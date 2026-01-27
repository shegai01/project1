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
		Level: slog.LevelInfo,
	}))

	r := mux.NewRouter()
	r.Use(internal.MWRecoverPanic(logger))

	storage := internal.NewStorage()
	h := internal.New(logger, r, storage)

	r.HandleFunc("/getstatus", h.GetStatus)
	r.HandleFunc("/getbyid", h.GetbyID)
	r.HandleFunc("/getlist", h.GetList)

	log.Fatal(http.ListenAndServe(":8080", r))

}
