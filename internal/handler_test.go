package internal

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gorilla/mux"
)

func newTestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func TestStoragePutAndGet(t *testing.T) {
	storage := NewStorage()

	links := []string{"http://example.com", "http://example.org"}
	id := storage.Put(links)

	if id != 0 {
		t.Fatalf("expected id 0, got %d", id)
	}

	got, err := storage.Get(uint64(id))
	if err != nil {
		t.Fatalf("unexpected error from Get: %v", err)
	}

	if len(got) != len(links) {
		t.Fatalf("expected %d links, got %d", len(links), len(got))
	}

	for i, l := range links {
		if got[i] != l {
			t.Errorf("expected link %q at index %d, got %q", l, i, got[i])
		}
	}
}

func TestGetStatus_Success(t *testing.T) {
	// HTTP-сервер, который всегда отвечает 200 OK,
	// чтобы Checker помечал ссылку как "available".
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	logger := newTestLogger()
	router := mux.NewRouter()
	storage := NewStorage()
	handler := New(logger, router, storage)

	bodyStruct := RequestLinks{Links: []string{ts.URL}}
	bodyBytes, err := json.Marshal(bodyStruct)
	if err != nil {
		t.Fatalf("failed to marshal request body: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/getstatus", bytes.NewReader(bodyBytes))
	rr := httptest.NewRecorder()

	handler.GetStatus(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d, body: %s", http.StatusOK, rr.Code, rr.Body.String())
	}

	if ct := rr.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %q", ct)
	}

	var resp map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	status, ok := resp[ts.URL]
	if !ok {
		t.Fatalf("expected key %q in response, got %v", ts.URL, resp)
	}

	if status != "available" {
		t.Errorf("expected status 'available' for %q, got %q", ts.URL, status)
	}
}

func TestGetStatus_BadJSON(t *testing.T) {
	logger := newTestLogger()
	router := mux.NewRouter()
	storage := NewStorage()
	handler := New(logger, router, storage)

	req := httptest.NewRequest(http.MethodPost, "/getstatus", bytes.NewBufferString("{invalid json"))
	rr := httptest.NewRecorder()

	handler.GetStatus(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestGetbyID_Success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	logger := newTestLogger()
	router := mux.NewRouter()
	storage := NewStorage()
	handler := New(logger, router, storage)

	// заранее кладём ссылки в сторадж и используем полученный id
	id := storage.Put([]string{ts.URL})

	req := httptest.NewRequest(http.MethodGet, "/getbyid?id="+strconv.FormatUint(uint64(id), 10), nil)
	rr := httptest.NewRecorder()

	handler.GetbyID(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d, body: %s", http.StatusOK, rr.Code, rr.Body.String())
	}

	if ct := rr.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %q", ct)
	}

	var resp map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	status, ok := resp[ts.URL]
	if !ok {
		t.Fatalf("expected key %q in response, got %v", ts.URL, resp)
	}

	if status != "available" {
		t.Errorf("expected status 'available' for %q, got %q", ts.URL, status)
	}
}

func TestGetList_Success(t *testing.T) {
	logger := newTestLogger()
	router := mux.NewRouter()
	storage := NewStorage()
	handler := New(logger, router, storage)

	first := []string{"http://one.example"}
	second := []string{"http://two.example", "http://three.example"}

	storage.Put(first)
	storage.Put(second)

	req := httptest.NewRequest(http.MethodGet, "/getlist", nil)
	rr := httptest.NewRecorder()

	handler.GetList(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d, body: %s", http.StatusOK, rr.Code, rr.Body.String())
	}

	if ct := rr.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %q", ct)
	}

	var resp [][]string
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if len(resp) != 2 {
		t.Fatalf("expected 2 items in list, got %d", len(resp))
	}

	if len(resp[0]) != len(first) || resp[0][0] != first[0] {
		t.Errorf("unexpected first element: %#v", resp[0])
	}
	if len(resp[1]) != len(second) {
		t.Errorf("unexpected second element length: got %d, want %d", len(resp[1]), len(second))
	}
}

