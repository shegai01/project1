package internal

import (
	"encoding/json"
	"errors"
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
		Links: make([][]string, 0, len(links)),
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

func (s *Storage) Checker(arr []string) *StatusResponse {
	status := make(map[string]string, len(arr))

	for _, url := range arr {
		response, err := http.Head(url)
		if err != nil {
			status[url] = "not available"
			continue
		}

		defer response.Body.Close()

		if response.StatusCode == http.StatusOK {
			status[url] = "available"
		} else {
			status[url] = "not available"
		}
	}

	stat := &StatusResponse{
		Links:   status,
		LinksID: uint64(len(s.Links) - 1),
	}

	return stat
}

func (s *Storage) Put(links []string) uint {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Links = append(s.Links, links)

	return uint(len(s.Links)) - 1

}

func (s *Storage) Get(id uint64) ([]string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if id > uint64(len(s.Links)) {
		slog.Info("incorretly input")
		return nil, errors.New("id < 0 || id > uint64(len(s.Links))")
	}

	links := s.Links[int(id)]

	return links, nil
}

func (h *HandlerLinks) GetStatus(w http.ResponseWriter, r *http.Request) {
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

	sample, err := h.storage.Get(uint64(id))
	if err != nil {
		slog.Info("incorrectly input")
		return
	}

	status := h.storage.Checker(sample)

	slog.Info("status")

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
		slog.Info("idconv")
		return
	}

	respo, err := h.storage.Get(uint64(idconv))
	if err != nil {
		return
	}

	statusResp := h.storage.Checker(respo)
	if err := json.NewEncoder(w).Encode(statusResp); err != nil {
		slog.Info("respo")
		return
	}

	slog.Info("get")
}
