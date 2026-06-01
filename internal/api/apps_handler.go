package api

import (
	"net/http"
)

func (s *Server) ListApps(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, s.apps.ListCatalog())
}

func (s *Server) DeployApp(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Key string `json:"key"`
	}
	if err := decodeJSON(r, &req); err != nil {
		http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
		return
	}

	port, err := s.apps.DeployApp(req.Key)
	if err != nil {
		// If port is non-zero, it deployed but Caddy mapping failed
		if port != 0 {
			writeJSON(w, http.StatusOK, map[string]interface{}{
				"status":  "deployed",
				"port":    port,
				"warning": err.Error(),
			})
			return
		}
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"status": "deployed",
		"port":   port,
	})
}
