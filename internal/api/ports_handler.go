package api

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/aevynt/goPanel/internal/ports"
)

func (s *Server) ListPorts(w http.ResponseWriter, r *http.Request) {
	startStr := r.URL.Query().Get("start")
	endStr := r.URL.Query().Get("end")

	var result []ports.PortInfo
	var err error

	if startStr != "" && endStr != "" {
		start, _ := strconv.Atoi(startStr)
		end, _ := strconv.Atoi(endStr)
		result, err = s.pm.ScanRange(start, end)
	} else {
		result, err = s.pm.ListListening()
	}
	if err != nil {
		http.Error(w, `{"error":"failed to scan ports"}`, http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (s *Server) CheckPort(w http.ResponseWriter, r *http.Request) {
	portStr := chi.URLParam(r, "port")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		http.Error(w, `{"error":"invalid port"}`, http.StatusBadRequest)
		return
	}
	available := s.pm.IsPortAvailable(port)
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"port":      port,
		"available": available,
	})
}

func (s *Server) FindPort(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Preferred int `json:"preferred"`
	}
	if err := decodeJSON(r, &req); err != nil {
		req.Preferred = 0
	}
	port, err := s.pm.FindAvailablePort(req.Preferred)
	if err != nil {
		http.Error(w, `{"error":"no available ports"}`, http.StatusServiceUnavailable)
		return
	}
	writeJSON(w, http.StatusOK, map[string]int{"port": port})
}
