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

func NewStorage() *Storage {
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

func (s *Storage) Checker(arr []string) *map[string]string {
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

	return &status
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
		slog.Error("id > uint64(len(s.Links))", "incorretly input", id)
		return nil, errors.New("id < 0 || id > uint64(len(s.Links))")
	}

	links := s.Links[int(id)]

	return links, nil
}

func (h *HandlerLinks) GetStatus(w http.ResponseWriter, r *http.Request) {
	var reqBody RequestLinks
	initContentType(w)

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		Error(w, http.StatusBadRequest)
		slog.Error("json.NewDecoder(r.Body).Decode(&reqBody);")
		return
	}

	id := h.storage.Put(reqBody.Links)
	if len(reqBody.Links) == 0 {
		slog.Error("h.storage.Put(reqBody.Links)", "id", id)
		return
	}

	sample, err := h.storage.Get(uint64(id))
	if err != nil {
		slog.Error("h.storage.Get", "err", err)
		return
	}

	status := h.storage.Checker(sample)

	slog.Info("status")

	if err := json.NewEncoder(w).Encode(status); err != nil {
		slog.Error("json.NewEncoder(w).encode", "err", err)
		return
	}

}

func (h *HandlerLinks) GetbyID(w http.ResponseWriter, r *http.Request) {
	initContentType(w)

	id := r.URL.Query().Get("id")

	// idconv, err := strconv.Atoi(id)
	idconv, err := strconv.ParseUint(id, 0, 10)
	if err != nil {
		slog.Error("strconv.Atoi(id)", "err", err)
		return
	}

	respo, err := h.storage.Get(uint64(idconv))
	if err != nil {
		slog.Error("h.storage.Get(uint64(idconv))", "err", err)
		return
	}

	statusResp := h.storage.Checker(respo)
	if err := json.NewEncoder(w).Encode(statusResp); err != nil {
		slog.Error("json.NewEncoder(w)", "err", err)
		return
	}

	slog.Info("get")
}

func (h *HandlerLinks) GetList(w http.ResponseWriter, r *http.Request) {
	initContentType(w)

}
