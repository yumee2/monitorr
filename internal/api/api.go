package api

import (
	"context"
	"encoding/json"
	"monitorr/internal/storage"
	"net/http"
	"strconv"
	"time"
)

type CheckReader interface {
	FindAll(ctx context.Context, interval time.Duration) (map[string][]storage.ServiceStatus, error)
	CountUptime(ctx context.Context, serviceName string, interval time.Duration) (int, error)
	FindTotalByName(ctx context.Context, serviceName string, interval time.Duration) ([]storage.ServiceStatus, error)
}

type Handler struct {
	repo CheckReader
}

func NewHandler(storage CheckReader) *Handler {
	return &Handler{repo: storage}
}

func (s *Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /status", s.handleGetAllServices)
	mux.HandleFunc("GET /status/{name}", s.handleGetStatusByName)
	return withCORS(mux)
}

// withCORS allows the web client (served from a different origin/port in
// dev) to call this read-only, unauthenticated API.
func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (s *Handler) handleGetAllServices(w http.ResponseWriter, r *http.Request) {
	interval := parseInterval(r)
	checks, err := s.repo.FindAll(r.Context(), interval)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(checks)
}

func (s *Handler) handleGetStatusByName(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	interval := parseInterval(r)
	ctx := r.Context()

	uptime, err := s.repo.CountUptime(ctx, name, interval)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	checks, err := s.repo.FindTotalByName(ctx, name, interval)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := struct {
		Name          string                  `json:"name"`
		UptimePercent int                     `json:"uptimePercent"`
		Checks        []storage.ServiceStatus `json:"checks"`
	}{
		Name:          name,
		UptimePercent: uptime,
		Checks:        checks,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// parseInterval reads from the query string, defaulting to 60.
func parseInterval(r *http.Request) time.Duration {
	minutes, err := strconv.Atoi(r.URL.Query().Get("interval"))
	if err != nil || minutes <= 0 {
		minutes = 60
	}
	return time.Duration(minutes) * time.Minute
}
