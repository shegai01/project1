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
	h := internal.NewHandler(r, logger)
	r.HandleFunc("/status", h.Status)
	r.HandleFunc("/getbyid", h.GetbyID)
	log.Fatal(http.ListenAndServe(":8080", r))

}
