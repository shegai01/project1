package internal

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type Handler struct {
	// logger *slog.Logger
	mx *mux.Router
}

type Result struct {
	URL     string `json:"url"`
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

type StatusReq struct {
	Results []Result `json:"results"`
}

func initContentType(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

func Error(w http.ResponseWriter, statusCode int) {
	http.Error(w, http.StatusText(statusCode), statusCode)
}

func NewHandler(router *mux.Router) *Handler {
	h := &Handler{}
	return h
}

func (h *Handler) Status(w http.ResponseWriter, r *http.Request) {
	initContentType(w)
	var req StatusReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "json.NewDecoder(r.Body", http.StatusInternalServerError)
		return
	}

	results := make([]Result, 0, len(req.Results))
	for _, url := range req.Results {
		results = append(results, CheckUrl(url.URL))
	}
	json.NewEncoder(w).Encode(StatusReq{Results: results})
}
