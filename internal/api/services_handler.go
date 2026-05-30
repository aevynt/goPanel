package api

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/lhqua/gopanel/internal/servicemanager"
)

func isSafeServiceName(name string) bool {
	if name == "" || len(name) > 64 || name == "." || name == ".." {
		return false
	}
	for _, char := range name {
		if !((char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') || char == '-' || char == '_') {
			return false
		}
	}
	return true
}

func (s *Server) ListServices(w http.ResponseWriter, r *http.Request) {
	services, err := s.sm.List()
	if err != nil {
		http.Error(w, `{"error":"failed to list services"}`, http.StatusInternalServerError)
		return
	}
	panelNames, err := s.db.Query("SELECT name FROM panel_services")
	if err == nil {
		defer panelNames.Close()
		for panelNames.Next() {
			var name string
			if err := panelNames.Scan(&name); err != nil {
				continue
			}
			for i := range services {
				if services[i].Name == name {
					services[i].PanelManaged = true
					break
				}
			}
		}
	}

	hasGoPanel := false
	for _, svc := range services {
		if svc.Name == "goPanel" {
			hasGoPanel = true
			break
		}
	}
	if !hasGoPanel {
		exePath := ""
		if p, err := os.Executable(); err == nil {
			exePath = p
		}
		services = append(services, servicemanager.Service{
			Name:         "goPanel",
			Description:  "goPanel server management panel",
			Status:       servicemanager.StatusActive,
			Port:         s.cfg.Port,
			BinaryPath:   exePath,
			PanelManaged: true,
		})
	}

	writeJSON(w, http.StatusOK, services)
}

func (s *Server) GetService(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	if !isSafeServiceName(name) {
		http.Error(w, `{"error":"invalid service name"}`, http.StatusBadRequest)
		return
	}
	svc, err := s.sm.Get(name)
	if err != nil {
		http.Error(w, `{"error":"service not found"}`, http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, svc)
}

func (s *Server) CreateService(w http.ResponseWriter, r *http.Request) {
	var spec servicemanager.ServiceSpec
	if err := json.NewDecoder(r.Body).Decode(&spec); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}
	if spec.Name == "" || spec.BinaryPath == "" {
		http.Error(w, `{"error":"name and binary_path are required"}`, http.StatusBadRequest)
		return
	}
	for _, char := range spec.Name {
		if !((char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') || char == '-' || char == '_') {
			http.Error(w, `{"error":"invalid service name: only alphanumeric, dashes, and underscores allowed"}`, http.StatusBadRequest)
			return
		}
	}
	if err := s.sm.Create(spec); err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}
	s.db.Exec("INSERT OR IGNORE INTO panel_services (name) VALUES (?)", spec.Name)
	writeJSON(w, http.StatusCreated, map[string]string{"status": "created", "name": spec.Name})
}

func (s *Server) StartService(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	if !isSafeServiceName(name) {
		http.Error(w, `{"error":"invalid service name"}`, http.StatusBadRequest)
		return
	}
	if err := s.sm.Start(name); err != nil {
		http.Error(w, `{"error":"failed to start service"}`, http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "started"})
}

func (s *Server) StopService(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	if !isSafeServiceName(name) {
		http.Error(w, `{"error":"invalid service name"}`, http.StatusBadRequest)
		return
	}
	if err := s.sm.Stop(name); err != nil {
		http.Error(w, `{"error":"failed to stop service"}`, http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "stopped"})
}

func (s *Server) RestartService(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	if !isSafeServiceName(name) {
		http.Error(w, `{"error":"invalid service name"}`, http.StatusBadRequest)
		return
	}
	if err := s.sm.Restart(name); err != nil {
		http.Error(w, `{"error":"failed to restart service"}`, http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "restarted"})
}

func (s *Server) EnableService(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	if !isSafeServiceName(name) {
		http.Error(w, `{"error":"invalid service name"}`, http.StatusBadRequest)
		return
	}
	if err := s.sm.Enable(name); err != nil {
		http.Error(w, `{"error":"failed to enable service"}`, http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "enabled"})
}

func (s *Server) DisableService(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	if !isSafeServiceName(name) {
		http.Error(w, `{"error":"invalid service name"}`, http.StatusBadRequest)
		return
	}
	if err := s.sm.Disable(name); err != nil {
		http.Error(w, `{"error":"failed to disable service"}`, http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "disabled"})
}

func (s *Server) RemoveService(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	if !isSafeServiceName(name) {
		http.Error(w, `{"error":"invalid service name"}`, http.StatusBadRequest)
		return
	}
	if err := s.sm.Remove(name); err != nil {
		http.Error(w, `{"error":"failed to remove service"}`, http.StatusInternalServerError)
		return
	}
	s.db.Exec("DELETE FROM panel_services WHERE name = ?", name)
	writeJSON(w, http.StatusOK, map[string]string{"status": "removed"})
}

func (s *Server) GetServiceLogs(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	if !isSafeServiceName(name) {
		http.Error(w, `{"error":"invalid service name"}`, http.StatusBadRequest)
		return
	}
	tailStr := r.URL.Query().Get("tail")
	tail := 50
	if tailStr != "" {
		if t, err := strconv.Atoi(tailStr); err == nil && t > 0 {
			tail = t
		}
	}
	logs, err := s.sm.Logs(name, tail)
	if err != nil {
		http.Error(w, `{"error":"failed to get logs"}`, http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, logs)
}
