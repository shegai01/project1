package internal

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type Handler struct {
	router *mux.Router
}
type Links struct {
	Urls []string `json:"urls"`
}

type StatusResponse struct {
	Links   map[string]string `json:"links"`
	LinksID uint64            `json:"links_num"`
}

type Storage struct {
	Array []StatusResponse
}

func initContentType(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

func Error(w http.ResponseWriter, statusCode int) {
	http.Error(w, http.StatusText(statusCode), statusCode)
}

func NewHandler(router *mux.Router) *Handler {
	h := &Handler{
		router: router,
	}
	return h
}

func (h *Handler) Status(w http.ResponseWriter, r *http.Request) {
	initContentType(w)

	var reqBody Links

	results := make(map[string]string)
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	var responseBody StatusResponse
	for _, val := range reqBody.Urls {
		response, err := http.Head(val)
		if err != nil {
			results[val] = "not available"
			continue
		}
		response.Body.Close()

		if response.StatusCode == http.StatusOK {
			results[val] = "available"
		} else {
			results[val] = "not available"
		}

	}

	responseBody = StatusResponse{
		Links:   results,
		LinksID: uint64(len(reqBody.Urls)),
	}

	var storage Storage
	storage.Array = append(storage.Array, responseBody)

	if err := json.NewEncoder(w).Encode(responseBody); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
