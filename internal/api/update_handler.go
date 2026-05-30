package api

import (
	"net/http"
)

func (s *Server) CheckUpdate(w http.ResponseWriter, r *http.Request) {
	result := s.updateChecker.Check()
	writeJSON(w, http.StatusOK, result)
}
