package internal

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"sync/atomic"

	"github.com/gorilla/mux"
)

var counter uint64

type Handler struct {
	logger  *slog.Logger
	router  *mux.Router
	storage *Storage
}
type Links struct {
	Urls []string `json:"urls"`
}

type StatusResponse struct {
	Links   map[string]string `json:"links"`
	LinksID uint64            `json:"links_num"`
}

func NewStatusResponse(status map[string]string) *StatusResponse {
	return &StatusResponse{
		Links:   status,
		LinksID: atomic.AddUint64(&counter, 1),
	}
}

type Storage struct {
	Saved [][]string
}

func initContentType(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

func Error(w http.ResponseWriter, statusCode int) {
	http.Error(w, http.StatusText(statusCode), statusCode)
}

func (h *Handler) Get(id int) map[string]string {

	if id < 0 || len(h.storage.Saved)-1 < id {
		return nil
	}
	res := make(map[string]string)
	urls := h.storage.Saved[id]
	for _, url := range urls {
		response, err := http.Head(url)
		if err != nil {
			res[url] = "not available"
			continue
		}
		defer response.Body.Close()

		if response.StatusCode == http.StatusOK {
			res[url] = "available"
		} else {
			res[url] = "not available"
		}
	}

	return res
}

func NewHandler(router *mux.Router, logger *slog.Logger) *Handler {
	h := &Handler{
		router:  router,
		logger:  logger,
		storage: &Storage{},
	}
	return h
}

func (h *Handler) Status(w http.ResponseWriter, r *http.Request) {
	initContentType(w)

	var reqBody Links

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		h.logger.Error("json.NewDecoder(r.Body)")
		Error(w, http.StatusInternalServerError)
		return
	}

	h.storage.Saved = append(h.storage.Saved, reqBody.Urls)
	id := len(h.storage.Saved) - 1
	resp := map[string]any{
		"message": "links saved",
		"id":      id,
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.logger.Error(" json.NewEncoder(w).Encode(resp)")
	}
}

func (h *Handler) GetbyID(w http.ResponseWriter, r *http.Request) {
	initContentType(w)

	idstr := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idstr)
	if err != nil {
		h.logger.Error(" strconv.Atoi(idstr)")
		Error(w, http.StatusNotFound)
		return
	}

	status := h.Get(id)
	if status == nil {
		Error(w, http.StatusNotFound)
		return
	}

	resp := NewStatusResponse(status)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		h.logger.Error("json.NewEncoder(w)")
		return
	}

}
