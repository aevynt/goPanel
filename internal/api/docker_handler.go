package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type DeployComposeRequest struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

func (s *Server) ListContainers(w http.ResponseWriter, r *http.Request) {
	containers, err := s.docker.ListContainers()
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, containers)
}

func (s *Server) StartContainer(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, `{"error":"id is required"}`, http.StatusBadRequest)
		return
	}
	if err := s.docker.StartContainer(id); err != nil {
		http.Error(w, `{"error":"failed to start container"}`, http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "started"})
}

func (s *Server) StopContainer(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, `{"error":"id is required"}`, http.StatusBadRequest)
		return
	}
	if err := s.docker.StopContainer(id); err != nil {
		http.Error(w, `{"error":"failed to stop container"}`, http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "stopped"})
}

func (s *Server) RestartContainer(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, `{"error":"id is required"}`, http.StatusBadRequest)
		return
	}
	if err := s.docker.RestartContainer(id); err != nil {
		http.Error(w, `{"error":"failed to restart container"}`, http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "restarted"})
}

func (s *Server) GetContainerLogs(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, `{"error":"id is required"}`, http.StatusBadRequest)
		return
	}
	logs, err := s.docker.GetContainerLogs(id)
	if err != nil {
		http.Error(w, `{"error":"failed to get logs"}`, http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"logs": logs})
}

func (s *Server) DeployCompose(w http.ResponseWriter, r *http.Request) {
	var req DeployComposeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}
	if req.Name == "" || req.Content == "" {
		http.Error(w, `{"error":"name and content are required"}`, http.StatusBadRequest)
		return
	}

	if err := s.docker.DeployCompose(req.Name, req.Content); err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusAccepted, map[string]string{"status": "deploying", "name": req.Name})
}
