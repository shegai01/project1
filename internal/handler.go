package internal

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
)

type RequestLinks struct {
	Links []string `json:"links"`
}

type StatusResponse struct {
	Links   map[string]string `json:"links"`
	LinksID uint64            `json:"links_id"`
}

type Storage struct {
	Links [][]string
	mu    sync.Mutex
}

type HandlerLinks struct {
	logger  *slog.Logger
	router  *mux.Router
	storage *Storage
}

func NewStorage(links []string) *Storage {
	return &Storage{
		Links: make([][]string, 0),
	}
}
func initContentType(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

func Error(w http.ResponseWriter, statusCode int) {
	http.Error(w, http.StatusText(statusCode), statusCode)
}

func New(logger *slog.Logger, r *mux.Router, storage *Storage) *HandlerLinks {
	return &HandlerLinks{
		logger:  logger,
		router:  r,
		storage: storage,
	}
}

func (s *Storage) Put(links []string) uint {
	s.mu.Lock()
	s.Links = append(s.Links, links)
	s.mu.Unlock()
	return uint(len(s.Links)) - 1

}

func (s *Storage) Get(id uint64) *StatusResponse {
	status := make(map[string]string)
	links := s.Links[int(id)]
	for _, url := range links {
		response, err := http.Head(url)
		if err != nil {
			status[url] = "not available"
			continue
		}

		defer response.Body.Close()

		if response.StatusCode == http.StatusOK {
			slog.Info("available")
			status[url] = "available"
		} else {
			slog.Info("not available")
			status[url] = "not available"
		}
	}

	stat := &StatusResponse{
		Links:   status,
		LinksID: uint64(len(s.Links) - 1),
	}
	slog.Any("id", stat)
	return stat
}

func (h *HandlerLinks) GetLinks(w http.ResponseWriter, r *http.Request) {
	var reqBody RequestLinks
	initContentType(w)
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		slog.Info("son.NewDecoder(r.Body).Decode(&reqBody);")
		return
	}
	id := h.storage.Put(reqBody.Links)
	if len(reqBody.Links) == 0 {
		slog.Info("h.storage.Put(reqBody.Links)")
		return
	}

	status := h.storage.Get(uint64(id))
	slog.Any("status", status)
	if err := json.NewEncoder(w).Encode(status); err != nil {
		slog.Info("json.NewEncoder(w).encode")
		return
	}

}

func (h *HandlerLinks) GetbyID(w http.ResponseWriter, r *http.Request) {
	initContentType(w)
	id := r.URL.Query().Get("id")

	idconv, err := strconv.Atoi(id)
	if err != nil {
		slog.Any("idconv", err)
		return
	}
	respo := h.storage.Get(uint64(idconv))
	if err := json.NewEncoder(w).Encode(respo); err != nil {
		slog.Any("respo", respo)
		return
	}

}
